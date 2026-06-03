package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
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

func (r *UserRepository) ExistsByEmail(email string) error {
	var user domain.User
	err := r.db.Model(&domain.User{}).Where("lower(email) = lower(?)", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return domain.ErrUserAlreadyExists
}

func (r *UserRepository) Create(user *domain.User) error {
	if user.ID == "" {
		user.ID = uuid.New().String()
	}
	if user.DeviceFingerprint == "" {
		user.DeviceFingerprint = uuid.New().String()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return r.db.Create(user).Error
}

func (r *UserRepository) UpdateLastLogin(userID string) error {
	return r.db.Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"last_login_at": time.Now(),
			"updated_at":    time.Now(),
		}).Error
}
