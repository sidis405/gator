-- name: CreatePost :exec
INSERT INTO posts
    (title, url, description, feed_id, published_at, created_at, updated_at)
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7
       ) RETURNING *;