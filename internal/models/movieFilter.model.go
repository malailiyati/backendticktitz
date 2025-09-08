package models

type MovieFilter struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Poster string `json:"poster"`
	Genres string `json:"genres"` // hasil string_agg
}
