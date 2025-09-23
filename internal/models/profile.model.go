package models

import (
	"mime/multipart"
	"time"
)

type Profile struct {
	ID             int       `db:"id" json:"id"`
	UserID         int       `db:"users_id" json:"user_id"`
	FirstName      string    `db:"firstname" json:"first_name"`
	LastName       string    `db:"lastname" json:"last_name"`
	Phone          string    `db:"phone" json:"phone"`
	Email          string    `db:"email" json:"email"`
	ProfilePicture string    `db:"profile_picture" json:"profile_picture"`
	CreatedAt      time.Time `db:"created_at" json:"created_at"`
	UpdatedAt      time.Time `db:"updated_at" json:"updated_at"`
}

// Untuk request update profile dengan upload file
type ProfileBody struct {
	FirstName          *string               `form:"first_name"`
	LastName           *string               `form:"last_name"`
	Phone              *string               `form:"phone"`
	Email              *string               `form:"email"`
	ProfilePictureFile *multipart.FileHeader `form:"profile_picture"`
}

// ProfileResponse: gabungan dari users + profile
type ProfileResponse struct {
	UserID         int       `json:"user_id"`
	Email          string    `json:"email"`
	Role           string    `json:"role"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Phone          string    `json:"phone"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
