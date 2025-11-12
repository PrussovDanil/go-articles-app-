package main

import (
	"context"
	"fmt"
	"go-articles-app/db"
	"go-articles-app/models"
	"go-articles-app/repository"
	"log"
	"time"

	"github.com/gosimple/slug"
)

func main() {
	// Загружаем конфигурацию из .env файла
	cfg, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gormDB, err := db.NewGormConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	// ШАГ 1: Запускаем AutoMigrate для всех моделей
	fmt.Println("\n🔄 Running auto-migration...")
	err = gormDB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Article{},
		&models.Comment{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	fmt.Println("✅ Migration completed")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(gormDB)
	categoryRepo := repository.NewCategoryRepository(gormDB)
	articleRepo := repository.NewArticleRepository(gormDB)
	commentRepo := repository.NewCommentRepository(gormDB)

	// Очистка данных
	fmt.Println("\n🗑️  Cleaning old data...")
	gormDB.Exec("DELETE FROM comments")
	gormDB.Exec("DELETE FROM articles")
	gormDB.Exec("DELETE FROM categories")
	gormDB.Exec("DELETE FROM users")
	fmt.Println("✅ Data cleaned")

	// ШАГ 2: Создаём 3 категории
	fmt.Println("\n📁 Creating categories...")
	categories := []*models.Category{
		{Name: "Go", Slug: slug.Make("Go"), Description: "Go programming language"},
		{Name: "PostgreSQL", Slug: slug.Make("PostgreSQL"), Description: "PostgreSQL database"},
		{Name: "DevOps", Slug: slug.Make("DevOps"), Description: "DevOps practices and tools"},
	}

	for _, cat := range categories {
		if err := categoryRepo.Create(ctx, cat); err != nil {
			log.Fatalf("Failed to create category: %v", err)
		}
		fmt.Printf("✅ Category created: %s (ID: %d)\n", cat.Name, cat.ID)
	}

	// ШАГ 3: Создаём 3 пользователей
	fmt.Println("\n👤 Creating users...")
	users := []*models.User{
		{Email: "alice@example.com", Name: "Alice Johnson"},
		{Email: "bob@example.com", Name: "Bob Smith"},
		{Email: "charlie@example.com", Name: "Charlie Brown"},
	}

	for _, user := range users {
		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		fmt.Printf("✅ User created: %s (ID: %d)\n", user.Name, user.ID)
	}

	// ШАГ 4: Создаём 5 статей в разных категориях (slug генерируется автоматически в BeforeCreate)
	fmt.Println("\n📰 Creating articles...")
	articles := []*models.Article{
		{
			Title:      "Getting Started with Go",
			Content:    "Go is a statically typed, compiled programming language designed at Google. It's known for its simplicity and efficiency.",
			AuthorID:   users[0].ID, // Alice
			CategoryID: &categories[0].ID, // Go
		},
		{
			Title:      "Advanced Go Patterns",
			Content:    "Learn about advanced design patterns in Go including interfaces, composition, and concurrency patterns.",
			AuthorID:   users[1].ID, // Bob
			CategoryID: &categories[0].ID, // Go
		},
		{
			Title:      "PostgreSQL Performance Tuning",
			Content:    "Optimize your PostgreSQL database with indexes, query optimization, and configuration tuning.",
			AuthorID:   users[0].ID, // Alice
			CategoryID: &categories[1].ID, // PostgreSQL
		},
		{
			Title:      "Docker and Kubernetes for Beginners",
			Content:    "Learn containerization with Docker and orchestration with Kubernetes in this comprehensive guide.",
			AuthorID:   users[2].ID, // Charlie
			CategoryID: &categories[2].ID, // DevOps
		},
		{
			Title:      "CI/CD Best Practices",
			Content:    "Continuous Integration and Continuous Deployment best practices for modern software development.",
			AuthorID:   users[1].ID, // Bob
			CategoryID: &categories[2].ID, // DevOps
		},
	}

	for _, article := range articles {
		if err := articleRepo.Create(ctx, article); err != nil {
			log.Fatalf("Failed to create article: %v", err)
		}
		fmt.Printf("✅ Article created: %s (ID: %d, Slug: %s)\n", article.Title, article.ID, article.Slug)
	}

	// ШАГ 5: Создаём 10 комментариев к статьям (включая вложенные)
	fmt.Println("\n💬 Creating comments...")

	// Основные комментарии
	comment1 := &models.Comment{
		Content:   "Great introduction! Very helpful for beginners.",
		ArticleID: articles[0].ID, // Getting Started with Go
		UserID:    users[1].ID,    // Bob
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment1); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment1.Content[:20], users[1].Name)

	comment2 := &models.Comment{
		Content:   "Thanks for sharing! I learned a lot.",
		ArticleID: articles[0].ID, // Getting Started with Go
		UserID:    users[2].ID,    // Charlie
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment2); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment2.Content[:20], users[2].Name)

	// Вложенный комментарий (ответ на comment1)
	comment3 := &models.Comment{
		Content:   "I agree! The examples are very clear.",
		ArticleID: articles[0].ID, // Getting Started with Go
		UserID:    users[2].ID,    // Charlie
		ParentID:  &comment1.ID,   // Ответ на комментарий Bob
	}
	if err := commentRepo.Create(ctx, comment3); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Reply created: '%s...' by %s (reply to comment #%d)\n", comment3.Content[:20], users[2].Name, comment1.ID)

	comment4 := &models.Comment{
		Content:   "Could you explain more about interfaces?",
		ArticleID: articles[1].ID, // Advanced Go Patterns
		UserID:    users[0].ID,    // Alice
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment4); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment4.Content[:20], users[0].Name)

	// Вложенный комментарий (ответ на comment4)
	comment5 := &models.Comment{
		Content:   "Sure! I'll write a follow-up article about interfaces.",
		ArticleID: articles[1].ID, // Advanced Go Patterns
		UserID:    users[1].ID,    // Bob (автор статьи)
		ParentID:  &comment4.ID,   // Ответ на Alice
	}
	if err := commentRepo.Create(ctx, comment5); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Reply created: '%s...' by %s (reply to comment #%d)\n", comment5.Content[:20], users[1].Name, comment4.ID)

	comment6 := &models.Comment{
		Content:   "Excellent performance tips! My queries are much faster now.",
		ArticleID: articles[2].ID, // PostgreSQL Performance Tuning
		UserID:    users[1].ID,    // Bob
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment6); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment6.Content[:20], users[1].Name)

	comment7 := &models.Comment{
		Content:   "Kubernetes was confusing, but this cleared things up!",
		ArticleID: articles[3].ID, // Docker and Kubernetes
		UserID:    users[0].ID,    // Alice
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment7); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment7.Content[:20], users[0].Name)

	// Вложенный комментарий (ответ на comment7)
	comment8 := &models.Comment{
		Content:   "Glad it helped! Check out my next article on Helm.",
		ArticleID: articles[3].ID, // Docker and Kubernetes
		UserID:    users[2].ID,    // Charlie (автор статьи)
		ParentID:  &comment7.ID,   // Ответ на Alice
	}
	if err := commentRepo.Create(ctx, comment8); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Reply created: '%s...' by %s (reply to comment #%d)\n", comment8.Content[:20], users[2].Name, comment7.ID)

	comment9 := &models.Comment{
		Content:   "CI/CD has transformed our development workflow!",
		ArticleID: articles[4].ID, // CI/CD Best Practices
		UserID:    users[2].ID,    // Charlie
		ParentID:  nil,
	}
	if err := commentRepo.Create(ctx, comment9); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Comment created: '%s...' by %s\n", comment9.Content[:20], users[2].Name)

	// Вложенный комментарий второго уровня (ответ на вложенный комментарий)
	comment10 := &models.Comment{
		Content:   "That would be amazing! Looking forward to it.",
		ArticleID: articles[1].ID, // Advanced Go Patterns
		UserID:    users[2].ID,    // Charlie
		ParentID:  &comment5.ID,   // Ответ на ответ Bob
	}
	if err := commentRepo.Create(ctx, comment10); err != nil {
		log.Fatalf("Failed to create comment: %v", err)
	}
	fmt.Printf("✅ Reply created: '%s...' by %s (reply to comment #%d)\n", comment10.Content[:20], users[2].Name, comment5.ID)

	// ШАГ 6: Публикуем 3 статьи
	fmt.Println("\n📢 Publishing articles...")
	articlesToPublish := []int{
		int(articles[0].ID), // Getting Started with Go
		int(articles[2].ID), // PostgreSQL Performance Tuning
		int(articles[3].ID), // Docker and Kubernetes
	}

	for _, articleID := range articlesToPublish {
		if err := articleRepo.Publish(ctx, articleID); err != nil {
			log.Fatalf("Failed to publish article: %v", err)
		}
		// Находим статью в массиве для вывода названия
		for _, article := range articles {
			if article.ID == uint(articleID) {
				fmt.Printf("✅ Article published: %s (ID: %d)\n", article.Title, article.ID)
				break
			}
		}
	}

	// ШАГ 7: Увеличиваем просмотры
	fmt.Println("\n👁️  Incrementing article views...")

	// Getting Started with Go - самая популярная (150 просмотров)
	for i := 0; i < 150; i++ {
		if err := articleRepo.IncrementViews(ctx, int(articles[0].ID)); err != nil {
			log.Fatalf("Failed to increment views: %v", err)
		}
	}
	fmt.Printf("✅ Views incremented: %s → 150 views\n", articles[0].Title)

	// PostgreSQL Performance Tuning - средняя популярность (85 просмотров)
	for i := 0; i < 85; i++ {
		if err := articleRepo.IncrementViews(ctx, int(articles[2].ID)); err != nil {
			log.Fatalf("Failed to increment views: %v", err)
		}
	}
	fmt.Printf("✅ Views incremented: %s → 85 views\n", articles[2].Title)

	// Docker and Kubernetes - высокая популярность (120 просмотров)
	for i := 0; i < 120; i++ {
		if err := articleRepo.IncrementViews(ctx, int(articles[3].ID)); err != nil {
			log.Fatalf("Failed to increment views: %v", err)
		}
	}
	fmt.Printf("✅ Views incremented: %s → 120 views\n", articles[3].Title)

	// Advanced Go Patterns - неопубликованная, но есть просмотры (30)
	for i := 0; i < 30; i++ {
		if err := articleRepo.IncrementViews(ctx, int(articles[1].ID)); err != nil {
			log.Fatalf("Failed to increment views: %v", err)
		}
	}
	fmt.Printf("✅ Views incremented: %s → 30 views (draft)\n", articles[1].Title)

	fmt.Println("\n🎉 Demo setup completed!")
}
