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
		"Golang",
		"Docker",
		"AI & Architecture",
		"Vibe Coding",
		"Viable Systems",
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

	return tx.Commit()
}

// SeedMockData populates initial test data to make the project instantly testable without manual intervention.
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

	// 2. Inject an initial post thread detailing our architectural mindset
	postQuery := `INSERT OR IGNORE INTO posts (id, user_id, title, content, created_at) 
	              VALUES (?, ?, ?, ?, ?);`
	mockPostID := "00000000-0000-4000-a000-000000000002"

	_, err = db.ExecContext(ctx, postQuery, mockPostID, mockUserID, "The Leap to Architect", "Stop writing strings blindly. Own your blueprint and build viable software configurations natively.", time.Now())
	if err != nil {
		return fmt.Errorf("failed to seed structural mock post: %w", err)
	}

	// 3. Connect the post to the 'Viable Systems' category tag index cleanly
	junctionQuery := `INSERT OR IGNORE INTO post_categories (post_id, category_id) 
	                  SELECT ?, id FROM categories WHERE name = ?;`
	_, err = db.ExecContext(ctx, junctionQuery, mockPostID, "Viable Systems")
	if err != nil {
		return fmt.Errorf("failed to link structural category associations: %w", err)
	}

	return nil
}
