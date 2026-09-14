package repository

import (
	"context"
	"database/sql"
	"forum/internal/models"
)

type InteractionRepository struct {
	db *sql.DB
}

func NewInteractionRepository(db *sql.DB) *InteractionRepository {
	return &InteractionRepository{db: db}
}

func (r *InteractionRepository) Set(ctx context.Context, interaction *models.Interaction) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existingValue int
	err = tx.QueryRowContext(ctx, `SELECT value FROM interactions WHERE user_id = ? AND target_id = ? AND target_type = ?`, interaction.UserID, interaction.TargetID, interaction.TargetType).Scan(&existingValue)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err == nil && existingValue == interaction.Value {
		_, err = tx.ExecContext(ctx, `DELETE FROM interactions WHERE user_id = ? AND target_id = ? AND target_type = ?`, interaction.UserID, interaction.TargetID, interaction.TargetType)
		if err != nil {
			return err
		}
		return tx.Commit()
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM interactions WHERE user_id = ? AND target_id = ? AND target_type = ?`, interaction.UserID, interaction.TargetID, interaction.TargetType); err != nil {
		return err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO interactions (id, user_id, target_id, target_type, value) VALUES (?, ?, ?, ?, ?)`, interaction.ID, interaction.UserID, interaction.TargetID, interaction.TargetType, interaction.Value); err != nil {
		return err
	}
	return tx.Commit()
}

func (r *InteractionRepository) HasLikedPost(ctx context.Context, userID, postID string) (bool, error) {
	var exists int
	err := r.db.QueryRowContext(ctx, `SELECT 1 FROM interactions WHERE user_id = ? AND target_id = ? AND target_type = 'post' AND value = 1 LIMIT 1`, userID, postID).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return err == nil && exists == 1, err
}
