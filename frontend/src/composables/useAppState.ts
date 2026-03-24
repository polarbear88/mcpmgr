import { reactive, ref } from 'vue'
import { ApplyToAllClients, DeleteServer, DisableClient, EnableClient, GetAppState, PreviewAppConfig, PreviewClientConfig, SaveServer } from '../../wailsjs/go/main/App'
import type { AppState, ApplyResult, ClientConfigPreview, ClientID, MCPServer, ServerForm } from '../types/app'

export function useAppState(t: (key: any, params?: Record<string, string | number>) => string) {
  const state = ref<AppState>({
    configPath: '',
    servers: [],
    clients: [],
  })

  const isLoading = ref(true)
  const isSaving = ref(false)
  const isApplying = ref(false)
  const loadError = ref('')
  const applyResult = ref<ApplyResult | null>(null)

  const editorOpen = ref(false)
  const editingId = ref<string | null>(null)
  const form = reactive<ServerForm>(createEmptyForm())

  function createEmptyForm(): ServerForm {
    return {
      id: '',
      name: '',
      type: 'stdio',
      command: '',
      url: '',
      argsText: '',
      envText: '',
      workingDir: '',
      notes: '',
    }
  }

  function resetForm() {
    Object.assign(form, createEmptyForm())
  }

  function formatEnv(env?: Record<string, string>) {
    if (!env) {
      return ''
    }
    return Object.entries(env)
      .sort(([left], [right]) => left.localeCompare(right))
      .map(([key, value]) => `${key}=${value}`)
      .join('\n')
  }

  async function loadState() {
    isLoading.value = true
    loadError.value = ''

    try {
      state.value = (await GetAppState()) as AppState
    } catch (error) {
      loadError.value = `${t('loadError')}: ${String(error)}`
    } finally {
      isLoading.value = false
    }
  }

  function openCreateEditor() {
    editingId.value = null
    resetForm()
    editorOpen.value = true
  }

  function openEditEditor(server: MCPServer) {
    editingId.value = server.id
    Object.assign(form, {
      id: server.id,
      name: server.name,
      type: server.type,
      command: server.command ?? '',
      url: server.url ?? '',
      argsText: (server.args ?? []).join('\n'),
      envText: formatEnv(server.env),
      workingDir: server.workingDir ?? '',
      notes: server.notes ?? '',
    })
    editorOpen.value = true
  }

  function closeEditor() {
    editorOpen.value = false
    editingId.value = null
    resetForm()
  }

  async function submitForm() {
    isSaving.value = true

    try {
      state.value = (await SaveServer({
        id: form.id,
        name: form.name,
        type: form.type,
        command: form.command,
        url: form.url,
        argsText: form.argsText,
        envText: form.envText,
        workingDir: form.workingDir,
        notes: form.notes,
      })) as AppState
      return state.value
    } catch (error) {
      await loadState()
      throw error
    } finally {
      isSaving.value = false
    }
  }

  async function removeServer(server: MCPServer) {
    try {
      state.value = (await DeleteServer(server.id)) as AppState
      return state.value
    } catch (error) {
      await loadState()
      throw error
    }
  }

  async function applyAllClients() {
    isApplying.value = true

    try {
      applyResult.value = (await ApplyToAllClients()) as ApplyResult
      await loadState()
      return applyResult.value
    } catch (error) {
      throw error
    } finally {
      isApplying.value = false
    }
  }

  async function enableClient(clientID: ClientID) {
    state.value = (await EnableClient(clientID)) as AppState
  }

  async function disableClient(clientID: ClientID, restoreBackup: boolean) {
    state.value = (await DisableClient(clientID, restoreBackup)) as AppState
  }

  async function previewClientConfig(clientID: ClientID) {
    return (await PreviewClientConfig(clientID)) as ClientConfigPreview
  }

  async function previewAppConfig() {
    return (await PreviewAppConfig()) as ClientConfigPreview
  }

  return {
    applyResult,
    closeEditor,
    editorOpen,
    editingId,
    form,
    isApplying,
    isLoading,
    isSaving,
    loadError,
    loadState,
    openCreateEditor,
    openEditEditor,
    removeServer,
    state,
    submitForm,
    applyAllClients,
    enableClient,
    disableClient,
    previewAppConfig,
    previewClientConfig,
  }
}
