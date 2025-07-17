-- name: GetUserByEmail :one
SELECT  * FROM lms_user WHERE  lms_user_email = $1;