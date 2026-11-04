<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  emoji?: string
}>()
</script>

<template>
  <header class="mb-10 pt-10 sm:pt-14">
    <div class="flex items-center gap-3">
      <span v-if="emoji" class="text-2xl">{{ emoji }}</span>
      <h1 class="font-hand text-3xl text-ink sm:text-4xl">{{ title }}</h1>
    </div>
    <p v-if="subtitle" class="mt-3 max-w-prose text-sm leading-7 text-ink-soft">{{ subtitle }}</p>
    <div class="mt-6 h-[2px] w-24 bg-gradient-to-r from-terracotta to-transparent" />
  </header>
</template>
