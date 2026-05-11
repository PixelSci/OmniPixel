package usecase

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"

	"omni-pixel/domain"
)

const defaultSessionTitle = "New Chat"
const sessionTitleMaxLength = 60

type SessionUsecase struct {
	sessionRepository domain.SessionRepository
	chatClient        domain.ChatCompletionClient
}

func NewSessionUsecase(sessionRepository domain.SessionRepository, chatClient domain.ChatCompletionClient) *SessionUsecase {
	return &SessionUsecase{sessionRepository: sessionRepository, chatClient: chatClient}
}

func (u *SessionUsecase) ListSessions(userID string) (*domain.SessionListResponse, error) {
	sessions, err := u.sessionRepository.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.SessionListResponse{Sessions: sessions}, nil
}

func (u *SessionUsecase) GetSession(sessionID string, userID string) (*domain.SessionDetailResponse, error) {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return nil, domain.ErrSessionNotFound
	}

	session, err := u.sessionRepository.FindByID(sessionID, userID)
	if err != nil {
		return nil, err
	}

	var messages []domain.ChatMessage
	if len(session.ChatContent) > 0 {
		if err := json.Unmarshal(session.ChatContent, &messages); err != nil {
			return nil, err
		}
	}
	if messages == nil {
		messages = make([]domain.ChatMessage, 0)
	}

	return &domain.SessionDetailResponse{
		ID:        session.ID,
		Title:     session.Title,
		Preview:   session.Preview,
		Model:     session.Model,
		Messages:  messages,
		CreatedAt: session.CreatedAt,
		UpdatedAt: session.UpdatedAt,
	}, nil
}

func (u *SessionUsecase) CreateSession(request domain.CreateSessionRequest) (*domain.Session, error) {
	sessionID := strings.TrimSpace(request.ID)
	if sessionID == "" {
		sessionID = uuid.NewString()
	} else if _, err := uuid.Parse(sessionID); err != nil {
		return nil, domain.ErrInvalidSessionID
	}

	title := strings.TrimSpace(request.Title)
	if title == "" {
		title = defaultSessionTitle
	}

	session := domain.Session{
		ID:          sessionID,
		UserID:      request.UserID,
		Title:       title,
		Preview:     strings.TrimSpace(request.Preview),
		Model:       strings.TrimSpace(request.Model),
		ChatContent: []byte("[]"),
	}

	return u.sessionRepository.Create(session)
}

func (u *SessionUsecase) SendSessionPrompt(request domain.SendSessionPromptRequest) (*domain.SendSessionPromptResponse, error) {
	prompt := strings.TrimSpace(request.Prompt)
	if prompt == "" {
		return nil, domain.ErrInvalidPrompt
	}

	sessionID := strings.TrimSpace(request.SessionID)
	if sessionID != "" {
		if _, err := uuid.Parse(sessionID); err != nil {
			return nil, domain.ErrInvalidSessionID
		}
	}

	messages := make([]domain.ChatMessage, 0)
	createdSession := false

	if sessionID != "" {
		session, err := u.sessionRepository.FindByID(sessionID, request.UserID)
		if err != nil && !errors.Is(err, domain.ErrSessionNotFound) {
			return nil, err
		}
		if session != nil {
			messages, err = sessionMessages(session)
			if err != nil {
				return nil, err
			}
		} else {
			created, err := u.CreateSession(domain.CreateSessionRequest{
				ID:      sessionID,
				UserID:  request.UserID,
				Title:   titleFromPrompt(prompt),
				Preview: prompt,
				Model:   strings.TrimSpace(request.Model),
			})
			if err != nil {
				return nil, err
			}
			sessionID = created.ID
			createdSession = true
		}
	} else {
		created, err := u.CreateSession(domain.CreateSessionRequest{
			UserID:  request.UserID,
			Title:   titleFromPrompt(prompt),
			Preview: prompt,
			Model:   strings.TrimSpace(request.Model),
		})
		if err != nil {
			return nil, err
		}
		sessionID = created.ID
		createdSession = true
	}

	userMessage := domain.ChatMessage{
		ID:      uuid.NewString(),
		Role:    "user",
		Content: prompt,
	}
	messages = append(messages, userMessage)

	answer, err := u.chatClient.Complete(domain.ChatCompletionRequest{
		Provider: strings.TrimSpace(request.Provider),
		Model:    strings.TrimSpace(request.Model),
		APIKey:   strings.TrimSpace(request.APIKey),
		Messages: messages,
	})
	if err != nil {
		return nil, err
	}

	assistantMessage := domain.ChatMessage{
		ID:      uuid.NewString(),
		Role:    "assistant",
		Content: answer,
		Model:   strings.TrimSpace(request.Model),
	}
	messages = append(messages, assistantMessage)

	if err := u.SaveSessionChatContent(domain.SaveSessionChatContentRequest{
		SessionID: sessionID,
		UserID:    request.UserID,
		Messages:  messages,
	}); err != nil {
		return nil, err
	}

	return &domain.SendSessionPromptResponse{
		SessionID:         sessionID,
		CreatedSession:   createdSession,
		Message:          userMessage,
		AssistantMessage: assistantMessage,
		Messages:         messages,
	}, nil
}

func (u *SessionUsecase) SaveSessionChatContent(request domain.SaveSessionChatContentRequest) error {
	if request.Messages == nil {
		request.Messages = make([]domain.ChatMessage, 0)
	}

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

func (u *SessionUsecase) DeleteSession(sessionID string, userID string) error {
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return domain.ErrSessionNotFound
	}

	return u.sessionRepository.Delete(sessionID, userID)
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

func sessionMessages(session *domain.Session) ([]domain.ChatMessage, error) {
	var messages []domain.ChatMessage
	if len(session.ChatContent) > 0 {
		if err := json.Unmarshal(session.ChatContent, &messages); err != nil {
			return nil, err
		}
	}
	if messages == nil {
		return make([]domain.ChatMessage, 0), nil
	}

	return messages, nil
}

func titleFromPrompt(prompt string) string {
	title := strings.TrimSpace(prompt)
	if title == "" {
		return defaultSessionTitle
	}
	if len([]rune(title)) <= sessionTitleMaxLength {
		return title
	}

	return string([]rune(title)[:sessionTitleMaxLength])
}
