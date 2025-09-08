package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
	"github.com/malailiyati/backend/internal/utils"
)

type MovieDetailRepository struct {
	db *pgxpool.Pool
}

func NewMovieDetailRepository(db *pgxpool.Pool) *MovieDetailRepository {
	return &MovieDetailRepository{db: db}
}

func (r *MovieDetailRepository) GetMovieDetail(ctx context.Context, movieID int) (*models.MovieDetail, error) {
	const q = `
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

	var detail models.MovieDetail
	var rawDuration pgtype.Interval //  scan ke interval

	err := r.db.QueryRow(ctx, q, movieID).Scan(
		&detail.ID,
		&detail.Title,
		&detail.Poster,
		&detail.BackgroundPoster,
		&detail.ReleaseDate,
		&rawDuration,
		&detail.Synopsis,
		&detail.Director,
		&detail.Genres,
		&detail.Casts,
	)
	if err != nil {
		return nil, err
	}

	// format dengan utils
	detail.Duration = utils.FormatIntervalToText(rawDuration)

	return &detail, nil
}
