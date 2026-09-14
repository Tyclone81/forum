package repository

import (
	"context"
	"database/sql"
	"forum/internal/models"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

// Create inserts a thread and binds categories within an atomic SQL Transaction
func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback() // Safe rollback fallback if execution breaks early

	// 1. Insert Core Post Entry
	postQuery := `INSERT INTO posts (id, user_id, title, content, created_at) VALUES (?, ?, ?, ?, ?);`
	_, err = tx.ExecContext(ctx, postQuery, post.ID, post.UserID, post.Title, post.Content, post.CreatedAt)
	if err != nil {
		return err
	}

	// 2. Map Multi-Category Relationships (Junction Table)
	catQuery := `INSERT OR IGNORE INTO post_categories (post_id, category_id) 
	             SELECT ?, id FROM categories WHERE name = ?;`
	for _, categoryName := range post.Categories {
		_, err = tx.ExecContext(ctx, catQuery, post.ID, categoryName)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *PostRepository) Update(ctx context.Context, postID, userID, title, content string, categories []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `UPDATE posts SET title = ?, content = ? WHERE id = ? AND user_id = ?`, title, content, postID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM post_categories WHERE post_id = ?`, postID); err != nil {
		return err
	}
	for _, categoryName := range categories {
		if _, err := tx.ExecContext(ctx, `INSERT OR IGNORE INTO post_categories (post_id, category_id) SELECT ?, id FROM categories WHERE name = ?`, postID, categoryName); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *PostRepository) Delete(ctx context.Context, postID, userID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM posts WHERE id = ? AND user_id = ?`, postID, userID)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

// GetAll fetches baseline threads decorated with computed likes/dislikes metrics.
func (r *PostRepository) GetAll(ctx context.Context) ([]*models.Post, error) {
	query := `
		SELECT p.id, p.user_id, u.username, p.title, p.content, p.created_at,
		       COALESCE((SELECT SUM(value) FROM interactions WHERE target_id = p.id AND target_type = 'post' AND value = 1), 0) as likes,
		       COALESCE((SELECT ABS(SUM(value)) FROM interactions WHERE target_id = p.id AND target_type = 'post' AND value = -1), 0) as dislikes
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC;`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	var posts []*models.Post
	for rows.Next() {
		var p models.Post
		err := rows.Scan(&p.ID, &p.UserID, &p.Username, &p.Title, &p.Content, &p.CreatedAt, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &p)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, err
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}

	for _, post := range posts {
		post.Categories, _ = r.GetCategoriesByPostID(ctx, post.ID)
	}
	return posts, nil
}

// GetCategoriesByPostID pulls structural categories mapped to a target post
func (r *PostRepository) GetCategoriesByPostID(ctx context.Context, postID string) ([]string, error) {
	query := `SELECT c.name FROM categories c 
	          JOIN post_categories pc ON c.id = pc.category_id 
	          WHERE pc.post_id = ?;`
	rows, err := r.db.QueryContext(ctx, query, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			categories = append(categories, name)
		}
	}
	return categories, nil
}
