-- Create function to prevent students from being tenant members
CREATE OR REPLACE FUNCTION prevent_student_tenant_membership()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM LMS_USER u
        JOIN LMS_USER_Role LUR ON LUR.lms_role_id = u.lms_role_id
        WHERE u.lms_user_id = NEW.lms_user_id
        AND LUR.lms_role_name = 'STUDENT'
    ) THEN
        RAISE EXCEPTION 'Student cannot be the tenant members';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create function to auto assign student to namespace_consumer
CREATE OR REPLACE FUNCTION auto_assign_student_to_namespace_consumer()
RETURNS TRIGGER AS $$
DECLARE
student_namespace VARCHAR(255);
BEGIN
    IF EXISTS (
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

-- Create trigger to prevent student tenant membership
CREATE TRIGGER trg_prevent_student_tenant_membership
    BEFORE INSERT OR UPDATE ON Tenants_Members
                         FOR EACH ROW
                         EXECUTE FUNCTION prevent_student_tenant_membership();

-- Create trigger to auto assign student to namespace_consumer
CREATE TRIGGER trg_auto_assign_student_insert
    AFTER INSERT ON LMS_USER
    FOR EACH ROW
    EXECUTE FUNCTION auto_assign_student_to_namespace_consumer();
