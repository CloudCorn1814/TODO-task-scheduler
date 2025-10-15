package api

import (
	"net/http"
	"todo-task-scheduler/internal/db"
)

const taskLimit = 10

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.SearchTasks(search, taskLimit)
	} else {
		tasks, err = db.Tasks(taskLimit)
	}

	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, http.StatusOK, TasksResp{Tasks: tasks})
}
