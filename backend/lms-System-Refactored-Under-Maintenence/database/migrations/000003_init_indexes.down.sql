-- LMS_USER indexes
DROP INDEX IF EXISTS idx_lms_user_email;
DROP INDEX IF EXISTS idx_lms_user_role;
DROP INDEX IF EXISTS idx_lms_user_created_at;
DROP INDEX IF EXISTS idx_lms_user_registration_date;
DROP INDEX IF EXISTS idx_lms_user_phone;
DROP INDEX IF EXISTS idx_lms_user_role_active;

-- Tenants indexes
DROP INDEX IF EXISTS idx_tenants_namespace;
DROP INDEX IF EXISTS idx_tenants_owner;
DROP INDEX IF EXISTS idx_tenants_active;
DROP INDEX IF EXISTS idx_tenants_created_at;
DROP INDEX IF EXISTS idx_tenants_active_true;

-- Tenants_Members indexes
DROP INDEX IF EXISTS idx_tenant_members_user;
DROP INDEX IF EXISTS idx_tenant_members_tenant;
DROP INDEX IF EXISTS idx_tenant_members_joined_date;
DROP INDEX IF EXISTS idx_tenant_members_active;
DROP INDEX IF EXISTS idx_tenant_members_active_user;
DROP INDEX IF EXISTS idx_tenant_members_active_tenant;
DROP INDEX IF EXISTS idx_tenant_members_active_true;

-- Course indexes
DROP INDEX IF EXISTS idx_course_instructor;
DROP INDEX IF EXISTS idx_course_category;
DROP INDEX IF EXISTS idx_course_owner;
DROP INDEX IF EXISTS idx_course_title;
DROP INDEX IF EXISTS idx_course_created_at;
DROP INDEX IF EXISTS idx_course_rating;
DROP INDEX IF EXISTS idx_course_status;
DROP INDEX IF EXISTS idx_course_instructor_category;
DROP INDEX IF EXISTS idx_course_owner_category;

-- Course Category indexes
DROP INDEX IF EXISTS idx_category_name;
DROP INDEX IF EXISTS idx_category_created_at;

-- namespace_consumer indexes
DROP INDEX IF EXISTS idx_namespace_consumer_user;
DROP INDEX IF EXISTS idx_namespace_consumer_namespace;
DROP INDEX IF EXISTS idx_namespace_consumer_active;
DROP INDEX IF EXISTS idx_namespace_consumer_joined_date;

-- Enrollment indexes
DROP INDEX IF EXISTS idx_enrollment_student;
DROP INDEX IF EXISTS idx_enrollment_course;
DROP INDEX IF EXISTS idx_enrollment_status;
DROP INDEX IF EXISTS idx_enrollment_date;

-- Rating indexes
DROP INDEX IF EXISTS idx_rating_user;
DROP INDEX IF EXISTS idx_rating_course;
DROP INDEX IF EXISTS idx_rating_count;

-- Module indexes
DROP INDEX IF EXISTS idx_module_course;
DROP INDEX IF EXISTS idx_module_name;

-- Quiz indexes
DROP INDEX IF EXISTS idx_quiz_module;

-- Student_Quiz indexes
DROP INDEX IF EXISTS idx_student_quiz_student;
DROP INDEX IF EXISTS idx_student_quiz_quiz;
DROP INDEX IF EXISTS idx_student_quiz_score;

-- Lesson indexes
DROP INDEX IF EXISTS idx_lesson_module;
DROP INDEX IF EXISTS idx_lesson_title;

-- Assignment indexes
DROP INDEX IF EXISTS idx_assignment_course;

-- Submission indexes
DROP INDEX IF EXISTS idx_submission_assignment;
DROP INDEX IF EXISTS idx_submission_student;
DROP INDEX IF EXISTS idx_submission_date;

-- Review indexes
DROP INDEX IF EXISTS idx_review_course;
DROP INDEX IF EXISTS idx_review_user;

-- Report indexes
DROP INDEX IF EXISTS idx_report_generated_by;
DROP INDEX IF EXISTS idx_report_date;
