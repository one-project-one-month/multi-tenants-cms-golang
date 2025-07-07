package types

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

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type OwnerCreateRequest struct {
	Name      string `json:"name" validate:"required,min=2"`
	Email     string `json:"email" validate:"required,email"`
	NameSpace string `json:"namespace" validate:"required,min=3"` // TODO : Gotta remove namespace later
	Password  string `json:"password" validate:"required,min=6"`
}
type OwnerUpdateRequest struct {
	Name      string `json:"name" validate:"required,min=2"`
	NameSpace string `json:"namespace" validate:"required,min=3"` // TODO : Gotta remove namespace later
}
