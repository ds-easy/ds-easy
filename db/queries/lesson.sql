-- name: FindLessons :many
SELECT * FROM lessons;

-- name: InsertLesson :one
INSERT INTO
    lessons (lesson_name)
VALUES (?) RETURNING *;

-- name: FindAllLessonNames :many
SELECT lesson_name FROM lessons;

-- name: FindLessonByName :one
SELECT * FROM lessons WHERE lessons.lesson_name = ? LIMIT 1;