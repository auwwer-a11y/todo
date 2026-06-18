package service

import (
	"context"
	"time"

	"github.com/auwwer-a11y/todo/internal/domain"
)

type mockUserRepo struct {
	users map[string]*domain.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*domain.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *domain.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if u, ok := m.users[email]; ok {
		return u, nil
	}
	return nil, nil
}

type mockCache struct{}

func (m *mockCache) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return nil
}

func (m *mockCache) Get(ctx context.Context, key string) (string, error) {
	return "", nil
}

func (m *mockCache) Delete(ctx context.Context, key string) error {
	return nil
}
