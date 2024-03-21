-- +goose Up

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    hashed_password CHAR(60) NOT NULL
);
ALTER TABLE users ADD CONSTRAINT email_unique UNIQUE (email);

-- +goose Down

DROP TABLE users;

