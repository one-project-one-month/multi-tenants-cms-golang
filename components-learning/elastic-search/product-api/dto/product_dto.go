package dto

import "github.com/google/uuid"

type CreateProductRequest struct {
	ProductName string `json:"product_name" validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1"`
}

type UpdateProductRequest struct {
	ProductName *string `json:"product_name,omitempty"`
	Description *string `json:"description,omitempty"`
}

type SearchRequest struct {
	Query string `json:"query" form:"query"`
	Page  int    `json:"page" form:"page"`
	Size  int    `json:"size" form:"size"`
}

type SearchResponse struct {
	Products []ProductResponse `json:"products"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	Size     int               `json:"size"`
}

type ProductResponse struct {
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	CreatedAt   string    `json:"created_at"`
}
