package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"forum/internal/models"
	"forum/internal/repository"
)

var (
	ErrEmptyPostContent = errors.New("post title and main content parameters cannot be empty")
	ErrNoCategories     = errors.New("a post must have at least one category tag mapped onto it")
	ErrPostNotFound     = errors.New("post not found or not owned by user")
)

type PostService struct {
	postRepo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{postRepo: repo}
}

func (s *PostService) CreatePost(ctx context.Context, userID, title, content string, categories []string) error {
	if title == "" || content == "" {
		return ErrEmptyPostContent
	}
	if len(categories) == 0 {
		return ErrNoCategories
	}

	newPost := &models.Post{
		ID:         uuid.New().String(),
		UserID:     userID,
		Title:      title,
		Content:    content,
		Categories: categories,
		CreatedAt:  time.Now(),
	}

	return s.postRepo.Create(ctx, newPost)
}

func (s *PostService) FetchAllPosts(ctx context.Context) ([]*models.Post, error) {
	return s.postRepo.GetAll(ctx)
}

func (s *PostService) UpdatePost(ctx context.Context, userID, postID, title, content string, categories []string) error {
	if title == "" || content == "" {
		return ErrEmptyPostContent
	}
	if len(categories) == 0 {
		return ErrNoCategories
	}
	if err := s.postRepo.Update(ctx, postID, userID, title, content, categories); errors.Is(err, sql.ErrNoRows) {
		return ErrPostNotFound
	} else if err != nil {
		return err
	}
	return nil
}

func (s *PostService) DeletePost(ctx context.Context, userID, postID string) error {
	if err := s.postRepo.Delete(ctx, postID, userID); errors.Is(err, sql.ErrNoRows) {
		return ErrPostNotFound
	} else if err != nil {
		return err
	}
	return nil
}
