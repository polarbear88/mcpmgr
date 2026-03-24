# Mcpmgr

`Mcpmgr` is a desktop MCP server manager for local AI development workflows.

The idea is simple: maintain your MCP servers once inside `Mcpmgr`, then sync them to the AI CLIs and desktop apps you have enabled, instead of editing multiple config files by hand.

[中文 README](/Users/polarbear/code/mcpmgr/README.md)

## Why This Project Exists

In real-world AI workflows, many people use multiple tools at the same time, such as:

- Codex
- Claude Code
- Claude Desktop
- Copilot CLI
- Visual Studio Code

Whenever a new local MCP server is added, you usually have to update every client configuration separately. That is repetitive, error-prone, and hard to maintain.

`Mcpmgr` solves that problem by:

- managing MCP servers in one place
- treating `Mcpmgr`'s own config as the source of truth
- syncing to enabled target clients automatically
- allowing MCP-specific backup and restore when a client is disabled

## Current Features

- Add, edit, and delete MCP servers
- Support both `stdio` and `HTTP` MCP server types
- Enable or disable client adapters individually
- Auto-sync to enabled clients when servers change
- Manual `Apply To All Clients`
- Modify only the MCP section of target config files
- Preserve non-MCP settings in target config files
- Preview and copy client config files inside the app
- Light mode, dark mode, and system theme
- Chinese and English UI
- Built with `Wails + Vue 3 + TypeScript + Go`

## Currently Supported Clients

- Codex
- Claude Code
- Claude Desktop
- Copilot CLI
- Visual Studio Code

## Platform Support

`Mcpmgr` is designed as a cross-platform desktop app for:

- macOS
- Linux
- Windows

Client path handling has been implemented for all three platforms using standard user configuration locations. Real-world validation across different environments is still welcome through issues and pull requests.

## How It Works

The core workflow is:

1. Maintain your MCP servers inside `Mcpmgr`
2. Choose which clients should be auto-synced
3. When servers are added, edited, or deleted, `Mcpmgr` syncs them to enabled clients
4. When a client is disabled, you can restore the previously backed up MCP configuration

Important behavior:

- `Mcpmgr` takes over the MCP section of the target client config
- non-MCP settings are preserved
- backup and restore also affect only the MCP section, not the entire file

## UI Notes

- Fixed left sidebar with `Servers`, `Clients`, and `Settings`
- Main workspace on the right
- Click a client config path to preview its file content in-app
- Confirmation dialogs for destructive or sensitive actions
- Toast for successful sync feedback, dialog for failures

## Development

### Requirements

- Go `1.23`
- Node.js
- `pnpm`
- Wails `v2`

### Install Frontend Dependencies

```bash
cd frontend
pnpm install
```

### Run In Development

```bash
wails dev
```

### Build Frontend

```bash
cd frontend
pnpm build
```

### Build Backend

```bash
go build ./...
```

## Project Structure

Backend files are now split by responsibility:

- [app.go](/Users/polarbear/code/mcpmgr/app.go): Wails entrypoints
- [service_app.go](/Users/polarbear/code/mcpmgr/service_app.go): main application workflows
- [service_sync.go](/Users/polarbear/code/mcpmgr/service_sync.go): sync, backup, restore, config writing
- [client_adapters.go](/Users/polarbear/code/mcpmgr/client_adapters.go): client adapters
- [config_store.go](/Users/polarbear/code/mcpmgr/config_store.go): local config persistence
- [service_types.go](/Users/polarbear/code/mcpmgr/service_types.go): shared types
- [service_helpers.go](/Users/polarbear/code/mcpmgr/service_helpers.go): helpers

Frontend code mainly lives in:

- [frontend/src/App.vue](/Users/polarbear/code/mcpmgr/frontend/src/App.vue)
- [frontend/src/views](/Users/polarbear/code/mcpmgr/frontend/src/views)
- [frontend/src/components](/Users/polarbear/code/mcpmgr/frontend/src/components)
- [frontend/src/composables](/Users/polarbear/code/mcpmgr/frontend/src/composables)
- [frontend/src/i18n](/Users/polarbear/code/mcpmgr/frontend/src/i18n)

## Developers

All code and both README files are written entirely by Codex, with direction and ideas provided by polarbear.

## Planned Next Steps

- Import existing MCP configs from supported clients
- More client adapters
- Sync history and result tracking
- Better release and installation workflow

## License

MIT
