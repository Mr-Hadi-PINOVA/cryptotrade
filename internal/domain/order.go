package domain

import (
    "errors"
    "time"
)

// OrderItem represents a product purchase entry within an order.
type OrderItem struct {
    ProductID string `json:"product_id"`
    Quantity  int    `json:"quantity"`
}

// Order represents a customer's purchase order.
type Order struct {
    ID        string      `json:"id"`
    UserID    string      `json:"user_id"`
    Items     []OrderItem `json:"items"`
    Total     float64     `json:"total"`
    CreatedAt time.Time   `json:"created_at"`
}

// Validate ensures the order is well formed.
func (o Order) Validate() error {
    if o.UserID == "" {
        return errors.New("user_id is required")
    }
    if len(o.Items) == 0 {
        return errors.New("order must contain at least one item")
    }
    for _, item := range o.Items {
        if item.ProductID == "" {
            return errors.New("product_id is required for each item")
        }
        if item.Quantity <= 0 {
            return errors.New("item quantity must be positive")
        }
    }
    return nil
}
