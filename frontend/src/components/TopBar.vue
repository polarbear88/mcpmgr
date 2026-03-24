<script lang="ts" setup>
defineProps<{
  title: string
  subtitle: string
}>()
</script>

<template>
  <header class="topbar">
    <div>
      <h1>{{ title }}</h1>
      <p>{{ subtitle }}</p>
    </div>
    <slot />
  </header>
</template>
