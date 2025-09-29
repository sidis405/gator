-- name: GetUser :one
SELECT * FROM users where name = $1;