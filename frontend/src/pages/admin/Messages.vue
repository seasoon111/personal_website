<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import Pagination from '@/components/Pagination.vue'
import { useToast } from '@/composables/useToast'
import type { ContactMessage, PageResult } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const result = ref<PageResult<ContactMessage> | null>(null)
const loading = ref(true)
const unreadOnly = ref(false)
const page = ref(1)

async function load() {
  loading.value = true
  try {
    result.value = await api.messages({
      page: page.value,
      page_size: 10,
      unread: unreadOnly.value ? 1 : undefined,
    })
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

async function toggleRead(message: ContactMessage) {
  try {
    await api.setMessageRead(message.id, !message.is_read)
    message.is_read = !message.is_read
    if (unreadOnly.value && message.is_read) await load()
  } catch (e) {
    toast.error((e as Error).message)
  }
}

async function remove(message: ContactMessage) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteMessage(message.id)
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
  <div class="mx-auto max-w-4xl">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-hand text-2xl">{{ t('admin.messages') }}</h1>
      <label class="flex items-center gap-2 text-sm text-ink-soft">
        <input v-model="unreadOnly" type="checkbox" class="accent-terracotta" @change="load" />
        {{ t('admin.unreadOnly') }}
      </label>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else-if="result?.items?.length" class="space-y-3">
      <article
        v-for="message in result.items"
        :key="message.id"
        class="rounded-xl border bg-white/70 p-4"
        :class="message.is_read ? 'border-line' : 'border-terracotta-soft'"
      >
        <div class="flex flex-wrap items-center gap-2 text-xs text-ink-faint">
          <span v-if="!message.is_read" class="rounded-full bg-terracotta-soft/40 px-2 py-0.5 text-terracotta-deep">
            NEW
          </span>
          <span class="font-medium text-ink">{{ message.name }}</span>
          <a :href="`mailto:${message.email}`" class="underline decoration-dashed underline-offset-4 hover:text-terracotta">
            {{ message.email }}
          </a>
          <time class="ml-auto">{{ formatDate(message.created_at, locale) }}</time>
        </div>

        <p v-if="message.subject" class="mt-2 font-hand text-lg text-ink">{{ message.subject }}</p>
        <p class="mt-2 whitespace-pre-wrap text-sm leading-6 text-ink/90">{{ message.content }}</p>

        <div class="mt-3 flex gap-3 text-xs">
          <button type="button" class="text-terracotta hover:underline" @click="toggleRead(message)">
            {{ message.is_read ? t('admin.markUnread') : t('admin.markRead') }}
          </button>
          <button type="button" class="text-ink-faint hover:text-terracotta-deep" @click="remove(message)">
            {{ t('common.delete') }}
          </button>
        </div>
      </article>

      <Pagination :page="result.page" :page-size="result.page_size" :total="result.total" @update:page="changePage" />
    </div>

    <p v-else class="py-16 text-center text-sm text-ink-faint">{{ t('admin.noMessages') }}</p>
  </div>
</template>
