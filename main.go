package main

import (
    "taskmanager/controllers"
    "github.com/gin-gonic/gin"
)
func main() {
    repo := repository.NewInMemoryTaskRepo()
    service := services.NewTaskService(repo)
    controllers.Setup(service)

    router := gin.Default()

    router.GET("/tasks", controllers.GetTasks)
    router.POST("/tasks", controllers.CreateTask)
    router.GET("/tasks/:id", controllers.GetTaskByID)
    router.PUT("/tasks/:id", controllers.UpdateTask)
    router.DELETE("/tasks/:id", controllers.DeleteTask)

    router.Run(":8080")
}
