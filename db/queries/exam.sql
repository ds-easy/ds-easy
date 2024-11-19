-- name: FindExams :many
SELECT * FROM exams;

-- name: InsertExam :one
INSERT INTO
    exams (
        date_of_passing,
        exam_number,
        professor_id,
        template_id
    )
VALUES (?, ?, ?, ?) RETURNING *;