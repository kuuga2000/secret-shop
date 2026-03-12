package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"

	appauth "secret-shop/internal/auth"
	"secret-shop/internal/domain"
	"secret-shop/internal/repository"
)

var ErrCustomerAlreadyExists = errors.New("customer already exists")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrCustomerInactive = errors.New("customer inactive")

type AuthService struct {
	customers      repository.CustomerRepository
	jwtSecret      string
	jwtIssuer      string
	accessTokenTTL time.Duration
}

type RegisterCustomerInput struct {
	Email     string
	Password  string
	Firstname string
	Lastname  string
	Phone     string
}

type LoginCustomerInput struct {
	Email    string
	Password string
}

type AuthResult struct {
	Customer    domain.Customer
	AccessToken string
}

func NewAuthService(customers repository.CustomerRepository, jwtSecret, jwtIssuer string, accessTokenTTL time.Duration) *AuthService {
	return &AuthService{
		customers:      customers,
		jwtSecret:      jwtSecret,
		jwtIssuer:      jwtIssuer,
		accessTokenTTL: accessTokenTTL,
	}
}

func (s *AuthService) Register(ctx context.Context, input RegisterCustomerInput) (AuthResult, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return AuthResult{}, fmt.Errorf("hash password: %w", err)
	}

	customer, err := s.customers.Create(ctx, domain.Customer{
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		PasswordHash: string(passwordHash),
		Firstname:    strings.TrimSpace(input.Firstname),
		Lastname:     strings.TrimSpace(input.Lastname),
		Phone:        strings.TrimSpace(input.Phone),
		IsActive:     true,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return AuthResult{}, ErrCustomerAlreadyExists
		}
		return AuthResult{}, fmt.Errorf("register customer: %w", err)
	}

	token, err := appauth.GenerateAccessToken(customer, s.jwtSecret, s.jwtIssuer, s.accessTokenTTL)
	if err != nil {
		return AuthResult{}, fmt.Errorf("generate access token: %w", err)
	}

	return AuthResult{
		Customer:    customer,
		AccessToken: token,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, input LoginCustomerInput) (AuthResult, error) {
	customer, err := s.customers.FindByEmail(ctx, input.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return AuthResult{}, ErrInvalidCredentials
		}
		return AuthResult{}, fmt.Errorf("load customer for login: %w", err)
	}

	if !customer.IsActive {
		return AuthResult{}, ErrCustomerInactive
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(input.Password)); err != nil {
		return AuthResult{}, ErrInvalidCredentials
	}

	if err := s.customers.UpdateLastLoginAt(ctx, customer.ID); err != nil {
		return AuthResult{}, fmt.Errorf("update last login: %w", err)
	}

	customer.LastLoginAt = ptrTime(time.Now().UTC())

	token, err := appauth.GenerateAccessToken(customer, s.jwtSecret, s.jwtIssuer, s.accessTokenTTL)
	if err != nil {
		return AuthResult{}, fmt.Errorf("generate access token: %w", err)
	}

	return AuthResult{
		Customer:    customer,
		AccessToken: token,
	}, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}

	return false
}

func ptrTime(t time.Time) *time.Time {
	return &t
}
