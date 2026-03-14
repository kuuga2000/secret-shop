package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	appauth "secret-shop/internal/auth"
)

const customerIDContextKey = "customerID"

func AuthMiddleware(secret, issuer string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := strings.TrimSpace(c.GetHeader("Authorization"))
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			return
		}

		tokenString, ok := strings.CutPrefix(header, "Bearer ")
		if !ok || strings.TrimSpace(tokenString) == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header"})
			return
		}

		claims, err := appauth.ParseAccessToken(strings.TrimSpace(tokenString), secret, issuer)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		c.Set(customerIDContextKey, claims.CustomerID)
		c.Next()
	}
}

func customerIDFromContext(c *gin.Context) (int64, bool) {
	value, ok := c.Get(customerIDContextKey)
	if !ok {
		return 0, false
	}

	customerID, ok := value.(int64)
	return customerID, ok
}
