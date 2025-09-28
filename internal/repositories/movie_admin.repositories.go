package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/utils"
	"github.com/redis/go-redis/v9"
)

type MovieAdminRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewMovieAdminRepository(db *pgxpool.Pool, rdb *redis.Client) *MovieAdminRepository {
	return &MovieAdminRepository{db: db, rdb: rdb}
}

// ADD: helper invalidasi cache
func (r *MovieAdminRepository) InvalidateMovieCache(ctx context.Context) { // FIXED
	// hapus upcoming
	r.rdb.Del(ctx, "lala:movie-upcoming")

	// hapus semua filter page=1
	keys, _ := r.rdb.Keys(ctx, "movies:filter:*:page:1").Result()
	for _, key := range keys {
		r.rdb.Del(ctx, key)
	}

	// Popular tidak dihapus manual → biarkan expired otomatis setelah 7 hari
}

func (r *MovieAdminRepository) GetAllMovies(ctx context.Context) ([]models.MovieAdmin, error) {
	sql := `
		SELECT id, title, director_id, poster, background_poster,
			   releasedate, duration, synopsis, popularity,
			   created_at, updated_at
		FROM movies
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieAdmin
	for rows.Next() {
		var m models.MovieAdmin
		var iv pgtype.Interval

		if err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.DirectorID,
			&m.Poster,
			&m.BackgroundPoster,
			&m.ReleaseDate,
			&iv,
			&m.Synopsis,
			&m.Popularity,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}

		// simpan interval mentah + string hasil convert
		m.Duration = iv
		m.DurationText = utils.FormatIntervalToText(iv)

		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieAdminRepository) DeleteMovie(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `
		UPDATE movies
		SET deleted_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND deleted_at IS NULL
	`, id)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return errors.New("movie not found or already deleted")
	}

	// ADD: invalidasi cache
	r.InvalidateMovieCache(ctx)

	return nil
}

func (r *MovieAdminRepository) UpdateMovie(ctx context.Context, id int, updates map[string]interface{}, genres []int, casts []int) (models.MovieAdmin, error) {
	if len(updates) == 0 && len(genres) == 0 && len(casts) == 0 {
		return models.MovieAdmin{}, fmt.Errorf("no fields to update")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.MovieAdmin{}, fmt.Errorf("begin tx failed: %w", err)
	}
	defer tx.Rollback(ctx) // rollback kalau gagal

	var m models.MovieAdmin
	var iv pgtype.Interval

	// --- Update movie fields ---
	if len(updates) > 0 {
		setParts := []string{}
		args := []interface{}{}
		i := 1
		for k, v := range updates {
			setParts = append(setParts, fmt.Sprintf("%s = $%d", k, i))
			args = append(args, v)
			i++
		}
		// always update updated_at
		setParts = append(setParts, "updated_at = NOW()")

		sql := fmt.Sprintf(`
			UPDATE movies
			SET %s
			WHERE id = $%d
			RETURNING id, title, director_id, poster, background_poster,
					  releaseDate, duration, synopsis, popularity,
					  created_at, updated_at
		`, strings.Join(setParts, ", "), i)

		args = append(args, id)

		err = tx.QueryRow(ctx, sql, args...).Scan(
			&m.ID,
			&m.Title,
			&m.DirectorID,
			&m.Poster,
			&m.BackgroundPoster,
			&m.ReleaseDate,
			&iv,
			&m.Synopsis,
			&m.Popularity,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return models.MovieAdmin{}, fmt.Errorf("update failed: %w", err)
		}
	} else {
		// kalau tidak ada update field, ambil data existing
		err = tx.QueryRow(ctx, `
			SELECT id, title, director_id, poster, background_poster,
			       releaseDate, duration, synopsis, popularity,
			       created_at, updated_at
			FROM movies WHERE id=$1
		`, id).Scan(
			&m.ID,
			&m.Title,
			&m.DirectorID,
			&m.Poster,
			&m.BackgroundPoster,
			&m.ReleaseDate,
			&iv,
			&m.Synopsis,
			&m.Popularity,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return models.MovieAdmin{}, fmt.Errorf("fetch movie failed: %w", err)
		}
	}

	m.Duration = iv
	m.DurationText = utils.FormatIntervalToText(iv)

	// --- Replace genres ---
	if len(genres) > 0 {
		_, err := tx.Exec(ctx, "DELETE FROM movie_genre WHERE movie_id=$1", id)
		if err != nil {
			return models.MovieAdmin{}, fmt.Errorf("delete genres failed: %w", err)
		}
		for _, gid := range genres {
			_, err := tx.Exec(ctx, "INSERT INTO movie_genre (movie_id, genre_id) VALUES ($1,$2)", id, gid)
			if err != nil {
				return models.MovieAdmin{}, fmt.Errorf("insert genres failed: %w", err)
			}
		}
		m.Genres = genres
	}

	// --- Replace casts ---
	if len(casts) > 0 {
		_, err := tx.Exec(ctx, "DELETE FROM movie_cast WHERE movie_id=$1", id)
		if err != nil {
			return models.MovieAdmin{}, fmt.Errorf("delete casts failed: %w", err)
		}
		for _, cid := range casts {
			_, err := tx.Exec(ctx, "INSERT INTO movie_cast (movie_id, cast_id) VALUES ($1,$2)", id, cid)
			if err != nil {
				return models.MovieAdmin{}, fmt.Errorf("insert casts failed: %w", err)
			}
		}
		m.Casts = casts
	}

	// --- Commit transaction ---
	if err := tx.Commit(ctx); err != nil {
		return models.MovieAdmin{}, fmt.Errorf("commit failed: %w", err)
	}

	r.InvalidateMovieCache(ctx)

	return m, nil
}

func (r *MovieAdminRepository) CreateMovie(ctx context.Context, movie models.Movie) (models.Movie, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models.Movie{}, err
	}
	defer tx.Rollback(ctx)

	const qMovie = `
        INSERT INTO movies (title, synopsis, releaseDate, duration, director_id, popularity, poster, background_poster)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
        RETURNING id, title, synopsis, releaseDate, duration, director_id, popularity, poster, background_poster, created_at, updated_at;
    `

	var m models.Movie
	var duration pgtype.Interval

	err = tx.QueryRow(ctx, qMovie,
		movie.Title,
		movie.Synopsis,
		movie.ReleaseDate,
		movie.Duration,
		movie.DirectorID,
		movie.Popularity,
		movie.Poster,
		movie.BackgroundPoster,
	).Scan(
		&m.ID,
		&m.Title,
		&m.Synopsis,
		&m.ReleaseDate,
		&duration,
		&m.DirectorID,
		&m.Popularity,
		&m.Poster,
		&m.BackgroundPoster,
		&m.CreatedAt,
		&m.UpdatedAt,
	)
	if err != nil {
		return m, fmt.Errorf("insert movie failed: %w", err)
	}
	m.DurationText = utils.FormatIntervalToText(m.Duration)

	// simpan duration ke struct
	m.Duration = duration

	// Insert ke movie_genre
	if len(movie.Genres) > 0 {
		const qGenre = `INSERT INTO movie_genre (movie_id, genre_id) VALUES ($1, $2)`
		for _, gid := range movie.Genres {
			_, err := tx.Exec(ctx, qGenre, m.ID, gid)
			if err != nil {
				return m, fmt.Errorf("insert genre failed: %w", err)
			}
		}
	}

	// Insert ke movie_cast
	if len(movie.Casts) > 0 {
		const qCast = `INSERT INTO movie_cast (movie_id, cast_id) VALUES ($1, $2)`
		for _, cid := range movie.Casts {
			_, err := tx.Exec(ctx, qCast, m.ID, cid)
			if err != nil {
				return m, fmt.Errorf("insert cast failed: %w", err)
			}
		}
	}

	// Ambil ulang genre IDs
	rows, err := tx.Query(ctx, "SELECT genre_id FROM movie_genre WHERE movie_id=$1", m.ID)
	if err != nil {
		return m, err
	}
	defer rows.Close()
	for rows.Next() {
		var gid int
		if err := rows.Scan(&gid); err != nil {
			return m, err
		}
		m.Genres = append(m.Genres, gid)
	}

	// Ambil ulang cast IDs
	rows2, err := tx.Query(ctx, "SELECT cast_id FROM movie_cast WHERE movie_id=$1", m.ID)
	if err != nil {
		return m, err
	}
	defer rows2.Close()
	for rows2.Next() {
		var cid int
		if err := rows2.Scan(&cid); err != nil {
			return m, err
		}
		m.Casts = append(m.Casts, cid)
	}

	// commit transaksi
	if err = tx.Commit(ctx); err != nil {
		return m, err
	}
	r.InvalidateMovieCache(ctx)

	return m, nil
}
