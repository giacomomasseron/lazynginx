package gui

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	// Border colors
	FocusedBorderColor   = lipgloss.Color("#7D56F4") // Purple - for focused panel
	UnfocusedBorderColor = lipgloss.Color("#4A4A4A") // Dim gray - for unfocused panels

	TitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			Bold(true)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Background(lipgloss.Color("#FAFAFA")).
			Padding(0, 1).
			Bold(true)

	ActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#4A4A4A")).
			Padding(0, 1).
			Bold(true)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FAFAFA")).
			Padding(0, 1)

	StatusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5555")).
			Bold(true)

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BD93F9"))

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(30).
			Height(20)

	DetailBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Width(60).
			Height(20)
)
