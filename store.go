package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
}

type Store struct {
	conn *sql.DB
}

func (s *Store) Init() error {
	var err error
	s.conn, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		return err
	}
	createTableStatement := `CREATE TABLE IF NOT EXISTS tasks (
	id INTEGER NOT NULL PRIMARY KEY,
	title TEXT NOT NULL,
	description TEXT,
	completed BOOLEAN DEFAULT FALSE
	)`

	if _, err = s.conn.Exec(createTableStatement); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetTasks() ([]Task, error) {
	row, err := s.conn.Query("SELECT * FROM tasks WHERE completed=FALSE")
	if err != nil {
		return nil, err
	}
	defer row.Close()

	tasks := []Task{}

	for row.Next() {
		task := Task{}
		row.Scan(&task.ID, &task.Title, &task.Description, &task.Completed)
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (s *Store) SaveTask(title, description string) (Task, error) {
	insertQuery := `
	INSERT INTO tasks(title,description)
	VALUES(?,?)
	`
	if _, err := s.conn.Exec(insertQuery, title, description); err != nil {
		return Task{}, err
	}

	return Task{Title: title, Description: description}, nil
}

func (s *Store) DeleteTask(taskID int) error {
	deleteQuery := `
	DELETE FROM tasks
	WHERE id=?
	`
	if _, err := s.conn.Exec(deleteQuery, taskID); err != nil {
		return err
	}

	return nil
}
