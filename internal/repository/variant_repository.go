package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"secret-shop/internal/domain"
)

type VariantRepository interface {
	FindBySKU(ctx context.Context, sku string) (domain.Variant, error)
}

type PostgresVariantRepository struct {
	db *pgxpool.Pool
}

func NewPostgresVariantRepository(db *pgxpool.Pool) *PostgresVariantRepository {
	return &PostgresVariantRepository{db: db}
}

func (r *PostgresVariantRepository) FindBySKU(ctx context.Context, sku string) (domain.Variant, error) {
	query := `
		SELECT
			pv.id,
			pv.product_id,
			pv.sku,
			pv.price,
			pv.stock,
			p.id,
			p.name,
			p.description,
			p.slug,
			p.is_active
		FROM product_variants pv
		JOIN products p ON p.id = pv.product_id
		WHERE pv.sku = $1
	`

	log.Printf("[sql] FindVariantBySKU sku=%s query=%q", sku, compactSQL(query))

	var v domain.Variant
	if err := r.db.QueryRow(ctx, query, sku).Scan(
		&v.ID,
		&v.ProductID,
		&v.SKU,
		&v.Price,
		&v.Stock,
		&v.Product.ID,
		&v.Product.Name,
		&v.Product.Description,
		&v.Product.Slug,
		&v.Product.IsActive,
	); err != nil {
		return domain.Variant{}, fmt.Errorf("query variant by sku: %w", err)
	}

	return v, nil
}
