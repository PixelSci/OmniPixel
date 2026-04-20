package repository

import (
	"context"
	"sort"
	"sync"

	"omni-pixel/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, u *model.User) error
	FindByID(ctx context.Context, id string) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	List(ctx context.Context, offset, limit int) ([]*model.User, int, error)
}

type inMemoryUserRepo struct {
	mu   sync.RWMutex
	byID map[string]*model.User
}

func NewInMemoryUserRepo() UserRepository {
	return &inMemoryUserRepo{byID: make(map[string]*model.User)}
}

func (r *inMemoryUserRepo) Create(_ context.Context, u *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, existing := range r.byID {
		if existing.Email == u.Email {
			return ErrDuplicate
		}
	}
	clone := *u
	r.byID[u.ID] = &clone
	return nil
}

func (r *inMemoryUserRepo) FindByID(_ context.Context, id string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.byID[id]
	if !ok {
		return nil, ErrNotFound
	}
	clone := *u
	return &clone, nil
}

func (r *inMemoryUserRepo) FindByEmail(_ context.Context, email string) (*model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.byID {
		if u.Email == email {
			clone := *u
			return &clone, nil
		}
	}
	return nil, ErrNotFound
}

func (r *inMemoryUserRepo) List(_ context.Context, offset, limit int) ([]*model.User, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*model.User, 0, len(r.byID))
	for _, u := range r.byID {
		clone := *u
		all = append(all, &clone)
	}
	sort.Slice(all, func(i, j int) bool { return all[i].CreatedAt.Before(all[j].CreatedAt) })

	total := len(all)
	if offset >= total {
		return []*model.User{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}
