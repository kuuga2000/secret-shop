package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"secret-shop/internal/domain"
	"secret-shop/internal/service"
)

type VariantService interface {
	GetVariantBySKU(ctx context.Context, sku string) (domain.Variant, error)
}

type VariantHandler struct {
	service VariantService
}

func NewVariantHandler(service VariantService) *VariantHandler {
	return &VariantHandler{service: service}
}

func (h *VariantHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	v1.GET("/variants/:sku", h.getVariantBySKU)
}

func (h *VariantHandler) getVariantBySKU(c *gin.Context) {
	sku := strings.TrimSpace(c.Param("sku"))
	if sku == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "sku is required"})
		return
	}

	variant, err := h.service.GetVariantBySKU(c.Request.Context(), sku)
	if err != nil {
		if errors.Is(err, service.ErrVariantNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "variant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"id":        variant.ID,
			"productId": variant.ProductID,
			"sku":       variant.SKU,
			"price":     variant.Price,
			"stock":     variant.Stock,
			"product": gin.H{
				"id":          variant.Product.ID,
				"name":        variant.Product.Name,
				"description": variant.Product.Description,
				"slug":        variant.Product.Slug,
				"isActive":    variant.Product.IsActive,
			},
		},
	})
}
