package query

import (
	"fmt"
	"github.com/SwanHtetAungPhyo/migration/types"
	"gorm.io/gorm"
	"log"
)

type (
	LoginReq struct {
		email    string
		password string
	}
)

func SystemUsersReportSimplified(db *gorm.DB) ([]*types.CustomerPurchaseReport, error) {
	var results []*types.CustomerPurchaseReport
	err := db.Table("cms_user cu").
		Select("cu.cms_user_name, cu.cms_user_email, cp.system_name, cp.purchase_date").
		Joins("JOIN cms_cus_purchase cp ON cu.cms_user_id = cp.cms_cus_id").
		Where("cu.cms_user_role = ?", string(types.CMSCustomer)).
		Find(&results).Error
	return results, err
}

func SystemUserReportBySystemNameWithCount(db *gorm.DB, systemName string, limit, offset int) ([]*types.CustomerPurchaseReport, int64, error) {
	var results []*types.CustomerPurchaseReport
	var total int64

	err := db.Table("cms_user cu").
		Joins("JOIN cms_cus_purchase cp ON cu.cms_user_id = cp.cms_cus_id").
		Where("cu.cms_user_role = ?", string(types.CMSCustomer)).
		Where("cp.system_name = ?", systemName).
		Count(&total).Error

	if err != nil {
		return nil, 0, err
	}

	err = db.Table("cms_user cu").
		Select("cu.cms_user_name, cu.cms_user_email, cp.system_name, cp.purchase_date").
		Joins("JOIN cms_cus_purchase cp ON cu.cms_user_id = cp.cms_cus_id").
		Where("cu.cms_user_role = ?", string(types.CMSCustomer)).
		Where("cp.system_name = ?", systemName).
		Limit(limit).
		Offset(offset).
		Find(&results).Error

	if err != nil {
		return nil, 0, err
	}

	return results, total, nil
}
func SystemUserReportBySystemNameOrdered(db *gorm.DB, systemName string, limit, offset int, orderBy string) ([]*types.CustomerPurchaseReport, error) {
	var results []*types.CustomerPurchaseReport

	query := db.Table("cms_user cu").
		Select("cu.cms_user_name, cu.cms_user_email, cp.system_name, cp.purchase_date").
		Joins("JOIN cms_cus_purchase cp ON cu.cms_user_id = cp.cms_cus_id").
		Where("cu.cms_user_role = ?", string(types.CMSCustomer)).
		Where("cp.system_name = ?", systemName).
		Limit(limit).
		Offset(offset)

	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("cp.purchase_date DESC")
	}

	err := query.Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func Login(db *gorm.DB, email string) {
	var loginReq LoginReq
	err := db.Table("cms_user cu").Select(" cu.cms_user_email, cp.password").Where("cu.email = ?", email).Find(&loginReq).Error
	if err != nil {
		log.Fatalf("Failed to query cms_user: %v", err)
		return
	}
	fmt.Println(loginReq)
}
