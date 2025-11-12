package repository

import (
	"context"
	"fmt"
	"go-articles-app/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment *models.Comment) error {
	result := r.db.WithContext(ctx).Create(comment)
	return result.Error
}

func (r *CommentRepository) GetByArticleID(ctx context.Context, articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	result := r.db.WithContext(ctx).Where("article_id = ?", articleID).Preload("Author").Find(&comments)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get comments by article id : %w", result.Error)
	}

	return comments, nil
}

func (r *CommentRepository) GetReplies(ctx context.Context, parentID uint) ([]models.Comment, error) {
	var replies []models.Comment

	result := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Preload("User").Order("created_at ASC").Find(&replies)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get replies: %w", result.Error)
	}

	return replies, nil
}

func (r *CommentRepository) Delete(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&models.Comment{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete comments: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("comments not found")
	}

	return nil
}
