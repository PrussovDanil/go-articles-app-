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
	// Загружаем конфигурацию из .env файла
	cfg, err := db.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gormDB, err := db.NewGormConnection(cfg)

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	fmt.Println("\n🔄 Running auto-migration...")
	err = gormDB.AutoMigrate(&models.User{}, &models.Article{})
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	fmt.Println("✅ Migration completed")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	userRepo := repository.NewUserRepository(gormDB)
	articleRepo := repository.NewArticleRepository(gormDB)

	//Удаление
	gormDB.Exec("DELETE FROM articles")
	gormDB.Exec("DELETE FROM users")

	fmt.Println("\n📝 Creating users...")

	alice := &models.User{Email: "alice@example.com", Name: "Alice"}
	userRepo.Create(ctx, alice)
	fmt.Printf("✅ Created: %s (ID: %d)\n", alice.Name, alice.ID)

	bob := &models.User{Email: "bob@example.com", Name: "Bob"}
	userRepo.Create(ctx, bob)
	fmt.Printf("✅ Created: %s (ID: %d)\n", bob.Name, bob.ID)

	fmt.Println("\n📰 Creating articles...")
	article1 := &models.Article{
		Title:    "Introduction to GORM",
		Content:  "GORM is a fantastic ORM library for Go...",
		AuthorID: alice.ID,
	}
	articleRepo.Create(ctx, article1)
	fmt.Printf("✅ Created article: \"%s\" (ID: %d)\n", article1.Title, article1.ID)

	fmt.Println("\n📖 Reading article with author...")
	foundArticle, _ := articleRepo.GetByID(ctx, article1.ID)
	fmt.Printf("Article: \"%s\"\n", foundArticle.Title)
	fmt.Printf("Author: %s (%s)\n", foundArticle.Author.Name, foundArticle.Author.Email)

	fmt.Println("\n👤 User with articles...")
	var userWithArticles models.User
	gormDB.Preload("Articles").First(&userWithArticles, alice.ID)
	fmt.Printf("User: %s\n", userWithArticles.Name)
	fmt.Printf("Articles: %d\n", len(userWithArticles.Articles))
	for _, a := range userWithArticles.Articles {
		fmt.Printf("  - %s\n", a.Title)
	}

	fmt.Println("\n🎉 GORM demo completed!")

}
