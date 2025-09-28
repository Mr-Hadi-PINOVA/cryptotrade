package memory

import (
    "context"
    "sync"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/repository"
)

// ProductRepository is an in-memory implementation of repository.ProductRepository.
type ProductRepository struct {
    mu       sync.RWMutex
    products map[string]domain.Product
}

// NewProductRepository constructs a new in-memory product repository.
func NewProductRepository() *ProductRepository {
    return &ProductRepository{products: make(map[string]domain.Product)}
}

func (r *ProductRepository) Create(_ context.Context, product domain.Product) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.products[product.ID]; exists {
        return repository.ErrConflict
    }

    r.products[product.ID] = product
    return nil
}

func (r *ProductRepository) Update(_ context.Context, product domain.Product) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.products[product.ID]; !ok {
        return repository.ErrNotFound
    }
    r.products[product.ID] = product
    return nil
}

func (r *ProductRepository) Delete(_ context.Context, id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, ok := r.products[id]; !ok {
        return repository.ErrNotFound
    }
    delete(r.products, id)
    return nil
}

func (r *ProductRepository) GetByID(_ context.Context, id string) (domain.Product, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    product, ok := r.products[id]
    if !ok {
        return domain.Product{}, repository.ErrNotFound
    }
    return product, nil
}

func (r *ProductRepository) List(_ context.Context) ([]domain.Product, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    products := make([]domain.Product, 0, len(r.products))
    for _, product := range r.products {
        products = append(products, product)
    }
    return products, nil
}

// UserRepository is an in-memory implementation of repository.UserRepository.
type UserRepository struct {
    mu    sync.RWMutex
    users map[string]domain.User
}

// NewUserRepository constructs a new in-memory user repository.
func NewUserRepository() *UserRepository {
    return &UserRepository{users: make(map[string]domain.User)}
}

func (r *UserRepository) Create(_ context.Context, user domain.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.users[user.ID]; exists {
        return repository.ErrConflict
    }

    for _, existing := range r.users {
        if existing.Email == user.Email {
            return repository.ErrConflict
        }
    }

    r.users[user.ID] = user
    return nil
}

func (r *UserRepository) GetByID(_ context.Context, id string) (domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    user, ok := r.users[id]
    if !ok {
        return domain.User{}, repository.ErrNotFound
    }
    return user, nil
}

func (r *UserRepository) GetByEmail(_ context.Context, email string) (domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    for _, user := range r.users {
        if user.Email == email {
            return user, nil
        }
    }
    return domain.User{}, repository.ErrNotFound
}

func (r *UserRepository) List(_ context.Context) ([]domain.User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    users := make([]domain.User, 0, len(r.users))
    for _, user := range r.users {
        users = append(users, user)
    }
    return users, nil
}

// OrderRepository is an in-memory implementation of repository.OrderRepository.
type OrderRepository struct {
    mu     sync.RWMutex
    orders map[string]domain.Order
}

// NewOrderRepository constructs a new in-memory order repository.
func NewOrderRepository() *OrderRepository {
    return &OrderRepository{orders: make(map[string]domain.Order)}
}

func (r *OrderRepository) Create(_ context.Context, order domain.Order) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.orders[order.ID]; exists {
        return repository.ErrConflict
    }

    r.orders[order.ID] = order
    return nil
}

func (r *OrderRepository) GetByID(_ context.Context, id string) (domain.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    order, ok := r.orders[id]
    if !ok {
        return domain.Order{}, repository.ErrNotFound
    }
    return order, nil
}

func (r *OrderRepository) List(_ context.Context) ([]domain.Order, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()

    orders := make([]domain.Order, 0, len(r.orders))
    for _, order := range r.orders {
        orders = append(orders, order)
    }
    return orders, nil
}
