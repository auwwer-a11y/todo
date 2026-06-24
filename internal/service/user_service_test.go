package service

import (
	"context"
	"testing"
	"time"

	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestRegister_Success(t *testing.T) {
	repo := newMockUserRepo()
	cache := &mockCache{}
	svc := NewUserService(repo, cache, nil, "secret", 24*time.Hour)

	user, err := svc.Register(context.Background(), "Test User", "test@test.com", "password123")

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@test.com", user.Email)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	repo := newMockUserRepo()
	cache := &mockCache{}
	svc := NewUserService(repo, cache, nil, "secret", 24*time.Hour)

	_, err := svc.Register(context.Background(), "Test User", "test@test.com", "password123")
	assert.NoError(t, err)

	_, err = svc.Register(context.Background(), "Test User 2", "test@test.com", "password123")
	assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
}
