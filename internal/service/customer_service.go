package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"

	"secret-shop/internal/domain"
	"secret-shop/internal/repository"
)

var ErrCustomerNotFound = errors.New("customer not found")
var ErrIncorrectPassword = errors.New("incorrect password")

type CustomerService struct {
	customers repository.CustomerRepository
}

type UpdateCustomerProfileInput struct {
	Firstname string
	Lastname  string
	Phone     string
}

type ChangeCustomerPasswordInput struct {
	CurrentPassword string
	NewPassword     string
}

func NewCustomerService(customers repository.CustomerRepository) *CustomerService {
	return &CustomerService{customers: customers}
}

func (s *CustomerService) GetByID(ctx context.Context, customerID int64) (domain.Customer, error) {
	customer, err := s.customers.FindByID(ctx, customerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Customer{}, ErrCustomerNotFound
		}
		return domain.Customer{}, fmt.Errorf("get customer by id: %w", err)
	}

	if !customer.IsActive {
		return domain.Customer{}, ErrCustomerInactive
	}

	return customer, nil
}

func (s *CustomerService) UpdateProfile(ctx context.Context, customerID int64, input UpdateCustomerProfileInput) (domain.Customer, error) {
	customer, err := s.GetByID(ctx, customerID)
	if err != nil {
		return domain.Customer{}, err
	}

	customer.Firstname = strings.TrimSpace(input.Firstname)
	customer.Lastname = strings.TrimSpace(input.Lastname)
	customer.Phone = strings.TrimSpace(input.Phone)

	updated, err := s.customers.UpdateProfile(ctx, customer)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Customer{}, ErrCustomerNotFound
		}
		return domain.Customer{}, fmt.Errorf("update customer profile: %w", err)
	}

	return updated, nil
}

func (s *CustomerService) ChangePassword(ctx context.Context, customerID int64, input ChangeCustomerPasswordInput) error {
	customer, err := s.GetByID(ctx, customerID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return ErrIncorrectPassword
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.customers.UpdatePassword(ctx, customerID, string(passwordHash)); err != nil {
		return fmt.Errorf("change customer password: %w", err)
	}

	return nil
}
