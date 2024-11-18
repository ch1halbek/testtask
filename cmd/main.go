package main

import (
	"log"
	"test_task/cmd/app"
	_ "test_task/docs"
	"test_task/internal/database"
)

func main() {
	err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	log.Println("Starting server on http://localhost:8080")
	app.StartServe()
}
