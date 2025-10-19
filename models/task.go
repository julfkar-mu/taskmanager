package models

import (
	"time"
	"taskmanager/constants"
	"taskmanager/errors"
)

// Task represents a task in the system
type Task struct {
	ID          string    `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Title       string    `json:"title" binding:"required" example:"Complete project documentation"`
	Description string    `json:"description,omitempty" example:"Write comprehensive documentation for the API"`
	Status      string    `json:"status" binding:"required" example:"Pending"`
	Priority    string    `json:"priority,omitempty" example:"High"`
	DueDate     *time.Time `json:"dueDate,omitempty" example:"2024-12-31T23:59:59Z"`
	CreatedAt   time.Time `json:"createdAt" example:"2024-01-01T00:00:00Z"`
	UpdatedAt   time.Time `json:"updatedAt" example:"2024-01-01T00:00:00Z"`
	AssignedTo  string    `json:"assignedTo,omitempty" example:"john.doe@example.com"`
}

// IsValidStatus checks if the status is valid
func (t *Task) IsValidStatus() bool {
	switch t.Status {
	case constants.StatusPending, constants.StatusInProgress, constants.StatusCompleted, constants.StatusCancelled:
		return true
	default:
		return false
	}
}

// IsValidPriority checks if the priority is valid
func (t *Task) IsValidPriority() bool {
	if t.Priority == "" {
		return true // Priority is optional
	}
	switch t.Priority {
	case constants.PriorityLow, constants.PriorityMedium, constants.PriorityHigh:
		return true
	default:
		return false
	}
}

// Validate performs validation on the task
func (t *Task) Validate() error {
	if t.Title == "" {
		return errors.NewValidationError("title", constants.ValidationTitleRequired)
	}
	if t.Status == "" {
		return errors.NewValidationError("status", constants.ValidationStatusRequired)
	}
	if !t.IsValidStatus() {
		return errors.NewValidationError("status", constants.ValidationInvalidStatus)
	}
	if !t.IsValidPriority() {
		return errors.NewValidationError("priority", constants.ValidationInvalidPriority)
	}
	return nil
}
