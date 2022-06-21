-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id serial PRIMARY KEY,
    name text NOT NULL,
    description text NOT NULL
);

INSERT INTO roles (name, description) 
VALUES ('admin', 'Administrator'), ('user', 'users');

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    username text NOT NULL,
    email text NOT NULL,
    password text NOT NULL,
    role_id integer NOT NULL REFERENCES roles(id),
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS roles CASCADE;
-- +goose StatementEnd
