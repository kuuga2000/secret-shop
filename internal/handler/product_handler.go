package handler

import (
	"context"
	"errors"
	"net/http"
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
	ListProducts(ctx context.Context, limit, offset int) ([]domain.Product, error)
	GetProductByID(ctx context.Context, id int64) (domain.Product, error)
}

type ProductHandler struct {
	service ProductService
}

type productResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Slug        string `json:"slug"`
	IsActive    bool   `json:"isActive"`
}

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

	products, err := h.service.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list products"})
		return
	}

	payload := make([]productResponse, 0, len(products))
	for _, product := range products {
		payload = append(payload, toProductResponse(product))
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

	product, err := h.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrProductNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toProductResponse(product)})
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

func toProductResponse(product domain.Product) productResponse {
	return productResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Slug:        product.Slug,
		IsActive:    product.IsActive,
	}
}
