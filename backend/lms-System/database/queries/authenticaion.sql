-- name: RegisterLMSUser :one
INSERT INTO lms_user (
    lms_user_name,
    lms_user_email,
    password,
    address,
    phone_number
) VALUES (
             $1, $2, $3, $4, $5
         )
RETURNING
    lms_user_id,
    lms_user_name,
    lms_user_email,
    address,
    phone_number,
    registration_date,
    updated_at
;


-- name: UpdateEmailVerification :exec
UPDATE lms_user
SET  email_verified = true , updated_at = now()
WHERE lms_user_id = $1;
