package tenants

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/multi-tenants/tenant_models"
	"gorm.io/gorm"
)

func CreateTenantSchema(db *gorm.DB, namespace string) error {
	createSchemaSQL := fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, namespace)
	return db.Exec(createSchemaSQL).Error
}

func MigrateTenantModels(db *gorm.DB, namespace string) error {
	return db.Table(fmt.Sprintf("%s.items", namespace)).AutoMigrate(&tenant_models.Item{})
}
