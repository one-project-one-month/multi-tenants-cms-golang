-- name: CreateNewTenant :one
INSERT INTO tenants (namespace, cms_owner_id)
VALUES ($1, $2)
RETURNING tenant_id, namespace, cms_owner_id, created_at, is_active;

-- name: GetTenantByNamespace :one
SELECT  * FROM tenants
WHERE  namespace = $1 AND  is_active = true;

-- name: GetTenantByID :one
SELECT tenant_id, namespace, cms_owner_id, created_at, is_active
FROM Tenants
WHERE tenant_id = $1 AND is_active = TRUE;