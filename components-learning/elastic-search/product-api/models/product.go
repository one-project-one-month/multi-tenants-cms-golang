package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ProductID   uuid.UUID `json:"product_id" database:"product_id"`
	ProductName string    `json:"product_name" database:"product_name"`
	Description string    `json:"description" database:"description"`
	CreatedAt   time.Time `json:"created_at" database:"created_at"`
}

// ProductESDoc ElasticsearchDocument for search indexing
type ProductESDoc struct {
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (p *Product) ToESDoc() *ProductESDoc {
	return &ProductESDoc{
		ProductID:   p.ProductID.String(),
		ProductName: p.ProductName,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}
