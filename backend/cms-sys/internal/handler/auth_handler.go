package handler

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"strings"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type AuthHandle interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
	GetMe(c *fiber.Ctx) error
	UpdateUserProfile(c *fiber.Ctx) error
	LoginWithMFA(c *fiber.Ctx) error
	VerifyMFASetup(c *fiber.Ctx) error
	SetupMFA(c *fiber.Ctx) error
	VerifyEmail(c *fiber.Ctx) error
}

type Handler struct {
	service   service.AuthService
	validator *validator.Validate
}

var _ AuthHandle = (*Handler)(nil)

func NewHandler(service service.AuthService) *Handler {
	return &Handler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req types.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	user, err := h.service.VerifyCredentials(req.Email, req.Password)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Invalid email or password")
	}

	if !user.Verified {
		return utils.UnauthorizedResponse(c, "Please verify your email before logging in")
	}

	//mfaEnabled, err := h.service.IsMFAEnabled(user.CMSUserID)
	//if err != nil {
	//	return utils.InternalServerErrorResponse(c, "Failed to check MFA status", err.Error())
	//}
	//
	//if mfaEnabled {
	//	return utils.SuccessResponse(c, "MFA verification required", map[string]interface{}{
	//		"mfa_required": true,
	//		"user_id":      user.CMSUserID,
	//		"message":      "Please provide MFA code to complete login",
	//	})
	//}

	authResponse, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.UnauthorizedResponse(c, "Login failed")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    authResponse.AccessToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResponse.RefreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})
	return utils.SuccessResponse(c, "Login successful", authResponse)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req types.RegisterRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	ipAddress := c.IP()
	authResponse, err := h.service.Register(&req, ipAddress)
	if err != nil {
		switch {
		case errors.Is(err, errors.New("email already exists")):
			return utils.ConflictResponse(c, "An account with this email already exists", nil)
		case strings.Contains(err.Error(), "failed to send verification email"):
			return utils.SuccessResponse(c, "Registration successful but verification email failed", authResponse)
		default:
			return utils.InternalServerErrorResponse(c, "Registration failed", err.Error())
		}
	}

	message := "Registration successful! Please check your email for verification code."
	if authResponse.User.Verified {
		message = "Registration successful! Your account is ready to use."
	}

	return utils.CreatedResponse(c, message, authResponse)
}

func (h *Handler) VerifyEmail(c *fiber.Ctx) error {
	var req types.EmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	err := h.service.VerifyEmail(req.Email, req.Code)
	if err != nil {
		switch err.Error() {
		case "verification code expired":
			return utils.BadRequestResponse(c, "Verification code has expired", nil)
		case "invalid verification code":
			return utils.BadRequestResponse(c, "Invalid verification code", nil)
		default:
			return utils.InternalServerErrorResponse(c, "Email verification failed", err.Error())
		}
	}

	return utils.SuccessResponse(c, "Email verified successfully! You can now log in to your account.", nil)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var accessToken string
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		accessToken = strings.TrimPrefix(authHeader, "Bearer ")
	}

	var refreshToken string
	var req struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.BodyParser(&req); err == nil {
		refreshToken = req.RefreshToken
	}

	if accessToken == "" && refreshToken == "" {
		return utils.BadRequestResponse(c, "No tokens provided for logout", nil)
	}

	if err := h.service.Logout(accessToken, refreshToken); err != nil {
		return utils.InternalServerErrorResponse(c, "Logout failed", err.Error())
	}

	return utils.SuccessResponse(c, "Logged out successfully", nil)
}

func (h *Handler) Refresh(c *fiber.Ctx) error {
	var req types.RefreshTokenRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	tokenResponse, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		switch err.Error() {
		case "invalid refresh token":
			return utils.UnauthorizedResponse(c, "Invalid or expired refresh token")
		case "user not found":
			return utils.UnauthorizedResponse(c, "User account not found")
		default:
			return utils.InternalServerErrorResponse(c, "Token refresh failed", err.Error())
		}
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", tokenResponse)
}

// Note: Update the authentication me method to get userId from JWT token
// By Swan Htet Aung Phyo

func (h *Handler) GetMe(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uuid.UUID)

	profileResponse, err := h.service.GetUserProfile(userID)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return utils.NotFoundResponse(c, "User profile not found")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to get user profile", err.Error())
		}
	}

	return utils.SuccessResponse(c, "User profile retrieved successfully", profileResponse)
}

func (h *Handler) UpdateUserProfile(c *fiber.Ctx) error {
	stringId := c.Params("id")
	id, err := uuid.Parse(stringId)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID format", err.Error())
	}
	userID := c.Locals("userID").(uuid.UUID)
	if userID != id {
		return utils.ForbiddenResponse(c, "You can only update your own profile")
	}

	var req types.UserUpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	updatedProfileResponse, err := h.service.UpdateUserProfile(id, req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			return utils.NotFoundResponse(c, "User not found")
		default:
			return utils.InternalServerErrorResponse(c, "Failed to update user profile", err.Error())
		}
	}

	return utils.SuccessResponse(c, "User profile updated successfully", updatedProfileResponse)
}

func (h *Handler) SetupMFA(c *fiber.Ctx) error {
	userID := c.Params("userid")

	uuidUserID, err := uuid.Parse(userID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID format", err.Error())
	}
	//mfaEnabled, err := h.service.IsMFAEnabled(userID)
	//if err != nil {
	//	return utils.InternalServerErrorResponse(c, "Failed to check MFA status", err.Error())
	//}

	//if mfaEnabled {
	//	return utils.ConflictResponse(c, "MFA is already enabled for this account", nil)
	//}

	mfaSetup, err := h.service.GenerateMFATokenSecret(uuidUserID)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to setup MFA", err.Error())
	}

	return utils.SuccessResponse(c, "MFA setup initiated successfully", map[string]interface{}{
		"setup_data": mfaSetup,
		"next_step":  "verify_mfa_setup",
		"message":    "Scan the QR code with your authenticator app and enter the verification code",
	})
}

func (h *Handler) VerifyMFASetup(c *fiber.Ctx) error {
	var req types.MFAVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	userID := c.Params("userid")
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID format", err.Error())
	}

	if err := h.service.VerifyMFASetup(userUUID, req.TokenID, req.Code); err != nil {
		switch err.Error() {
		case "mfa token expired":
			return utils.BadRequestResponse(c, "MFA setup token has expired. Please restart the setup process.", nil)
		case "mfa token is invalid":
			return utils.BadRequestResponse(c, "Invalid MFA verification code", nil)
		case "failed to get mfa token":
			return utils.NotFoundResponse(c, "MFA setup token not found")
		default:
			return utils.InternalServerErrorResponse(c, "MFA verification failed", err.Error())
		}
	}

	return utils.SuccessResponse(c, "MFA setup completed successfully! Your account is now protected with two-factor authentication.", nil)
}

func (h *Handler) LoginWithMFA(c *fiber.Ctx) error {
	var req types.MFALoginRequest
	userId := c.Params("userid")
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return utils.BadRequestResponse(c, "Invalid user ID format", err.Error())
	}
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := h.validator.Struct(&req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed", err.Error())
	}

	//user, err := h.service.VerifyCredentials(req.Email, req.Password)
	//if err != nil {
	//	return utils.UnauthorizedResponse(c, "Invalid email or password")
	//}
	//
	//if !user.Verified {
	//	return utils.UnauthorizedResponse(c, "Please verify your email before logging in")
	//}

	//mfaEnabled, err := h.service.IsMFAEnabled(user.CMSUserID)
	//if err != nil {
	//	return utils.InternalServerErrorResponse(c, "Failed to check MFA status", err.Error())
	//}

	//if !mfaEnabled {
	//	return utils.BadRequestResponse(c, "MFA is not enabled for this account", nil)
	//}

	if req.MFACode == "" {
		return utils.BadRequestResponse(c, "MFA verification code is required", nil)
	}

	authResponse, err := h.service.VerifyMFALogin(userUUID, req.MFACode)
	if err != nil {
		switch err.Error() {
		case "mfa token is invalid":
			return utils.UnauthorizedResponse(c, "Invalid MFA verification code")
		case "failed to get mfa token":
			return utils.InternalServerErrorResponse(c, "MFA verification failed", nil)
		default:
			return utils.UnauthorizedResponse(c, "MFA verification failed")
		}
	}

	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    authResponse.RefreshToken,
		Path:     "/",
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
	})

	authResponse.RefreshToken = ""

	return utils.SuccessResponse(c, "Login successful with MFA verification", authResponse)
}
