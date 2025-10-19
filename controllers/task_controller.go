package controllers

import (
	"net/http"
	"taskmanager/constants"
	"taskmanager/errors"
	"taskmanager/models"
	"taskmanager/services"

	"github.com/gin-gonic/gin"
)

var taskService services.TaskService

func init() {
	// Initialize with nil to avoid nil pointer panic
	// This will be properly set via Setup() function
	taskService = nil
}

// Dependency injection setup
func Setup(taskSvc services.TaskService) {
	taskService = taskSvc
}

// GetTasks retrieves all tasks
// @Summary Get all tasks
// @Description Get a list of all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Task
// @Router /tasks [get]
func GetTasks(c *gin.Context) {
	tasks := taskService.GetTasks()
	c.JSON(http.StatusOK, gin.H{
		"data": tasks,
		"count": len(tasks),
	})
}

// GetTaskByID retrieves a task by ID
// @Summary Get task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} models.Task
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [get]
func GetTaskByID(c *gin.Context) {
	id := c.Param("id")
	task, err := taskService.GetTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": task})
}

// CreateTask creates a new task
// @Summary Create a new task
// @Description Create a new task with the provided information
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task information"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Router /tasks [post]
func CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	created, err := taskService.CreateTask(task)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"data": created,
		"message": constants.MessageTaskCreated,
	})
}

// UpdateTask updates an existing task
// @Summary Update a task
// @Description Update an existing task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param task body models.Task true "Updated task information"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	id := c.Param("id")
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	updated, err := taskService.UpdateTask(id, task)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": updated,
		"message": constants.MessageTaskUpdated,
	})
}

// DeleteTask deletes a task
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	err := taskService.DeleteTask(id)
	if err != nil {
		handleError(c, err)
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": constants.MessageTaskDeleted})
}

// handleError handles different types of errors and returns appropriate HTTP responses
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *errors.ValidationError:
		c.JSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
			"field": e.Field,
		})
	case *errors.AppError:
		c.JSON(e.Code, gin.H{"error": e.Message})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": constants.MessageInternalError,
		})
	}
}
