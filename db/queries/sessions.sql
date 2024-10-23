-- name: DeleteSession :exec
DELETE FROM sessions WHERE id = ?;

-- name: FindSessionById :one
SELECT * FROM sessions WHERE id = ?;

-- name: CreateSession :one
INSERT INTO sessions (
	id,
	user_id
)
VALUES (?, ?) RETURNING *;
