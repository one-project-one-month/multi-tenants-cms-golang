package module

import (
	modulepb "github.com/multi-tenants-cms-golang/lms-sys/protogen/module"
)

type ModuleService struct {
	modulepb.UnimplementedModuleServiceServer
}

func NewModuleService() *ModuleService {
	return &ModuleService{}
}
