package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"

	"forum/internal/database"
	"forum/internal/repository"
)

func TestPostServiceValidationAndCreate(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash'); INSERT INTO categories (name) VALUES ('Technology')`); err != nil {
		t.Fatal(err)
	}
	svc := NewPostService(repository.NewPostRepository(db))
	ctx := context.Background()
	if !errors.Is(svc.CreatePost(ctx, "user-1", "", "body", []string{"Technology"}), ErrEmptyPostContent) {
		t.Fatal("expected empty post error")
	}
	if !errors.Is(svc.CreatePost(ctx, "user-1", "Title", "body", nil), ErrNoCategories) {
		t.Fatal("expected category error")
	}
	if err := svc.CreatePost(ctx, "user-1", "Title", "body", []string{"Technology"}); err != nil {
		t.Fatal(err)
	}
	posts, err := svc.FetchAllPosts(ctx)
	if err != nil || len(posts) != 1 || posts[0].Title != "Title" {
		t.Fatalf("posts=%#v err=%v", posts, err)
	}
}
