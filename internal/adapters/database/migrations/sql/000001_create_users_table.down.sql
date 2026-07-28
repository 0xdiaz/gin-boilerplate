-- Rollback users table creation
-- Migration: 000001_create_users_table
-- This will drop the users table and all associated data

DROP INDEX IF EXISTS idx_users_deleted_at;
DROP INDEX IF EXISTS idx_users_email;
DROP TABLE IF EXISTS users;
