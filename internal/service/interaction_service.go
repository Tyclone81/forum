package service

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"forum/internal/models"
	"forum/internal/repository"
)

var (
	ErrEmptyInteractionTarget   = errors.New("interaction target is required")
	ErrInvalidInteractionTarget = errors.New("interaction target must be post or comment")
	ErrInvalidInteractionValue  = errors.New("interaction value must be 1 or -1")
)

type InteractionService struct {
	repo *repository.InteractionRepository
}

func NewInteractionService(repo *repository.InteractionRepository) *InteractionService {
	return &InteractionService{repo: repo}
}

func (s *InteractionService) SetReaction(ctx context.Context, userID, targetID, targetType string, value int) error {
	if targetID == "" {
		return ErrEmptyInteractionTarget
	}
	if targetType != "post" && targetType != "comment" {
		return ErrInvalidInteractionTarget
	}
	if value != 1 && value != -1 {
		return ErrInvalidInteractionValue
	}
	return s.repo.Set(ctx, &models.Interaction{
		ID:         uuid.New().String(),
		UserID:     userID,
		TargetID:   targetID,
		TargetType: targetType,
		Value:      value,
	})
}
