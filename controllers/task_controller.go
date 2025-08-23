package controllers

import (
    "net/http"
    "taskmanager/models"
    "taskmanager/services"

    "github.com/gin-gonic/gin"
)

var taskService services.TaskService

func init() {
    repo := services.NewTaskService(nil) // temporarily nil repo to avoid nil pointer panic
    taskService = repo
}

// Dependency injection setup
func Setup(taskSvc services.TaskService) {
    taskService = taskSvc
}

func GetTasks(c *gin.Context) {
    tasks := taskService.GetTasks()
    c.JSON(http.StatusOK, tasks)
}

func GetTaskByID(c *gin.Context) {
    id := c.Param("id")
    task, err := taskService.GetTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, task)
}

func CreateTask(c *gin.Context) {
    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    created, err := taskService.CreateTask(task)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, created)
}

func UpdateTask(c *gin.Context) {
    id := c.Param("id")
    var task models.Task
    if err := c.ShouldBindJSON(&task); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    updated, err := taskService.UpdateTask(id, task)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, updated)
}

func DeleteTask(c *gin.Context) {
    id := c.Param("id")
    err := taskService.DeleteTask(id)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
