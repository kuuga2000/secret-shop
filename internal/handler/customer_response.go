package handler

import (
	"github.com/gin-gonic/gin"

	"secret-shop/internal/domain"
)

func toCustomerResponse(customer domain.Customer) gin.H {
	response := gin.H{
		"id":        customer.ID,
		"email":     customer.Email,
		"firstname": customer.Firstname,
		"lastname":  customer.Lastname,
		"phone":     customer.Phone,
		"isActive":  customer.IsActive,
	}

	if customer.LastLoginAt != nil {
		response["lastLoginAt"] = customer.LastLoginAt
	}

	if customer.EmailVerifiedAt != nil {
		response["emailVerifiedAt"] = customer.EmailVerifiedAt
	}

	return response
}
