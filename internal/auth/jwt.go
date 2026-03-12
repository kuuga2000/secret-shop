package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"secret-shop/internal/domain"
)

func GenerateAccessToken(customer domain.Customer, secret, issuer string, ttl time.Duration) (string, error) {
	now := time.Now().UTC()

	claims := jwt.MapClaims{
		"sub":   fmt.Sprintf("%d", customer.ID),
		"email": customer.Email,
		"iss":   issuer,
		"iat":   now.Unix(),
		"exp":   now.Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
