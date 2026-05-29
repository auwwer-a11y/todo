package postgres

import (
	"github.com/jmoiron/sqlx"
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	_, err := r.db.ExecContext(ctx,`
		INSERT INTO users (id, name, password_hash, email, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`, user.ID, user.Name, user.PasswordHash, user.Email, user.CreatedAt)
	return err
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, name, password_hash, email, created_at
		FROM users
		WHERE id = $1
		`, id)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.GetContext(ctx, &user, `
		SELECT id, name, password_hash, email, created_at
		FROM users
		WHERE email = $1
		`, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}