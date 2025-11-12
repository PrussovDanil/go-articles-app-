package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Email     string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Name      string         `gorm:"not null;size:255"  json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Articles  []Article      `gorm:"foreignKey:AuthorID constraint:OnDelete:CASCADE" json:"articles,omitempty"`
	Comments  []Comment      `gorm:"foreignKey:UserID constraint:OnDelete:CASCADE" json:"comments,omitempty"`
}

func (User) TableName() string {
	return "users"
}
