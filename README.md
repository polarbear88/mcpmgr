# Mcpmgr

`Mcpmgr` 是一个面向本地 AI 开发工作流的 MCP Server 管理器。

它的目标很直接：你只需要在 `Mcpmgr` 里维护一份 MCP 服务配置，然后把它自动同步到你已经启用的 AI CLI 或桌面应用里，不再反复手改多个客户端配置文件。

[English README](/Users/polarbear/code/mcpmgr/README.en.md)

## 项目背景

在实际使用 AI 开发工具时，很多人会同时用到多个客户端，例如：

- Codex
- Claude Code
- Claude Desktop
- Copilot CLI
- Visual Studio Code

一旦本地新增了一个 MCP Server，往往就要分别去这些客户端各自的配置文件里再加一遍，既重复又容易出错。

`Mcpmgr` 就是为了解决这个问题而做的：

- 在应用内统一管理 MCP 服务
- 以 `Mcpmgr` 自己的配置作为真实源
- 自动同步到已启用的目标客户端
- 在关闭客户端同步时，支持恢复之前备份的 MCP 配置

## 当前特性

- 支持新增、编辑、删除 MCP 服务
- 支持 `stdio` 和 `HTTP` 两种 MCP 类型
- 支持多客户端适配器启用/关闭
- 启用客户端时自动同步当前 MCP 配置
- 支持手动“同步到全部客户端”
- 只修改目标配置文件中的 MCP 部分，不影响其他配置项
- 支持客户端配置文件预览与复制
- 支持浅色、深色、跟随系统
- 支持中文和英文界面
- 基于 `Wails + Vue 3 + TypeScript + Go`

## 当前支持的客户端

- Codex
- Claude Code
- Claude Desktop
- Copilot CLI
- Visual Studio Code

## 平台支持

`Mcpmgr` 以跨平台桌面应用方式实现，目标平台为：

- macOS
- Linux
- Windows

其中客户端适配路径已经按三平台常见目录实现。不同客户端在不同平台上的真实安装和用户配置习惯可能略有差异，因此欢迎通过 issue 或 PR 补充验证结果。

## 工作方式

`Mcpmgr` 的核心思路是：

1. 你在 `Mcpmgr` 里维护 MCP 服务列表
2. 你选择哪些客户端处于“已启用自动同步”状态
3. 当服务发生增删改时，`Mcpmgr` 会自动同步到这些已启用客户端
4. 若关闭某个客户端，还可以选择恢复该客户端之前备份的 MCP 配置

几个重要约束：

- `Mcpmgr` 同步时会接管目标客户端中的 MCP 配置部分
- 非 MCP 的其他配置项会被保留
- 备份与恢复也只针对 MCP 部分，不覆盖整个配置文件

## 界面与交互

- 左侧为固定导航，包含 `服务`、`客户端`、`设置`
- 右侧为主工作区
- 点击客户端配置路径，可直接在应用内预览配置文件内容
- 删除、启用、关闭等高风险操作会要求确认
- 同步成功使用 toast 提示，失败使用 dialog 提示

## 开发

### 依赖

- Go `1.23`
- Node.js
- `pnpm`
- Wails `v2`

### 前端安装

```bash
cd frontend
pnpm install
```

### 开发模式

```bash
wails dev
```

### 构建前端

```bash
cd frontend
pnpm build
```

### 构建后端

```bash
go build ./...
```

## 代码结构

后端已经按职责做了拆分：

- [app.go](/Users/polarbear/code/mcpmgr/app.go): Wails 暴露入口
- [service_app.go](/Users/polarbear/code/mcpmgr/service_app.go): 应用核心业务流程
- [service_sync.go](/Users/polarbear/code/mcpmgr/service_sync.go): 同步、备份、恢复、配置写入
- [client_adapters.go](/Users/polarbear/code/mcpmgr/client_adapters.go): 各客户端适配器
- [config_store.go](/Users/polarbear/code/mcpmgr/config_store.go): 本地配置存储
- [service_types.go](/Users/polarbear/code/mcpmgr/service_types.go): 类型定义
- [service_helpers.go](/Users/polarbear/code/mcpmgr/service_helpers.go): 工具函数

前端主要位于：

- [frontend/src/App.vue](/Users/polarbear/code/mcpmgr/frontend/src/App.vue)
- [frontend/src/views](/Users/polarbear/code/mcpmgr/frontend/src/views)
- [frontend/src/components](/Users/polarbear/code/mcpmgr/frontend/src/components)
- [frontend/src/composables](/Users/polarbear/code/mcpmgr/frontend/src/composables)
- [frontend/src/i18n](/Users/polarbear/code/mcpmgr/frontend/src/i18n)

## 开发人员

代码、README完全由Codex编写，由polarbear指挥和提供ideas。

## 规划中的后续能力

- 导入已有客户端中的 MCP 配置
- 更多客户端适配器
- 同步历史与结果追踪
- 更完善的发布与安装流程

## License

MIT
