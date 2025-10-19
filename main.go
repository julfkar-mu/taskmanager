package main

import (
	"log"
	"taskmanager/controllers"
	"taskmanager/repository"
	"taskmanager/services"

	"github.com/gin-gonic/gin"
)

func main() {
	repo := repository.NewInMemoryTaskRepo()
	service := services.NewTaskService(repo)
	controllers.Setup(service)

	router := gin.Default()

	// API routes
	api := router.Group("/api/v1")
	{
		api.GET("/tasks", controllers.GetTasks)
		api.POST("/tasks", controllers.CreateTask)
		api.GET("/tasks/:id", controllers.GetTaskByID)
		api.PUT("/tasks/:id", controllers.UpdateTask)
		api.DELETE("/tasks/:id", controllers.DeleteTask)
	}

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
