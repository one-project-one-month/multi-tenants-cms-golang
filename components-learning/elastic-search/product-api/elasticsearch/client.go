package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/SwanHtetAungPhyo/product-api/dto"
	"github.com/SwanHtetAungPhyo/product-api/models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/google/uuid"
	"io"
	"log"
	"strings"
)

type Client struct {
	es *elasticsearch.Client
}

func NewClient(url string) (*Client, error) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			url,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	client := &Client{
		es: es,
	}

	if err := client.createIndex(); err != nil {
		log.Fatal(err.Error())
	}
	return client, nil
}

func (c *Client) createIndex() error {
	indexMapping := `{
        "mappings": {
            "properties": {
                "product_id": {"type": "keyword"},
                "product_name": {
                    "type": "text",
                    "analyzer": "standard",
                    "fields": {
                        "keyword": {"type": "keyword"}
                    }
                },
                "description": {
                    "type": "text",
                    "analyzer": "standard"
                },
                "created_at": {"type": "date"}
            }
        }
    }`

	req := esapi.IndexRequest{
		Index: "products",
		Body:  strings.NewReader(indexMapping),
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(res.Body)

	return nil
}

func (c *Client) IndexProduct(product *models.Product) error {
	doc := product.ToESDoc()
	docJson, err := json.Marshal(doc)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	req := esapi.IndexRequest{
		Index:      "products",
		DocumentID: product.ProductID.String(),
		Body:       strings.NewReader(string(docJson)),
	}

	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(res.Body)
	return nil
}

func (c *Client) DeleteProduct(productID string) error {
	req := esapi.DeleteRequest{
		Index:      "products",
		DocumentID: productID,
	}
	res, err := req.Do(context.Background(), c.es)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(res.Body)
	return nil
}

func (c *Client) SearchProduct(searchReq *dto.SearchRequest) (*dto.SearchResponse, error) {
	query := c.buildSearchQuery(searchReq)
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}
	res, err := c.es.Search(
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex("products"),
		c.es.Search.WithBody(&buf),
		c.es.Search.WithFrom(searchReq.Page*searchReq.Size),
		c.es.Search.WithSize(searchReq.Size),
		c.es.Search.WithSort("created_at:desc"),
	)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(res.Body)
	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		}
	}
	return c.parseSearchResponse(res, searchReq)
}
func (c *Client) buildSearchQuery(searchQuery *dto.SearchRequest) map[string]interface{} {
	if searchQuery.Query == "" {
		return map[string]interface{}{
			"query": map[string]interface{}{
				"match_all": map[string]interface{}{},
			},
		}
	}
	return map[string]interface{}{
		"query": map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  searchQuery.Query,
				"fields": []string{"product_name^2", "description"},
				"type":   "best_fields",
			},
		},
		"highlight": map[string]interface{}{
			"fields": map[string]interface{}{
				"product_name": map[string]interface{}{},
				"description":  map[string]interface{}{},
			},
		},
	}

}

func (c *Client) parseSearchResponse(res *esapi.Response, req *dto.SearchRequest) (*dto.SearchResponse, error) {
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	hits := result["hits"].(map[string]interface{})
	total := int64(hits["total"].(map[string]interface{})["value"].(float64))
	products := []dto.ProductResponse{}
	for _, hit := range hits["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"].(map[string]interface{})

		product := dto.ProductResponse{
			ProductName: source["product_name"].(string),
			Description: source["description"].(string),
			CreatedAt:   source["created_at"].(string),
		}

		if id, ok := source["product_id"].(string); ok {
			if parsed, err := uuid.Parse(id); err == nil {
				product.ProductID = parsed
			}
		}

		products = append(products, product)
	}
	return &dto.SearchResponse{
		Products: products,
		Total:    total,
		Page:     req.Page,
		Size:     req.Size,
	}, nil
}
