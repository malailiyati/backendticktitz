package models

import "time"

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
