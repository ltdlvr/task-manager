CREATE TABLE users(
    id serial,
    name varchar(100) NOT NULL,
    password varchar(512) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (name)
);

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

CREATE TABLE board_members (
    board_id int NOT NULL REFERENCES boards(id) ON DELETE CASCADE,
    user_id int NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role varchar(20) NOT NULL,

    PRIMARY KEY (board_id, user_id)
);