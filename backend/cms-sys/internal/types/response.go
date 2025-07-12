package types

import (
	"github.com/google/uuid"
	"time"
)

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token,omitempty"`
	RefreshToken string       `json:"refresh_token,omitempty"`
	ExpiresAt    time.Time    `json:"expires_at,omitempty"`
}
type MFASetupResponse struct {
	Secret      string `json:"secret"`
	QRCodeURL   string `json:"qr_code_url"`
	QRCodeImage string `json:"qr_code_image"` // Base64 encoded PNG
	TokenID     uint   `json:"token_id"`
	ManualEntry string `json:"manual_entry"`
}
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TokenResponse struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type OwnerResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	NameSpace string    `json:"name_space"`
	Verified  bool      `json:"verified"`
}

type PageRequestResponse struct {
	ID          uuid.UUID `json:"id"`
	OwnerID     uuid.UUID `json:"ownerId"`
	RequestType string    `json:"requestType"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status		RequestStatus	`json:"status"`
	PageUrl     *string   `json:"pageUrl"`
	LogoUrl     *string   `json:"logoUrl"`
	CreatedAt   time.Time `json:"createdAt"`
}
