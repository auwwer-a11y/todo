package service

import (
	"github.com/auwwer-a11y/todo/internal/domain"
	"log/slog"
	"context"
	"github.com/auwwer-a11y/todo/internal/ports"
	"github.com/google/uuid"
	"time"
)

type NoteService struct {
	noteRepo ports.NoteRepository
	taskRepo ports.TaskRepository
	logger *slog.Logger
}

func NewNoteService(noteRepo ports.NoteRepository, taskRepo ports.TaskRepository, logger *slog.Logger) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
		taskRepo: taskRepo,
		logger: logger,
	}
}

func (s *NoteService) Create(ctx context.Context, taskID string, text string, authorID string, meta map[string]interface{}) (*domain.Note, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	note := &domain.Note{
		ID: uuid.New().String(),
		TaskID: taskID,
		Text: text,
		CreatedAt: time.Now(),
		Meta: meta,
		AuthorID: authorID,
	}

	if err := s.noteRepo.Create(ctx, note); err != nil {
		return nil, err
	}

	return note, nil
}

func (s *NoteService) GetAllByTaskID(ctx context.Context, taskID string) ([]*domain.Note, error) {
	task, err := s.taskRepo.GetByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, domain.ErrTaskNotFound
	}

	return s.noteRepo.GetAllByTaskID(ctx, taskID)
}

func (s *NoteService) Delete( ctx context.Context, noteID string) error {
	return s.noteRepo.Delete(ctx, noteID)
}