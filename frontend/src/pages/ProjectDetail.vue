<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import type { Project } from '@/types'
import { renderMarkdown } from '@/utils/markdown'

const route = useRoute()
const { t } = useI18n()
const project = ref<Project | null>(null)
const loading = ref(true)
const error = ref('')

const html = computed(() => project.value?.content_html || renderMarkdown(project.value?.content_md))
const tech = computed(() => project.value?.tech_stack ?? [])

async function load() {
  loading.value = true
  error.value = ''
  try {
    project.value = await api.project(route.params.slug as string)
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}

watch(() => route.params.slug, load)
onMounted(load)
</script>

<template>
  <div class="container-page">
    <div class="pt-8">
      <router-link to="/projects" class="btn-hand">← 项目</router-link>
    </div>

    <div v-if="loading" class="mt-10 space-y-4">
      <div class="h-10 w-1/2 animate-pulse rounded bg-white/60" />
      <div class="h-72 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <EmptyState v-else-if="error" emoji="🕳️" :text="error" :hint="t('common.notFoundHint')" />

    <article v-else-if="project" class="mx-auto mt-8 max-w-3xl">
      <h1 class="font-hand text-3xl text-ink sm:text-4xl">{{ project.title }}</h1>
      <p v-if="project.summary" class="mt-3 text-sm leading-7 text-ink-soft">{{ project.summary }}</p>

      <div class="mt-5 flex flex-wrap items-center gap-2">
        <span
          v-for="item in tech"
          :key="item"
          class="rounded-full bg-paper-deep px-3 py-1 text-xs text-ink-soft"
        >
          {{ item }}
        </span>
        <a v-if="project.repo_url" :href="project.repo_url" target="_blank" rel="noopener noreferrer" class="btn-hand">
          源码 ↗
        </a>
        <a v-if="project.demo_url" :href="project.demo_url" target="_blank" rel="noopener noreferrer" class="btn-primary !py-1.5">
          演示 ↗
        </a>
      </div>

      <img
        v-if="project.cover_url"
        :src="project.cover_url"
        :alt="project.title"
        class="mt-8 w-full rounded-2xl border border-line object-cover"
      />

      <div class="article-body mt-8" v-html="html" />
    </article>
  </div>
</template>
