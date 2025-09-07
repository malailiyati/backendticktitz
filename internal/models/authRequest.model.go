package models

type Register struct {
	Email    string  `json:"email" binding:"required"`
	Password string  `json:"password" binding:"required,min=8"`
	Role     *string `json:"role,omitempty" extensions:"x-omitempty"`
}

type Login struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=8"`
}
