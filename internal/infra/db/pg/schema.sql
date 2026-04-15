CREATE TABLE users(
    id serial,
    name varchar(100) NOT NULL,
    password varchar(512) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (name)
);
