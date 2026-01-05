package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// Messages
type StatusMsg struct {
	Status string
}

type OutputMsg struct {
	Output string
}

// Commands
func CheckNginxStatus() tea.Msg {
	// Check if nginx binary exists
	_, err := exec.LookPath("nginx")
	if err != nil {
		return StatusMsg{Status: "Nginx binary not found in PATH"}
	}

	// Windows: Check if nginx process is running
	cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq nginx.exe")
	output, err := cmd.Output()
	if err == nil {
		if strings.Contains(strings.ToLower(string(output)), "nginx.exe") {
			return StatusMsg{Status: "Nginx is running"}
		}
	}

	// Unix/Linux: Try systemctl first (most common)
	cmd = exec.Command("systemctl", "is-active", "nginx")
	output, err = cmd.Output()
	if err == nil {
		status := strings.TrimSpace(string(output))
		if status == "active" {
			return StatusMsg{Status: "Nginx is running (systemd)"}
		} else if status == "inactive" {
			return StatusMsg{Status: "Nginx is not running (systemd)"}
		}
	}

	// Unix/Linux: Try pgrep
	cmd = exec.Command("pgrep", "-x", "nginx")
	err = cmd.Run()
	if err == nil {
		return StatusMsg{Status: "Nginx is running"}
	}

	// Try ps command as fallback
	cmd = exec.Command("ps", "aux")
	output, err = cmd.Output()
	if err == nil {
		if strings.Contains(strings.ToLower(string(output)), "nginx") {
			return StatusMsg{Status: "Nginx is running"}
		}
	}

	return StatusMsg{Status: "Nginx is not running"}
}

func StartNginx() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Try different methods based on OS
	// Windows
	cmd = exec.Command("net", "start", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx started successfully\n\n" + string(output)}
	}
	windowsErr := string(output)

	// Unix/Linux with systemd (with sudo)
	cmd = exec.Command("sudo", "systemctl", "start", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx started successfully\n\n" + string(output)}
	}
	systemdSudoErr := string(output)

	// Unix/Linux with systemd (without sudo)
	cmd = exec.Command("systemctl", "start", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx started successfully\n\n" + string(output)}
	}

	// Direct nginx command (with sudo)
	cmd = exec.Command("sudo", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx started successfully\n\n" + string(output)}
	}
	nginxSudoErr := string(output)

	// Direct nginx command (without sudo)
	cmd = exec.Command("nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx started successfully\n\n" + string(output)}
	}
	nginxErr := string(output)

	// All methods failed, show detailed error
	errorMsg := fmt.Sprintf("Failed to start nginx. Tried the following methods:\n\n")
	errorMsg += fmt.Sprintf("1. Windows (net start): %s\n", windowsErr)
	errorMsg += fmt.Sprintf("2. Systemd (sudo): %s\n", systemdSudoErr)
	errorMsg += fmt.Sprintf("3. Direct nginx (sudo): %s\n", nginxSudoErr)
	errorMsg += fmt.Sprintf("4. Direct nginx (no sudo): %s\n\n", nginxErr)
	errorMsg += "Note: You may need to:\n"
	errorMsg += "- Run lazynginx with sudo: sudo ./lazynginx\n"
	errorMsg += "- Or configure passwordless sudo for nginx commands\n"
	errorMsg += "- Or run as administrator on Windows"

	return OutputMsg{Output: errorMsg}
}

func StopNginx() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Windows
	cmd = exec.Command("net", "stop", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx stopped successfully\n\n" + string(output)}
	}
	windowsErr := string(output)

	// Unix/Linux with systemd
	cmd = exec.Command("sudo", "systemctl", "stop", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx stopped successfully\n\n" + string(output)}
	}
	systemdErr := string(output)

	// Try without sudo for systemd (if user has permissions)
	cmd = exec.Command("systemctl", "stop", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx stopped successfully\n\n" + string(output)}
	}

	// Direct nginx command with sudo
	cmd = exec.Command("sudo", "nginx", "-s", "stop")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx stopped successfully\n\n" + string(output)}
	}
	sudoNginxErr := string(output)

	// Try without sudo for nginx command
	cmd = exec.Command("nginx", "-s", "stop")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx stopped successfully\n\n" + string(output)}
	}
	nginxErr := string(output)

	// All methods failed, show detailed error
	errorMsg := fmt.Sprintf("Failed to stop nginx. Tried the following methods:\n\n")
	errorMsg += fmt.Sprintf("1. Windows (net stop): %s\n", windowsErr)
	errorMsg += fmt.Sprintf("2. Systemd (sudo): %s\n", systemdErr)
	errorMsg += fmt.Sprintf("3. Direct nginx (sudo): %s\n", sudoNginxErr)
	errorMsg += fmt.Sprintf("4. Direct nginx (no sudo): %s\n\n", nginxErr)
	errorMsg += "Note: You may need to:\n"
	errorMsg += "- Run lazynginx with sudo: sudo ./lazynginx\n"
	errorMsg += "- Or configure passwordless sudo for nginx commands\n"
	errorMsg += "- Or run as administrator on Windows"

	return OutputMsg{Output: errorMsg}
}

func RestartNginx() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Windows
	cmd = exec.Command("net", "stop", "nginx")
	cmd.Run()
	cmd = exec.Command("net", "start", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx restarted successfully\n\n" + string(output)}
	}
	windowsErr := string(output)

	// Unix/Linux with systemd (with sudo)
	cmd = exec.Command("sudo", "systemctl", "restart", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx restarted successfully\n\n" + string(output)}
	}
	systemdSudoErr := string(output)

	// Unix/Linux with systemd (without sudo)
	cmd = exec.Command("systemctl", "restart", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx restarted successfully\n\n" + string(output)}
	}

	// Direct nginx command (reload with sudo)
	cmd = exec.Command("sudo", "nginx", "-s", "reload")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx reloaded successfully\n\n" + string(output)}
	}
	nginxSudoErr := string(output)

	// Direct nginx command (reload without sudo)
	cmd = exec.Command("nginx", "-s", "reload")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx reloaded successfully\n\n" + string(output)}
	}
	nginxErr := string(output)

	// All methods failed, show detailed error
	errorMsg := fmt.Sprintf("Failed to restart nginx. Tried the following methods:\n\n")
	errorMsg += fmt.Sprintf("1. Windows (net restart): %s\n", windowsErr)
	errorMsg += fmt.Sprintf("2. Systemd (sudo): %s\n", systemdSudoErr)
	errorMsg += fmt.Sprintf("3. Direct nginx reload (sudo): %s\n", nginxSudoErr)
	errorMsg += fmt.Sprintf("4. Direct nginx reload (no sudo): %s\n\n", nginxErr)
	errorMsg += "Note: You may need to:\n"
	errorMsg += "- Run lazynginx with sudo: sudo ./lazynginx\n"
	errorMsg += "- Or configure passwordless sudo for nginx commands\n"
	errorMsg += "- Or run as administrator on Windows"

	return OutputMsg{Output: errorMsg}
}

func ReloadNginx() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Windows
	cmd = exec.Command("nginx", "-s", "reload")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx configuration reloaded successfully\n\n" + string(output)}
	}

	// Unix/Linux with systemd
	cmd = exec.Command("sudo", "systemctl", "reload", "nginx")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx configuration reloaded successfully\n\n" + string(output)}
	}

	// Direct nginx command
	cmd = exec.Command("sudo", "nginx", "-s", "reload")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Nginx configuration reloaded successfully\n\n" + string(output)}
	}

	return OutputMsg{Output: fmt.Sprintf("Failed to reload nginx:\n%s\n\nNote: You may need to run with sudo/administrator privileges", string(output))}
}

func TestNginxConfig() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Test configuration
	cmd = exec.Command("nginx", "-t")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Configuration test passed!\n\n" + string(output)}
	}

	// Try with sudo
	cmd = exec.Command("sudo", "nginx", "-t")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: "Configuration test passed!\n\n" + string(output)}
	}

	return OutputMsg{Output: fmt.Sprintf("Configuration test failed:\n%s", string(output))}
}

func ViewNginxConfig() tea.Msg {
	var cmd *exec.Cmd
	var output []byte
	var err error

	// Try common config locations
	configPaths := []string{
		"/etc/nginx/nginx.conf",
		"C:\\nginx\\conf\\nginx.conf",
		"/usr/local/nginx/conf/nginx.conf",
	}

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			content, err := os.ReadFile(path)
			if err == nil {
				return OutputMsg{Output: fmt.Sprintf("Nginx Configuration (%s):\n\n%s", path, string(content))}
			}
		}
	}

	// Try to get config path from nginx -V
	cmd = exec.Command("nginx", "-V")
	output, err = cmd.CombinedOutput()
	if err == nil {
		return OutputMsg{Output: fmt.Sprintf("Nginx version and configuration paths:\n\n%s\n\nNote: Use the --prefix path shown above to locate nginx.conf", string(output))}
	}

	return OutputMsg{Output: "Could not locate nginx configuration file.\n\nCommon locations:\n- /etc/nginx/nginx.conf (Linux)\n- C:\\nginx\\conf\\nginx.conf (Windows)\n- /usr/local/nginx/conf/nginx.conf (macOS/Unix)"}
}

func ViewErrorLogs() tea.Msg {
	var cmd *exec.Cmd
	var output []byte

	// Try common log locations
	logPaths := []string{
		"/var/log/nginx/error.log",
		"C:\\nginx\\logs\\error.log",
		"/usr/local/nginx/logs/error.log",
	}

	for _, path := range logPaths {
		if _, err := os.Stat(path); err == nil {
			cmd = exec.Command("tail", "-n", "50", path)
			output, err = cmd.CombinedOutput()
			if err == nil {
				return OutputMsg{Output: fmt.Sprintf("Last 50 lines of error log (%s):\n\n%s", path, string(output))}
			}
			// If tail doesn't work, try reading file directly
			content, err := os.ReadFile(path)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				start := 0
				if len(lines) > 50 {
					start = len(lines) - 50
				}
				return OutputMsg{Output: fmt.Sprintf("Last 50 lines of error log (%s):\n\n%s", path, strings.Join(lines[start:], "\n"))}
			}
		}
	}

	return OutputMsg{Output: "Could not locate nginx error log file.\n\nCommon locations:\n- /var/log/nginx/error.log (Linux)\n- C:\\nginx\\logs\\error.log (Windows)\n- /usr/local/nginx/logs/error.log (macOS/Unix)"}
}

func ViewAccessLogs() tea.Msg {
	var cmd *exec.Cmd
	var output []byte

	// Try common log locations
	logPaths := []string{
		"/var/log/nginx/access.log",
		"C:\\nginx\\logs\\access.log",
		"/usr/local/nginx/logs/access.log",
	}

	for _, path := range logPaths {
		if _, err := os.Stat(path); err == nil {
			cmd = exec.Command("tail", "-n", "50", path)
			output, err = cmd.CombinedOutput()
			if err == nil {
				return OutputMsg{Output: fmt.Sprintf("Last 50 lines of access log (%s):\n\n%s", path, string(output))}
			}
			// If tail doesn't work, try reading file directly
			content, err := os.ReadFile(path)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				start := 0
				if len(lines) > 50 {
					start = len(lines) - 50
				}
				return OutputMsg{Output: fmt.Sprintf("Last 50 lines of access log (%s):\n\n%s", path, strings.Join(lines[start:], "\n"))}
			}
		}
	}

	return OutputMsg{Output: "Could not locate nginx access log file.\n\nCommon locations:\n- /var/log/nginx/access.log (Linux)\n- C:\\nginx\\logs\\access.log (Windows)\n- /usr/local/nginx/logs/access.log (macOS/Unix)"}
}

// Sites commands
// ModelInterface defines the methods needed from the model
type ModelInterface interface {
	SetSubMenus(index int, items []string)
}

func LoadSites(m ModelInterface) tea.Cmd {
	return func() tea.Msg {
		var sites []string

		// Try common sites directories
		sitePaths := []string{
			"/etc/nginx/sites-available",
			"/etc/nginx/sites-enabled",
			"C:\\nginx\\conf\\sites-available",
			"/usr/local/nginx/sites-available",
		}

		for _, path := range sitePaths {
			entries, err := os.ReadDir(path)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() && entry.Name() != "default" {
						sites = append(sites, entry.Name())
					}
				}
				if len(sites) > 0 {
					// Prepend "Add site" to the sites list
					m.SetSubMenus(2, append([]string{"Add site"}, sites...))
					return StatusMsg{Status: fmt.Sprintf("Found %d sites", len(sites))}
				}
			}
		}

		// If no sites found, keep "Add site" option
		m.SetSubMenus(2, []string{"Add site", "No sites found"})

		return StatusMsg{Status: "No sites configured"}
	}
}

func ViewSiteConfig(siteName string) tea.Msg {
	if siteName == "No sites found" || siteName == "Loading sites..." {
		return OutputMsg{Output: "No site selected"}
	}

	// Try to find the site config file
	sitePaths := []string{
		"/etc/nginx/sites-available/" + siteName,
		"/etc/nginx/sites-enabled/" + siteName,
		"C:\\nginx\\conf\\sites-available\\" + siteName,
		"/usr/local/nginx/sites-available/" + siteName,
	}

	for _, path := range sitePaths {
		if _, err := os.Stat(path); err == nil {
			content, err := os.ReadFile(path)
			if err == nil {
				return OutputMsg{Output: fmt.Sprintf("Site Configuration: %s\n\nPath: %s\n\n%s", siteName, path, string(content))}
			}
		}
	}

	return OutputMsg{Output: fmt.Sprintf("Could not locate configuration file for site: %s\n\nSearched in:\n- /etc/nginx/sites-available/\n- /etc/nginx/sites-enabled/\n- C:\\nginx\\conf\\sites-available\\", siteName)}
}

func LoadReverseProxies(m ModelInterface) tea.Cmd {
	return func() tea.Msg {
		var proxies []string

		// Try to find nginx config and parse for proxy_pass directives
		configPaths := []string{
			"/etc/nginx/nginx.conf",
			"C:\\nginx\\conf\\nginx.conf",
			"/usr/local/nginx/conf/nginx.conf",
		}

		// Also check sites-available and sites-enabled directories
		siteDirs := []string{
			"/etc/nginx/sites-available",
			"/etc/nginx/sites-enabled",
			"/etc/nginx/conf.d",
			"C:\\nginx\\conf\\sites-available",
			"/usr/local/nginx/sites-available",
		}

		// Parse config files for proxy_pass directives
		var allConfigs []string

		// Add main config
		for _, path := range configPaths {
			if _, err := os.Stat(path); err == nil {
				allConfigs = append(allConfigs, path)
			}
		}

		// Add site configs
		for _, dir := range siteDirs {
			entries, err := os.ReadDir(dir)
			if err == nil {
				for _, entry := range entries {
					if !entry.IsDir() {
						allConfigs = append(allConfigs, dir+"/"+entry.Name())
					}
				}
			}
		}

		// Parse configs for proxy_pass
		proxyMap := make(map[string]bool)
		for _, configPath := range allConfigs {
			content, err := os.ReadFile(configPath)
			if err == nil {
				lines := strings.Split(string(content), "\n")
				for _, line := range lines {
					trimmed := strings.TrimSpace(line)
					if strings.Contains(trimmed, "proxy_pass") {
						// Extract the proxy target
						parts := strings.Fields(trimmed)
						for i, part := range parts {
							if part == "proxy_pass" && i+1 < len(parts) {
								target := strings.TrimRight(parts[i+1], ";")
								if !proxyMap[target] {
									proxyMap[target] = true
									proxies = append(proxies, target)
								}
							}
						}
					}
				}
			}
		}

		if len(proxies) > 0 {
			// Prepend "Add Reverse Proxy" to the proxies list
			m.SetSubMenus(3, append([]string{"Add Reverse Proxy"}, proxies...))
			return StatusMsg{Status: fmt.Sprintf("Found %d reverse proxies", len(proxies))}
		}

		// If no proxies found, keep "Add Reverse Proxy" option
		m.SetSubMenus(3, []string{"Add Reverse Proxy", "No reverse proxies found"})
		return StatusMsg{Status: "No reverse proxies configured"}
	}
}

func AddSite(siteType string, siteName string) tea.Msg {
	var configContent string
	var actualSiteName string

	if siteType == "Laravel" {
		actualSiteName = siteName
		if actualSiteName == "" {
			actualSiteName = "laravel-site"
		}
		configContent = fmt.Sprintf(`server {
    listen 80;
    server_name %s.local;
    root /var/www/%s/public;

    add_header X-Frame-Options "SAMEORIGIN";
    add_header X-Content-Type-Options "nosniff";

    index index.php;

    charset utf-8;

    location / {
        try_files $uri $uri/ /index.php?$query_string;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    error_page 404 /index.php;

    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }

    location ~ /\.(?!well-known).* {
        deny all;
    }
}`, actualSiteName, actualSiteName)
	} else if siteType == "Static" {
		actualSiteName = siteName
		if actualSiteName == "" {
			actualSiteName = "static-site"
		}
		configContent = fmt.Sprintf(`server {
    listen 80;
    server_name %s.local;
    root /var/www/%s;

    index index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }
}`, actualSiteName, actualSiteName)
	} else if siteType == "VanillaPHP" {
		actualSiteName = siteName
		if actualSiteName == "" {
			actualSiteName = "php-site"
		}
		configContent = fmt.Sprintf(`server {
    listen 80;
    server_name %s.local;
    root /var/www/%s;

    index index.php index.html index.htm;

    location / {
        try_files $uri $uri/ =404;
    }

    location = /favicon.ico { access_log off; log_not_found off; }
    location = /robots.txt  { access_log off; log_not_found off; }

    location ~ \.php$ {
        fastcgi_pass unix:/var/run/php/php8.1-fpm.sock;
        fastcgi_param SCRIPT_FILENAME $realpath_root$fastcgi_script_name;
        include fastcgi_params;
    }
}`, actualSiteName, actualSiteName)
	} else {
		actualSiteName = siteName
		if actualSiteName == "" {
			actualSiteName = "custom-site"
		}
		configContent = fmt.Sprintf(`server {
    listen 80;
    server_name %s.local;
    root /var/www/%s;

    index index.html index.htm index.php;

    location / {
        try_files $uri $uri/ =404;
    }
}`, actualSiteName, actualSiteName)
	}

	// Try to write to sites-available
	sitePaths := []string{
		"/etc/nginx/sites-available/" + actualSiteName,
		"C:\\nginx\\conf\\sites-available\\" + actualSiteName,
		"/usr/local/nginx/sites-available/" + actualSiteName,
	}

	for _, path := range sitePaths {
		dir := path[:strings.LastIndex(path, string(os.PathSeparator))]

		// Check if directory exists
		if _, err := os.Stat(dir); err == nil {
			// Try to write the file
			err := os.WriteFile(path, []byte(configContent), 0644)
			if err == nil {
				// Try to create symlink to sites-enabled
				enabledPath := strings.Replace(path, "sites-available", "sites-enabled", 1)
				os.Symlink(path, enabledPath)

				return OutputMsg{Output: fmt.Sprintf("Site '%s' created successfully!\n\nConfiguration file: %s\n\nType: %s\n\nNext steps:\n1. Create directory: /var/www/%s\n2. Reload nginx: sudo systemctl reload nginx\n3. Add to /etc/hosts: 127.0.0.1 %s.local", actualSiteName, path, siteType, actualSiteName, actualSiteName)}
			}
			return OutputMsg{Output: fmt.Sprintf("Failed to create site: %s\n\nYou may need sudo/administrator privileges", err.Error())}
		}
	}

	return OutputMsg{Output: "Could not locate nginx sites directory.\n\nPlease ensure Nginx is properly installed.\n\nCommon locations:\n- /etc/nginx/sites-available/ (Linux)\n- C:\\nginx\\conf\\sites-available\\ (Windows)"}
}

func AddProxy(proxyType string, proxyConfig string) tea.Msg {
	// Parse proxy configuration
	// Expected format: domain.com:port -> backend:port
	// For load balanced: domain.com:port -> backend1:port,backend2:port

	parts := strings.Split(proxyConfig, "->")
	if len(parts) != 2 {
		return OutputMsg{Output: "Invalid proxy configuration format.\n\nExpected format:\n- Simple: domain.com:port -> http://backend:port\n- Load Balanced: domain.com:port -> backend1:port,backend2:port"}
	}

	frontend := strings.TrimSpace(parts[0])
	backends := strings.TrimSpace(parts[1])

	if frontend == "" || backends == "" {
		return OutputMsg{Output: "Invalid proxy configuration. Frontend and backend cannot be empty."}
	}

	// Parse frontend (domain:port or just domain)
	var serverName, listenPort string
	if strings.Contains(frontend, ":") {
		frontendParts := strings.Split(frontend, ":")
		serverName = frontendParts[0]
		listenPort = frontendParts[1]
	} else {
		serverName = frontend
		listenPort = "80"
	}

	var configContent string
	var configName string

	if proxyType == "Simple" {
		// Simple reverse proxy
		configName = fmt.Sprintf("proxy-%s", serverName)

		// Ensure backend has http:// prefix
		backend := backends
		if !strings.HasPrefix(backend, "http://") && !strings.HasPrefix(backend, "https://") {
			backend = "http://" + backend
		}

		configContent = fmt.Sprintf(`# Simple Reverse Proxy for %s
server {
    listen %s;
    server_name %s;

    location / {
        proxy_pass %s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`, serverName, listenPort, serverName, backend)

	} else {
		// Load balanced reverse proxy
		configName = fmt.Sprintf("proxy-lb-%s", serverName)

		// Parse backend servers
		backendList := strings.Split(backends, ",")
		upstreamName := strings.ReplaceAll(serverName, ".", "_") + "_backend"

		// Build upstream block
		upstreamBlock := fmt.Sprintf("upstream %s {\n", upstreamName)
		for _, backend := range backendList {
			backend = strings.TrimSpace(backend)
			if !strings.HasPrefix(backend, "http://") && !strings.HasPrefix(backend, "https://") {
				backend = "http://" + backend
			}
			// Remove http:// prefix for upstream servers
			backend = strings.TrimPrefix(backend, "http://")
			backend = strings.TrimPrefix(backend, "https://")
			upstreamBlock += fmt.Sprintf("    server %s;\n", backend)
		}
		upstreamBlock += "}\n\n"

		configContent = fmt.Sprintf(`# Load Balanced Reverse Proxy for %s
%sserver {
    listen %s;
    server_name %s;

    location / {
        proxy_pass http://%s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}`, serverName, upstreamBlock, listenPort, serverName, upstreamName)
	}

	// Try to write to sites-available or conf.d
	configPaths := []string{
		"/etc/nginx/sites-available/" + configName,
		"/etc/nginx/conf.d/" + configName + ".conf",
		"C:\\nginx\\conf\\sites-available\\" + configName,
		"/usr/local/nginx/sites-available/" + configName,
	}

	for _, path := range configPaths {
		dir := path[:strings.LastIndex(path, string(os.PathSeparator))]

		// Check if directory exists
		if _, err := os.Stat(dir); err == nil {
			// Try to write the file
			err := os.WriteFile(path, []byte(configContent), 0644)
			if err == nil {
				// Try to create symlink to sites-enabled (if applicable)
				if strings.Contains(path, "sites-available") {
					enabledPath := strings.Replace(path, "sites-available", "sites-enabled", 1)
					os.Symlink(path, enabledPath)
				}

				return OutputMsg{Output: fmt.Sprintf("Reverse proxy '%s' created successfully!\n\nConfiguration file: %s\n\nType: %s\n\nFrontend: %s:%s\nBackend(s): %s\n\nNext steps:\n1. Test configuration: sudo nginx -t\n2. Reload nginx: sudo systemctl reload nginx\n3. Add to /etc/hosts if needed: 127.0.0.1 %s", configName, path, proxyType, serverName, listenPort, backends, serverName)}
			}
			return OutputMsg{Output: fmt.Sprintf("Failed to create reverse proxy: %s\n\nYou may need sudo/administrator privileges", err.Error())}
		}
	}

	return OutputMsg{Output: "Could not locate nginx configuration directory.\n\nPlease ensure Nginx is properly installed.\n\nCommon locations:\n- /etc/nginx/sites-available/ (Linux)\n- /etc/nginx/conf.d/ (Linux)\n- C:\\nginx\\conf\\sites-available\\ (Windows)"}
}

func DeleteSite(siteName string) tea.Msg {
	if siteName == "" || siteName == "Add site" || siteName == "Loading sites..." || siteName == "No sites found" {
		return OutputMsg{Output: "Invalid site name"}
	}

	// Try to find and delete the site config file
	sitePaths := []string{
		"/etc/nginx/sites-available/" + siteName,
		"C:\\nginx\\conf\\sites-available\\" + siteName,
		"/usr/local/nginx/sites-available/" + siteName,
	}

	var configPath string
	var foundPath bool

	// Find the config file
	for _, path := range sitePaths {
		if _, err := os.Stat(path); err == nil {
			configPath = path
			foundPath = true
			break
		}
	}

	if !foundPath {
		return OutputMsg{Output: fmt.Sprintf("Could not locate configuration file for site: %s\n\nSearched in:\n- /etc/nginx/sites-available/\n- C:\\nginx\\conf\\sites-available\\", siteName)}
	}

	// Try to remove symlink from sites-enabled first
	enabledPath := strings.Replace(configPath, "sites-available", "sites-enabled", 1)
	if _, err := os.Stat(enabledPath); err == nil {
		err := os.Remove(enabledPath)
		if err != nil {
			return OutputMsg{Output: fmt.Sprintf("Failed to remove symlink from sites-enabled: %s\n\nYou may need sudo/administrator privileges", err.Error())}
		}
	}

	// Delete the config file
	err := os.Remove(configPath)
	if err != nil {
		return OutputMsg{Output: fmt.Sprintf("Failed to delete site configuration: %s\n\nYou may need sudo/administrator privileges", err.Error())}
	}

	return OutputMsg{Output: fmt.Sprintf("Site '%s' deleted successfully!\n\nRemoved: %s\n\nNext steps:\n1. Reload nginx: sudo systemctl reload nginx\n2. Remove site directory if needed: /var/www/%s", siteName, configPath, siteName)}
}
