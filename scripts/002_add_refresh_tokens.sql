-- Add refresh_tokens table for token refresh functionality

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(500) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_token ON refresh_tokens(token);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- Add avatar_url column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS avatar_url VARCHAR(500);

COMMENT ON TABLE refresh_tokens IS 'Stores refresh tokens for JWT token refresh functionality';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Timestamp when token was revoked (NULL if still valid)';
COMMENT ON COLUMN users.avatar_url IS 'URL to user avatar image';
