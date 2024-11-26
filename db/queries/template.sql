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

-- name: FindTemplateByName :one
SELECT * FROM templates WHERE template_name = ? LIMIT 1;

-- name: FindAllTemplateNames :many
SELECT template_name FROM templates;