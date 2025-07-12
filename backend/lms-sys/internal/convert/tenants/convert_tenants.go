package tenants

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/multi-tenants-cms-golang/lms-sys/internal/convert/global"
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	tpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenants"
)

func ConvertTenantsSQLTOProto(tenants db.Tenant) *tpb.CreateTenantResponse {
	return &tpb.CreateTenantResponse{
		CreatedTenant: &tpb.Tenants{
			TenantsId:  tenants.TenantID.String(),
			Namespace:  tenants.Namespace,
			CmsOwnerId: tenants.CmsOwnerID.String(),
			CreatedAt:  global.ConvertTimestamp(pgtype.Timestamptz(tenants.CreatedAt)),
			IsActive:   tenants.IsActive.Bool,
		},
	}
}
