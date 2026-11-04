<script setup lang="ts">
import { useI18n } from 'vue-i18n'

withDefaults(
  defineProps<{
    emoji?: string
    text?: string
    hint?: string
  }>(),
  { emoji: '🍂' },
)

const { t } = useI18n()
</script>

<template>
  <div class="flex flex-col items-center justify-center rounded-2xl border border-dashed border-line bg-white/40 px-6 py-16 text-center">
    <div class="text-4xl">{{ emoji }}</div>
    <p class="mt-4 font-hand text-lg text-ink-soft">{{ text ?? t('common.empty') }}</p>
    <p v-if="hint" class="mt-1 text-sm text-ink-faint">{{ hint }}</p>
    <slot />
  </div>
</template>
