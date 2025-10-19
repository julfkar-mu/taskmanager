package repository

import (
	"testing"
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/testutils"
)

func TestInMemoryTaskRepo_GetAll(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	
	// Test empty repository
	tasks := repo.GetAll()
	if len(tasks) != 0 {
		t.Errorf("GetAll() on empty repo = %v, want empty slice", tasks)
	}

	// Test with tasks
	task1 := testutils.CreateTestTask()
	task1.ID = "1"
	task2 := testutils.CreateTestTask()
	task2.ID = "2"
	task2.Status = constants.StatusCompleted

	repo.Save(task1)
	repo.Save(task2)

	tasks = repo.GetAll()
	if len(tasks) != 2 {
		t.Errorf("GetAll() = %v, want 2 tasks", len(tasks))
	}
}

func TestInMemoryTaskRepo_GetByID(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	task := testutils.CreateTestTask()
	task.ID = "test-id"

	// Test getting non-existent task
	_, err := repo.GetByID("non-existent")
	if err != ErrTaskNotFound {
		t.Errorf("GetByID() error = %v, want %v", err, ErrTaskNotFound)
	}

	// Test getting existing task
	repo.Save(task)
	retrieved, err := repo.GetByID("test-id")
	if err != nil {
		t.Errorf("GetByID() unexpected error: %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("GetByID() = %v, want %v", retrieved.ID, task.ID)
	}
}

func TestInMemoryTaskRepo_Save(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	task := testutils.CreateTestTask()
	task.ID = "test-id"

	// Test saving task
	saved := repo.Save(task)
	if saved.ID != task.ID {
		t.Errorf("Save() = %v, want %v", saved.ID, task.ID)
	}

	// Verify task was saved
	retrieved, err := repo.GetByID("test-id")
	if err != nil {
		t.Errorf("GetByID() after save unexpected error: %v", err)
	}
	if retrieved.ID != task.ID {
		t.Errorf("GetByID() after save = %v, want %v", retrieved.ID, task.ID)
	}
}

func TestInMemoryTaskRepo_Update(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	task := testutils.CreateTestTask()
	task.ID = "test-id"
	repo.Save(task)

	// Test updating existing task
	updatedTask := testutils.CreateTestTask()
	updatedTask.ID = "test-id"
	updatedTask.Title = "Updated Title"
	updatedTask.Status = constants.StatusCompleted

	updated, err := repo.Update("test-id", updatedTask)
	if err != nil {
		t.Errorf("Update() unexpected error: %v", err)
	}
	if updated.Title != "Updated Title" {
		t.Errorf("Update() title = %v, want %v", updated.Title, "Updated Title")
	}
	if updated.Status != constants.StatusCompleted {
		t.Errorf("Update() status = %v, want %v", updated.Status, constants.StatusCompleted)
	}

	// Test updating non-existent task
	_, err = repo.Update("non-existent", updatedTask)
	if err != ErrTaskNotFound {
		t.Errorf("Update() error = %v, want %v", err, ErrTaskNotFound)
	}
}

func TestInMemoryTaskRepo_Delete(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	task := testutils.CreateTestTask()
	task.ID = "test-id"
	repo.Save(task)

	// Test deleting existing task
	err := repo.Delete("test-id")
	if err != nil {
		t.Errorf("Delete() unexpected error: %v", err)
	}

	// Verify task was deleted
	_, err = repo.GetByID("test-id")
	if err != ErrTaskNotFound {
		t.Errorf("GetByID() after delete = %v, want %v", err, ErrTaskNotFound)
	}

	// Test deleting non-existent task
	err = repo.Delete("non-existent")
	if err != ErrTaskNotFound {
		t.Errorf("Delete() error = %v, want %v", err, ErrTaskNotFound)
	}
}

func TestInMemoryTaskRepo_Concurrency(t *testing.T) {
	repo := NewInMemoryTaskRepo()
	
	// Test concurrent writes
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			task := testutils.CreateTestTask()
			task.ID = string(rune('0' + i))
			repo.Save(task)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all tasks were saved
	tasks := repo.GetAll()
	if len(tasks) != 10 {
		t.Errorf("Concurrent saves resulted in %v tasks, want 10", len(tasks))
	}
}
