-- name: FindTemplates :many
SELECT * FROM templates;

-- name: InsertTemplate :one
INSERT INTO
    templates (
        uploaded_by,
        pb_file_id,
        template_name
    )
VALUES (?, ?, ?) RETURNING *;