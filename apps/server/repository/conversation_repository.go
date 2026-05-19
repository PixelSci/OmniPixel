package repository

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"omni-pixel/domain"
)

type ConversationRepository struct {
	db      *gorm.DB
	timeout time.Duration
}

func NewConversationRepository(db *gorm.DB, timeout time.Duration) *ConversationRepository {
	return &ConversationRepository{db: db, timeout: timeout}
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
