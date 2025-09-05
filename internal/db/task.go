package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

const dateFormat = "20060102"
const taskLimit = 10

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Repeat  string `json:"repeat"`
	Comment string `json:"comment"`
}

func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, repeat, comment) VALUES (?, ?, ?, ?)`
	res, err := DataBase.Exec(query, task.Date, task.Title, task.Repeat, task.Comment)
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err
}

func queryTasks(query string, args ...any) ([]*Task, error) {
	rows, err := DataBase.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var t Task
		var id int64
		if err := rows.Scan(&id, &t.Date, &t.Title, &t.Repeat, &t.Comment); err != nil {
			return nil, err
		}
		t.ID = strconv.FormatInt(id, taskLimit)
		tasks = append(tasks, &t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

func SearchTasks(search string, limit int) ([]*Task, error) { // *5
	if t, err := time.Parse("02.01.2006", search); err == nil {
		date := t.Format(dateFormat)
		query := `SELECT id, date, title, repeat, comment 
                  FROM scheduler 
                  WHERE date = ? 
                  ORDER BY date 
                  LIMIT ?`
		return queryTasks(query, date, limit)
	}

	like := "%" + search + "%"
	query := `SELECT id, date, title, repeat, comment 
              FROM scheduler 
              WHERE title LIKE ? OR comment LIKE ? 
              ORDER BY date 
              LIMIT ?`
	return queryTasks(query, like, like, limit)
}

func Tasks(limit int) ([]*Task, error) {
	query := `SELECT id, date, title, repeat, comment 
              FROM scheduler 
              ORDER BY date 
              LIMIT ?`
	return queryTasks(query, limit)
}

func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, repeat, comment FROM scheduler WHERE id = ?`

	t := &Task{}
	var numId int64
	err := DataBase.QueryRow(query, id).Scan(&numId, &t.Date, &t.Title, &t.Repeat, &t.Comment)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}
	t.ID = strconv.FormatInt(numId, taskLimit)
	return t, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler
	           SET date = ?, title = ?, repeat = ?, comment = ?
	           WHERE id = ?`

	res, err := DataBase.Exec(query, task.Date, task.Title, task.Repeat, task.Comment, task.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := DataBase.Exec(query, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}

func UpdateDate(next string, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := DataBase.Exec(query, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return fmt.Errorf("Task not found")
	}
	return nil
}
