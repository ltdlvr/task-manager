-- +goose Up
ALTER TABLE users
ADD created_at timestamptz NOT NULL DEFAULT NOW();

-- +goose Down
ALTER TABLE users
DROP COLUMN created_at;
