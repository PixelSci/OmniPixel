package domain

import (
	"time"

	"github.com/google/uuid"
)

type Model struct {
	ID         uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid;column:id"`
	ProviderID uuid.UUID  `json:"provider_id" gorm:"type:uuid;column:provider_id"`
	ModelName  string     `json:"model_name" gorm:"column:model_name"`
	IsEnabled  bool       `json:"is_enabled" gorm:"column:is_enabled;default:true"`
	ExpireTime *time.Time `json:"expire_time" gorm:"column:expire_time"`
	CreatedAt  time.Time  `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"column:updated_at"`

	Provider *Provider `json:"provider,omitempty" gorm:"foreignKey:ProviderID"`
}
