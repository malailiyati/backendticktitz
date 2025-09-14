package repositories

import (
	"context"
	"fmt"

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

	// loop kursi
	for _, seatID := range req.SeatIDs {
		var exists bool
		// cek apakah seat sudah dipakai di order lain (yang sudah dibayar)
		err := tx.QueryRow(ctx, `
			SELECT EXISTS (
				SELECT 1
				FROM order_seat os
				JOIN orders o ON o.id = os.order_id
				WHERE os.seat_id = $1
				  AND o.schedule_id = $2
				  AND o.isPaid = true
			)
		`, seatID, req.ScheduleID).Scan(&exists)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("seat %d sudah terjual", seatID)
		}

		// insert kursi ke order_seat
		_, err = tx.Exec(ctx,
			`INSERT INTO order_seat (order_id, seat_id) VALUES ($1, $2)`,
			order.ID, seatID,
		)
		if err != nil {
			return nil, err
		}
	}

	// commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}
	return &order, nil
}
