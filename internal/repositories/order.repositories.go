package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, req models.CreateOrderRequest) (*models.Order, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	var order models.Order
	err = tx.QueryRow(ctx, `
		INSERT INTO orders 
		(users_id, schedule_id, payment_id, totalPrice, fullName, email, phone, isPaid, qr_code, created_at, updated_at) 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,CURRENT_TIMESTAMP,CURRENT_TIMESTAMP)
		RETURNING id, users_id, schedule_id, payment_id, totalPrice, fullName, email, phone, isPaid, qr_code, created_at, updated_at
	`,
		req.UserID, req.ScheduleID, req.PaymentID, req.TotalPrice,
		req.FullName, req.Email, req.Phone, req.IsPaid, req.QRCode,
	).Scan(
		&order.ID, &order.UserID, &order.ScheduleID, &order.PaymentID,
		&order.TotalPrice, &order.FullName, &order.Email, &order.Phone,
		&order.IsPaid, &order.QRCode, &order.CreatedAt, &order.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// insert kursi ke order_seat
	for _, seatID := range req.SeatIDs {
		_, err := tx.Exec(ctx,
			`INSERT INTO order_seat (order_id, seat_id) VALUES ($1, $2)`,
			order.ID, seatID,
		)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &order, nil
}
