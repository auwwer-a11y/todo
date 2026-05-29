package domain

import (
	"errors"
	"time"
)

type Note struct {
	ID string `json:"id" bson:"_id"`
	TaskID string `json:"task_id" bson:"task_id"`
	AuthorID string `json:"author_id" bson:"author_id"`
	Text string `json:"text" bson:"text"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`

	Meta map[string]interface{} `json:"meta,omitempty" bson:"meta,omitempty"`
}

var ErrNoteNotFound = errors.New("note not found")