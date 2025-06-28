-- name: FindExercises :many
SELECT * FROM exercises;

-- name: FindPublicExercises :many
SELECT * FROM exercises WHERE is_public = true;

-- name: InsertExercise :one
INSERT INTO
    exercises (
        exercise_name,
        exercise_path,
        lesson_id,
        uploaded_by,
        is_public
    )
VALUES (?, ?, ?, ?, ?) RETURNING *;

-- name: FindExercisesByLessonName :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ?;

-- name: FindPublicExercisesByLessonName :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ? AND e.is_public = true;

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

-- name: FindRandomPublicExercisesByLessonNameWithLimit :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ? AND e.is_public = true
ORDER BY RANDOM()
LIMIT ?;