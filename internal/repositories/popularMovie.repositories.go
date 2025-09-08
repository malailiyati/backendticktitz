package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type MoviePopularRepository struct {
	db *pgxpool.Pool
}

func NewMoviePopularRepository(db *pgxpool.Pool) *MoviePopularRepository {
	return &MoviePopularRepository{db: db}
}

func (r *MoviePopularRepository) GetPopularMovies(ctx context.Context, limit int) ([]models.Movie, error) {
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
