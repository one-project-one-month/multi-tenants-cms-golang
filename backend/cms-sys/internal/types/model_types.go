package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RoleType string
type SystemType string

type PageStatus string

type RequestStatus string

const (
	RootAdmin   RoleType = "ROOT_ADMIN"
	CMSCustomer RoleType = "CMS_CUSTOMER"
)

const (
	LMS SystemType = "LMS"
	EMS SystemType = "EMS"
)

//const (
//	PageStatusDraft     PageStatus = "DRAFT"
//	PageStatusPublished PageStatus = "PUBLISHED"
//	PageStatusArchived  PageStatus = "ARCHIVED"
//)
//
const (
	RequestStatusPending  RequestStatus = "PENDING"
	RequestStatusApproved RequestStatus = "APPROVED"
	RequestStatusRejected RequestStatus = "REJECTED"
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
	MFAEnabled   bool            `gorm:"default:false" json:"mfa_enabled"`

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
	PageType          PageType          `gorm:"type:varchar(20);not null" json:"pagetype"`
	Status            PageRequestStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`

	User CMSUser `gorm:"foreignKey:UserID;references:CMSUserID" json:"user,omitempty"`
}

func (UserPageRequest) TableName() string {
	return "user_page_request"
}

type Page struct {
	PageID             uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"page_id"`
	PageRequestID      uuid.UUID  `gorm:"type:uuid;not null" json:"page_request_id"`
	Title              string     `gorm:"type:varchar(255);not null" json:"title"`
	Content            string     `gorm:"type:text" json:"content"`
	ImageURL           *string    `gorm:"type:varchar(255)" json:"image_url,omitempty"`
	Status             PageStatus `gorm:"type:varchar(50);default:'DRAFT'" json:"status"`
	OwnerID            uuid.UUID  `gorm:"type:uuid;not null" json:"owner_id"`
	PublishedByStaffID *uuid.UUID `gorm:"type:uuid" json:"published_by_staff_id,omitempty"`
	CreatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	Owner            CMSUser     `gorm:"foreignKey:OwnerID;references:CMSUserID" json:"owner,omitempty"`
	PublishedByStaff *CMSUser    `gorm:"foreignKey:PublishedByStaffID;references:CMSUserID" json:"published_by_staff,omitempty"`
	PageRequest      PageRequest `gorm:"foreignKey:PageRequestID;references:RequestID" json:"page_request,omitempty"`
}

func (Page) TableName() string {
	return "cms_page"
}

func (p *Page) BeforeCreate(tx *gorm.DB) error {
	if p.PageID == uuid.Nil {
		p.PageID = uuid.New()
	}
	return nil
}

type PageRequest struct {
	RequestID   uuid.UUID     `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"request_id"`
	OwnerID     uuid.UUID     `gorm:"type:uuid;not null" json:"owner_id"`
	RequestType string        `gorm:"type:varchar(100);not null" json:"request_type"`
	Title       string        `gorm:"type:varchar(100);not null" json:"title"`
	Description *string       `gorm:"type:text" json:"description,omitempty"`
	PageUrl     *string       `gorm:"type:varchar(100)" json:"page_url"`
	LogoUrl     *string       `gorm:"type:text" json:"logo_url,omitempty"`
	Status      RequestStatus `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	AdminID     *uuid.UUID    `gorm:"type:uuid" json:"admin_id,omitempty"`
	CreatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Owner       CMSUser       `gorm:"foreignKey:OwnerID;references:CMSUserID" json:"owner,omitempty"`
	Admin       *CMSUser      `gorm:"foreignKey:AdminID;references:CMSUserID" json:"admin,omitempty"`
}

func (PageRequest) TableName() string {
	return "cms_page_request"
}

func (pr *PageRequest) BeforeCreate(tx *gorm.DB) error {
	if pr.RequestID == uuid.Nil {
		pr.RequestID = uuid.New()
	}
	return nil
}
