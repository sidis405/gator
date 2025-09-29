-- name: GetPostsForUser :many
SELECT * FROM posts
WHERE feed_id IN (
    select feed_id from feed_follows where user_id = $1
    ) ORDER BY published_at DESC;