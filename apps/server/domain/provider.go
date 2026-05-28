package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrProviderNotFound = errors.New("provider not found")
)

type Provider struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;column:id"`
	Name      string    `json:"name" gorm:"column:name"`
	BaseURL   string    `json:"base_url" gorm:"column:base_url"`
	APIKey    string    `json:"api_key" gorm:"column:api_key"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type ProviderResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	BaseURL   string    `json:"base_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ProviderRepository interface {
	FindAll() ([]Provider, error)
}
