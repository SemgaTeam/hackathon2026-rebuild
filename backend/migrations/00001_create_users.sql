-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
  id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  login             VARCHAR(50) NOT NULL UNIQUE,
  full_name         VARCHAR(150) NOT NULL,
  password_hash     TEXT NOT NULL,
  registration_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_deleted        BOOLEAN NOT NULL DEFAULT FALSE
);

-- Indexes
CREATE INDEX idx_users_login ON users(login);
CREATE INDEX idx_users_registration_date ON users(registration_date);
CREATE INDEX idx_users_is_deleted ON users(is_deleted);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
-- +goose StatementEnd
