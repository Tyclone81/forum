package repository

import (
	"context"
	"path/filepath"
	"testing"

	"forum/internal/database"
	"forum/internal/models"
)

func TestInteractionRepositoryReplacesReactionAndChecksLikes(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash'); INSERT INTO posts (id, user_id, title, content) VALUES ('post-1', 'user-1', 'Title', 'Content')`); err != nil {
		t.Fatal(err)
	}

	repo := NewInteractionRepository(db)
	ctx := context.Background()
	first := &models.Interaction{ID: "reaction-1", UserID: "user-1", TargetID: "post-1", TargetType: "post", Value: 1}
	if err := repo.Set(ctx, first); err != nil {
		t.Fatal(err)
	}
	liked, err := repo.HasLikedPost(ctx, "user-1", "post-1")
	if err != nil || !liked {
		t.Fatalf("liked=%v err=%v", liked, err)
	}

	second := &models.Interaction{ID: "reaction-2", UserID: "user-1", TargetID: "post-1", TargetType: "post", Value: -1}
	if err := repo.Set(ctx, second); err != nil {
		t.Fatal(err)
	}
	liked, err = repo.HasLikedPost(ctx, "user-1", "post-1")
	if err != nil || liked {
		t.Fatalf("expected dislike to replace like, liked=%v err=%v", liked, err)
	}
	if err := repo.Set(ctx, second); err != nil {
		t.Fatal(err)
	}
	liked, err = repo.HasLikedPost(ctx, "user-1", "post-1")
	if err != nil || liked {
		t.Fatalf("expected repeated dislike to clear reaction, liked=%v err=%v", liked, err)
	}
}
