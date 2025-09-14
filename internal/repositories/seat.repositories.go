package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type SeatRepository struct {
	db *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{db: db}
}

// Ambil kursi untuk schedule tertentu
func (r *SeatRepository) GetAvailableSeats(ctx context.Context, scheduleID int) ([]models.Seat, error) {
	const q = `
		SELECT s.id, s.seat_number
		FROM seats s
		WHERE s.id IN (
		SELECT os.seat_id
		FROM order_seat os
		JOIN orders o ON o.id = os.order_id
		WHERE o.schedule_id = $1
		AND o.isPaid = true
		)
		ORDER BY s.seat_number
	`
	rows, err := r.db.Query(ctx, q, scheduleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var seats []models.Seat
	for rows.Next() {
		var seat models.Seat
		if err := rows.Scan(&seat.ID, &seat.SeatNumber); err != nil {
			return nil, err
		}
		seats = append(seats, seat)
	}
	return seats, nil
}
