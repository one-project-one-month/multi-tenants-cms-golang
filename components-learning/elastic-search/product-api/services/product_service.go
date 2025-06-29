package services

import (
	"database/sql"
	"fmt"
	"github.com/SwanHtetAungPhyo/product-api/dto"
	"github.com/SwanHtetAungPhyo/product-api/elasticsearch"
	"github.com/SwanHtetAungPhyo/product-api/models"
	"github.com/google/uuid"

	"strings"
)

type ProductService struct {
	db *sql.DB
	es *elasticsearch.Client
}

func NewProductService(db *sql.DB, es *elasticsearch.Client) *ProductService {
	return &ProductService{db: db, es: es}
}

func (s *ProductService) CreateProduct(req *dto.CreateProductRequest) (*models.Product, error) {
	product := &models.Product{
		ProductID:   uuid.New(),
		ProductName: req.ProductName,
		Description: req.Description,
	}

	query := `
        INSERT INTO products (product_id, product_name, description) 
        VALUES ($1, $2, $3) 
        RETURNING created_at`

	err := s.db.QueryRow(query, product.ProductID, product.ProductName, product.Description).
		Scan(&product.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Index in Elasticsearch
	if err := s.es.IndexProduct(product); err != nil {
		fmt.Printf("Warning: Failed to index product in Elasticsearch: %v\n", err)
	}

	return product, nil
}

func (s *ProductService) GetProductByID(id string) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	var product models.Product
	query := `SELECT product_id, product_name, description, created_at FROM products WHERE product_id = $1`

	err = s.db.QueryRow(query, productID).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Description,
		&product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

func (s *ProductService) UpdateProduct(id string, req *dto.UpdateProductRequest) (*models.Product, error) {
	productID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid product ID: %w", err)
	}

	// Build dynamic query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.ProductName != nil {
		setParts = append(setParts, fmt.Sprintf("product_name = $%d", argIndex))
		args = append(args, *req.ProductName)
		argIndex++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIndex))
		args = append(args, *req.Description)
		argIndex++
	}

	if len(setParts) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	args = append(args, productID)
	query := fmt.Sprintf(`
        UPDATE products 
        SET %s 
        WHERE product_id = $%d 
        RETURNING product_id, product_name, description, created_at`,
		strings.Join(setParts, ", "), argIndex)

	var product models.Product
	err = s.db.QueryRow(query, args...).Scan(
		&product.ProductID,
		&product.ProductName,
		&product.Description,
		&product.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	if err := s.es.IndexProduct(&product); err != nil {
		fmt.Printf("Warning: Failed to update product in Elasticsearch: %v\n", err)
	}

	return &product, nil
}

func (s *ProductService) DeleteProduct(id string) error {
	productID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid product ID: %w", err)
	}

	query := `DELETE FROM products WHERE product_id = $1`
	result, err := s.db.Exec(query, productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product not found")
	}

	// Delete from Elasticsearch
	if err := s.es.DeleteProduct(id); err != nil {
		fmt.Printf("Warning: Failed to delete product from Elasticsearch: %v\n", err)
	}

	return nil
}

func (s *ProductService) GetAllProducts(page, size int) ([]models.Product, int64, error) {
	if page < 0 {
		page = 0
	}
	if size <= 0 || size > 100 {
		size = 10
	}

	offset := page * size

	// Get total count
	var total int64
	countQuery := `SELECT COUNT(*) FROM products`
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	// Get products
	query := `
        SELECT product_id, product_name, description, created_at 
        FROM products 
        ORDER BY created_at DESC 
        LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(query, size, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get products: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("Failed to close rows: %v\n", err)
		}
	}(rows)

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.Description,
			&product.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	return products, total, nil
}

func (s *ProductService) SearchProducts(req *dto.SearchRequest) (*dto.SearchResponse, error) {
	// Set defaults
	if req.Page < 0 {
		req.Page = 0
	}
	if req.Size <= 0 || req.Size > 100 {
		req.Size = 10
	}

	return s.es.SearchProduct(req)
}
