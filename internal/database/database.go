package database

import (
	"database/sql"
	"fmt"
	"log"
	_ "modernc.org/sqlite"
	"os"
	"test_task/internal/config"
)

var database *sql.DB

func InitDB() error {
	cfg, err := config.LoadDBConfig()
	if err != nil {
		return fmt.Errorf("failed to load database config: %w", err)
	}

	dbDir := "./internal/database"
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create database directory: %w", err)
		}
	}

	database, err = sql.Open("sqlite", cfg.DBPath)
	if err != nil {
		return fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		description TEXT,
		due_date TEXT,
		overdue BOOLEAN DEFAULT FALSE,
		completed BOOLEAN DEFAULT FALSE
	);`
	_, err = database.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("failed to create tasks table: %w", err)
	}

	log.Println("Database initialized successfully!")
	return nil
}

func GetDB() *sql.DB {
	return database
}
