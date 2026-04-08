package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

func (s *AppService) applyEnabledClients(config *AppConfig) (ApplyResult, error) {
	results := make([]ApplyTargetResult, 0, len(config.EnabledClients))
	enabledClients := toClientSet(config.EnabledClients)
	var failures []string

	for _, adapter := range s.adapters {
		if !enabledClients[adapter.ID()] {
			continue
		}

		result := adapter.Apply(config.Servers, config.LastAppliedByClient[adapter.ID()])
		results = append(results, result)

		if result.Success {
			config.LastAppliedByClient[adapter.ID()] = adapter.ManagedServerNames(config.Servers)
			continue
		}

		failures = append(failures, fmt.Sprintf("%s: %s", adapter.Status().Name, result.Message))
	}

	config.UpdatedAt = nowRFC3339()
	if err := s.store.Save(*config); err != nil {
		return ApplyResult{}, err
	}

	result := ApplyResult{
		AppliedAt:      nowRFC3339(),
		AppliedServers: len(config.Servers),
		Results:        results,
	}

	if len(failures) > 0 {
		return result, errors.New(strings.Join(failures, "\n"))
	}

	return result, nil
}

func renderCodexManagedBlock(servers []MCPServer) string {
	if len(servers) == 0 {
		return ""
	}

	var b strings.Builder
	for _, server := range servers {
		b.WriteString(fmt.Sprintf("[mcp_servers.%s]\n", tomlBareKey(server.Name)))
		if server.Type == "http" {
			b.WriteString("enabled = true\n")
			b.WriteString(fmt.Sprintf("url = %s\n", tomlString(server.URL)))
		} else {
			b.WriteString(fmt.Sprintf("command = %s\n", tomlString(server.Command)))
			if len(server.Args) > 0 {
				b.WriteString(fmt.Sprintf("args = %s\n", tomlArray(server.Args)))
			}
		}

		if server.WorkingDir != "" {
			b.WriteString(fmt.Sprintf("cwd = %s\n", tomlString(server.WorkingDir)))
		}

		if len(server.Env) > 0 {
			sectionName := "env"
			if server.Type == "http" {
				sectionName = "http_headers"
			}
			b.WriteString(fmt.Sprintf("[mcp_servers.%s.%s]\n", tomlBareKey(server.Name), sectionName))
			keys := make([]string, 0, len(server.Env))
			for key := range server.Env {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			for _, key := range keys {
				b.WriteString(fmt.Sprintf("%s = %s\n", tomlBareKey(key), tomlString(server.Env[key])))
			}
		}

		b.WriteString("\n")
	}

	return strings.TrimSpace(b.String())
}

func replaceManagedBlock(content, marker, replacement string) string {
	startMarker := "# BEGIN " + marker
	endMarker := "# END " + marker

	start := strings.Index(content, startMarker)
	if start >= 0 {
		end := strings.Index(content[start:], endMarker)
		if end >= 0 {
			endIndex := start + end + len(endMarker)
			if endIndex < len(content) && content[endIndex] == '\n' {
				endIndex++
			}
			content = strings.TrimSpace(content[:start] + content[endIndex:])
		}
	}

	if strings.TrimSpace(replacement) == "" {
		if content == "" {
			return ""
		}
		return strings.TrimSpace(content) + "\n"
	}

	block := startMarker + "\n" + replacement + "\n" + endMarker
	if strings.TrimSpace(content) == "" {
		return block + "\n"
	}

	return strings.TrimSpace(content) + "\n\n" + block + "\n"
}

func replaceCodexMCPServers(content, replacement string) string {
	base, _ := splitCodexMCPServers(content)
	return joinCodexMCPServers(base, replacement)
}

func splitCodexMCPServers(content string) (string, string) {
	cleaned := replaceManagedBlock(content, "MCPMGR MANAGED MCP SERVERS", "")
	lines := strings.Split(strings.ReplaceAll(cleaned, "\r\n", "\n"), "\n")
	baseLines := make([]string, 0, len(lines))
	mcpLines := make([]string, 0, len(lines))
	inMCPSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]") {
			inMCPSection = strings.HasPrefix(trimmed, "[mcp_servers")
		}

		if inMCPSection {
			mcpLines = append(mcpLines, line)
		} else {
			baseLines = append(baseLines, line)
		}
	}

	return strings.TrimSpace(strings.Join(baseLines, "\n")), strings.TrimSpace(strings.Join(mcpLines, "\n"))
}

func joinCodexMCPServers(baseContent, mcpContent string) string {
	base := strings.TrimSpace(baseContent)
	mcp := strings.TrimSpace(mcpContent)

	if mcp == "" {
		if base == "" {
			return ""
		}
		return base + "\n"
	}

	if base == "" {
		return mcp + "\n"
	}

	return base + "\n\n" + mcp + "\n"
}

func applyJSONConfig(path string, servers []MCPServer, previous []string) error {
	root := map[string]any{}

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	_ = previous
	mcpServers := map[string]any{}

	for _, server := range servers {
		mcpServers[server.Name] = jsonConfigForServer(server)
	}

	root["mcpServers"] = mcpServers

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	return writeTextFile(path, string(output)+"\n")
}

func applyCopilotCLIConfig(path string, servers []MCPServer, previous []string) error {
	root := map[string]any{}

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	_ = previous
	mcpServers := map[string]any{}

	for _, server := range servers {
		mcpServers[server.Name] = copilotCLIConfigForServer(server)
	}

	root["mcpServers"] = mcpServers

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	return writeTextFile(path, string(output)+"\n")
}

func applyVSCodeConfig(path string, servers []MCPServer, previous []string) error {
	root := map[string]any{}

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	_ = previous
	serverConfigs := map[string]any{}
	for _, server := range servers {
		serverConfigs[server.Name] = jsonConfigForServer(server)
	}

	root["servers"] = serverConfigs

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	return writeTextFile(path, string(output)+"\n")
}

func jsonConfigForServer(server MCPServer) map[string]any {
	config := map[string]any{}

	if server.Type == "http" {
		config["type"] = "http"
		config["url"] = server.URL
		if len(server.Env) > 0 {
			config["headers"] = server.Env
		}
	} else {
		config["type"] = "stdio"
		config["command"] = server.Command
		if len(server.Args) > 0 {
			config["args"] = server.Args
		}
		if len(server.Env) > 0 {
			config["env"] = server.Env
		}
	}
	if server.WorkingDir != "" {
		config["cwd"] = server.WorkingDir
	}

	return config
}

func copilotCLIConfigForServer(server MCPServer) map[string]any {
	config := map[string]any{
		"tools": []string{"*"},
	}

	if server.Type == "http" {
		config["type"] = "http"
		config["url"] = server.URL
		if len(server.Env) > 0 {
			config["headers"] = server.Env
		}
		return config
	}

	config["type"] = "stdio"
	config["command"] = server.Command
	if len(server.Args) > 0 {
		config["args"] = server.Args
	}
	if len(server.Env) > 0 {
		config["env"] = server.Env
	}

	return config
}

func readOptionalText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func captureClientBackup(clientID, path string) (ClientBackup, error) {
	switch clientID {
	case "codex":
		content, err := readOptionalText(path)
		if err != nil {
			return ClientBackup{}, err
		}
		_, mcpContent := splitCodexMCPServers(content)
		return ClientBackup{
			Kind:      "codex_mcp",
			Existed:   strings.TrimSpace(mcpContent) != "",
			Content:   strings.TrimSpace(mcpContent),
			UpdatedAt: nowRFC3339(),
		}, nil
	case "vscode":
		return captureJSONSectionBackup(path, "servers")
	case "opencode":
		return captureJSONSectionBackup(path, "mcp")
	default:
		return captureJSONSectionBackup(path, "mcpServers")
	}
}

func restoreClientBackup(clientID, path string, backup ClientBackup) error {
	switch clientID {
	case "codex":
		content, err := readOptionalText(path)
		if err != nil {
			return err
		}
		baseContent, _ := splitCodexMCPServers(content)
		next := joinCodexMCPServers(baseContent, backup.Content)
		if strings.TrimSpace(next) == "" {
			if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				return err
			}
			return nil
		}
		return writeTextFile(path, next)
	case "vscode":
		return restoreJSONSectionBackup(path, "servers", backup)
	case "opencode":
		return restoreJSONSectionBackup(path, "mcp", backup)
	default:
		return restoreJSONSectionBackup(path, "mcpServers", backup)
	}
}

func captureJSONSectionBackup(path, key string) (ClientBackup, error) {
	root := map[string]any{}
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return ClientBackup{
			Kind:      "json_" + key,
			Existed:   false,
			Content:   "",
			UpdatedAt: nowRFC3339(),
		}, nil
	}
	if err != nil {
		return ClientBackup{}, err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return ClientBackup{}, fmt.Errorf("parse %s: %w", path, err)
		}
	}

	node, ok := root[key]
	if !ok {
		return ClientBackup{
			Kind:      "json_" + key,
			Existed:   false,
			Content:   "",
			UpdatedAt: nowRFC3339(),
		}, nil
	}

	serialized, err := json.MarshalIndent(node, "", "  ")
	if err != nil {
		return ClientBackup{}, err
	}

	return ClientBackup{
		Kind:      "json_" + key,
		Existed:   true,
		Content:   string(serialized),
		UpdatedAt: nowRFC3339(),
	}, nil
}

func restoreJSONSectionBackup(path, key string, backup ClientBackup) error {
	root := map[string]any{}
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	if backup.Existed {
		var node any
		if err := json.Unmarshal([]byte(backup.Content), &node); err != nil {
			return err
		}
		root[key] = node
	} else {
		delete(root, key)
	}

	if len(root) == 0 {
		if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
		return nil
	}

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return writeTextFile(path, string(output)+"\n")
}

func writeTextFile(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func tomlString(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"\"", "\\\"",
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	)
	return "\"" + replacer.Replace(value) + "\""
}

func tomlArray(values []string) string {
	escaped := make([]string, 0, len(values))
	for _, value := range values {
		escaped = append(escaped, tomlString(value))
	}
	return "[" + strings.Join(escaped, ", ") + "]"
}

func tomlBareKey(value string) string {
	if value == "" {
		return "\"\""
	}

	for _, char := range value {
		if !(char == '-' || char == '_' || char >= 'a' && char <= 'z' || char >= 'A' && char <= 'Z' || char >= '0' && char <= '9') {
			return tomlString(value)
		}
	}

	return value
}

func applyOpenCodeConfig(path string, servers []MCPServer, previous []string) error {
	root := map[string]any{}

	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if len(data) > 0 {
		if err := json.Unmarshal(data, &root); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	}

	_ = previous
	mcpServers := map[string]any{}
	for _, server := range servers {
		mcpServers[server.Name] = openCodeConfigForServer(server)
	}

	root["mcp"] = mcpServers

	output, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	return writeTextFile(path, string(output)+"\n")
}

func openCodeConfigForServer(server MCPServer) map[string]any {
	config := map[string]any{}

	if server.Type == "http" {
		config["type"] = "remote"
		config["url"] = server.URL
		if len(server.Env) > 0 {
			config["headers"] = server.Env
		}
	} else {
		config["type"] = "local"
		command := []string{server.Command}
		if len(server.Args) > 0 {
			command = append(command, server.Args...)
		}
		config["command"] = command
		if len(server.Env) > 0 {
			config["environment"] = server.Env
		}
	}
	if server.WorkingDir != "" {
		config["cwd"] = server.WorkingDir
	}

	return config
}
