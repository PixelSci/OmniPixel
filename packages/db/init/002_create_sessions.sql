CREATE TABLE IF NOT EXISTS sessions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  preview TEXT,
  model VARCHAR(100),
  chat_content JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  last_chat_at TIMESTAMPTZ,
  CONSTRAINT sessions_title_check CHECK (length(trim(title)) > 0)
);

CREATE INDEX IF NOT EXISTS sessions_user_updated_idx
  ON sessions (user_id, updated_at DESC);
