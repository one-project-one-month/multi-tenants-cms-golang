package tenants

import (
	tenantspb "github.com/multi-tenants-cms-golang/lms-system/protogen/tenants"
)

type TenantsService struct {
	tenantspb.UnimplementedTenantsServiceServer
}

func NewTenantsService() *TenantsService {
	return &TenantsService{}
}
