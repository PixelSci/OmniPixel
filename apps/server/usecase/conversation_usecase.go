package usecase

import (
	"time"

	"github.com/google/uuid"

	"omni-pixel/domain"
)

const defaultConversationTitle = "New Chat"
const conversationTitleMaxLength = 60

type ConversationUseCase struct {
	conversationRepository domain.ConversationRepository
	aiProvider             domain.AIProvider
}

func NewConversationUseCase(conversationRepository domain.ConversationRepository, aiProvider domain.AIProvider) *ConversationUseCase {
	return &ConversationUseCase{conversationRepository: conversationRepository, aiProvider: aiProvider}
}

func (u *ConversationUseCase) DeleteConversation(userID, conversationID uuid.UUID) error {
	return u.conversationRepository.Delete(conversationID, userID)
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

func (u *ConversationUseCase) Chat(userID uuid.UUID, request domain.ChatRequest, writer domain.StreamWriter) error {
	var conversationID uuid.UUID

	if request.ConversationID == nil {
		title := request.Message
		if len([]rune(title)) > conversationTitleMaxLength {
			title = string([]rune(title)[:conversationTitleMaxLength])
		}
		conversation := &domain.Conversation{
			ID:        uuid.New(),
			UserID:    userID,
			Title:     title,
			IsVisible: true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := u.conversationRepository.Insert(conversation); err != nil {
			return err
		}
		conversationID = conversation.ID
	} else {
		conv, err := u.conversationRepository.FindByID(*request.ConversationID, userID)
		if err != nil {
			return err
		}
		conversationID = conv.ID
	}

	now := time.Now()
	modelUUID, _ := uuid.Parse(request.ModelID)

	userMessage := &domain.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		UserId:         userID,
		Content:        request.Message,
		ModelID:        modelUUID,
		Type:           0,
		CreatedAt:      now,
	}
	if err := u.conversationRepository.InsertMessage(userMessage); err != nil {
		return err
	}

	history, err := u.conversationRepository.ListMessagesByConversationID(conversationID)
	if err != nil {
		return err
	}

	aiMessages := make([]domain.AIChatMessage, 0, len(history)+1)
	for _, m := range history {
		role := "user"
		if m.Type == 1 {
			role = "assistant"
		}
		aiMessages = append(aiMessages, domain.AIChatMessage{Role: role, Content: m.Content})
	}

	ch, err := u.aiProvider.ChatStream(aiMessages, request.ModelID)
	if err != nil {
		return err
	}

	var fullContent string
	for chunk := range ch {
		if chunk.Done {
			break
		}
		fullContent += chunk.Token
		if err := writer.WriteToken(chunk.Token); err != nil {
			return err
		}
	}

	aiMessage := &domain.Message{
		ID:             uuid.New(),
		ConversationID: conversationID,
		UserId:         userID,
		Content:        fullContent,
		ModelID:        modelUUID,
		Type:           1,
		CreatedAt:      time.Now(),
	}
	if err := u.conversationRepository.InsertMessage(aiMessage); err != nil {
		return err
	}

	return writer.WriteDone(conversationID, aiMessage.ID)
}
