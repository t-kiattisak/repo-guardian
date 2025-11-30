package repository

import (
	"context"
	"errors"
	"sync"

	"repo-guardian/internal/domain"
)

type memoryUserRepository struct {
	mu    sync.RWMutex
	users map[int64]*domain.User
}

func NewMemoryUserRepository() domain.UserRepository {
	return &memoryUserRepository{
		users: make(map[int64]*domain.User),
	}
}

func (r *memoryUserRepository) Create(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; exists {
		return errors.New("user already exists")
	}

	r.users[user.ID] = user
	return nil
}

func (r *memoryUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}

	return user, nil
}
