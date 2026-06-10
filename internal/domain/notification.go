package domain

import "time"

type Notification struct {
	ID string `json:"id"`
	TaskID string `json:"task_id"`
	UserID string `json:"user_id"`
	Type string `json:"type"`
	Message string `json:"message"`
	CreatedAt time.Time `json:"created_at"`
}