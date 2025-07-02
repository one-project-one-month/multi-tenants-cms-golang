package types

import "time"

type CustomerPurchaseReport struct {
	CMSUserName  string     `json:"cms_user_name"`
	CMSUserEmail string     `json:"cms_user_email"`
	SystemName   SystemType `json:"system_name"`
	PurchaseDate time.Time  `json:"purchase_date"`
}

const (
	CMSTABLE = "cms_user cu "
)
