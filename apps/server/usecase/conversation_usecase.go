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

func (u *ConversationUseCase) GetConversation(userID, conversationID uuid.UUID) (*domain.ConversationDetailResponse, error) {
	conversation, err := u.conversationRepository.FindByID(conversationID, userID)
	if err != nil {
		return nil, err
	}

	messages, err := u.conversationRepository.ListMessagesByConversationID(conversationID)
	if err != nil {
		return nil, err
	}

	return &domain.ConversationDetailResponse{
		ID:         conversation.ID,
		Title:      conversation.Title,
		IsVisible:  conversation.IsVisible,
		IsArchived: conversation.IsArchived,
		CreatedAt:  conversation.CreatedAt,
		UpdatedAt:  conversation.UpdatedAt,
		Messages:   messages,
	}, nil
}
