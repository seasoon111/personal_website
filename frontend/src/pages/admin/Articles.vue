<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import Pagination from '@/components/Pagination.vue'
import type { Article, PageResult } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const result = ref<PageResult<Article> | null>(null)
const loading = ref(true)
const status = ref('')
const keyword = ref('')
const page = ref(1)

async function load() {
  loading.value = true
  try {
    result.value = await api.adminArticles({
      page: page.value,
      page_size: 10,
      status: status.value || undefined,
      keyword: keyword.value || undefined,
    })
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

function search() {
  page.value = 1
  load()
}

async function toggleStatus(article: Article) {
  try {
    const next = article.status === 'published' ? 'draft' : 'published'
    await api.setArticleStatus(article.id, next)
    article.status = next
    toast.success(next === 'published' ? '已发布' : '已转为草稿')
  } catch (e) {
    toast.error((e as Error).message)
  }
}

async function togglePin(article: Article) {
  try {
    await api.setArticlePin(article.id, !article.is_pinned)
    article.is_pinned = !article.is_pinned
    toast.success('已更新')
  } catch (e) {
    toast.error((e as Error).message)
  }
}

async function remove(article: Article) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteArticle(article.id)
    toast.success(t('admin.deleted'))
    load()
  } catch (e) {
    toast.error((e as Error).message)
  }
}

function changePage(p: number) {
  page.value = p
  load()
}

onMounted(load)
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-hand text-2xl">{{ t('admin.articles') }}</h1>
      <router-link to="/admin/articles/new" class="btn-primary">+ {{ t('admin.newArticle') }}</router-link>
    </div>

    <div class="mb-5 flex flex-wrap items-center gap-2">
      <input
        v-model="keyword"
        class="input-hand max-w-xs"
        placeholder="搜索标题 / 摘要"
        @keyup.enter="search"
      />
      <select v-model="status" class="input-hand w-32" @change="search">
        <option value="">{{ t('common.all') }}</option>
        <option value="published">{{ t('admin.published') }}</option>
        <option value="draft">{{ t('admin.draft') }}</option>
      </select>
      <button type="button" class="btn-hand" @click="search">搜索</button>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 4" :key="i" class="h-20 animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else-if="result?.items?.length" class="space-y-2">
      <div
        v-for="article in result.items"
        :key="article.id"
        class="flex flex-col gap-3 rounded-xl border border-line bg-white/70 p-4 sm:flex-row sm:items-center"
      >
        <div class="min-w-0 flex-1">
          <div class="flex flex-wrap items-center gap-2 text-xs text-ink-faint">
            <span
              class="rounded-full px-2 py-0.5"
              :class="article.status === 'published' ? 'bg-sage/20 text-ink' : 'bg-paper-deep text-ink-soft'"
            >
              {{ article.status === 'published' ? t('admin.published') : t('admin.draft') }}
            </span>
            <span v-if="article.is_pinned">📌</span>
            <time>{{ formatDate(article.updated_at, locale) }}</time>
            <span v-if="article.category">· {{ article.category.name }}</span>
          </div>
          <p class="mt-1.5 truncate font-hand text-lg text-ink">{{ article.title }}</p>
        </div>

        <div class="flex flex-wrap gap-1.5 text-xs">
          <router-link :to="`/admin/articles/${article.id}/edit`" class="btn-hand">{{ t('common.edit') }}</router-link>
          <button type="button" class="btn-hand" @click="toggleStatus(article)">
            {{ article.status === 'published' ? t('admin.unpublish') : t('admin.publish') }}
          </button>
          <button type="button" class="btn-hand" @click="togglePin(article)">
            {{ article.is_pinned ? t('admin.unpin') : t('admin.pin') }}
          </button>
          <button type="button" class="btn-hand !text-terracotta-deep" @click="remove(article)">
            {{ t('common.delete') }}
          </button>
        </div>
      </div>

      <Pagination :page="result.page" :page-size="result.page_size" :total="result.total" @update:page="changePage" />
    </div>

    <p v-else class="py-16 text-center text-sm text-ink-faint">还没有文章，点右上角开始写第一篇。</p>
  </div>
</template>
