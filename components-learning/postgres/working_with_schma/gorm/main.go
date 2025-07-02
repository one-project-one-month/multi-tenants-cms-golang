package main

import (
	"fmt"
	"gorm.io/gorm/logger"
	"os"
	"os/signal"
	"time"

	"github.com/SwanHtetAungPhyo/multi-tenants/tenant_models"
	"github.com/SwanHtetAungPhyo/multi-tenants/tenants"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=postgres user=cms_user password=cms_pass dbname=cms_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil
	}
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	return db
}

func insertItemsForTenant(db *gorm.DB, schema string, item tenant_models.Item) error {

	if err := db.Table(fmt.Sprintf("%s.items", schema)).Create(&item).Error; err != nil {
		return err
	}

	return nil
}

func main() {
	database := InitDB()

	tenantList := []tenants.TenantList{
		{NameSpace: "swan"},
		{NameSpace: "kyaw"},
		{NameSpace: "kaung"},
	}

	err := database.Transaction(func(tx *gorm.DB) error {
		var items []tenant_models.Item

		for _, tenant := range tenantList {
			if err := tenants.CreateTenantSchema(tx, tenant.NameSpace); err != nil {
				return fmt.Errorf("failed to create schema %s: %w", tenant.NameSpace, err)
			}

			if err := tenants.MigrateTenantModels(tx, tenant.NameSpace); err != nil {
				return fmt.Errorf("failed to migrate models for %s: %w", tenant.NameSpace, err)
			}

			item := tenant_models.Item{
				Name: fmt.Sprintf("Sample item for %s", tenant.NameSpace),
			}
			items = append(items, item)
		}
		for i, tenant := range tenantList {
			if err := insertItemsForTenant(tx, tenant.NameSpace, items[i]); err != nil {
				return fmt.Errorf("failed to insert item for %s: %w", tenant.NameSpace, err)
			}
		}

		return nil
	})

	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Schemas created, models migrated, and items inserted per tenant.")

	osChan := make(chan os.Signal, 1)
	signal.Notify(osChan, os.Interrupt)
	<-osChan
}
