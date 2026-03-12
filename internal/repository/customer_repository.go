package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"secret-shop/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer domain.Customer) (domain.Customer, error)
	FindByEmail(ctx context.Context, email string) (domain.Customer, error)
	UpdateLastLoginAt(ctx context.Context, id int64) error
}

type PostgresCustomerRepository struct {
	db *pgxpool.Pool
}

func NewPostgresCustomerRepository(db *pgxpool.Pool) *PostgresCustomerRepository {
	return &PostgresCustomerRepository{db: db}
}

func (r *PostgresCustomerRepository) Create(ctx context.Context, customer domain.Customer) (domain.Customer, error) {
	query := `
		INSERT INTO customers (email, password_hash, firstname, lastname, phone, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, email, password_hash, firstname, lastname, phone, is_active, email_verified_at, last_login_at, created_at, updated_at
	`

	email := strings.ToLower(strings.TrimSpace(customer.Email))

	var created domain.Customer
	if err := r.db.QueryRow(
		ctx,
		query,
		email,
		customer.PasswordHash,
		customer.Firstname,
		customer.Lastname,
		customer.Phone,
		customer.IsActive,
	).Scan(
		&created.ID,
		&created.Email,
		&created.PasswordHash,
		&created.Firstname,
		&created.Lastname,
		&created.Phone,
		&created.IsActive,
		&created.EmailVerifiedAt,
		&created.LastLoginAt,
		&created.CreatedAt,
		&created.UpdatedAt,
	); err != nil {
		return domain.Customer{}, fmt.Errorf("create customer: %w", err)
	}

	return created, nil
}

func (r *PostgresCustomerRepository) FindByEmail(ctx context.Context, email string) (domain.Customer, error) {
	query := `
		SELECT id, email, password_hash, firstname, lastname, phone, is_active, email_verified_at, last_login_at, created_at, updated_at
		FROM customers
		WHERE LOWER(email) = LOWER($1)
	`

	var customer domain.Customer
	if err := r.db.QueryRow(ctx, query, strings.TrimSpace(email)).Scan(
		&customer.ID,
		&customer.Email,
		&customer.PasswordHash,
		&customer.Firstname,
		&customer.Lastname,
		&customer.Phone,
		&customer.IsActive,
		&customer.EmailVerifiedAt,
		&customer.LastLoginAt,
		&customer.CreatedAt,
		&customer.UpdatedAt,
	); err != nil {
		return domain.Customer{}, fmt.Errorf("find customer by email: %w", err)
	}

	return customer, nil
}

func (r *PostgresCustomerRepository) UpdateLastLoginAt(ctx context.Context, id int64) error {
	query := `
		UPDATE customers
		SET last_login_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1
	`

	if _, err := r.db.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("update customer last login: %w", err)
	}

	return nil
}
