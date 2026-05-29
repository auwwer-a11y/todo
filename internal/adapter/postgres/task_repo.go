package postgres

import (
	"github.com/jmoiron/sqlx"
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
)

type TaskRepo struct {
	db *sqlx.DB
}

func NewTaskRepo(db *sqlx.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tasks (id, user_id, title, description, status, deadline, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`, task.ID, task.UserID, task.Title, task.Description, task.Status, task.Deadline, task.CreatedAt, task.UpdatedAt)
	return err
}

func (r *TaskRepo) GetByID(ctx context.Context, id string) (*domain.Task, error) {
	var task domain.Task
	err := r.db.GetContext(ctx, &task, `
		SELECT id, user_id, title, description, status, deadline, created_at, updated_at
		FROM tasks
		WHERE id = $1
		`, id)
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepo) GetAllByUserID(ctx context.Context, userID string) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err:= r.db.SelectContext(ctx, &tasks, `
	SELECT id, user_id, title, description, status, deadline, created_at, updated_at
	FROM tasks
	WHERE user_id = $1
	`, userID)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (r *TaskRepo) Update(ctx context.Context, task *domain.Task) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, deadline = $4, updated_at = $5
		WHERE id = $6
	`, task.Title, task.Description, task.Status, task.Deadline, task.UpdatedAt, task.ID)
	return err
}

func (r *TaskRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM tasks
		WHERE id = $1
	`, id)
	return err
}