package repository

import (
    "context"
    "errors"

    "cryptotrade/internal/domain"
)

// ErrNotFound is returned when an entity cannot be located.
var ErrNotFound = errors.New("entity not found")

// ErrConflict is returned when an entity would violate uniqueness constraints.
var ErrConflict = errors.New("entity already exists")

// ProductRepository describes the persistence operations for products.
type ProductRepository interface {
    Create(ctx context.Context, product domain.Product) error
    Update(ctx context.Context, product domain.Product) error
    Delete(ctx context.Context, id string) error
    GetByID(ctx context.Context, id string) (domain.Product, error)
    List(ctx context.Context) ([]domain.Product, error)
}

// UserRepository describes persistence operations for users.
type UserRepository interface {
    Create(ctx context.Context, user domain.User) error
    GetByID(ctx context.Context, id string) (domain.User, error)
    GetByEmail(ctx context.Context, email string) (domain.User, error)
    List(ctx context.Context) ([]domain.User, error)
}

// OrderRepository describes persistence operations for orders.
type OrderRepository interface {
    Create(ctx context.Context, order domain.Order) error
    GetByID(ctx context.Context, id string) (domain.Order, error)
    List(ctx context.Context) ([]domain.Order, error)
}
