package repository

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"forum/internal/database"
	"forum/internal/models"
)

func TestUserRepositoryCreateAndLookup(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	repo := NewUserRepository(db)
	want := &models.User{ID: "user-1", Username: "builder", Email: "builder@example.com", PasswordHash: "hash", CreatedAt: time.Now()}
	if err := repo.Create(context.Background(), want); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByEmail(context.Background(), want.Email)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil || got.ID != want.ID || got.Username != want.Username {
		t.Fatalf("unexpected user: %#v", got)
	}
	byID, err := repo.GetByID(context.Background(), want.ID)
	if err != nil {
		t.Fatal(err)
	}
	if byID == nil || byID.Email != want.Email {
		t.Fatalf("unexpected user by id: %#v", byID)
	}
}
