package tenant

import (
	tenantpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenant"
)

type TenantsService struct {
	tenantpb.UnimplementedTenantServiceServer
}

func NewTenantService() *TenantsService {
	return &TenantsService{}
}
