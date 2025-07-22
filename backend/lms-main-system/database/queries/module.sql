-- name: CreateModule :one
INSERT INTO Module (module_name,
                    course_id,
                    description,
                    updated_at)
VALUES ($1, $2, $3, CURRENT_TIMESTAMP) RETURNING *;

-- name: IsCourseOwnedByTenant :one
SELECT EXISTS (
    SELECT 1 
    FROM Course c
    JOIN Tenants t ON c.owned_by = t.tenant_id
    WHERE c.course_id = $1 AND t.namespace = $2
) AS is_owned;

-- name: GetModuleByIDWithTenant :one
SELECT m.*
FROM Module m
JOIN Course c ON c.course_id = m.course_id
JOIN Tenants t ON c.owned_by = t.tenant_id
WHERE m.module_id = $1 AND t.namespace = $2;

-- name: ListModulesWithTenant :many
SELECT m.*
FROM Module m
JOIN Course c ON c.course_id = m.course_id
JOIN Tenants t ON c.owned_by = t.tenant_id
WHERE t.namespace = $1
ORDER BY m.created_at DESC;

-- name: UpdateModuleByID :one
UPDATE Module
SET module_name = COALESCE($2, module_name),
    course_id   = COALESCE($3, course_id),
    description = COALESCE($4, description),
    updated_at  = CURRENT_TIMESTAMP
WHERE module_id = $1 RETURNING *;

-- name: IsModuleOwnedByTenant :one
SELECT EXISTS (
    SELECT 1
    FROM Module m
    JOIN Course c ON c.course_id = m.course_id
    JOIN Tenants t ON c.owned_by = t.tenant_id
    WHERE m.module_id = $1 AND t.namespace = $2
) as is_owned;

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
