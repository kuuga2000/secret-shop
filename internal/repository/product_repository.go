package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"secret-shop/internal/domain"
)

type ProductRepository interface {
	FindAll(ctx context.Context, limit, offset int, fields []string) ([]domain.Product, error)
	FindByID(ctx context.Context, id int64, fields []string) (domain.Product, error)
}

type PostgresProductRepository struct {
	db *pgxpool.Pool
}

type productFieldBinding struct {
	column string
	target func(*domain.Product) any
}

var productFieldBindings = map[string]productFieldBinding{
	"id": {
		column: "id",
		target: func(p *domain.Product) any { return &p.ID },
	},
	"name": {
		column: "name",
		target: func(p *domain.Product) any { return &p.Name },
	},
	"description": {
		column: "description",
		target: func(p *domain.Product) any { return &p.Description },
	},
	"slug": {
		column: "slug",
		target: func(p *domain.Product) any { return &p.Slug },
	},
	"isActive": {
		column: "is_active",
		target: func(p *domain.Product) any { return &p.IsActive },
	},
}

func NewPostgresProductRepository(db *pgxpool.Pool) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) FindAll(ctx context.Context, limit, offset int, fields []string) ([]domain.Product, error) {
	selectColumns, err := buildProductSelectColumns(fields)
	if err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM products
		ORDER BY id
		LIMIT $1 OFFSET $2
	`, strings.Join(selectColumns, ", "))

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("query products: %w", err)
	}
	defer rows.Close()

	products := make([]domain.Product, 0, limit)
	for rows.Next() {
		var p domain.Product
		scanTargets := buildScanTargets(&p, fields)
		if err := rows.Scan(scanTargets...); err != nil {
			return nil, fmt.Errorf("scan product row: %w", err)
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate product rows: %w", err)
	}

	return products, nil
}

func (r *PostgresProductRepository) FindByID(ctx context.Context, id int64, fields []string) (domain.Product, error) {
	selectColumns, err := buildProductSelectColumns(fields)
	if err != nil {
		return domain.Product{}, err
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM products
		WHERE id = $1
	`, strings.Join(selectColumns, ", "))

	var p domain.Product
	scanTargets := buildScanTargets(&p, fields)
	if err := r.db.QueryRow(ctx, query, id).Scan(scanTargets...); err != nil {
		return domain.Product{}, fmt.Errorf("query product by id: %w", err)
	}

	return p, nil
}

func buildProductSelectColumns(fields []string) ([]string, error) {
	columns := make([]string, 0, len(fields))
	for _, field := range fields {
		binding, ok := productFieldBindings[field]
		if !ok {
			return nil, fmt.Errorf("unsupported product field: %s", field)
		}
		columns = append(columns, binding.column)
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("at least one projected field is required")
	}

	return columns, nil
}

func buildScanTargets(product *domain.Product, fields []string) []any {
	targets := make([]any, 0, len(fields))
	for _, field := range fields {
		binding, ok := productFieldBindings[field]
		if ok {
			targets = append(targets, binding.target(product))
		}
	}

	return targets
}
