package models

import "time"

type Schedule struct {
	ID         int       `json:"id"`
	MovieID    int       `json:"movie_id"`
	TimeID     int       `json:"time_id"`
	LocationID int       `json:"location_id"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type ScheduleDetail struct {
	ID       int       `json:"id"`
	MovieID  int       `json:"movie_id"`
	Title    string    `json:"title"`
	Date     time.Time `json:"date"`
	Time     string    `json:"time"`
	Location string    `json:"location"`
	Cinema   string    `json:"cinema"`
	Price    float64   `json:"price"`
}
