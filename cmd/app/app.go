package app

import (
	"context"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"test_task/internal/routers"
	"test_task/internal/service"
	"time"
)

func StartBackgroundTask(stopChan chan struct{}, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(1 * time.Minute) // Проверяем раз в минуту
		defer ticker.Stop()

		for {
			select {
			case <-stopChan:
				log.Println("Stopping background task...")
				return
			case <-ticker.C:
				if err := service.UpdateOverdueTasks(); err != nil {
					log.Printf("Error updating overdue tasks: %v", err)
				}
			}
		}
	}()
}

func StartServe() {
	router := setupRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Порт по умолчанию
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	stopChan := make(chan struct{})
	var wg sync.WaitGroup

	StartBackgroundTask(stopChan, &wg)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server is running on port :%s", port)

	gracefulShutdown(srv, stopChan, &wg, 5*time.Second)
}

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the API"})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	routers.ToDoList(api)

	return router
}

func gracefulShutdown(srv *http.Server, stopChan chan struct{}, wg *sync.WaitGroup, timeout time.Duration) {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Println("Shutting down server...")

	close(stopChan)
	wg.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v\n", err)
		os.Exit(1)
	} else {
		log.Println("Server shutdown gracefully")
		os.Exit(0)
	}
}
