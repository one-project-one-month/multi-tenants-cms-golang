package module

import (
	db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"
	mpb "github.com/multi-tenants-cms-golang/lms-sys/protogen/modules"
	"github.com/sirupsen/logrus"
)

type ModulesService struct {
	mpb.UnimplementedModuleServiceServer
	store  *db.Store
	logger *logrus.Logger
}

func NewModuleService(store *db.Store, logger *logrus.Logger) *ModulesService {
	return &ModulesService{
		store:  store,
		logger: logger,
	}
}
