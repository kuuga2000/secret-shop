package repository

import (
	"context"
	"fmt"
	"log"
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

	normalizedQuery := strings.Join(strings.Fields(query), " ")
	log.Printf(
		"projection: endpoint=listProducts sql_query=%s fields=%v limit=%d offset=%d",
		normalizedQuery,
		fields,
		limit,
		offset,
	)

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

	normalizedQuery := strings.Join(strings.Fields(query), " ")
	log.Printf(
		"projection: endpoint=getProductByID sql_query=%s fields=%v id=%d",
		normalizedQuery,
		fields,
		id,
	)

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
		switch field {
		case "id":
			columns = append(columns, "id")
		case "name":
			columns = append(columns, "name")
		case "description":
			columns = append(columns, "description")
		case "slug":
			columns = append(columns, "slug")
		case "isActive":
			columns = append(columns, "is_active")
		default:
			return nil, fmt.Errorf("unsupported product field: %s", field)
		}
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("at least one projected field is required")
	}

	return columns, nil
}

func buildScanTargets(product *domain.Product, fields []string) []any {
	targets := make([]any, 0, len(fields))
	for _, field := range fields {
		switch field {
		case "id":
			targets = append(targets, &product.ID)
		case "name":
			targets = append(targets, &product.Name)
		case "description":
			targets = append(targets, &product.Description)
		case "slug":
			targets = append(targets, &product.Slug)
		case "isActive":
			targets = append(targets, &product.IsActive)
		}
	}

	return targets
}
