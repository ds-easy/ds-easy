-- name: FindAllUsers :many
SELECT * FROM users;

-- name: AddUser :one
INSERT INTO
    users (
        pb_id,
        first_name,
        last_name,
        email,
        admin
    )
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: FindUserById :one
SELECT * FROM users where id = ?;

-- name: FindUserByPBId :one
SELECT * FROM users where pb_id = ?;

-- name: FindUserByEmail :one
SELECT * FROM users where email = ?;
