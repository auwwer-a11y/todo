package domain

import (
	"time"
	"errors"
)

type TaskStatus string

const (
	StatusPending TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted TaskStatus = "completed"
)

type Task struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Status TaskStatus `json:"status"`
	Deadline *time.Time `json:"deadline"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

var ErrTaskNotFound = errors.New("task not found")

var ErrTaskForbidden = errors.New("forbidden")