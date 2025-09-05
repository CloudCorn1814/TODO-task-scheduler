package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"
	"todo-task-scheduler/internal/db"
)

func nextDateHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	query := req.URL.Query()

	nowStr := query.Get("now")
	dateStr := query.Get("date")
	repeat := query.Get("repeat")

	if dateStr == "" || repeat == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "require query parameters 'date' and 'repeat'"})
		return
	}

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		t, err := time.Parse(dateFormat, nowStr)
		if err != nil {
			writeJson(w, http.StatusBadRequest, map[string]string{"error": "incorrect 'now' parameter"})
			return
		}
		now = t
	}

	next, err := NextDate(now, dateStr, repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte(next))
}

func taskHandler(w http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodPost:
		addTaskHandler(w, req)
	case http.MethodGet:
		getTaskHandler(w, req)
	case http.MethodPut:
		updateTaskHandler(w, req)
	case http.MethodDelete:
		deleteTaskHandler(w, req)
	default:
		w.Header().Set("Allow", "GET, POST, PUT, DELETE")
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
	}
}

func getTaskHandler(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(req.URL.Query().Get("id"))
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}
	t, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, t)
}

func updateTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPut {
		w.Header().Set("Allow", http.MethodPut)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	var t db.Task
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&t); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "invalid JSON: " + err.Error()})
		return
	}

	if t.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "title is required"})
		return
	}

	if strings.TrimSpace(t.ID) == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}

	if err := checkDate(&t); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, map[string]string{"result": "ok"})
}

func deleteTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodDelete {
		w.Header().Set("Allow", http.MethodDelete)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id := strings.TrimSpace(req.URL.Query().Get("id"))
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]any{})
}

func doneTaskHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeJson(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}

	id := strings.TrimSpace(req.URL.Query().Get("id"))
	if id == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "ID not specified"})
		return
	}

	t, err := db.GetTask(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if strings.TrimSpace(t.Repeat) == "" {
		if err := db.DeleteTask(id); err != nil {
			writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJson(w, http.StatusOK, map[string]any{})
		return
	}

	next, err := NextDate(time.Now(), t.Date, t.Repeat)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJson(w, http.StatusNotFound, map[string]string{"error": "task not found"})
			return
		}
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJson(w, http.StatusOK, map[string]any{})
}
