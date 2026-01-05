# LazyNginx AI Development Guide

## Project Overview
LazyNginx is a terminal-based Nginx manager built with Go and Bubble Tea (TUI framework). It provides cross-platform menu-driven nginx management without requiring users to remember commands.

## Architecture

### Core Components
- **main.go**: Entry point that initializes Bubble Tea program with `tea.WithAltScreen()` for full-screen TUI
- **model.go**: Bubble Tea model implementing Update/View/Init pattern with two display modes:
  - `mode: "menu"` - Shows interactive menu with cursor navigation
  - `mode: "output"` - Displays command results (switched via Enter key)
- **commands.go**: Command execution functions returning `tea.Msg` types for async operations
- **messages.go**: Custom message types (`statusMsg`, `outputMsg`) for Bubble Tea's Elm architecture
- **styles.go**: Lipgloss style definitions for consistent terminal UI theming

### Bubble Tea Pattern
The app uses Bubble Tea's command pattern where:
1. User selects menu item → `handleSelection()` returns a `tea.Cmd`
2. Command executes asynchronously → returns a message (`statusMsg` or `outputMsg`)
3. `Update()` receives message → updates model state
4. `View()` renders based on current mode

## Platform-Specific Behavior

Commands try multiple approaches in sequence (returns first success):

**Start/Stop/Restart Operations:**
1. Windows: `net start/stop nginx`
2. Linux/systemd: `sudo systemctl [action] nginx`
3. Direct: `sudo nginx [-s stop]`

**Status Check Cascading:**
1. Windows: `tasklist /FI "IMAGENAME eq nginx.exe"`
2. systemd: `systemctl is-active nginx`
3. Unix: `pgrep -x nginx`
4. Fallback: `ps aux | grep nginx`

**File Discovery:**
Config/logs checked in order:
- `/etc/nginx/` (Linux)
- `C:\nginx\` (Windows)
- `/usr/local/nginx/` (macOS/Unix)

## Development Patterns

### Adding New Commands
1. Create function in `commands.go` returning `tea.Msg`:
   ```go
   func myCommand() tea.Msg {
       output, err := exec.Command("mycommand").CombinedOutput()
       if err != nil {
           return outputMsg{output: "Error: " + string(output)}
       }
       return outputMsg{output: string(output)}
   }
   ```

2. Add menu item to `initialModel()` choices slice in [model.go](model.go#L21-L32)

3. Map cursor position to command in `handleSelection()` switch statement in [model.go](model.go#L147-L163)

### Navigation Keys
- Vim-style (`j`/`k`) and arrow keys handled in [model.go](model.go#L64-L78) Update() switch
- `q` and `Ctrl+C` handled differently based on `mode` (quit vs return to menu)

### Styling Convention
All UI styles use Lipgloss with hex color codes defined as package variables in [styles.go](styles.go). Status messages conditionally use `statusStyle` (green) or `errorStyle` (red) based on text content matching in [model.go](model.go#L108-L114).

## Building and Running

```bash
# Build
go build -o lazynginx

# Run (requires sudo for service operations on Linux/macOS)
sudo ./lazynginx
```

No tests exist. Manual testing requires nginx installed and verifying menu navigation + command execution on target OS.

## Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework using Elm architecture
- `github.com/charmbracelet/lipgloss` - Terminal styling (colors, padding, bold)

Binary output: `lazynginx` (Unix) or `lazynginx.exe` (Windows)
