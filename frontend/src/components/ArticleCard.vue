<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Article } from '@/types'
import { formatDate } from '@/utils/date'

const props = defineProps<{ article: Article }>()
const { t, locale } = useI18n()

const date = computed(() => formatDate(props.article.published_at ?? props.article.created_at, locale.value))
const tags = computed(() => (props.article.tags ?? []).slice(0, 3))
</script>

<template>
  <router-link
    :to="`/articles/${article.slug}`"
    class="group block rounded-2xl border border-line bg-white/70 p-5 shadow-hand transition duration-300 hover:-translate-y-0.5 hover:border-terracotta-soft hover:shadow-hand-lg"
  >
    <div class="flex gap-4">
      <div class="min-w-0 flex-1">
        <div class="flex flex-wrap items-center gap-2 text-xs text-ink-faint">
          <span v-if="article.is_pinned" class="rounded-full bg-terracotta-soft/40 px-2 py-0.5 text-terracotta-deep">
            📌 {{ t('admin.pinned') }}
          </span>
          <time>{{ date }}</time>
          <span v-if="article.category">· {{ article.category.name }}</span>
        </div>

        <h3 class="mt-2 font-hand text-xl leading-snug text-ink transition group-hover:text-terracotta-deep">
          {{ article.title }}
        </h3>

        <p v-if="article.summary" class="mt-2 line-clamp-2 text-sm leading-6 text-ink-soft">
          {{ article.summary }}
        </p>

        <div v-if="tags.length" class="mt-3 flex flex-wrap gap-2">
          <span v-for="tag in tags" :key="tag.id" class="text-xs text-ink-faint">#{{ tag.name }}</span>
        </div>
      </div>

      <img
        v-if="article.cover_url"
        :src="article.cover_url"
        :alt="article.title"
        loading="lazy"
        class="hidden h-24 w-32 shrink-0 rounded-xl border border-line object-cover sm:block"
      />
    </div>
  </router-link>
</template>
