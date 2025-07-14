package main

import (
	// "github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			MarginLeft(1).
			MarginTop(1).
		// Foreground(lipgloss.Color("#7C3AED")).
		Foreground(lipgloss.Color("#8aadf4")).
		Bold(true)

	helpStyle = lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#626262")).
		Foreground(lipgloss.Color("#C3E88D")).
		MarginTop(2).
		MarginLeft(2)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#C3E88D")).
			MarginTop(2).
			MarginLeft(2)

	// errorStyle = lipgloss.NewStyle().
	// 		Foreground(lipgloss.Color("#FF5F87")).
	// 		MarginTop(1).
	// 		MarginLeft(2)

	inputStyle = lipgloss.NewStyle().
			MarginTop(2).
			MarginLeft(2)

	// Help component styles
	helpKeyStyle = lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#8aadf4")). // Same blue as your title
		Foreground(lipgloss.Color("#c098fe")). // Same blue as your title
		Bold(true)

	helpDescStyle = lipgloss.NewStyle().
		// Foreground(lipgloss.Color("#C3E88D")) // Same green as your help text
		Foreground(lipgloss.Color("#7d66ab")) // Same green as your help text

	helpSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#5C5C5C")) // Subdued separator

	// selected task styles
	selectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#c098fe")).
				Foreground(lipgloss.Color("#c098fe")). // Blue for selected title
				Bold(true).
				Padding(0, 0, 0, 1)

	selectedDescStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color("#c098fe")).
				Foreground(lipgloss.Color("#7d66ab")).
				Padding(0, 0, 0, 1)

	normalTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#DDDDDD")). // White for normal title
				Padding(0, 0, 0, 2)

	normalDescStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")). // Gray for normal description
			Padding(0, 0, 0, 2)
)
