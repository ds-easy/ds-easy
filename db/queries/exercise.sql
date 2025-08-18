-- name: FindExercises :many
SELECT * FROM exercises;

-- name: FindPublicExercises :many
SELECT * FROM exercises WHERE is_public = true;

-- name: FindAccessibleExercises :many
SELECT *
FROM exercises e
WHERE
  e.is_public = true OR e.uploaded_by = ?;


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

-- name: FindExercisesByName :one
SELECT * FROM exercises WHERE exercises.exercise_name = ? LIMIT 1;

-- name: FindPublicExercisesByName :one
SELECT * FROM exercises WHERE exercises.exercise_name = ? AND is_public = true LIMIT 1;

-- name: FindPublicExercisesByLessonName :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ? AND e.is_public = true;

-- name: FindAccessibleExercisesByLessonName :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ?
    AND (e.is_public = true OR e.uploaded_by = ?);

-- name: FindRandomPublicExercisesByLessonNameWithLimit :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ? AND e.is_public = true
ORDER BY RANDOM()
LIMIT ?;

-- name: FindRandomAccessibleExercisesByLessonNameWithLimit :many
SELECT e.*
FROM exercises e
    LEFT JOIN lessons l ON e.lesson_id = l.id
WHERE
    l.lesson_name = ?
    AND (e.is_public = true OR e.uploaded_by = ?)
ORDER BY RANDOM()
LIMIT ?;