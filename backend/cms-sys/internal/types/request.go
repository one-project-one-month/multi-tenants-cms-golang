package types

import "mime/multipart"

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role,omitempty"`
}

type MFAVerificationRequest struct {
	TokenID uint   `json:"token_id" validate:"required"`
	Code    string `json:"code" validate:"required,len=6,numeric"`
}

type MFALoginRequest struct {
	//Email    string `json:"email" validate:"required,email"`
	//Password string `json:"password" validate:"required"`
	MFACode string `json:"mfa_code,omitempty" validate:"omitempty,len=6,numeric"`
}
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type OwnerCreateRequest struct {
	Name      string `json:"name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	NameSpace string `json:"namespaces" validate:"required,min=3"` // TODO : Gotta remove namespace later
	Password  string `json:"password" validate:"required,min=6"`
}
type OwnerUpdateRequest struct {
	Name      string `json:"name" validate:"required,min=2"`
	NameSpace string `json:"namespace" validate:"required,min=3"` // TODO : Gotta remove namespace later
}

type GetMeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserUpdateRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type PaginateRequest struct {
	Page  int `query:"page"`
	Limit int `query:"limit"`
}

type CreatePageRequest struct {
	OwnerID     string                `json:"ownerId" validate:"required,uuid4"`
	RequestType string                `json:"requestType"`
	Title       string                `json:"title"`
	Description string                `json:"description"`
	PageUrl     *string               `json:"pageUrl"`
	LogoFile    *multipart.FileHeader `json:"logo"`
}

type OwnerDeleteRequest struct {
	IDs         []string `json:"ids" validate:"required"`
	ForceDelete bool     `json:"forceDelete"`
}
type EmailVerificationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

type ChangeStatusPageRequest struct {
	RequestID string        `json:"requestId" validate:"required,uuid4"`
	Status    RequestStatus `json:"status" validate:"required"`
}
