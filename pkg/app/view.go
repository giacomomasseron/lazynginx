package app

import (
	"lazynginx/pkg/gui"

	"github.com/charmbracelet/lipgloss"
	"github.com/jesseduffield/lazycore/pkg/boxlayout"
)

func (m Model) View() string {
	// Ensure we have valid dimensions
	if m.WindowWidth < 40 || m.WindowHeight < 10 {
		return "Terminal too small. Please resize."
	}

	// Calculate available space
	// Reserve 1 line for footer
	footerHeight := 1
	// Reserve an extra line to prevent the top border from being cut off
	contentHeight := m.WindowHeight - footerHeight - 1

	// Create horizontal layout with 3 boxes: weights 1, 1, 2 (25%, 25%, 50%)
	root := &boxlayout.Box{
		Direction: boxlayout.COLUMN,
		Children: []*boxlayout.Box{
			{Window: "mainmenu", Weight: 1},
			{Window: "submenu", Weight: 1},
			{Window: "details", Weight: 2},
		},
	}

	// Arrange windows to get dimensions
	dimensions := boxlayout.ArrangeWindows(root, 0, 0, m.WindowWidth, contentHeight)

	// Render each panel with its dimensions
	mainMenuView := gui.ViewMainMenuWithDim(m, dimensions["mainmenu"])
	subMenuView := gui.ViewSubMenuWithDim(m, dimensions["submenu"])
	detailsView := gui.ViewDetailsWithDim(m, dimensions["details"])

	// Create footer with keybindings
	footer := gui.ViewFooter(m, m.WindowWidth)

	// Join panels horizontally
	panels := lipgloss.JoinHorizontal(lipgloss.Top, mainMenuView, subMenuView, detailsView)

	// Join panels and footer vertically
	view := lipgloss.JoinVertical(lipgloss.Left, panels, footer)

	// Render modal as overlay if active
	if m.ShowModal {
		modalView := gui.ViewModal(m)

		// Place modal over the base view (this fills the entire window with background)
		return lipgloss.Place(m.WindowWidth, m.WindowHeight,
			lipgloss.Center, lipgloss.Center,
			modalView,
			lipgloss.WithWhitespaceChars("█"),
			lipgloss.WithWhitespaceForeground(lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}))
	}

	return view
}
