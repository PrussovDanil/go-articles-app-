package main

import (
	"context"
	"fmt"
	"go-articles-app/db"
	"go-articles-app/models"
	"go-articles-app/repository"
	"log"
	"time"
)

func main() {
	cfg := db.Config{
		Host:     "localhost",
		Port:     5432,
		User:     "gouser",
		Password: "gopass",
		DBName:   "go_article_app",
		SSLMode:  "disable",
	}

	database, err := db.NewConnection(cfg)
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer database.Close()
	fmt.Println("✅ Connected to PostgreSQL")

	userRepo := repository.NewUserRepository(database)
	articleRepo := repository.NewArticleRepository(database)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Очистка таблиц для демонстрации
	_, _ = database.ExecContext(ctx, "DELETE FROM articles")
	_, _ = database.ExecContext(ctx, "DELETE FROM users")

	// 1. Создание пользователей
	fmt.Println("\n📝 Creating users...")

	alice := &models.User{Email: "alice@example.com", Name: "Alice"}
	if err := userRepo.Create(ctx, alice); err != nil {
		log.Fatalf("Failed to create Alice: %v", err)
	}
	fmt.Printf("✅ Created user: %s (%s)\n", alice.Name, alice.Email)

	bob := &models.User{Email: "bob@example.com", Name: "Bob"}
	if err := userRepo.Create(ctx, bob); err != nil {
		log.Fatalf("Failed to create Bob: %v", err)
	}
	fmt.Printf("✅ Created user: %s (%s)\n", bob.Name, bob.Email)

	charlie := &models.User{Email: "charlie@example.com", Name: "Charlie"}
	if err := userRepo.Create(ctx, charlie); err != nil {
		log.Fatalf("Failed to create Charlie: %v", err)
	}
	fmt.Printf("✅ Created user: %s (%s)\n", charlie.Name, charlie.Email)

	// 2. Создание статей
	fmt.Println("\n📰 Creating articles with CreateArticleWithAuthor...")

	userAlice, article1, err := articleRepo.CreateArticleWithAuthor(ctx, "Alice", "alice@example.com", "Introduction to Go", "Go is a statically typed, compiled language...")
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
	fmt.Printf(`✅ Created article "%s" by %s`+"\n", article1.Title, userAlice.Name)

	_, article2, err := articleRepo.CreateArticleWithAuthor(ctx, "Alice", "alice@example.com", "PostgreSQL Basics", "PostgreSQL is a powerful database...")
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
	fmt.Printf(`✅ Created article "%s" by %s`+"\n", article2.Title, "Alice")

	userBob, article3, err := articleRepo.CreateArticleWithAuthor(ctx, "Bob", "bob@example.com", "Web Development in Go", "Building web applications in Go...")
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
	fmt.Printf(`✅ Created article "%s" by %s`+"\n", article3.Title, userBob.Name)

	_, article4, err := articleRepo.CreateArticleWithAuthor(ctx, "Bob", "bob@example.com", "Docker for Beginners", "Docker simplifies deployment...")
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
	fmt.Printf(`✅ Created article "%s" by %s`+"\n", article4.Title, "Bob")

	userDiana, article5, err := articleRepo.CreateArticleWithAuthor(ctx, "Diana", "diana@example.com", "Microservices Architecture", "Microservices pattern explained...")
	if err != nil {
		log.Fatalf("Failed to create article: %v", err)
	}
	fmt.Printf(`✅ Created article "%s" by %s (new user created)`+"\n", article5.Title, userDiana.Name)

	// 3. Публикация статей
	fmt.Println("\n📢 Publishing articles...")

	if err := articleRepo.Publish(ctx, article1.ID); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
	fmt.Printf(`✅ Published: "%s"`+"\n", article1.Title)

	if err := articleRepo.Publish(ctx, article3.ID); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
	fmt.Printf(`✅ Published: "%s"`+"\n", article3.Title)

	if err := articleRepo.Publish(ctx, article4.ID); err != nil {
		log.Fatalf("Failed to publish: %v", err)
	}
	fmt.Printf(`✅ Published: "%s"`+"\n", article4.Title)

	// 4. Увеличение просмотров
	fmt.Println("\n👁️  Incrementing views...")
	for i := 0; i < 5; i++ {
		if err := articleRepo.IncrementViews(ctx, article1.ID); err != nil {
			log.Fatalf("Failed to increment views: %v", err)
		}
	}
	updatedArticle1, _ := articleRepo.GetByID(ctx, article1.ID)
	fmt.Printf(`✅ "%s" views: 0 → %d`+"\n", article1.Title, updatedArticle1.Views)

	// 5. Статистика
	fmt.Println("\n📊 Statistics:")

	allUsers, err := userRepo.GetAll(ctx)
	if err != nil {
		log.Fatalf("Failed to get users: %v", err)
	}
	fmt.Printf("  - Total users: %d\n", len(allUsers))

	// Получаем все статьи (нужно добавить метод GetAll в ArticleRepository)
	allArticlesQuery := `SELECT COUNT(*) FROM articles`
	var totalArticles int
	database.QueryRowContext(ctx, allArticlesQuery).Scan(&totalArticles)
	fmt.Printf("  - Total articles: %d\n", totalArticles)

	publishedArticles, err := articleRepo.GetPublished(ctx)
	if err != nil {
		log.Fatalf("Failed to get published articles: %v", err)
	}
	fmt.Printf("  - Published articles: %d\n", len(publishedArticles))

	// 6. Статьи Alice
	fmt.Println("\n📚 Articles by Alice:")

	aliceArticles, err := articleRepo.GetByAuthorID(ctx, alice.ID)
	if err != nil {
		log.Fatalf("Failed to get Alice's articles: %v", err)
	}
	for i, article := range aliceArticles {
		status := "draft"
		if article.Published {
			status = "published"
		}
		fmt.Printf("  %d. \"%s\" (%s, %d views)\n", i+1, article.Title, status, article.Views)
	}

	// 7. Все опубликованные статьи
	fmt.Println("\n🌐 All published articles:")

	for i, article := range publishedArticles {
		author, _ := userRepo.GetByID(ctx, article.AuthorID)
		authorName := "Unknown"
		if author != nil {
			authorName = author.Name
		}
		fmt.Printf("  %d. \"%s\" by %s (%d views)\n", i+1, article.Title, authorName, article.Views)
	}

	// 8. Обновление статьи
	fmt.Println("\n✏️  Updating article...")

	article2.Title = "Advanced PostgreSQL"
	article2.Content = "Advanced PostgreSQL features and optimization..."
	if err := articleRepo.Update(ctx, article2); err != nil {
		log.Fatalf("Failed to update article: %v", err)
	}
	fmt.Printf(`✅ Updated: "PostgreSQL Basics" → "Advanced PostgreSQL"` + "\n")

	// 9. Удаление пользователя Bob
	fmt.Println("\n🗑️  Deleting user Bob...")

	bobArticles, _ := articleRepo.GetByAuthorID(ctx, bob.ID)
	bobArticlesCount := len(bobArticles)
	if err := userRepo.Delete(ctx, bob.ID); err != nil {
		log.Fatalf("Failed to delete Bob: %v", err)
	}
	fmt.Printf("✅ Deleted user Bob (%d articles deleted automatically via CASCADE)\n", bobArticlesCount)

	// 10. Финальная статистика
	fmt.Println("\n📊 Final statistics:")

	allUsers, _ = userRepo.GetAll(ctx)
	fmt.Printf("  - Total users: %d\n", len(allUsers))

	database.QueryRowContext(ctx, allArticlesQuery).Scan(&totalArticles)
	fmt.Printf("  - Total articles: %d\n", totalArticles)

	publishedArticles, _ = articleRepo.GetPublished(ctx)
	fmt.Printf("  - Published articles: %d\n", len(publishedArticles))

	articleWithAuthor, err := articleRepo.GetArticleWithAuthor(ctx, 43)
	if err != nil {
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("Title: %s\n", articleWithAuthor.Article.Title)
	fmt.Printf("Author: %s (%s)\n", articleWithAuthor.AuthorName, articleWithAuthor.AuthorEmail)
	fmt.Printf("Views: %d\n", articleWithAuthor.Article.Views)

	fmt.Println("\n🎉 All operations completed successfully!")
}
