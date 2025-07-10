package types

import (
	"github.com/google/uuid"
	"time"
)

// CMSUserResponse DTO structs for API responses
type CMSUserResponse struct {
	CMSUserID    uuid.UUID `json:"cms_user_id"`
	CMSUserName  string    `json:"cms_user_name"`
	CMSUserEmail string    `json:"cms_user_email"`
	CMSNameSpace *string   `json:"cms_name_space,omitempty"`
	RoleName     RoleType  `json:"role_name"`
	Verified     bool      `json:"verified"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CMSUserCreateRequest struct {
	CMSUserName   string    `json:"cms_user_name" binding:"required,min=2,max=100"`
	CMSUserEmail  string    `json:"cms_user_email" binding:"required,email,max=150"`
	CMSNameSpace  *string   `json:"cms_name_space,omitempty" binding:"omitempty,max=100"`
	Password      string    `json:"password" binding:"required,min=8,max=90"`
	CMSUserRoleID uuid.UUID `json:"cms_user_role_id" binding:"required"`
}

type CMSUserUpdateRequest struct {
	CMSUserName  *string `json:"cms_user_name,omitempty" binding:"omitempty,min=2,max=100"`
	CMSNameSpace *string `json:"cms_name_space,omitempty" binding:"omitempty,max=100"`
	Verified     *bool   `json:"verified,omitempty"`
}

type PurchaseRequest struct {
	CMSCusID   uuid.UUID  `json:"cms_cus_id" binding:"required"`
	SystemName SystemType `json:"system_name" binding:"required"`
}

type PurchaseResponse struct {
	RelationID   uuid.UUID  `json:"relation_id"`
	CustomerName string     `json:"customer_name"`
	SystemName   SystemType `json:"system_name"`
	PurchaseDate time.Time  `json:"purchase_date"`
}
