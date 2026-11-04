<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import type { Moment } from '@/types'
import ImageGrid from '@/components/ImageGrid.vue'

const props = defineProps<{ moment: Moment }>()
const { locale } = useI18n()

const date = new Date(props.moment.created_at)
const day = computed(() => String(date.getDate()).padStart(2, '0'))
const month = computed(() =>
  locale.value === 'en-US'
    ? date.toLocaleDateString('en-US', { month: 'short' })
    : `${date.getMonth() + 1} 月`,
)
const urls = computed(() =>
  (props.moment.images ?? [])
    .map((item) => item.image?.url)
    .filter((url): url is string => Boolean(url)),
)
</script>

<template>
  <article class="flex gap-4 sm:gap-6">
    <div class="w-12 shrink-0 pt-1 text-center sm:w-16">
      <div class="font-hand text-2xl leading-none text-terracotta sm:text-3xl">{{ day }}</div>
      <div class="mt-1 text-xs text-ink-faint">{{ month }}</div>
    </div>

    <div class="min-w-0 flex-1 rounded-2xl border border-line bg-white/70 p-5 shadow-hand">
      <p v-if="moment.content" class="whitespace-pre-wrap break-words text-sm leading-7 text-ink/90">
        {{ moment.content }}
      </p>
      <ImageGrid v-if="urls.length" :urls="urls" class="mt-3" />
      <div class="mt-3 flex items-center justify-between text-xs text-ink-faint">
        <time>{{ new Date(moment.created_at).toLocaleTimeString(locale, { hour: '2-digit', minute: '2-digit' }) }}</time>
        <span v-if="!moment.is_public" class="rounded-full bg-paper-deep px-2 py-0.5">仅自己可见</span>
      </div>
    </div>
  </article>
</template>
