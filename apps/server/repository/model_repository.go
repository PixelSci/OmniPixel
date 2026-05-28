package repository

import (
	"time"

	"omni-pixel/domain"

	"gorm.io/gorm"
)

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) FindAll(filter domain.ModelFilter) ([]domain.Model, error) {
	var models []domain.Model

	query := r.db.Preload("Provider")

	if filter.ProviderID != nil {
		query = query.Where("provider_id = ?", *filter.ProviderID)
	}

	if !filter.IncludeDisabled {
		query = query.Where("is_enabled = ?", true).
			Where("expire_time IS NULL OR expire_time > ?", time.Now())
	}

	err := query.Find(&models).Error
	return models, err
}
