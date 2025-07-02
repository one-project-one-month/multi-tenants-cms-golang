-- Create schema and set search path
CREATE SCHEMA IF NOT EXISTS learning_management;
SET search_path TO learning_management;

-- Create custom types
CREATE TYPE lms_role_type AS ENUM ('LMS_ADMIN', 'STUDENT', 'INSTRUCTOR');
CREATE TYPE system_type AS ENUM ('LMS', 'EMS');
CREATE TYPE enrollment_type AS ENUM('ENROLLED', 'COMPLETED', 'DROPPED');
CREATE TYPE material_type AS ENUM('Video', 'PDF', 'Slide', 'Link');
CREATE TYPE course_status AS ENUM('Pending', 'Published', 'Unpublished', 'Archived');

-- Create tables in proper dependency order

-- 1. Role table (no dependencies)
CREATE TABLE LMS_USER_Role (
                               lms_role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               lms_role_name lms_role_type NOT NULL UNIQUE
);

-- 2. Tenants table (no dependencies)
CREATE TABLE Tenants (
                         tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         namespace VARCHAR(255) UNIQUE NOT NULL,
                         cms_owner_id UUID NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         is_active BOOLEAN DEFAULT TRUE
);

-- 3. User table (depends on Role and Tenants)
CREATE TABLE LMS_USER (
                          lms_user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          lms_user_email VARCHAR(255) UNIQUE NOT NULL,
                          password VARCHAR(255) NOT NULL,
                          lms_role_id UUID NOT NULL,
                          tenant_id UUID,
                          address TEXT,
                          phone_number VARCHAR(100),
                          registration_date DATE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_lms_user_role
                              FOREIGN KEY (lms_role_id)
                                  REFERENCES LMS_USER_Role(lms_role_id) ON DELETE RESTRICT,
                          CONSTRAINT fk_lms_user_tenant
                              FOREIGN KEY (tenant_id)
                                  REFERENCES Tenants(tenant_id) ON DELETE SET NULL
);

-- 4. Tenants_Members table (depends on User and Tenants)
CREATE TABLE Tenants_Members (
                                 tm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 lms_user_id UUID NOT NULL,
                                 tenant_id UUID NOT NULL,
                                 joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 is_active BOOLEAN DEFAULT TRUE,
                                 CONSTRAINT fk_tenant_member_user
                                     FOREIGN KEY (lms_user_id)
                                         REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                 CONSTRAINT fk_tenant_member_tenant
                                     FOREIGN KEY (tenant_id)
                                         REFERENCES Tenants(tenant_id) ON DELETE CASCADE,
                                 UNIQUE(lms_user_id, tenant_id)
);

-- 5. Course Category table (no dependencies)
CREATE TABLE Course_Category (
                                 category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 category_name VARCHAR(100) NOT NULL,
                                 description TEXT,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP
);

-- 6. Course table (depends on User, Tenants, and Course_Category)
CREATE TABLE Course (
                        course_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        course_title VARCHAR(150) NOT NULL,
                        description TEXT,
                        instructor_id UUID NOT NULL,
                        overall_rating INT,
                        course_category UUID NOT NULL,
                        status course_status DEFAULT 'Pending',
                        duration_day_count INT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        owned_by UUID NOT NULL,
                        CONSTRAINT fk_course_category
                            FOREIGN KEY (course_category)
                                REFERENCES Course_Category(category_id),
                        CONSTRAINT fk_course_instructor
                            FOREIGN KEY (instructor_id)
                                REFERENCES LMS_USER(lms_user_id),
                        CONSTRAINT fk_course_tenant
                            FOREIGN KEY (owned_by)
                                REFERENCES Tenants(tenant_id)
);

-- 7. Namespace Consumer table (depends on User)
CREATE TABLE namespace_consumer (
                                    consumer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    lms_user_id UUID NOT NULL,
                                    namespace VARCHAR(255) NOT NULL DEFAULT 'default_consumer_namespace',
                                    joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    is_active BOOLEAN DEFAULT TRUE,
                                    CONSTRAINT fk_namespace_consumer_user
                                        FOREIGN KEY (lms_user_id)
                                            REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                    UNIQUE(lms_user_id, namespace)
);

-- 8. Enrollment table
CREATE TABLE enrollment (
                            enrollment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            student_id UUID NOT NULL,
                            course_id UUID NOT NULL,
                            enrollment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            progress DECIMAL,
                            status enrollment_type NOT NULL,
                            due_date TIMESTAMP,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP,
                            CONSTRAINT fk_enrollment_student
                                FOREIGN KEY (student_id)
                                    REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                            CONSTRAINT fk_enrollment_course
                                FOREIGN KEY (course_id)
                                    REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 9. Certificate table
CREATE TABLE Certificate (
                             certificate_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             enrollment_id UUID UNIQUE NOT NULL,
                             issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             certificate_url VARCHAR(500),
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP,
                             CONSTRAINT fk_certificate_enrollment
                                 FOREIGN KEY (enrollment_id)
                                     REFERENCES enrollment(enrollment_id) ON DELETE CASCADE
);

-- 10. Rating table
CREATE TABLE Rating (
                        rating_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NOT NULL,
                        course_id UUID NOT NULL,
                        rating_count INT CHECK (rating_count >= 1 AND rating_count <= 5),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_rating_user
                            FOREIGN KEY (user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                        CONSTRAINT fk_rating_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE,
                        UNIQUE(user_id, course_id)
);

-- 11. Module table
CREATE TABLE Module (
                        module_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        module_name VARCHAR(150) NOT NULL,
                        course_id UUID NOT NULL,
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_module_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 12. Quiz table
CREATE TABLE Quiz (
                      quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      question TEXT NOT NULL,
                      answer TEXT NOT NULL,
                      module_id UUID NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP,
                      CONSTRAINT fk_quiz_module
                          FOREIGN KEY (module_id)
                              REFERENCES Module(module_id) ON DELETE CASCADE
);

-- 13. Student_Quiz table
CREATE TABLE Student_Quiz (
                              student_quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              student_id UUID NOT NULL,
                              quiz_id UUID NOT NULL,
                              score INT,
                              attempt INT DEFAULT 1,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP,
                              CONSTRAINT fk_student_quiz_student
                                  FOREIGN KEY (student_id)
                                      REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                              CONSTRAINT fk_student_quiz_quiz
                                  FOREIGN KEY (quiz_id)
                                      REFERENCES Quiz(quiz_id) ON DELETE CASCADE
);

-- 14. Lesson table
CREATE TABLE Lesson (
                        lesson_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(200) NOT NULL,
                        content TEXT,
                        material_type material_type,
                        module_id UUID NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_lesson_module
                            FOREIGN KEY (module_id)
                                REFERENCES Module(module_id) ON DELETE CASCADE
);

-- 15. Assignment table
CREATE TABLE Assignment (
                            assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            course_id UUID NOT NULL,
                            title VARCHAR(200) NOT NULL,
                            instructions TEXT,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP,
                            CONSTRAINT fk_assignment_course
                                FOREIGN KEY (course_id)
                                    REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 16. Submission table
CREATE TABLE Submission (
                            submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            assignment_id UUID NOT NULL,
                            student_id UUID NOT NULL,
                            submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            file_url VARCHAR(500),
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP,
                            CONSTRAINT fk_submission_assignment
                                FOREIGN KEY (assignment_id)
                                    REFERENCES Assignment(assignment_id) ON DELETE CASCADE,
                            CONSTRAINT fk_submission_student
                                FOREIGN KEY (student_id)
                                    REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

-- 17. Review table
CREATE TABLE Review (
                        review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        course_id UUID NOT NULL,
                        user_id UUID NOT NULL,
                        title VARCHAR(200),
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_review_course
                            FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE,
                        CONSTRAINT fk_review_user
                            FOREIGN KEY (user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

-- 18. Report table
CREATE TABLE Report (
                        report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        report_name VARCHAR(200) NOT NULL,
                        generated_by_user_id UUID NOT NULL,
                        generated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        data_snapshot TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_report_user
                            FOREIGN KEY (generated_by_user_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

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

-- Create trigger functions

-- Function to prevent student tenant membership
CREATE OR REPLACE FUNCTION prevent_student_tenant_membership()
    RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM LMS_USER u
                 JOIN LMS_USER_Role r ON u.lms_role_id = r.lms_role_id
        WHERE u.lms_user_id = NEW.lms_user_id
          AND r.lms_role_name = 'STUDENT'
    ) THEN
        RAISE EXCEPTION 'Students cannot be tenant members';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Function to auto-assign students to namespace consumer
CREATE OR REPLACE FUNCTION auto_assign_student_to_namespace_consumer()
    RETURNS TRIGGER AS $$
DECLARE
    student_namespace VARCHAR(255);
BEGIN
    IF EXISTS(
        SELECT 1
        FROM LMS_USER_Role r
        WHERE r.lms_role_id = NEW.lms_role_id
          AND r.lms_role_name = 'STUDENT'
    ) THEN
        IF NEW.tenant_id IS NOT NULL THEN
            SELECT namespace INTO student_namespace
            FROM Tenants
            WHERE tenant_id = NEW.tenant_id
              AND is_active = TRUE;

            IF student_namespace IS NOT NULL THEN
                INSERT INTO namespace_consumer (lms_user_id, namespace)
                VALUES (NEW.lms_user_id, student_namespace)
                ON CONFLICT (lms_user_id, namespace) DO NOTHING;
            END IF;
        END IF;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers

-- Trigger to prevent student tenant membership
CREATE TRIGGER trg_prevent_student_tenant_membership
    BEFORE INSERT OR UPDATE ON Tenants_Members
    FOR EACH ROW
EXECUTE FUNCTION prevent_student_tenant_membership();

-- Trigger to auto-assign students to namespace consumer
CREATE TRIGGER trg_auto_assign_student_insert
    AFTER INSERT ON LMS_USER
    FOR EACH ROW
EXECUTE FUNCTION auto_assign_student_to_namespace_consumer();
