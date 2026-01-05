package app

import (
	"lazynginx/pkg/commands"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	MainMenu     []string
	SubMenus     map[int][]string
	MainCursor   int
	SubCursor    int
	ActivePanel  int // 0: main menu, 1: sub menu, 2: details
	Status       string
	DetailOutput string
	WindowWidth  int
	WindowHeight int
	ShowModal    bool
	ModalType    string // "site-type", "custom-input"
	ModalCursor  int
	TextInput    string
	MainScroll   int // Scroll position for main menu
	SubScroll    int // Scroll position for submenu
	DetailScroll int // Scroll position for details panel
	IsAdmin      bool // Whether running with admin/root privileges
}

// Implement interface methods for commands.ModelInterface
func (m *Model) SetSubMenus(index int, items []string) {
	m.SubMenus[index] = items
}

// Implement interface methods for gui.ModelView
func (m Model) GetMainMenu() []string         { return m.MainMenu }
func (m Model) GetSubMenus() map[int][]string { return m.SubMenus }
func (m Model) GetMainCursor() int            { return m.MainCursor }
func (m Model) GetSubCursor() int             { return m.SubCursor }
func (m Model) GetActivePanel() int           { return m.ActivePanel }
func (m Model) GetStatus() string             { return m.Status }
func (m Model) GetDetailOutput() string       { return m.DetailOutput }
func (m Model) GetWindowWidth() int           { return m.WindowWidth }
func (m Model) GetWindowHeight() int          { return m.WindowHeight }
func (m Model) GetShowModal() bool            { return m.ShowModal }
func (m Model) GetModalType() string          { return m.ModalType }
func (m Model) GetModalCursor() int           { return m.ModalCursor }
func (m Model) GetTextInput() string          { return m.TextInput }
func (m Model) GetMainScroll() int            { return m.MainScroll }
func (m Model) GetSubScroll() int             { return m.SubScroll }
func (m Model) GetDetailScroll() int          { return m.DetailScroll }
func (m Model) GetIsAdmin() bool              { return m.IsAdmin }

// getAdminWarning returns the admin warning message if not admin
func (m Model) getAdminWarning() string {
	if !m.IsAdmin {
		return "\n\n" + strings.Repeat("─", 50) + "\n\n" +
			"⚠️  WARNING: Not running with administrator privileges\n\n" +
			"Some operations (start, stop, restart, reload) may fail.\n" +
			"Please restart LazyNginx with elevated permissions:\n\n" +
			"Windows: Run as Administrator\n" +
			"Linux/macOS: Use sudo"
	}
	return ""
}

func NewModel() Model {
	subMenus := make(map[int][]string)
	subMenus[0] = []string{"Check Status", "Test Configuration"}               // Status & Monitoring
	subMenus[1] = []string{"Start", "Stop", "Restart", "Reload Configuration"} // Service Control
	subMenus[2] = []string{"Add site", "Loading sites..."}                     // Sites - populated dynamically
	subMenus[3] = []string{"Add Reverse Proxy", "Loading reverse proxies..."}  // Reverse Proxies - populated dynamically
	subMenus[4] = []string{}                                                   // Configuration - auto-loads config file
	subMenus[5] = []string{"View Error Log", "View Access Log"}                // Logs
	subMenus[6] = []string{"Exit Application"}                                 // Quit

	// Check for admin permissions
	isAdmin := commands.IsAdmin()

	// Set initial detail message with warning if needed
	initialDetail := "Select an option from the menu"
	if !isAdmin {
		initialDetail += "\n\n" + strings.Repeat("─", 50) + "\n\n" +
			"⚠️  WARNING: Not running with administrator privileges\n\n" +
			"Some operations (start, stop, restart, reload) may fail.\n" +
			"Please restart LazyNginx with elevated permissions:\n\n" +
			"Windows: Run as Administrator\n" +
			"Linux/macOS: Use sudo"
	}

	return Model{
		MainMenu: []string{
			"Status & Monitoring",
			"Service Control",
			"Sites",
			"Reverse Proxies",
			"Configuration",
			"Logs",
			"Quit",
		},
		SubMenus:     subMenus,
		MainCursor:   0,
		SubCursor:    0,
		ActivePanel:  0,
		Status:       "",
		DetailOutput: initialDetail,
		WindowWidth:  120,
		WindowHeight: 30,
		ShowModal:    false,
		ModalType:    "",
		ModalCursor:  0,
		TextInput:    "",
		IsAdmin:      isAdmin,
	}
}

func (m Model) Init() tea.Cmd {
	return commands.CheckNginxStatus
}
