package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"secret-shop/internal/service"
)

type AuthService interface {
	Register(ctx context.Context, input service.RegisterCustomerInput) (service.AuthResult, error)
	Login(ctx context.Context, input service.LoginCustomerInput) (service.AuthResult, error)
}

type AuthHandler struct {
	service AuthService
}

type registerRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Phone     string `json:"phone"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(service AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) RegisterRoutes(v1 *gin.RouterGroup) {
	auth := v1.Group("/auth")
	auth.POST("/register", h.register)
	auth.POST("/login", h.login)
}

func (h *AuthHandler) register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}

	result, err := h.service.Register(c.Request.Context(), service.RegisterCustomerInput{
		Email:     req.Email,
		Password:  req.Password,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Phone:     req.Phone,
	})
	if err != nil {
		if errors.Is(err, service.ErrCustomerAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": "customer already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register customer"})
		return
	}

	c.JSON(http.StatusCreated, authResponse(result))
}

func (h *AuthHandler) login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	result, err := h.service.Login(c.Request.Context(), service.LoginCustomerInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidCredentials):
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case errors.Is(err, service.ErrCustomerInactive):
			c.JSON(http.StatusForbidden, gin.H{"error": "customer is inactive"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		}
		return
	}

	c.JSON(http.StatusOK, authResponse(result))
}

func authResponse(result service.AuthResult) gin.H {
	response := gin.H{
		"accessToken": result.AccessToken,
		"customer": gin.H{
			"id":        result.Customer.ID,
			"email":     result.Customer.Email,
			"firstname": result.Customer.Firstname,
			"lastname":  result.Customer.Lastname,
			"phone":     result.Customer.Phone,
			"isActive":  result.Customer.IsActive,
		},
	}

	if result.Customer.LastLoginAt != nil {
		response["customer"].(gin.H)["lastLoginAt"] = result.Customer.LastLoginAt
	}

	return gin.H{"data": response}
}
