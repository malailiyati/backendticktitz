package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type MovieAdmin struct {
	ID               int             `json:"id"`
	Title            string          `json:"title"`
	DirectorID       int             `json:"director_id"`
	Poster           string          `json:"poster"`
	BackgroundPoster string          `json:"background_poster"`
	ReleaseDate      time.Time       `json:"release_date"`
	Duration         pgtype.Interval `json:"-"`        // simpan raw interval
	DurationText     string          `json:"duration"` // untuk response JSON
	Synopsis         string          `json:"synopsis"`
	Popularity       int             `json:"popularity"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}
