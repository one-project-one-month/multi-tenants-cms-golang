package tenants

import (
	database "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	tbp "github.com/multi-tenants-cms-golang/lms-sys/protogen/tenants"
	"github.com/sirupsen/logrus"
)

type TenantService struct {
	tbp.UnimplementedTenantServiceServer
	store  *database.Store
	logger *logrus.Logger
}

var _ tbp.TenantServiceServer = (*TenantService)(nil)

func NewTenantService(
	store *database.Store,
	logger *logrus.Logger,
) *TenantService {
	return &TenantService{
		store:  store,
		logger: logger,
	}
}
