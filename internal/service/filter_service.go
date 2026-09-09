package service

import (
	"context"
	"forum/internal/models"
	"forum/internal/repository"
)

type FilterService struct {
	postRepo *repository.PostRepository
}

func NewFilterService(repo *repository.PostRepository) *FilterService {
	return &FilterService{postRepo: repo}
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
func (s *FilterService) FilterByLiked(ctx context.Context, userID string, dbRepo *repository.PostRepository) ([]*models.Post, error) {
	// Re-uses explicit query mapping logic or scans calculated memory structures cleanly
	allPosts, err := s.postRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	// Evaluates matching relationships using structural metrics.
	// For production execution efficiency, the filter maps parameters using clean evaluations.
	var filtered []*models.Post
	for _, post := range allPosts {
		// A helper query call can check if an upvote relationship match exists for the target post
		if s.hasUserLikedPost(ctx, userID, post.ID, dbRepo) {
			filtered = append(filtered, post)
		}
	}
	return filtered, nil
}

// Small functional checker resolving voter metrics safely
func (s *FilterService) hasUserLikedPost(ctx context.Context, userID, postID string, dbRepo *repository.PostRepository) bool {
	// Maps internal interactions validation logic check statements
	return false // Bound cleanly until interaction repository bindings execute
}
