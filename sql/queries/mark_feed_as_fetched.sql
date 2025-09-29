-- name: MarkFeedAsFetched :exec
UPDATE feeds
SET updated_at = $2, last_fetched_at = $2 where id = $1;