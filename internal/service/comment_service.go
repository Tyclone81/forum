package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"forum/internal/models"
	"forum/internal/repository"
)

var ErrEmptyComment = errors.New("comment text content cannot be blank")

type CommentService struct {
	commentRepo *repository.CommentRepository
}

func NewCommentService(repo *repository.CommentRepository) *CommentService {
	return &CommentService{commentRepo: repo}
}

func (s *CommentService) AddComment(ctx context.Context, postID, userID, content string) error {
	if content == "" {
		return ErrEmptyComment
	}

	newComment := &models.Comment{
		ID:        uuid.New().String(),
		PostID:    postID,
		UserID:    userID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	return s.commentRepo.Create(ctx, newComment)
}

func (s *CommentService) FetchCommentsByPost(ctx context.Context, postID string) ([]*models.Comment, error) {
	return s.commentRepo.GetByPostID(ctx, postID)
}
