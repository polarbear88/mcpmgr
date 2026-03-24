<script lang="ts" setup>
import { computed, onMounted, ref } from 'vue'
import AppSidebar from './components/AppSidebar.vue'
import ServerEditorModal from './components/ServerEditorModal.vue'
import { useAppState } from './composables/useAppState'
import { usePreferences } from './composables/usePreferences'
import ClientsView from './views/ClientsView.vue'
import SettingsView from './views/SettingsView.vue'
import ServersView from './views/ServersView.vue'
import type { ClientConfigPreview, ClientStatus, Locale, MCPServer, ThemeMode, ViewID } from './types/app'

type ConfirmOption = {
  value: string
  label: string
  primary?: boolean
  danger?: boolean
}

const activeView = ref<ViewID>('servers')

const { locale, themeMode, themeLabel, t } = usePreferences()
const app = useAppState(t)

const isEditing = computed(() => Boolean(app.editingId.value))
const previewTitle = ref('')
const previewPath = ref('')
const previewCode = ref('')
const previewVisible = computed(() => previewCode.value !== '')
const previewLines = computed(() => previewCode.value.split('\n'))
const toastMessage = ref('')
const toastVisible = ref(false)
const dialogTitle = ref('')
const dialogMessage = ref('')
const dialogVisible = ref(false)
const confirmTitle = ref('')
const confirmMessage = ref('')
const confirmVisible = ref(false)
const confirmOptions = ref<ConfirmOption[]>([])
let toastTimer: ReturnType<typeof setTimeout> | null = null
let confirmResolver: ((value: string | null) => void) | null = null

function showToast(message: string) {
  toastMessage.value = message
  toastVisible.value = true

  if (toastTimer) {
    clearTimeout(toastTimer)
  }

  toastTimer = setTimeout(() => {
    toastVisible.value = false
  }, 2200)
}

function showErrorDialog(message: string) {
  dialogTitle.value = locale.value === 'zh-CN' ? '同步失败' : 'Sync Failed'
  dialogMessage.value = message
  dialogVisible.value = true
}

function askConfirm(title: string, message: string, options: ConfirmOption[]) {
  confirmTitle.value = title
  confirmMessage.value = message
  confirmOptions.value = options
  confirmVisible.value = true

  return new Promise<string | null>((resolve) => {
    confirmResolver = resolve
  })
}

function closeConfirm(value: string | null) {
  confirmVisible.value = false
  confirmTitle.value = ''
  confirmMessage.value = ''
  confirmOptions.value = []

  if (confirmResolver) {
    confirmResolver(value)
    confirmResolver = null
  }
}

function showPreview(preview: ClientConfigPreview, clientName: string) {
  previewTitle.value = clientName
  previewPath.value = preview.path
  previewCode.value = preview.content
}

function closePreview() {
  previewTitle.value = ''
  previewPath.value = ''
  previewCode.value = ''
}

function handleNavigate(view: ViewID) {
  activeView.value = view
  closePreview()
}

async function handleApplyAll() {
  try {
    const result = await app.applyAllClients()
    showToast(t('applySummary', { count: result.appliedServers }))
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handleSaveServer() {
  try {
    await app.submitForm()
    app.closeEditor()
    showToast(t('syncedEnabledClients'))
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handleDeleteServer(server: MCPServer) {
  const choice = await askConfirm(t('deleteConfirm'), server.name, [
    { value: 'cancel', label: t('cancelAction') },
    { value: 'delete', label: t('delete'), danger: true },
  ])

  if (choice !== 'delete') {
    return
  }

  try {
    await app.removeServer(server)
    showToast(t('syncedEnabledClients'))
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handlePreviewAppConfig() {
  try {
    const preview = await app.previewAppConfig()
    showPreview(preview, t('appName'))
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handleClientToggle(client: ClientStatus, enabled: boolean) {
  if (enabled) {
    const choice = await askConfirm(t('confirmEnableTitle'), t('confirmEnableMessage'), [
      { value: 'cancel', label: t('cancelAction') },
      { value: 'enable', label: t('enableNow'), primary: true },
    ])

    if (choice !== 'enable') {
      return
    }

    try {
      await app.enableClient(client.id)
      showToast(`${client.name} · ${t('clientEnabled')}`)
    } catch (error) {
      showErrorDialog(String(error))
    }

    return
  }

  const message = client.hasBackup ? t('confirmDisableMessage') : t('confirmDisableNoBackupMessage')
  const options: ConfirmOption[] = [{ value: 'cancel', label: t('cancelAction') }]

  if (client.hasBackup) {
    options.push({ value: 'restore', label: t('restoreBackup'), primary: true })
  }
  options.push({ value: 'disable', label: t('disableOnly'), primary: !client.hasBackup })

  const choice = await askConfirm(t('confirmDisableTitle'), message, options)
  if (choice === 'cancel' || choice === null) {
    return
  }

  try {
    await app.disableClient(client.id, choice === 'restore')
    showToast(`${client.name} · ${t('clientDisabled')}`)
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handlePreviewClient(client: ClientStatus) {
  try {
    const preview = await app.previewClientConfig(client.id)
    showPreview(preview, client.name)
  } catch (error) {
    showErrorDialog(String(error))
  }
}

async function handleCopyPreviewCode() {
  if (!previewCode.value) {
    return
  }

  try {
    await navigator.clipboard.writeText(previewCode.value)
    showToast(t('copied'))
  } catch (error) {
    showErrorDialog(String(error))
  }
}

onMounted(async () => {
  await app.loadState()
})
</script>

<template>
  <div class="desktop-shell">
    <AppSidebar
      :active-view="activeView"
      :app-name="t('appName')"
      :nav-clients="t('navClients')"
      :nav-settings="t('navSettings')"
      :nav-servers="t('navServers')"
      @navigate="handleNavigate($event as ViewID)"
    />

    <main class="workspace">
      <section v-if="previewVisible" class="preview-page panel">
        <div class="preview-page__header">
          <div>
            <h2>{{ previewTitle }}</h2>
            <p>{{ previewPath }}</p>
          </div>
          <div class="preview-page__actions">
            <button class="button button--ghost button--compact" type="button" @click="handleCopyPreviewCode">
              {{ t('copyContent') }}
            </button>
            <button class="button button--primary button--compact" type="button" @click="closePreview">
              {{ t('previewClose') }}
            </button>
          </div>
        </div>
        <div class="preview-page__body">
          <div class="preview-page__lines" aria-hidden="true">
            <span v-for="(_, index) in previewLines" :key="`preview-line-${index}`">{{ index + 1 }}</span>
          </div>
          <pre class="preview-page__content"><code>{{ previewCode }}</code></pre>
        </div>
      </section>

      <ServersView
        v-else-if="activeView === 'servers'"
        :is-loading="app.isLoading.value"
        :load-error="app.loadError.value"
        :locale="locale"
        :state="app.state.value"
        :t="t"
        @add="app.openCreateEditor"
        @edit="app.openEditEditor"
        @preview-config="handlePreviewAppConfig"
        @remove="handleDeleteServer"
      />

      <ClientsView
        v-else-if="activeView === 'clients'"
        :is-applying="app.isApplying.value"
        :locale="locale"
        :state="app.state.value"
        :t="t"
        @apply-all="handleApplyAll"
        @preview-client="handlePreviewClient($event as ClientStatus)"
        @toggle-client="handleClientToggle($event.client as ClientStatus, $event.enabled as boolean)"
        @rescan="app.loadState"
      />

      <SettingsView
        v-else
        :locale="locale"
        :t="t"
        :theme-label="themeLabel"
        :theme-mode="themeMode"
        @locale-change="locale = $event as Locale"
        @theme-change="themeMode = $event as ThemeMode"
      />
    </main>

    <ServerEditorModal
      v-if="app.editorOpen.value"
      :form="app.form"
      :is-editing="isEditing"
      :is-saving="app.isSaving.value"
      :locale="locale"
      :t="t"
      @close="app.closeEditor"
      @submit="handleSaveServer"
    />

    <div v-if="toastVisible" class="toast">{{ toastMessage }}</div>

    <section v-if="confirmVisible" class="dialog-overlay" @click.self="closeConfirm('cancel')">
      <div class="dialog">
        <h3>{{ confirmTitle }}</h3>
        <p>{{ confirmMessage }}</p>
        <div class="dialog__actions">
          <button
            v-for="option in confirmOptions"
            :key="option.value"
            class="button"
            :class="option.danger ? 'button--danger' : option.primary ? 'button--primary' : 'button--ghost'"
            type="button"
            @click="closeConfirm(option.value)"
          >
            {{ option.label }}
          </button>
        </div>
      </div>
    </section>

    <section v-if="dialogVisible" class="dialog-overlay" @click.self="dialogVisible = false">
      <div class="dialog">
        <h3>{{ dialogTitle }}</h3>
        <p>{{ dialogMessage }}</p>
        <div class="dialog__actions">
          <button class="button button--primary" type="button" @click="dialogVisible = false">
            {{ locale === 'zh-CN' ? '关闭' : 'Close' }}
          </button>
        </div>
      </div>
    </section>
  </div>
</template>
