package repository

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"omni-pixel/domain"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Model(&domain.User{}).Where("lower(email) = lower(?)", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		}).Error
}
