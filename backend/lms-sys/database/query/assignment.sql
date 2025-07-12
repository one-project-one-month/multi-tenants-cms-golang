-- name: CreateUser :one
INSERT INTO LMS_USER (
    lms_user_email,
    password,
    lms_role_id,
    tenant_id,
    address,
    phone_number,
    registration_date
) VALUES (
             $1, $2, $3, $4, $5, $6, $7
         ) RETURNING *;

-- name: GetUser :one
SELECT * FROM LMS_USER
WHERE lms_user_id = $1 LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM LMS_USER
WHERE lms_user_email = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM LMS_USER
ORDER BY lms_user_email
    LIMIT $1
OFFSET $2;

-- name: UpdateUser :one
UPDATE LMS_USER
SET
    lms_user_email = COALESCE($2, lms_user_email),
    password = COALESCE($3, password),
    lms_role_id = COALESCE($4, lms_role_id),
    tenant_id = COALESCE($5, tenant_id),
    address = COALESCE($6, address),
    phone_number = COALESCE($7, phone_number),
    updated_at = CURRENT_TIMESTAMP
WHERE lms_user_id = $1
    RETURNING *;

-- name: DeleteUser :exec
DELETE FROM LMS_USER
WHERE lms_user_id = $1;

-- name: GetUsersByTenant :many
SELECT * FROM LMS_USER
WHERE tenant_id = $1
ORDER BY lms_user_email;

-- name: GetUsersByRole :many
SELECT * FROM LMS_USER
WHERE lms_role_id = $1
ORDER BY lms_user_email;

-- name: UpdateUserPassword :exec
UPDATE LMS_USER
SET password = $2, updated_at = CURRENT_TIMESTAMP
WHERE lms_user_id = $1;