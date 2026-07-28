-- Add refresh/reset token columns required by the auth flow
-- Migration: 000002_add_user_auth_tokens
-- These columns exist on models.User but were missing from 000001, so
-- AutoMigrate silently created them while the SQL migrations stayed incomplete.

ALTER TABLE users ADD COLUMN IF NOT EXISTS refresh_token VARCHAR(500);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_reset_token VARCHAR(255);
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_reset_expiry TIMESTAMP;

CREATE INDEX IF NOT EXISTS idx_users_refresh_token ON users(refresh_token);
CREATE INDEX IF NOT EXISTS idx_users_password_reset_token ON users(password_reset_token);

COMMENT ON COLUMN users.refresh_token IS 'Active refresh token for JWT renewal';
COMMENT ON COLUMN users.password_reset_token IS 'One-time token for forgot-password flow';
COMMENT ON COLUMN users.password_reset_expiry IS 'Expiry timestamp for password_reset_token';
