<script setup lang="ts">
defineProps<{
  title: string
  subtitle?: string
  to?: string
  linkLabel?: string
}>()
</script>

<template>
  <div class="mb-6 flex flex-wrap items-end justify-between gap-3">
    <div>
      <h2 class="dashed-title">{{ title }}</h2>
      <p v-if="subtitle" class="mt-2 text-sm text-ink-faint">{{ subtitle }}</p>
    </div>
    <router-link v-if="to" :to="to" class="btn-hand">
      {{ linkLabel ?? '查看全部' }}
      <span aria-hidden="true">→</span>
    </router-link>
  </div>
</template>
