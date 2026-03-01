package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"secret-shop/internal/domain"
)

type ProductRepository interface {
	FindAll(ctx context.Context, limit, offset int) ([]domain.Product, error)
	FindByID(ctx context.Context, id int64) (domain.Product, error)
}

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) FindAll(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	const query = `
		SELECT id, name, description, slug, is_active
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0, limit)
	for rows.Next() {
		var p domain.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Slug, &p.IsActive); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return products, nil
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64) (domain.Product, error) {
	const query = `
		SELECT id, name, description, slug, is_active
		FROM products
		WHERE id = $1
	`

	var p domain.Product
	if err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.Slug, &p.IsActive); err != nil {
		return domain.Product{}, fmt.Errorf("query product by id: %w", err)
	}

	return p, nil
}
