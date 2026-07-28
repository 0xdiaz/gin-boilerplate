-- Rollback auth token columns
-- Migration: 000002_add_user_auth_tokens

DROP INDEX IF EXISTS idx_users_password_reset_token;
DROP INDEX IF EXISTS idx_users_refresh_token;

ALTER TABLE users DROP COLUMN IF EXISTS password_reset_expiry;
ALTER TABLE users DROP COLUMN IF EXISTS password_reset_token;
ALTER TABLE users DROP COLUMN IF EXISTS refresh_token;
