CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY,
  username VARCHAR(50) NOT NULL UNIQUE,
  email VARCHAR(255) NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  device_fingerprint VARCHAR(255) NOT NULL,
  display_name VARCHAR(100),
  avatar_url TEXT,
  status VARCHAR(20) NOT NULL DEFAULT 'active',
  email_verified_at TIMESTAMPTZ,
  last_login_at TIMESTAMPTZ,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT users_device_fingerprint_check CHECK (length(trim(device_fingerprint)) > 0),
  CONSTRAINT users_status_check CHECK (status IN ('active', 'inactive', 'suspended'))
);
