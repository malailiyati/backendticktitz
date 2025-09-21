package repositories

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/malailiyati/backend/internal/models"
)

type ProfileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) *ProfileRepository {
	return &ProfileRepository{db: db}
}

func (r *ProfileRepository) UpdateProfile(ctx context.Context, userID int, updates map[string]interface{}) (*models.Profile, error) {
	if len(updates) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	setClause := []string{}
	args := []interface{}{}
	i := 1

	for k, v := range updates {
		setClause = append(setClause, fmt.Sprintf("%s = $%d", k, i))
		args = append(args, v)
		i++
	}

	query := fmt.Sprintf(`
		UPDATE profile
		SET %s, updated_at = CURRENT_TIMESTAMP
		WHERE users_id = $%d
		RETURNING id, users_id, firstname, lastname, phone, profile_picture, created_at, updated_at
	`, strings.Join(setClause, ", "), i)

	args = append(args, userID)

	var p models.Profile
	err = tx.QueryRow(ctx, query, args...).Scan(
		&p.ID, &p.UserID, &p.FirstName, &p.LastName,
		&p.Phone, &p.ProfilePicture,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("update profile failed: %w", err)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProfileRepository) GetProfileByUserID(ctx context.Context, userID int) (*models.ProfileResponse, error) {
	const q = `
		SELECT u.id AS user_id, u.email, u.role,
		       p.firstName, p.lastName, p.phone, 
		       p.profile_picture, p.created_at, p.updated_at
		FROM users u
		LEFT JOIN profile p ON p.users_id = u.id
		WHERE u.id = $1
	`

	var p models.ProfileResponse
	err := r.db.QueryRow(ctx, q, userID).Scan(
		&p.UserID,
		&p.Email,
		&p.Role,
		&p.FirstName,
		&p.LastName,
		&p.Phone,
		&p.ProfilePicture,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (r *ProfileRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	const q = `
		SELECT id, email, role, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var u models.User
	err := r.db.QueryRow(ctx, q, userID).Scan(
		&u.Id,
		&u.Email,
		&u.Role,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *ProfileRepository) UpdatePassword(ctx context.Context, userID int, hashed string) error {
	const q = `
		UPDATE users
		SET password = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	_, err := r.db.Exec(ctx, q, hashed, userID)
	return err
}
