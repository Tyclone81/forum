package repository

import (
	"context"
	"database/sql"
	"errors"
	"forum/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create inserts a new registered user account. (Fulfills INSERT query instruction)
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, username, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, user.ID, user.Username, user.Email, user.PasswordHash, user.CreatedAt)
	return err
}

// GetByEmail searches for a user by email to handle credential checks. (Fulfills SELECT query instruction)
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = ?;`
	row := r.db.QueryRowContext(ctx, query, email)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return explicit nil if no record matches
		}
		return nil, err
	}
	return &user, nil
}

// GetByID retrieves structural profile details by user ID context.
func (r *UserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?;`
	row := r.db.QueryRowContext(ctx, query, id)

	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
