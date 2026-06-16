package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888")).
			Padding(0, 1)

	errStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Padding(0, 1)

	metaStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#66CCFF")).
			Padding(0, 1)
)
