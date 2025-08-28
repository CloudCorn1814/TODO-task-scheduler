package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
id INTEGER PRIMARY KEY AUTOINCREMENT,
date CHAR(8) NOT NULL DEFAULT "",
title VARCHAR(128) NOT NULL DEFAULT "",
comment TEXT NOT NULL DEFAULT "",
repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler(date);
`

var DataBase *sql.DB

func Init(dbFile string) error {
	if DataBase != nil {
		return fmt.Errorf("db already initialized")
	}

	var install bool

	if _, err := os.Stat(dbFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			install = true
		} else {
			return err
		}
	}

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("open %s: %w", dbFile, err)
	}

	if err = db.Ping(); err != nil {
		_ = db.Close()
		return fmt.Errorf("ping sqlite: %w", err)
	}

	if install {
		if _, err = db.Exec(schema); err != nil {
			_ = db.Close()
			return fmt.Errorf("apply schema: %w", err)
		}
	}
	DataBase = db
	return nil
}
