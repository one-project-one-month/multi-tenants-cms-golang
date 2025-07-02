package tenant_models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Name  string  `gorm:"index;not null"`
	Price float64 `gorm:"default:0;check:price >= 0"`
}

func (item Item) TableName() string {
	return "items"
}
