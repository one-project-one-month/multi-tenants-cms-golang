-- name: CreateLesson :one
INSERT INTO Lesson (title, content, material_type, module_id)
VALUES ($1, $2, $3, $4)
    RETURNING *;

-- name: GetLessonById :one
SELECT
    l.*,
    m.module_id AS module_id,
    m.module_name AS module_name
FROM Lesson l
         JOIN Module m ON l.module_id = m.module_id
WHERE l.lesson_id = $1;

-- name: GetLessonsByModuleId :many
SELECT
    l.*,
    m.module_id AS module_id,
    m.module_name AS module_name
FROM Lesson l
         JOIN Module m ON l.module_id = m.module_id
WHERE l.module_id = $1
ORDER BY l.created_at DESC;

-- name: UpdateLesson :one
UPDATE Lesson
SET
    title = $2,
    content = $3,
    material_type = $4,
    module_id = $5,
    updated_at = CURRENT_TIMESTAMP
WHERE lesson_id = $1
    RETURNING *;

-- name: DeleteLessons :many
DELETE FROM Lesson
WHERE lesson_id = ANY($1::uuid[])
    RETURNING lesson_id;
