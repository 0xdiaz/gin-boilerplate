-- Create users table for authentication
-- Migration: 000001_create_users_table
-- Created: 2025-11-08

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on deleted_at for soft delete queries
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Add comments
COMMENT ON TABLE users IS 'User accounts for authentication';
COMMENT ON COLUMN users.password IS 'Bcrypt hashed password';
COMMENT ON COLUMN users.deleted_at IS 'Soft delete timestamp';
