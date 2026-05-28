package repository

import (
	"omni-pixel/domain"

	"gorm.io/gorm"
)

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

func (r *ProviderRepository) FindAll() ([]domain.Provider, error) {
	var providers []domain.Provider
	err := r.db.Find(&providers).Error
	return providers, err
}
