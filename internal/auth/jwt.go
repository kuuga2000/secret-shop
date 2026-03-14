package auth

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"secret-shop/internal/domain"
)

var ErrInvalidToken = errors.New("invalid token")

type AccessTokenClaims struct {
	CustomerID int64
	Email      string
}

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

func ParseAccessToken(tokenString, secret, issuer string) (AccessTokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	}, jwt.WithIssuer(issuer))
	if err != nil {
		return AccessTokenClaims{}, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return AccessTokenClaims{}, ErrInvalidToken
	}

	subject, err := claims.GetSubject()
	if err != nil {
		return AccessTokenClaims{}, ErrInvalidToken
	}

	customerID, err := strconv.ParseInt(subject, 10, 64)
	if err != nil || customerID <= 0 {
		return AccessTokenClaims{}, ErrInvalidToken
	}

	email, _ := claims["email"].(string)

	return AccessTokenClaims{
		CustomerID: customerID,
		Email:      email,
	}, nil
}
