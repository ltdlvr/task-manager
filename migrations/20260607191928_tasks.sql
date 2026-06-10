-- +goose Up
CREATE TABLE tasks(
    id serial PRIMARY KEY,
    column_id int NOT NULL REFERENCES columns(id) ON DELETE CASCADE,
    title varchar(255) NOT NULL,
    "description" text NOT NULL DEFAULT '',
    position int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE tasks;
