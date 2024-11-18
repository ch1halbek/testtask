package repository

import (
	"fmt"
	"log"
	"test_task/internal/database"
	"test_task/internal/models"
	"time"
)

func Create(task *models.TaskRequest, cTask *models.TaskResponse) error {
	db := database.GetDB()

	query := `INSERT INTO tasks (title, description, due_date) 
			  VALUES (?, ?, ?)`

	result, err := db.Exec(query, task.Title, task.Description, task.DueDate)
	if err != nil {
		return fmt.Errorf("error inserting task into database: %w", err)
	}

	taskID, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("error retrieving task ID: %w", err)
	}

	cTask.ID = int(taskID)
	cTask.Title = task.Title
	cTask.Description = task.Description
	cTask.DueDate = task.DueDate
	cTask.Overdue = false
	cTask.Completed = false

	return nil
}

func GetAllTasks() ([]models.TaskResponse, error) {
	db := database.GetDB()

	rows, err := db.Query("SELECT id, title, description, due_date, overdue, completed FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("error fetching tasks: %w", err)
	}
	defer rows.Close()

	var tasks []models.TaskResponse

	for rows.Next() {
		var task models.TaskResponse
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Overdue, &task.Completed); err != nil {
			return nil, fmt.Errorf("error scanning task row: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return tasks, nil
}

func DeleteTask(id int) error {
	db := database.GetDB()

	if !taskExists(id) {
		return fmt.Errorf("задача с id %s не найдена", id)
	}

	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %w", err)
	}

	return nil
}

func taskExists(id int) bool {
	db := database.GetDB()

	query := `SELECT 1 FROM tasks WHERE id = ? LIMIT 1`

	var exists int
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		log.Println("Ошибка проверки:", err)
		return false
	}

	return exists == 1
}

func GetTaskByID(id int) (*models.TaskResponse, error) {
	db := database.GetDB()

	var task models.TaskResponse
	query := `SELECT id, title, description, due_date, overdue, completed FROM tasks WHERE id = ?`

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Description, &task.DueDate, &task.Overdue, &task.Completed)
	if err != nil {
		return nil, fmt.Errorf("error fetching task by ID: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *models.TaskResponse) error {
	db := database.GetDB()

	query := `UPDATE tasks SET title = ?, description = ?, due_date = ? WHERE id = ?`
	_, err := db.Exec(query, task.Title, task.Description, task.DueDate, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	return nil
}

func CompleteTask(task *models.TaskResponse) error {
	db := database.GetDB()

	query := `UPDATE tasks SET completed = ? WHERE id = ?`

	result, err := db.Exec(query, task.Completed, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", task.ID)
	}

	return nil
}

func PastedDueDate(dueDate string) bool {
	layout := "2006-01-02"
	parsedDate, err := time.Parse(layout, dueDate)
	if err != nil {
		log.Printf("Invalid due date format: %v", err)
		return false
	}

	now := time.Now()
	log.Printf("Checking dueDate: %s, parsedDate: %s, current time: %s", dueDate, parsedDate.Format(layout), now.Format(layout))

	return now.After(parsedDate)
}

func SetOverdue(task *models.TaskResponse) error {
	db := database.GetDB()

	query := `UPDATE tasks SET overdue = ? WHERE id = ?`
	_, err := db.Exec(query, task.Overdue, task.ID)
	if err != nil {
		return fmt.Errorf("error updating task in database: %w", err)
	}

	return nil
}
