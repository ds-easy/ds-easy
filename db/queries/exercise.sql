-- name: FindExercises :many
SELECT * FROM exercises;

-- name: InsertExercise :one
INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by
    )
VALUES (?, ?, ?, ?) RETURNING *;

-- name: FindExercisesByLessonName :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ?;

-- name: FindExercisesByName :one
SELECT * FROM exercises WHERE exercises.exercise_name = ? LIMIT 1;

-- name: FindRandomExercisesByLessonNameWithLimit :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ?
ORDER BY RANDOM()
LIMIT ?;