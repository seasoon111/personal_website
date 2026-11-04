<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import ImagePicker from '@/components/ImagePicker.vue'
import Pagination from '@/components/Pagination.vue'
import { useToast } from '@/composables/useToast'
import type { ImageItem, Moment, PageResult } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const result = ref<PageResult<Moment> | null>(null)
const loading = ref(true)
const saving = ref(false)
const page = ref(1)
const pickerOpen = ref(false)
const editingId = ref<number | null>(null)

const form = ref({ content: '', is_public: true })
const selected = ref<ImageItem[]>([])

const selectedUrls = computed(() => selected.value.map((item) => item.url))

async function load() {
  loading.value = true
  try {
    result.value = await api.adminMoments({ page: page.value, page_size: 10 })
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

function onPicked(items: ImageItem[]) {
  selected.value = items
}

function removeSelected(id: number) {
  selected.value = selected.value.filter((item) => item.id !== id)
}

function resetForm() {
  editingId.value = null
  form.value = { content: '', is_public: true }
  selected.value = []
}

async function submit() {
  if (!form.value.content.trim() && !selected.value.length) {
    toast.error('写点什么，或者配一张图')
    return
  }
  saving.value = true
  try {
    const payload = {
      content: form.value.content.trim(),
      image_ids: selected.value.map((item) => item.id),
      is_public: form.value.is_public,
    }
    if (editingId.value) await api.updateMoment(editingId.value, payload)
    else await api.createMoment(payload)
    toast.success(t('admin.saved'))
    resetForm()
    await load()
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

function startEdit(moment: Moment) {
  editingId.value = moment.id
  form.value = { content: moment.content, is_public: moment.is_public }
  selected.value = (moment.images ?? [])
    .map((item) => item.image)
    .filter((image): image is ImageItem => Boolean(image))
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

async function remove(moment: Moment) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteMoment(moment.id)
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
  <div class="mx-auto max-w-3xl">
    <h1 class="mb-6 font-hand text-2xl">{{ t('admin.moments') }}</h1>

    <div class="rounded-2xl border border-line bg-white/70 p-5">
      <textarea
        v-model="form.content"
        rows="4"
        class="input-hand resize-y"
        :placeholder="t('admin.momentContent')"
        maxlength="2000"
      />

      <div v-if="selectedUrls.length" class="mt-3 grid grid-cols-3 gap-2">
        <div v-for="(url, index) in selectedUrls" :key="url" class="relative">
          <img :src="url" class="aspect-square w-full rounded-lg border border-line object-cover" />
          <button
            type="button"
            class="absolute right-1 top-1 rounded-full bg-ink/70 px-1.5 text-xs text-white"
            @click="removeSelected(selected[index].id)"
          >
            ✕
          </button>
        </div>
      </div>

      <div class="mt-3 flex flex-wrap items-center gap-3">
        <button type="button" class="btn-hand text-xs" @click="pickerOpen = true">
          🖼 {{ t('admin.chooseImages') }}
        </button>
        <label class="flex items-center gap-1.5 text-xs text-ink-soft">
          <input v-model="form.is_public" type="checkbox" class="accent-terracotta" />
          {{ t('admin.momentPublic') }}
        </label>
        <span v-if="selectedUrls.length" class="text-xs text-ink-faint">
          {{ t('admin.selected', { count: selectedUrls.length }) }}
        </span>
        <div class="ml-auto flex gap-2">
          <button v-if="editingId" type="button" class="btn-hand text-xs" @click="resetForm">
            {{ t('common.cancel') }}
          </button>
          <button type="button" class="btn-primary !py-1.5 text-xs" :disabled="saving" @click="submit">
            {{ editingId ? t('common.save') : t('admin.newMoment') }}
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="mt-6 space-y-3">
      <div v-for="i in 3" :key="i" class="h-24 animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else-if="result?.items?.length" class="mt-6 space-y-3">
      <div v-for="moment in result.items" :key="moment.id" class="rounded-xl border border-line bg-white/70 p-4">
        <div class="flex items-center justify-between text-xs text-ink-faint">
          <time>{{ formatDate(moment.created_at, locale) }}</time>
          <div class="flex gap-2">
            <button type="button" class="hover:text-terracotta" @click="startEdit(moment)">{{ t('common.edit') }}</button>
            <button type="button" class="hover:text-terracotta-deep" @click="remove(moment)">{{ t('common.delete') }}</button>
          </div>
        </div>
        <p class="mt-2 whitespace-pre-wrap text-sm leading-6 text-ink/90">{{ moment.content }}</p>
        <div v-if="moment.images?.length" class="mt-2 flex gap-1.5">
          <img
            v-for="image in moment.images"
            :key="image.id"
            :src="image.image?.url"
            class="h-16 w-16 rounded-lg border border-line object-cover"
          />
        </div>
        <p v-if="!moment.is_public" class="mt-2 text-xs text-terracotta">仅自己可见</p>
      </div>

      <Pagination :page="result.page" :page-size="result.page_size" :total="result.total" @update:page="changePage" />
    </div>

    <p v-else class="mt-10 text-center text-sm text-ink-faint">还没有动态。</p>

    <ImagePicker v-model:open="pickerOpen" multiple :max="9" @confirm="onPicked" />
  </div>
</template>
