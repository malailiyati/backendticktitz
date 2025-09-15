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

func (r *ScheduleRepository) GetSchedules(ctx context.Context, date string, timeID, locationID, cinemaID int) ([]models.ScheduleDetail, error) {
	const q = `
		SELECT s.id, s.movie_id, m.title, s.date, t.time, l.location, c.name, c.price
		FROM schedule s
		JOIN movies m ON m.id = s.movie_id
		JOIN time t ON t.id = s.time_id
		JOIN location l ON l.id = s.location_id
		JOIN cinema c ON c.id = s.cinema_id
		WHERE s.date = $1
		  AND s.time_id = $2
		  AND s.location_id = $3
		  AND s.cinema_id = $4
		ORDER BY m.title;
	`

	rows, err := r.db.Query(ctx, q, date, timeID, locationID, cinemaID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.ScheduleDetail
	for rows.Next() {
		var s models.ScheduleDetail
		if err := rows.Scan(&s.ID, &s.MovieID, &s.Title, &s.Date, &s.Time, &s.Location, &s.Cinema, &s.Price); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}
