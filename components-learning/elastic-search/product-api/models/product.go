package models

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ProductID   uuid.UUID `json:"product_id" db:"product_id"`
	ProductName string    `json:"product_name" db:"product_name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
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
