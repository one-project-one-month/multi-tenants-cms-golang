package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type AuthHandle interface {
	Login(c *fiber.Ctx) error
	Register(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
	Refresh(c *fiber.Ctx) error
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

	authResponse, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return utils.UnauthorizedResponse(c, err.Error())
	}

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

	authResponse, err := h.service.Register(&req)
	if err != nil {
		if err.Error() == "email already exists" {
			return utils.ConflictResponse(c, err.Error(), nil)
		}
		return utils.InternalServerErrorResponse(c, err.Error(), nil)
	}

	return utils.CreatedResponse(c, "User registered successfully", authResponse)
}

func (h *Handler) Logout(c *fiber.Ctx) error {
	return utils.SuccessResponse(c, "Logout successful", nil)
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
		return utils.UnauthorizedResponse(c, err.Error())
	}

	return utils.SuccessResponse(c, "Token refreshed successfully", tokenResponse)
}
