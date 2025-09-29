-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
                             user_id, feed_id, created_at, updated_at
        ) VALUES ($1, $2, $3, $4)
           RETURNING *
)

SELECT inserted_feed_follow.*,
       feeds.name as feed_name,
       users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds on feeds.id = inserted_feed_follow.feed_id
INNER JOIN users on users.id = inserted_feed_follow.user_id;