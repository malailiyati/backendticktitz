package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type MovieRepository struct {
	db *pgxpool.Pool
}

func NewUpcomingMovieRepository(db *pgxpool.Pool) *MovieRepository {
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
