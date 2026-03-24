package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

func normalizeServerInput(input ServerInput) (MCPServer, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return MCPServer{}, errors.New("name is required")
	}

	serverType := strings.TrimSpace(input.Type)
	if serverType != "stdio" && serverType != "http" {
		return MCPServer{}, errors.New("type must be stdio or http")
	}

	server := MCPServer{
		ID:         input.ID,
		Name:       name,
		Type:       serverType,
		Command:    strings.TrimSpace(input.Command),
		Args:       parseLines(input.ArgsText),
		URL:        strings.TrimSpace(input.URL),
		Env:        parseEnvText(input.EnvText),
		WorkingDir: strings.TrimSpace(input.WorkingDir),
		Notes:      strings.TrimSpace(input.Notes),
	}

	if server.Type == "stdio" && server.Command == "" {
		return MCPServer{}, errors.New("command is required for stdio servers")
	}

	if server.Type == "http" && server.URL == "" {
		return MCPServer{}, errors.New("url is required for http servers")
	}

	return server, nil
}

func normalizeClientIDs(input []string) []string {
	seen := map[string]bool{}
	targets := make([]string, 0, len(input))
	for _, value := range input {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" || seen[trimmed] {
			continue
		}
		seen[trimmed] = true
		targets = append(targets, trimmed)
	}
	sort.Strings(targets)
	return targets
}

func setClientEnabled(enabledClients []string, clientID string, enabled bool) []string {
	current := toClientSet(enabledClients)
	if enabled {
		current[clientID] = true
	} else {
		delete(current, clientID)
	}

	next := make([]string, 0, len(current))
	for clientID := range current {
		next = append(next, clientID)
	}
	sort.Strings(next)
	return next
}

func toClientSet(clientIDs []string) map[string]bool {
	set := make(map[string]bool, len(clientIDs))
	for _, clientID := range normalizeClientIDs(clientIDs) {
		set[clientID] = true
	}
	return set
}

func parseLines(input string) []string {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	values := make([]string, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

func parseEnvText(input string) map[string]string {
	lines := strings.Split(strings.ReplaceAll(input, "\r\n", "\n"), "\n")
	values := map[string]string{}
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		parts := strings.SplitN(trimmed, "=", 2)
		key := strings.TrimSpace(parts[0])
		if key == "" {
			continue
		}

		value := ""
		if len(parts) == 2 {
			value = strings.TrimSpace(parts[1])
		}
		values[key] = value
	}
	return values
}

func nowRFC3339() string {
	return time.Now().Format(time.RFC3339)
}

func newServerID() string {
	return fmt.Sprintf("srv_%d", time.Now().UnixNano())
}
