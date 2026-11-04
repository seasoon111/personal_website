<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import type { Category } from '@/types'

const { t } = useI18n()
const toast = useToast()

const categories = ref<Category[]>([])
const loading = ref(true)
const saving = ref(false)
const editId = ref<number | null>(null)

const form = reactive({ name: '', slug: '', sort: 0 })

async function load() {
  loading.value = true
  try {
    categories.value = await api.categories()
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
  form.sort = 0
}

function startEdit(category: Category) {
  editId.value = category.id
  form.name = category.name
  form.slug = category.slug
  form.sort = category.sort
}

async function submit() {
  if (!form.name.trim()) {
    toast.error('分类名称不能为空')
    return
  }
  saving.value = true
  try {
    if (editId.value) await api.updateCategory(editId.value, { ...form })
    else await api.createCategory({ ...form })
    toast.success(t('admin.saved'))
    reset()
    await load()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

async function remove(category: Category) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteCategory(category.id)
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
    <h1 class="mb-6 font-hand text-2xl">{{ t('admin.categories') }}</h1>

    <div class="rounded-2xl border border-line bg-white/70 p-5">
      <div class="grid gap-3 sm:grid-cols-[1fr_1fr_80px_auto]">
        <input v-model="form.name" class="input-hand" placeholder="分类名称" maxlength="64" />
        <input v-model="form.slug" class="input-hand" placeholder="slug（可留空）" maxlength="64" />
        <input v-model.number="form.sort" type="number" class="input-hand" placeholder="排序" />
        <div class="flex gap-2">
          <button type="button" class="btn-primary flex-1" :disabled="saving" @click="submit">
            {{ editId ? t('common.save') : t('common.create') }}
          </button>
          <button v-if="editId" type="button" class="btn-hand" @click="reset">{{ t('common.cancel') }}</button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="mt-5 space-y-2">
      <div v-for="i in 3" :key="i" class="h-14 animate-pulse rounded-xl bg-white/60" />
    </div>

    <ul v-else class="mt-5 space-y-2">
      <li
        v-for="category in categories"
        :key="category.id"
        class="flex items-center gap-3 rounded-xl border border-line bg-white/70 px-4 py-3"
      >
        <span class="font-hand text-lg text-ink">{{ category.name }}</span>
        <span class="text-xs text-ink-faint">/{{ category.slug }}</span>
        <span class="text-xs text-ink-faint">{{ category.article_count }} 篇</span>
        <span class="ml-auto flex gap-2 text-xs">
          <button type="button" class="btn-hand" @click="startEdit(category)">{{ t('common.edit') }}</button>
          <button type="button" class="btn-hand !text-terracotta-deep" @click="remove(category)">{{ t('common.delete') }}</button>
        </span>
      </li>
    </ul>

    <p v-if="!loading && !categories.length" class="mt-8 text-center text-sm text-ink-faint">
      还没有分类，先建一个吧。
    </p>
  </div>
</template>
