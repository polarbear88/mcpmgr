import type { Locale } from '../types/app'

export type MessageKey =
  | 'appName'
  | 'appTagline'
  | 'navServers'
  | 'navClients'
  | 'navSettings'
  | 'serversTitle'
  | 'serversSubtitle'
  | 'clientsTitle'
  | 'clientsSubtitle'
  | 'settingsTitle'
  | 'settingsSubtitle'
  | 'addServer'
  | 'applyAll'
  | 'syncNow'
  | 'managedServers'
  | 'activeNow'
  | 'connectedClients'
  | 'detectedCount'
  | 'sourceOfTruth'
  | 'configPath'
  | 'emptyStateTitle'
  | 'emptyStateSummary'
  | 'edit'
  | 'delete'
  | 'enabled'
  | 'draft'
  | 'stdio'
  | 'http'
  | 'updatedAt'
  | 'noNotes'
  | 'clientAdapters'
  | 'rescanClients'
  | 'statusDetected'
  | 'statusMissing'
  | 'ready'
  | 'notFound'
  | 'format'
  | 'syncSummary'
  | 'applySummary'
  | 'resultSuccess'
  | 'resultFailed'
  | 'editorTitleCreate'
  | 'editorTitleEdit'
  | 'editorSummary'
  | 'name'
  | 'type'
  | 'command'
  | 'url'
  | 'args'
  | 'argsHint'
  | 'env'
  | 'envHint'
  | 'headers'
  | 'headersHint'
  | 'workingDir'
  | 'notes'
  | 'save'
  | 'saving'
  | 'cancel'
  | 'themeSection'
  | 'theme'
  | 'themeHint'
  | 'languageSection'
  | 'language'
  | 'languageHint'
  | 'auto'
  | 'light'
  | 'dark'
  | 'autoDark'
  | 'autoLight'
  | 'chinese'
  | 'english'
  | 'loadError'
  | 'deleteConfirm'
  | 'loading'
  | 'syncedEnabledClients'
  | 'saveSuccess'
  | 'deleteSuccess'
  | 'clientAutoSync'
  | 'clientEnabled'
  | 'clientDisabled'
  | 'copyContent'
  | 'copied'
  | 'confirmEnableTitle'
  | 'confirmEnableMessage'
  | 'confirmDisableTitle'
  | 'confirmDisableMessage'
  | 'confirmDisableNoBackupMessage'
  | 'enableNow'
  | 'disableOnly'
  | 'restoreBackup'
  | 'cancelAction'
  | 'previewClose'

export const messages: Record<Locale, Record<MessageKey, string>> = {
  'zh-CN': {
    appName: 'Mcpmgr',
    appTagline: '统一管理 MCP 服务并同步到各个 AI 客户端',
    navServers: '服务',
    navClients: '客户端',
    navSettings: '设置',
    serversTitle: 'MCP 服务',
    serversSubtitle: '新增 MCP 服务器后将自动同步到已启用的目标客户端。',
    clientsTitle: '客户端适配器',
    clientsSubtitle: '选择需要自动同步的目标客户端，必要时可以手动全量同步。',
    settingsTitle: '设置',
    settingsSubtitle: '管理界面主题和语言偏好。',
    addServer: '新增服务',
    applyAll: '同步到全部客户端',
    syncNow: '立即同步',
    managedServers: '服务总数',
    activeNow: '启用中 {count} 个',
    connectedClients: '客户端',
    detectedCount: '检测到 {count} 个',
    sourceOfTruth: '真实源配置',
    configPath: '配置路径',
    emptyStateTitle: '还没有 MCP 服务',
    emptyStateSummary: '先新增一条服务定义，后续就能自动同步到多个客户端。',
    edit: '编辑',
    delete: '删除',
    enabled: '已启用',
    draft: '草稿',
    stdio: '命令行启动',
    http: '远程 HTTP',
    updatedAt: '更新时间',
    noNotes: '暂无备注',
    clientAdapters: '客户端适配器',
    rescanClients: '重新扫描',
    statusDetected: '已检测到',
    statusMissing: '未找到',
    ready: '就绪',
    notFound: '未检测到',
    format: '格式',
    syncSummary: '把当前启用的服务写入所有已支持客户端的配置文件。',
    applySummary: '本次同步写入了 {count} 个启用中的 MCP 服务。',
    resultSuccess: '同步成功',
    resultFailed: '同步失败',
    editorTitleCreate: '新增 MCP 服务',
    editorTitleEdit: '编辑 MCP 服务',
    editorSummary: '新增 MCP 服务器后将自动同步到已启用的目标客户端。',
    name: '名称',
    type: '类型',
    command: '命令',
    url: 'URL',
    args: '参数',
    argsHint: '每行一个参数',
    env: '环境变量',
    envHint: '每行一个 KEY=VALUE',
    headers: 'Headers',
    headersHint: '每行一个 Header=Value',
    workingDir: '工作目录',
    notes: '备注',
    save: '保存',
    saving: '保存中...',
    cancel: '取消',
    themeSection: '外观',
    theme: '主题',
    themeHint: '自动模式会跟随系统明暗外观。',
    languageSection: '语言',
    language: '语言',
    languageHint: '默认跟随系统语言，也可以手动切换。',
    auto: '自动',
    light: '浅色',
    dark: '深色',
    autoDark: '自动 · 深色',
    autoLight: '自动 · 浅色',
    chinese: '中文',
    english: 'English',
    loadError: '加载应用状态失败',
    deleteConfirm: '确定删除这个 MCP 服务吗？',
    loading: '正在加载...',
    syncedEnabledClients: '已同步到启用的客户端。',
    saveSuccess: '{name} 已保存。',
    deleteSuccess: '{name} 已删除。',
    clientAutoSync: '自动同步',
    clientEnabled: '已启用',
    clientDisabled: '未启用',
    copyContent: '复制内容',
    copied: '已复制',
    confirmEnableTitle: '启用自动同步',
    confirmEnableMessage: '启用后将立即把当前配置的 MCP 服务器同步到这个客户端。',
    confirmDisableTitle: '关闭自动同步',
    confirmDisableMessage: '关闭后你可以恢复之前备份的原始配置，或者只关闭自动同步并保留当前文件内容。',
    confirmDisableNoBackupMessage: '这个客户端没有可恢复的备份，只能关闭自动同步或取消操作。',
    enableNow: '立即启用',
    disableOnly: '仅关闭',
    restoreBackup: '恢复备份',
    cancelAction: '取消',
    previewClose: '关闭预览',
  },
  'en-US': {
    appName: 'Mcpmgr',
    appTagline: 'Manage MCP servers once and sync them across AI clients',
    navServers: 'Servers',
    navClients: 'Clients',
    navSettings: 'Settings',
    serversTitle: 'MCP Servers',
    serversSubtitle: 'New MCP servers will automatically sync to enabled target clients.',
    clientsTitle: 'Client Adapters',
    clientsSubtitle: 'Choose which target clients receive automatic sync, then use full sync when needed.',
    settingsTitle: 'Settings',
    settingsSubtitle: 'Control theme and language preferences for the desktop app.',
    addServer: 'Add Server',
    applyAll: 'Apply To All Clients',
    syncNow: 'Sync Now',
    managedServers: 'Managed servers',
    activeNow: '{count} active now',
    connectedClients: 'Clients',
    detectedCount: '{count} detected',
    sourceOfTruth: 'Source of truth',
    configPath: 'Config path',
    emptyStateTitle: 'No MCP servers yet',
    emptyStateSummary: 'Create your first server definition here, then sync it out to multiple clients.',
    edit: 'Edit',
    delete: 'Delete',
    enabled: 'Enabled',
    draft: 'Draft',
    stdio: 'CLI / stdio',
    http: 'Remote HTTP',
    updatedAt: 'Updated',
    noNotes: 'No notes yet',
    clientAdapters: 'Client adapters',
    rescanClients: 'Rescan',
    statusDetected: 'Detected',
    statusMissing: 'Missing',
    ready: 'Ready',
    notFound: 'Not found',
    format: 'Format',
    syncSummary: 'Write the currently enabled servers into every supported client config.',
    applySummary: 'This sync wrote {count} enabled MCP server(s).',
    resultSuccess: 'Sync succeeded',
    resultFailed: 'Sync failed',
    editorTitleCreate: 'Add MCP Server',
    editorTitleEdit: 'Edit MCP Server',
    editorSummary: 'New MCP servers will automatically sync to enabled target clients.',
    name: 'Name',
    type: 'Type',
    command: 'Command',
    url: 'URL',
    args: 'Arguments',
    argsHint: 'One argument per line',
    env: 'Environment',
    envHint: 'One KEY=VALUE per line',
    headers: 'Headers',
    headersHint: 'One Header=Value per line',
    workingDir: 'Working directory',
    notes: 'Notes',
    save: 'Save',
    saving: 'Saving...',
    cancel: 'Cancel',
    themeSection: 'Appearance',
    theme: 'Theme',
    themeHint: 'Auto mode follows your system appearance.',
    languageSection: 'Language',
    language: 'Language',
    languageHint: 'Defaults to your system language, but you can override it.',
    auto: 'Auto',
    light: 'Light',
    dark: 'Dark',
    autoDark: 'Auto · Dark',
    autoLight: 'Auto · Light',
    chinese: '中文',
    english: 'English',
    loadError: 'Failed to load application state',
    deleteConfirm: 'Delete this MCP server?',
    loading: 'Loading...',
    syncedEnabledClients: 'Synced to enabled clients.',
    saveSuccess: '{name} saved.',
    deleteSuccess: '{name} deleted.',
    clientAutoSync: 'Auto sync',
    clientEnabled: 'Enabled',
    clientDisabled: 'Disabled',
    copyContent: 'Copy',
    copied: 'Copied',
    confirmEnableTitle: 'Enable Auto Sync',
    confirmEnableMessage: 'Enabling this client will immediately sync the current MCP servers to it.',
    confirmDisableTitle: 'Disable Auto Sync',
    confirmDisableMessage: 'You can restore the saved backup now, or just disable auto sync and keep the current file as-is.',
    confirmDisableNoBackupMessage: 'No backup is available for this client. You can only disable auto sync or cancel.',
    enableNow: 'Enable Now',
    disableOnly: 'Disable Only',
    restoreBackup: 'Restore Backup',
    cancelAction: 'Cancel',
    previewClose: 'Close Preview',
  },
}
