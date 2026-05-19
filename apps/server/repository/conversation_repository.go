package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *ConversationRepository) ListByUserID(userID string) ([]domain.ConversationListItem, error) {
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

	conversations := make([]domain.ConversationListItem, 0)
	for rows.Next() {
		var conversation domain.ConversationListItem
		if err := rows.Scan(
			&conversation.ID,
			&conversation.Title,
			&conversation.Preview,
			&conversation.Model,
			&conversation.UpdatedAt,
			&conversation.LastChatAt,
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

func (r *ConversationRepository) FindByID(conversationID string, userID string) (*domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		SELECT id, user_id, title, preview, model, chat_content, created_at, updated_at, last_chat_at
		FROM conversations
		WHERE id = $1 AND user_id = $2
		LIMIT 1
	`

	var conversation domain.Conversation
	err := r.db.QueryRow(ctx, query, conversationID, userID).Scan(
		&conversation.ID,
		&conversation.UserID,
		&conversation.Title,
		&conversation.Preview,
		&conversation.Model,
		&conversation.ChatContent,
		&conversation.CreatedAt,
		&conversation.UpdatedAt,
		&conversation.LastChatAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrConversationNotFound
	}
	if err != nil {
		return nil, err
	}

	return &conversation, nil
}

func (r *ConversationRepository) Create(conversation domain.Conversation) (*domain.Conversation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		INSERT INTO conversations (id, user_id, title, preview, model, chat_content)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
		RETURNING id, user_id, title, preview, model, chat_content, created_at, updated_at, last_chat_at
	`

	var created domain.Conversation
	err := r.db.QueryRow(
		ctx,
		query,
		conversation.ID,
		conversation.UserID,
		conversation.Title,
		conversation.Preview,
		conversation.Model,
		string(conversation.ChatContent),
	).Scan(
		&created.ID,
		&created.UserID,
		&created.Title,
		&created.Preview,
		&created.Model,
		&created.ChatContent,
		&created.CreatedAt,
		&created.UpdatedAt,
		&created.LastChatAt,
	)
	if err != nil {
		return nil, err
	}

	return &created, nil
}

func (r *ConversationRepository) SaveChatContent(conversationID string, userID string, chatContent []byte, preview string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	commandTag, err := r.db.Exec(
		ctx,
		`
			UPDATE conversations
			SET chat_content = $3::jsonb, preview = $4, last_chat_at = now(), updated_at = now()
			WHERE id = $1 AND user_id = $2
		`,
		conversationID,
		userID,
		string(chatContent),
		preview,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrConversationNotFound
	}

	return nil
}

func (r *ConversationRepository) Delete(conversationID string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	commandTag, err := r.db.Exec(
		ctx,
		`DELETE FROM conversations WHERE id = $1 AND user_id = $2`,
		conversationID,
		userID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrConversationNotFound
	}

	return nil
}
