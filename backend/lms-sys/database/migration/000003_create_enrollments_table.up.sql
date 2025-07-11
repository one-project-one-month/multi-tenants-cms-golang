CREATE TYPE enrollment_status AS ENUM (
  'ENROLLED',
  'IN_PROGRESS',
  'COMPLETED',
  'DROPPED',
  'WITHDRAWN'
);

CREATE TABLE "Enrollments" (
                               "enrollment_id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
                               "student_id" UUID NOT NULL,
                               "course_id" UUID NOT NULL,
                               "enrollment_date" TIMESTAMPTZ NOT NULL DEFAULT now(),
                               "progress" DECIMAL(5,2) NOT NULL DEFAULT 0.00,
                               "status" enrollment_status NOT NULL DEFAULT 'ENROLLED',
                               "due_date" TIMESTAMPTZ,
                               "created_at" TIMESTAMPTZ NOT NULL DEFAULT now(),
                               "updated_at" TIMESTAMPTZ,

                               CONSTRAINT fk_student
                                   FOREIGN KEY ("student_id")
                                       REFERENCES "LMS_USER" ("user_id") ON DELETE CASCADE,
                               CONSTRAINT fk_course
                                   FOREIGN KEY ("course_id")
                                       REFERENCES "Course" ("course_id") ON DELETE CASCADE,
                               CONSTRAINT unique_student_course UNIQUE ("student_id", "course_id")
);

-- Add indexes for foreign keys
CREATE INDEX IF NOT EXISTS idx_enrollments_student_id ON "Enrollments" ("student_id");
CREATE INDEX IF NOT EXISTS idx_enrollments_course_id ON "Enrollments" ("course_id");