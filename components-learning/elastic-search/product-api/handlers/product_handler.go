package handlers

import (
	"github.com/SwanHtetAungPhyo/product-api/dto"
	"github.com/SwanHtetAungPhyo/product-api/services"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type ProductHandler struct {
	service *services.ProductService
}

func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req dto.CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.ProductName == "" || req.Description == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Product name and description are required",
		})
	}

	product, err := h.service.CreateProduct(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(201).JSON(product)
}

func (h *ProductHandler) GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	product, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var req dto.UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	product, err := h.service.UpdateProduct(id, &req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(product)
}

func (h *ProductHandler) DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := h.service.DeleteProduct(id); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(204).Send(nil)
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "0"))
	size, _ := strconv.Atoi(c.Query("size", "10"))

	products, total, err := h.service.GetAllProducts(page, size)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"total":    total,
		"page":     page,
		"size":     size,
	})
}

func (h *ProductHandler) SearchProducts(c *fiber.Ctx) error {
	var req dto.SearchRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	result, err := h.service.SearchProducts(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(result)
}
