package module

import db "github.com/multi-tenants-cms-golang/lms-sys/internal/repo"

type ModulesService struct {
	mpb.UnimplementedModulesServer
	store *db.Store
}
