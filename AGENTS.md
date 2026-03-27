# LazyNginx Development Guide for Coding Agents

## Project Overview
LazyNginx is a cross-platform terminal-based Nginx manager built with Go and Bubble Tea (TUI framework). It provides an interactive menu-driven interface for managing nginx without remembering commands.

## Build, Test, and Run Commands

### Building
```bash
# Build for current platform
go build -o lazynginx

# Build with specific output name
go build -o lazynginx main.go

# The binary output is: lazynginx (Unix) or lazynginx.exe (Windows)
```

### Running
```bash
# Run directly (requires sudo for service operations on Linux/macOS)
sudo ./lazynginx

# Run without building
go run main.go

# Windows: Run as Administrator
lazynginx.exe
```

### Testing
```bash
# No tests currently exist in this project
# Manual testing requires nginx installed and verifying menu navigation + command execution

# To test manually:
# 1. Build: go build -o lazynginx
# 2. Run with sudo: sudo ./lazynginx
# 3. Test menu navigation with j/k or arrow keys
# 4. Verify each command executes correctly on your OS
```

### Dependencies
```bash
# Download dependencies
go mod download

# Tidy dependencies
go mod tidy

# Update dependencies
go get -u ./...
```

### Linting
```bash
# No golangci-lint configuration exists
# Use standard Go tools:
go fmt ./...
go vet ./...

# Install and run golangci-lint manually if needed:
# golangci-lint run
```

## Architecture and Code Organization

### Directory Structure
```
lazynginx/
├── main.go                  # Entry point - initializes Bubble Tea program
├── pkg/
│   ├── app/                 # Bubble Tea model and handlers
│   │   ├── app.go          # Model definition and initialization
│   │   ├── handlers.go     # Business logic handlers
│   │   ├── update.go       # Bubble Tea Update() implementation
│   │   └── view.go         # Bubble Tea View() rendering
│   ├── commands/           # Command execution functions
│   │   └── commands.go     # Nginx operations (start/stop/status/logs)
│   ├── supervisor/         # Process supervisor detection and control
│   │   └── supervisor.go   # Multi-supervisor support (systemd/s6/supervisord/runit)
│   ├── gui/                # UI styling and views
│   │   ├── styles.go       # Lipgloss style definitions
│   │   └── views.go        # Reusable view components
│   └── utils/              # Helper utilities
│       └── helpers.go      # Math helpers (Max, Min)
└── go.mod                   # Dependencies
```

### Core Dependencies
- `github.com/charmbracelet/bubbletea` - TUI framework using Elm architecture
- `github.com/charmbracelet/lipgloss` - Terminal styling (colors, padding, bold)
- `github.com/jesseduffield/lazycore` - Core utilities

### Bubble Tea Pattern
The app follows Bubble Tea's Elm architecture:
1. User interaction → `Update()` receives `tea.Msg`
2. `Update()` modifies `Model` state and returns `tea.Cmd`
3. Commands execute asynchronously → return custom messages (`StatusMsg`, `OutputMsg`, `ConfigViewMsg`)
4. `View()` renders based on current model state

### Display Modes and Panels
- `ActivePanel`: 0=main menu, 1=submenu, 2=details
- Modal types: "site-type", "custom-input"
- Scroll tracking: `MainScroll`, `SubScroll`, `DetailScroll`

## Code Style Guidelines

### Import Organization
```go
import (
    // Standard library first
    "fmt"
    "os"
    "strings"
    
    // Third-party packages
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
    
    // Local packages
    "lazynginx/pkg/commands"
    "lazynginx/pkg/gui"
)
```

### Naming Conventions
- **Types**: PascalCase (e.g., `Model`, `StatusMsg`, `OutputMsg`)
- **Functions**: camelCase for unexported, PascalCase for exported (e.g., `NewModel()`, `handleSelection()`)
- **Variables**: camelCase (e.g., `mainCursor`, `detailOutput`)
- **Constants**: PascalCase with lipgloss colors (e.g., `FocusedBorderColor`)
- **Package names**: lowercase, single word (e.g., `app`, `commands`, `gui`, `utils`)

### Function Structure
```go
// Commands return tea.Msg types for Bubble Tea integration
func CommandName() tea.Msg {
    var cmd *exec.Cmd
    var output []byte
    var err error
    
    // Try platform-specific approaches in sequence
    // Windows → systemd → direct command
    
    // Return appropriate message type
    return OutputMsg{Output: "result"}
}
```

### Error Handling
- Try multiple platform-specific methods before failing
- Accumulate errors from each attempt
- Return comprehensive error messages showing all attempts
- Include helpful suggestions (e.g., "run with sudo", "run as administrator")
- Use `CombinedOutput()` to capture both stdout and stderr

### Platform Detection
```go
// Use runtime.GOOS for platform-specific logic
if runtime.GOOS == "windows" {
    // Windows-specific code
} else {
    // Unix/Linux/macOS code
}
```

### Process Supervisor Support
LazyNginx detects and supports multiple process supervisors:

**Supported Supervisors:**
1. **systemd** - Standard Linux init system
2. **s6/s6-rc** - Lightweight supervision suite (s6-overlay compatible)
3. **supervisord** - Python-based process control system
4. **runit** - Unix init scheme with service supervision
5. **none** - Direct nginx control without supervisor

**Detection Order:**
1. Check if systemd is active and managing nginx
2. Check if nginx parent process is s6-supervise
3. Check for supervisord managing nginx
4. Check for runit service directories
5. Fall back to direct nginx commands

**Command Execution Pattern:**
Commands use supervisor-specific methods with graceful fallbacks:

For **systemd**:
- `systemctl start/stop/restart/reload nginx`

For **s6**:
- Try `s6-svc` command (e.g., `s6-svc -u /run/service/nginx`)
- Try `s6-rc` command (e.g., `s6-rc -u change nginx`)
- Write control bytes to `/proc/[s6-supervise-pid]/fd/5` (FIFO via /proc)
- Write to traditional FIFO path (e.g., `/run/service/nginx/supervise/control`)
- Fallback with warnings about s6 auto-restart behavior

For **supervisord**:
- `supervisorctl start/stop/restart nginx`

For **runit**:
- `sv start/stop/restart/reload nginx`

For **none** (direct control):
1. Windows: `net start/stop nginx`
2. Direct: `sudo nginx [-s stop/reload/quit]`
3. Fallback: Try without sudo

**S6 Control Bytes:**
- `u` = up/start
- `d` = down/stop
- `t` = terminate (kill and auto-restart)
- `h` = HUP signal (reload)
- `k` = kill (SIGKILL)

### Styling Guidelines
- All styles defined in `pkg/gui/styles.go` using lipgloss
- Use package-level style variables (e.g., `TitleStyle`, `SelectedStyle`, `ErrorStyle`)
- Color scheme: Purple primary (#7D56F4), Green success (#50FA7B), Red error (#FF5555)
- Border styles: Green (#2) for focused, Dim gray (#8) for unfocused
- Apply styles consistently across all views

### Type Definitions
```go
// Message types for Bubble Tea
type StatusMsg struct {
    Status string
}

type OutputMsg struct {
    Output string
}

type ConfigViewMsg struct {
    Output   string
    Path     string
    Type     string // "main" or "site"
    SiteName string
}
```

### File Operations
- Check multiple common paths for config files (Linux, Windows, macOS)
- Use `os.Stat()` to verify file existence before operations
- Use `os.ReadFile()` for reading config files
- Use `os.WriteFile()` with 0644 permissions for creating configs
- Use `os.Symlink()` for creating sites-enabled links

## Adding New Features

### Adding a New Command
1. Add function in `pkg/commands/commands.go` returning `tea.Msg`:
```go
func NewCommand() tea.Msg {
    cmd := exec.Command("command", "args")
    output, err := cmd.CombinedOutput()
    if err != nil {
        return OutputMsg{Output: "Error: " + string(output)}
    }
    return OutputMsg{Output: string(output)}
}
```

2. Add menu item to submenu in `pkg/app/app.go` `NewModel()` function
3. Wire up selection handler in `pkg/app/handlers.go` or `update.go`

### Navigation Keys
- Vim-style: `j` (down), `k` (up)
- Arrow keys: `↑`, `↓`, `←`, `→`
- Select: `Enter`
- Quit: `q` or `Ctrl+C`
- Tab: Switch between panels

## Best Practices

### DO
- Use Bubble Tea's command pattern for async operations
- Return `tea.Msg` types from command functions
- Check admin/root privileges before operations
- Try multiple platform-specific approaches
- Provide detailed error messages with all failed attempts
- Use lipgloss styles from `pkg/gui/styles.go`
- Test on all target platforms (Linux, Windows, macOS)

### DON'T
- Don't block the UI thread with long-running operations
- Don't assume single platform (always handle cross-platform)
- Don't ignore errors - accumulate and show all failures
- Don't hard-code file paths - check multiple common locations
- Don't add emojis without explicit user request
- Don't use panic - return error messages instead

## Platform-Specific Notes

### File Paths
- Linux: `/etc/nginx/`, `/var/log/nginx/`
- Windows: `C:\nginx\`, `C:\nginx\logs\`
- macOS/Unix: `/usr/local/nginx/`

### Permissions
- Linux/macOS: Requires `sudo` for service operations
- Windows: Requires "Run as Administrator"
- Check with: `commands.IsAdmin()` (returns bool)

### Status Check Cascade
1. Windows: `tasklist /FI "IMAGENAME eq nginx.exe"`
2. systemd: `systemctl is-active nginx`
3. Unix: `pgrep -x nginx`
4. Fallback: `ps aux | grep nginx`

## Testing Checklist for New Features
- [ ] Build succeeds: `go build -o lazynginx`
- [ ] No compilation errors
- [ ] Menu navigation works (j/k and arrow keys)
- [ ] Command executes successfully on Linux
- [ ] Command executes successfully on Windows
- [ ] Command executes successfully on macOS (if applicable)
- [ ] Error messages are helpful and actionable
- [ ] UI renders correctly at different terminal sizes
- [ ] Admin/sudo requirements are clear in messages
- [ ] No race conditions or deadlocks in Bubble Tea updates
