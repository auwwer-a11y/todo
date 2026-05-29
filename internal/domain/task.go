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
	ID string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Title string `json:"title" db:"title"`
	Description string `json:"description" db:"description"`
	Status TaskStatus `json:"status" db:"status"`
	Deadline *time.Time `json:"deadline" db:"deadline"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

var ErrTaskNotFound = errors.New("task not found")

var ErrTaskForbidden = errors.New("forbidden")