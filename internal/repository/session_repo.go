package repository

import (
	"context"
	"database/sql"
	"errors"
	"forum/internal/models"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(ctx context.Context, sess *models.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, sess.ID, sess.UserID, sess.ExpiresAt)
	return err
}

func (r *SessionRepository) GetByID(ctx context.Context, id string) (*models.Session, error) {
	query := `SELECT id, user_id, expires_at FROM sessions WHERE id = ?;`
	row := r.db.QueryRowContext(ctx, query, id)

	var sess models.Session
	err := row.Scan(&sess.ID, &sess.UserID, &sess.ExpiresAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &sess, nil
}

func (r *SessionRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM sessions WHERE id = ?;`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *SessionRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = ?`, userID)
	return err
}
