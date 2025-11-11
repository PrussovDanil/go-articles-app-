package repository

import (
	"context"
	"database/sql"
	"fmt"
	"go-articles-app/models"
	"time"
)

type ArticleRepository struct {
	db *sql.DB
}

func NewArticleRepository(db *sql.DB) *ArticleRepository {
	return &ArticleRepository{db: db}
}

func (r *ArticleRepository) Create(ctx context.Context, article *models.Article) error {
	query := `
		INSERT INTO articles (title, content, author_id)
		VALUES ($1, $2, $3)
		RETURNING id, published, views, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		article.Title,
		article.Content,
		article.AuthorID,
		article.Published,
		0,
		article.CreatedAt,
		article.UpdatedAt,
	).Scan(&article.ID,
		&article.Published,
		&article.Views,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create article: %w", err)
	}
	return nil
}

func (r *ArticleRepository) GetByID(ctx context.Context, id int) (*models.Article, error) {
	query := `SELECT id, title, content, author_id, published, views, created_at, updated_at FROM articles WHERE id = $1`

	article := &models.Article{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.AuthorID,
		&article.Published,
		&article.Views,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return article, nil
}

func (r *ArticleRepository) GetByAuthorID(ctx context.Context, authorID int) ([]*models.Article, error) {
	query := `
		SELECT id, title, content, author_id, published, views, created_at, updated_at
		FROM articles 
		WHERE author_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, authorID)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.AuthorID,
			&article.Published,
			&article.Views,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to query articles: %w", err)
		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) GetPublished(ctx context.Context) ([]*models.Article, error) {
	query := `
		SELECT id, title, content, author_id, published, views, created_at, updated_at
		FROM articles
		WHERE published = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query articles: %w", err)
	}
	defer rows.Close()

	var articles []*models.Article
	for rows.Next() {
		article := &models.Article{}
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.AuthorID,
			&article.Published,
			&article.Views,
			&article.CreatedAt,
			&article.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan article: %w", err)

		}
		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return articles, nil
}

func (r *ArticleRepository) Update(ctx context.Context, article *models.Article) error {
	query := `UPDATE articles SET 
		title = $1, 
		content = $2,
		author_id = $3,
		published = $4,
		updated_at = $5
		WHERE id = $6
	`

	article.UpdatedAt = time.Now()
	result, err := r.db.ExecContext(
		ctx,
		query,
		article.Title,
		article.Content,
		article.AuthorID,
		article.Published,
		article.UpdatedAt,
		article.ID,
	)

	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *ArticleRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM articles WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *ArticleRepository) Publish(ctx context.Context, id int) error {
	query := `UPDATE articles SET published = true, updated_at = NOW() WHERE id = $1 AND published = false`
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to publish article: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	if rowsAffected == 0 {
		var exists bool
		r.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM articles WHERE id = $1)", id).Scan(&exists)
		if !exists {
			return fmt.Errorf("article not found")
		}
		return nil
	}

	return nil
}

func (r *ArticleRepository) IncrementViews(ctx context.Context, id int) error {
	query := `UPDATE articles SET views = views + 1 WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *ArticleRepository) CreateArticleWithAuthor(
	ctx context.Context,
	userName, userEmail, articleTitle, articleContent string,
) (*models.User, *models.Article, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var user models.User
	queryUser := `SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1`
	err = tx.QueryRowContext(ctx, queryUser, userEmail).Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

	if err == sql.ErrNoRows {

		err = tx.QueryRowContext(
			ctx,
			`INSERT INTO users (email, name, created_at, updated_at)
			VALUES ($1, $2, NOW(), NOW())
			RETURNING id
			`, userEmail, userName,
		).Scan(&user.ID)

		if err != nil {
			return nil, nil, fmt.Errorf("create user: %w", err)
		}

		user.Name = userName
		user.Email = userEmail
	} else if err != nil {

		return nil, nil, fmt.Errorf("check user: %w", err)
	}

	var article models.Article
	err = tx.QueryRowContext(
		ctx,
		`INSERT INTO articles (title, content, author_id, created_at, updated_at)
		  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, title, content, author_id, published, views, created_at, updated_at`,
		articleTitle, articleContent, user.ID).Scan(
		&article.ID,
		&article.Title,
		&article.Content,
		&article.AuthorID,
		&article.Published,
		&article.Views,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		return nil, nil, fmt.Errorf("create article: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("commit: %w", err)
	}

	return &user, &article, nil
}

type ArticleWithAuthor struct {
	Article     *models.Article
	AuthorName  string
	AuthorEmail string
}

// Используй JOIN
func (r *ArticleRepository) GetArticleWithAuthor(ctx context.Context, id int) (*ArticleWithAuthor, error) {
	query := ` SELECT 	
			articles.id, 
			articles.title, 
			articles.content, 
			articles.author_id,
			articles.published,
			articles.views,
			articles.created_at,
			articles.updated_at,
			users.name,
			users.email
	      FROM articles JOIN users ON articles.author_id = users.id
	      WHERE articles.id = $1`

	var result ArticleWithAuthor
	result.Article = &models.Article{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&result.Article.ID,
		&result.Article.Title,
		&result.Article.Content,
		&result.Article.AuthorID,
		&result.Article.Published,
		&result.Article.Views,
		&result.Article.CreatedAt,
		&result.Article.UpdatedAt,
		&result.AuthorName,
		&result.AuthorEmail,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("article not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get article with author: %w", err)
	}

	return &result, nil

}
