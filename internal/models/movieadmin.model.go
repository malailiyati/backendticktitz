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
	Genres           []int           `json:"genres"`
	Casts            []int           `json:"casts"`
	CreatedAt        *time.Time      `json:"created_at"`
	UpdatedAt        *time.Time      `json:"updated_at"`
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
	Genres           string                `form:"genres"`
	Casts            string                `form:"casts"`
}

type CreateMovieAdminBody struct {
	Title            string                `form:"title" binding:"required"`
	Synopsis         string                `form:"synopsis" binding:"required"`
	ReleaseDate      string                `form:"release_date" binding:"required"`
	Duration         string                `form:"duration" binding:"required"`
	DirectorID       int                   `form:"director_id" binding:"required"`
	Popularity       int                   `form:"popularity"`
	Poster           *multipart.FileHeader `form:"poster"`
	BackgroundPoster *multipart.FileHeader `form:"background_poster"`
	Genres           string                `form:"genres"`
	Casts            string                `form:"casts"`
}

type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Cast struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
