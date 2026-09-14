package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"forum/internal/database"
	"forum/internal/repository"
)

func TestInteractionServiceValidatesReactions(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash'); INSERT INTO posts (id, user_id, title, content) VALUES ('post-1', 'user-1', 'Title', 'Content')`); err != nil {
		t.Fatal(err)
	}

	svc := NewInteractionService(repository.NewInteractionRepository(db))
	ctx := context.Background()
	if !errors.Is(svc.SetReaction(ctx, "user-1", "post-1", "profile", 1), ErrInvalidInteractionTarget) {
		t.Fatal("expected invalid target error")
	}
	if !errors.Is(svc.SetReaction(ctx, "user-1", "post-1", "post", 0), ErrInvalidInteractionValue) {
		t.Fatal("expected invalid value error")
	}
	if err := svc.SetReaction(ctx, "user-1", "post-1", "post", 1); err != nil {
		t.Fatal(err)
	}
}
