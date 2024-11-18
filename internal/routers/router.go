package routers

import (
	"github.com/gin-gonic/gin"
	"test_task/internal/service"
)

func ToDoList(router *gin.RouterGroup) {
	router.POST("/tasks", service.CreateTask) // Создание задачи
	router.GET("/tasks", service.GetTasks)
	router.DELETE("/tasks/:id", service.DeleteTask)
	router.PUT("/tasks/:id", service.UpdateTask)
	router.PATCH("/tasks/:id/complete", service.CompleteTask)
}
