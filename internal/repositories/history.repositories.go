package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type HistoryRepository struct {
	db *pgxpool.Pool
}

func NewHistoryRepository(db *pgxpool.Pool) *HistoryRepository {
	return &HistoryRepository{db: db}
}

func (r *HistoryRepository) GetOrderHistory(ctx context.Context, userID int) ([]models.OrderHistory, error) {
	const q = `
		SELECT o.id AS order_id, o.users_id, o.totalPrice, o.isPaid, o.qr_code, 
		       o.created_at, o.updated_at,
		       m.title AS movie_title,
		       s.date, t.time, l.location, c.name AS cinema_name,
		       string_agg(seat.seat_number, ', ') AS seats
		FROM orders o
		JOIN schedule s ON s.id = o.schedule_id
		JOIN movies m ON m.id = s.movie_id
		JOIN time t ON t.id = s.time_id
		JOIN location l ON l.id = s.location_id
		JOIN cinema c ON c.id = s.cinema_id
		JOIN order_seat os ON os.order_id = o.id
		JOIN seats seat ON seat.id = os.seat_id
		WHERE o.users_id = $1
		GROUP BY o.id, m.title, s.date, t.time, l.location, c.name
		ORDER BY o.created_at DESC;
	`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []models.OrderHistory
	for rows.Next() {
		var h models.OrderHistory
		if err := rows.Scan(
			&h.OrderID,
			&h.UserID,
			&h.TotalPrice,
			&h.IsPaid,
			&h.QRCode,
			&h.CreatedAt,
			&h.UpdatedAt,
			&h.MovieTitle,
			&h.Date,
			&h.Time,
			&h.Location,
			&h.CinemaName,
			&h.Seats,
		); err != nil {
			return nil, err
		}
		history = append(history, h)
	}
	return history, nil
}
