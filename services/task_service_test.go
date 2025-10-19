package services

import (
	"testing"
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/models"
	"taskmanager/repository"
	"taskmanager/testutils"
)

// MockTaskRepository is a mock implementation of TaskRepository for testing
type MockTaskRepository struct {
	tasks map[string]models.Task
}

func NewMockTaskRepository() *MockTaskRepository {
	return &MockTaskRepository{
		tasks: make(map[string]models.Task),
	}
}

func (m *MockTaskRepository) GetAll() []models.Task {
	var result []models.Task
	for _, task := range m.tasks {
		result = append(result, task)
	}
	return result
}

func (m *MockTaskRepository) GetByID(id string) (models.Task, error) {
	task, exists := m.tasks[id]
	if !exists {
		return models.Task{}, repository.ErrTaskNotFound
	}
	return task, nil
}

func (m *MockTaskRepository) Save(task models.Task) models.Task {
	m.tasks[task.ID] = task
	return task
}

func (m *MockTaskRepository) Update(id string, task models.Task) (models.Task, error) {
	if _, exists := m.tasks[id]; !exists {
		return models.Task{}, repository.ErrTaskNotFound
	}
	task.ID = id
	m.tasks[id] = task
	return task, nil
}

func (m *MockTaskRepository) Delete(id string) error {
	if _, exists := m.tasks[id]; !exists {
		return repository.ErrTaskNotFound
	}
	delete(m.tasks, id)
	return nil
}

func TestTaskService_GetTasks(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	service := NewTaskService(mockRepo)

	// Test empty repository
	tasks := service.GetTasks()
	if len(tasks) != 0 {
		t.Errorf("GetTasks() on empty repo = %v, want empty slice", tasks)
	}

	// Test with tasks
	task1 := testutils.CreateTestTask()
	task1.ID = "1"
	task2 := testutils.CreateTestTask()
	task2.ID = "2"

	mockRepo.Save(task1)
	mockRepo.Save(task2)

	tasks = service.GetTasks()
	if len(tasks) != 2 {
		t.Errorf("GetTasks() = %v, want 2 tasks", len(tasks))
	}
}

func TestTaskService_GetTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	service := NewTaskService(mockRepo)
	task := testutils.CreateTestTask()
	task.ID = "test-id"
	mockRepo.Save(task)

	// Test getting existing task
	retrieved, err := service.GetTask("test-id")
	if err != nil {
		t.Errorf("GetTask() unexpected error: %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("GetTask() = %v, want %v", retrieved.ID, task.ID)
	}

	// Test getting non-existent task
	_, err = service.GetTask("non-existent")
	if err != repository.ErrTaskNotFound {
		t.Errorf("GetTask() error = %v, want %v", err, repository.ErrTaskNotFound)
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	service := NewTaskService(mockRepo)

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
			name: "Empty status - should set default",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Status = ""
				return task
			}(),
			wantError: false,
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := service.CreateTask(tt.task)
			if tt.wantError {
				if err == nil {
					t.Errorf("CreateTask() expected error but got none")
					return
				}
				if tt.errorType == "ValidationError" {
					if _, ok := err.(*errors.ValidationError); !ok {
						t.Errorf("CreateTask() expected ValidationError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("CreateTask() unexpected error: %v", err)
					return
				}
				if created.ID == "" {
					t.Errorf("CreateTask() ID should be set")
				}
				if created.CreatedAt.IsZero() {
					t.Errorf("CreateTask() CreatedAt should be set")
				}
				if created.UpdatedAt.IsZero() {
					t.Errorf("CreateTask() UpdatedAt should be set")
				}
				if tt.task.Status == "" && created.Status != constants.StatusPending {
					t.Errorf("CreateTask() default status = %v, want %v", created.Status, constants.StatusPending)
				}
			}
		})
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	service := NewTaskService(mockRepo)
	existingTask := testutils.CreateTestTask()
	existingTask.ID = "test-id"
	mockRepo.Save(existingTask)

	tests := []struct {
		name      string
		id        string
		task      models.Task
		wantError bool
		errorType string
	}{
		{
			name: "Valid update",
			id:   "test-id",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Title = "Updated Title"
				task.Status = constants.StatusCompleted
				return task
			}(),
			wantError: false,
		},
		{
			name: "Non-existent task",
			id:   "non-existent",
			task: testutils.CreateTestTask(),
			wantError: true,
			errorType: "AppError",
		},
		{
			name: "Invalid task data",
			id:   "test-id",
			task: func() models.Task {
				task := testutils.CreateTestTask()
				task.Title = ""
				return task
			}(),
			wantError: true,
			errorType: "ValidationError",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updated, err := service.UpdateTask(tt.id, tt.task)
			if tt.wantError {
				if err == nil {
					t.Errorf("UpdateTask() expected error but got none")
					return
				}
				if tt.errorType == "ValidationError" {
					if _, ok := err.(*errors.ValidationError); !ok {
						t.Errorf("UpdateTask() expected ValidationError, got %T", err)
					}
				} else if tt.errorType == "AppError" {
					if _, ok := err.(*errors.AppError); !ok {
						t.Errorf("UpdateTask() expected AppError, got %T", err)
					}
				}
			} else {
				if err != nil {
					t.Errorf("UpdateTask() unexpected error: %v", err)
					return
				}
				if updated.ID != tt.id {
					t.Errorf("UpdateTask() ID = %v, want %v", updated.ID, tt.id)
				}
				if updated.Title != tt.task.Title {
					t.Errorf("UpdateTask() Title = %v, want %v", updated.Title, tt.task.Title)
				}
			}
		})
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	mockRepo := NewMockTaskRepository()
	service := NewTaskService(mockRepo)
	task := testutils.CreateTestTask()
	task.ID = "test-id"
	mockRepo.Save(task)

	// Test deleting existing task
	err := service.DeleteTask("test-id")
	if err != nil {
		t.Errorf("DeleteTask() unexpected error: %v", err)
	}

	// Test deleting non-existent task
	err = service.DeleteTask("non-existent")
	if err != repository.ErrTaskNotFound {
		t.Errorf("DeleteTask() error = %v, want %v", err, repository.ErrTaskNotFound)
	}
}
