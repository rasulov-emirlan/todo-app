-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS todos (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title text NOT NULL,
    description text,
    completed boolean NOT NULL DEFAULT false,
    deadline timestamp,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS todos CASCADE;
-- +goose StatementEnd
