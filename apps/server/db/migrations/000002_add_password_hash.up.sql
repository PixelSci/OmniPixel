-- password_hash is NULL for identities created via OAuth (github / google).
ALTER TABLE users ADD COLUMN password_hash TEXT;
