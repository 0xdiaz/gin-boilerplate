-- Replace single plaintext users.refresh_token column with a dedicated
-- refresh_tokens table: hashed tokens, per-token expiry, and a family_id
-- chain that enables rotation + reuse (theft) detection.
-- Migration: 000004_create_refresh_tokens_table

CREATE TABLE IF NOT EXISTS refresh_tokens (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(64) NOT NULL UNIQUE,
    family_id VARCHAR(64) NOT NULL,
    replaced_by_hash VARCHAR(64),
    revoked_at TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_family_id ON refresh_tokens(family_id);

COMMENT ON TABLE refresh_tokens IS 'Hashed refresh tokens with rotation chain (family_id) for reuse/theft detection';
COMMENT ON COLUMN refresh_tokens.token_hash IS 'SHA-256 hex digest of the raw token; raw token is never stored';
COMMENT ON COLUMN refresh_tokens.family_id IS 'Stable id shared by every token in a rotation chain (one per login session)';
COMMENT ON COLUMN refresh_tokens.replaced_by_hash IS 'token_hash of the token that replaced this one via rotation, if any';
COMMENT ON COLUMN refresh_tokens.revoked_at IS 'Set when rotated, logged out, or force-revoked (reuse detection/logout-all)';

DROP INDEX IF EXISTS idx_users_refresh_token;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token;
