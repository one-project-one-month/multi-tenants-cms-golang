BEGIN;

DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'course_status') THEN
DROP TYPE course_status;
END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'material_type') THEN
DROP TYPE material_type;
END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'enrollment_type') THEN
DROP TYPE enrollment_type;
END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'system_type') THEN
DROP TYPE system_type;
END IF;

    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'lms_role_type') THEN
DROP TYPE lms_role_type;
END IF;
END
$$;

COMMIT;
