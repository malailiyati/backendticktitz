package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type MovieFilterRepository struct {
	db *pgxpool.Pool
}

func NewMovieFilterRepository(db *pgxpool.Pool) *MovieFilterRepository {
	return &MovieFilterRepository{db: db}
}

func (r *MovieFilterRepository) GetMoviesByFilter(ctx context.Context, title, genre string, limit, offset int) ([]models.MovieFilter, error) {
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
