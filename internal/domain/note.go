package domain

import (
	"errors"
	"time"
)

type Note struct {
	ID string `json:"id"`
	TaskID string `json:"task_id"`
	AuthorID string `json:"author_id"`
	Text string `json:"text"`
	CreatedAt time.Time `json:"created_at"`

	Meta map[string]interface{} `json:"meta,omitempty"`
}

var ErrNoteNotFound = errors.New("note not found")