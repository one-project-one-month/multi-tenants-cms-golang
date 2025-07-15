-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS Report CASCADE;
DROP TABLE IF EXISTS Review CASCADE;
DROP TABLE IF EXISTS Submission CASCADE;
DROP TABLE IF EXISTS Assignment CASCADE;
DROP TABLE IF EXISTS Lesson CASCADE;
DROP TABLE IF EXISTS Student_Quiz CASCADE;
DROP TABLE IF EXISTS Quiz CASCADE;
DROP TABLE IF EXISTS Module CASCADE;
DROP TABLE IF EXISTS Rating CASCADE;
DROP TABLE IF EXISTS Certificate CASCADE;
DROP TABLE IF EXISTS enrollment CASCADE;
DROP TABLE IF EXISTS namespace_consumer CASCADE;
DROP TABLE IF EXISTS Course CASCADE;
DROP TABLE IF EXISTS Course_Category CASCADE;
DROP TABLE IF EXISTS Tenants_Members CASCADE;
DROP TABLE IF EXISTS LMS_USER CASCADE;
DROP TABLE IF EXISTS Tenants CASCADE;
DROP TABLE IF EXISTS LMS_USER_Role CASCADE;

-- Drop indexes (if they exist independently)
DROP INDEX IF EXISTS idx_lms_user_email;
DROP INDEX IF EXISTS idx_lms_user_role;
DROP INDEX IF EXISTS idx_lms_user_tenant;
DROP INDEX IF EXISTS idx_tenants_namespace;
DROP INDEX IF EXISTS idx_course_instructor;
DROP INDEX IF EXISTS idx_course_category;
DROP INDEX IF EXISTS idx_course_tenant;
DROP INDEX IF EXISTS idx_course_status;
DROP INDEX IF EXISTS idx_enrollment_student;
DROP INDEX IF EXISTS idx_enrollment_course;
DROP INDEX IF EXISTS idx_enrollment_status;
DROP INDEX IF EXISTS idx_enrollments_student_id;
DROP INDEX IF EXISTS idx_enrollments_course_id;
DROP INDEX IF EXISTS idx_module_course;
DROP INDEX IF EXISTS idx_lesson_module;
DROP INDEX IF EXISTS idx_quiz_module;
DROP INDEX IF EXISTS idx_assignment_course;
DROP INDEX IF EXISTS idx_submission_assignment;
DROP INDEX IF EXISTS idx_submission_student;
DROP INDEX IF EXISTS idx_rating_course;
DROP INDEX IF EXISTS idx_rating_user;

-- Drop functions and triggers
DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;

-- Drop custom types
DROP TYPE IF EXISTS course_status CASCADE;
DROP TYPE IF EXISTS material_type CASCADE;
DROP TYPE IF EXISTS enrollment_type CASCADE;
DROP TYPE IF EXISTS system_type CASCADE;
DROP TYPE IF EXISTS lms_role_type CASCADE;
-- Drop all existing indexes
DROP INDEX IF EXISTS idx_lms_user_email;
DROP INDEX IF EXISTS idx_lms_user_role;
DROP INDEX IF EXISTS idx_lms_user_tenant;
DROP INDEX IF EXISTS idx_tenants_namespace;
DROP INDEX IF EXISTS idx_course_instructor;
DROP INDEX IF EXISTS idx_course_category;
DROP INDEX IF EXISTS idx_course_tenant;
DROP INDEX IF EXISTS idx_course_status;
DROP INDEX IF EXISTS idx_enrollment_student;
DROP INDEX IF EXISTS idx_enrollment_course;
DROP INDEX IF EXISTS idx_enrollment_status;
DROP INDEX IF EXISTS idx_enrollments_student_id;
DROP INDEX IF EXISTS idx_enrollments_course_id;
DROP INDEX IF EXISTS idx_module_course;
DROP INDEX IF EXISTS idx_lesson_module;
DROP INDEX IF EXISTS idx_lesson_title;
DROP INDEX IF EXISTS idx_quiz_module;
DROP INDEX IF EXISTS idx_assignment_course;
DROP INDEX IF EXISTS idx_submission_assignment;
DROP INDEX IF EXISTS idx_submission_student;
DROP INDEX IF EXISTS idx_submission_date;
DROP INDEX IF EXISTS idx_rating_course;
DROP INDEX IF EXISTS idx_rating_user;
DROP INDEX IF EXISTS idx_review_course;
DROP INDEX IF EXISTS idx_review_user;
DROP INDEX IF EXISTS idx_report_generated_by;
DROP INDEX IF EXISTS idx_report_date;