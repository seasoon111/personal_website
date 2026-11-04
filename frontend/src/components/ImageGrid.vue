<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

const props = defineProps<{ urls: string[] }>()

const active = ref<string | null>(null)

const cols = computed(() => {
  const n = props.urls.length
  if (n === 1) return 'grid-cols-1'
  if (n === 2 || n === 4) return 'grid-cols-2'
  return 'grid-cols-3'
})

const ratio = computed(() => (props.urls.length === 1 ? 'aspect-[4/3]' : 'aspect-square'))

function onKey(event: KeyboardEvent) {
  if (event.key === 'Escape') active.value = null
}

onMounted(() => window.addEventListener('keydown', onKey))
onBeforeUnmount(() => window.removeEventListener('keydown', onKey))
</script>

<template>
  <div>
    <div class="grid gap-1.5" :class="cols">
      <button
        v-for="(url, index) in urls"
        :key="url"
        type="button"
        class="group relative overflow-hidden rounded-lg border border-line bg-paper-deep"
        :class="[ratio, urls.length === 1 ? 'max-h-[420px]' : '']"
        @click="active = url"
      >
        <img
          :src="url"
          :alt="`图片 ${index + 1}`"
          loading="lazy"
          class="h-full w-full object-cover transition duration-500 group-hover:scale-105"
        />
      </button>
    </div>

    <teleport to="body">
      <div
        v-if="active"
        class="fixed inset-0 z-50 flex items-center justify-center bg-ink/85 p-4 backdrop-blur-sm"
        @click.self="active = null"
      >
        <button
          type="button"
          class="absolute right-5 top-4 text-2xl text-paper/80 transition hover:text-paper"
          aria-label="close"
          @click="active = null"
        >
          ✕
        </button>
        <img :src="active" alt="" class="max-h-[88vh] max-w-full rounded-lg object-contain shadow-2xl" />
      </div>
    </teleport>
  </div>
</template>
