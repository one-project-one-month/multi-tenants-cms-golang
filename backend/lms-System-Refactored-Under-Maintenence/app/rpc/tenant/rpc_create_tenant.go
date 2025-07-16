package tenant

import (
	"context"
	tenantpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *TenantsService) CreateTenant(ctx context.Context, req *tenantpb.CreateTenantRequest) (*tenantpb.CreateTenantResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request cannot be nil")
	}
	resp := &tenantpb.CreateTenantResponse{}
	return resp, nil
}
