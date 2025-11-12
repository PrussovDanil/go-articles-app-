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

}
