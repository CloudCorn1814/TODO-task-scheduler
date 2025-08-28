package server

import "os"

type Config struct {
	Port   string
	WebDir string
	DBFile string
}

func DefaultConfig() Config {
	port := os.Getenv("TODO_PORT") //*1
	if port == "" {
		port = "7540"
	}

	dbFile := os.Getenv("TODO_DBFILE") //*2
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	return Config{
		Port:   port,
		WebDir: "./web",
		DBFile: dbFile,
	}
}
