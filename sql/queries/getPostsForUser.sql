-- name: GetPostsForUser :many
SELECT * FROM posts JOIN feeds ON posts.feed_id = feeds.id WHERE feeds.user_id = $1 ORDER BY posts.created_at DESC LIMIT $2;  