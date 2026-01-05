package app

import (
	"lazynginx/pkg/commands"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) handleModalInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Close modal
		m.ShowModal = false
		m.ModalType = ""
		m.TextInput = ""
		return m, nil

	case "up", "k":
		// Only handle navigation for selection modals, not text input modals
		if m.ModalType == "custom-input" || m.ModalType == "laravel-input" ||
			m.ModalType == "static-input" || m.ModalType == "vanilla-php-input" ||
			m.ModalType == "proxy-input" || m.ModalType == "proxy-input-lb" {
			// For text input modals, let these keys fall through to default handler
			if msg.String() == "k" {
				key := msg.String()
				if len(key) == 1 {
					m.TextInput += key
				}
			}
			return m, nil
		}
		if m.ModalType == "confirm-stop" && m.ModalCursor > 0 {
			m.ModalCursor--
		} else if m.ModalType == "confirm-delete-site" && m.ModalCursor > 0 {
			m.ModalCursor--
		} else if m.ModalType == "site-type" && m.ModalCursor > 0 {
			m.ModalCursor--
		} else if m.ModalType == "proxy-type" && m.ModalCursor > 0 {
			m.ModalCursor--
		}
		return m, nil

	case "down", "j":
		// Only handle navigation for selection modals, not text input modals
		if m.ModalType == "custom-input" || m.ModalType == "laravel-input" ||
			m.ModalType == "static-input" || m.ModalType == "vanilla-php-input" ||
			m.ModalType == "proxy-input" || m.ModalType == "proxy-input-lb" {
			// For text input modals, let these keys fall through to default handler
			if msg.String() == "j" {
				key := msg.String()
				if len(key) == 1 {
					m.TextInput += key
				}
			}
			return m, nil
		}
		if m.ModalType == "confirm-stop" && m.ModalCursor < 1 {
			m.ModalCursor++
		} else if m.ModalType == "confirm-delete-site" && m.ModalCursor < 1 {
			m.ModalCursor++
		} else if m.ModalType == "site-type" && m.ModalCursor < 3 {
			m.ModalCursor++
		} else if m.ModalType == "proxy-type" && m.ModalCursor < 1 {
			m.ModalCursor++
		}
		return m, nil

	case "enter":
		if m.ModalType == "confirm-stop" {
			if m.ModalCursor == 0 {
				// Yes selected - execute stop
				m.ShowModal = false
				m.ModalType = ""
				return m, commands.StopNginx
			} else {
				// No selected - cancel
				m.ShowModal = false
				m.ModalType = ""
				return m, nil
			}
		} else if m.ModalType == "confirm-delete-site" {
			if m.ModalCursor == 0 {
				// Yes selected - execute delete
				subItems := m.SubMenus[m.MainCursor]
				siteName := subItems[m.SubCursor]
				m.ShowModal = false
				m.ModalType = ""
				return m, func() tea.Msg {
					return commands.DeleteSite(siteName)
				}
			} else {
				// No selected - cancel
				m.ShowModal = false
				m.ModalType = ""
				return m, nil
			}
		} else if m.ModalType == "site-type" {
			if m.ModalCursor == 0 {
				// Laravel selected - show text input modal for Laravel site name
				m.ModalType = "laravel-input"
				m.TextInput = ""
				return m, nil
			} else if m.ModalCursor == 1 {
				// Static Website selected - show text input modal
				m.ModalType = "static-input"
				m.TextInput = ""
				return m, nil
			} else if m.ModalCursor == 2 {
				// Vanilla PHP selected - show text input modal
				m.ModalType = "vanilla-php-input"
				m.TextInput = ""
				return m, nil
			} else {
				// Custom selected - show text input modal
				m.ModalType = "custom-input"
				m.TextInput = ""
				return m, nil
			}
		} else if m.ModalType == "laravel-input" {
			// Submit Laravel site name
			m.ShowModal = false
			siteName := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddSite("Laravel", siteName)
			}
		} else if m.ModalType == "static-input" {
			// Submit Static Website site name
			m.ShowModal = false
			siteName := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddSite("Static", siteName)
			}
		} else if m.ModalType == "vanilla-php-input" {
			// Submit Vanilla PHP site name
			m.ShowModal = false
			siteName := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddSite("VanillaPHP", siteName)
			}
		} else if m.ModalType == "proxy-type" {
			if m.ModalCursor == 0 {
				// Simple Proxy selected - show text input modal for custom configuration
				m.ModalType = "proxy-input"
				m.TextInput = ""
				return m, nil
			} else {
				// Load Balanced selected
				m.ModalType = "proxy-input-lb"
				m.TextInput = ""
				return m, nil
			}
		} else if m.ModalType == "proxy-input" {
			// Submit proxy configuration
			m.ShowModal = false
			proxyConfig := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddProxy("Simple", proxyConfig)
			}
		} else if m.ModalType == "proxy-input-lb" {
			// Submit load balanced proxy configuration
			m.ShowModal = false
			proxyConfig := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddProxy("LoadBalanced", proxyConfig)
			}
		} else if m.ModalType == "custom-input" {
			// Submit custom site name
			m.ShowModal = false
			siteName := m.TextInput
			m.TextInput = ""
			return m, func() tea.Msg {
				return commands.AddSite("Custom", siteName)
			}
		}
		return m, nil

	case "backspace":
		if m.ModalType == "custom-input" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		} else if m.ModalType == "laravel-input" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		} else if m.ModalType == "static-input" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		} else if m.ModalType == "vanilla-php-input" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		} else if m.ModalType == "proxy-input" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		} else if m.ModalType == "proxy-input-lb" && len(m.TextInput) > 0 {
			m.TextInput = m.TextInput[:len(m.TextInput)-1]
		}
		return m, nil

	default:
		// Handle text input for custom site name
		if m.ModalType == "custom-input" {
			// Only accept alphanumeric, dash, underscore, and dot
			key := msg.String()
			if len(key) == 1 {
				m.TextInput += key
			}
		} else if m.ModalType == "laravel-input" {
			// Only accept alphanumeric, dash, underscore, and dot
			key := msg.String()
			if len(key) == 1 {
				m.TextInput += key
			}
		} else if m.ModalType == "static-input" {
			// Only accept alphanumeric, dash, underscore, and dot
			key := msg.String()
			if len(key) == 1 {
				m.TextInput += key
			}
		} else if m.ModalType == "vanilla-php-input" {
			// Only accept alphanumeric, dash, underscore, and dot
			key := msg.String()
			if len(key) == 1 {
				m.TextInput += key
			}
		} else if m.ModalType == "proxy-input" || m.ModalType == "proxy-input-lb" {
			// Accept any printable characters for proxy config
			key := msg.String()
			if len(key) == 1 {
				m.TextInput += key
			}
		}
		return m, nil
	}
}

func (m Model) handleSelection() tea.Cmd {
	// Main menu indices:
	// 0=Status & Monitoring, 1=Service Control, 2=Sites, 3=Reverse Proxies, 4=Configuration, 5=Logs, 6=Quit
	switch m.MainCursor {
	case 0: // Status & Monitoring
		switch m.SubCursor {
		case 0:
			return commands.CheckNginxStatus
		case 1:
			return commands.TestNginxConfig
		}
	case 1: // Service Control
		switch m.SubCursor {
		case 0:
			return commands.StartNginx
		case 1:
			return commands.StopNginx
		case 2:
			return commands.RestartNginx
		case 3:
			return commands.ReloadNginx
		}
	case 2: // Sites
		subItems := m.SubMenus[m.MainCursor]
		// Skip index 0 (Add site) - that's handled in the enter key
		if m.SubCursor > 0 && m.SubCursor < len(subItems) {
			siteName := subItems[m.SubCursor]
			return func() tea.Msg {
				return commands.ViewSiteConfig(siteName)
			}
		}
	case 3: // Reverse Proxies
		// TODO: Implement reverse proxy viewing
		return nil
	case 4: // Configuration
		// Auto-loaded, but can also be triggered manually
		return commands.ViewNginxConfig
	case 5: // Logs
		switch m.SubCursor {
		case 0:
			return commands.ViewErrorLogs
		case 1:
			return commands.ViewAccessLogs
		}
	case 6: // Quit
		return tea.Quit
	}
	return nil
}
