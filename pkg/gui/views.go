package gui

import (
	"lazynginx/pkg/utils"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/jesseduffield/lazycore/pkg/boxlayout"
)

// ModelView defines the interface needed from the model for rendering
type ModelView interface {
	GetMainMenu() []string
	GetSubMenus() map[int][]string
	GetMainCursor() int
	GetSubCursor() int
	GetActivePanel() int
	GetStatus() string
	GetDetailOutput() string
	GetWindowWidth() int
	GetWindowHeight() int
	GetShowModal() bool
	GetModalType() string
	GetModalCursor() int
	GetTextInput() string
	GetMainScroll() int
	GetSubScroll() int
	GetDetailScroll() int
}

func ViewMainMenuWithDim(m ModelView, dim boxlayout.Dimensions) string {
	// Calculate dimensions from the box
	boxWidth := dim.X1 - dim.X0 + 1
	boxHeight := dim.Y1 - dim.Y0 + 1

	s := strings.Builder{}
	s.WriteString(TitleStyle.Render(" Main Menu ") + "\n\n")

	// Calculate available height for content
	// Border takes 2 lines (top + bottom), content already includes padding
	contentHeight := boxHeight - 2
	if contentHeight < 5 {
		contentHeight = 5
	}
	availableLines := contentHeight - 2 // Reserve space for title and spacing

	mainMenu := m.GetMainMenu()
	mainCursor := m.GetMainCursor()
	mainScroll := m.GetMainScroll()
	activePanel := m.GetActivePanel()

	// Clamp scroll position
	maxScroll := len(mainMenu) - availableLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	scrollPos := mainScroll
	if scrollPos > maxScroll {
		scrollPos = maxScroll
	}
	if scrollPos < 0 {
		scrollPos = 0
	}

	// Calculate visible range
	startLine := scrollPos
	endLine := startLine + availableLines
	if endLine > len(mainMenu) {
		endLine = len(mainMenu)
	}

	// Check if scrollbar is needed
	showScrollbar := len(mainMenu) > availableLines

	// Scrollbar is OUTSIDE the box, so we need to reduce box width by 1 char
	actualBoxWidth := boxWidth
	if showScrollbar {
		actualBoxWidth = boxWidth - 1
	}

	// Content width = box width - border (2) - padding (2)
	contentWidth := actualBoxWidth - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Render visible items (no scrollbar inside content)
	for idx, i := range make([]int, endLine-startLine) {
		i = startLine + idx
		choice := mainMenu[i]
		cursor := "  "
		var line string
		if mainCursor == i && activePanel == 0 {
			cursor = "▶ "
			line = SelectedStyle.Render(cursor + choice)
		} else if mainCursor == i {
			cursor = "▶ "
			line = ActiveStyle.Render(cursor + choice)
		} else {
			line = NormalStyle.Render(cursor + choice)
		}

		s.WriteString(line + "\n")
	}

	// Fill remaining lines with empty space to ensure consistent height
	for len := endLine - startLine; len < availableLines; len++ {
		s.WriteString("\n")
	}

	// Ensure content has exact number of lines for consistent height
	contentStr := s.String()
	builtLines := strings.Split(contentStr, "\n")

	expectedLines := contentHeight
	if len(builtLines) > expectedLines {
		contentStr = strings.Join(builtLines[:expectedLines], "\n")
	} else if len(builtLines) < expectedLines {
		for i := len(builtLines); i < expectedLines; i++ {
			contentStr += "\n"
		}
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(actualBoxWidth - 2).
		Height(boxHeight - 2)

	renderedBox := boxStyle.Render(contentStr)

	// If scrollbar needed, build it separately and join horizontally
	if showScrollbar {
		scrollbarStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

		scrollbarBuilder := strings.Builder{}

		// Calculate scrollbar thumb position
		scrollbarHeight := boxHeight
		thumbSize := utils.Max(1, scrollbarHeight*availableLines/len(mainMenu))
		thumbStart := scrollbarHeight * scrollPos / len(mainMenu)
		thumbEnd := thumbStart + thumbSize

		// Build scrollbar for full box height
		for i := 0; i < boxHeight; i++ {
			var scrollChar string
			// Offset by 1 for top border
			if i >= thumbStart+1 && i < thumbEnd+1 {
				scrollChar = scrollbarStyle.Render("█")
			} else {
				scrollChar = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("░")
			}
			scrollbarBuilder.WriteString(scrollChar)
			if i < boxHeight-1 {
				scrollbarBuilder.WriteString("\n")
			}
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, renderedBox, scrollbarBuilder.String())
	}

	return renderedBox
}

func ViewSubMenuWithDim(m ModelView, dim boxlayout.Dimensions) string {
	// Calculate dimensions from the box
	boxWidth := dim.X1 - dim.X0 + 1
	boxHeight := dim.Y1 - dim.Y0 + 1

	s := strings.Builder{}
	s.WriteString(TitleStyle.Render(" Options ") + "\n\n")

	// Calculate available height for content
	// Border takes 2 lines (top + bottom), content already includes padding
	contentHeight := boxHeight - 2
	if contentHeight < 5 {
		contentHeight = 5
	}
	availableLines := contentHeight - 2 // Reserve space for title and spacing

	subMenus := m.GetSubMenus()
	mainCursor := m.GetMainCursor()
	subCursor := m.GetSubCursor()
	subScroll := m.GetSubScroll()
	activePanel := m.GetActivePanel()
	subItems := subMenus[mainCursor]

	// Clamp scroll position
	maxScroll := len(subItems) - availableLines
	if maxScroll < 0 {
		maxScroll = 0
	}
	scrollPos := subScroll
	if scrollPos > maxScroll {
		scrollPos = maxScroll
	}
	if scrollPos < 0 {
		scrollPos = 0
	}

	// Calculate visible range
	startLine := scrollPos
	endLine := startLine + availableLines
	if endLine > len(subItems) {
		endLine = len(subItems)
	}

	// Check if scrollbar is needed
	showScrollbar := len(subItems) > availableLines

	// Scrollbar is OUTSIDE the box, so we need to reduce box width by 1 char
	actualBoxWidth := boxWidth
	if showScrollbar {
		actualBoxWidth = boxWidth - 1
	}

	// Content width = box width - border (2) - padding (2)
	contentWidth := actualBoxWidth - 4
	if contentWidth < 10 {
		contentWidth = 10
	}

	// Render visible items (no scrollbar inside content)
	for idx, i := range make([]int, endLine-startLine) {
		i = startLine + idx
		choice := subItems[i]
		cursor := "  "
		var line string
		if subCursor == i && activePanel == 1 {
			cursor = "▶ "
			line = SelectedStyle.Render(cursor + choice)
		} else if subCursor == i {
			cursor = "▶ "
			line = ActiveStyle.Render(cursor + choice)
		} else {
			line = NormalStyle.Render(cursor + choice)
		}

		s.WriteString(line + "\n")
	}

	// Fill remaining lines with empty space to ensure consistent height
	for len := endLine - startLine; len < availableLines; len++ {
		s.WriteString("\n")
	}

	// Ensure content has exact number of lines for consistent height
	contentStr := s.String()
	builtLines := strings.Split(contentStr, "\n")

	expectedLines := contentHeight
	if len(builtLines) > expectedLines {
		contentStr = strings.Join(builtLines[:expectedLines], "\n")
	} else if len(builtLines) < expectedLines {
		for i := len(builtLines); i < expectedLines; i++ {
			contentStr += "\n"
		}
	}

	boxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(actualBoxWidth - 2).
		Height(boxHeight - 2)

	renderedBox := boxStyle.Render(contentStr)

	// If scrollbar needed, build it separately and join horizontally
	if showScrollbar {
		scrollbarStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

		scrollbarBuilder := strings.Builder{}

		// Calculate scrollbar thumb position
		scrollbarHeight := boxHeight
		thumbSize := utils.Max(1, scrollbarHeight*availableLines/len(subItems))
		thumbStart := scrollbarHeight * scrollPos / len(subItems)
		thumbEnd := thumbStart + thumbSize

		// Build scrollbar for full box height
		for i := 0; i < boxHeight; i++ {
			var scrollChar string
			// Offset by 1 for top border
			if i >= thumbStart+1 && i < thumbEnd+1 {
				scrollChar = scrollbarStyle.Render("█")
			} else {
				scrollChar = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("░")
			}
			scrollbarBuilder.WriteString(scrollChar)
			if i < boxHeight-1 {
				scrollbarBuilder.WriteString("\n")
			}
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, renderedBox, scrollbarBuilder.String())
	}

	return renderedBox
}

func ViewDetailsWithDim(m ModelView, dim boxlayout.Dimensions) string {
	// Calculate dimensions from the box
	boxWidth := dim.X1 - dim.X0 + 1
	boxHeight := dim.Y1 - dim.Y0 + 1

	s := strings.Builder{}
	s.WriteString(TitleStyle.Render(" Details ") + "\n\n")

	detailOutput := m.GetDetailOutput()
	detailScroll := m.GetDetailScroll()

	// Split content into lines for scrolling
	contentLines := strings.Split(detailOutput, "\n")

	// Calculate available height for content
	// Border takes 2 lines (top + bottom), content already includes padding
	contentHeight := boxHeight - 2
	if contentHeight < 5 {
		contentHeight = 5
	}

	availableLines := contentHeight - 2
	if availableLines < 1 {
		availableLines = 1
	}

	// Check if vertical scrollbar is needed
	showVScrollbar := len(contentLines) > availableLines

	// Scrollbar is OUTSIDE the box, so we need to reduce box width by 1 char
	actualBoxWidth := boxWidth
	if showVScrollbar {
		actualBoxWidth = boxWidth - 1
	}

	// Content width = box width - border (2) - padding (2)
	contentWidth := actualBoxWidth - 4
	if contentWidth < 20 {
		contentWidth = 20
	}

	// Calculate scroll boundaries for vertical scrolling
	maxScroll := len(contentLines) - availableLines
	if maxScroll < 0 {
		maxScroll = 0
	}

	// Clamp vertical scroll position for display
	scrollPos := detailScroll
	if scrollPos > maxScroll {
		scrollPos = maxScroll
	}
	if scrollPos < 0 {
		scrollPos = 0
	}

	// Get visible lines based on vertical scroll position
	startLine := scrollPos
	endLine := startLine + availableLines
	if endLine > len(contentLines) {
		endLine = len(contentLines)
	}

	// Render visible content (no scrollbar inside content)
	for idx := 0; idx < endLine-startLine; idx++ {
		line := contentLines[startLine+idx]

		// Replace tabs with spaces for consistent rendering
		line = strings.ReplaceAll(line, "\t", "    ")

		// Truncate the line to contentWidth
		runes := []rune(line)
		if len(runes) > contentWidth-3 {
			line = string(runes[:contentWidth-3]) + "..."
		}

		s.WriteString(line + "\n")
	}

	// Fill remaining lines with empty space to ensure consistent height
	for len := endLine - startLine; len < availableLines; len++ {
		s.WriteString("\n")
	}

	// Ensure content has exact number of lines for consistent height
	contentStr := s.String()
	builtLines := strings.Split(contentStr, "\n")

	expectedLines := contentHeight
	if len(builtLines) > expectedLines {
		contentStr = strings.Join(builtLines[:expectedLines], "\n")
	} else if len(builtLines) < expectedLines {
		for i := len(builtLines); i < expectedLines; i++ {
			contentStr += "\n"
		}
	}

	detailBoxStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7D56F4")).
		Padding(0, 1).
		Width(actualBoxWidth - 2).
		Height(boxHeight - 2)

	renderedBox := detailBoxStyle.Render(contentStr)

	// If scrollbar needed, build it separately and join horizontally
	if showVScrollbar {
		scrollbarStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)

		scrollbarBuilder := strings.Builder{}

		// Calculate scrollbar thumb position
		scrollbarHeight := boxHeight
		thumbSize := utils.Max(1, scrollbarHeight*availableLines/len(contentLines))
		thumbStart := scrollbarHeight * scrollPos / len(contentLines)
		thumbEnd := thumbStart + thumbSize

		// Build scrollbar for full box height
		for i := 0; i < boxHeight; i++ {
			var scrollChar string
			// Offset by 1 for top border
			if i >= thumbStart+1 && i < thumbEnd+1 {
				scrollChar = scrollbarStyle.Render("█")
			} else {
				scrollChar = lipgloss.NewStyle().Foreground(lipgloss.Color("#7D56F4")).Render("░")
			}
			scrollbarBuilder.WriteString(scrollChar)
			if i < boxHeight-1 {
				scrollbarBuilder.WriteString("\n")
			}
		}

		return lipgloss.JoinHorizontal(lipgloss.Top, renderedBox, scrollbarBuilder.String())
	}

	return renderedBox
}

func ViewFooter(m ModelView, windowWidth int) string {
	// Build keybindings based on active panel (like lazydocker)
	var keybindings string

	activePanel := m.GetActivePanel()
	mainCursor := m.GetMainCursor()
	subCursor := m.GetSubCursor()

	switch activePanel {
	case 0: // Main menu
		keybindings = "[↑↓/jk] scroll [→/l/tab] next panel [enter] select [mouse] scroll/click [q] quit"
	case 1: // Sub menu
		// Check if we're in Sites menu with a site selected (not "Add site")
		if mainCursor == 2 && subCursor > 0 {
			keybindings = "[↑↓/jk] scroll [←/h] prev panel [→/l/tab] next panel [enter] execute [d] delete [mouse] scroll/click [q] quit"
		} else {
			keybindings = "[↑↓/jk] scroll [←/h] prev panel [→/l/tab] next panel [enter] execute [mouse] scroll/click [q] quit"
		}
	case 2: // Details
		keybindings = "[↑↓/jk] scroll [←/h] prev panel [mouse] scroll/click [q] quit"
	}

	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#1E3A8A")).
		Width(windowWidth).
		Padding(0, 1).
		Bold(true)

	return footerStyle.Render(keybindings)
}

func ViewModal(m ModelView) string {
	var content string

	modalType := m.GetModalType()
	modalCursor := m.GetModalCursor()
	textInput := m.GetTextInput()

	if modalType == "confirm-stop" {
		title := " Confirm Stop "
		options := []string{"Yes", "No"}

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Are you sure you want to stop Nginx?\n\n")

		for i, opt := range options {
			cursor := "  "
			if modalCursor == i {
				cursor = "▶ "
				s.WriteString(SelectedStyle.Render(cursor+opt) + "\n")
			} else {
				s.WriteString(NormalStyle.Render(cursor+opt) + "\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(InfoStyle.Render("↑/↓: Navigate | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "confirm-delete-site" {
		title := " Confirm Delete Site "
		options := []string{"Yes", "No"}

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Are you sure you want to delete this site?\n")
		s.WriteString("This will remove the configuration file.\n\n")

		for i, opt := range options {
			cursor := "  "
			if modalCursor == i {
				cursor = "▶ "
				s.WriteString(SelectedStyle.Render(cursor+opt) + "\n")
			} else {
				s.WriteString(NormalStyle.Render(cursor+opt) + "\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(InfoStyle.Render("↑/↓: Navigate | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "site-type" {
		title := " Add New Site "
		options := []string{"Laravel", "Static Website", "Vanilla PHP", "Custom"}

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Select site type:\n\n")

		for i, opt := range options {
			cursor := "  "
			if modalCursor == i {
				cursor = "▶ "
				s.WriteString(SelectedStyle.Render(cursor+opt) + "\n")
			} else {
				s.WriteString(NormalStyle.Render(cursor+opt) + "\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(InfoStyle.Render("↑/↓: Navigate | Enter: Select | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "custom-input" {
		title := " Custom Site Name "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Enter site name:\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type site name | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "laravel-input" {
		title := " Laravel Site Name "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Enter Laravel site name:\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type site name | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "static-input" {
		title := " Static Website Name "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Enter static website name:\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type site name | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "vanilla-php-input" {
		title := " Vanilla PHP Site Name "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Enter PHP site name:\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type site name | Enter: Confirm | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "proxy-type" {
		title := " Add Reverse Proxy "
		options := []string{"Simple Proxy", "Load Balanced"}

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString("Select proxy type:\n\n")

		for i, opt := range options {
			cursor := "  "
			if modalCursor == i {
				cursor = "▶ "
				s.WriteString(SelectedStyle.Render(cursor+opt) + "\n")
			} else {
				s.WriteString(NormalStyle.Render(cursor+opt) + "\n")
			}
		}

		s.WriteString("\n")
		s.WriteString(InfoStyle.Render("↑/↓: Navigate | Enter: Select | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "proxy-location-input" {
		title := " Simple Reverse Proxy - Step 1/2 "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString(NormalStyle.Render("Enter nginx location path:") + "\n")
		s.WriteString(InfoStyle.Render("Example: / or /api or /app") + "\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type location | Enter: Next | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "proxy-location-input-lb" {
		title := " Load Balanced Proxy - Step 1/2 "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString(NormalStyle.Render("Enter nginx location path:") + "\n")
		s.WriteString(InfoStyle.Render("Example: / or /api or /app") + "\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type location | Enter: Next | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "proxy-host-input" {
		title := " Simple Reverse Proxy - Step 2/2 "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString(NormalStyle.Render("Enter backend host:") + "\n")
		s.WriteString(InfoStyle.Render("Example: http://localhost:3000 or http://192.168.1.10:8080") + "\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type host | Enter: Create | Esc: Cancel") + "\n")
		content = s.String()
	} else if modalType == "proxy-host-input-lb" {
		title := " Load Balanced Proxy - Step 2/2 "

		s := strings.Builder{}
		s.WriteString(TitleStyle.Render(title) + "\n\n")
		s.WriteString(NormalStyle.Render("Enter backend hosts (comma-separated):") + "\n")
		s.WriteString(InfoStyle.Render("Example: localhost:3000,localhost:3001,localhost:3002") + "\n\n")
		s.WriteString(SelectedStyle.Render(" "+textInput+"█ ") + "\n\n")
		s.WriteString(InfoStyle.Render("Type hosts | Enter: Create | Esc: Cancel") + "\n")
		content = s.String()
	}

	modalStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#FF79C6")).
		Padding(1, 2).
		Width(50).
		AlignHorizontal(lipgloss.Center)

	return modalStyle.Render(content)
}
