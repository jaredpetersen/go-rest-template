package task

import (
	"time"

	"github.com/google/uuid"
)

// Task represents something that must be done.
type Task struct {
	ID          string     `json:"id"`
	Description string     `json:"description"`
	DateDue     *time.Time `json:"dateDue,string"`
	DateCreated time.Time  `json:"dateCreated,string"`
	DateUpdated time.Time  `json:"dateUpdated"`
}

// New creates a new task with default values. The returned pointer will never be nil.
func New() *Task {
	now := time.Now()
	return &Task{ID: uuid.New().String(), DateCreated: now, DateUpdated: now}
}
