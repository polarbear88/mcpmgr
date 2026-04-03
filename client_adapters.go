package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

func (c clientBase) ID() string {
	return c.id
}

func (c clientBase) Status() ClientStatus {
	_, err := os.Stat(c.path)
	return ClientStatus{
		ID:        c.id,
		Name:      c.name,
		Path:      c.path,
		Format:    c.format,
		Detected:  err == nil,
		Supported: c.support,
		Notes:     c.notes,
	}
}

func (c clientBase) ManagedServerNames(servers []MCPServer) []string {
	names := make([]string, 0, len(servers))
	for _, server := range servers {
		names = append(names, server.Name)
	}
	sort.Strings(names)
	return names
}

func NewCodexAdapter() *CodexAdapter {
	home, _ := os.UserHomeDir()
	return &CodexAdapter{
		clientBase: clientBase{
			id:      "codex",
			name:    "Codex",
			path:    filepath.Join(home, ".codex", "config.toml"),
			format:  "toml",
			support: true,
			notes:   "Writes a managed [mcp_servers.*] block into ~/.codex/config.toml.",
		},
	}
}

func (a *CodexAdapter) Apply(servers []MCPServer, _ []string) ApplyTargetResult {
	selected := servers
	content, err := readOptionalText(a.path)
	if err != nil {
		return a.failure(err)
	}

	block := renderCodexManagedBlock(selected)
	next := replaceCodexMCPServers(content, block)
	if err := writeTextFile(a.path, next); err != nil {
		return a.failure(err)
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func (a *CodexAdapter) failure(err error) ApplyTargetResult {
	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    false,
		Message:    err.Error(),
	}
}

func NewClaudeCodeAdapter() *ClaudeCodeAdapter {
	home, _ := os.UserHomeDir()
	return &ClaudeCodeAdapter{
		clientBase: clientBase{
			id:      "claude_code",
			name:    "Claude Code",
			path:    filepath.Join(home, ".claude.json"),
			format:  "json",
			support: true,
			notes:   "Updates user-scoped MCP entries inside ~/.claude.json.",
		},
	}
}

func (a *ClaudeCodeAdapter) Apply(servers []MCPServer, previous []string) ApplyTargetResult {
	selected := servers
	if err := applyJSONConfig(a.path, selected, previous); err != nil {
		return ApplyTargetResult{
			ClientID:   a.id,
			ClientName: a.name,
			Path:       a.path,
			Success:    false,
			Message:    err.Error(),
		}
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func NewClaudeDesktopAdapter() *ClaudeDesktopAdapter {
	return &ClaudeDesktopAdapter{
		clientBase: clientBase{
			id:      "claude_desktop",
			name:    "Claude Desktop",
			path:    claudeDesktopPath(),
			format:  "json",
			support: true,
			notes:   claudeDesktopNotes(),
		},
	}
}

func (a *ClaudeDesktopAdapter) Apply(servers []MCPServer, previous []string) ApplyTargetResult {
	selected := servers
	if err := applyJSONConfig(a.path, selected, previous); err != nil {
		return ApplyTargetResult{
			ClientID:   a.id,
			ClientName: a.name,
			Path:       a.path,
			Success:    false,
			Message:    err.Error(),
		}
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func NewCopilotCLIAdapter() *CopilotCLIAdapter {
	home, _ := os.UserHomeDir()
	return &CopilotCLIAdapter{
		clientBase: clientBase{
			id:      "copilot_cli",
			name:    "Copilot CLI",
			path:    filepath.Join(home, ".copilot", "mcp-config.json"),
			format:  "json",
			support: true,
			notes:   "Writes GitHub Copilot CLI MCP config to ~/.copilot/mcp-config.json.",
		},
	}
}

func NewVSCodeAdapter() *VSCodeAdapter {
	return &VSCodeAdapter{
		clientBase: clientBase{
			id:      "vscode",
			name:    "Visual Studio Code",
			path:    vscodeMCPPath(),
			format:  "json",
			support: true,
			notes:   "Writes VS Code MCP servers into the user profile mcp.json file.",
		},
	}
}

func (a *CopilotCLIAdapter) Apply(servers []MCPServer, previous []string) ApplyTargetResult {
	selected := servers
	if err := applyCopilotCLIConfig(a.path, selected, previous); err != nil {
		return ApplyTargetResult{
			ClientID:   a.id,
			ClientName: a.name,
			Path:       a.path,
			Success:    false,
			Message:    err.Error(),
		}
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func (a *VSCodeAdapter) Apply(servers []MCPServer, previous []string) ApplyTargetResult {
	selected := servers
	if err := applyVSCodeConfig(a.path, selected, previous); err != nil {
		return ApplyTargetResult{
			ClientID:   a.id,
			ClientName: a.name,
			Path:       a.path,
			Success:    false,
			Message:    err.Error(),
		}
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func claudeDesktopPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Claude", "claude_desktop_config.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return filepath.Join(home, "AppData", "Roaming", "Claude", "claude_desktop_config.json")
		}
		return filepath.Join(appData, "Claude", "claude_desktop_config.json")
	default:
		return filepath.Join(home, ".config", "Claude", "claude_desktop_config.json")
	}
}

func vscodeMCPPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Application Support", "Code", "User", "mcp.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return filepath.Join(home, "AppData", "Roaming", "Code", "User", "mcp.json")
		}
		return filepath.Join(appData, "Code", "User", "mcp.json")
	default:
		return filepath.Join(home, ".config", "Code", "User", "mcp.json")
	}
}

func claudeDesktopNotes() string {
	if runtime.GOOS == "linux" {
		return "Uses ~/.config/Claude/claude_desktop_config.json as the Linux target path."
	}
	return "Updates claude_desktop_config.json while preserving unmanaged MCP entries."
}

func NewOpenCodeAdapter() *OpenCodeAdapter {
	return &OpenCodeAdapter{
		clientBase: clientBase{
			id:      "opencode",
			name:    "OpenCode",
			path:    openCodeConfigPath(),
			format:  "json",
			support: true,
			notes:   "Writes OpenCode MCP servers into the top-level mcp object in opencode.json.",
		},
	}
}

func (a *OpenCodeAdapter) Apply(servers []MCPServer, previous []string) ApplyTargetResult {
	selected := servers
	if err := applyOpenCodeConfig(a.path, selected, previous); err != nil {
		return ApplyTargetResult{
			ClientID:   a.id,
			ClientName: a.name,
			Path:       a.path,
			Success:    false,
			Message:    err.Error(),
		}
	}

	return ApplyTargetResult{
		ClientID:   a.id,
		ClientName: a.name,
		Path:       a.path,
		Success:    true,
		Message:    fmt.Sprintf("Applied %d server(s).", len(selected)),
	}
}

func openCodeConfigPath() string {
	home, _ := os.UserHomeDir()
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, ".config", "opencode", "opencode.json")
	case "windows":
		appData := os.Getenv("APPDATA")
		if appData == "" {
			return filepath.Join(home, "AppData", "Roaming", "opencode", "opencode.json")
		}
		return filepath.Join(appData, "opencode", "opencode.json")
	default:
		return filepath.Join(home, ".config", "opencode", "opencode.json")
	}
}
