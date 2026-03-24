<script lang="ts" setup>
import type { Locale, ThemeMode } from '../types/app'
import TopBar from '../components/TopBar.vue'

defineProps<{
  locale: Locale
  themeMode: ThemeMode
  themeLabel: string
  t: (key: any, params?: Record<string, string | number>) => string
}>()

const emit = defineEmits<{
  (event: 'theme-change', value: ThemeMode): void
  (event: 'locale-change', value: Locale): void
}>()
</script>

<template>
  <div class="view">
    <TopBar :subtitle="t('settingsSubtitle')" :title="t('settingsTitle')" />

    <div class="settings-grid">
      <section class="panel settings-card">
        <div class="panel__header">
          <div>
            <span class="section-label">{{ t('themeSection') }}</span>
            <h2>{{ t('theme') }}</h2>
          </div>
          <span class="status-pill status-pill--neutral">{{ themeLabel }}</span>
        </div>

        <div class="segment-control" role="group" aria-label="Theme mode">
          <button
            v-for="mode in ['system', 'light', 'dark'] as ThemeMode[]"
            :key="mode"
            class="segment-control__item"
            :class="{ 'segment-control__item--active': themeMode === mode }"
            type="button"
            @click="emit('theme-change', mode)"
          >
            {{ mode === 'system' ? t('auto') : mode === 'light' ? t('light') : t('dark') }}
          </button>
        </div>

        <p class="settings-card__hint">{{ t('themeHint') }}</p>
      </section>

      <section class="panel settings-card">
        <div class="panel__header">
          <div>
            <span class="section-label">{{ t('languageSection') }}</span>
            <h2>{{ t('language') }}</h2>
          </div>
        </div>

        <div class="segment-control segment-control--two" role="group" aria-label="Language">
          <button
            class="segment-control__item"
            :class="{ 'segment-control__item--active': locale === 'zh-CN' }"
            type="button"
            @click="emit('locale-change', 'zh-CN')"
          >
            {{ t('chinese') }}
          </button>
          <button
            class="segment-control__item"
            :class="{ 'segment-control__item--active': locale === 'en-US' }"
            type="button"
            @click="emit('locale-change', 'en-US')"
          >
            {{ t('english') }}
          </button>
        </div>

        <p class="settings-card__hint">{{ t('languageHint') }}</p>
      </section>
    </div>
  </div>
</template>
