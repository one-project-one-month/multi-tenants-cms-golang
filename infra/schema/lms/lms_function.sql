-- Create functions
CREATE OR REPLACE FUNCTION prevent_student_tenant_membership()
    RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS(
        SELECT 1
        FROM LMS_USER u
                 JOIN LMS_USER_Role LUR on LUR.lms_role_id = u.lms_role_id
        WHERE u.lms_user_id = NEW.lms_user_id
          AND LUR.lms_role_name = 'STUDENT'
    ) THEN
        RAISE EXCEPTION 'Student cannot be the tenant members';
END IF;
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION auto_assign_student_to_namespace_consumer()
    RETURNS TRIGGER AS $$
DECLARE
student_namespace varchar(255);
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
CREATE TRIGGER trg_prevent_student_tenant_membership
    BEFORE INSERT OR UPDATE ON Tenants_Members
                         FOR EACH ROW
                         EXECUTE FUNCTION prevent_student_tenant_membership();

CREATE TRIGGER trg_auto_assign_student_insert
    AFTER INSERT ON LMS_USER
    FOR EACH ROW
    EXECUTE FUNCTION auto_assign_student_to_namespace_consumer();