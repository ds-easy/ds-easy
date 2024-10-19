-- name: FindAllUsers :many
SELECT * FROM users;

-- name: AddUser :one
INSERT INTO
    users (
        first_name,
        last_name,
        email,
        password,
        admin
    )
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: FindUserById :one
SELECT * FROM users where id = ?;

-- name: FindUserByEmail :one
SELECT * FROM users where email = ?;