package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Content   string         `gorm:"type:text;not null" json:"content"`
	ArticleID uint           `gorm:"not null;index" json:"article_id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	ParentID  *uint          `gorm:"index" json:"parent_id,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Article   Article        `gorm:"foreignKey:ArticleID;constraint:OnDelete:CASCADE" json:"article,omitempty"`
	User      User           `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Parent    *Comment       `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"parent,omitempty"`
	Replies   []Comment      `gorm:"foreignKey:ParentID" json:"replies,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

func (a *Comment) BeforeCreate(tx *gorm.DB) error {

	if len(a.Content) < 3 {
		return fmt.Errorf("content must be at least 3 characters")
	}

	return nil
}
