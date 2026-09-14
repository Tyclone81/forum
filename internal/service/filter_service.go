package service

import (
	"context"
	"errors"
	"forum/internal/models"
	"forum/internal/repository"
)

type FilterService struct {
	postRepo        *repository.PostRepository
	interactionRepo *repository.InteractionRepository
}

func NewFilterService(repo *repository.PostRepository, interactionRepos ...*repository.InteractionRepository) *FilterService {
	service := &FilterService{postRepo: repo}
	if len(interactionRepos) > 0 {
		service.interactionRepo = interactionRepos[0]
	}
	return service
}

// FilterByCategory returns posts containing a specific target subforum topic string.
func (s *FilterService) FilterByCategory(ctx context.Context, categoryName string) ([]*models.Post, error) {
	allPosts, err := s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*models.Post
	for _, post := range allPosts {
		for _, cat := range post.Categories {
			if cat == categoryName {
				filtered = append(filtered, post)
				break
			}
		}
	}
	return filtered, nil
}

// FilterByCreated returns posts authored by a specific user profile identifier.
func (s *FilterService) FilterByCreated(ctx context.Context, userID string) ([]*models.Post, error) {
	allPosts, err := s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*models.Post
	for _, post := range allPosts {
		if post.UserID == userID {
			filtered = append(filtered, post)
		}
	}
	return filtered, nil
}

// FilterByLiked filters posts that the target user evaluated with an active upvote interaction.
func (s *FilterService) FilterByLiked(ctx context.Context, userID string) ([]*models.Post, error) {
	if s.interactionRepo == nil {
		return nil, errors.New("interaction repository is not configured")
	}
	allPosts, err := s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	var filtered []*models.Post
	for _, post := range allPosts {
		liked, err := s.interactionRepo.HasLikedPost(ctx, userID, post.ID)
		if err != nil {
			return nil, err
		}
		if liked {
			filtered = append(filtered, post)
		}
	}
	return filtered, nil
}
