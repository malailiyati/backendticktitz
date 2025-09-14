package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/utils"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewMovieRepository(db *pgxpool.Pool) *MovieRepository {
	return &MovieRepository{db: db}
}

func (r *MovieRepository) GetUpcomingMovies(ctx context.Context) ([]models.Movie, error) {
	const q = `
		SELECT id, title, director_id, poster, background_poster,
       releaseDate, duration::text, synopsis, popularity, created_at, updated_at
	   FROM movies
	   WHERE releaseDate > CURRENT_DATE
	   ORDER BY releaseDate ASC;
	`

	rows, err := r.db.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.DirectorID,
			&m.Poster,
			&m.BackgroundPoster,
			&m.ReleaseDate,
			&m.Duration,
			&m.Synopsis,
			&m.Popularity,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetPopularMovies(ctx context.Context, limit int) ([]models.Movie, error) {
	const q = `
		SELECT id, title, director_id, poster, background_poster,
		       releaseDate, duration, synopsis, popularity, created_at, updated_at
		FROM movies
		ORDER BY popularity DESC
		LIMIT $1
	`
	rows, err := r.db.Query(ctx, q, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.DirectorID,
			&m.Poster,
			&m.BackgroundPoster,
			&m.ReleaseDate,
			&m.Duration,
			&m.Synopsis,
			&m.Popularity,
			&m.CreatedAt,
			&m.UpdatedAt,
		); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetMoviesByFilter(ctx context.Context, title, genre string, limit, offset int) ([]models.MovieFilter, error) {
	q := `
		SELECT m.id, m.title, m.poster,
		       COALESCE(string_agg(g.name, ','), '') AS genres
		FROM movies m
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		WHERE 1=1
	`
	var args []interface{}
	idx := 1

	if title != "" {
		q += fmt.Sprintf(" AND m.title ILIKE $%d", idx)
		args = append(args, "%"+title+"%")
		idx++
	}

	if genre != "" {
		q += fmt.Sprintf(" AND g.name ILIKE $%d", idx)
		args = append(args, "%"+genre+"%")
		idx++
	}

	q += ` GROUP BY m.id, m.title, m.poster
	       LIMIT $` + fmt.Sprint(idx) + ` OFFSET $` + fmt.Sprint(idx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.MovieFilter
	for rows.Next() {
		var m models.MovieFilter
		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &m.Genres); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

func (r *MovieRepository) GetMovieDetailByID(ctx context.Context, id int) (*models.MovieDetail, error) {
	sql := `
   SELECT m.id, m.title, m.poster, m.background_poster, 
		       m.releaseDate, m.duration,  
		       m.synopsis, d.name AS director,
		       COALESCE(string_agg(DISTINCT g.name, ', '), '') AS genres,
		       COALESCE(string_agg(DISTINCT c.name, ', '), '') AS casts
		FROM movies m
		LEFT JOIN director d ON d.id = m.director_id
		LEFT JOIN movie_genre mg ON mg.movie_id = m.id
		LEFT JOIN genres g ON g.id = mg.genre_id
		LEFT JOIN movie_cast mc ON mc.movie_id = m.id
		LEFT JOIN casts c ON c.id = mc.cast_id
		WHERE m.id = $1
		GROUP BY m.id, m.title, m.poster, m.background_poster, 
		         m.releaseDate, m.duration, m.synopsis, d.name
    `

	row := r.db.QueryRow(ctx, sql, id)

	var md models.MovieDetail
	var iv pgtype.Interval

	err := row.Scan(
		&md.ID,
		&md.Title,
		&md.Poster,
		&md.BackgroundPoster,
		&md.ReleaseDate,
		&iv,
		&md.Synopsis,
		&md.Director,
		&md.Genres,
		&md.Casts,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	md.Duration = utils.FormatIntervalToText(iv)

	return &md, nil
}
