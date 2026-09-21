# LazyNginx 🚀

A beautiful terminal-based Nginx manager built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

[![Release](https://img.shields.io/github/v/release/giacomomasseron/lazynginx)](https://github.com/giacomomasseron/lazynginx/releases/latest) [![GitHub Releases](https://img.shields.io/github/downloads/giacomomasseron/lazynginx/total)](https://github.com/giacomomasseron/lazynginx/releases) [![Lint](https://github.com/giacomomasseron/lazynginx/actions/workflows/lint.yml/badge.svg)](https://github.com/giacomomasseron/lazynginx/actions/workflows/lint.yml)

![lazynginx screen](docs/screens/lazynginx_screen_1.png?raw=true "laxynging screen")

![lazynginx screen](docs/screens/lazynginx_screen_2.png?raw=true "laxynging screen 2")

## Features

- ✅ Check Nginx status
- 🚀 Start/Stop/Restart Nginx
- 🔄 Reload configuration
- ✅ Test configuration
- 🌐 Manage sites: list, add, view and delete
- 🔀 Manage reverse proxies and load balancers
- 📄 View and edit configuration files in your `$EDITOR`
- 📊 View error logs
- 📈 View access logs
- 🎨 Beautiful terminal UI

## Installation

### Binary Releases (Recommended)

For Windows, macOS (10.12+) or Linux, you can download a binary release from the [releases page](https://github.com/giacomomasseron/lazynginx/releases).

**Linux / macOS:**
```bash
# Download the latest release for your platform
# For example, Linux x86_64:
LAZYNGINX_VERSION=$(curl -s "https://api.github.com/repos/giacomomasseron/lazynginx/releases/latest" | grep -Po '"tag_name": *"v\K[^"]*')
curl -Lo lazynginx.tar.gz "https://github.com/giacomomasseron/lazynginx/releases/download/v${LAZYNGINX_VERSION}/lazynginx_${LAZYNGINX_VERSION}_linux_amd64.tar.gz"
tar xf lazynginx.tar.gz lazynginx
sudo install lazynginx -D -t /usr/local/bin/
```

**Windows:**
Download the `.zip` file from the [releases page](https://github.com/giacomomasseron/lazynginx/releases), extract it, and run `lazynginx.exe`.

### Homebrew (macOS and Linux)

```bash
# Coming soon
brew install lazynginx
```

### Go

```bash
go install github.com/giacomomasseron/lazynginx@latest
```

**Note:** If you get an error claiming that lazynginx cannot be found or is not defined, you may need to add `~/go/bin` to your `$PATH` (macOS/Linux), or `%HOME%\go\bin` (Windows).

### Build from Source

**Prerequisites:**
- Go 1.21 or later
- Nginx installed on your system

```bash
git clone https://github.com/giacomomasseron/lazynginx.git
cd lazynginx
go mod download
go build -o lazynginx
```

You can also use `go run main.go` to compile and run in one go.

### Running LazyNginx

```bash
# Linux/macOS (may require sudo for service operations)
sudo ./lazynginx

# Or if installed globally
sudo lazynginx
```

On Windows, run as Administrator:
```powershell
lazynginx.exe
```

## Usage

### Navigation

- `↑` / `↓` or `k` / `j`: Navigate menu
- `→` / `l` / `Tab`: Move to the next panel
- `←` / `h`: Move back to the previous panel
- `Enter`: Select option
- `e`: Edit the selected configuration file
- `d`: Delete the selected site
- `q` or `Ctrl+C`: Quit application

### Available Commands

The menu is grouped into categories:

1. **Status & Monitoring** - Check Status, Test Configuration
2. **Service Control** - Start, Stop, Restart, Reload Configuration
3. **Sites** - Add a site, then view, edit or delete the existing ones
4. **Reverse Proxies** - Add Reverse Proxy, Add Load Balancer, and view the existing ones
5. **Configuration** - Display and edit the Nginx configuration file
6. **Logs** - View Error Log, View Access Log (last 50 lines)
7. **Quit** - Exit the application

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

### Linting

This project uses [golangci-lint](https://golangci-lint.run/). The rules live in
[`.golangci.yml`](.golangci.yml) and run on every push and pull request.

Install it ([other methods](https://golangci-lint.run/docs/welcome/install/)):

```bash
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.13.2
```

Then run it from the repository root:

```bash
golangci-lint run ./...

# Apply the fixes it can make automatically (including formatting)
golangci-lint run --fix ./...
```

Please make sure it reports no issues before opening a Pull Request.

