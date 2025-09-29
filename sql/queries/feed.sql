-- name: GetFeed :one
SELECT * FROM feeds where url = $1;