-- +goose Up
CREATE TABLE posts
(
    id           INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title        TEXT      NOT NULL,
    url          TEXT      NOT NULL UNIQUE,
    description  TEXT      NOT NULL,
    feed_id      INTEGER   NOT NULL REFERENCES feeds (id) ON DELETE CASCADE,
    published_at TIMESTAMP NOT NULL,
    created_at   TIMESTAMP NOT NULL,
    updated_at   TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE posts;