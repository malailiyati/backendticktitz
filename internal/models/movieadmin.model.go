package models

import (
	"mime/multipart"
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

type UpdateMovieAdminBody struct {
	Title            string                `form:"title"`
	DirectorID       string                `form:"director_id"`
	Poster           *multipart.FileHeader `form:"poster"`
	BackgroundPoster *multipart.FileHeader `form:"background_poster"`
	ReleaseDate      string                `form:"release_date"` // YYYY-MM-DD
	Duration         string                `form:"duration"`     // e.g. 02:35
	Synopsis         string                `form:"synopsis"`
	Popularity       string                `form:"popularity"`
}
