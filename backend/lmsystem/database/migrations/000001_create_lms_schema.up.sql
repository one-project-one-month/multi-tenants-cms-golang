-- Create custom types (must be created first)
CREATE TYPE lms_role_type AS ENUM ('LMS_ADMIN', 'STUDENT', 'INSTRUCTOR');
CREATE TYPE enrollment_type AS ENUM('ENROLLED', 'COMPLETED', 'DROPPED');
CREATE TYPE material_type AS ENUM('VIDEO', 'PDF', 'SLIDE', 'LINK');
CREATE TYPE course_status AS ENUM('PENDING', 'PUBLISHED', 'UNPUBLISHED', 'ARCHIVED');

-- Create tables in proper dependency order

-- 1. Role table (no dependencies)
CREATE TABLE LMS_USER_Role (
                               lms_role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               lms_role_name lms_role_type NOT NULL UNIQUE,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 2. Tenants table (no dependencies initially)
CREATE TABLE Tenants (
                         tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         namespace VARCHAR(255) UNIQUE NOT NULL,
                         cms_owner_id UUID, -- Will be constrained later to avoid circular dependency
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         is_active BOOLEAN DEFAULT TRUE
);

-- 3. User table (depends on Role and Tenants)
CREATE TABLE LMS_USER (
                          lms_user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          lms_user_email VARCHAR(255) UNIQUE NOT NULL,
                          password VARCHAR(255) NOT NULL,
                          first_name VARCHAR(100),
                          last_name VARCHAR(100),
                          lms_role_id UUID NOT NULL,
                          tenant_id UUID,
                          address TEXT,
                          phone_number VARCHAR(20),
                          registration_date DATE DEFAULT CURRENT_DATE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          is_active BOOLEAN DEFAULT TRUE,
                          CONSTRAINT fk_lms_user_role
                              FOREIGN KEY (lms_role_id)
                                  REFERENCES LMS_USER_Role(lms_role_id) ON DELETE RESTRICT,
                          CONSTRAINT fk_lms_user_tenant
                              FOREIGN KEY (tenant_id)
                                  REFERENCES Tenants(tenant_id) ON DELETE SET NULL
);

-- 4. Add the circular dependency constraint after LMS_USER is created
ALTER TABLE Tenants
    ADD CONSTRAINT fk_tenants_cms_owner
        FOREIGN KEY (cms_owner_id)
            REFERENCES LMS_USER(lms_user_id) ON DELETE SET NULL;

-- 5. Tenants_Members table (depends on User and Tenants)
CREATE TABLE Tenants_Members (
                                 tm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 lms_user_id UUID NOT NULL,
                                 tenant_id UUID NOT NULL,
                                 joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 is_active BOOLEAN DEFAULT TRUE,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 CONSTRAINT fk_tenant_member_user
                                     FOREIGN KEY (lms_user_id)
                                         REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                 CONSTRAINT fk_tenant_member_tenant
                                     FOREIGN KEY (tenant_id)
                                         REFERENCES Tenants(tenant_id) ON DELETE CASCADE,
                                 UNIQUE(lms_user_id, tenant_id)
);

-- 6. Course Category table (no dependencies)
CREATE TABLE Course_Category (
                                 category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 category_name VARCHAR(100) NOT NULL UNIQUE,
                                 description TEXT,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 is_active BOOLEAN DEFAULT TRUE
);

-- 7. Course table (depends on User, Tenants, and Course_Category)
CREATE TABLE Course (
                        course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        course_title VARCHAR(150) NOT NULL,
                        description TEXT,
                        instructor_id UUID NOT NULL,
                        overall_rating DECIMAL(3,2) CHECK (overall_rating >= 1.0 AND overall_rating <= 5.0),
                        course_category UUID NOT NULL,
                        status course_status DEFAULT 'PENDING',
                        duration_day_count INT CHECK (duration_day_count > 0),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        owned_by UUID NOT NULL,
                        is_active BOOLEAN DEFAULT TRUE,
                        CONSTRAINT fk_course_category
                            FOREIGN KEY (course_category)
                                REFERENCES Course_Category(category_id) ON DELETE RESTRICT,
                        CONSTRAINT fk_course_instructor
                            FOREIGN KEY (instructor_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE RESTRICT,
                        CONSTRAINT fk_course_tenant
                            FOREIGN KEY (owned_by)
                                REFERENCES Tenants(tenant_id) ON DELETE CASCADE
);

-- 8. Namespace Consumer table (depends on User)
CREATE TABLE namespace_consumer (
                                    consumer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    lms_user_id UUID NOT NULL,
                                    namespace VARCHAR(255) NOT NULL DEFAULT 'default_consumer_namespace',
                                    joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    is_active BOOLEAN DEFAULT TRUE,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT fk_namespace_consumer_user
                                        FOREIGN KEY (lms_user_id)
                                            REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                    UNIQUE(lms_user_id, namespace)
);

-- 9. Enrollment table
CREATE TABLE enrollment (
                            enrollment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            student_id UUID NOT NULL,
                            course_id UUID NOT NULL,
                            enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            progress DECIMAL(5,2) DEFAULT 0.0 CHECK (progress >= 0.0 AND progress <= 100.0),
                            status enrollment_type NOT NULL DEFAULT 'ENROLLED',
                            due_date TIMESTAMP,
                            completed_at TIMESTAMP,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            CONSTRAINT fk_enrollment_student
                                FOREIGN KEY (student_id)
                                    REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                            CONSTRAINT fk_enrollment_course
                                FOREIGN KEY (course_id)
                                    REFERENCES Course(course_id) ON DELETE CASCADE,
                            CONSTRAINT check_due_date CHECK (due_date > enrollment_date),
                            UNIQUE(student_id, course_id)
);

-- 10. Certificate table
CREATE TABLE Certificate (
                             certificate_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             enrollment_id UUID UNIQUE NOT NULL,
                             issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             certificate_url VARCHAR(500),
                             certificate_hash VARCHAR(255) UNIQUE,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             is_valid BOOLEAN DEFAULT TRUE,
                             CONSTRAINT fk_certificate_enrollment
                                 FOREIGN KEY (enrollment_id)
                                     REFERENCES enrollment(enrollment_id) ON DELETE CASCADE
);

-- 11. Rating table
CREATE TABLE Rating (
                        rating_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NOT NULL,
                        course_id UUID NOT NULL,
                        rating_count INT NOT NULL CHECK (rating_count >= 1 AND rating_count <= 5),
                        review_text TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT fk_rating_user
                            FOREIGN KEY (user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                        CONSTRAINT fk_rating_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE,
                        UNIQUE(user_id, course_id)
);

-- 12. Module table
CREATE TABLE Module (
                        module_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        module_name VARCHAR(150) NOT NULL,
                        course_id UUID NOT NULL,
                        description TEXT,
                        order_number INT NOT NULL CHECK (order_number > 0),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        is_active BOOLEAN DEFAULT TRUE,
                        CONSTRAINT fk_module_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE,
                        UNIQUE(course_id, order_number)
);

-- 13. Lesson table
CREATE TABLE Lesson (
                        lesson_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(200) NOT NULL,
                        content TEXT,
                        material_type material_type NOT NULL,
                        module_id UUID NOT NULL,
                        order_number INT NOT NULL CHECK (order_number > 0),
                        duration_minutes INT CHECK (duration_minutes > 0),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        is_active BOOLEAN DEFAULT TRUE,
                        CONSTRAINT fk_lesson_module
                            FOREIGN KEY (module_id)
                                REFERENCES Module(module_id) ON DELETE CASCADE,
                        UNIQUE(module_id, order_number)
);

-- 14. Quiz table
CREATE TABLE Quiz (
                      quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      title VARCHAR(200) NOT NULL,
                      question TEXT NOT NULL,
                      answer TEXT NOT NULL,
                      module_id UUID NOT NULL,
                      max_attempts INT DEFAULT 3 CHECK (max_attempts > 0),
                      time_limit_minutes INT CHECK (time_limit_minutes > 0),
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      is_active BOOLEAN DEFAULT TRUE,
                      CONSTRAINT fk_quiz_module
                          FOREIGN KEY (module_id)
                              REFERENCES Module(module_id) ON DELETE CASCADE
);

-- 15. Student_Quiz table
CREATE TABLE Student_Quiz (
                              student_quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              student_id UUID NOT NULL,
                              quiz_id UUID NOT NULL,
                              score DECIMAL(5,2) CHECK (score >= 0.0 AND score <= 100.0),
                              attempt INT DEFAULT 1 CHECK (attempt > 0),
                              started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              completed_at TIMESTAMP,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              CONSTRAINT fk_student_quiz_student
                                  FOREIGN KEY (student_id)
                                      REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                              CONSTRAINT fk_student_quiz_quiz
                                  FOREIGN KEY (quiz_id)
                                      REFERENCES Quiz(quiz_id) ON DELETE CASCADE,
                              CONSTRAINT check_completion_time CHECK (completed_at IS NULL OR completed_at >= started_at)
);

-- 16. Assignment table
CREATE TABLE Assignment (
                            assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            course_id UUID NOT NULL,
                            title VARCHAR(200) NOT NULL,
                            instructions TEXT,
                            due_date TIMESTAMP,
                            max_points DECIMAL(5,2) DEFAULT 100.0 CHECK (max_points > 0),
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            is_active BOOLEAN DEFAULT TRUE,
                            CONSTRAINT fk_assignment_course
                                FOREIGN KEY (course_id)
                                    REFERENCES Course(course_id) ON DELETE CASCADE,
                            CONSTRAINT check_assignment_due_date CHECK (due_date > created_at)
);

-- 17. Submission table
CREATE TABLE Submission (
                            submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            assignment_id UUID NOT NULL,
                            student_id UUID NOT NULL,
                            submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            file_url VARCHAR(500),
                            submission_text TEXT,
                            grade DECIMAL(5,2) CHECK (grade >= 0),
                            feedback TEXT,
                            graded_at TIMESTAMP,
                            graded_by UUID,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            CONSTRAINT fk_submission_assignment
                                FOREIGN KEY (assignment_id)
                                    REFERENCES Assignment(assignment_id) ON DELETE CASCADE,
                            CONSTRAINT fk_submission_student
                                FOREIGN KEY (student_id)
                                    REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                            CONSTRAINT fk_submission_grader
                                FOREIGN KEY (graded_by)
                                    REFERENCES LMS_USER(lms_user_id) ON DELETE SET NULL,
                            CONSTRAINT check_grading_time CHECK (graded_at IS NULL OR graded_at >= submitted_at),
                            UNIQUE(assignment_id, student_id)
);

-- 18. Review table
CREATE TABLE Review (
                        review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        course_id UUID NOT NULL,
                        user_id UUID NOT NULL,
                        title VARCHAR(200),
                        description TEXT,
                        is_helpful BOOLEAN DEFAULT TRUE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT fk_review_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE,
                        CONSTRAINT fk_review_user
                            FOREIGN KEY (user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                        UNIQUE(course_id, user_id)
);

-- 19. Report table
CREATE TABLE Report (
                        report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        report_name VARCHAR(200) NOT NULL,
                        generated_by_user_id UUID NOT NULL,
                        generated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        data_snapshot JSONB,
                        report_type VARCHAR(50),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        CONSTRAINT fk_report_user
                            FOREIGN KEY (generated_by_user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

-- Insert default roles
INSERT INTO LMS_USER_Role (lms_role_name) VALUES
                                              ('LMS_ADMIN'),
                                              ('STUDENT'),
                                              ('INSTRUCTOR');

-- Create a function to update the updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at columns
CREATE TRIGGER update_lms_user_role_updated_at BEFORE UPDATE ON LMS_USER_Role FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tenants_updated_at BEFORE UPDATE ON Tenants FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_lms_user_updated_at BEFORE UPDATE ON LMS_USER FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tenants_members_updated_at BEFORE UPDATE ON Tenants_Members FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_course_category_updated_at BEFORE UPDATE ON Course_Category FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_course_updated_at BEFORE UPDATE ON Course FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_enrollment_updated_at BEFORE UPDATE ON enrollment FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_certificate_updated_at BEFORE UPDATE ON Certificate FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_rating_updated_at BEFORE UPDATE ON Rating FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_module_updated_at BEFORE UPDATE ON Module FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_lesson_updated_at BEFORE UPDATE ON Lesson FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_quiz_updated_at BEFORE UPDATE ON Quiz FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_student_quiz_updated_at BEFORE UPDATE ON Student_Quiz FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_assignment_updated_at BEFORE UPDATE ON Assignment FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_submission_updated_at BEFORE UPDATE ON Submission FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_review_updated_at BEFORE UPDATE ON Review FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_report_updated_at BEFORE UPDATE ON Report FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


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
CREATE INDEX idx_enrollment_student ON enrollment(student_id);
CREATE INDEX idx_enrollment_course ON enrollment(course_id);
CREATE INDEX idx_enrollment_status ON enrollment(status);
CREATE INDEX idx_enrollment_date ON enrollment(enrollment_date);

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