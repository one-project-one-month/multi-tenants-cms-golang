-- name: CreateEnrollment :one
INSERT INTO Enrollment (
    student_id,   
    course_id,
    status,
    due_date,
    updated_at
)
SELECT
    $1,
    $2,
    $3,
    $4,
    CURRENT_TIMESTAMP
FROM Course c
WHERE c.course_id = $2 AND c.status = 'Active'
RETURNING enrollment_id;

-- name: ListEnrollmentsByAdminFilters :many
SELECT * FROM enrollment_details
WHERE namespace = $1 
    AND ($2::uuid = '00000000-0000-0000-0000-000000000000' OR course_id = $2) 
    AND ($3::uuid = '00000000-0000-0000-0000-000000000000' OR category_id = $3)
    AND ($4::text = '' OR student_email = $4);

-- name: ListEnrollmentsByInstructorFilters :many
SELECT * FROM enrollment_details
WHERE namespace = $1 
    AND instructor_id = $2 
    AND ($3::uuid = '00000000-0000-0000-0000-000000000000' OR course_id = $3) 
    AND ($4::uuid = '00000000-0000-0000-0000-000000000000' OR category_id = $4)
    AND ($5::text = '' OR student_email = $5);

-- name: DeleteEnrollmentByID :exec
DELETE FROM Enrollment
WHERE enrollment_id = $1;

-- name: IsEnrollmentOwnedByTenant :one
SELECT EXISTS (
    SELECT 1
    From Enrollment e
    JOIN Course c ON c.course_id = e.course_id
    JOIN Tenants t ON t.tenant_id = c.owned_by
    WHERE t.namespace = $1 AND e.enrollment_id = $2
) as is_owned;

-- name: GetEnrollmentByIDAndTenant :one
SELECT * FROM enrollment_details
WHERE namespace = $1 AND enrollment_id = $2;

-- name: UpdateEnrollmentStatus :one
UPDATE Enrollment
SET status      = COALESCE($2, status),
    updated_at  = CURRENT_TIMESTAMP
WHERE enrollment_id = $1
RETURNING *;

-- name: UpdateEnrollmentProgress :one
UPDATE Enrollment
SET progress    = COALESCE($2, progress),
    updated_at  = CURRENT_TIMESTAMP
WHERE enrollment_id = $1
RETURNING *;

-- name: IsEnrollmentExistUnderCourseID :one
SELECT EXISTS (
    SELECT 1
    FROM Enrollment
    WHERE student_id = $1 AND course_id = $2
) as exists;

-- name: IsEnrollmentExists :one
SELECT EXISTS (
    SELECT 1
    FROM Enrollment
    WHERE enrollment_id = $1
) as exists;

-- name: IsEnrollmentBelongsToStudent :one
SELECT EXISTS (
    SELECT 1
    FROM Enrollment
    WHERE student_id = $1 AND enrollment_id = $2
) as exists;

-- name: IsUserInRole :one
SELECT EXISTS (
    SELECT 1 
    FROM user_role_details
    WHERE namespace = $1 
        AND lms_user_id = $2 
        AND lms_role_name = $3
        AND tenant_is_active = true
        AND role_map_is_active = true
) as in_role;
