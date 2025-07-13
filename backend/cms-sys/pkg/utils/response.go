package utils

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success    bool        `json:"success"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func SuccessResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusOK).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(http.StatusCreated).JSON(Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func BadRequestResponse(c *fiber.Ctx, message string, err interface{}) error {
	return c.Status(http.StatusBadRequest).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func UnauthorizedResponse(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusUnauthorized).JSON(Response{
		Success: false,
		Message: message,
	})
}

func ForbiddenResponse(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusForbidden).JSON(Response{
		Success: false,
		Message: message,
	})
}

func NotFoundResponse(c *fiber.Ctx, message string) error {
	return c.Status(http.StatusNotFound).JSON(Response{
		Success: false,
		Message: message,
	})
}

func ConflictResponse(c *fiber.Ctx, message string, err interface{}) error {
	return c.Status(http.StatusConflict).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func InternalServerErrorResponse(c *fiber.Ctx, message string, err interface{}) error {
	return c.Status(http.StatusInternalServerError).JSON(Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}

func PaginatedSuccessResponse(c *fiber.Ctx, message string, data interface{}, pagination Pagination) error {
	return c.Status(http.StatusOK).JSON(PaginatedResponse{
		Success:    true,
		Message:    message,
		Data:       data,
		Pagination: pagination,
	})
}
