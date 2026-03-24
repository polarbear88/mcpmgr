package main

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
)

func NewAppService() *AppService {
	return &AppService{
		store: NewConfigStore(),
		adapters: []ClientAdapter{
			NewCodexAdapter(),
			NewClaudeCodeAdapter(),
			NewClaudeDesktopAdapter(),
			NewCopilotCLIAdapter(),
			NewVSCodeAdapter(),
		},
	}
}

func (s *AppService) GetState() (AppState, error) {
	config, err := s.store.Load()
	if err != nil {
		return AppState{}, err
	}

	return AppState{
		ConfigPath: s.store.Path(),
		Servers:    config.Servers,
		Clients:    s.clientStatuses(config),
	}, nil
}

func (s *AppService) SaveServer(input ServerInput) (AppState, error) {
	config, err := s.store.Load()
	if err != nil {
		return AppState{}, err
	}
	previousConfig := cloneConfig(config)

	server, err := normalizeServerInput(input)
	if err != nil {
		return AppState{}, err
	}

	index := -1
	for i, existing := range config.Servers {
		if existing.Name == server.Name && existing.ID != server.ID {
			return AppState{}, fmt.Errorf("server name %q already exists", server.Name)
		}
		if existing.ID == server.ID {
			index = i
		}
	}

	if input.ID == "" {
		server.ID = newServerID()
	}

	server.UpdatedAt = nowRFC3339()

	if index >= 0 {
		config.Servers[index] = server
	} else {
		config.Servers = append(config.Servers, server)
	}

	sort.Slice(config.Servers, func(i, j int) bool {
		return strings.ToLower(config.Servers[i].Name) < strings.ToLower(config.Servers[j].Name)
	})

	config.UpdatedAt = nowRFC3339()

	if _, err := s.applyEnabledClients(&config); err != nil {
		if rollbackErr := s.rollbackConfig(previousConfig); rollbackErr != nil {
			return AppState{}, fmt.Errorf("%w; rollback failed: %v", err, rollbackErr)
		}
		return AppState{}, err
	}

	return s.GetState()
}

func (s *AppService) DeleteServer(id string) (AppState, error) {
	config, err := s.store.Load()
	if err != nil {
		return AppState{}, err
	}
	previousConfig := cloneConfig(config)

	nextServers := make([]MCPServer, 0, len(config.Servers))
	found := false
	for _, server := range config.Servers {
		if server.ID == id {
			found = true
			continue
		}
		nextServers = append(nextServers, server)
	}

	if !found {
		return AppState{}, errors.New("server not found")
	}

	config.Servers = nextServers
	config.UpdatedAt = nowRFC3339()

	if _, err := s.applyEnabledClients(&config); err != nil {
		if rollbackErr := s.rollbackConfig(previousConfig); rollbackErr != nil {
			return AppState{}, fmt.Errorf("%w; rollback failed: %v", err, rollbackErr)
		}
		return AppState{}, err
	}

	return s.GetState()
}

func (s *AppService) ApplyToAllClients() (ApplyResult, error) {
	config, err := s.store.Load()
	if err != nil {
		return ApplyResult{}, err
	}

	return s.applyEnabledClients(&config)
}

func (s *AppService) EnableClient(clientID string) (AppState, error) {
	config, err := s.store.Load()
	if err != nil {
		return AppState{}, err
	}

	var adapter ClientAdapter
	for _, item := range s.adapters {
		if item.ID() == clientID {
			adapter = item
			break
		}
	}
	if adapter == nil {
		return AppState{}, fmt.Errorf("client %q not found", clientID)
	}

	if toClientSet(config.EnabledClients)[clientID] {
		return s.GetState()
	}

	if _, ok := config.ClientBackups[clientID]; !ok {
		backup, err := captureClientBackup(clientID, adapter.Status().Path)
		if err != nil {
			return AppState{}, err
		}
		config.ClientBackups[clientID] = backup
	}

	result := adapter.Apply(config.Servers, config.LastAppliedByClient[adapter.ID()])
	if !result.Success {
		return AppState{}, errors.New(result.Message)
	}

	config.EnabledClients = setClientEnabled(config.EnabledClients, clientID, true)
	config.LastAppliedByClient[adapter.ID()] = adapter.ManagedServerNames(config.Servers)
	config.UpdatedAt = nowRFC3339()

	if err := s.store.Save(config); err != nil {
		return AppState{}, err
	}

	return s.GetState()
}

func (s *AppService) rollbackConfig(previous AppConfig) error {
	if _, err := s.applyEnabledClients(&previous); err != nil {
		return err
	}
	return nil
}

func (s *AppService) DisableClient(clientID string, restoreBackup bool) (AppState, error) {
	config, err := s.store.Load()
	if err != nil {
		return AppState{}, err
	}

	var adapter ClientAdapter
	for _, item := range s.adapters {
		if item.ID() == clientID {
			adapter = item
			break
		}
	}
	if adapter == nil {
		return AppState{}, fmt.Errorf("client %q not found", clientID)
	}

	if !toClientSet(config.EnabledClients)[clientID] {
		return s.GetState()
	}

	if restoreBackup {
		backup, ok := config.ClientBackups[clientID]
		if !ok {
			return AppState{}, errors.New("no backup available for this client")
		}

		if err := restoreClientBackup(clientID, adapter.Status().Path, backup); err != nil {
			return AppState{}, err
		}

		delete(config.ClientBackups, clientID)
		delete(config.LastAppliedByClient, clientID)
	}

	config.EnabledClients = setClientEnabled(config.EnabledClients, clientID, false)
	config.UpdatedAt = nowRFC3339()

	if err := s.store.Save(config); err != nil {
		return AppState{}, err
	}

	return s.GetState()
}

func (s *AppService) PreviewClientConfig(clientID string) (ClientConfigPreview, error) {
	var adapter ClientAdapter
	for _, item := range s.adapters {
		if item.ID() == clientID {
			adapter = item
			break
		}
	}

	if adapter == nil {
		return ClientConfigPreview{}, fmt.Errorf("client %q not found", clientID)
	}

	status := adapter.Status()
	data, err := os.ReadFile(status.Path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ClientConfigPreview{}, fmt.Errorf("config file not found: %s", status.Path)
		}
		return ClientConfigPreview{}, err
	}

	return ClientConfigPreview{
		ClientID: clientID,
		Path:     status.Path,
		Content:  string(data),
	}, nil
}

func (s *AppService) PreviewAppConfig() (ClientConfigPreview, error) {
	path := s.store.Path()
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ClientConfigPreview{}, fmt.Errorf("config file not found: %s", path)
		}
		return ClientConfigPreview{}, err
	}

	return ClientConfigPreview{
		ClientID: "mcpmgr",
		Path:     path,
		Content:  string(data),
	}, nil
}

func (s *AppService) clientStatuses(config AppConfig) []ClientStatus {
	enabledClients := toClientSet(config.EnabledClients)
	statuses := make([]ClientStatus, 0, len(s.adapters))
	for _, adapter := range s.adapters {
		status := adapter.Status()
		status.Enabled = enabledClients[adapter.ID()]
		_, status.HasBackup = config.ClientBackups[adapter.ID()]
		statuses = append(statuses, status)
	}
	return statuses
}
