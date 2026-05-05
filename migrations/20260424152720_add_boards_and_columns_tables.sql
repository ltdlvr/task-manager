-- +goose Up
CREATE TABLE boards (
    id serial PRIMARY KEY,
    name varchar(100) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE columns (
    id serial PRIMARY KEY,
    board_id int NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    name varchar(100) NOT NULL,
    position int NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE columns, boards;
