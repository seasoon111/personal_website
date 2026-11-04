<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import Pagination from '@/components/Pagination.vue'
import type { Article, Category, PageResult, Tag } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const result = ref<PageResult<Article> | null>(null)
const loading = ref(true)

const activeCategory = computed(() => (route.query.category as string) ?? '')
const activeTag = computed(() => (route.query.tag as string) ?? '')
const page = computed(() => Number(route.query.page ?? 1) || 1)

async function load() {
  loading.value = true
  try {
    const [list, cats, tagList] = await Promise.all([
      api.articles({
        page: page.value,
        page_size: 8,
        category: activeCategory.value || undefined,
        tag: activeTag.value || undefined,
      }),
      api.categories(),
      api.tags(),
    ])
    result.value = list
    categories.value = cats
    tags.value = tagList
  } finally {
    loading.value = false
  }
}

function setQuery(patch: Record<string, string | number | undefined>) {
  const query = { ...route.query, ...patch }
  Object.keys(query).forEach((key) => {
    if (query[key] === '' || query[key] === undefined) delete query[key]
  })
  router.push({ query })
}

function selectCategory(slug: string) {
  setQuery({ category: slug || undefined, page: 1 })
}

function clearFilter() {
  router.push({ query: {} })
}

watch(() => route.query, load)
onMounted(load)
</script>

<template>
  <div class="container-page">
    <PageHeader :title="t('articles.title')" :subtitle="t('articles.subtitle')" emoji="✍️" />

    <!-- 分类筛选 -->
    <div class="mb-8 flex flex-wrap items-center gap-2">
      <button
        class="rounded-full border px-3.5 py-1 text-sm transition"
        :class="
          !activeCategory && !activeTag
            ? 'border-terracotta bg-terracotta text-white'
            : 'border-line bg-white/60 text-ink-soft hover:border-terracotta-soft'
        "
        @click="clearFilter"
      >
        {{ t('common.all') }}
      </button>
      <button
        v-for="category in categories"
        :key="category.id"
        class="rounded-full border px-3.5 py-1 text-sm transition"
        :class="
          activeCategory === category.slug
            ? 'border-terracotta bg-terracotta text-white'
            : 'border-line bg-white/60 text-ink-soft hover:border-terracotta-soft'
        "
        @click="selectCategory(category.slug)"
      >
        {{ category.name }}
        <span class="ml-1 text-xs opacity-70">{{ category.article_count }}</span>
      </button>
    </div>

    <div v-if="activeTag" class="mb-6 flex items-center gap-3 text-sm text-ink-soft">
      <span>{{ t('articles.byTag', { tag: activeTag }) }}</span>
      <button class="text-terracotta underline decoration-dashed underline-offset-4" @click="setQuery({ tag: undefined, page: 1 })">
        {{ t('articles.clearFilter') }}
      </button>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 4" :key="i" class="h-28 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <div v-else-if="result?.items?.length" class="space-y-3">
      <ArticleCard v-for="(article, index) in result.items" :key="article.id" v-reveal="index * 40" :article="article" />
      <Pagination
        :page="result.page"
        :page-size="result.page_size"
        :total="result.total"
        @update:page="(p) => setQuery({ page: p })"
      />
    </div>

    <EmptyState v-else :text="t('articles.empty')" :hint="t('common.emptyHint')" />

    <div v-if="tags.length" class="mt-12 flex flex-wrap items-center gap-2 text-xs text-ink-faint">
      <span class="mr-1">{{ t('common.tags') }}：</span>
      <button
        v-for="tag in tags"
        :key="tag.id"
        class="rounded-full bg-paper-deep px-2.5 py-1 transition hover:text-terracotta"
        @click="setQuery({ tag: tag.slug, page: 1 })"
      >
        #{{ tag.name }} <span class="opacity-60">{{ tag.article_count }}</span>
      </button>
    </div>
  </div>
</template>
