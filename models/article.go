package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gosimple/slug"
	"gorm.io/gorm"
)

type Article struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Title         string         `gorm:"not null;size:255;index" json:"title"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Slug          string         `gorm:"uniqueIndex;not null" json:"slug"`
	FeaturedImage string         `gorm:"size:500" json:"featured_image,omitempty"`
	Published     bool           `gorm:"default:false;index" json:"published"`
	Views         uint           `gorm:"default:0" json:"views"`
	AuthorID      uint           `gorm:"not null;index" json:"author_id"`
	CategoryID    *uint          `gorm:"index" json:"category_id,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Author        User           `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE" json:"author,omitempty"`
	Category      *Category      `gorm:"foreignKey:CategoryID;constraint:OnDelete:SET NULL" json:"category,omitempty"`
	Comments      []Comment      `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}

func (a *Article) BeforeCreate(tx *gorm.DB) error {

	if len(a.Title) < 3 {
		return fmt.Errorf("title must be least 3 characters")
	}

	if len(a.Content) < 10 {
		return fmt.Errorf("content must be at least 10 characters")
	}
	a.Slug = slug.Make(a.Title)
	a.Title = strings.TrimSpace(a.Title)
	a.Content = strings.TrimSpace(a.Content)

	return nil
}

func (a *Article) AfterCreate(tx *gorm.DB) error {
	log.Printf("✅ New article created: %s (ID: %d)", a.Title, a.ID)

	return nil
}

func (a *Article) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("AuthorID") {
		return fmt.Errorf("cannot change article author")
	}

	return nil
}

func (a *Article) BeforeDelete(tx *gorm.DB) error {
	log.Printf("⚠️  Deleting article: %s (ID: %d)", a.Title, a.ID)
	return nil
}

func (a *Article) AfterFind(tx *gorm.DB) error {
	log.Printf("📖 Article viewed: %s", a.Title)
	return nil
}
