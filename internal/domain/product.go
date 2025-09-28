package domain

import "errors"

// Product represents a product that can be purchased.
type Product struct {
    ID          string  `json:"id"`
    Name        string  `json:"name"`
    Description string  `json:"description"`
    Price       float64 `json:"price"`
    Stock       int     `json:"stock"`
}

// Validate ensures the product is well formed before persistence.
func (p Product) Validate() error {
    if p.Name == "" {
        return errors.New("name is required")
    }
    if p.Price <= 0 {
        return errors.New("price must be positive")
    }
    if p.Stock < 0 {
        return errors.New("stock cannot be negative")
    }
    return nil
}
