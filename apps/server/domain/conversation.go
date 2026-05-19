package domain

import (
	"time"

	"github.com/google/uuid"
)

// var (
//
//	ErrConversationNotFound  = errors.New("conversation not found")
//	ErrInvalidConversationID = errors.New("invalid conversation id")
//	ErrInvalidPrompt         = errors.New("invalid prompt")
//	ErrInvalidAIConfig       = errors.New("invalid ai config")
//	ErrUnsupportedAIProvider = errors.New("unsupported ai provider")
//
// )

type Conversation struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;column:id"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	Title      string    `json:"title" gorm:"column:title"`
	IsVisible  bool      `json:"is_visible" gorm:"column:is_visible;default:true"`
	IsArchived bool      `json:"is_archived" gorm:"column:is_archived;default:false"`
	CreatedAt  time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at"`
}

type Message struct {
	ID             uuid.UUID `json:"id" gorm:"primaryKey;type:uuid;column:id"`
	ConversationID uuid.UUID `json:"conversation_id" gorm:"type:uuid;column:conversation_id"`
	UserId         uuid.UUID `json:"user_id" gorm:"type:uuid;column:user_id"`
	Content        string    `json:"content" gorm:"column:content"`
	ModelID        uuid.UUID `json:"model_id" gorm:"type:uuid;column:model_id"`
	Type           uint8     `json:"type" gorm:"column:type"`
	CreatedAt      time.Time `json:"created_at" gorm:"column:created_at"`
}

type ConversationRepository interface {
	ListByUserID(userID uuid.UUID) ([]Conversation, error)
}
