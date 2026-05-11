CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  preview TEXT NOT NULL DEFAULT '',
  model VARCHAR(100) NOT NULL DEFAULT '',
  chat_content JSONB NOT NULL DEFAULT '[]'::jsonb,
  last_chat_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT sessions_title_check CHECK (length(trim(title)) > 0),
  CONSTRAINT sessions_chat_content_check CHECK (jsonb_typeof(chat_content) = 'array')
);

CREATE INDEX IF NOT EXISTS sessions_user_recent_idx
  ON sessions (user_id, (COALESCE(last_chat_at, updated_at)) DESC);
