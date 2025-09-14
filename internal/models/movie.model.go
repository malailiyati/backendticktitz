package models

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Movie struct {
	ID               int             `json:"id"`
	Title            string          `json:"title"`
	DirectorID       int             `json:"director_id"`
	Poster           string          `json:"poster"`
	BackgroundPoster string          `json:"background_poster"`
	ReleaseDate      time.Time       `json:"release_date"`
	Duration         pgtype.Interval `json:"-"`
	Synopsis         string          `json:"synopsis"`
	Popularity       int             `json:"popularity"`
	CreatedAt        time.Time       `json:"created_at"`
	UpdatedAt        time.Time       `json:"updated_at"`
}

// MovieDetail: hasil query join movie + director + genre + cast
type MovieDetail struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Poster           string    `json:"poster"`
	BackgroundPoster string    `json:"background_poster"`
	ReleaseDate      time.Time `json:"release_date"`
	Duration         string    `json:"duration"` // cast dari interval ke text
	Synopsis         string    `json:"synopsis"`
	Director         string    `json:"director"`
	Genres           string    `json:"genres"`
	Casts            string    `json:"casts"`
}

type MovieFilter struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Poster string `json:"poster"`
	Genres string `json:"genres"` // hasil string_agg
}

// response struct khusus API

type MovieResponse struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	DirectorID       int       `json:"director_id"`
	Poster           string    `json:"poster"`
	BackgroundPoster string    `json:"background_poster"`
	ReleaseDate      time.Time `json:"release_date"`
	Duration         string    `json:"duration"`
	Synopsis         string    `json:"synopsis"`
	Popularity       int       `json:"popularity"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
