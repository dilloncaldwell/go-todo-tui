package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	_ "github.com/mattn/go-sqlite3"
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
	return fmt.Sprintf("ID: %d | Created: %s", t.ID, t.CreatedAt.Format("2006-01-02 15:04"))
}

// App modes
type mode int

const (
	modeList mode = iota
	modeAdd
	modeDelete
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

// Key bindings
type keyMap struct {
	up      key.Binding
	down    key.Binding
	add     key.Binding
	delete  key.Binding
	toggle  key.Binding
	filter  key.Binding
	confirm key.Binding
	cancel  key.Binding
	quit    key.Binding
	help    key.Binding
}

var keys = keyMap{
	up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add task"),
	),
	delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	toggle: key.NewBinding(
		key.WithKeys("t", "space"),
		key.WithHelp("t/space", "toggle done"),
	),
	filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter views"),
	),
	confirm: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "confirm"),
	),
	cancel: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "cancel"),
	),
	quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "help"),
	),
}

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginTop(1).
			Foreground(lipgloss.Color("#7C3AED")).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			MarginTop(1).
			MarginLeft(2)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575")).
			MarginTop(1).
			MarginLeft(2)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F87")).
			MarginTop(1).
			MarginLeft(2)

	inputStyle = lipgloss.NewStyle().
			MarginTop(1).
			MarginLeft(2)
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

func loadTodos(db *sql.DB) ([]list.Item, error) {
	rows, err := db.Query("SELECT id, text, done, created_at FROM todos ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []list.Item
	for rows.Next() {
		var todo Todo
		var done int
		var createdAt string

		err := rows.Scan(&todo.ID, &todo.Text, &done, &createdAt)
		if err != nil {
			return nil, err
		}

		todo.Done = done == 1

		// Try multiple timestamp formats that SQLite might use
		if parsedTime, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			todo.CreatedAt = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05Z", createdAt); err == nil {
			todo.CreatedAt = parsedTime
		} else if parsedTime, err := time.Parse("2006-01-02T15:04:05.000Z", createdAt); err == nil {
			todo.CreatedAt = parsedTime
		} else {
			// If all parsing fails, use current time as fallback
			todo.CreatedAt = time.Now()
		}

		todos = append(todos, todo)
	}

	return todos, nil
}

// Filter functions
func (m *Model) applyFilter() {
	switch m.filter {
	case filterAll:
		m.filteredTodos = m.todos
		m.list.Title = "📋 Todo List - All"
	case filterOpen:
		m.filteredTodos = []list.Item{}
		for _, item := range m.todos {
			if todo, ok := item.(Todo); ok && !todo.Done {
				m.filteredTodos = append(m.filteredTodos, item)
			}
		}
		m.list.Title = "📋 Todo List - Open"
	case filterCompleted:
		m.filteredTodos = []list.Item{}
		for _, item := range m.todos {
			if todo, ok := item.(Todo); ok && todo.Done {
				m.filteredTodos = append(m.filteredTodos, item)
			}
		}
		m.list.Title = "📋 Todo List - Completed"
	}
	m.list.SetItems(m.filteredTodos)
}

func (m *Model) cycleFilter() {
	switch m.filter {
	case filterAll:
		m.filter = filterOpen
	case filterOpen:
		m.filter = filterCompleted
	case filterCompleted:
		m.filter = filterAll
	}
	m.applyFilter()
}

func (m *Model) getFilterStatus() string {
	switch m.filter {
	case filterAll:
		return "All"
	case filterOpen:
		return "Open"
	case filterCompleted:
		return "Completed"
	default:
		return "All"
	}
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

	// Reload todos
	m.todos, err = loadTodos(m.db)
	if err != nil {
		return err
	}

	m.applyFilter()
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

	// Reload todos
	m.todos, err = loadTodos(m.db)
	if err != nil {
		return err
	}

	m.applyFilter()
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

	// Reload todos
	m.todos, err = loadTodos(m.db)
	if err != nil {
		return err
	}

	m.applyFilter()
	return nil
}

// Bubble Tea methods
func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height - 4)
		return m, nil

	case tea.KeyMsg:
		switch m.mode {
		case modeList:
			// Handle space key first, before passing to list
			if msg.String() == " " {
				if len(m.filteredTodos) > 0 {
					selectedItem := m.list.SelectedItem()
					if selectedItem != nil {
						if todo, ok := selectedItem.(Todo); ok {
							if err := m.toggleTodo(todo.ID); err != nil {
								m.message = fmt.Sprintf("Error: %v", err)
							} else {
								status := "done"
								if todo.Done {
									status = "undone"
								}
								m.message = fmt.Sprintf("Marked '%s' as %s", todo.Text, status)
							}
						}
					}
				}
				return m, nil
			}

			switch {
			case key.Matches(msg, keys.quit):
				m.quitting = true
				return m, tea.Quit

			case key.Matches(msg, keys.add):
				m.mode = modeAdd
				m.input.Reset()
				m.input.Focus()
				m.message = ""
				return m, nil

			case key.Matches(msg, keys.filter):
				m.cycleFilter()
				m.message = fmt.Sprintf("Filter: %s", m.getFilterStatus())
				return m, nil

			case key.Matches(msg, keys.delete):
				if len(m.filteredTodos) > 0 {
					selectedItem := m.list.SelectedItem()
					if selectedItem != nil {
						if todo, ok := selectedItem.(Todo); ok {
							m.mode = modeDelete
							m.message = fmt.Sprintf("Delete '%s'? (y/N)", todo.Text)
							return m, nil
						}
					}
				}

			case key.Matches(msg, keys.toggle):
				if len(m.filteredTodos) > 0 {
					selectedItem := m.list.SelectedItem()
					if selectedItem != nil {
						if todo, ok := selectedItem.(Todo); ok {
							if err := m.toggleTodo(todo.ID); err != nil {
								m.message = fmt.Sprintf("Error: %v", err)
							} else {
								status := "done"
								if todo.Done {
									status = "undone"
								}
								m.message = fmt.Sprintf("Marked '%s' as %s", todo.Text, status)
							}
						}
					}
				}
				return m, nil
			}

		case modeAdd:
			switch {
			case key.Matches(msg, keys.confirm):
				text := strings.TrimSpace(m.input.Value())
				if text != "" {
					if err := m.addTodo(text); err != nil {
						m.message = fmt.Sprintf("Error: %v", err)
					} else {
						m.message = fmt.Sprintf("Added: %s", text)
					}
				}
				m.mode = modeList
				return m, nil

			case key.Matches(msg, keys.cancel):
				m.mode = modeList
				m.message = ""
				return m, nil
			}

		case modeDelete:
			switch msg.String() {
			case "y", "Y":
				selectedItem := m.list.SelectedItem()
				if selectedItem != nil {
					if todo, ok := selectedItem.(Todo); ok {
						if err := m.deleteTodo(todo.ID); err != nil {
							m.message = fmt.Sprintf("Error: %v", err)
						} else {
							m.message = fmt.Sprintf("Deleted: %s", todo.Text)
						}
					}
				}
				m.mode = modeList
				return m, nil

			case "n", "N", "esc":
				m.mode = modeList
				m.message = ""
				return m, nil
			}
		}
	}

	// Update components
	var cmd tea.Cmd
	switch m.mode {
	case modeList:
		m.list, cmd = m.list.Update(msg)
	case modeAdd:
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	if m.quitting {
		return "Thanks for using todo! 👋\n"
	}

	var s strings.Builder

	switch m.mode {
	case modeList:
		s.WriteString(m.list.View())

		// Help text
		help := "\n" + helpStyle.Render("Press 'a' to add, 'd' to delete, 't/space' to toggle, 'f' to filter, 'q' to quit")
		s.WriteString(help)

	case modeAdd:
		s.WriteString(titleStyle.Render("📝 Add New Task"))
		s.WriteString("\n\n")
		s.WriteString(inputStyle.Render(m.input.View()))
		s.WriteString("\n\n")
		s.WriteString(helpStyle.Render("Press Enter to add, Esc to cancel"))

	case modeDelete:
		s.WriteString(titleStyle.Render("🗑️  Delete Task"))
		s.WriteString("\n\n")
		s.WriteString(inputStyle.Render(m.message))
	}

	// Show message
	if m.message != "" && m.mode == modeList {
		s.WriteString("\n")
		s.WriteString(messageStyle.Render(m.message))
	}

	return s.String()
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
