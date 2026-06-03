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
		Order("updated_at DESC").
		Find(&conversations).Error
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

func (r *ConversationRepository) Insert(conversation *domain.Conversation) error {
	return r.db.Create(conversation).Error
}

func (r *ConversationRepository) InsertMessage(message *domain.Message) error {
	return r.db.Create(message).Error
}

func (r *ConversationRepository) UpdateTitle(conversationID, userID uuid.UUID, title string) error {
	result := r.db.
		Model(&domain.Conversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Update("title", title)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrConversationNotFound
	}
	return nil
}

func (r *ConversationRepository) Delete(conversationID, userID uuid.UUID) error {
	result := r.db.
		Where("id = ? AND user_id = ?", conversationID, userID).
		Delete(&domain.Conversation{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrConversationNotFound
	}
	return nil
}
