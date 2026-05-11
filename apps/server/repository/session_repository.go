package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"omni-pixel/domain"
)

type SessionRepository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func NewSessionRepository(db *pgxpool.Pool, timeout time.Duration) *SessionRepository {
	return &SessionRepository{db: db, timeout: timeout}
}

func (r *SessionRepository) ListByUserID(userID string) ([]domain.SessionListItem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		SELECT id, title, preview, model, updated_at, last_chat_at
		FROM sessions
		WHERE user_id = $1
		ORDER BY COALESCE(last_chat_at, updated_at) DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]domain.SessionListItem, 0)
	for rows.Next() {
		var session domain.SessionListItem
		if err := rows.Scan(
			&session.ID,
			&session.Title,
			&session.Preview,
			&session.Model,
			&session.UpdatedAt,
			&session.LastChatAt,
		); err != nil {
			return nil, err
		}
		sessions = append(sessions, session)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sessions, nil
}

func (r *SessionRepository) Create(session domain.Session) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	query := `
		INSERT INTO sessions (id, user_id, title, preview, model, chat_content)
		VALUES ($1, $2, $3, $4, $5, $6::jsonb)
		RETURNING id, user_id, title, preview, model, chat_content, created_at, updated_at, last_chat_at
	`

	var created domain.Session
	err := r.db.QueryRow(
		ctx,
		query,
		session.ID,
		session.UserID,
		session.Title,
		session.Preview,
		session.Model,
		string(session.ChatContent),
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

func (r *SessionRepository) SaveChatContent(sessionID string, userID string, chatContent []byte, preview string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	commandTag, err := r.db.Exec(
		ctx,
		`
			UPDATE sessions
			SET chat_content = $3::jsonb, preview = $4, last_chat_at = now(), updated_at = now()
			WHERE id = $1 AND user_id = $2
		`,
		sessionID,
		userID,
		string(chatContent),
		preview,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrSessionNotFound
	}

	return nil
}

func (r *SessionRepository) Delete(sessionID string, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	commandTag, err := r.db.Exec(
		ctx,
		`DELETE FROM sessions WHERE id = $1 AND user_id = $2`,
		sessionID,
		userID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return domain.ErrSessionNotFound
	}

	return nil
}
