package usecase

import (
	"encoding/json"
	"strings"

	"github.com/google/uuid"

	"omni-pixel/domain"
)

const defaultSessionTitle = "New Chat"

type SessionUsecase struct {
	sessionRepository domain.SessionRepository
}

func NewSessionUsecase(sessionRepository domain.SessionRepository) *SessionUsecase {
	return &SessionUsecase{sessionRepository: sessionRepository}
}

func (u *SessionUsecase) ListSessions(userID string) (*domain.SessionListResponse, error) {
	sessions, err := u.sessionRepository.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.SessionListResponse{Sessions: sessions}, nil
}

func (u *SessionUsecase) CreateSession(request domain.CreateSessionRequest) (*domain.Session, error) {
	title := strings.TrimSpace(request.Title)
	if title == "" {
		title = defaultSessionTitle
	}

	session := domain.Session{
		ID:          uuid.NewString(),
		UserID:      request.UserID,
		Title:       title,
		Preview:     strings.TrimSpace(request.Preview),
		Model:       strings.TrimSpace(request.Model),
		ChatContent: []byte("[]"),
	}

	return u.sessionRepository.Create(session)
}

func (u *SessionUsecase) SaveSessionChatContent(request domain.SaveSessionChatContentRequest) error {
	chatContent, err := json.Marshal(request.Messages)
	if err != nil {
		return err
	}

	return u.sessionRepository.SaveChatContent(
		request.SessionID,
		request.UserID,
		chatContent,
		sessionPreview(request.Messages),
	)
}

func sessionPreview(messages []domain.ChatMessage) string {
	for i := len(messages) - 1; i >= 0; i-- {
		content := strings.TrimSpace(messages[i].Content)
		if content != "" {
			return content
		}
	}

	return ""
}
