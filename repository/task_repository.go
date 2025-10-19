package repository

import (
	"sync"
	"taskmanager/errors"
	"taskmanager/models"
	"time"
)

var (
	ErrTaskNotFound = errors.NewNotFoundError("Task")
)

type TaskRepository interface {
    GetAll() []models.Task
    GetByID(id string) (models.Task, error)
    Save(task models.Task) models.Task
    Update(id string, task models.Task) (models.Task, error)
    Delete(id string) error
}

type InMemoryTaskRepo struct {
    tasks map[string]models.Task
    mu    sync.RWMutex
}

func NewInMemoryTaskRepo() *InMemoryTaskRepo {
    return &InMemoryTaskRepo{
        tasks: make(map[string]models.Task),
    }
}

func (r *InMemoryTaskRepo) GetAll() []models.Task {
    r.mu.RLock()
    defer r.mu.RUnlock()
    result := make([]models.Task, 0, len(r.tasks))
    for _, task := range r.tasks {
        result = append(result, task)
    }
    return result
}

func (r *InMemoryTaskRepo) GetByID(id string) (models.Task, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    task, ok := r.tasks[id]
    if !ok {
        return models.Task{}, ErrTaskNotFound
    }
    return task, nil
}

func (r *InMemoryTaskRepo) Save(task models.Task) models.Task {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.tasks[task.ID] = task
    return task
}

func (r *InMemoryTaskRepo) Update(id string, task models.Task) (models.Task, error) {
    r.mu.Lock()
    defer r.mu.Unlock()
    _, ok := r.tasks[id]
    if !ok {
        return models.Task{}, ErrTaskNotFound
    }
    task.ID = id
    task.UpdatedAt = time.Now()
    r.tasks[id] = task
    return task, nil
}

func (r *InMemoryTaskRepo) Delete(id string) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    if _, ok := r.tasks[id]; !ok {
        return ErrTaskNotFound
    }
    delete(r.tasks, id)
    return nil
}
