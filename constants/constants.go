package constants

// Task status constants
const (
	StatusPending    = "Pending"
	StatusInProgress = "InProgress"
	StatusCompleted  = "Completed"
	StatusCancelled  = "Cancelled"
)

// Task priority constants
const (
	PriorityLow    = "Low"
	PriorityMedium = "Medium"
	PriorityHigh   = "High"
)

// HTTP status messages
const (
	MessageTaskCreated    = "Task created successfully"
	MessageTaskUpdated    = "Task updated successfully"
	MessageTaskDeleted    = "Task deleted successfully"
	MessageTaskNotFound   = "Task not found"
	MessageInvalidInput   = "Invalid input"
	MessageInternalError  = "Internal server error"
)

// Validation messages
const (
	ValidationTitleRequired  = "title is required"
	ValidationStatusRequired = "status is required"
	ValidationInvalidStatus  = "invalid status value"
	ValidationInvalidPriority = "invalid priority value"
)
