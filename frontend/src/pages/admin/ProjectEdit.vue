<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import ImagePicker from '@/components/ImagePicker.vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useToast } from '@/composables/useToast'
import type { ImageItem } from '@/types'

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
  tech_text: '',
  repo_url: '',
  demo_url: '',
  sort: 0,
  status: 'draft' as 'draft' | 'published',
})

const loading = ref(false)
const saving = ref(false)
const coverPicker = ref(false)

function onCoverPicked(items: ImageItem[]) {
  if (items[0]) form.cover_url = items[0].url
}

async function save() {
  if (!form.title.trim()) {
    toast.error('项目名称不能为空')
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
      tech_stack: form.tech_text
        .split(/[,，\n]/)
        .map((item) => item.trim())
        .filter(Boolean),
      repo_url: form.repo_url.trim(),
      demo_url: form.demo_url.trim(),
      sort: Number(form.sort) || 0,
      status: form.status,
    }
    if (id.value) await api.updateProject(id.value, payload)
    else await api.createProject(payload)
    toast.success(t('admin.saved'))
    router.push('/admin/projects')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  if (!id.value) return
  loading.value = true
  try {
    const project = await api.adminProject(id.value)
    form.title = project.title
    form.slug = project.slug
    form.summary = project.summary
    form.content_md = project.content_md ?? ''
    form.cover_url = project.cover_url
    form.tech_text = (project.tech_stack ?? []).join(', ')
    form.repo_url = project.repo_url
    form.demo_url = project.demo_url
    form.sort = project.sort
    form.status = project.status
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
        <router-link to="/admin/projects" class="btn-hand">←</router-link>
        <h1 class="font-hand text-2xl">{{ isEdit ? t('common.edit') : t('admin.newProject') }}</h1>
      </div>
      <button type="button" class="btn-primary" :disabled="saving" @click="save">
        {{ saving ? t('common.saving') : t('common.save') }}
      </button>
    </div>

    <div v-if="loading" class="h-64 animate-pulse rounded-2xl bg-white/60" />

    <div v-else class="grid gap-6 lg:grid-cols-[1fr_300px]">
      <div class="space-y-4">
        <input v-model="form.title" class="input-hand !py-3 font-hand text-xl" placeholder="项目名称" maxlength="200" />
        <textarea
          v-model="form.summary"
          rows="2"
          class="input-hand resize-y"
          placeholder="一句话介绍（卡片上显示）"
          maxlength="500"
        />
        <MarkdownEditor v-model="form.content_md" :rows="16" placeholder="项目详情…" />
      </div>

      <aside class="space-y-4">
        <div class="rounded-xl border border-line bg-white/70 p-4">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.cover') }}</label>
          <img
            v-if="form.cover_url"
            :src="form.cover_url"
            class="aspect-[16/10] w-full rounded-lg border border-line object-cover"
          />
          <button type="button" class="btn-hand mt-2 w-full justify-center text-xs" @click="coverPicker = true">
            {{ t('admin.selectCover') }}
          </button>
          <button
            v-if="form.cover_url"
            type="button"
            class="mt-1 w-full text-center text-xs text-ink-faint hover:text-terracotta"
            @click="form.cover_url = ''"
          >
            移除封面
          </button>
        </div>

        <div class="space-y-3 rounded-xl border border-line bg-white/70 p-4 text-sm">
          <div>
            <label class="mb-1 block text-xs text-ink-faint">{{ t('admin.status') }}</label>
            <select v-model="form.status" class="input-hand">
              <option value="draft">{{ t('admin.draft') }}</option>
              <option value="published">{{ t('admin.published') }}</option>
            </select>
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">技术栈（逗号分隔）</label>
            <input v-model="form.tech_text" class="input-hand" placeholder="Go, Vue3, MySQL" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">源码地址</label>
            <input v-model="form.repo_url" class="input-hand" placeholder="https://github.com/…" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">演示地址</label>
            <input v-model="form.demo_url" class="input-hand" placeholder="https://…" />
          </div>
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="mb-1 block text-xs text-ink-faint">排序</label>
              <input v-model.number="form.sort" type="number" class="input-hand" />
            </div>
            <div>
              <label class="mb-1 block text-xs text-ink-faint">{{ t('admin.slug') }}</label>
              <input v-model="form.slug" class="input-hand text-xs" placeholder="留空自动" />
            </div>
          </div>
        </div>
      </aside>
    </div>

    <ImagePicker v-model:open="coverPicker" @confirm="onCoverPicked" />
  </div>
</template>
