<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import Pagination from '@/components/Pagination.vue'
import { useToast } from '@/composables/useToast'
import type { ImageItem, PageResult } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const result = ref<PageResult<ImageItem> | null>(null)
const loading = ref(true)
const uploading = ref(false)
const page = ref(1)

async function load() {
  loading.value = true
  try {
    result.value = await api.images({ page: page.value, page_size: 24 })
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

async function onUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  if (!files.length) return
  uploading.value = true
  try {
    await api.uploadImages(files)
    toast.success('上传成功')
    page.value = 1
    await load()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    uploading.value = false
    input.value = ''
  }
}

async function copy(url: string) {
  try {
    await navigator.clipboard.writeText(url)
    toast.success(t('admin.copied'))
  } catch {
    const textarea = document.createElement('textarea')
    textarea.value = url
    document.body.appendChild(textarea)
    textarea.select()
    document.execCommand('copy')
    textarea.remove()
    toast.success(t('admin.copied'))
  }
}

async function remove(image: ImageItem) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteImage(image.id)
    toast.success(t('admin.deleted'))
    await load()
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
    <div class="mb-2 flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-hand text-2xl">{{ t('admin.imageLibrary') }}</h1>
      <label class="btn-primary cursor-pointer text-sm">
        {{ uploading ? t('admin.uploading') : '+ ' + t('admin.upload') }}
        <input type="file" accept="image/*" multiple class="hidden" :disabled="uploading" @change="onUpload" />
      </label>
    </div>
    <p class="mb-6 text-sm text-ink-faint">{{ t('admin.imageLibraryHint') }}</p>

    <div v-if="loading" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div v-for="i in 8" :key="i" class="aspect-square animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else-if="result?.items?.length" class="grid grid-cols-2 gap-3 sm:grid-cols-4">
      <div
        v-for="image in result.items"
        :key="image.id"
        class="group overflow-hidden rounded-xl border border-line bg-white/70"
      >
        <img :src="image.url" :alt="image.original_name" class="aspect-square w-full object-cover" loading="lazy" />
        <div class="p-2">
          <p class="truncate text-xs text-ink-soft">{{ image.original_name }}</p>
          <p class="text-[11px] text-ink-faint">
            {{ formatDate(image.created_at, locale) }} · {{ Math.round(image.size / 1024) }} KB
          </p>
          <div class="mt-2 flex gap-2 text-xs">
            <button type="button" class="text-terracotta hover:underline" @click="copy(image.url)">
              {{ t('admin.copyUrl') }}
            </button>
            <button type="button" class="text-ink-faint hover:text-terracotta-deep" @click="remove(image)">
              {{ t('common.delete') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <p v-else class="py-16 text-center text-sm text-ink-faint">{{ t('admin.emptyImages') }}</p>

    <Pagination
      v-if="result"
      :page="result.page"
      :page-size="result.page_size"
      :total="result.total"
      @update:page="changePage"
    />
  </div>
</template>
