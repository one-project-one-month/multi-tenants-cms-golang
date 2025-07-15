package tenants

import (
	"context"
	tenantspb "github.com/multi-tenants-cms-golang/lms-system/protogen/tenants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TenantsService) CreateTenant(ctx context.Context, req *tenantspb.CreateTenantRequest) (*tenantspb.CreateTenantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &tenantspb.CreateTenantResponse{}
	return resp, nil
}
