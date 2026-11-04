<script setup lang="ts">
import { computed } from 'vue'
import type { Project } from '@/types'

const props = defineProps<{ project: Project }>()

const tech = computed(() => (props.project.tech_stack ?? []).slice(0, 5))
</script>

<template>
  <router-link
    :to="`/projects/${project.slug}`"
    class="group flex flex-col overflow-hidden rounded-2xl border border-line bg-white/70 shadow-hand transition duration-300 hover:-translate-y-1 hover:border-terracotta-soft hover:shadow-hand-lg"
  >
    <div class="relative aspect-[16/10] overflow-hidden bg-paper-deep">
      <img
        v-if="project.cover_url"
        :src="project.cover_url"
        :alt="project.title"
        loading="lazy"
        class="h-full w-full object-cover transition duration-500 group-hover:scale-[1.03]"
      />
      <div v-else class="flex h-full w-full items-center justify-center">
        <span class="font-hand text-4xl text-terracotta/50">{{ project.title.slice(0, 1) }}</span>
      </div>
    </div>

    <div class="flex flex-1 flex-col p-5">
      <h3 class="font-hand text-xl text-ink transition group-hover:text-terracotta-deep">{{ project.title }}</h3>
      <p v-if="project.summary" class="mt-2 line-clamp-2 text-sm leading-6 text-ink-soft">{{ project.summary }}</p>

      <div v-if="tech.length" class="mt-auto flex flex-wrap gap-1.5 pt-4">
        <span
          v-for="item in tech"
          :key="item"
          class="rounded-full bg-paper-deep px-2.5 py-0.5 text-xs text-ink-soft"
        >
          {{ item }}
        </span>
      </div>
    </div>
  </router-link>
</template>
