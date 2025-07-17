BEGIN;

-- Drop triggers first
DROP TRIGGER IF EXISTS trg_auto_assign_student_update ON LMS_USER;
DROP TRIGGER IF EXISTS trg_auto_assign_student_insert ON LMS_USER;
DROP TRIGGER IF EXISTS trg_prevent_student_tenant_membership ON Tenants_Members;

-- Drop functions
DROP FUNCTION IF EXISTS auto_assign_student_to_namespace_consumer();
DROP FUNCTION IF EXISTS prevent_student_tenant_membership();

-- Drop indexes
DROP INDEX IF EXISTS idx_report_date;
DROP INDEX IF EXISTS idx_report_generated_by;
DROP INDEX IF EXISTS idx_review_user;
DROP INDEX IF EXISTS idx_review_course;
DROP INDEX IF EXISTS idx_submission_date;
DROP INDEX IF EXISTS idx_submission_student;
DROP INDEX IF EXISTS idx_submission_assignment;
DROP INDEX IF EXISTS idx_assignment_course;
DROP INDEX IF EXISTS idx_lesson_title;
DROP INDEX IF EXISTS idx_lesson_module;
DROP INDEX IF EXISTS idx_student_quiz_score;
DROP INDEX IF EXISTS idx_student_quiz_quiz;
DROP INDEX IF EXISTS idx_student_quiz_student;
DROP INDEX IF EXISTS idx_quiz_module;
DROP INDEX IF EXISTS idx_module_name;
DROP INDEX IF EXISTS idx_module_course;
DROP INDEX IF EXISTS idx_rating_count;
DROP INDEX IF EXISTS idx_rating_course;
DROP INDEX IF EXISTS idx_rating_user;
DROP INDEX IF EXISTS idx_enrollment_date;
DROP INDEX IF EXISTS idx_enrollment_status;
DROP INDEX IF EXISTS idx_enrollment_course;
DROP INDEX IF EXISTS idx_enrollment_student;
DROP INDEX IF EXISTS idx_namespace_consumer_joined_date;
DROP INDEX IF EXISTS idx_namespace_consumer_active;
DROP INDEX IF EXISTS idx_namespace_consumer_namespace;
DROP INDEX IF EXISTS idx_namespace_consumer_user;
DROP INDEX IF EXISTS idx_category_created_at;
DROP INDEX IF EXISTS idx_category_name;
DROP INDEX IF EXISTS idx_course_owner_category;
DROP INDEX IF EXISTS idx_course_instructor_category;
DROP INDEX IF EXISTS idx_course_status;
DROP INDEX IF EXISTS idx_course_rating;
DROP INDEX IF EXISTS idx_course_created_at;
DROP INDEX IF EXISTS idx_course_title;
DROP INDEX IF EXISTS idx_course_owner;
DROP INDEX IF EXISTS idx_course_category;
DROP INDEX IF EXISTS idx_course_instructor;
DROP INDEX IF EXISTS idx_tenant_members_active_true;
DROP INDEX IF EXISTS idx_tenant_members_active_tenant;
DROP INDEX IF EXISTS idx_tenant_members_active_user;
DROP INDEX IF EXISTS idx_tenant_members_active;
DROP INDEX IF EXISTS idx_tenant_members_joined_date;
DROP INDEX IF EXISTS idx_tenant_members_tenant;
DROP INDEX IF EXISTS idx_tenant_members_user;
DROP INDEX IF EXISTS idx_tenants_active_true;
DROP INDEX IF EXISTS idx_tenants_created_at;
DROP INDEX IF EXISTS idx_tenants_active;
DROP INDEX IF EXISTS idx_tenants_owner;
DROP INDEX IF EXISTS idx_tenants_namespace;
DROP INDEX IF EXISTS idx_lms_user_role_active;
DROP INDEX IF EXISTS idx_lms_user_phone;
DROP INDEX IF EXISTS idx_lms_user_registration_date;
DROP INDEX IF EXISTS idx_lms_user_created_at;
DROP INDEX IF EXISTS idx_lms_user_role;
DROP INDEX IF EXISTS idx_lms_user_email;

-- Drop tables in reverse order of creation
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

-- Drop types
DROP TYPE IF EXISTS material_type;
DROP TYPE IF EXISTS enrollment_type;
DROP TYPE IF EXISTS course_status;
DROP TYPE IF EXISTS lms_role_type;

COMMIT;