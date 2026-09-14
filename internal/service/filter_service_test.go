package service

import (
	"context"
	"path/filepath"
	"testing"

	"forum/internal/database"
	"forum/internal/repository"
)

func TestFilterServiceByCategoryAndAuthor(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'one', 'one@example.com', 'hash'), ('user-2', 'two', 'two@example.com', 'hash'); INSERT INTO categories (name) VALUES ('Technology'), ('Sports'); INSERT INTO posts (id, user_id, title, content) VALUES ('post-1', 'user-1', 'Technology', 'body'), ('post-2', 'user-2', 'Sports', 'body'); INSERT INTO post_categories (post_id, category_id) SELECT 'post-1', id FROM categories WHERE name = 'Technology'; INSERT INTO post_categories (post_id, category_id) SELECT 'post-2', id FROM categories WHERE name = 'Sports'`)
	if err != nil {
		t.Fatal(err)
	}
	svc := NewFilterService(repository.NewPostRepository(db))
	ctx := context.Background()
	byCategory, err := svc.FilterByCategory(ctx, "Technology")
	if err != nil || len(byCategory) != 1 || byCategory[0].ID != "post-1" {
		t.Fatalf("category result=%#v err=%v", byCategory, err)
	}
	byAuthor, err := svc.FilterByCreated(ctx, "user-2")
	if err != nil || len(byAuthor) != 1 || byAuthor[0].ID != "post-2" {
		t.Fatalf("author result=%#v err=%v", byAuthor, err)
	}
}
