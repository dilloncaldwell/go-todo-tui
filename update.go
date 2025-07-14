package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

// Filter functions
func (m *Model) refreshListItems() {
	var selectedID int = -1
	if m.list.SelectedItem() != nil {
		if todo, ok := m.list.SelectedItem().(Todo); ok {
			selectedID = todo.ID
		}
	}

	var filtered []list.Item
	switch m.filter {
	case filterAll:
		filtered = append([]list.Item(nil), m.todos...)
	case filterOpen:
		for _, item := range m.todos {
			if todo, ok := item.(Todo); ok && !todo.Done {
				filtered = append(filtered, todo)
			}
		}
	case filterCompleted:
		for _, item := range m.todos {
			if todo, ok := item.(Todo); ok && todo.Done {
				filtered = append(filtered, todo)
			}
		}
	}
	m.filteredTodos = filtered
	m.list.SetItems(filtered)

	if selectedID != -1 {
		foundIndex := -1
		for i, item := range filtered {
			if todo, ok := item.(Todo); ok && todo.ID == selectedID {
				foundIndex = i
				break
			}
		}
		if foundIndex != -1 {
			m.list.Select(foundIndex)
		} else if len(filtered) > 0 {
			m.list.Select(0)
		}
	} else if len(filtered) > 0 {
		m.list.Select(0)
	}
}

func (m *Model) applyFilter() {
	m.refreshListItems()
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

// Sort functions
func (m *Model) applySort() {
	var err error
	m.todos, err = loadTodos(m.db, m.sort)
	if err != nil {
			m.message = fmt.Sprintf("Error loading todos: %v", err)
	}
}

func (m *Model) cycleSort() {
	switch m.sort {
	case sortByID:
		m.sort = sortByCreatedAt
	case sortByCreatedAt:
		m.sort = sortByID
	}
	m.applySort()
	m.refreshListItems()
}

func (m *Model) getSortStatus() string {
	switch m.sort {
	case sortByID:
		return "ID"
	case sortByCreatedAt:
		return "Creation Date"
	default:
		return "ID"
	}
}

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
		// fmt.Printf("DEBUG: Key pressed: %s\n", msg.String())
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
							var err error
							m.todos, err = loadTodos(m.db, m.sort)
							if err != nil {
								m.message = fmt.Sprintf("Error loading todos: %v", err)
							}
							m.refreshListItems()
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

			case key.Matches(msg, keys.sort):
				m.cycleSort()
				m.message = fmt.Sprintf("Sort: %s", m.getSortStatus())
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

			case key.Matches(msg, keys.edit):
				if len(m.filteredTodos) > 0 {
					selectedItem := m.list.SelectedItem()
					if selectedItem != nil {
						if todo, ok := selectedItem.(Todo); ok {
							m.mode = modeEdit
							m.input.SetValue(todo.Text)
							m.input.Focus()
							m.message = fmt.Sprintf("Editing '%s'", todo.Text)
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
							var err error
							m.todos, err = loadTodos(m.db, m.sort)
							if err != nil {
								m.message = fmt.Sprintf("Error loading todos: %v", err)
							}
							m.refreshListItems()
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
						var err error
						m.todos, err = loadTodos(m.db, m.sort)
						if err != nil {
							m.message = fmt.Sprintf("Error loading todos: %v", err)
						}
						m.refreshListItems()
					}
				}
				m.mode = modeList
				return m, nil

			case key.Matches(msg, keys.cancel):
				m.mode = modeList
				m.message = ""
				return m, nil
			}

		case modeEdit:
			switch {
			case key.Matches(msg, keys.confirm):
				text := strings.TrimSpace(m.input.Value())
				if text != "" {
					selectedItem := m.list.SelectedItem()
					if selectedItem != nil {
						if todo, ok := selectedItem.(Todo); ok {
															if err := m.editTodo(todo.ID, text); err != nil {
								m.message = fmt.Sprintf("Error: %v", err)
							} else {
								m.message = fmt.Sprintf("Edited: %s", text)
								var err error
								m.todos, err = loadTodos(m.db, m.sort)
								if err != nil {
									m.message = fmt.Sprintf("Error loading todos: %v", err)
								}
								m.refreshListItems()
							}
						}
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
							var err error
							m.todos, err = loadTodos(m.db, m.sort)
							if err != nil {
								m.message = fmt.Sprintf("Error loading todos: %v", err)
							}
							m.refreshListItems()
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
	case modeEdit:
		m.input, cmd = m.input.Update(msg)
	}

	return m, cmd
}
