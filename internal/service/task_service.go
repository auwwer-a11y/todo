package service

import (
	"github.com/auwwer-a11y/todo/internal/ports"
	"log/slog"
	"context"
	"time"
	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/google/uuid"
)

type TaskService struct {
	taskRepo ports.TaskRepository
	noteRepo ports.NoteRepository
	eventPublisher ports.EventPublisher
	logger *slog.Logger
}

func NewTaskService(taskRepo ports.TaskRepository, noteRepo ports.NoteRepository, eventPublisher ports.EventPublisher, logger *slog.Logger) *TaskService {
	return &TaskService{
		taskRepo: taskRepo,
		noteRepo: noteRepo,
		eventPublisher: eventPublisher,
		logger: logger,
	}
}

func (s *TaskService) Create(ctx context.Context, userID string, title string, description string, deadline *time.Time) (*domain.Task, error) {
	task := &domain.Task {
		ID: uuid.New().String(),
		Title: title,
		Description: description,
		Status: domain.StatusPending,
		Deadline: deadline,
		UserID: userID,
		CreatedAt: time.Now(),
	}
	
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, err
	}

	if err := s.eventPublisher.Publish(ctx, "task.created", task); err != nil {
    s.logger.Error("Failed to publish event", "error", err)
	}
	return task, nil
}

func (s *TaskService) GetByID(ctx context.Context, userID string, taskID string) (*domain.Task, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task.UserID != userID {
		return nil, domain.ErrTaskForbidden
	}

	return task, nil
}

func (s *TaskService) GetAllByUserID(ctx context.Context, userID string) ([]*domain.Task, error) {
	return s.taskRepo.GetAllByUserID(ctx, userID)
}

func (s *TaskService) Update(ctx context.Context, userID string, taskID string, title string, description string, status domain.TaskStatus, deadline *time.Time) (*domain.Task, error) {
	task, err := s.GetByID(ctx, userID, taskID)
	if err != nil {
		return nil, err
	}

	oldStatus := task.Status

	task.Title = title
	task.Description = description
	task.Status = status
	task.Deadline = deadline
	task.UpdatedAt = time.Now()

	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, err
	}
	
	if oldStatus != status {
    if err := s.eventPublisher.Publish(ctx, "task.status_changed", task); err != nil {
        s.logger.Error("Failed to publish event", "error", err)
    	}
	}
	return task, nil
}

func (s *TaskService) Delete(ctx context.Context, userID string, taskID string) error {
	task, err := s.GetByID(ctx, userID, taskID)
	if err != nil {
		return err
	}

	if err := s.taskRepo.Delete(ctx, taskID); err != nil {
		return err
	}

	if err := s.eventPublisher.Publish(ctx, "task.deleted", task); err != nil {
    	s.logger.Error("Failed to publish event", "error", err)
	}

	return nil
}