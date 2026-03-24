import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { messages, type MessageKey } from '../i18n/messages'
import type { Locale, ThemeMode } from '../types/app'

const themeStorageKey = 'mcpmgr-theme-mode'
const localeStorageKey = 'mcpmgr-locale'

export function usePreferences() {
  const themeMode = ref<ThemeMode>('system')
  const locale = ref<Locale>('zh-CN')
  const systemPrefersDark = ref(false)

  const mediaQuery =
    typeof window !== 'undefined'
      ? window.matchMedia('(prefers-color-scheme: dark)')
      : null

  const copy = computed(() => messages[locale.value])

  const effectiveTheme = computed<'light' | 'dark'>(() => {
    if (themeMode.value === 'system') {
      return systemPrefersDark.value ? 'dark' : 'light'
    }
    return themeMode.value
  })

  const themeLabel = computed(() => {
    if (themeMode.value === 'system') {
      return systemPrefersDark.value ? copy.value.autoDark : copy.value.autoLight
    }
    return themeMode.value === 'dark' ? copy.value.dark : copy.value.light
  })

  function t(key: MessageKey, params?: Record<string, string | number>) {
    let text = copy.value[key]
    if (!params) {
      return text
    }
    for (const [paramKey, value] of Object.entries(params)) {
      text = text.replace(`{${paramKey}}`, String(value))
    }
    return text
  }

  function applyTheme() {
    document.documentElement.dataset.theme = effectiveTheme.value
    document.documentElement.dataset.themeMode = themeMode.value
    document.documentElement.lang = locale.value === 'zh-CN' ? 'zh-CN' : 'en'
  }

  function detectLocale(): Locale {
    const preferred = window.navigator.language.toLowerCase()
    return preferred.startsWith('zh') ? 'zh-CN' : 'en-US'
  }

  function updateSystemTheme(event?: MediaQueryListEvent) {
    systemPrefersDark.value = event?.matches ?? mediaQuery?.matches ?? false
  }

  onMounted(() => {
    const savedTheme = window.localStorage.getItem(themeStorageKey) as ThemeMode | null
    const savedLocale = window.localStorage.getItem(localeStorageKey) as Locale | null

    if (savedTheme === 'light' || savedTheme === 'dark' || savedTheme === 'system') {
      themeMode.value = savedTheme
    }
    locale.value = savedLocale === 'zh-CN' || savedLocale === 'en-US' ? savedLocale : detectLocale()

    updateSystemTheme()
    applyTheme()
    mediaQuery?.addEventListener('change', updateSystemTheme)
  })

  onBeforeUnmount(() => {
    mediaQuery?.removeEventListener('change', updateSystemTheme)
  })

  watch(themeMode, (value) => {
    window.localStorage.setItem(themeStorageKey, value)
    applyTheme()
  })

  watch(locale, (value) => {
    window.localStorage.setItem(localeStorageKey, value)
    applyTheme()
  })

  watch(systemPrefersDark, () => {
    applyTheme()
  })

  return {
    locale,
    themeMode,
    themeLabel,
    t,
  }
}
