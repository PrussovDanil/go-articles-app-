package repository

import (
	"context"
	"fmt"
	"go-articles-app/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Create(user)

	if result.Error != nil {
		if isDuplicateKeyError(result.Error) {
			return fmt.Errorf("user with email %s already exist", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", result.Error)
	}
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).First(&user, id)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user not found")
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, fmt.Errorf("failed to get user: %w", result.Error)
	}

	return &user, nil
}

func (r *UserRepository) GetAll(ctx context.Context) ([]models.User, error) {
	var users []models.User

	result := r.db.WithContext(ctx).Order("id").Find(&users)

	if result.Error != nil {
		return nil, fmt.Errorf("failed to get users: %w", result.Error)
	}

	return users, nil
}

func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	result := r.db.WithContext(ctx).Save(user)

	if result.Error != nil {
		return fmt.Errorf("failed to update user:%w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, id)

	if result.Error != nil {
		return fmt.Errorf("failed to delete user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "ERROR: duplicate key value violates unique constraint" ||
		contains(err.Error(), "duplicate key") ||
		contains(err.Error(), "23505"))

}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr
}
