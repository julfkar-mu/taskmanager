package services

import (
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/models"
	"taskmanager/repository"
	"time"

	"github.com/google/uuid"
)

type TaskService interface {
    GetTasks() []models.Task
    GetTask(id string) (models.Task, error)
    CreateTask(task models.Task) (models.Task, error)
    UpdateTask(id string, task models.Task) (models.Task, error)
    DeleteTask(id string) error
}

type taskService struct {
    repo repository.TaskRepository
}

func NewTaskService(r repository.TaskRepository) TaskService {
    return &taskService{repo: r}
}

func (s *taskService) GetTasks() []models.Task {
    return s.repo.GetAll()
}

func (s *taskService) GetTask(id string) (models.Task, error) {
    return s.repo.GetByID(id)
}

func (s *taskService) CreateTask(task models.Task) (models.Task, error) {
	// Validate the task
	if err := task.Validate(); err != nil {
		return models.Task{}, err
	}

	// Set default values
	task.ID = uuid.NewString()
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	// Set default status if not provided
	if task.Status == "" {
		task.Status = constants.StatusPending
	}

	return s.repo.Save(task), nil
}

func (s *taskService) UpdateTask(id string, task models.Task) (models.Task, error) {
	// Get existing task
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return models.Task{}, err
	}

	// Validate the updated task
	if err := task.Validate(); err != nil {
		return models.Task{}, err
	}

	// Only update allowed fields (SOLID - Single Responsibility)
	existing.Title = task.Title
	existing.Description = task.Description
	existing.Status = task.Status
	existing.Priority = task.Priority
	existing.DueDate = task.DueDate
	existing.AssignedTo = task.AssignedTo
	existing.UpdatedAt = time.Now()

	return s.repo.Update(id, existing)
}

func (s *taskService) DeleteTask(id string) error {
    return s.repo.Delete(id)
}
