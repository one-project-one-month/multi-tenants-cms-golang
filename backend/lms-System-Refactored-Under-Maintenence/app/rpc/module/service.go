package module

import (
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
)

type ModulesService struct {
	modulepb.UnimplementedModuleServiceServer
}

func NewModuleService() *ModulesService {
	return &ModulesService{}
}
