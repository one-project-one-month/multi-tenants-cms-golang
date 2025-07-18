-- name: CreateModule :one
INSERT INTO Module (module_name,
                    course_id,
                    description,
                    updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP) RETURNING *;

-- name: GetModuleByID :one
SELECT *
FROM Module
WHERE module_id = $1;

-- name: ListModules :many
SELECT *
FROM Module
ORDER BY created_at DESC;

-- name: UpdateModuleByID :one
UPDATE Module
SET module_name = COALESCE($2, module_name),
    course_id   = COALESCE($3, course_id),
    description = COALESCE($4, description),
    updated_at  = CURRENT_TIMESTAMP
WHERE module_id = $1 RETURNING *;

-- name: DeleteModule :exec
DELETE
FROM Module
WHERE module_id = $1;

-- name: DeleteModules :exec
DELETE
FROM Module
WHERE module_id = ANY ($1::uuid[]);

-- name: ModuleHasAssociations :one
SELECT (
           (SELECT COUNT(*) FROM Quiz q WHERE q.module_id = $1) > 0
               OR
           (SELECT COUNT(*) FROM Lesson l WHERE l.module_id = $1) > 0
           ) AS exists;

-- name: DeleteAssociatedQuizzes :exec
DELETE
FROM Quiz
WHERE module_id = $1;

-- name: DeleteAssociatedLessons :exec
DELETE
FROM Lesson
WHERE module_id = $1;
