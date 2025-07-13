package main

import (
	"strings"
)

func (m Model) View() string {
	if m.quitting {
		return "Thanks for using todo! 👋"
	}

	var s strings.Builder

	switch m.mode {
	case modeList:
		s.WriteString(m.list.View())
		if m.message != "" {
			s.WriteString("")
			s.WriteString(messageStyle.Render(m.message))
		}
	case modeAdd:
		s.WriteString(titleStyle.Render("📝 Add New Task"))
		s.WriteString("")
		s.WriteString(inputStyle.Render(m.input.View()))
		s.WriteString("")
		s.WriteString(helpStyle.Render("Press Enter to add, Esc to cancel"))

	case modeEdit:
		s.WriteString(titleStyle.Render("✍️ Edit Task"))
		s.WriteString("")
		s.WriteString(inputStyle.Render(m.input.View()))
		s.WriteString("")
		s.WriteString(helpStyle.Render("Press Enter to save, Esc to cancel"))

	case modeDelete:
		s.WriteString(titleStyle.Render("🗑️  Delete Task"))
		s.WriteString("")
		s.WriteString(inputStyle.Render(m.message))
	}

	return s.String()
}
