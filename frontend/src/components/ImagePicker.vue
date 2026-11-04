<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import api from '@/api'
import type { ImageItem } from '@/types'
import { useToast } from '@/composables/useToast'

const props = withDefaults(
  defineProps<{
    open: boolean
    multiple?: boolean
    max?: number
  }>(),
  { multiple: false, max: 9 },
)

const emit = defineEmits<{
  (e: 'update:open', value: boolean): void
  (e: 'confirm', items: ImageItem[]): void
}>()

const toast = useToast()
const images = ref<ImageItem[]>([])
const selected = ref<ImageItem[]>([])
const loading = ref(false)
const uploading = ref(false)

const canConfirm = computed(() => selected.value.length > 0)

async function load() {
  loading.value = true
  try {
    const res = await api.images({ page: 1, page_size: 60 })
    images.value = res.items ?? []
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

function close() {
  emit('update:open', false)
}

function pick(item: ImageItem) {
  if (!props.multiple) {
    emit('confirm', [item])
    close()
    return
  }
  const index = selected.value.findIndex((s) => s.id === item.id)
  if (index >= 0) {
    selected.value.splice(index, 1)
    return
  }
  if (selected.value.length >= props.max) {
    toast.error(`最多选择 ${props.max} 张`)
    return
  }
  selected.value.push(item)
}

function isSelected(item: ImageItem) {
  return selected.value.some((s) => s.id === item.id)
}

async function onUpload(event: Event) {
  const input = event.target as HTMLInputElement
  const files = Array.from(input.files ?? [])
  if (!files.length) return
  uploading.value = true
  try {
    const uploaded = await api.uploadImages(files)
    images.value = [...uploaded, ...images.value]
    toast.success('上传成功')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    uploading.value = false
    input.value = ''
  }
}

function confirm() {
  emit('confirm', [...selected.value])
  close()
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      selected.value = []
      load()
    }
  },
)
</script>

<template>
  <teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center bg-ink/70 p-4" @click.self="close">
      <div class="flex max-h-[86vh] w-full max-w-3xl flex-col overflow-hidden rounded-2xl border border-line bg-paper shadow-2xl">
        <header class="flex items-center justify-between border-b border-line px-5 py-3">
          <h3 class="font-hand text-lg">图片库</h3>
          <div class="flex items-center gap-2">
            <label class="btn-hand cursor-pointer text-xs">
              {{ uploading ? '上传中…' : '上传图片' }}
              <input type="file" accept="image/*" multiple class="hidden" :disabled="uploading" @change="onUpload" />
            </label>
            <button type="button" class="text-lg text-ink-faint hover:text-ink" @click="close">✕</button>
          </div>
        </header>

        <div class="flex-1 overflow-y-auto p-4">
          <div v-if="loading" class="grid grid-cols-3 gap-3 sm:grid-cols-4">
            <div v-for="i in 8" :key="i" class="aspect-square animate-pulse rounded-lg bg-white/60" />
          </div>

          <div v-else-if="images.length" class="grid grid-cols-3 gap-3 sm:grid-cols-4">
            <button
              v-for="image in images"
              :key="image.id"
              type="button"
              class="group relative aspect-square overflow-hidden rounded-lg border-2 transition"
              :class="isSelected(image) ? 'border-terracotta' : 'border-transparent hover:border-terracotta-soft'"
              @click="pick(image)"
            >
              <img :src="image.url" :alt="image.original_name" class="h-full w-full object-cover" loading="lazy" />
              <span
                v-if="isSelected(image)"
                class="absolute right-1 top-1 flex h-5 w-5 items-center justify-center rounded-full bg-terracotta text-xs text-white"
              >
                ✓
              </span>
            </button>
          </div>

          <p v-else class="py-16 text-center text-sm text-ink-faint">图片库还是空的，先上传几张吧</p>
        </div>

        <footer v-if="multiple" class="flex items-center justify-between border-t border-line px-5 py-3 text-sm">
          <span class="text-ink-faint">已选 {{ selected.length }} / {{ max }}</span>
          <div class="flex gap-2">
            <button type="button" class="btn-hand" @click="close">取消</button>
            <button type="button" class="btn-primary" :disabled="!canConfirm" @click="confirm">确定</button>
          </div>
        </footer>
      </div>
    </div>
  </teleport>
</template>
