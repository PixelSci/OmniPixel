package domain

import (
	"errors"
	"time"
)

var (
	ErrConversationNotFound      = errors.New("conversation not found")
	ErrInvalidConversationID     = errors.New("invalid conversation id")
	ErrInvalidPrompt             = errors.New("invalid prompt")
	ErrInvalidAIConfig           = errors.New("invalid ai config")
	ErrUnsupportedAIProvider     = errors.New("unsupported ai provider")
)

type ChatMessage struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
}

type Conversation struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Title       string     `json:"title"`
	Preview     string     `json:"preview,omitempty"`
	Model       string     `json:"model,omitempty"`
	ChatContent []byte     `json:"-"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	LastChatAt  *time.Time `json:"last_chat_at,omitempty"`
}

type ConversationListItem struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Preview   string     `json:"preview,omitempty"`
	Model     string     `json:"model,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastChatAt *time.Time `json:"last_chat_at,omitempty"`
}

type ConversationListResponse struct {
	Conversations []ConversationListItem `json:"conversations"`
}

type ConversationDetailResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Preview   string        `json:"preview,omitempty"`
	Model     string        `json:"model,omitempty"`
	Messages  []ChatMessage `json:"messages"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CreateConversationRequest struct {
	ID      string `json:"conversation_id"`
	UserID  string `json:"-"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Model   string `json:"model"`
}

type SaveConversationChatContentRequest struct {
	ConversationID string        `json:"-"`
	UserID         string        `json:"-"`
	Messages       []ChatMessage `json:"messages"`
}

type SendConversationPromptRequest struct {
	ConversationID string `json:"conversation_id"`
	UserID         string `json:"-"`
	Prompt         string `json:"prompt"`
	Provider       string `json:"provider"`
	Model          string `json:"model"`
	APIKey         string `json:"api_key"`
}

type SendConversationPromptResponse struct {
	ConversationID      string        `json:"conversation_id"`
	CreatedConversation bool          `json:"created_conversation"`
	Message             ChatMessage   `json:"message"`
	AssistantMessage    ChatMessage   `json:"assistant_message"`
	Messages            []ChatMessage `json:"messages"`
}

type ChatCompletionRequest struct {
	Provider string
	Model    string
	APIKey   string
	Messages []ChatMessage
}

type ChatCompletionClient interface {
	Complete(request ChatCompletionRequest) (string, error)
}

type ConversationRepository interface {
	ListByUserID(userID string) ([]ConversationListItem, error)
	FindByID(conversationID string, userID string) (*Conversation, error)
	Create(conversation Conversation) (*Conversation, error)
	SaveChatContent(conversationID string, userID string, chatContent []byte, preview string) error
	Delete(conversationID string, userID string) error
}
