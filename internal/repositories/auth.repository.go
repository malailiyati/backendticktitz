package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: db}
}

func (a *AuthRepository) GetUserWithPasswordAndRole(rctx context.Context, email string) (models.User, error) {
	// validasi user
	// ambil data user berdasarkan input user
	sql := "SELECT id, email, password, role FROM users WHERE email = $1"

	var user models.User
	if err := a.db.QueryRow(rctx, sql, email).Scan(&user.Id, &user.Email, &user.Password, &user.Role); err != nil {
		if err == pgx.ErrNoRows {
			return models.User{}, errors.New("user not found")
		}
		log.Println("Internal Server Error.\nCause: ", err.Error())
		return models.User{}, err
	}
	return user, nil
}

func (a *AuthRepository) CreateUser(ctx context.Context, email, hashedPass string, role *string) (models.User, error) {
	sql := `
        INSERT INTO users (email, password, role) 
        VALUES ($1, $2, $3)
        RETURNING id, email, password, role, created_at, updated_at;
    `

	var user models.User
	if err := a.db.QueryRow(ctx, sql, email, hashedPass, role).Scan(
		&user.Id,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {
		return models.User{}, err
	}
	return user, nil
}
