package service

import (
	"time"
	"log/slog"
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/auwwer-a11y/todo/internal/ports"
	"github.com/auwwer-a11y/todo/pkg/jwt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)


type UserService struct {
	userRepo ports.UserRepository
	cacheRepo ports.CacheRepository
	logger *slog.Logger
	jwtSecret string
	jwtTTL time.Duration
}

func NewUserService(userRepo ports.UserRepository, cacheRepo ports.CacheRepository, logger *slog.Logger, jwtSecret string, jwtTTL time.Duration) *UserService {
	return &UserService{
		userRepo: userRepo,
		cacheRepo: cacheRepo,
		logger: logger,
		jwtSecret: jwtSecret,
		jwtTTL: jwtTTL,
	}
}

func (s *UserService) Register(ctx context.Context, name string, email string, password string) (*domain.User, error) {
	if user, _ := s.userRepo.GetByEmail(ctx, email); user != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		ID: uuid.New().String(),
		Name: name,
		Email: email,
		PasswordHash: string(hashedPassword),
		CreatedAt: time.Now(),
	}
	return s.userRepo.Create(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", domain.ErrUserNotFound
	}

	passwordErr := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if passwordErr != nil {
		return "", domain.ErrInvalidCredentials
	}

	token, err := jwt.Generate(user.ID, s.jwtSecret, s.jwtTTL)
	if err != nil {
		return "", err
	}

	return token, nil
}
