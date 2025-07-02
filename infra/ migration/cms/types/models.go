package types

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

// RoleType Enums as string constants
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

// CMSWholeSysRole represents the cms_whole_sys_role table
type CMSWholeSysRole struct {
	RoleID   uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"role_id"`
	RoleName string    `gorm:"type:varchar(15);not null;unique" json:"role_name"`

	// One-to-many relationship: One role can have many users
	Users []CMSUser `gorm:"foreignKey:CMSUserRole;references:RoleName" json:"users,omitempty"`
}

func (CMSWholeSysRole) TableName() string {
	return "cms_whole_sys_role"
}

// CMSUser represents the cms_user table
type CMSUser struct {
	CMSUserID    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"cms_user_id"`
	CMSUserName  string    `gorm:"size:100;not null" json:"cms_user_name"`
	CMSUserEmail string    `gorm:"size:150;not null;unique" json:"cms_user_email"`
	CMSNameSpace *string   `gorm:"size:100" json:"cms_name_space,omitempty"`
	Password     string    `gorm:"size:90;not null" json:"-"`
	CMSUserRole  string    `gorm:"size:15;not null" json:"cms_user_role"`
	Verified     bool      `gorm:"default:false" json:"verified"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Associations
	Role      CMSWholeSysRole  `gorm:"foreignKey:CMSUserRole;references:RoleName" json:"role,omitempty"`
	Purchases []CMSCusPurchase `gorm:"foreignKey:CMSCusID" json:"purchases,omitempty"`
	MFATokens []MFAToken       `gorm:"foreignKey:UserID" json:"mfa_tokens,omitempty"`
}

func (CMSUser) TableName() string {
	return "cms_user"
}

// BeforeCreate hook to set default values
func (u *CMSUser) BeforeCreate(tx *gorm.DB) error {
	if u.CMSUserID == uuid.Nil {
		u.CMSUserID = uuid.New()
	}
	return nil
}

// CMSCusPurchase represents the cms_cus_purchase table
type CMSCusPurchase struct {
	RelationID   uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"relation_id"`
	CMSCusID     uuid.UUID  `gorm:"type:uuid;not null" json:"cms_cus_id"`
	SystemName   SystemType `gorm:"type:system_type;not null" json:"system_name"`
	PurchaseDate time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"purchase_date"`
	CreatedAt    time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Associations
	Customer CMSUser `gorm:"foreignKey:CMSCusID;references:CMSUserID" json:"customer,omitempty"`
}

func (CMSCusPurchase) TableName() string {
	return "cms_cus_purchase"
}

// BeforeCreate hook to set default values
func (p *CMSCusPurchase) BeforeCreate(tx *gorm.DB) error {
	if p.RelationID == uuid.Nil {
		p.RelationID = uuid.New()
	}
	return nil
}

// MFAToken represents the mfa_token table
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

// IsAdmin Custom methods for CMSUser
func (u *CMSUser) IsAdmin() bool {
	return u.CMSUserRole == string(RootAdmin)
}

func (u *CMSUser) IsCustomer() bool {
	return u.CMSUserRole == string(CMSCustomer)
}

func (u *CMSUser) HasPurchased(systemType SystemType) bool {
	for _, purchase := range u.Purchases {
		if purchase.SystemName == systemType {
			return true
		}
	}
	return false
}

// SeedRoles Seed data functions
func SeedRoles(db *gorm.DB) error {
	roles := []CMSWholeSysRole{
		{RoleName: string(RootAdmin)},
		{RoleName: string(CMSCustomer)},
	}

	for _, role := range roles {
		if err := db.Where("role_name = ?", role.RoleName).FirstOrCreate(&role).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedUsers(db *gorm.DB) error {
	users := []CMSUser{
		{
			CMSUserName:  "System Administrator",
			CMSUserEmail: "admin@system.com",
			Password:     "$2a$10$hashedpassword1",
			CMSUserRole:  string(RootAdmin),
			Verified:     true,
		},
		{
			CMSUserName:  "John Doe",
			CMSUserEmail: "john.doe@customer.com",
			Password:     "$2a$10$hashedpassword2",
			CMSUserRole:  string(CMSCustomer),
			Verified:     true,
		},
		{
			CMSUserName:  "Jane Smith",
			CMSUserEmail: "jane.smith@customer.com",
			Password:     "$2a$10$hashedpassword3",
			CMSUserRole:  string(CMSCustomer),
			Verified:     false,
		},
	}

	for _, user := range users {
		if err := db.Where("cms_user_email = ?", user.CMSUserEmail).FirstOrCreate(&user).Error; err != nil {
			return err
		}
	}
	return nil
}

func SeedPurchases(db *gorm.DB) error {
	// First, get some users to create purchases for
	var customers []CMSUser
	if err := db.Where("cms_user_role = ?", string(CMSCustomer)).Find(&customers).Error; err != nil {
		return err
	}

	if len(customers) == 0 {
		return nil
	}

	purchases := []CMSCusPurchase{
		{
			CMSCusID:     customers[0].CMSUserID,
			SystemName:   LMS,
			PurchaseDate: time.Now().AddDate(0, -1, 0),
		},
		{
			CMSCusID:     customers[0].CMSUserID,
			SystemName:   EMS,
			PurchaseDate: time.Now().AddDate(0, 0, -15),
		},
	}

	if len(customers) > 1 {
		purchases = append(purchases, CMSCusPurchase{
			CMSCusID:     customers[1].CMSUserID,
			SystemName:   LMS,
			PurchaseDate: time.Now().AddDate(0, 0, -7),
		})
	}

	for _, purchase := range purchases {
		if err := db.Create(&purchase).Error; err != nil {
			return err
		}
	}
	return nil
}

// SeedAll runs all seed functions
func SeedAll(db *gorm.DB) error {
	if err := SeedRoles(db); err != nil {
		return err
	}
	if err := SeedUsers(db); err != nil {
		return err
	}
	if err := SeedPurchases(db); err != nil {
		return err
	}
	return nil
}
