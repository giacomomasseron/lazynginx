package app

import (
	"lazynginx/pkg/commands"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseLeft:
			// Calculate which panel was clicked based on x position (horizontal layout)
			footerHeight := 1
			contentHeight := m.WindowHeight - footerHeight

			// Calculate box widths using the same logic as View()
			box1Weight := 1
			box2Weight := 1
			box3Weight := 2
			totalWeight := box1Weight + box2Weight + box3Weight

			box1Width := (m.WindowWidth * box1Weight) / totalWeight
			box2Width := (m.WindowWidth * box2Weight) / totalWeight

			panel1End := box1Width
			panel2End := panel1End + box2Width

			if msg.X < panel1End {
				// Clicked on main menu
				m.ActivePanel = 0
				// Calculate which menu item (account for border, padding, and scroll)
				availableLines := contentHeight - 5 // Account for border and title
				if availableLines < 1 {
					availableLines = 1
				}
				if msg.Y >= 3 && msg.Y < 3+availableLines {
					clickedLine := msg.Y - 3
					newCursor := m.MainScroll + clickedLine
					if newCursor < len(m.MainMenu) && newCursor != m.MainCursor {
						m.MainCursor = newCursor
						m.SubCursor = 0
						m.SubScroll = 0
						m.DetailOutput = ""
						m.DetailScroll = 0
						// Auto-load status when Status menu selected
						if m.MainCursor == 0 {
							return m, commands.CheckNginxStatus
						}
						// Auto-load sites when Sites menu selected
						if m.MainCursor == 2 {
							return m, commands.LoadSites(&m)
						}
						// Auto-load reverse proxies when Reverse Proxies menu selected
						if m.MainCursor == 3 {
							return m, commands.LoadReverseProxies(&m)
						}
						// Auto-load config when Configuration menu selected
						if m.MainCursor == 4 {
							return m, func() tea.Msg { return commands.ViewNginxConfig() }
						}
					}
				}
			} else if msg.X < panel2End {
				// Clicked on sub menu
				m.ActivePanel = 1
				subItems := m.SubMenus[m.MainCursor]
				availableLines := contentHeight - 5 // Account for border and title
				if availableLines < 1 {
					availableLines = 1
				}
				if msg.Y >= 3 && msg.Y < 3+availableLines {
					clickedLine := msg.Y - 3
					newSubCursor := m.SubScroll + clickedLine
					if newSubCursor < len(subItems) && newSubCursor != m.SubCursor {
						m.SubCursor = newSubCursor
						m.DetailOutput = ""
						m.DetailScroll = 0
						// Auto-load logs when in Logs menu
						if m.MainCursor == 5 {
							if m.SubCursor == 0 {
								return m, func() tea.Msg { return commands.ViewErrorLogs() }
							} else if m.SubCursor == 1 {
								return m, func() tea.Msg { return commands.ViewAccessLogs() }
							}
						}
						// Auto-load site config when in Sites menu (skip "Add site")
						if m.MainCursor == 2 && m.SubCursor > 0 {
							siteName := m.SubMenus[m.MainCursor][m.SubCursor]
							return m, func() tea.Msg { return commands.ViewSiteConfig(siteName) }
						}
					}
				}
			} else {
				// Clicked on details panel
				m.ActivePanel = 2
			}
		}
		return m, nil

	case tea.KeyMsg:
		// Handle modal input first
		if m.ShowModal {
			return m.handleModalInput(msg)
		}

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "left", "h":
			if m.ActivePanel > 0 {
				m.ActivePanel--
			}
			return m, nil

		case "right", "l", "tab":
			if m.ActivePanel < 2 {
				m.ActivePanel++
				if m.ActivePanel == 1 {
					m.SubCursor = 0
				}
			}
			return m, nil

		case "up", "k":
			if m.ActivePanel == 0 {
				if m.MainCursor > 0 {
					m.MainCursor--
					m.SubCursor = 0
					m.DetailOutput = ""
					m.DetailScroll = 0
					// Adjust scroll to keep cursor visible
					if m.MainCursor < m.MainScroll {
						m.MainScroll = m.MainCursor
					}
					// Auto-load status when Status menu selected
					if m.MainCursor == 0 {
						return m, commands.CheckNginxStatus
					}
					// Auto-load sites when Sites menu selected
					if m.MainCursor == 2 {
						return m, commands.LoadSites(&m)
					}
					// Auto-load reverse proxies when Reverse Proxies menu selected
					if m.MainCursor == 3 {
						return m, commands.LoadReverseProxies(&m)
					}
					// Auto-load config when Configuration menu selected
					if m.MainCursor == 4 {
						return m, func() tea.Msg { return commands.ViewNginxConfig() }
					}
				}
			} else if m.ActivePanel == 1 {
				if m.SubCursor > 0 {
					m.SubCursor--
					m.DetailOutput = ""
					m.DetailScroll = 0
					// Adjust scroll to keep cursor visible
					if m.SubCursor < m.SubScroll {
						m.SubScroll = m.SubCursor
					}
					// Auto-load logs when in Logs menu
					if m.MainCursor == 5 {
						if m.SubCursor == 0 {
							return m, func() tea.Msg { return commands.ViewErrorLogs() }
						} else if m.SubCursor == 1 {
							return m, func() tea.Msg { return commands.ViewAccessLogs() }
						}
					}
					// Auto-load site config when in Sites menu (skip "Add site")
					if m.MainCursor == 2 && m.SubCursor > 0 {
						siteName := m.SubMenus[m.MainCursor][m.SubCursor]
						return m, func() tea.Msg { return commands.ViewSiteConfig(siteName) }
					}
				}
			} else if m.ActivePanel == 2 {
				// Scroll up in details panel
				if m.DetailScroll > 0 {
					m.DetailScroll--
				}
			}
			return m, nil

		case "down", "j":
			if m.ActivePanel == 0 {
				if m.MainCursor < len(m.MainMenu)-1 {
					m.MainCursor++
					m.SubCursor = 0
					m.DetailOutput = ""
					m.DetailScroll = 0
					// Adjust scroll to keep cursor visible
					contentHeight := m.WindowHeight - 4
					if contentHeight < 5 {
						contentHeight = 5
					}
					availableLines := contentHeight - 3
					if m.MainCursor >= m.MainScroll+availableLines {
						m.MainScroll = m.MainCursor - availableLines + 1
					}
					// Auto-load status when Status menu selected
					if m.MainCursor == 0 {
						return m, commands.CheckNginxStatus
					}
					// Auto-load sites when Sites menu selected
					if m.MainCursor == 2 {
						return m, commands.LoadSites(&m)
					}
					// Auto-load reverse proxies when Reverse Proxies menu selected
					if m.MainCursor == 3 {
						return m, commands.LoadReverseProxies(&m)
					}
					// Auto-load config when Configuration menu selected
					if m.MainCursor == 4 {
						return m, func() tea.Msg { return commands.ViewNginxConfig() }
					}
				}
			} else if m.ActivePanel == 1 {
				subItems := m.SubMenus[m.MainCursor]
				if m.SubCursor < len(subItems)-1 {
					m.SubCursor++
					m.DetailOutput = ""
					m.DetailScroll = 0
					// Adjust scroll to keep cursor visible
					contentHeight := m.WindowHeight - 4
					if contentHeight < 5 {
						contentHeight = 5
					}
					availableLines := contentHeight - 3
					if m.SubCursor >= m.SubScroll+availableLines {
						m.SubScroll = m.SubCursor - availableLines + 1
					}
					// Auto-load logs when in Logs menu
					if m.MainCursor == 5 {
						if m.SubCursor == 0 {
							return m, func() tea.Msg { return commands.ViewErrorLogs() }
						} else if m.SubCursor == 1 {
							return m, func() tea.Msg { return commands.ViewAccessLogs() }
						}
					}
					// Auto-load site config when in Sites menu (skip "Add site")
					if m.MainCursor == 2 && m.SubCursor > 0 {
						siteName := m.SubMenus[m.MainCursor][m.SubCursor]
						return m, func() tea.Msg { return commands.ViewSiteConfig(siteName) }
					}
				}
			} else if m.ActivePanel == 2 {
				// Scroll down in details panel
				m.DetailScroll++
			}
			return m, nil

		case "enter":
			if m.ActivePanel == 1 {
				// Check if it's "Stop" in Service Control menu
				if m.MainCursor == 1 && m.SubCursor == 1 {
					m.ShowModal = true
					m.ModalType = "confirm-stop"
					m.ModalCursor = 0
					return m, nil
				}
				// Check if it's "Add site" in Sites menu
				if m.MainCursor == 2 && m.SubCursor == 0 {
					m.ShowModal = true
					m.ModalType = "site-type"
					m.ModalCursor = 0
					return m, nil
				}
				// Check if it's "Add Reverse Proxy" in Reverse Proxies menu
				if m.MainCursor == 3 && m.SubCursor == 0 {
					m.ShowModal = true
					m.ModalType = "proxy-type"
					m.ModalCursor = 0
					return m, nil
				}
				// Otherwise execute the selection
				return m, m.handleSelection()
			}
			return m, nil

		case "d":
			// Delete key - only works in Sites submenu for actual sites (not "Add site")
			if m.ActivePanel == 1 && m.MainCursor == 2 && m.SubCursor > 0 {
				subItems := m.SubMenus[m.MainCursor]
				if m.SubCursor < len(subItems) {
					siteName := subItems[m.SubCursor]
					// Don't allow deleting placeholder items
					if siteName != "Loading sites..." && siteName != "No sites found" {
						m.ShowModal = true
						m.ModalType = "confirm-delete-site"
						m.ModalCursor = 0
						return m, nil
					}
				}
			}
			return m, nil
		}

	case commands.StatusMsg:
		m.Status = msg.Status
		m.DetailOutput = msg.Status + m.getAdminWarning() // Also display in details panel with warning
		m.DetailScroll = 0                                // Reset scroll on new content
		return m, nil

	case commands.OutputMsg:
		m.DetailOutput = msg.Output + m.getAdminWarning()
		m.DetailScroll = 0 // Reset scroll on new content

		// Check if we need to reload sites after add/delete operations
		if m.MainCursor == 2 && (strings.Contains(msg.Output, "Site '") && (strings.Contains(msg.Output, "created successfully") || strings.Contains(msg.Output, "deleted successfully"))) {
			// Reload sites list after successful add or delete
			return m, commands.LoadSites(&m)
		}

		// Check if we need to reload reverse proxies after add operation
		if m.MainCursor == 3 && strings.Contains(msg.Output, "Reverse proxy") && strings.Contains(msg.Output, "created successfully") {
			// Reload reverse proxies list after successful add
			return m, commands.LoadReverseProxies(&m)
		}

		return m, nil

	case tea.WindowSizeMsg:
		m.WindowWidth = msg.Width
		m.WindowHeight = msg.Height
		return m, nil
	}

	return m, nil
}
