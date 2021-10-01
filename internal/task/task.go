package task

import (
	"time"

	"github.com/google/uuid"
)

// Task represents something that must be done.
type Task struct {
	Id          string     `json:"id"`
	Description string     `json:"description"`
	DateDue     *time.Time `json:"date_due,string"`
	DateCreated time.Time  `json:"date_created,string"`
	DateUpdated time.Time  `json:"date_updated"`
}

// New creates a new task with default values. The returned pointer will never be nil.
func New() *Task {
	now := time.Now()
	return &Task{Id: uuid.New().String(), DateCreated: now, DateUpdated: now}
}
