package models

import "time"

// dipakai untuk request body
type CreateOrderRequest struct {
	UserID     int    `json:"user_id"`
	ScheduleID int    `json:"schedule_id"`
	PaymentID  int    `json:"payment_id"`
	TotalPrice int    `json:"total_price"`
	FullName   string `json:"full_name"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	IsPaid     bool   `json:"is_paid"`
	QRCode     string `json:"qr_code"`
	SeatIDs    []int  `json:"seat_ids"`
}

// dipakai untuk response
type Order struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	ScheduleID int       `json:"schedule_id"`
	PaymentID  int       `json:"payment_id"`
	TotalPrice int       `json:"total_price"`
	FullName   string    `json:"full_name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	IsPaid     bool      `json:"is_paid"`
	QRCode     string    `json:"qr_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
