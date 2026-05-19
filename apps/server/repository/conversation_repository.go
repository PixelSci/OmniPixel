package repository

import (
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
