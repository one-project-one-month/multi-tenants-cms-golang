-- name: CreateEnrollment :one
INSERT INTO Enrollment (
    student_id,   
    course_id,
    status,
    due_date,
    updated_at
) VALUES (
    $1,
    $2,
    $3,
    $4,
    CURRENT_TIMESTAMP
) RETURNING enrollment_id;

-- name: ListEnrollmentsWithTenant :many
SELECT * FROM enrollment_details
WHERE namespace = $1;

-- name: ListEnrollmentsWithStudentEmails :many
SELECT * FROM enrollment_details
WHERE namespace = $1 AND student_email = $2;

-- name: ListEnrollmentsWithCourseID :many
SELECT * FROM enrollment_details
WHERE namespace = $1 AND course_id = $2;

-- name: ListEnrollmentsWithCategoryID :many
SELECT * FROM enrollment_details
WHERE namespace = $1 AND category_id = $2;

-- name: DeleteEnrollmentsByID :exec
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

-- name: IsEnrollmentExist :one
SELECT EXISTS (
    SELECT 1
    FROM Enrollment e
    WHERE student_id = $1 AND course_id = $2
) as exists;

-- name: IsEnrolleeStudent :one
SELECT EXISTS (
    SELECT 1 
    FROM LMS_USER u
    JOIN lms_user_roles_map rm ON rm.lms_user_id = u.lms_user_id
    JOIN LMS_USER_Role r ON r.lms_role_id = rm.lms_role_id
    JOIN Tenants t ON t.tenant_id = u.tenant_id
    WHERE t.namespace = $1 
        AND t.is_active = true
        AND u.lms_user_id = $2 
        AND r.lms_role_name = 'STUDENT'
        AND rm.is_active = true
) as is_student;
