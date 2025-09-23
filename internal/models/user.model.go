package models

import "time"

// import "mime/multipart"

type User struct {
	Id        int       `db:"id" json:"id"`
	Email     string    `db:"email" json:"email"`
	Role      string    `db:"role" json:"role,omitempty"`
	Password  string    `db:"password" json:"password,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
	// Images   string `db:"images" json:"image"`
}

type UserBody struct {
	User
	// Images *multipart.FileHeader `form:"image"`
}

type UserAuth struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldpassword" binding:"required"`
	NewPassword string `json:"newpassword" binding:"required,min=8"`
	// ConfirmPassword string `json:"confirm_password" binding:"required"`
}
