-- name: GetFeedFollowsForUser :many
select feed_follows.*, u2.name as user_name, feeds.name as feed_name
from feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
INNER JOIN users u2 ON feeds.user_id = feeds.user_id
WHERE feed_follows.user_id = $1;