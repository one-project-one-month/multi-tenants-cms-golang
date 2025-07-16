-- Drop triggers
DROP TRIGGER IF EXISTS trg_prevent_student_tenant_membership ON Tenants_Members;
DROP TRIGGER IF EXISTS trg_auto_assign_student_insert ON LMS_USER;

-- Drop functions
DROP FUNCTION IF EXISTS prevent_student_tenant_membership();
DROP FUNCTION IF EXISTS auto_assign_student_to_namespace_consumer();
