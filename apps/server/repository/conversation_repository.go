package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"omni-pixel/domain"
)

type ConversationRepository struct {
	db *gorm.DB
}

func NewConversationRepository(db *gorm.DB) *ConversationRepository {
	return &ConversationRepository{db: db}
}

func (r *ConversationRepository) ListByUserID(userID uuid.UUID) ([]domain.Conversation, error) {
	var conversations []domain.Conversation
	err := r.db.
		Where("user_id = ?", userID).
		Order("COALESCE(last_chat_at, updated_at) DESC").
		Find(&domain.Conversation{}).Error
	if err != nil {
		return nil, err
	}
	return conversations, nil
}

func (r *ConversationRepository) FindByID(conversationID, userID uuid.UUID) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.db.
		Where("id = ? AND user_id = ? AND is_visible = ?", conversationID, userID, true).
		First(&conversation).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, domain.ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}
	return &conversation, nil
}

func (r *ConversationRepository) ListMessagesByConversationID(conversationID uuid.UUID) ([]domain.Message, error) {
	var messages []domain.Message
	err := r.db.
		Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error
	if err != nil {
		return nil, err
	}
	return messages, nil
}
