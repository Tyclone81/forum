package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"forum/internal/database"
	"forum/internal/models"
)

func TestCommentRepositoryCreateAndList(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	_, err = db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash'); INSERT INTO posts (id, user_id, title, content) VALUES ('post-1', 'user-1', 'Title', 'Content')`)
	if err != nil {
		t.Fatal(err)
	}

	repo := NewCommentRepository(db)
	want := &models.Comment{ID: "comment-1", PostID: "post-1", UserID: "user-1", Content: "Useful thought", CreatedAt: time.Now()}
	if err := repo.Create(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	comments, err := repo.GetByPostID(context.Background(), "post-1")
	if err != nil || len(comments) != 1 {
		t.Fatalf("comments=%#v err=%v", comments, err)
	}
	if comments[0].Username != "builder" || comments[0].Content != want.Content {
		t.Fatalf("unexpected comment: %#v", comments[0])
	}
}
