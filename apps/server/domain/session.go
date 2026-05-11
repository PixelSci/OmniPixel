package domain

import (
	"errors"
	"time"
)

var ErrSessionNotFound = errors.New("session not found")

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

type CreateSessionRequest struct {
	UserID  string
	Title   string
	Preview string
	Model   string
}

type SaveSessionChatContentRequest struct {
	SessionID string
	UserID    string
	Messages  []ChatMessage
}

type SessionRepository interface {
	ListByUserID(userID string) ([]SessionListItem, error)
	Create(session Session) (*Session, error)
	SaveChatContent(sessionID string, userID string, chatContent []byte, preview string) error
}
