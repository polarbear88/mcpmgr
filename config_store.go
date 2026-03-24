package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func NewConfigStore() *ConfigStore {
	return &ConfigStore{
		path: appConfigPath(),
	}
}

func (s *ConfigStore) Path() string {
	return s.path
}

func (s *ConfigStore) Load() (AppConfig, error) {
	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return AppConfig{}, err
	}

	data, err := os.ReadFile(s.path)
	if errors.Is(err, os.ErrNotExist) {
		return defaultConfig(), nil
	}
	if err != nil {
		return AppConfig{}, err
	}

	var config AppConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return AppConfig{}, fmt.Errorf("parse config: %w", err)
	}

	if config.Version == 0 {
		config.Version = configVersion
	}
	config.EnabledClients = normalizeClientIDs(config.EnabledClients)
	if config.ClientBackups == nil {
		config.ClientBackups = map[string]ClientBackup{}
	}
	if config.LastAppliedByClient == nil {
		config.LastAppliedByClient = map[string][]string{}
	}

	return config, nil
}

func (s *ConfigStore) Save(config AppConfig) error {
	if config.LastAppliedByClient == nil {
		config.LastAppliedByClient = map[string][]string{}
	}
	config.EnabledClients = normalizeClientIDs(config.EnabledClients)
	if config.ClientBackups == nil {
		config.ClientBackups = map[string]ClientBackup{}
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return err
	}

	return os.WriteFile(s.path, append(data, '\n'), 0o644)
}

func defaultConfig() AppConfig {
	return AppConfig{
		Version:             configVersion,
		UpdatedAt:           nowRFC3339(),
		Servers:             []MCPServer{},
		EnabledClients:      []string{},
		ClientBackups:       map[string]ClientBackup{},
		LastAppliedByClient: map[string][]string{},
	}
}

func cloneConfig(config AppConfig) AppConfig {
	data, err := json.Marshal(config)
	if err != nil {
		return config
	}

	var cloned AppConfig
	if err := json.Unmarshal(data, &cloned); err != nil {
		return config
	}

	if cloned.ClientBackups == nil {
		cloned.ClientBackups = map[string]ClientBackup{}
	}
	if cloned.LastAppliedByClient == nil {
		cloned.LastAppliedByClient = map[string][]string{}
	}
	return cloned
}

func appConfigPath() string {
	configRoot, err := os.UserConfigDir()
	if err != nil {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return "mcpmgr.json"
		}
		return filepath.Join(home, ".config", "mcpmgr", "config.json")
	}

	return filepath.Join(configRoot, "mcpmgr", "config.json")
}
