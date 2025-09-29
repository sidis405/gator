-- +goose Up
CREATE TABLE feeds
(
    id         INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name       TEXT NOT NULL,
    url        TEXT NOT NULL UNIQUE,
    user_id    uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE feeds;