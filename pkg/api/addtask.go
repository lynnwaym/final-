package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"final_project/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeError(w, "title is required", http.StatusBadRequest)
		return
	}

	nowDate := time.Now()
	now := time.Date(nowDate.Year(), nowDate.Month(), nowDate.Day(), 0, 0, 0, 0, time.UTC)

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)
	if err != nil {
		writeError(w, "incorrect date format", http.StatusBadRequest)
		return
	}
	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}
	id, err := db.AddTask(&task)
	if err != nil {
		writeError(w, "database error", http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]any{"id": id}, http.StatusCreated)
}

func writeJson(w http.ResponseWriter, data any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

func writeError(w http.ResponseWriter, errMsg string, status int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": errMsg})
}
