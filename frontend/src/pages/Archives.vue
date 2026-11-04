<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import type { ArchiveGroup } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const groups = ref<ArchiveGroup[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    groups.value = await api.archives()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container-page">
    <PageHeader :title="t('archives.title')" :subtitle="t('archives.subtitle')" emoji="🗄️" />

    <div v-if="loading" class="space-y-4">
      <div v-for="i in 3" :key="i" class="h-24 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <div v-else-if="groups.length" class="space-y-12">
      <section v-for="group in groups" :key="group.year">
        <div class="mb-4 flex items-baseline gap-3">
          <h2 class="font-hand text-2xl text-terracotta">{{ group.year }}</h2>
          <span class="text-xs text-ink-faint">{{ t('archives.count', { count: group.articles.length }) }}</span>
          <div class="h-[1px] flex-1 bg-line" />
        </div>

        <ul class="space-y-1">
          <li v-for="article in group.articles" :key="article.id">
            <router-link
              :to="`/articles/${article.slug}`"
              class="group flex items-baseline justify-between gap-4 rounded-lg px-2 py-2 transition hover:bg-white/60"
            >
              <span class="min-w-0 truncate text-sm text-ink-soft transition group-hover:text-terracotta-deep">
                {{ article.title }}
              </span>
              <time class="shrink-0 text-xs text-ink-faint">
                {{ formatDate(article.published_at ?? article.created_at, locale, false) }}
              </time>
            </router-link>
          </li>
        </ul>
      </section>
    </div>

    <EmptyState v-else :hint="t('common.emptyHint')" />
  </div>
</template>
