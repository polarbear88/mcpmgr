<script lang="ts" setup>
import appLogo from '../assets/mcpmgr-logo.svg'
import type { ViewID } from '../types/app'

defineProps<{
  activeView: ViewID
  appName: string
  navClients: string
  navServers: string
  navSettings: string
}>()

const emit = defineEmits<{
  (event: 'navigate', view: ViewID): void
}>()
</script>

<template>
  <aside class="sidebar">
    <div class="sidebar__brand">
      <img :alt="appName" :src="appLogo" class="sidebar__brand-logo" />
      <strong>{{ appName }}</strong>
    </div>

    <nav class="sidebar__nav">
      <button
        class="nav-item"
        :class="{ 'nav-item--active': activeView === 'servers' }"
        type="button"
        @click="emit('navigate', 'servers')"
      >
        {{ navServers }}
      </button>
      <button
        class="nav-item"
        :class="{ 'nav-item--active': activeView === 'clients' }"
        type="button"
        @click="emit('navigate', 'clients')"
      >
        {{ navClients }}
      </button>
      <button
        class="nav-item"
        :class="{ 'nav-item--active': activeView === 'settings' }"
        type="button"
        @click="emit('navigate', 'settings')"
      >
        {{ navSettings }}
      </button>
    </nav>
  </aside>
</template>
