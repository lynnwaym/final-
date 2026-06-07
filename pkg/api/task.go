package api

import (
	"encoding/json"
	"net/http"
	"time"

	"final_project/pkg/db"
)

const TasksLimit = 50

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// обработка других методов будет добавлена на следующих шагах
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPut:
		updTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
	}

}
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "no id specified", http.StatusBadRequest)
		return
	}

	err := db.DeleteTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJson(w, map[string]any{}, http.StatusOK)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "no id specified", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}
	writeJson(w, task, http.StatusOK)
}
func updTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		writeError(w, "ivalid JSON", http.StatusBadRequest)
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
	err = db.UpdateTask(&task)
	if err != nil {
		writeError(w, "database error", http.StatusInternalServerError)
		return
	}
	writeJson(w, map[string]any{}, http.StatusOK)
}

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := db.Tasks(TasksLimit)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJson(w, TasksResp{Tasks: tasks}, http.StatusOK)

}

func doneTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		writeError(w, "no id specified", http.StatusBadRequest)
		return
	}
	task, err := db.GetTask(id)
	if err != nil {
		writeError(w, err.Error(), http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		err := db.DeleteTask(id)
		if err != nil {
			writeError(w, err.Error(), http.StatusNotFound)
			return
		}
		writeJson(w, map[string]any{}, http.StatusOK)
		return
	}

	now := time.Now()
	nextDate, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = db.UpdateDate(nextDate, id)
	if err != nil {
		writeError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeJson(w, map[string]any{}, http.StatusOK)
}
