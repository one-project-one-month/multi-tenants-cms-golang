-- name: CreateAnswer :one
INSERT INTO Answer (answer, is_correct, question_id)
VALUES ($1, $2, $3)
    RETURNING *;

-- name: UpdateAnswer :one
UPDATE Answer
SET answer = $1, is_correct = $2, updated_at = now()
WHERE answer_id = $3
    RETURNING *;

-- name: GetAnswersByQuestionID :many
SELECT * FROM Answer
WHERE question_id = $1;

-- name: DeleteAnswersByQuestionIDs :exec
DELETE FROM Answer
WHERE question_id = ANY($1::uuid[]);

-- name: DeleteAnswersByQuestionId :exec
DELETE FROM Answer
WHERE question_id = $1;