package models

import "time"

type Task struct {
    ID          string    `json:"id"`
    Title       string    `json:"title" binding:"required"`
    Description string    `json:"description,omitempty"`
    Status      string    `json:"status" binding:"required"` // e.g. Pending, Completed
    Priority    string    `json:"priority,omitempty"`         // Optional: Low, Medium, High
    DueDate     time.Time `json:"dueDate,omitempty"`
    CreatedAt   time.Time `json:"createdAt"`
    UpdatedAt   time.Time `json:"updatedAt"`
    AssignedTo  string    `json:"assignedTo,omitempty"`
}
