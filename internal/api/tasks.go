package api

import (
	"net/http"
	"todo-task-scheduler/internal/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, req *http.Request) {
	search := req.URL.Query().Get("search")

	var tasks []*db.Task
	var err error

	if search != "" {
		tasks, err = db.SearchTasks(search, 10)
	} else {
		tasks, err = db.Tasks(10)
	}

	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()})
		return
	}
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJson(w, TasksResp{Tasks: tasks})
}
