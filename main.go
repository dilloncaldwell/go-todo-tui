package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	_ "github.com/mattn/go-sqlite3"
)

// Initialize the model
func initialModel(dbPath string) (Model, error) {
	db, err := initDB(dbPath)
	if err != nil {
		return Model{}, err
	}

	todos, err := loadTodos(db)
	if err != nil {
		return Model{}, err
	}

	// Create list
	l := list.New(todos, list.NewDefaultDelegate(), 0, 0)
	l.Title = "📋 Todo List"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleStyle
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.add,
			keys.delete,
			keys.edit,
			keys.toggle,
			keys.filter,
		}
	}
	l.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{
			keys.add,
			keys.delete,
			keys.edit,
			keys.toggle,
			keys.filter,
		}
	}

	// Create text input
	ti := textinput.New()
	ti.Placeholder = "Enter task..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50

	return Model{
		db:            db,
		dbPath:        dbPath,
		todos:         todos,
		filteredTodos: todos,
		list:          l,
		input:         ti,
		mode:          modeList,
		filter:        filterAll,
	}, nil
}

// Get database path
func getDBPath(global bool, configPath string) string {
	if configPath != "" {
		return configPath
	}

	if global {
		dataHome := os.Getenv("XDG_DATA_HOME")
		if dataHome == "" {
			home := os.Getenv("HOME")
			if home == "" {
				log.Fatal("Could not determine home directory")
			}
			dataHome = filepath.Join(home, ".local", "share")
		}
		return filepath.Join(dataHome, "todos.db")
	}

	return ".todos.db"
}

func main() {
	var (
		global     = flag.Bool("g", false, "Use global todo list")
		configPath = flag.String("config", "", "Custom database path")
	)
	flag.Parse()

	dbPath := getDBPath(*global, *configPath)

	model, err := initialModel(dbPath)
	if err != nil {
		log.Fatal(err)
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}