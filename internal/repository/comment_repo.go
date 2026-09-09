package repository

import (
	"context"
	"database/sql"
	"forum/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	query := `INSERT INTO comments (id, post_id, user_id, content, created_at) VALUES (?, ?, ?, ?, ?);`
	_, err := r.db.ExecContext(ctx, query, comment.ID, comment.PostID, comment.UserID, comment.Content, comment.CreatedAt)
	return err
}

func (r *CommentRepository) GetByPostID(ctx context.Context, postID string) ([]*models.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, u.username, c.content, c.created_at,
		       COALESCE((SELECT SUM(value) FROM interactions WHERE target_id = c.id AND target_type = 'comment' AND value = 1), 0) as likes,
		       COALESCE((SELECT ABS(SUM(value)) FROM interactions WHERE target_id = c.id AND target_type = 'comment' AND value = -1), 0) as dislikes
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC;`

	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Username, &c.Content, &c.CreatedAt, &c.Likes, &c.Dislikes)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &c)
	}
	return comments, nil
}
