package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"secret-shop/internal/domain"
	"secret-shop/internal/repository"
)

var ErrProductNotFound = errors.New("product not found")

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) ListProducts(ctx context.Context, limit, offset int) ([]domain.Product, error) {
	return s.repo.FindAll(ctx, limit, offset)
}

func (s *ProductService) GetProductByID(ctx context.Context, id int64) (domain.Product, error) {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Product{}, ErrProductNotFound
		}
		return domain.Product{}, fmt.Errorf("get product by id: %w", err)
	}

	return product, nil
}
