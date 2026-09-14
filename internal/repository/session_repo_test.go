package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"forum/internal/database"
	"forum/internal/models"
)

func TestSessionRepositoryLifecycle(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.Exec(`INSERT INTO users (id, username, email, password_hash) VALUES ('user-1', 'builder', 'builder@example.com', 'hash')`); err != nil {
		t.Fatal(err)
	}

	repo := NewSessionRepository(db)
	want := &models.Session{ID: "session-1", UserID: "user-1", ExpiresAt: time.Now().Add(time.Hour)}
	if err := repo.Create(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByID(context.Background(), want.ID)
	if err != nil || got == nil || got.UserID != want.UserID {
		t.Fatalf("unexpected session: %#v, %v", got, err)
	}
	if err := repo.Delete(context.Background(), want.ID); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetByID(context.Background(), want.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got != nil {
		t.Fatalf("expected deleted session, got %#v", got)
	}
}
