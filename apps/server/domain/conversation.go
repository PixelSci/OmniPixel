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
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	Title      string    `json:"title"`
	IsVisible  bool      `json:"is_visible" comment:"是否可见(软删)"`
	IsArchived bool      `json:"is_archived" comment:"是否归档"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Message struct {
	ID             uuid.UUID `json:"id"`
	ConversationID uuid.UUID `json:"conversation_id"`
	UserId         uuid.UUID `json:"user_id"`
	Content        string    `json:"content"`
	ModelID        uuid.UUID `json:"model_id"`
	Type           uint8     `json:"type" comment:"0: question, 1: reply"`
	CreatedAt      time.Time `json:"created_at"`
}

type ConversationRepository interface {
	ListByUserID(userID uuid.UUID) ([]Conversation, error)
}
