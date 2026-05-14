package app

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func Tasks(w http.ResponseWriter, r *http.Request) {
	slog.Info("Tasks Endpoint")

	tasks := map[string]string{
		"1": "Task1",
		"2": "Task2",
		"3": "Task3",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tasks)
}
