package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"forum/internal/models"
	"forum/internal/repository"
)

var (
	ErrEmailTaken    = errors.New("the provided email address is already registered")
	ErrInvalidAuth   = errors.New("invalid email or password credentials provided")
	ErrMissingFields = errors.New("all structural registration input fields are required")
)

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(u *repository.UserRepository, s *repository.SessionRepository) *AuthService {
	return &AuthService{userRepo: u, sessionRepo: s}
}

// Register processes structural signups. (Fulfills unique email validation constraints)
func (s *AuthService) Register(ctx context.Context, username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return ErrMissingFields
	}

	// 1. Verify if user email unique context constraint exists
	existingUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("failed during registration lookup validation: %w", err)
	}
	if existingUser != nil {
		return ErrEmailTaken
	}

	// 2. Encrypt the password using allowed standard cryptographic hashing (Bcrypt)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed executing password encryption routines: %w", err)
	}

	// 3. Populate new domain entity model record using clean modern UUID generation
	newUser := &models.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		CreatedAt:    time.Now(),
	}

	return s.userRepo.Create(ctx, newUser)
}

// Login verifies identity and issues an authenticated session token.
func (s *AuthService) Login(ctx context.Context, email, password string, sessionDuration time.Duration) (*models.Session, error) {
	// 1. Retrieve the registered user object profile matching target credentials
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed lookup check during login phase: %w", err)
	}
	if user == nil {
		return nil, ErrInvalidAuth
	}

	// 2. Validate passwords against encrypted database record content
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidAuth // Hand back clear uninformative user error status
	}

	// 3. Establish single unique tracker session utilizing an isolated token string
	newSession := &models.Session{
		ID:        uuid.New().String(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(sessionDuration),
	}

	if err := s.sessionRepo.Create(ctx, newSession); err != nil {
		return nil, fmt.Errorf("failed saving user session reference metadata: %w", err)
	}

	return newSession, nil
}

// Logout tears down active authorization references
func (s *AuthService) Logout(ctx context.Context, token string) error {
	return s.sessionRepo.Delete(ctx, token)
}

// Authenticate verifies if a session token string is valid and active.
func (s *AuthService) Authenticate(ctx context.Context, token string) (*models.Session, error) {
	session, err := s.sessionRepo.GetByID(ctx, token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session matching reference not found")
	}

	// Validate timeline timeout lifespan values
	if time.Now().After(session.ExpiresAt) {
		_ = s.sessionRepo.Delete(ctx, token) // Silently clean expired reference tracking blocks
		return nil, errors.New("authenticated session expired")
	}

	return session, nil
}
