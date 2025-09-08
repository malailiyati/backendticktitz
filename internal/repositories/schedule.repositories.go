package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type ScheduleRepository struct {
	db *pgxpool.Pool
}

func NewScheduleRepository(db *pgxpool.Pool) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) GetSchedulesByMovie(ctx context.Context, movieID int) ([]models.ScheduleDetail, error) {
	const q = `
		SELECT s.id, s.movie_id, s.date, t.time, l.location, c.name, c.price
		FROM schedule s
		JOIN time t ON t.id = s.time_id
		JOIN location l ON l.id = s.location_id
		JOIN cinema c ON c.id = s.cinema_id
		INNER JOIN movies mv ON mv.id = s.movie_id
		WHERE mv.id = $1
		ORDER BY s.date, t.time
	`
	rows, err := r.db.Query(ctx, q, movieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.ScheduleDetail
	for rows.Next() {
		var s models.ScheduleDetail
		if err := rows.Scan(&s.ID, &s.MovieID, &s.Date, &s.Time, &s.Location, &s.Cinema, &s.Price); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}
