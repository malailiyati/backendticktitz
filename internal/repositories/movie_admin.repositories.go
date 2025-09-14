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
)

type MovieAdminRepository struct {
	db *pgxpool.Pool
}

func NewMovieAdminRepository(db *pgxpool.Pool) *MovieAdminRepository {
	return &MovieAdminRepository{db: db}
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

	return nil
}

func (r *MovieAdminRepository) UpdateMovie(ctx context.Context, id int, updates map[string]interface{}) (models.MovieAdmin, error) {
	if len(updates) == 0 {
		return models.MovieAdmin{}, fmt.Errorf("no fields to update")
	}

	setParts := []string{}
	args := []interface{}{}
	i := 1
	for k, v := range updates {
		setParts = append(setParts, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}
	// updated_at selalu di-update
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

	var m models.MovieAdmin
	var iv pgtype.Interval

	err := r.db.QueryRow(ctx, sql, args...).Scan(
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
		// tampilkan error asli, jangan ditiban
		return models.MovieAdmin{}, fmt.Errorf("update failed: %w", err)
	}

	// convert interval -> string human readable
	m.Duration = iv
	m.DurationText = utils.FormatIntervalToText(iv)

	return m, nil
}
