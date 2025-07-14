package main

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
)

// Database operations
func initDB(dbPath string) (*sql.DB, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Create table - use datetime() function for consistent format
	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		text TEXT NOT NULL,
		done INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT (datetime('now', 'localtime'))
	);`

	_, err = db.Exec(query)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func loadTodos(db *sql.DB, sort sortMode) ([]list.Item, error) {
	var query string
	switch sort {
	case sortByID:
		query = "SELECT id, text, done, created_at FROM todos ORDER BY id"
	case sortByCreatedAt:
		query = "SELECT id, text, done, created_at FROM todos ORDER BY created_at DESC"
	default:
		query = "SELECT id, text, done, created_at FROM todos ORDER BY id"
	}
	// fmt.Printf("DEBUG: Executing query: %s with sort mode: %v\n", query, sort)
	os.Stderr.Sync()
	rows, err := db.Query(query)


	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []list.Item
	for rows.Next() {
		var todo Todo
		var done int

		err := rows.Scan(&todo.ID, &todo.Text, &done, &todo.CreatedAt)
		if err != nil {
			return nil, err
		}

		todo.Done = done == 1

		todos = append(todos, todo)
	}

	return todos, nil
}

func (m *Model) addTodo(text string) error {
	stmt, err := m.db.Prepare("INSERT INTO todos (text) VALUES (?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(text)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) toggleTodo(id int) error {
	stmt, err := m.db.Prepare("UPDATE todos SET done = 1 - done WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) deleteTodo(id int) error {
	stmt, err := m.db.Prepare("DELETE FROM todos WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Model) editTodo(id int, text string) error {
	stmt, err := m.db.Prepare("UPDATE todos SET text = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(text, id)
	if err != nil {
		return err
	}

	return nil
}
