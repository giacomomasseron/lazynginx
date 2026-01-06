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

## Layout System

The project uses [lazycore's boxlayout](https://github.com/jesseduffield/lazycore) for managing panel layouts. This provides:
- Flexible box-based layouts with ROW and COLUMN directions
- Dynamic weight-based sizing (e.g., 1:1:2 ratio for the three panels)
- Static size support for fixed-width/height boxes
- Automatic space distribution and responsive layout
- Full-screen layout that adapts to terminal dimensions

The main view uses a horizontal (COLUMN) layout with three panels:
- Main Menu (weight: 1) - 25% width
- Sub Menu (weight: 1) - 25% width  
- Details Panel (weight: 2) - 50% width

## Data Flow

```
User Input → Update() → Command Execution → Message → Update() → View()
```

1. User presses key (j/k/Enter)
2. `Update()` modifies cursor or calls `handleSelection()`
3. Command runs asynchronously, returns `tea.Msg`
4. `Update()` receives message, updates model state
5. `View()` renders current state using lazycore boxlayout