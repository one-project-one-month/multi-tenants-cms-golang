-- name: CreateNewCourse :one
INSERT  INTO course (course_title)
VALUES  ($1)
RETURNING  *;


-- name: GetAllCourse :many
SELECT * FROM  course;