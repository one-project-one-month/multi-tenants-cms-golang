package tenants

import "gorm.io/gorm"

type TenantList struct {
	gorm.Model
	NameSpace string `gorm:"column:name_space;unique;not null;index"`
}
