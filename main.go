package main

import (
	"log"
	"todo-task-scheduler/internal/db"
	"todo-task-scheduler/internal/server"
)

func main() {
	cfg := server.DefaultConfig()

	if err := db.Init(cfg.DBFile); err != nil {
		log.Fatal("db init:", err)
	}
	defer db.DataBase.Close()

	srv := server.NewServer(cfg, db.DataBase)
	log.Printf("server starting on port :%s", cfg.Port)
	log.Fatal(srv.Server.ListenAndServe())
}
