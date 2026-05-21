package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrConversationNotFound  = errors.New("conversation not found")
	ErrInvalidConversationID = errors.New("invalid conversation id")
)

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

type ConversationDetailResponse struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	IsVisible  bool      `json:"is_visible"`
	IsArchived bool      `json:"is_archived"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Messages   []Message `json:"messages"`
}

type AIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AIStreamChunk struct {
	Token string
	Done  bool
}

type ChatRequest struct {
	ConversationID *uuid.UUID `json:"conversation_id"`
	Message        string     `json:"message"`
	ModelID        string     `json:"model_id"`
}

type StreamWriter interface {
	WriteToken(token string) error
	WriteDone(conversationID, messageID uuid.UUID) error
}

type AIProvider interface {
	ChatStream(messages []AIChatMessage, modelID string) (<-chan AIStreamChunk, error)
}

type ConversationRepository interface {
	ListByUserID(userID uuid.UUID) ([]Conversation, error)
	FindByID(conversationID, userID uuid.UUID) (*Conversation, error)
	ListMessagesByConversationID(conversationID uuid.UUID) ([]Message, error)
	Insert(conversation *Conversation) error
	InsertMessage(message *Message) error
}
