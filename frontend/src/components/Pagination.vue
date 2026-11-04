<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  page: number
  pageSize: number
  total: number
}>()

const emit = defineEmits<{ (e: 'update:page', page: number): void }>()

const pages = computed(() => Math.max(1, Math.ceil(props.total / props.pageSize)))

const visible = computed(() => {
  const total = pages.value
  const current = props.page
  const start = Math.max(1, Math.min(current - 2, total - 4))
  const end = Math.min(total, start + 4)
  const list: number[] = []
  for (let i = start; i <= end; i++) list.push(i)
  return list
})

function go(page: number) {
  if (page < 1 || page > pages.value || page === props.page) return
  emit('update:page', page)
}
</script>

<template>
  <div v-if="pages > 1" class="mt-8 flex items-center justify-center gap-1.5 text-sm">
    <button class="btn-hand !px-3" :disabled="page <= 1" @click="go(page - 1)">←</button>
    <button
      v-for="item in visible"
      :key="item"
      class="h-8 w-8 rounded-full border text-xs transition"
      :class="
        item === page
          ? 'border-terracotta bg-terracotta text-white'
          : 'border-line bg-white/60 text-ink-soft hover:border-terracotta-soft'
      "
      @click="go(item)"
    >
      {{ item }}
    </button>
    <button class="btn-hand !px-3" :disabled="page >= pages" @click="go(page + 1)">→</button>
  </div>
</template>
