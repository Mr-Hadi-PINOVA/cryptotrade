package domain

import (
    "errors"
    "regexp"
)

var emailRegex = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)

// User represents a customer account in the system.
type User struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}

// Validate ensures the user is well formed.
func (u User) Validate() error {
    if u.Name == "" {
        return errors.New("name is required")
    }
    if !emailRegex.MatchString(u.Email) {
        return errors.New("email is invalid")
    }
    return nil
}
