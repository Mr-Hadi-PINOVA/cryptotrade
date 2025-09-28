package service

import (
    "context"
    "fmt"

    "github.com/google/uuid"

    "cryptotrade/internal/domain"
    "cryptotrade/internal/repository"
)

// UserService contains the business logic for users.
type UserService struct {
    repo repository.UserRepository
}

// NewUserService creates a new UserService.
func NewUserService(repo repository.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// CreateUser registers a new user.
func (s *UserService) CreateUser(ctx context.Context, input domain.User) (domain.User, error) {
    user := domain.User{
        ID:    uuid.NewString(),
        Name:  input.Name,
        Email: input.Email,
    }

    if err := user.Validate(); err != nil {
        return domain.User{}, fmt.Errorf("%w: %w", ErrValidation, err)
    }

    if _, err := s.repo.GetByEmail(ctx, user.Email); err == nil {
        return domain.User{}, repository.ErrConflict
    } else if err != repository.ErrNotFound {
        return domain.User{}, err
    }

    if err := s.repo.Create(ctx, user); err != nil {
        return domain.User{}, err
    }

    return user, nil
}

// GetUser returns a user by ID.
func (s *UserService) GetUser(ctx context.Context, id string) (domain.User, error) {
    return s.repo.GetByID(ctx, id)
}

// ListUsers returns all registered users.
func (s *UserService) ListUsers(ctx context.Context) ([]domain.User, error) {
    return s.repo.List(ctx)
}
