package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"secret-shop/internal/domain"
	"secret-shop/internal/repository"
)

var ErrVariantNotFound = errors.New("variant not found")

type VariantService struct {
	repo repository.VariantRepository
}

func NewVariantService(repo repository.VariantRepository) *VariantService {
	return &VariantService{repo: repo}
}

func (s *VariantService) GetVariantBySKU(ctx context.Context, sku string) (domain.Variant, error) {
	variant, err := s.repo.FindBySKU(ctx, sku)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Variant{}, ErrVariantNotFound
		}
		return domain.Variant{}, fmt.Errorf("get variant by sku: %w", err)
	}

	return variant, nil
}
