package handler

import (
	"context"
	"errors"
	"log"
	"net/http"
	"slices"
	"strings"
	"strconv"

	"github.com/gin-gonic/gin"

	"secret-shop/internal/domain"
	"secret-shop/internal/service"
)

const (
	defaultListLimit = 20
	maxListLimit     = 100
)

type ProductService interface {
	ListProducts(ctx context.Context, limit, offset int, fields []string) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id int64, fields []string) (domain.Product, error)
}

type ProductHandler struct {
	service ProductService
}

var allowedProductFields = []string{"id", "name", "description", "slug", "isActive"}

func NewProductHandler(service ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	v1.GET("/products", h.listProducts)
	v1.GET("/products/:id", h.getProductByID)
}

func (h *ProductHandler) listProducts(c *gin.Context) {
	limit, err := parsePositiveIntWithDefault(c.Query("limit"), defaultListLimit)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
		return
	}
	if limit > maxListLimit {
		limit = maxListLimit
	}

	offset, err := parseNonNegativeIntWithDefault(c.Query("offset"), 0)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "offset must be a non-negative integer"})
		return
	}

	fields, err := parseProjectedFields(c.Query("fields"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("debug projection: endpoint=listProducts raw_query=%q raw_fields=%q parsed_fields=%v", c.Request.URL.RawQuery, c.Query("fields"), fields)

	products, err := h.service.ListProducts(c.Request.Context(), limit, offset, fields)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	payload := make([]gin.H, 0, len(products))
	for _, product := range products {
		payload = append(payload, toProductResponse(product, fields))
	}

	c.JSON(http.StatusOK, gin.H{
		"data": payload,
		"meta": gin.H{
			"limit":  limit,
			"offset": offset,
			"count":  len(products),
		},
	})
}

func (h *ProductHandler) getProductByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
		return
	}

	fields, err := parseProjectedFields(c.Query("fields"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("debug projection: endpoint=getProductByID product_id=%d raw_query=%q raw_fields=%q parsed_fields=%v", id, c.Request.URL.RawQuery, c.Query("fields"), fields)

	product, err := h.service.GetProductByID(c.Request.Context(), id, fields)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toProductResponse(product, fields)})
}

func parsePositiveIntWithDefault(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, errors.New("invalid positive integer")
	}

	return value, nil
}

func parseNonNegativeIntWithDefault(raw string, defaultValue int) (int, error) {
	if raw == "" {
		return defaultValue, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 0 {
		return 0, errors.New("invalid non-negative integer")
	}

	return value, nil
}

func parseProjectedFields(raw string) ([]string, error) {
	if strings.TrimSpace(raw) == "" {
		return allowedProductFields, nil
	}

	parts := strings.Split(raw, ",")
	fields := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))

	for _, part := range parts {
		field := strings.TrimSpace(part)
		if field == "" {
			continue
		}
		if !slices.Contains(allowedProductFields, field) {
			return nil, errors.New("invalid fields parameter")
		}
		if _, exists := seen[field]; exists {
			continue
		}
		seen[field] = struct{}{}
		fields = append(fields, field)
	}

	if len(fields) == 0 {
		return nil, errors.New("fields parameter must include at least one field")
	}

	return fields, nil
}

func toProductResponse(product domain.Product, fields []string) gin.H {
	resp := gin.H{}
	for _, field := range fields {
		switch field {
		case "id":
			resp["id"] = product.ID
		case "name":
			resp["name"] = product.Name
		case "description":
			resp["description"] = product.Description
		case "slug":
			resp["slug"] = product.Slug
		case "isActive":
			resp["isActive"] = product.IsActive
		}
	}

	return resp
}
