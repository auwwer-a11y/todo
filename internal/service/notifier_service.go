package service

import (
	"context"
	"github.com/auwwer-a11y/todo/internal/domain"
	"github.com/auwwer-a11y/todo/internal/ports"
	"log/slog"
	"time"
	"github.com/google/uuid"
	"fmt"
)

type NotifierService struct {
	notificationRepo ports.NotificationRepository
	taskRepo ports.TaskRepository
	logger *slog.Logger
}

func NewNotifierService(notificationRepo ports.NotificationRepository, taskRepo ports.TaskRepository, logger *slog.Logger) *NotifierService {
	return &NotifierService{
		notificationRepo: notificationRepo,
		taskRepo: taskRepo,
		logger: logger,
	}
}

func (s *NotifierService) HandleStatusChanged(ctx context.Context, taskID, userID, newStatus string) error {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		s.logger.Error("Failed to get task", "error", err)
		return err
	}

	if task.UserID != userID {
		s.logger.Warn("User does not own the task", "userID", userID, "taskID", taskID)
		return fmt.Errorf("user does not own the task")
	}

	notification := &domain.Notification{
		ID: uuid.New().String(),
		UserID: userID,
		Message: "Task '" + task.Title + "' status changed to " + newStatus,
		CreatedAt: time.Now(),
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		s.logger.Error("Failed to create notification", "error", err)
	} else {
		s.logger.Info("Notification created", "taskID", taskID, "userID", userID)
	}
	return nil
}

func (s *NotifierService) CheckDeadlines(ctx context.Context) error {
	deadline := time.Now().Add(24 * time.Hour)	
	tasks, err := s.taskRepo.GetTasksWithDeadlineBefore(ctx, deadline)
	if err != nil {
		return err
	}

	for _, task := range tasks {
		notification := &domain.Notification{
			ID: uuid.New().String(),
			UserID: task.UserID,
			Message: "Task '" + task.Title + "' is approaching its deadline",
			CreatedAt: time.Now(),
		}

		if err := s.notificationRepo.Create(ctx, notification); err != nil {
			s.logger.Error("Failed to create notification", "error", err)
		}
	}
	return nil
}

