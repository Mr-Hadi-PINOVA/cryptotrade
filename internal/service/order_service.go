package service

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/repository"
)

// OrderService contains the business logic for orders.
type OrderService struct {
    orders   repository.OrderRepository
    users    repository.UserRepository
    products repository.ProductRepository
}

// NewOrderService creates a new OrderService.
func NewOrderService(orderRepo repository.OrderRepository, userRepo repository.UserRepository, productRepo repository.ProductRepository) *OrderService {
    return &OrderService{orders: orderRepo, users: userRepo, products: productRepo}
}

// CreateOrder creates a new order for the supplied user and items.
func (s *OrderService) CreateOrder(ctx context.Context, userID string, items []domain.OrderItem) (domain.Order, error) {
    order := domain.Order{
        ID:     uuid.NewString(),
        UserID: userID,
        Items:  items,
    }

    if err := order.Validate(); err != nil {
        return domain.Order{}, fmt.Errorf("%w: %w", ErrValidation, err)
    }

    if _, err := s.users.GetByID(ctx, userID); err != nil {
        return domain.Order{}, err
    }

    var total float64
    updatedProducts := make([]domain.Product, 0, len(items))

    for _, item := range items {
        product, err := s.products.GetByID(ctx, item.ProductID)
        if err != nil {
            return domain.Order{}, err
        }
        if product.Stock < item.Quantity {
            return domain.Order{}, fmt.Errorf("%w: insufficient stock for product %s", ErrValidation, product.ID)
        }

        product.Stock -= item.Quantity
        total += product.Price * float64(item.Quantity)
        updatedProducts = append(updatedProducts, product)
    }

    for _, product := range updatedProducts {
        if err := s.products.Update(ctx, product); err != nil {
            return domain.Order{}, err
        }
    }

    order.Total = total
    order.CreatedAt = time.Now().UTC()

    if err := s.orders.Create(ctx, order); err != nil {
        return domain.Order{}, err
    }

    return order, nil
}

// GetOrder retrieves an order by ID.
func (s *OrderService) GetOrder(ctx context.Context, id string) (domain.Order, error) {
    return s.orders.GetByID(ctx, id)
}

// ListOrders returns all orders.
func (s *OrderService) ListOrders(ctx context.Context) ([]domain.Order, error) {
    return s.orders.List(ctx)
}
