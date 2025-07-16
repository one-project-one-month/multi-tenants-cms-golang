-- name: CreateNewCourse :one
INSERT  INTO course (course_title)
VALUES  ($1)
RETURNING  *;


-- name: GetAllCourse :many
SELECT * FROM  course;



SELECT *
FROM "Course"
WHERE id = $1;  

UPDATE "Course"
SET
    course_title = $2,
    course_description = $3 
WHERE id = $1
RETURNING *;

DELETE FROM "Course"
WHERE id = $1;