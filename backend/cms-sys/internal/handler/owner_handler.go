package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/service"
	"github.com/multi-tenants-cms-golang/cms-sys/internal/types"
	"github.com/multi-tenants-cms-golang/cms-sys/pkg/utils"
)

type OwnerHandle interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	GetOwnerByID(c *fiber.Ctx) error
}
type OwnerHandler struct {
	service   service.OwnerService
	validator *validator.Validate
}

var _ OwnerHandle = (*OwnerHandler)(nil)

func NewOwnerHandler(service service.OwnerService) OwnerHandle {
	return &OwnerHandler{service: service, validator: validator.New()}
}

func (o OwnerHandler) Create(c *fiber.Ctx) error {
	var req types.OwnerCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}

	if err := o.validator.Struct(req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed!", err.Error())
	}

	ownerResponse, err := o.service.Create(req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to create owner", err.Error())
	}

	return utils.CreatedResponse(c, "Owner created", ownerResponse)
}

func (o OwnerHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req types.OwnerUpdateRequest

	if err := c.BodyParser(&req); err != nil {
		return utils.BadRequestResponse(c, "Invalid request body", err.Error())
	}
	if err := o.validator.Struct(req); err != nil {
		return utils.BadRequestResponse(c, "Validation failed!", err.Error())
	}

	ownerResponse, err := o.service.Update(id, req)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to update owner", err.Error())
	}

	return utils.SuccessResponse(c, "Owner updated", ownerResponse)
}

func (o OwnerHandler) GetOwnerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	ownerResponse, err := o.service.GetOwnerByID(id)
	if err != nil {
		return utils.InternalServerErrorResponse(c, "Failed to fetch owner", err.Error())
	}

	return utils.SuccessResponse(c, "Fetched owner", ownerResponse)
}
