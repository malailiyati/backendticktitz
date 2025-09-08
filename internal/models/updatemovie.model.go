package models

import (
	"mime/multipart"
)

// Body request edit movie
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
