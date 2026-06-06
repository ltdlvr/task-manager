-- +goose Up
CREATE TABLE board_members (
    board_id int NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role varchar(20) NOT NULL,

    PRIMARY KEY (board_id, user_id)
);

-- +goose Down
DROP TABLE board_members;
