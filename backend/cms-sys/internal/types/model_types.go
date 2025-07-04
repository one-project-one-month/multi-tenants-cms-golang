package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RoleType string
type SystemType string

const (
	RootAdmin   RoleType = "ROOT_ADMIN"
	CMSCustomer RoleType = "CMS_CUSTOMER"
)

const (
	LMS SystemType = "LMS"
	EMS SystemType = "EMS"
)

type CMSWholeSysRole struct {
	RoleID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"role_id"`
	RoleName string    `gorm:"type:varchar(15);not null;unique" json:"role_name"`
	Users    []CMSUser `gorm:"foreignKey:CMSUserRole;references:RoleName" json:"users,omitempty"`
}

func (CMSWholeSysRole) TableName() string {
	return "cms_whole_sys_role"
}

type CMSUser struct {
	CMSUserID    uuid.UUID        `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"cms_user_id"`
	CMSUserName  string           `gorm:"size:100;not null" json:"cms_user_name"`
	CMSUserEmail string           `gorm:"size:150;not null;unique" json:"cms_user_email"`
	CMSNameSpace *string          `gorm:"size:100" json:"cms_name_space,omitempty"`
	Password     string           `gorm:"size:90;not null" json:"-"`
	CMSUserRole  string           `gorm:"size:15;not null" json:"cms_user_role"`
	Verified     bool             `gorm:"default:false" json:"verified"`
	CreatedAt    time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Role         CMSWholeSysRole  `gorm:"foreignKey:CMSUserRole;references:RoleName" json:"role,omitempty"`
	Purchases    []CMSCusPurchase `gorm:"foreignKey:CMSCusID" json:"purchases,omitempty"`
	MFATokens    []MFAToken       `gorm:"foreignKey:UserID" json:"mfa_tokens,omitempty"`
}

func (CMSUser) TableName() string {
	return "cms_user"
}

func (u *CMSUser) BeforeCreate(tx *gorm.DB) error {
	if u.CMSUserID == uuid.Nil {
		u.CMSUserID = uuid.New()
	}
	return nil
}

type CMSCusPurchase struct {
	RelationID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"relation_id"`
	CMSCusID     uuid.UUID `gorm:"type:uuid;not null" json:"cms_cus_id"`
	SystemName   string    `gorm:"type:varchar(100);not null" json:"system_name"`
	PurchaseDate time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"purchase_date"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	Customer     CMSUser   `gorm:"foreignKey:CMSCusID;references:CMSUserID" json:"customer,omitempty"`
}

func (CMSCusPurchase) TableName() string {
	return "cms_cus_purchase"
}

func (p *CMSCusPurchase) BeforeCreate(tx *gorm.DB) error {
	if p.RelationID == uuid.Nil {
		p.RelationID = uuid.New()
	}
	return nil
}

type MFAToken struct {
	TokenID   uint       `gorm:"primaryKey;autoIncrement" json:"token_id"`
	MFAToken  string     `gorm:"not null" json:"mfa_token"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	CreatedAt time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	User      CMSUser    `gorm:"foreignKey:UserID;references:CMSUserID" json:"user,omitempty"`
}

func (MFAToken) TableName() string {
	return "mfa_token"
}

type PageRequestStatus string
type PageType string

const (
	PageStatusPending  PageRequestStatus = "PENDING"
	PageStatusApproved PageRequestStatus = "APPROVED"
	PageStatusRejected PageRequestStatus = "REJECTED"
)

type UserPageRequest struct {
	UserPageRequestID uint              `gorm:"primaryKey;autoIncrement" json:"user_page_request_id"`
	UserID            uuid.UUID         `gorm:"type:uuid;not null;unique" json:"user_id"`
	PageType          PageType          `gorm:"type:system_check_domain;not null" json:"pagetype"`
	Status            PageRequestStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`

	User CMSUser `gorm:"foreignKey:UserID;references:CMSUserID" json:"user,omitempty"`
}

func (UserPageRequest) TableName() string {
	return "user_page_request"
}
