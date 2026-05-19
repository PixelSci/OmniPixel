package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/domain"
)

type ConversationRepository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewConversationRepository(db *pgxpool.Pool, timeout time.Duration) *ConversationRepository {
	return &ConversationRepository{db: db, timeout: timeout}
}

func (r *ConversationRepository) ListByUserID(userID uuid.UUID) ([]domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		SELECT id, title, preview, model, updated_at, last_chat_at
		FROM conversations
		WHERE user_id = $1
		ORDER BY COALESCE(last_chat_at, updated_at) DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	conversations := make([]domain.Conversation, 0)
	for rows.Next() {
		var conversation domain.Conversation
		if err := rows.Scan(
			&conversation.ID,
			&conversation.Title,
			&conversation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		conversations = append(conversations, conversation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return conversations, nil
}
