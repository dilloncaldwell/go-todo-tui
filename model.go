package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
)

// Todo represents a single todo item
type Todo struct {
	ID        int
	Text      string
	Done      bool
	CreatedAt time.Time
}

// Implement list.Item interface
func (t Todo) FilterValue() string { return t.Text }
func (t Todo) Title() string {
	status := "[ ]"
	if t.Done {
		status = "[x]"
	}
	return fmt.Sprintf("%s %s", status, t.Text)
}
func (t Todo) Description() string {
	return fmt.Sprintf("ID: %d | Created: %s", t.ID, t.CreatedAt.Local().Format("2006-01-02 15:04"))
}

// App modes
type mode int

const (
	modeList mode = iota
	modeAdd
	modeDelete
	modeEdit
)

// Filter modes
type filterMode int

const (
	filterAll filterMode = iota
	filterOpen
	filterCompleted
)

// Model represents the application state
type Model struct {
	db            *sql.DB
	dbPath        string
	todos         []list.Item
	filteredTodos []list.Item
	list          list.Model
	input         textinput.Model
	mode          mode
	filter        filterMode
	message       string
	quitting      bool
}
