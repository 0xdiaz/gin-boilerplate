-- Create examples table
-- Migration: 000003_create_examples_table
-- Created to match internal/domain/models/example_model.go

CREATE TABLE IF NOT EXISTS examples (
    id SERIAL PRIMARY KEY,
    data TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

COMMENT ON TABLE examples IS 'Example CRUD resource for boilerplate demonstration';
