<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import hljs from 'highlight.js/lib/common'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import type { ArticleDetail } from '@/types'
import { formatDate } from '@/utils/date'
import { renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const { t, locale } = useI18n()

const data = ref<ArticleDetail | null>(null)
const loading = ref(true)
const error = ref('')
const bodyEl = ref<HTMLElement | null>(null)
const toc = ref<{ id: string; text: string; level: number }[]>([])

const article = computed(() => data.value?.article)
const html = computed(() => article.value?.content_html || renderMarkdown(article.value?.content_md))
const date = computed(() =>
  formatDate(article.value?.published_at ?? article.value?.created_at, locale.value),
)

async function load() {
  loading.value = true
  error.value = ''
  try {
    data.value = await api.article(route.params.slug as string)
  } catch (e) {
    error.value = (e as Error).message
    data.value = null
  } finally {
    loading.value = false
  }
  // 等正文渲染完成后再生成目录、做代码高亮
  await nextTick()
  buildToc()
  highlightCode()
}

function buildToc() {
  toc.value = []
  if (!bodyEl.value) return
  const headings = bodyEl.value.querySelectorAll('h2, h3')
  headings.forEach((heading, index) => {
    const el = heading as HTMLElement
    if (!el.id) el.id = `heading-${index}`
    toc.value.push({
      id: el.id,
      text: el.textContent ?? '',
      level: el.tagName === 'H2' ? 2 : 3,
    })
  })
}

function highlightCode() {
  bodyEl.value?.querySelectorAll('pre code').forEach((block) => {
    hljs.highlightElement(block as HTMLElement)
  })
}

watch(() => route.params.slug, load)
onMounted(load)
</script>

<template>
  <div class="container-page">
    <div class="pt-8">
      <router-link to="/articles" class="btn-hand">← {{ t('nav.articles') }}</router-link>
    </div>

    <div v-if="loading" class="mt-10 space-y-4">
      <div class="h-10 w-2/3 animate-pulse rounded bg-white/60" />
      <div class="h-4 w-1/3 animate-pulse rounded bg-white/50" />
      <div class="h-64 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <EmptyState v-else-if="error" emoji="🕳️" :text="error" :hint="t('common.notFoundHint')" />

    <article v-else-if="article" class="mx-auto mt-8 max-w-prose">
      <header>
        <div class="flex flex-wrap items-center gap-2 text-xs text-ink-faint">
          <span v-if="article.is_pinned" class="rounded-full bg-terracotta-soft/40 px-2 py-0.5 text-terracotta-deep">
            📌 {{ t('admin.pinned') }}
          </span>
          <time>{{ date }}</time>
          <span v-if="article.category">· {{ article.category.name }}</span>
        </div>

        <h1 class="mt-3 font-hand text-3xl leading-tight text-ink sm:text-4xl">{{ article.title }}</h1>

        <p v-if="article.summary" class="mt-3 text-sm leading-7 text-ink-soft">{{ article.summary }}</p>

        <div v-if="article.tags?.length" class="mt-4 flex flex-wrap gap-2">
          <router-link
            v-for="tag in article.tags"
            :key="tag.id"
            :to="`/articles?tag=${tag.slug}`"
            class="rounded-full bg-paper-deep px-2.5 py-1 text-xs text-ink-soft transition hover:text-terracotta"
          >
            #{{ tag.name }}
          </router-link>
        </div>
      </header>

      <img
        v-if="article.cover_url"
        :src="article.cover_url"
        :alt="article.title"
        class="mt-8 w-full rounded-2xl border border-line object-cover"
      />

      <div class="mt-8 grid gap-10 xl:grid-cols-[1fr_180px]">
        <div ref="bodyEl" class="article-body min-w-0" v-html="html" />

        <aside v-if="toc.length" class="hidden xl:block">
          <div class="sticky top-24 border-l border-dashed border-line pl-4 text-sm">
            <p class="mb-2 font-hand text-ink-soft">目录</p>
            <a
              v-for="item in toc"
              :key="item.id"
              :href="`#${item.id}`"
              class="block py-1 leading-6 text-ink-faint transition hover:text-terracotta"
              :class="item.level === 3 ? 'pl-3 text-xs' : ''"
            >
              {{ item.text }}
            </a>
          </div>
        </aside>
      </div>

      <nav class="mt-14 grid gap-3 border-t border-dashed border-line pt-8 sm:grid-cols-2">
        <router-link
          v-if="data?.prev"
          :to="`/articles/${data.prev.slug}`"
          class="card-hand group p-4 transition hover:border-terracotta-soft"
        >
          <p class="text-xs text-ink-faint">← {{ t('common.prev') }}</p>
          <p class="mt-1 line-clamp-1 text-sm text-ink group-hover:text-terracotta-deep">{{ data.prev.title }}</p>
        </router-link>
        <div v-else />

        <router-link
          v-if="data?.next"
          :to="`/articles/${data.next.slug}`"
          class="card-hand group p-4 text-right transition hover:border-terracotta-soft"
        >
          <p class="text-xs text-ink-faint">{{ t('common.next') }} →</p>
          <p class="mt-1 line-clamp-1 text-sm text-ink group-hover:text-terracotta-deep">{{ data.next.title }}</p>
        </router-link>
      </nav>
    </article>
  </div>
</template>
