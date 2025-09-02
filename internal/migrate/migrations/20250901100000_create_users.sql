-- +goose Up
CREATE TABLE IF NOT EXISTS users (
     id UUID PRIMARY KEY,
     email TEXT NOT NULL UNIQUE,
     created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS users;
