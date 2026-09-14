package service

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"

	"forum/internal/database"
	"forum/internal/repository"
)

func TestAuthServiceRegisterLoginAndAuthenticate(t *testing.T) {
	db, err := database.InitDB(filepath.Join(t.TempDir(), "forum.db"), filepath.Join("..", "database", "schema.sql"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	svc := NewAuthService(repository.NewUserRepository(db), repository.NewSessionRepository(db))
	ctx := context.Background()
	if err := svc.Register(ctx, "builder", "builder@example.com", "secret"); err != nil {
		t.Fatal(err)
	}
	if err := svc.Register(ctx, "builder2", "builder@example.com", "secret"); !errors.Is(err, ErrEmailTaken) {
		t.Fatalf("expected duplicate email error, got %v", err)
	}
	session, err := svc.Login(ctx, "builder@example.com", "secret", time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Authenticate(ctx, session.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Login(ctx, "builder@example.com", "wrong", time.Hour); !errors.Is(err, ErrInvalidAuth) {
		t.Fatalf("expected invalid auth, got %v", err)
	}
	if err := svc.Logout(ctx, session.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.Authenticate(ctx, session.ID); err == nil {
		t.Fatal("expected logged-out session to be invalid")
	}
}

func TestAuthServiceValidation(t *testing.T) {
	svc := NewAuthService(nil, nil)
	if !errors.Is(svc.Register(context.Background(), "", "", ""), ErrMissingFields) {
		t.Fatal("expected missing fields error")
	}
}
