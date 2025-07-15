package tenant

import (
	tenantpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenant"
)

type TenantService struct {
	tenantpb.UnimplementedTenantServiceServer
}

func NewTenantService() *TenantService {
	return &TenantService{}
}
