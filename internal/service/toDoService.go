package service

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"test_task/internal/models"
	"test_task/internal/repository"
)

// @Summary Create task
// @Tags Tasks
// @Description Create a new task in the database
// @ID create-task
// @Accept json
// @Produce json
// @Param task body models.TaskRequest true "Task data"
// @Success 201 {object} models.TaskResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tasks [post]
func CreateTask(c *gin.Context) {
	var response models.Response
	var task models.TaskRequest
	var createdTask models.TaskResponse

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("неверные параметры тела запроса"))
		return
	}

	if task.Title == "" {
		c.JSON(http.StatusBadRequest, response.ErrorResponse("поле Title является обязательным"))
		return
	}

	if err := repository.Create(&task, &createdTask); err != nil {
		log.Println("Error creating task:", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("ошибка создания задачи"))
		return
	}

	c.JSON(http.StatusCreated, createdTask)
}

// @Summary Get tasks
// @Tags Tasks
// @Description Retrieve all tasks from the database
// @ID get-tasks
// @Accept json
// @Produce json
// @Success 200 {array} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tasks [get]
func GetTasks(c *gin.Context) {
	var response models.Response
	tasks, err := repository.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("ошибка получения задач"))
		return
	}

	if tasks == nil {
		c.JSON(http.StatusOK, response.NewWithMessage(nil, "список задач пуст"))
	}

	c.JSON(http.StatusOK, response.NewWithMessage(tasks, "ваши задачи"))
}

// @Summary Delete task
// @Tags Tasks
// @Description Delete task
// @ID delete-tasks
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Success 200 {array} models.Response
// @Failure 400 {object} models.Response
// @Router /api/tasks/{id} [delete]
func DeleteTask(c *gin.Context) {
	var response models.Response

	idStr := c.Param("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "ID не может быть пустым"))
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "неверный формат id"))
		return
	}

	err = repository.DeleteTask(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("ошибка при удалении задачи"))
		return
	}

	c.JSON(http.StatusOK, response.NewWithMessage(nil, "задача успешно удалена"))
}

// @Summary Update task
// @Tags Tasks
// @Description Update an existing task
// @ID update-task
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param task body models.TaskRequest true "Updated task data"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tasks/{id} [put]
func UpdateTask(c *gin.Context) {
	var response models.Response
	var taskRequest models.TaskRequest

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "неверный формат id"))
		return
	}

	if err := c.ShouldBindJSON(&taskRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "неверный формат данных"))
		return
	}

	task, err := repository.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewWithMessage(nil, "задача не найдена"))
		return
	}

	if taskRequest.Title != "" {
		task.Title = taskRequest.Title
	}
	if taskRequest.Description != "" {
		task.Description = taskRequest.Description
	}
	if taskRequest.DueDate != "" {
		task.DueDate = taskRequest.DueDate
	}

	err = repository.UpdateTask(task)
	if err != nil {
		log.Println("error updating task:", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("ошибка обновления задачи"))
		return
	}

	c.JSON(http.StatusOK, models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Overdue:     task.Overdue,
		Completed:   task.Completed,
	})
}

// @Summary Mark task as completed
// @Tags Tasks
// @Description Mark a task as completed or not
// @ID complete-task
// @Accept json
// @Produce json
// @Param id path int true "ID"
// @Param task body models.CompleteRequest true "Task completion status"
// @Success 200 {object} models.TaskResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/tasks/{id}/complete [patch]
func CompleteTask(c *gin.Context) {
	var response models.Response
	var completeRequest models.CompleteRequest

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "неверный формат id"))
		return
	}

	if err := c.ShouldBindJSON(&completeRequest); err != nil {
		c.JSON(http.StatusBadRequest, response.NewWithMessage(nil, "неверный формат данных"))
		return
	}

	task, err := repository.GetTaskByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, response.NewWithMessage(nil, "Задача не найдена"))
		return
	}

	task.Completed = completeRequest.Completed

	err = repository.CompleteTask(task)
	if err != nil {
		log.Println("error updating task:", err)
		c.JSON(http.StatusInternalServerError, response.ErrorResponse("ошибка обновления задачи"))
		return
	}

	c.JSON(http.StatusOK, models.TaskResponse{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		DueDate:     task.DueDate,
		Overdue:     task.Overdue,
		Completed:   task.Completed,
	})
}

func UpdateOverdueTasks() error {
	tasks, err := repository.GetAllTasks()
	if err != nil {
		return err
	}

	for _, task := range tasks {
		if task.DueDate != "" && repository.PastedDueDate(task.DueDate) && !task.Overdue {
			task.Overdue = true
			if err := repository.SetOverdue(&task); err != nil {
				log.Printf("Error updating overdue status for task ID %d: %v", task.ID, err)
			}
		}
	}

	return nil
}
