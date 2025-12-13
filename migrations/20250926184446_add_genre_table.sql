-- +goose Up
-- +goose StatementBegin
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    slug VARCHAR(40) UNIQUE NOT NULL,
    name VARCHAR(40),
    created_at TIMESTAMP with time zone DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE genres;
-- +goose StatementEnd
