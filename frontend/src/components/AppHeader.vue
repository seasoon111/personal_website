<script setup lang="ts">
import { computed, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import { useSiteStore } from '@/stores/site'

const { t, locale } = useI18n()
const site = useSiteStore()
const route = useRoute()
const open = ref(false)

const nav = computed(() => [
  { to: '/', label: t('nav.home') },
  { to: '/articles', label: t('nav.articles') },
  { to: '/projects', label: t('nav.projects') },
  { to: '/moments', label: t('nav.moments') },
  { to: '/about', label: t('nav.about') },
])

function isActive(to: string) {
  return to === '/' ? route.path === '/' : route.path.startsWith(to)
}

function toggleLocale() {
  const next = locale.value === 'zh-CN' ? 'en-US' : 'zh-CN'
  locale.value = next
  localStorage.setItem('pw_locale', next)
}
</script>

<template>
  <header class="sticky top-0 z-40 border-b border-line/70 bg-paper/85 backdrop-blur">
    <div class="container-page flex h-16 items-center justify-between gap-4">
      <router-link to="/" class="flex min-w-0 items-center gap-2">
        <span class="text-xl">📔</span>
        <span class="truncate font-hand text-xl tracking-wide">{{ site.info.name }}</span>
      </router-link>

      <nav class="hidden items-center gap-1 md:flex">
        <router-link
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="rounded-full px-3 py-1.5 text-sm text-ink-soft transition hover:bg-paper-deep hover:text-ink"
          :class="isActive(item.to) ? 'bg-paper-deep text-ink' : ''"
        >
          {{ item.label }}
        </router-link>
      </nav>

      <div class="flex items-center gap-2">
        <button
          class="rounded-full border border-line px-3 py-1 text-xs text-ink-soft transition hover:border-terracotta-soft hover:text-terracotta-deep"
          :title="t('common.langSwitchTitle')"
          @click="toggleLocale"
        >
          {{ t('common.langSwitch') }}
        </button>
        <button
          class="rounded-full border border-line px-3 py-1 text-sm text-ink-soft md:hidden"
          aria-label="menu"
          @click="open = !open"
        >
          {{ open ? '✕' : '☰' }}
        </button>
      </div>
    </div>

    <nav v-if="open" class="border-t border-line/70 bg-paper md:hidden">
      <div class="container-page flex flex-col py-2">
        <router-link
          v-for="item in nav"
          :key="item.to"
          :to="item.to"
          class="rounded-lg px-3 py-2.5 text-sm text-ink-soft"
          :class="isActive(item.to) ? 'bg-paper-deep text-ink' : ''"
          @click="open = false"
        >
          {{ item.label }}
        </router-link>
      </div>
    </nav>
  </header>
</template>
