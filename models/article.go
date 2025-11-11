package models

import (
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Article struct {
	ID        uint   `gorm:"primaryKey"`
	Title     string `gorm:"not null;size:255;index"`
	Content   string `gorm:"not null;type:text"`
	AuthorID  uint   `gorm:"not null;index"`
	Author    User   `gorm:"foreignKey:AuthorID"`
	Published bool   `gorm:"default:false"`
	Views     int    `gorm:"default:0"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
