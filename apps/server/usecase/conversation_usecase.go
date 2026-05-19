package usecase

import (
	"omni-pixel/domain"

	"github.com/google/uuid"
)

const defaultConversationTitle = "New Chat"
const conversationTitleMaxLength = 60

type ConversationUseCase struct {
	conversationRepository domain.ConversationRepository
}

func NewConversationUseCase(conversationRepository domain.ConversationRepository) *ConversationUseCase {
	return &ConversationUseCase{conversationRepository}
}

func (u *ConversationUseCase) ListConversations(userID uuid.UUID) (*[]domain.Conversation, error) {
	conversations, err := u.conversationRepository.ListByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &conversations, nil
}
