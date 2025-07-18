-- name: CreateUser :one
INSERT INTO lms_user (
    lms_user_name,
    lms_user_email,
    password,
    address,
    tenant_id,
    phone_number,
    registration_date,
    email_verified
) VALUES (
             $1, $2, $3, $4, $5, $6,CURRENT_DATE, FALSE
         )
RETURNING lms_user_id, lms_user_name, lms_user_email, address, phone_number, registration_date, email_verified, mfa_enable, created_at, updated_at;
-- name: GetTenantIdByNameSpace :one
SELECT tenant_id FROM tenants WHERE  namespace = $1;

-- name: CheckNameSpaceFromMetaData :one
SELECT  1 FROM tenants WHERE  namespace = $1;

-- name: CheckTenantsMemberExistenceByEmail :one
SELECT  1 FROM tenants_members WHERE lms_user_email = $1;

-- name: GetRoleIdByName :one
SELECT lms_role_id FROM lms_user_role WHERE  lms_role_name = $1;

-- name: GetDefaultRoleIDs :many
SELECT lms_role_id FROM lms_user_role
WHERE lms_role_name IN ('STUDENT', 'VIEWER');

-- name: AssignRolesToUser :exec
INSERT INTO lms_user_roles_map (lms_user_id, lms_role_id)
VALUES ($1, $2);

-- name: AssignMultipleRolesToUser :exec
INSERT INTO lms_user_roles_map (lms_user_id, lms_role_id)
SELECT $1, unnest($2::int[]);


-- name: UpdateEmailVerification :exec
UPDATE lms_user
SET email_verified = true, updated_at = now()
WHERE lms_user_email = $1;

-- name: GetUserByEmail :one
SELECT  * FROM lms_user WHERE  lms_user_email = $1 ;