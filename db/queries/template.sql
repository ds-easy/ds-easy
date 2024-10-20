-- name: FindTemplates :many
SELECT * FROM template;

-- name: InsertTemplate :one
INSERT INTO
    template (uploaded_by, pb_file_id)
VALUES (?, ?) RETURNING *;