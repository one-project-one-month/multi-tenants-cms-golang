
-- Create indexes for better performance

-- LMS_USER indexes
CREATE INDEX idx_lms_user_email ON LMS_USER(lms_user_email);
CREATE INDEX idx_lms_user_role ON LMS_USER(lms_role_id);
CREATE INDEX idx_lms_user_created_at ON LMS_USER(created_at);
CREATE INDEX idx_lms_user_registration_date ON LMS_USER(registration_date);
CREATE INDEX idx_lms_user_phone ON LMS_USER(phone_number);
CREATE INDEX idx_lms_user_role_active ON LMS_USER(lms_role_id, created_at);

-- Tenants indexes
CREATE INDEX idx_tenants_namespace ON Tenants(namespace);
CREATE INDEX idx_tenants_owner ON Tenants(cms_owner_id);
CREATE INDEX idx_tenants_active ON Tenants(is_active);
CREATE INDEX idx_tenants_created_at ON Tenants(created_at);
CREATE INDEX idx_tenants_active_true ON Tenants(tenant_id) WHERE is_active = TRUE;

-- Tenants_Members indexes
CREATE INDEX idx_tenant_members_user ON Tenants_Members(lms_user_id);
CREATE INDEX idx_tenant_members_tenant ON Tenants_Members(tenant_id);
CREATE INDEX idx_tenant_members_joined_date ON Tenants_Members(joined_date);
CREATE INDEX idx_tenant_members_active ON Tenants_Members(is_active);
CREATE INDEX idx_tenant_members_active_user ON Tenants_Members(is_active, lms_user_id);
CREATE INDEX idx_tenant_members_active_tenant ON Tenants_Members(is_active, tenant_id);
CREATE INDEX idx_tenant_members_active_true ON Tenants_Members(lms_user_id, tenant_id) WHERE is_active = TRUE;

-- Course indexes
CREATE INDEX idx_course_instructor ON Course(instructor_id);
CREATE INDEX idx_course_category ON Course(course_category);
CREATE INDEX idx_course_owner ON Course(owned_by);
CREATE INDEX idx_course_title ON Course(course_title);
CREATE INDEX idx_course_created_at ON Course(created_at);
CREATE INDEX idx_course_rating ON Course(overall_rating);
CREATE INDEX idx_course_status ON Course(status);
CREATE INDEX idx_course_instructor_category ON Course(instructor_id, course_category);
CREATE INDEX idx_course_owner_category ON Course(owned_by, course_category);

-- Course Category indexes
CREATE INDEX idx_category_name ON Course_Category(category_name);
CREATE INDEX idx_category_created_at ON Course_Category(created_at);

-- namespace_consumer indexes
CREATE INDEX idx_namespace_consumer_user ON namespace_consumer(lms_user_id);
CREATE INDEX idx_namespace_consumer_namespace ON namespace_consumer(namespace);
CREATE INDEX idx_namespace_consumer_active ON namespace_consumer(is_active);
CREATE INDEX idx_namespace_consumer_joined_date ON namespace_consumer(joined_date);

-- Enrollment indexes
CREATE INDEX idx_enrollment_student ON Enrollment(student_id);
CREATE INDEX idx_enrollment_course ON Enrollment(course_id);
CREATE INDEX idx_enrollment_status ON Enrollment(status);
CREATE INDEX idx_enrollment_date ON Enrollment(enrollment_date);

-- Rating indexes
CREATE INDEX idx_rating_user ON Rating(user_id);
CREATE INDEX idx_rating_course ON Rating(course_id);
CREATE INDEX idx_rating_count ON Rating(rating_count);

-- Module indexes
CREATE INDEX idx_module_course ON Module(course_id);
CREATE INDEX idx_module_name ON Module(module_name);

-- Quiz indexes
CREATE INDEX idx_quiz_module ON Quiz(module_id);

-- Student_Quiz indexes
CREATE INDEX idx_student_quiz_student ON Student_Quiz(student_id);
CREATE INDEX idx_student_quiz_quiz ON Student_Quiz(quiz_id);
CREATE INDEX idx_student_quiz_score ON Student_Quiz(score);

-- Lesson indexes
CREATE INDEX idx_lesson_module ON Lesson(module_id);
CREATE INDEX idx_lesson_title ON Lesson(title);

-- Assignment indexes
CREATE INDEX idx_assignment_course ON Assignment(course_id);

-- Submission indexes
CREATE INDEX idx_submission_assignment ON Submission(assignment_id);
CREATE INDEX idx_submission_student ON Submission(student_id);
CREATE INDEX idx_submission_date ON Submission(submitted_at);

-- Review indexes
CREATE INDEX idx_review_course ON Review(course_id);
CREATE INDEX idx_review_user ON Review(user_id);

-- Report indexes
CREATE INDEX idx_report_generated_by ON Report(generated_by_user_id);
CREATE INDEX idx_report_date ON Report(generated_date);