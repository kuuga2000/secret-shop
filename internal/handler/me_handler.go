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

type CustomerService interface {
	GetByID(ctx context.Context, customerID int64) (domain.Customer, error)
	UpdateProfile(ctx context.Context, customerID int64, input service.UpdateCustomerProfileInput) (domain.Customer, error)
	ChangePassword(ctx context.Context, customerID int64, input service.ChangeCustomerPasswordInput) error
}

type MeHandler struct {
	service CustomerService
}

type updateMeRequest struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type changePasswordRequest struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

func NewMeHandler(service CustomerService) *MeHandler {
	return &MeHandler{service: service}
}

func (h *MeHandler) RegisterRoutes(v1 *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	me := v1.Group("/me")
	me.Use(authMiddleware)
	me.GET("", h.getMe)
	me.PATCH("", h.updateMe)
	me.PATCH("/password", h.changePassword)
}

func (h *MeHandler) getMe(c *gin.Context) {
	customerID, ok := customerIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	customer, err := h.service.GetByID(c.Request.Context(), customerID)
	if err != nil {
		writeCustomerServiceError(c, err, "failed to load profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toCustomerResponse(customer)})
}

func (h *MeHandler) updateMe(c *gin.Context) {
	customerID, ok := customerIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req updateMeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	customer, err := h.service.UpdateProfile(c.Request.Context(), customerID, service.UpdateCustomerProfileInput{
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Phone:     req.Phone,
	})
	if err != nil {
		writeCustomerServiceError(c, err, "failed to update profile")
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": toCustomerResponse(customer)})
}

func (h *MeHandler) changePassword(c *gin.Context) {
	customerID, ok := customerIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.CurrentPassword) == "" || strings.TrimSpace(req.NewPassword) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "currentPassword and newPassword are required"})
		return
	}

	if len(req.NewPassword) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "newPassword must be at least 8 characters"})
		return
	}

	err := h.service.ChangePassword(c.Request.Context(), customerID, service.ChangeCustomerPasswordInput{
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrIncorrectPassword):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "current password is incorrect"})
		case errors.Is(err, service.ErrCustomerInactive):
			c.JSON(http.StatusForbidden, gin.H{"error": "customer is inactive"})
		case errors.Is(err, service.ErrCustomerNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to change password"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": gin.H{"message": "password updated"}})
}

func writeCustomerServiceError(c *gin.Context, err error, fallbackMessage string) {
	switch {
	case errors.Is(err, service.ErrCustomerInactive):
		c.JSON(http.StatusForbidden, gin.H{"error": "customer is inactive"})
	case errors.Is(err, service.ErrCustomerNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": fallbackMessage})
	}
}
