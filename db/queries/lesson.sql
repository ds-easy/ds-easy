-- name: FindLessons :many
SELECT * FROM lessons;

-- name: InsertLesson :one
INSERT INTO
    lessons (lesson_name, year, subject)
VALUES (?, ?, ?) RETURNING *;

-- name: FindAllLessonNames :many
SELECT lesson_name FROM lessons;