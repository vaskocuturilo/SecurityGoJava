package app

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Tasks(c *gin.Context) {
	slog.Info("Tasks Endpoint")

	tasks := map[string]string{
		"1": "Task1",
		"2": "Task2",
		"3": "Task3",
	}

	c.JSON(http.StatusOK, tasks)
}
