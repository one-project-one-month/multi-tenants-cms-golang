BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'lms_role_type') THEN
CREATE TYPE lms_role_type AS ENUM ('LMS_ADMIN', 'STUDENT', 'INSTRUCTOR');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'system_type') THEN
CREATE TYPE system_type AS ENUM ('LMS', 'EMS');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'enrollment_type') THEN
CREATE TYPE enrollment_type AS ENUM('ENROLLED', 'COMPLETED', 'DROPPED');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_type') THEN
CREATE TYPE material_type AS ENUM('Video', 'PDF', 'Slide', 'Link');
END IF;

    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'course_status') THEN
CREATE TYPE course_status AS ENUM('Pending', 'Published', 'Unpublished', 'Archived');
END IF;
END
$$;

COMMIT;
