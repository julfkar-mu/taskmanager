package models

import (
	"testing"
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/testutils"
)

func TestTask_IsValidStatus(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Valid Pending", constants.StatusPending, true},
		{"Valid InProgress", constants.StatusInProgress, true},
		{"Valid Completed", constants.StatusCompleted, true},
		{"Valid Cancelled", constants.StatusCancelled, true},
		{"Invalid status", "InvalidStatus", false},
		{"Empty status", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := testutils.CreateTestTask()
			task.Status = tt.status
			result := task.IsValidStatus()
			if result != tt.expected {
				t.Errorf("IsValidStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTask_IsValidPriority(t *testing.T) {
	tests := []struct {
		name     string
		priority string
		expected bool
	}{
		{"Valid Low", constants.PriorityLow, true},
		{"Valid Medium", constants.PriorityMedium, true},
		{"Valid High", constants.PriorityHigh, true},
		{"Empty priority", "", true}, // Priority is optional
		{"Invalid priority", "InvalidPriority", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := testutils.CreateTestTask()
			task.Priority = tt.priority
			result := task.IsValidPriority()
			if result != tt.expected {
				t.Errorf("IsValidPriority() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTask_Validate(t *testing.T) {
	tests := []struct {
		name      string
		task      models.Task
		wantError bool
		errorType string
	}{
		{
			name:      "Valid task",
			task:      testutils.CreateTestTask(),
			wantError: false,
		},
		{
			name: "Empty title",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Title = ""
				return task
			}(),
			wantError: true,
			errorType: "ValidationError",
		},
		{
			name: "Empty status",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Status = ""
				return task
			}(),
			wantError: true,
			errorType: "ValidationError",
		},
		{
			name: "Invalid status",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Status = "InvalidStatus"
				return task
			}(),
			wantError: true,
			errorType: "ValidationError",
		},
		{
			name: "Invalid priority",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Priority = "InvalidPriority"
				return task
			}(),
			wantError: true,
			errorType: "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.task.Validate()
			if tt.wantError {
				if err == nil {
					t.Errorf("Validate() expected error but got none")
					return
				}
				if tt.errorType == "ValidationError" {
					if _, ok := err.(*errors.ValidationError); !ok {
						t.Errorf("Validate() expected ValidationError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error: %v", err)
				}
			}
		})
	}
}
