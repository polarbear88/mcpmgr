<script lang="ts" setup>
import { computed } from 'vue'
import type { ServerForm } from '../types/app'

const props = defineProps<{
  form: ServerForm
  isSaving: boolean
  isEditing: boolean
  locale: string
  t: (key: any, params?: Record<string, string | number>) => string
}>()

const emit = defineEmits<{
  (event: 'close'): void
  (event: 'submit'): void
}>()

const title = computed(() =>
  props.isEditing ? props.t('editorTitleEdit') : props.t('editorTitleCreate'),
)
</script>

<template>
  <section class="editor-overlay" @click.self="emit('close')">
    <form class="editor-panel" @submit.prevent="emit('submit')">
      <div class="editor-panel__header">
        <div>
          <span class="section-label">{{ t('sourceOfTruth') }}</span>
          <h2>{{ title }}</h2>
          <p>{{ t('editorSummary') }}</p>
        </div>
        <button class="button button--ghost button--compact" type="button" @click="emit('close')">
          {{ t('cancel') }}
        </button>
      </div>

      <div class="editor-grid">
        <label class="field">
          <span>{{ t('name') }}</span>
          <input v-model="form.name" required type="text" />
        </label>

        <label class="field">
          <span>{{ t('type') }}</span>
          <select v-model="form.type">
            <option value="stdio">{{ t('stdio') }}</option>
            <option value="http">{{ t('http') }}</option>
          </select>
        </label>

        <label v-if="form.type === 'stdio'" class="field field--full">
          <span>{{ t('command') }}</span>
          <input v-model="form.command" required type="text" />
        </label>

        <label v-else class="field field--full">
          <span>{{ t('url') }}</span>
          <input v-model="form.url" required type="url" />
        </label>

        <label v-if="form.type === 'stdio'" class="field">
          <span>{{ t('args') }}</span>
          <textarea v-model="form.argsText" rows="5" />
          <small>{{ t('argsHint') }}</small>
        </label>

        <label class="field" :class="{ 'field--full': form.type === 'http' }">
          <span>{{ t('env') }}</span>
          <textarea v-model="form.envText" rows="5" />
          <small>{{ t('envHint') }}</small>
        </label>

        <label class="field">
          <span>{{ t('workingDir') }}</span>
          <input v-model="form.workingDir" type="text" />
        </label>

        <label class="field">
          <span>{{ t('notes') }}</span>
          <input v-model="form.notes" type="text" />
        </label>
      </div>

      <div class="editor-panel__footer">
        <button class="button button--ghost" type="button" @click="emit('close')">
          {{ t('cancel') }}
        </button>
        <button class="button button--primary" :disabled="isSaving" type="submit">
          {{ isSaving ? t('saving') : t('save') }}
        </button>
      </div>
    </form>
  </section>
</template>
