-- +goose Up
CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(250) NOT NULL
);



-- +goose Down
DROP TABLE users;
