<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import type { DashboardStats } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const stats = ref<DashboardStats | null>(null)
const loading = ref(true)

const cards = computed(() => {
  const s = stats.value
  return [
    { label: t('admin.stats.articles'), value: s?.articles ?? 0, icon: '✍️', to: '/admin/articles' },
    { label: t('admin.stats.drafts'), value: s?.drafts ?? 0, icon: '📝', to: '/admin/articles' },
    { label: t('admin.stats.projects'), value: s?.projects ?? 0, icon: '🧩', to: '/admin/projects' },
    { label: t('admin.stats.moments'), value: s?.moments ?? 0, icon: '💬', to: '/admin/moments' },
    { label: t('admin.stats.images'), value: s?.images ?? 0, icon: '🖼️', to: '/admin/images' },
    { label: t('admin.stats.unread'), value: s?.unread_messages ?? 0, icon: '📮', to: '/admin/messages' },
  ]
})

const passwordForm = reactive({ old_password: '', new_password: '', confirm: '' })
const changing = ref(false)

async function load() {
  loading.value = true
  try {
    stats.value = await api.dashboard()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

async function changePassword() {
  if (passwordForm.new_password.length < 8) {
    toast.error('新密码至少 8 位')
    return
  }
  if (passwordForm.new_password !== passwordForm.confirm) {
    toast.error('两次输入的新密码不一致')
    return
  }
  changing.value = true
  try {
    await api.changePassword({
      old_password: passwordForm.old_password,
      new_password: passwordForm.new_password,
    })
    toast.success(t('admin.passwordChanged'))
    passwordForm.old_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm = ''
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    changing.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="mb-6">
      <h1 class="font-hand text-2xl text-ink">{{ t('admin.dashboard') }}</h1>
      <p class="mt-1 text-sm text-ink-faint">{{ t('admin.dashboardHello') }}</p>
    </div>

    <div class="grid grid-cols-2 gap-3 sm:grid-cols-3 lg:grid-cols-6">
      <router-link
        v-for="card in cards"
        :key="card.label"
        :to="card.to"
        class="rounded-2xl border border-line bg-white/70 p-4 shadow-hand transition hover:-translate-y-0.5 hover:border-terracotta-soft"
      >
        <div class="text-lg">{{ card.icon }}</div>
        <div class="mt-2 font-hand text-2xl text-ink">{{ card.value }}</div>
        <div class="text-xs text-ink-faint">{{ card.label }}</div>
      </router-link>
    </div>

    <div class="mt-6 grid gap-6 lg:grid-cols-[1.4fr_1fr]">
      <section class="rounded-2xl border border-line bg-white/70 p-5">
        <div class="flex items-center justify-between">
          <h2 class="font-hand text-lg">{{ t('admin.recentDrafts') }}</h2>
          <router-link to="/admin/articles/new" class="btn-primary !px-3 !py-1 text-xs">
            + {{ t('admin.newArticle') }}
          </router-link>
        </div>

        <div v-if="loading" class="mt-4 space-y-2">
          <div v-for="i in 3" :key="i" class="h-12 animate-pulse rounded-lg bg-paper-deep" />
        </div>

        <ul v-else-if="stats?.recent_drafts?.length" class="mt-4 divide-y divide-line">
          <li v-for="draft in stats.recent_drafts" :key="draft.id" class="flex items-center justify-between gap-3 py-3">
            <router-link :to="`/admin/articles/${draft.id}/edit`" class="min-w-0 flex-1 truncate text-sm text-ink hover:text-terracotta">
              {{ draft.title }}
            </router-link>
            <time class="shrink-0 text-xs text-ink-faint">{{ formatDate(draft.updated_at, locale) }}</time>
          </li>
        </ul>

        <p v-else class="mt-6 text-sm text-ink-faint">{{ t('admin.noDrafts') }}</p>
      </section>

      <section class="rounded-2xl border border-line bg-white/70 p-5">
        <h2 class="font-hand text-lg">{{ t('admin.changePassword') }}</h2>
        <form class="mt-4 space-y-3" @submit.prevent="changePassword">
          <input v-model="passwordForm.old_password" type="password" class="input-hand" :placeholder="t('admin.oldPassword')" />
          <input v-model="passwordForm.new_password" type="password" class="input-hand" :placeholder="t('admin.newPassword')" />
          <input v-model="passwordForm.confirm" type="password" class="input-hand" placeholder="再次输入新密码" />
          <button type="submit" class="btn-primary w-full" :disabled="changing">
            {{ changing ? t('common.saving') : t('admin.changePassword') }}
          </button>
        </form>
      </section>
    </div>
  </div>
</template>
