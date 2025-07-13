package main

import "github.com/charmbracelet/lipgloss"

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
