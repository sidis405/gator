-- name: GetNextFeedToFetch :one
select * from feeds order by last_fetched_at ASC NULLS FIRST LIMIT 1;