package main

const configVersion = 1

type AppConfig struct {
	Version             int                     `json:"version"`
	UpdatedAt           string                  `json:"updatedAt"`
	Servers             []MCPServer             `json:"servers"`
	EnabledClients      []string                `json:"enabledClients"`
	ClientBackups       map[string]ClientBackup `json:"clientBackups"`
	LastAppliedByClient map[string][]string     `json:"lastAppliedByClient"`
}

type ClientBackup struct {
	Kind      string `json:"kind"`
	Existed   bool   `json:"existed"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updatedAt"`
}

type MCPServer struct {
	ID         string            `json:"id"`
	Name       string            `json:"name"`
	Type       string            `json:"type"`
	Command    string            `json:"command,omitempty"`
	Args       []string          `json:"args,omitempty"`
	URL        string            `json:"url,omitempty"`
	Env        map[string]string `json:"env,omitempty"`
	WorkingDir string            `json:"workingDir,omitempty"`
	Notes      string            `json:"notes,omitempty"`
	UpdatedAt  string            `json:"updatedAt"`
}

type ServerInput struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	Command    string `json:"command"`
	ArgsText   string `json:"argsText"`
	URL        string `json:"url"`
	EnvText    string `json:"envText"`
	WorkingDir string `json:"workingDir"`
	Notes      string `json:"notes"`
}

type ClientStatus struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Format    string `json:"format"`
	Detected  bool   `json:"detected"`
	Supported bool   `json:"supported"`
	Enabled   bool   `json:"enabled"`
	HasBackup bool   `json:"hasBackup"`
	Notes     string `json:"notes"`
}

type AppState struct {
	ConfigPath string         `json:"configPath"`
	Servers    []MCPServer    `json:"servers"`
	Clients    []ClientStatus `json:"clients"`
}

type ApplyTargetResult struct {
	ClientID   string `json:"clientId"`
	ClientName string `json:"clientName"`
	Path       string `json:"path"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
}

type ApplyResult struct {
	AppliedAt      string              `json:"appliedAt"`
	AppliedServers int                 `json:"appliedServers"`
	Results        []ApplyTargetResult `json:"results"`
}

type ClientConfigPreview struct {
	ClientID string `json:"clientId"`
	Path     string `json:"path"`
	Content  string `json:"content"`
}

type AppService struct {
	store    *ConfigStore
	adapters []ClientAdapter
}

type ConfigStore struct {
	path string
}

type ClientAdapter interface {
	ID() string
	Status() ClientStatus
	Apply(servers []MCPServer, previous []string) ApplyTargetResult
	ManagedServerNames(servers []MCPServer) []string
}

type clientBase struct {
	id      string
	name    string
	path    string
	format  string
	support bool
	notes   string
}

type CodexAdapter struct {
	clientBase
}

type ClaudeCodeAdapter struct {
	clientBase
}

type ClaudeDesktopAdapter struct {
	clientBase
}

type CopilotCLIAdapter struct {
	clientBase
}

type VSCodeAdapter struct {
	clientBase
}

type OpenCodeAdapter struct {
	clientBase
}
