BEGIN;

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
                          lms_user_email VARCHAR(255) UNIQUE NOT NULL,
                          password VARCHAR(255) NOT NULL,
                          lms_role_id UUID NOT NULL,
                          tenant_id UUID,
                          address TEXT,
                          phone_number VARCHAR(100),
                          registration_date DATE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          CONSTRAINT fk_lms_user_role FOREIGN KEY (lms_role_id)
                              REFERENCES LMS_USER_Role(lms_role_id) ON DELETE RESTRICT,
                          CONSTRAINT fk_lms_user_tenant FOREIGN KEY (tenant_id)
                              REFERENCES Tenants(tenant_id) ON DELETE SET NULL
);

-- 4. Tenant Members
CREATE TABLE Tenants_Members (
                                 tm_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 lms_user_id UUID NOT NULL,
                                 tenant_id UUID NOT NULL,
                                 joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 is_active BOOLEAN DEFAULT TRUE,
                                 CONSTRAINT fk_tenant_member_user FOREIGN KEY (lms_user_id)
                                     REFERENCES LMS_USER(lms_user_id) ON DELETE CASCADE,
                                 CONSTRAINT fk_tenant_member_tenant FOREIGN KEY (tenant_id)
                                     REFERENCES Tenants(tenant_id) ON DELETE CASCADE,
                                 UNIQUE(lms_user_id, tenant_id)
);

-- 5. Course Category
CREATE TABLE Course_Category (
                                 category_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                                 category_name VARCHAR(100) NOT NULL,
                                 description TEXT,
                                 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                 updated_at TIMESTAMP
);

-- 6. Course
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

-- 7. Namespace Consumer
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

-- 8. Enrollment
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

-- 9. Certificate
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

-- 10. Rating
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

-- 11. Module
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

-- 12. Quiz
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

-- 13. Student Quiz
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

-- 14. Lesson
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

-- 15. Assignment
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

-- 16. Submission
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

-- 17. Review
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

-- 18. Report
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
