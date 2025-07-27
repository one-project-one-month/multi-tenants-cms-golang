-- name: CreateQuiz :one
INSERT INTO Quiz (module_id, title)
VALUES ($1, $2)
RETURNING *;

-- name: GetQuizzesByModule :many
SELECT * FROM Quiz
WHERE module_id = $1;

-- name: DeleteQuizzes :exec
DELETE FROM Quiz
WHERE quiz_id = ANY($1::uuid[]);

-- name: GetQuizById :one
SELECT * FROM Quiz
WHERE quiz_id = $1;

-- name: DeleteQuizById :exec
DELETE FROM Quiz
WHERE quiz_id = $1;