package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"

	"omni-pixel/internal/apperr"
	"omni-pixel/internal/model"
	"omni-pixel/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, in CreateUserInput) (*model.User, error)
	Get(ctx context.Context, id string) (*model.User, error)
	List(ctx context.Context, page, pageSize int) ([]*model.User, int, error)
}

type CreateUserInput struct {
	Name  string
	Email string
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(ctx context.Context, in CreateUserInput) (*model.User, error) {
	email := strings.ToLower(strings.TrimSpace(in.Email))
	now := time.Now().UTC()
	u := &model.User{
		ID:        uuid.NewString(),
		Name:      strings.TrimSpace(in.Name),
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.repo.Create(ctx, u); err != nil {
		if errors.Is(err, repository.ErrDuplicate) {
			return nil, apperr.ErrUserExists.WithDetail("email", email).WithCause(err)
		}
		return nil, apperr.ErrInternal.WithCause(err)
	}
	return u, nil
}

func (s *userService) Get(ctx context.Context, id string) (*model.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, apperr.ErrUserNotFound.WithDetail("id", id).WithCause(err)
		}
		return nil, apperr.ErrInternal.WithCause(err)
	}
	return u, nil
}

func (s *userService) List(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	users, total, err := s.repo.List(ctx, offset, pageSize)
	if err != nil {
		return nil, 0, apperr.ErrInternal.WithCause(err)
	}
	return users, total, nil
}
