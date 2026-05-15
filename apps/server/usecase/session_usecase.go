package usecase

import (
	"omni-pixel/domain"
)

const defaultSessionTitle = "New Chat"
const sessionTitleMaxLength = 60

type SessionUseCase struct {
	sessionRepository domain.SessionRepository
}

func NewSessionUseCase(sessionRepository domain.SessionRepository) *SessionUseCase {
	return &SessionUseCase{sessionRepository}
}

func (u *SessionUseCase) ListSessions(userID string) (*domain.SessionListResponse, error) {
	sessions, err := u.sessionRepository.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &domain.SessionListResponse{Sessions: sessions}, nil
}
