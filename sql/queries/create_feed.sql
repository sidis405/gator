-- name: CreateFeed :one
INSERT INTO feeds (name,
                   url,
                   user_id,
                   created_at,
                   updated_at)
VALUES ($1,
        $2,
        $3,
        $4,
        $5)
RETURNING *;