# LazyNginx - Program Functions

## Overview
LazyNginx is a terminal-based Nginx management tool that provides an interactive menu interface for common Nginx operations without requiring command memorization.

## Menu Voices

### Status & Monitoring
- **Check Status** - Verifies if Nginx is running using multiple detection methods (process checks, systemctl, tasklist)
- **Test Configuration** - Validates nginx.conf syntax without applying changes (nginx -t)

### Service Control

- **Start** - Starts the Nginx service using platform-appropriate commands (systemctl, net start, or direct nginx binary)
- **Stop** - Stops the running Nginx service gracefully
- **Restart** - Performs a full restart of the Nginx service
- **Reload Configuration** - Reloads Nginx configuration without dropping connections (nginx -s reload)

### Sites

This menu voice shows the sites list of nginx in the sub-menu box.  
When you choose a site in the list, the third box shows the detail of the config file of the site.

- **Add site** - This function open a modal to add new nginx site, with some choices: Laravel, Custom.  
It you click on "Custom", another modal opens with text input.

### Reverse Proxies

This menu voice reads the nginx config file and lists all reverse proxies defined in it. And it shows them in the second box on the right.

### Configuration

This menu voice automatically shows the config filein the third box on the right.

### Logs
- **View Error Log** - Shows recent Nginx error log entries
- **View Access Log** - Displays recent access log entries

### Core Functions

### Navigation
- **Interactive Menu** - Cursor-based navigation using arrow keys or Vim-style (j/k) controls
- **Output Viewing** - Dedicated mode for viewing command results with ability to return to menu
- **Quit** - Exit the application

## Platform Support
All functions automatically adapt to the host operating system:
- **Windows** - Uses `net start/stop` and checks `C:\nginx\`
- **Linux** - Prefers systemd commands, checks `/etc/nginx/`
- **macOS/Unix** - Uses direct nginx commands, checks `/usr/local/nginx/`

## User Experience Features
- Full-screen terminal interface with clean styling
- Color-coded status messages (green for success, red for errors)
- No command-line arguments needed - all operations via interactive menu
- Sudo/admin handling automatic where required