<script lang="ts" setup>
import TopBar from '../components/TopBar.vue'
import type { AppState, MCPServer } from '../types/app'

defineProps<{
  isLoading: boolean
  loadError: string
  locale: string
  state: AppState
  t: (key: any, params?: Record<string, string | number>) => string
}>()

const emit = defineEmits<{
  (event: 'add'): void
  (event: 'edit', server: MCPServer): void
  (event: 'remove', server: MCPServer): void
  (event: 'preview-config'): void
}>()
</script>

<template>
  <div class="view">
    <TopBar :subtitle="t('serversSubtitle')" :title="t('serversTitle')">
      <div class="topbar__actions">
        <button class="button button--secondary" type="button" @click="emit('add')">
          {{ t('addServer') }}
        </button>
      </div>
    </TopBar>

    <div class="stats-grid stats-grid--single">
      <article class="stat-card">
        <span class="stat-card__label">{{ t('managedServers') }}</span>
        <strong>{{ state.servers.length }}</strong>
      </article>
    </div>

    <section class="panel">
      <div class="config-strip">
        <span>{{ t('configPath') }}</span>
        <button class="config-strip__button" type="button" @click="emit('preview-config')">
          <code>{{ state.configPath }}</code>
        </button>
      </div>

      <p v-if="loadError" class="flash flash--error">{{ loadError }}</p>

      <div v-if="isLoading" class="empty-state">
        <h3>{{ t('loading') }}</h3>
      </div>

      <div v-else-if="state.servers.length === 0" class="empty-state">
        <h3>{{ t('emptyStateTitle') }}</h3>
        <p>{{ t('emptyStateSummary') }}</p>
        <button class="button button--primary" type="button" @click="emit('add')">
          {{ t('addServer') }}
        </button>
      </div>

      <div v-else class="server-list">
        <article v-for="server in state.servers" :key="server.id" class="server-card">
          <div class="server-card__top">
            <div>
              <h3>{{ server.name }}</h3>
              <div class="server-meta">
                <span class="status-pill status-pill--neutral">
                  {{ server.type === 'stdio' ? t('stdio') : t('http') }}
                </span>
              </div>
            </div>

            <div class="card-actions">
              <button class="button button--ghost button--compact" type="button" @click="emit('edit', server)">
                {{ t('edit') }}
              </button>
              <button class="button button--ghost button--compact" type="button" @click="emit('remove', server)">
                {{ t('delete') }}
              </button>
            </div>
          </div>

          <code class="server-card__command">
            {{ server.type === 'stdio' ? server.command : server.url }}
          </code>
          <div class="server-footer">
            <span>{{ t('updatedAt') }} · {{ new Date(server.updatedAt).toLocaleString(locale) }}</span>
            <span>{{ server.notes || t('noNotes') }}</span>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>
