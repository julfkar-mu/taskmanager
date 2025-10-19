package testutils

import (
	"taskmanager/constants"
	"taskmanager/models"
	"time"
)

// CreateTestTask creates a test task with default values
func CreateTestTask() models.Task {
	now := time.Now()
	return models.Task{
		ID:          "test-id-123",
		Title:       "Test Task",
		Description: "Test Description",
		Status:      constants.StatusPending,
		Priority:    constants.PriorityMedium,
		DueDate:     &now,
		CreatedAt:   now,
		UpdatedAt:   now,
		AssignedTo:  "test@example.com",
	}
}

// CreateTestTaskWithStatus creates a test task with specific status
func CreateTestTaskWithStatus(status string) models.Task {
	task := CreateTestTask()
	task.Status = status
	return task
}

// CreateTestTaskWithPriority creates a test task with specific priority
func CreateTestTaskWithPriority(priority string) models.Task {
	task := CreateTestTask()
	task.Priority = priority
	return task
}

// CreateInvalidTask creates a task with invalid data for testing validation
func CreateInvalidTask() models.Task {
	return models.Task{
		ID:          "",
		Title:       "", // Invalid: empty title
		Description: "Test Description",
		Status:      "InvalidStatus", // Invalid: invalid status
		Priority:    "InvalidPriority", // Invalid: invalid priority
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
