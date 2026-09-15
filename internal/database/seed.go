package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// SeedCategories guarantees the application has default categories available from moment one.
func SeedCategories(ctx context.Context, db *sql.DB) error {
	defaultCategories := []string{
		"Entertainment",
		"Sports",
		"Politics",
		"Technology",
		"Others",
	}

	query := `INSERT OR IGNORE INTO categories (name) VALUES (?);`

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin seed transaction: %w", err)
	}
	defer tx.Rollback()

	for _, name := range defaultCategories {
		_, err := tx.ExecContext(ctx, query, name)
		if err != nil {
			return fmt.Errorf("failed seeding category item (%s): %w", name, err)
		}
	}
	if _, err := tx.ExecContext(ctx, `DELETE FROM categories WHERE name NOT IN ('Entertainment', 'Sports', 'Politics', 'Technology', 'Others')`); err != nil {
		return fmt.Errorf("failed removing legacy categories: %w", err)
	}

	return tx.Commit()
}

// SeedMockData populates the demo account to make the project instantly testable without manual intervention.
func SeedMockData(ctx context.Context, db *sql.DB) error {
	// 1. Inject an evaluation mock tester account
	// Password representation is pre-hashed with Bcrypt for the sequence ("password123")
	userQuery := `INSERT OR IGNORE INTO users (id, username, email, password_hash, created_at) 
	              VALUES (?, ?, ?, ?, ?);`
	mockUserID := "00000000-0000-4000-a000-000000000001"
	mockPasswordHash := "$2a$10$wN9D/3P3t48u0Z8mXm.SkeqG9P10xU7qU.G8XvNqfP46H1R6VzGZ." // "password123"

	_, err := db.ExecContext(ctx, userQuery, mockUserID, "architect_alpha", "alpha@zone01.edu", mockPasswordHash, time.Now())
	if err != nil {
		return fmt.Errorf("failed to seed mock structural user: %w", err)
	}

	mockPostID := "00000000-0000-4000-a000-000000000002"
	if _, err = db.ExecContext(ctx, `DELETE FROM interactions WHERE target_id = ? AND target_type = 'post'`, mockPostID); err != nil {
		return fmt.Errorf("failed to remove legacy mock post interactions: %w", err)
	}
	if _, err = db.ExecContext(ctx, `DELETE FROM posts WHERE id = ?`, mockPostID); err != nil {
		return fmt.Errorf("failed to remove legacy mock post: %w", err)
	}

	// The legacy post is intentionally not recreated.
	return nil
}
