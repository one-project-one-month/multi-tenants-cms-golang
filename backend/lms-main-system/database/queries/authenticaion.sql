-- name: CreateUser :one
INSERT INTO lms_user (lms_user_name,
                      lms_user_email,
                      password,
                      address,
                      tenant_id,
                      phone_number,
                      registration_date,
                      email_verified,
                      namespace_domain
                      )
VALUES ($1, $2, $3, $4, $5, $6, CURRENT_DATE,
        FALSE, $7) RETURNING lms_user_id, lms_user_name, lms_user_email, address, phone_number, registration_date, email_verified, namespace_domain, mfa_enable, created_at, updated_at;
-- name: GetTenantIdByNameSpace :one
SELECT tenant_id
FROM tenants
WHERE namespace = $1;

-- name: CheckNameSpaceFromMetaData :one
SELECT 1
FROM tenants
WHERE namespace = $1;

-- name: CheckTenantsMemberExistenceByEmail :one
SELECT 1
FROM tenants_members
WHERE lms_user_email = $1;

-- name: GetRoleIdByName :one
SELECT lms_role_id
FROM lms_user_role
WHERE lms_role_name = $1;

-- name: GetDefaultRoleIDs :many
SELECT lms_role_id
FROM lms_user_role
WHERE lms_role_name IN ('STUDENT', 'VIEWER');

-- name: AssignRolesToUser :exec
INSERT INTO lms_user_roles_map (lms_user_id, lms_role_id)
VALUES ($1, $2);

-- name: AssignMultipleRolesToUser :exec
INSERT INTO lms_user_roles_map (lms_user_id, lms_role_id)
SELECT $1, unnest($2::int[]);


-- name: UpdateEmailVerification :exec
UPDATE lms_user
SET email_verified = true,
    updated_at     = now()
WHERE lms_user_email = $1;

-- name: GetUserByEmail :one
SELECT *
FROM lms_user
WHERE lms_user_email = $1;

-- name: GetUserPasswordMfaByNameSpaceDomain :one
SELECT lms_user_id, password,mfa_enable
FROM lms_user WHERE  namespace_domain = $1 ;

-- name: SetUpMFA :one
INSERT  INTO  lms_user_mfa (mfa_secret, lms_user_domain_email)
VALUES ($1, $2) RETURNING  mfa_secret_id;