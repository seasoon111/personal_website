<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import type { Tag } from '@/types'

const { t } = useI18n()
const toast = useToast()

const tags = ref<Tag[]>([])
const loading = ref(true)
const saving = ref(false)
const editId = ref<number | null>(null)

const form = reactive({ name: '', slug: '' })

async function load() {
  loading.value = true
  try {
    tags.value = await api.tags()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

function reset() {
  editId.value = null
  form.name = ''
  form.slug = ''
}

function startEdit(tag: Tag) {
  editId.value = tag.id
  form.name = tag.name
  form.slug = tag.slug
}

async function submit() {
  if (!form.name.trim()) {
    toast.error('标签名称不能为空')
    return
  }
  saving.value = true
  try {
    if (editId.value) await api.updateTag(editId.value, { ...form })
    else await api.createTag({ ...form })
    toast.success(t('admin.saved'))
    reset()
    await load()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(tag: Tag) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteTag(tag.id)
    toast.success(t('admin.deleted'))
    await load()
  } catch (e) {
    toast.error((e as Error).message)
  }
}

onMounted(load)
</script>

<template>
  <div class="mx-auto max-w-3xl">
    <h1 class="mb-6 font-hand text-2xl">{{ t('admin.tags') }}</h1>

    <div class="rounded-2xl border border-line bg-white/70 p-5">
      <div class="grid gap-3 sm:grid-cols-[1fr_1fr_auto]">
        <input v-model="form.name" class="input-hand" placeholder="标签名称" maxlength="64" />
        <input v-model="form.slug" class="input-hand" placeholder="slug（可留空）" maxlength="64" />
        <div class="flex gap-2">
          <button type="button" class="btn-primary" :disabled="saving" @click="submit">
            {{ editId ? t('common.save') : t('common.create') }}
          </button>
          <button v-if="editId" type="button" class="btn-hand" @click="reset">{{ t('common.cancel') }}</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="mt-5 space-y-2">
      <div v-for="i in 4" :key="i" class="h-12 animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else class="mt-5 flex flex-wrap gap-2">
      <div
        v-for="tag in tags"
        :key="tag.id"
        class="flex items-center gap-2 rounded-full border border-line bg-white/70 px-3 py-1.5 text-sm"
      >
        <span class="text-ink">#{{ tag.name }}</span>
        <span class="text-xs text-ink-faint">{{ tag.article_count }}</span>
        <button type="button" class="text-xs text-ink-faint hover:text-terracotta" @click="startEdit(tag)">
          {{ t('common.edit') }}
        </button>
        <button type="button" class="text-xs text-ink-faint hover:text-terracotta-deep" @click="remove(tag)">
          ✕
        </button>
      </div>
    </div>

    <p v-if="!loading && !tags.length" class="mt-8 text-center text-sm text-ink-faint">还没有标签。</p>
  </div>
</template>
