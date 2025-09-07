package models

import "time"

type OrderHistory struct {
	OrderID    int       `json:"order_id"`
	UserID     int       `json:"user_id"`
	MovieTitle string    `json:"movie_title"`
	CinemaName string    `json:"cinema_name"`
	Location   string    `json:"location"`
	Date       time.Time `json:"date"`
	Time       string    `json:"time"`
	Seats      string    `json:"seats"`
	TotalPrice int       `json:"total_price"`
	IsPaid     bool      `json:"is_paid"`
	QRCode     string    `json:"qr_code"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
