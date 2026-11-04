<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import ToastHost from '@/components/ToastHost.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const items = computed(() => [
  { to: '/admin', label: t('admin.dashboard'), icon: '🏠', exact: true },
  { to: '/admin/articles', label: t('admin.articles'), icon: '✍️' },
  { to: '/admin/categories', label: t('admin.categories'), icon: '🗂️' },
  { to: '/admin/tags', label: t('admin.tags'), icon: '🏷️' },
  { to: '/admin/projects', label: t('admin.projects'), icon: '🧩' },
  { to: '/admin/moments', label: t('admin.moments'), icon: '💬' },
  { to: '/admin/about', label: t('admin.about'), icon: '🙋' },
  { to: '/admin/images', label: t('admin.images'), icon: '🖼️' },
  { to: '/admin/messages', label: t('admin.messages'), icon: '📮' },
])

function isActive(item: { to: string; exact?: boolean }) {
  return item.exact ? route.path === item.to : route.path.startsWith(item.to)
}

function logout() {
  auth.logout()
  router.push('/admin/login')
}

onMounted(() => {
  if (!auth.user) auth.fetchMe()
})
</script>

<template>
  <div class="flex min-h-screen bg-paper">
    <aside class="hidden w-60 shrink-0 border-r border-line bg-white/50 lg:block">
      <div class="sticky top-0 flex h-screen flex-col p-4">
        <router-link to="/admin" class="flex items-center gap-2 px-2 py-3">
          <span class="text-xl">📔</span>
          <span class="font-hand text-lg">管理后台</span>
        </router-link>

        <nav class="mt-2 flex-1 space-y-1 overflow-y-auto">
          <router-link
            v-for="item in items"
            :key="item.to"
            :to="item.to"
            class="flex items-center gap-2 rounded-xl px-3 py-2 text-sm transition"
            :class="
              isActive(item)
                ? 'bg-terracotta-soft/35 font-medium text-terracotta-deep'
                : 'text-ink-soft hover:bg-paper-deep'
            "
          >
            <span>{{ item.icon }}</span>
            {{ item.label }}
          </router-link>
        </nav>

        <div class="border-t border-line pt-3 text-sm">
          <router-link to="/" class="block rounded-xl px-3 py-2 text-ink-soft hover:bg-paper-deep">
            ← {{ t('admin.backToSite') }}
          </router-link>
          <button class="mt-1 block w-full rounded-xl px-3 py-2 text-left text-ink-soft hover:bg-paper-deep" @click="logout">
            {{ t('admin.logout') }}
          </button>
        </div>
      </div>
    </aside>

    <div class="flex min-w-0 flex-1 flex-col">
      <header class="sticky top-0 z-30 flex h-14 items-center justify-between border-b border-line bg-paper/90 px-4 backdrop-blur sm:px-6">
        <div class="flex items-center gap-2 lg:hidden">
          <span class="font-hand text-lg">📔 后台</span>
        </div>
        <div class="ml-auto flex items-center gap-3 text-sm text-ink-soft">
          <span class="hidden sm:inline">{{ auth.user?.username ?? '' }}</span>
          <router-link to="/" class="rounded-full border border-line px-3 py-1 text-xs hover:border-terracotta-soft">
            {{ t('admin.backToSite') }}
          </router-link>
          <button class="rounded-full border border-line px-3 py-1 text-xs hover:border-terracotta-soft" @click="logout">
            {{ t('admin.logout') }}
          </button>
        </div>
      </header>

      <nav class="flex gap-1 overflow-x-auto border-b border-line bg-white/40 px-3 py-2 lg:hidden">
        <router-link
          v-for="item in items"
          :key="item.to"
          :to="item.to"
          class="shrink-0 rounded-full px-3 py-1 text-xs"
          :class="isActive(item) ? 'bg-terracotta text-white' : 'bg-paper-deep text-ink-soft'"
        >
          {{ item.label }}
        </router-link>
      </nav>

      <main class="flex-1 p-4 sm:p-6">
        <router-view />
      </main>
    </div>

    <ToastHost />
  </div>
</template>
