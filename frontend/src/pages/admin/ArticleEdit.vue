<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import ImagePicker from '@/components/ImagePicker.vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useToast } from '@/composables/useToast'
import type { Category, ImageItem, Tag } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const toast = useToast()

const id = computed(() => (route.params.id ? Number(route.params.id) : null))
const isEdit = computed(() => id.value !== null)

const form = reactive({
  title: '',
  slug: '',
  summary: '',
  content_md: '',
  cover_url: '',
  category_id: null as number | null,
  tag_ids: [] as number[],
  status: 'draft' as 'draft' | 'published',
  is_pinned: false,
})

const categories = ref<Category[]>([])
const tags = ref<Tag[]>([])
const loading = ref(false)
const saving = ref(false)
const editor = ref<InstanceType<typeof MarkdownEditor> | null>(null)
const coverPicker = ref(false)
const contentPicker = ref(false)

function toggleTag(tagId: number) {
  const index = form.tag_ids.indexOf(tagId)
  if (index >= 0) form.tag_ids.splice(index, 1)
  else form.tag_ids.push(tagId)
}

function onCoverPicked(items: ImageItem[]) {
  if (items[0]) form.cover_url = items[0].url
}

function onContentPicked(items: ImageItem[]) {
  const item = items[0]
  if (!item) return
  editor.value?.insert(`\n![${item.original_name}](${item.url})\n`)
}

async function save(nextStatus?: 'draft' | 'published') {
  if (!form.title.trim()) {
    toast.error('标题不能为空')
    return
  }
  saving.value = true
  try {
    const payload = {
      title: form.title.trim(),
      slug: form.slug.trim(),
      summary: form.summary.trim(),
      content_md: form.content_md,
      cover_url: form.cover_url,
      category_id: form.category_id,
      tag_ids: form.tag_ids,
      status: nextStatus ?? form.status,
      is_pinned: form.is_pinned,
    }
    if (id.value) await api.updateArticle(id.value, payload)
    else await api.createArticle(payload)
    toast.success(t('admin.saved'))
    router.push('/admin/articles')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  loading.value = true
  try {
    const [cats, tagList] = await Promise.all([api.categories(), api.tags()])
    categories.value = cats
    tags.value = tagList

    if (id.value) {
      const article = await api.adminArticle(id.value)
      form.title = article.title
      form.slug = article.slug
      form.summary = article.summary
      form.content_md = article.content_md ?? ''
      form.cover_url = article.cover_url
      form.category_id = article.category_id
      form.tag_ids = (article.tags ?? []).map((tag) => tag.id)
      form.status = article.status
      form.is_pinned = article.is_pinned
    }
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="mx-auto max-w-5xl">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
      <div class="flex items-center gap-3">
        <router-link to="/admin/articles" class="btn-hand">←</router-link>
        <h1 class="font-hand text-2xl">{{ isEdit ? t('common.edit') : t('admin.newArticle') }}</h1>
      </div>
      <div class="flex gap-2">
        <button type="button" class="btn-hand" :disabled="saving" @click="save('draft')">
          {{ t('admin.draft') }}
        </button>
        <button type="button" class="btn-primary" :disabled="saving" @click="save('published')">
          {{ saving ? t('common.saving') : t('admin.publish') }}
        </button>
      </div>
    </div>

    <div v-if="loading" class="h-64 animate-pulse rounded-2xl bg-white/60" />

    <div v-else class="grid gap-6 lg:grid-cols-[1fr_300px]">
      <div class="space-y-4">
        <input
          v-model="form.title"
          class="input-hand !py-3 font-hand text-xl"
          placeholder="文章标题"
          maxlength="200"
        />
        <textarea
          v-model="form.summary"
          rows="2"
          class="input-hand resize-y"
          placeholder="摘要（列表页显示，可不填）"
          maxlength="500"
        />

        <div class="flex items-center gap-2">
          <button type="button" class="btn-hand text-xs" @click="contentPicker = true">🖼 插入图片</button>
          <span class="text-xs text-ink-faint">支持 Markdown 语法，右侧可切换预览</span>
        </div>

        <MarkdownEditor ref="editor" v-model="form.content_md" :rows="22" placeholder="开始写…" />
      </div>

      <aside class="space-y-4">
        <div class="rounded-xl border border-line bg-white/70 p-4">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.cover') }}</label>
          <div v-if="form.cover_url" class="relative">
            <img :src="form.cover_url" class="aspect-[16/10] w-full rounded-lg border border-line object-cover" />
            <button
              type="button"
              class="absolute right-1 top-1 rounded-full bg-ink/70 px-2 py-0.5 text-xs text-white"
              @click="form.cover_url = ''"
            >
              ✕
            </button>
          </div>
          <button type="button" class="btn-hand mt-2 w-full justify-center text-xs" @click="coverPicker = true">
            {{ t('admin.selectCover') }}
          </button>
        </div>

        <div class="rounded-xl border border-line bg-white/70 p-4">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.category') }}</label>
          <select v-model="form.category_id" class="input-hand">
            <option :value="null">未分类</option>
            <option v-for="category in categories" :key="category.id" :value="category.id">
              {{ category.name }}
            </option>
          </select>
        </div>

        <div class="rounded-xl border border-line bg-white/70 p-4">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.tagsLabel') }}</label>
          <div class="flex flex-wrap gap-1.5">
            <button
              v-for="tag in tags"
              :key="tag.id"
              type="button"
              class="rounded-full border px-2.5 py-1 text-xs transition"
              :class="
                form.tag_ids.includes(tag.id)
                  ? 'border-terracotta bg-terracotta text-white'
                  : 'border-line bg-white text-ink-soft hover:border-terracotta-soft'
              "
              @click="toggleTag(tag.id)"
            >
              {{ tag.name }}
            </button>
            <span v-if="!tags.length" class="text-xs text-ink-faint">还没有标签，去「标签」页新建</span>
          </div>
        </div>

        <div class="rounded-xl border border-line bg-white/70 p-4 text-sm">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.slug') }}</label>
          <input v-model="form.slug" class="input-hand text-xs" placeholder="留空自动生成" />

          <label class="mt-3 flex items-center gap-2 text-sm">
            <input v-model="form.is_pinned" type="checkbox" class="accent-terracotta" />
            {{ t('admin.pinned') }}
          </label>
        </div>
      </aside>
    </div>

    <ImagePicker v-model:open="coverPicker" @confirm="onCoverPicked" />
    <ImagePicker v-model:open="contentPicker" @confirm="onContentPicked" />
  </div>
</template>
