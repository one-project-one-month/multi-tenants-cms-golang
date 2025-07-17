-- name: RegisterUserWithRoles :one
WITH inserted_user AS (
    INSERT INTO lms_user (
                          lms_user_name,
                          lms_user_email,
                          password,
                          address,
                          phone_number,
                          registration_date,
                          email_verified
        ) VALUES (
                     $1, $2, $3, $4, $5, CURRENT_DATE, FALSE
                 ) RETURNING lms_user_id, lms_user_name, lms_user_email,address,phone_number,registration_date,created_at,updated_at
),
     role_ids AS (
         SELECT lms_role_id
         FROM LMS_USER_Role
         WHERE lms_role_name IN ('STUDENT', 'VIEWER')
     ),
     role_assignments AS (
         INSERT INTO lms_user_roles_map (lms_user_id, lms_role_id)
             SELECT iu.lms_user_id, ri.lms_role_id
             FROM inserted_user iu
                      CROSS JOIN role_ids ri
             RETURNING lms_user_id
     )
SELECT lms_user_id, lms_user_name, lms_user_email,address,registration_date,created_at,updated_at,phone_number
FROM inserted_user;

-- name: UpdateEmailVerification :exec
UPDATE  lms_user
SET  email_verified = true ,updated_at = now()
WHERE  lms_user_email = $1;