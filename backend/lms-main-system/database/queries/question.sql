-- name: CreateQuestion :one
INSERT INTO Question (question, quiz_id)
VALUES ($1, $2)
    RETURNING *;

-- name: UpdateQuestion :one
UPDATE Question
SET question = $1, updated_at = now()
WHERE question_id = $2
    RETURNING *;

-- name: GetQuestionsByQuizID :many
SELECT * FROM Question
WHERE quiz_id = $1;

-- name: DeleteQuestionsByQuizIDs :exec
DELETE FROM Question
WHERE quiz_id = ANY($1::uuid[]);

-- name: GetQuestionByID :one
SELECT * FROM Question
WHERE question_id = $1;

-- name: DeleteQuestionsByQuizId :exec
DELETE FROM Question
WHERE quiz_id = $1;