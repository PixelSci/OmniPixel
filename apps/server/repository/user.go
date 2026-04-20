package repository

import (
	"context"

	"omni-pixel/model"
)

// UserRepository is the storage contract used by the service layer.
// Implementations must translate backend-specific errors into the sentinel
// values declared in errors.go (ErrNotFound, ErrDuplicate).
type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]*model.User, int, error)
}
