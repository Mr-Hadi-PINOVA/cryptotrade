package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/repository"
)

// ProductService contains the business logic for products.
type ProductService struct {
    repo repository.ProductRepository
}

// NewProductService creates a new ProductService.
func NewProductService(repo repository.ProductRepository) *ProductService {
    return &ProductService{repo: repo}
}

// CreateProduct persists a new product.
func (s *ProductService) CreateProduct(ctx context.Context, input domain.Product) (domain.Product, error) {
    product := domain.Product{
        ID:          uuid.NewString(),
        Name:        input.Name,
        Description: input.Description,
        Price:       input.Price,
        Stock:       input.Stock,
    }

    if err := product.Validate(); err != nil {
        return domain.Product{}, fmt.Errorf("%w: %w", ErrValidation, err)
    }

    if err := s.repo.Create(ctx, product); err != nil {
        return domain.Product{}, err
    }

    return product, nil
}

// UpdateProduct updates an existing product by ID.
func (s *ProductService) UpdateProduct(ctx context.Context, id string, input domain.Product) (domain.Product, error) {
    product, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return domain.Product{}, err
    }

    product.Name = input.Name
    product.Description = input.Description
    product.Price = input.Price
    product.Stock = input.Stock

    if err := product.Validate(); err != nil {
        return domain.Product{}, fmt.Errorf("%w: %w", ErrValidation, err)
    }

    if err := s.repo.Update(ctx, product); err != nil {
        return domain.Product{}, err
    }

    return product, nil
}

// DeleteProduct removes a product by ID.
func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
    return s.repo.Delete(ctx, id)
}

// GetProduct returns a product by ID.
func (s *ProductService) GetProduct(ctx context.Context, id string) (domain.Product, error) {
    return s.repo.GetByID(ctx, id)
}

// ListProducts returns all products.
func (s *ProductService) ListProducts(ctx context.Context) ([]domain.Product, error) {
    return s.repo.List(ctx)
}
