<script lang="ts" setup>
import TopBar from '../components/TopBar.vue'
import openaiLogo from '../assets/client-logos/openai.svg'
import claudeLogo from '../assets/client-logos/claude-color.svg'
import githubCopilotLogo from '../assets/client-logos/githubcopilot.svg'
import vscodeLogo from '../assets/client-logos/vscode.svg'
import opencodeLogo from '../assets/client-logos/opencode.svg'
import type { AppState, ClientStatus } from '../types/app'

defineProps<{
  isApplying: boolean
  locale: string
  state: AppState
  t: (key: any, params?: Record<string, string | number>) => string
}>()

const emit = defineEmits<{
  (event: 'apply-all'): void
  (event: 'toggle-client', payload: { client: ClientStatus; enabled: boolean }): void
  (event: 'preview-client', client: ClientStatus): void
  (event: 'rescan'): void
}>()

const clientLogos: Record<ClientStatus['id'], string> = {
  codex: openaiLogo,
  claude_code: claudeLogo,
  claude_desktop: claudeLogo,
  copilot_cli: githubCopilotLogo,
  vscode: vscodeLogo,
  opencode: opencodeLogo,
}

const monochromeLogoClients = new Set<ClientStatus['id']>(['codex', 'copilot_cli'])
</script>

<template>
  <div class="view">
    <TopBar :subtitle="t('clientsSubtitle')" :title="t('clientsTitle')">
      <div class="topbar__actions">
        <button class="button button--ghost" type="button" @click="emit('rescan')">
          {{ t('rescanClients') }}
        </button>
        <button class="button button--primary" :disabled="isApplying" type="button" @click="emit('apply-all')">
          {{ t('applyAll') }}
        </button>
      </div>
    </TopBar>

    <section class="panel">
      <div class="client-list client-list--compact">
        <article v-for="client in state.clients" :key="client.id" class="client-item">
          <div class="client-item__content">
            <div class="client-item__title">
              <div class="client-logo-wrap">
                <img
                  :src="clientLogos[client.id]"
                  :alt="client.name"
                  class="client-logo"
                  :class="{ 'client-logo--monochrome': monochromeLogoClients.has(client.id) }"
                />
              </div>
              <h3>{{ client.name }}</h3>
            </div>
            <button
              :title="client.path"
              class="client-card__path client-card__path--single client-card__path-button"
              type="button"
              @click="emit('preview-client', client)"
            >
              {{ client.path }}
            </button>
            <div class="client-card__meta client-card__meta--compact">
              <span>{{ t('format') }} · {{ client.format.toUpperCase() }}</span>
            </div>
          </div>

          <div class="client-item__actions">
            <label class="client-switch">
              <span class="client-switch__label">
                {{ client.enabled ? t('clientEnabled') : t('clientDisabled') }}
              </span>
              <input
                :checked="client.enabled"
                :disabled="isApplying"
                type="checkbox"
                @click.prevent="emit('toggle-client', { client, enabled: !client.enabled })"
              />
              <span class="client-switch__track" aria-hidden="true"></span>
            </label>
            <span
              class="status-pill"
              :class="client.detected ? 'status-pill--active' : 'status-pill--draft'"
            >
              {{ client.detected ? t('statusDetected') : t('statusMissing') }}
            </span>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>
