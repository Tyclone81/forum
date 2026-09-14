package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"forum/internal/database"
	"forum/internal/models"
)

func TestPostRepositoryCreateAndGetAll(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash'); INSERT INTO categories (name) VALUES ('Technology')`); err != nil {
		t.Fatal(err)
	}

	repo := NewPostRepository(db)
	want := &models.Post{ID: "post-1", UserID: "user-1", Title: "Title", Content: "Content", Categories: []string{"Technology"}, CreatedAt: time.Now()}
	if err := repo.Create(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	posts, err := repo.GetAll(context.Background())
	if err != nil || len(posts) != 1 {
		t.Fatalf("posts=%#v err=%v", posts, err)
	}
	if posts[0].Username != "builder" || len(posts[0].Categories) != 1 || posts[0].Categories[0] != "Technology" {
		t.Fatalf("unexpected post: %#v", posts[0])
	}
}
