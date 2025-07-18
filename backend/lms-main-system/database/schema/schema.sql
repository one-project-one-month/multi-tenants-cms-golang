eBEGIN;

-- Create ENUM types
CREATE TYPE lms_role_type AS ENUM ('ADMIN', 'INSTRUCTOR', 'STUDENT', 'VIEWER');
CREATE TYPE course_status AS ENUM ('Pending', 'Active', 'Completed', 'Cancelled');
CREATE TYPE enrollment_type AS ENUM ('Active', 'Completed', 'Dropped', 'Pending');
CREATE TYPE material_type AS ENUM ('Video', 'Document', 'Quiz', 'Assignment', 'Link');

-- 1. Create Role Type Table
CREATE TABLE LMS_USER_Role (
                               lms_role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               lms_role_name lms_role_type NOT NULL UNIQUE
);

-- 2. Create Tenants Table
CREATE TABLE Tenants (
                         tenant_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                         namespace VARCHAR(255) UNIQUE NOT NULL,
                         cms_owner_id UUID NOT NULL,
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                         is_active BOOLEAN DEFAULT TRUE
);

-- 3. Create Users Table
CREATE TABLE LMS_USER (
                          lms_user_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                          lms_user_name VARCHAR(100) NOT NULL,
                          lms_user_email VARCHAR(255) UNIQUE NOT NULL,
                          password VARCHAR(255) NOT NULL,
                          tenant_id UUID,
                          address TEXT,
                          phone_number VARCHAR(100),
                          mfa_enable BOOLEAN DEFAULT FALSE,
                          email_verified BOOLEAN DEFAULT FALSE,
                          registration_date DATE DEFAULT CURRENT_DATE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_lms_user_tenant FOREIGN KEY (tenant_id)
                              REFERENCES Tenants(tenant_id) ON DELETE SET NULL
);

-- 4. User Roles Mapping Table (Many-to-Many)
CREATE TABLE lms_user_roles_map (
                                    user_role_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    lms_user_id UUID NOT NULL,
                                    lms_role_id UUID NOT NULL,
                                    assigned_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    is_active BOOLEAN DEFAULT TRUE,
                                    lms_user_email VARCHAR(255) NOT NULL,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    CONSTRAINT fk_user_roles_user FOREIGN KEY (lms_user_id)
                                        REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                    CONSTRAINT fk_user_roles_role FOREIGN KEY (lms_role_id)
                                        REFERENCES LMS_USER_Role(lms_role_id) ON DELETE CASCADE,
                                    UNIQUE(lms_user_id, lms_role_id)


);

-- 5. Tenant Members
create table tenants_members
(
    tm_id          uuid      default gen_random_uuid() not null
        primary key,
    lms_user_id    uuid                                not null
        constraint fk_tenant_member_user
            references lms_user
            on delete cascade,
    tenant_id      uuid                                not null
        constraint fk_tenant_member_tenant
            references tenants
            on delete cascade,
    joined_date    timestamp default CURRENT_TIMESTAMP,
    is_active      boolean   default true,
    lms_user_email varchar(255)                        not null
        constraint fk_tenant_member_user_email
            references lms_user (lms_user_email)
            on delete cascade,
    constraint unique_tenant_member
        unique (lms_user_id, tenant_id, lms_user_email)
);-- 6. Course Category
CREATE TABLE Course_Category (
                                 category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 category_name VARCHAR(100) NOT NULL,
                                 description TEXT,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP
);

-- 7. Course
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
                        CONSTRAINT fk_course_category FOREIGN KEY (course_category)
                            REFERENCES Course_Category(category_id),
                        CONSTRAINT fk_course_instructor FOREIGN KEY (instructor_id)
                            REFERENCES LMS_USER(lms_user_id),
                        CONSTRAINT fk_course_tenant FOREIGN KEY (owned_by)
                            REFERENCES Tenants(tenant_id)
);

-- 8. Namespace Consumer
CREATE TABLE namespace_consumer (
                                    consumer_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                    lms_user_id UUID NOT NULL,
                                    namespace VARCHAR(255) NOT NULL DEFAULT 'default_consumer_namespace',
                                    joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    is_active BOOLEAN DEFAULT TRUE,
                                    CONSTRAINT fk_namespace_consumer_user FOREIGN KEY (lms_user_id)
                                        REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                    UNIQUE(lms_user_id, namespace)
);

-- 9. Enrollment
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
                            CONSTRAINT fk_enrollment_student FOREIGN KEY (student_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                            CONSTRAINT fk_enrollment_course FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 10. Certificate
CREATE TABLE Certificate (
                             certificate_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                             enrollment_id UUID UNIQUE NOT NULL,
                             issue_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             certificate_url VARCHAR(500),
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             updated_at TIMESTAMP,
                             CONSTRAINT fk_certificate_enrollment FOREIGN KEY (enrollment_id)
                                 REFERENCES enrollment(enrollment_id) ON DELETE CASCADE
);

-- 11. Rating
CREATE TABLE Rating (
                        rating_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        user_id UUID NOT NULL,
                        course_id UUID NOT NULL,
                        rating_count INT CHECK (rating_count >= 1 AND rating_count <= 5),
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_rating_user FOREIGN KEY (user_id)
                            REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                        CONSTRAINT fk_rating_course FOREIGN KEY (course_id)
                            REFERENCES Course(course_id) ON DELETE CASCADE,
                        UNIQUE(user_id, course_id)
);

-- 12. Module
CREATE TABLE Module (
                        module_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        module_name VARCHAR(150) NOT NULL,
                        course_id UUID NOT NULL,
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_module_course FOREIGN KEY (course_id)
                            REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 13. Quiz
CREATE TABLE Quiz (
                      quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                      question TEXT NOT NULL,
                      answer TEXT NOT NULL,
                      module_id UUID NOT NULL,
                      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMP,
                      CONSTRAINT fk_quiz_module FOREIGN KEY (module_id)
                          REFERENCES Module(module_id) ON DELETE CASCADE
);

-- 14. Student Quiz
CREATE TABLE Student_Quiz (
                              student_quiz_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                              student_id UUID NOT NULL,
                              quiz_id UUID NOT NULL,
                              score INT,
                              attempt INT DEFAULT 1,
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              updated_at TIMESTAMP,
                              CONSTRAINT fk_student_quiz_student FOREIGN KEY (student_id)
                                  REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                              CONSTRAINT fk_student_quiz_quiz FOREIGN KEY (quiz_id)
                                  REFERENCES Quiz(quiz_id) ON DELETE CASCADE
);

-- 15. Lesson
CREATE TABLE Lesson (
                        lesson_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        title VARCHAR(200) NOT NULL,
                        content TEXT,
                        material_type material_type,
                        module_id UUID NOT NULL,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_lesson_module FOREIGN KEY (module_id)
                            REFERENCES Module(module_id) ON DELETE CASCADE
);

-- 16. Assignment
CREATE TABLE Assignment (
                            assignment_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            course_id UUID NOT NULL,
                            title VARCHAR(200) NOT NULL,
                            instructions TEXT,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP,
                            CONSTRAINT fk_assignment_course FOREIGN KEY (course_id)
                                REFERENCES Course(course_id) ON DELETE CASCADE
);

-- 17. Submission
CREATE TABLE Submission (
                            submission_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                            assignment_id UUID NOT NULL,
                            student_id UUID NOT NULL,
                            submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            file_url VARCHAR(500),
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            updated_at TIMESTAMP,
                            CONSTRAINT fk_submission_assignment FOREIGN KEY (assignment_id)
                                REFERENCES Assignment(assignment_id) ON DELETE CASCADE,
                            CONSTRAINT fk_submission_student FOREIGN KEY (student_id)
                                REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

-- 18. Review
CREATE TABLE Review (
                        review_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        course_id UUID NOT NULL,
                        user_id UUID NOT NULL,
                        title VARCHAR(200),
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_review_course FOREIGN KEY (course_id)
                            REFERENCES Course(course_id) ON DELETE CASCADE,
                        CONSTRAINT fk_review_user FOREIGN KEY (user_id)
                            REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

-- 19. Report
CREATE TABLE Report (
                        report_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                        report_name VARCHAR(200) NOT NULL,
                        generated_by_user_id UUID NOT NULL,
                        generated_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        data_snapshot TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP,
                        CONSTRAINT fk_report_user FOREIGN KEY (generated_by_user_id)
                            REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE
);

COMMIT;

-- Create indexes for better performance
CREATE INDEX idx_lms_user_email ON LMS_USER(lms_user_email);
CREATE INDEX idx_lms_user_created_at ON LMS_USER(created_at);
CREATE INDEX idx_lms_user_registration_date ON LMS_USER(registration_date);
CREATE INDEX idx_lms_user_phone ON LMS_USER(phone_number);
CREATE INDEX idx_lms_user_tenant ON LMS_USER(tenant_id);

-- Indexes for user roles mapping
CREATE INDEX idx_user_roles_map_user ON lms_user_roles_map(lms_user_id);
CREATE INDEX idx_user_roles_map_role ON lms_user_roles_map(lms_role_id);
CREATE INDEX idx_user_roles_map_active ON lms_user_roles_map(is_active);
CREATE INDEX idx_user_roles_map_user_active ON lms_user_roles_map(lms_user_id, is_active);

CREATE INDEX idx_tenants_namespace ON Tenants(namespace);
CREATE INDEX idx_tenants_owner ON Tenants(cms_owner_id);
CREATE INDEX idx_tenants_active ON Tenants(is_active);
CREATE INDEX idx_tenants_created_at ON Tenants(created_at);
CREATE INDEX idx_tenants_active_true ON Tenants(tenant_id) WHERE is_active = TRUE;

CREATE INDEX idx_tenant_members_user ON Tenants_Members(lms_user_id);
CREATE INDEX idx_tenant_members_tenant ON Tenants_Members(tenant_id);
CREATE INDEX idx_tenant_members_joined_date ON Tenants_Members(joined_date);
CREATE INDEX idx_tenant_members_active ON Tenants_Members(is_active);
CREATE INDEX idx_tenant_members_active_user ON Tenants_Members(is_active, lms_user_id);
CREATE INDEX idx_tenant_members_active_tenant ON Tenants_Members(is_active, tenant_id);
CREATE INDEX idx_tenant_members_active_true ON Tenants_Members(lms_user_id, tenant_id) WHERE is_active = TRUE;

CREATE INDEX idx_course_instructor ON Course(instructor_id);
CREATE INDEX idx_course_category ON Course(course_category);
CREATE INDEX idx_course_owner ON Course(owned_by);
CREATE INDEX idx_course_title ON Course(course_title);
CREATE INDEX idx_course_created_at ON Course(created_at);
CREATE INDEX idx_course_rating ON Course(overall_rating);
CREATE INDEX idx_course_status ON Course(status);
CREATE INDEX idx_course_instructor_category ON Course(instructor_id, course_category);
CREATE INDEX idx_course_owner_category ON Course(owned_by, course_category);

CREATE INDEX idx_category_name ON Course_Category(category_name);
CREATE INDEX idx_category_created_at ON Course_Category(created_at);

CREATE INDEX idx_namespace_consumer_user ON namespace_consumer(lms_user_id);
CREATE INDEX idx_namespace_consumer_namespace ON namespace_consumer(namespace);
CREATE INDEX idx_namespace_consumer_active ON namespace_consumer(is_active);
CREATE INDEX idx_namespace_consumer_joined_date ON namespace_consumer(joined_date);

CREATE INDEX idx_enrollment_student ON enrollment(student_id);
CREATE INDEX idx_enrollment_course ON enrollment(course_id);
CREATE INDEX idx_enrollment_status ON enrollment(status);
CREATE INDEX idx_enrollment_date ON enrollment(enrollment_date);

CREATE INDEX idx_rating_user ON Rating(user_id);
CREATE INDEX idx_rating_course ON Rating(course_id);
CREATE INDEX idx_rating_count ON Rating(rating_count);

CREATE INDEX idx_module_course ON Module(course_id);
CREATE INDEX idx_module_name ON Module(module_name);

CREATE INDEX idx_quiz_module ON Quiz(module_id);

CREATE INDEX idx_student_quiz_student ON Student_Quiz(student_id);
CREATE INDEX idx_student_quiz_quiz ON Student_Quiz(quiz_id);
CREATE INDEX idx_student_quiz_score ON Student_Quiz(score);

CREATE INDEX idx_lesson_module ON Lesson(module_id);
CREATE INDEX idx_lesson_title ON Lesson(title);

CREATE INDEX idx_assignment_course ON Assignment(course_id);

CREATE INDEX idx_submission_assignment ON Submission(assignment_id);
CREATE INDEX idx_submission_student ON Submission(student_id);
CREATE INDEX idx_submission_date ON Submission(submitted_at);

CREATE INDEX idx_review_course ON Review(course_id);
CREATE INDEX idx_review_user ON Review(user_id);

CREATE INDEX idx_report_generated_by ON Report(generated_by_user_id);
CREATE INDEX idx_report_date ON Report(generated_date);

-- Create function to prevent students from being tenant members
CREATE OR REPLACE FUNCTION prevent_student_tenant_membership()
    RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM lms_user_roles_map urm
        JOIN LMS_USER_Role r ON r.lms_role_id = urm.lms_role_id
        WHERE urm.lms_user_id = NEW.lms_user_id
          AND r.lms_role_name = 'STUDENT'
          AND urm.is_active = TRUE
    ) THEN
        RAISE EXCEPTION 'Students cannot be tenant members';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to auto assign student to namespace_consumer when role is assigned
CREATE OR REPLACE FUNCTION auto_assign_student_to_namespace_consumer()
    RETURNS TRIGGER AS $$
DECLARE
student_namespace VARCHAR(255);
    user_tenant_id UUID;
BEGIN
    -- Check if the role being assigned is STUDENT
    IF EXISTS (
        SELECT 1
        FROM LMS_USER_Role r
        WHERE r.lms_role_id = NEW.lms_role_id
          AND r.lms_role_name = 'STUDENT'
    ) THEN
        -- Get user's tenant_id
SELECT tenant_id INTO user_tenant_id
FROM LMS_USER
WHERE lms_user_id = NEW.lms_user_id;

-- If user has a tenant, get the namespace
IF user_tenant_id IS NOT NULL THEN
SELECT namespace INTO student_namespace
FROM Tenants
WHERE tenant_id = user_tenant_id
  AND is_active = TRUE;

-- Insert into namespace_consumer if namespace found
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

-- Create trigger to prevent student tenant membership
CREATE TRIGGER trg_prevent_student_tenant_membership
    BEFORE INSERT OR UPDATE ON Tenants_Members
                         FOR EACH ROW
                         EXECUTE FUNCTION prevent_student_tenant_membership();

-- Create trigger to auto assign student to namespace_consumer when role is assigned
CREATE TRIGGER trg_auto_assign_student_role
    AFTER INSERT ON lms_user_roles_map
    FOR EACH ROW
    EXECUTE FUNCTION auto_assign_student_to_namespace_consumer();