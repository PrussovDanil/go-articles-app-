package repository

import (
	"context"
	"fmt"
	"go-articles-app/models"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(ctx context.Context, article *models.Article) error {
	result := r.db.WithContext(ctx).Create(article)
	return result.Error
}

func (r *ArticleRepository) GetByID(ctx context.Context, id uint) (*models.Article, error) {
	var article models.Article

	result := r.db.WithContext(ctx).Preload("Author").First(&article, id)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("article not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get article: %w", result.Error)
	}

	return &article, nil
}

func (r *ArticleRepository) GetByAuthorID(ctx context.Context, authorID uint) ([]models.Article, error) {
	var articles []models.Article

	result := r.db.WithContext(ctx).Where("author_id = ?", authorID).Order("created_at DESC").Find(&articles)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get articles: %w", result.Error)
	}

	return articles, nil
}

func (r *ArticleRepository) GetPublished(ctx context.Context) ([]models.Article, error) {
	var articles []models.Article

	result := r.db.WithContext(ctx).
		Where("publish = ?", true).
		Preload("Author").
		Order("created_at DESC").
		Find(&articles)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get published articles: %w", result.Error)
	}
	return articles, nil
}

func (r *ArticleRepository) Update(ctx context.Context, article *models.Article) error {
	result := r.db.WithContext(ctx).Save(article)

	if result.Error != nil {
		return fmt.Errorf("failed to update article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.Article{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *ArticleRepository) Publish(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Article{}).
		Where("id = ? AND published = &", id, false).
		Update("published", true)

	if result.Error != nil {
		return fmt.Errorf("failed to publish article: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("article not found or already published")
	}

	return nil
}

func (r *ArticleRepository) IncrementViews(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).
		Model(&models.Article{}).
		Where("id = ?", id).
		UpdateColumn("views", gorm.Expr("views + 1"))

	if result.Error != nil {
		return fmt.Errorf("failed to increment views: %w", result.Error)
	}

	return nil
}

// func (r *ArticleRepository) CreateArticleWithAuthor(
// 	ctx context.Context,
// 	userName, userEmail, articleTitle, articleContent string,
// ) (*models.User, *models.Article, error) {

// 	var user models.User
// 	var article models.Article

// 	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
// 		result := tx.Where("email = ?", userEmail).First(&user)

// 		if result.Error == gorm.ErrRecordNotFound {
// 			user = models.User{
// 				Email: userEmail,
// 				Name:  userName,
// 			}
// 			if err := tx.Create(&user).Error; err != nil {
// 				return fmt.Errorf("create user: %w", err)
// 			}
// 		} else if result.Error != nil {
// 			// Другая ошибка при поиске
// 			return fmt.Errorf("check user: %w", result.Error)
// 		}

// 		article = models.Article{
// 			Title:     articleTitle,
// 			Content:   articleContent,
// 			AuthorID:  user.ID,
// 			Published: false,
// 			Views:     0,
// 		}
// 		if err := tx.Create(&article).Error; err != nil {
// 			return fmt.Errorf("create article: %w", err)
// 		}
// 		return nil
// 	})
// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return &user, &article, nil
// }
