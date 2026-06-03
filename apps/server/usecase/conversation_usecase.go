package usecase

import (
	"time"

	"github.com/google/uuid"

	"omni-pixel/domain"
	"omni-pixel/internal/generation"
)

const defaultConversationTitle = "New Chat"
const conversationTitleMaxLength = 60

type ConversationUseCase struct {
	conversationRepository domain.ConversationRepository
	aiProvider             domain.AIProvider
	genManager             *generation.Manager
}

func NewConversationUseCase(
	conversationRepository domain.ConversationRepository,
	aiProvider domain.AIProvider,
	genManager *generation.Manager,
) *ConversationUseCase {
	return &ConversationUseCase{
		conversationRepository: conversationRepository,
		aiProvider:             aiProvider,
		genManager:             genManager,
	}
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

	resp := &domain.ConversationDetailResponse{
		ID:         conversation.ID,
		Title:      conversation.Title,
		IsVisible:  conversation.IsVisible,
		IsArchived: conversation.IsArchived,
		CreatedAt:  conversation.CreatedAt,
		UpdatedAt:  conversation.UpdatedAt,
		Messages:   messages,
	}

	if gen, ok := u.genManager.Get(conversationID); ok {
		resp.Generating = true
		content := gen.Content()
		if content != "" {
			resp.Messages = append(resp.Messages, domain.Message{
				ID:             uuid.Nil,
				ConversationID: conversationID,
				Content:        content,
				Type:           1,
				CreatedAt:      time.Now(),
			})
		}
	}

	return resp, nil
}

func (u *ConversationUseCase) Chat(userID uuid.UUID, request domain.ChatRequest, writer domain.StreamWriter) error {
	// ── Resume path: existing generation ──────────────────────────
	if request.ConversationID != nil && u.genManager.Has(*request.ConversationID) {
		ch, buffer, ok := u.genManager.Subscribe(*request.ConversationID)
		if !ok {
			// Generation just finished between Has and Subscribe, treat as done
			gen, _ := u.genManager.Get(*request.ConversationID)
			if gen == nil {
				conv, err := u.conversationRepository.FindByID(*request.ConversationID, userID)
				if err != nil {
					return err
				}
				return writer.WriteDone(conv.ID, uuid.Nil)
			}
		}

		// Replay buffered tokens
		for _, token := range buffer {
			if err := writer.WriteToken(token); err != nil {
				u.genManager.Finish(*request.ConversationID) // cleanup subscriber
				return nil
			}
		}

		// Stream live tokens
		for token := range ch {
			if err := writer.WriteToken(token); err != nil {
				return nil // client disconnected, generation continues
			}
		}

		return writer.WriteDone(*request.ConversationID, uuid.Nil)
	}

	// ── New message path ──────────────────────────────────────────
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

	gen := u.genManager.Start(conversationID)
	subCh := gen.Subscribe()

	// Start background goroutine to consume AI stream into Manager
	go func() {
		var fullContent string
		for chunk := range ch {
			if chunk.Done {
				break
			}
			fullContent += chunk.Token
			gen.Append(chunk.Token)
		}

		u.genManager.Finish(conversationID)

		aiMessage := &domain.Message{
			ID:             uuid.New(),
			ConversationID: conversationID,
			UserId:         userID,
			Content:        fullContent,
			ModelID:        modelUUID,
			Type:           1,
			CreatedAt:      time.Now(),
		}
		_ = u.conversationRepository.InsertMessage(aiMessage)
	}()

	// Stream tokens to the SSE client from the Manager's subscriber channel
	for token := range subCh {
		if err := writer.WriteToken(token); err != nil {
			return nil // client disconnected, generation continues in background
		}
	}

	return writer.WriteDone(conversationID, uuid.Nil)
}
