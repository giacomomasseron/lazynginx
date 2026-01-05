# LazyNginx Architecture

## Project Structure

```
lazynginx/
├── .github/
│   └── copilot-instructions.md    # AI development guide
├── docs/
│   └── architecture.md            # This file - project structure documentation
├── pkg/                           # Folder that contains all package files - initializes Bubble Tea TUI
├── pkg/app/                       # Folder for app.go file, that contains the main app of the project
├── pkg/commands/                  # Folder that contains go file with commands
├── pkg/utils/                     # Folder that contains go file with utils functions
├── pkg/gui/                       # Folder that contains go file for styles
```

## Data Flow

```
User Input → Update() → Command Execution → Message → Update() → View()
```

1. User presses key (j/k/Enter)
2. `Update()` modifies cursor or calls `handleSelection()`
3. Command runs asynchronously, returns `tea.Msg`
4. `Update()` receives message, updates model state
5. `View()` renders current state