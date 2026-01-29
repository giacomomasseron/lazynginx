# LazyNginx 🚀

A beautiful terminal-based Nginx manager built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

[![GitHub Releases](https://img.shields.io/github/downloads/giacomomasseron/lazynginx/total)](https://github.com/giacomomasseron/lazynginx/releases) [![Go Report Card](https://goreportcard.com/badge/github.com/giacomomasseron/lazynginx)](https://goreportcard.com/report/github.com/giacomomasseron/lazynginx)

![lazynginx screen](docs/screens/lazynginx_screen_1.png?raw=true "laxynging screen")

![lazynginx screen](docs/screens/lazynginx_screen_2.png?raw=true "laxynging screen 2")

## Features

- ✅ Check Nginx status
- 🚀 Start/Stop/Restart Nginx
- 🔄 Reload configuration
- ✅ Test configuration
- 📄 View configuration file
- 📊 View error logs
- 📈 View access logs
- 🎨 Beautiful terminal UI

## Installation

### Prerequisites

- Go 1.21 or later
- Nginx installed on your system

### Build from Source

```bash
git clone <repository-url>
cd lazynginx
go mod download
go build -o lazynginx
```

### Run

```bash
./lazynginx
```

Or on Windows:
```bash
lazynginx.exe
```

## Usage

### Navigation

- `↑` / `↓` or `k` / `j`: Navigate menu
- `Enter`: Select option
- `q` or `Ctrl+C`: Quit application

### Available Commands

1. **Check Status** - Check if Nginx is running
2. **Start Nginx** - Start the Nginx service
3. **Stop Nginx** - Stop the Nginx service
4. **Restart Nginx** - Restart the Nginx service
5. **Reload Configuration** - Reload Nginx configuration without downtime
6. **Test Configuration** - Test Nginx configuration for syntax errors
7. **View Configuration** - Display Nginx configuration file
8. **View Error Logs** - Show last 50 lines of error log
9. **View Access Logs** - Show last 50 lines of access log
10. **Quit** - Exit the application

## Platform Support

The application automatically detects your platform and uses the appropriate commands:

- **Linux**: Uses `systemctl` when available, falls back to direct `nginx` commands
- **Windows**: Uses `net start/stop` commands
- **macOS/Unix**: Uses direct `nginx` commands

## Permissions

Some operations (start, stop, restart, reload) may require administrator/sudo privileges depending on your system configuration.

### Linux/macOS
```bash
sudo ./lazynginx
```

### Windows
Run as Administrator

## Configuration

The application automatically searches for Nginx in common locations:

- `/etc/nginx/nginx.conf` (Linux)
- `C:\nginx\conf\nginx.conf` (Windows)
- `/usr/local/nginx/conf/nginx.conf` (macOS/Unix)

## Logs

The application looks for logs in:

- `/var/log/nginx/` (Linux)
- `C:\nginx\logs\` (Windows)
- `/usr/local/nginx/logs/` (macOS/Unix)

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

