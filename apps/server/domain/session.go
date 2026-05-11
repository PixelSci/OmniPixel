package domain

import (
	"errors"
	"time"
)

var (
	ErrSessionNotFound      = errors.New("session not found")
	ErrInvalidSessionID     = errors.New("invalid session id")
	ErrInvalidPrompt        = errors.New("invalid prompt")
	ErrInvalidAIConfig      = errors.New("invalid ai config")
	ErrUnsupportedAIProvider = errors.New("unsupported ai provider")
)

type ChatMessage struct {
	ID      string `json:"id"`
	Role    string `json:"role"`
	Content string `json:"content"`
	Model   string `json:"model,omitempty"`
}

type Session struct {
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

type SessionListItem struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Preview   string     `json:"preview,omitempty"`
	Model     string     `json:"model,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
	LastChatAt *time.Time `json:"last_chat_at,omitempty"`
}

type SessionListResponse struct {
	Sessions []SessionListItem `json:"sessions"`
}

type SessionDetailResponse struct {
	ID        string        `json:"id"`
	Title     string        `json:"title"`
	Preview   string        `json:"preview,omitempty"`
	Model     string        `json:"model,omitempty"`
	Messages  []ChatMessage `json:"messages"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type CreateSessionRequest struct {
	ID      string `json:"session_id"`
	UserID  string `json:"-"`
	Title   string `json:"title"`
	Preview string `json:"preview"`
	Model   string `json:"model"`
}

type SaveSessionChatContentRequest struct {
	SessionID string        `json:"-"`
	UserID    string        `json:"-"`
	Messages  []ChatMessage `json:"messages"`
}

type SendSessionPromptRequest struct {
	SessionID string `json:"session_id"`
	UserID    string `json:"-"`
	Prompt    string `json:"prompt"`
	Provider  string `json:"provider"`
	Model     string `json:"model"`
	APIKey    string `json:"api_key"`
}

type SendSessionPromptResponse struct {
	SessionID        string        `json:"session_id"`
	CreatedSession   bool          `json:"created_session"`
	Message          ChatMessage   `json:"message"`
	AssistantMessage ChatMessage   `json:"assistant_message"`
	Messages         []ChatMessage `json:"messages"`
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

type SessionRepository interface {
	ListByUserID(userID string) ([]SessionListItem, error)
	FindByID(sessionID string, userID string) (*Session, error)
	Create(session Session) (*Session, error)
	SaveChatContent(sessionID string, userID string, chatContent []byte, preview string) error
	Delete(sessionID string, userID string) error
}
