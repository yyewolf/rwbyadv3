-- name: GetPlayer :one
SELECT * FROM players WHERE id = $1 LIMIT 1;
