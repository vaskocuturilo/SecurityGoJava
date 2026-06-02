package app

import (
	"log/slog"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	tasks = make(map[string]string)
	mu    sync.RWMutex
)

func Tasks(c *gin.Context) {
	slog.Info("Get all Tasks Endpoint")

	tasks["1"] = "Task1"
	tasks["2"] = "Task2"
	tasks["3"] = "Task3"

	c.JSON(http.StatusOK, tasks)
}

func CreateTask(c *gin.Context) {
	slog.Info("Create a Task Endpoint")

	var newTask = make(map[string]string)

	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}
	mu.Lock()
	for key, value := range newTask {
		tasks[key] = value
	}
	mu.Unlock()
	c.JSON(http.StatusCreated, gin.H{"message": "Task added", "tasks": tasks})
}
