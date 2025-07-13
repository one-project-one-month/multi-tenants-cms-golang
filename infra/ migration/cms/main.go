package main

import (
	"encoding/json"
	"fmt"
	"github.com/SwanHtetAungPhyo/migration/query"
	"github.com/SwanHtetAungPhyo/migration/types"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
)

func main() {

	dsn := "host=ep-cool-poetry-a2e2a9w7-pooler.eu-central-1.aws.neon.tech user=neondb_owner password=npg_c4KGnQbY2JUN dbname=neondb port=5432 sslmode=require "
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	db.AutoMigrate(&types.CMSWholeSysRole{}, &types.CMSUser{}, &types.CMSCusPurchase{}, &types.MFAToken{})
	if err := types.SeedAll(db); err != nil {
		log.Fatal("Failed to seed database:", err)
	}
	var count int64
	db.Table("cms_user cu").
		Joins("JOIN cms_whole_sys_role cr ON cu.cms_user_role = cr.role_name").
		Where("cr.role_name = ?", "CMS_CUSTOMER").
		Count(&count)
	fmt.Printf("Total customers: %d\n", count)
	howManyUserUsedTheSystem, err := query.SystemUsersReportSimplified(db)
	if err != nil {
		log.Fatal(err.Error())
	}
	JSONPretty(howManyUserUsedTheSystem)
	usersBySystem, count, err := query.SystemUserReportBySystemNameWithCount(db, string(types.LMS), 10, 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(count, usersBySystem)
	JSONPretty(usersBySystem)

	ordered, err := query.SystemUserReportBySystemNameOrdered(db, string(types.LMS), 10, 0, "")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	JSONPretty(ordered)

}

func JSONPretty(data interface{}) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(string(jsonData))
}
