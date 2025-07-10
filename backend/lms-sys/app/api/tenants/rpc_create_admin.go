package tenants

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/db"
	tpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenants"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (t *TenantService) CreateTenant(ctx context.Context, req *tpb.CreateTenantRequest) (*tpb.CreateTenantResponse, error) {
	// Validate required fields
	if req.GetNamespace() == "" {
		return nil, status.Error(codes.InvalidArgument, "namespace is required")
	}

	if req.GetCmsOwnerId() == "" {
		return nil, status.Error(codes.InvalidArgument, "cms_owner_id is required")
	}

	// Parse UUID from string
	ownerUUID, err := uuid.Parse(req.GetCmsOwnerId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid cms_owner_id format: %v", err)
	}

	// Create pgtype.UUID
	pgUUID := pgtype.UUID{
		Bytes: ownerUUID,
		Valid: true,
	}

	tenant, err := t.store.CreateNewTenant(ctx, db.CreateNewTenantParams{
		Namespace:  req.GetNamespace(),
		CmsOwnerID: pgUUID,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create tenant: %v", err)
	}

	var createdAtPb *timestamppb.Timestamp
	if tenant.CreatedAt.Valid {
		createdAtPb = timestamppb.New(tenant.CreatedAt.Time)
	} else {
		createdAtPb = timestamppb.New(time.Time{})
	}

	tenantIDStr := uuid.UUID(tenant.TenantID.Bytes).String()
	ownerIDStr := uuid.UUID(tenant.CmsOwnerID.Bytes).String()

	return &tpb.CreateTenantResponse{
		CreatedTenant: &tpb.Tenants{
			TenantsId:  tenantIDStr,
			Namespace:  tenant.Namespace,
			CmsOwnerId: ownerIDStr,
			CreatedAt:  createdAtPb,
		},
	}, nil
}
