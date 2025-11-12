package repository

import (
	"context"
	"fmt"
	"go-articles-app/models"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(ctx context.Context, category *models.Category) error {
	result := r.db.WithContext(ctx).Create(category)

	return result.Error
}

func (r *CategoryRepository) GetAll(ctx context.Context) ([]models.Category, error) {
	var categories []models.Category
	result := r.db.WithContext(ctx).Order("created_at DESC").Find(&categories)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to delete comments: %w", result.Error)
	}

	return categories, nil
}

func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*models.Category, error) {
	var category models.Category
	result := r.db.WithContext(ctx).Where("slug = ?", slug).First(&category)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("category not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get category: %w", result.Error)
	}

	return &category, nil
}

type CategoryStats struct {
	CategoryID    uint    `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	ArticlesCount int     `json:"articles_count"`
	TotalViews    int     `json:"total_views"`
	TotalComments int     `json:"total_comments"`
	AvgViews      float64 `json:"avg_views"`
}

func (r *CategoryRepository) GetWithArticlesCount(ctx context.Context) ([]CategoryStats, error) {
	var stats []CategoryStats

	result := r.db.WithContext(ctx).
		Select(`
			categories.id AS category_id,
			categories.name AS category_name,
			COUNT(DISTINCT articles.id) as articles_count, 
			COALESCE(SUM(articles.views), 0) AS total_views,
			COALESCE(COUNT(DISTINCT comments.id), 0) AS total_comments,
			COALESCE(AVG(articles.views), 0) AS avg_views`).
		Joins("LEFT JOIN articles ON articles.category_id = categories.id AND articles.deleted_at IS NULL").
		Joins("LEFT JOIN comments ON comments.article_id = articles.id AND comments.deleted_at IS NULL").
		Group("categories.id, categories.name").
		Order("articles_count DESC").
		Scan(&stats)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get category stats: %w", result.Error)
	}

	return stats, nil
}
