package api

import (
	"encoding/json"
	"net/http"
	"time"
	"todo-task-scheduler/internal/db"
)

func addTaskHandler(w http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()

	var task db.Task
	if err := json.NewDecoder(req.Body).Decode(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusCreated, struct {
		ID int64 `json:"id"`
	}{ID: id})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	t, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return err
	}

	var next string
	if len(task.Repeat) != 0 {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return err
		}
	}

	if now.Format(dateFormat) > t.Format(dateFormat) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(dateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
