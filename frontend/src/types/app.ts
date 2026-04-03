export type ThemeMode = 'system' | 'light' | 'dark'
export type Locale = 'zh-CN' | 'en-US'
export type ServerType = 'stdio' | 'http'
export type ClientID = 'codex' | 'claude_code' | 'claude_desktop' | 'copilot_cli' | 'vscode' | 'opencode'
export type ViewID = 'servers' | 'clients' | 'settings'

export type MCPServer = {
  id: string
  name: string
  type: ServerType
  command?: string
  args?: string[]
  url?: string
  env?: Record<string, string>
  workingDir?: string
  notes?: string
  updatedAt: string
}

export type ClientStatus = {
  id: ClientID
  name: string
  path: string
  format: string
  detected: boolean
  supported: boolean
  enabled: boolean
  hasBackup: boolean
  notes: string
}

export type AppState = {
  configPath: string
  servers: MCPServer[]
  clients: ClientStatus[]
}

export type ApplyTargetResult = {
  clientId: ClientID
  clientName: string
  path: string
  success: boolean
  message: string
}

export type ApplyResult = {
  appliedAt: string
  appliedServers: number
  results: ApplyTargetResult[]
}

export type ClientConfigPreview = {
  clientId: ClientID
  path: string
  content: string
}

export type ServerForm = {
  id: string
  name: string
  type: ServerType
  command: string
  url: string
  argsText: string
  envText: string
  workingDir: string
  notes: string
}
