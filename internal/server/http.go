package server

import (
	"database/sql"
	"net/http"

	"todo-task-scheduler/internal/api"
)

type Server struct {
	Server *http.Server
	DB     *sql.DB
}

func NewServer(cfg Config, db *sql.DB) *Server {
	mux := http.NewServeMux()
	api.Init(mux)

	fileService := http.FileServer(http.Dir(cfg.WebDir))
	mux.Handle("/", fileService)

	serverObject := &http.Server{
		Addr:    "localhost:" + cfg.Port,
		Handler: mux,
	}

	return &Server{Server: serverObject, DB: db}
}
