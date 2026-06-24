package ports

import (
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	GetByID(ctx context.Context, id string) (*domain.Task, error)
	GetAllByUserID(ctx context.Context, userID string) ([]*domain.Task, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string) error
	GetTasksWithDeadlineBefore(ctx context.Context, deadline time.Time) ([]*domain.Task, error)
}

type NoteRepository interface {
	Create(ctx context.Context, note *domain.Note) error
	GetAllByTaskID(ctx context.Context, taskID string) ([]*domain.Note, error)
	Delete(ctx context.Context, id string) error
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *domain.Notification) error
}