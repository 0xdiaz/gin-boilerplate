-- Rollback: restore plaintext users.refresh_token column, drop refresh_tokens table
-- Migration: 000004_create_refresh_tokens_table

ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(500);
CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);

DROP TABLE IF EXISTS refresh_tokens;
