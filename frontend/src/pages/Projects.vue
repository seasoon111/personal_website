<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import PageHeader from '@/components/PageHeader.vue'
import ProjectCard from '@/components/ProjectCard.vue'
import type { Project } from '@/types'

const { t } = useI18n()
const projects = ref<Project[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    projects.value = await api.projects()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container-page">
    <PageHeader :title="t('projects.title')" :subtitle="t('projects.subtitle')" emoji="🧩" />

    <div v-if="loading" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-72 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <div v-else-if="projects.length" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
      <ProjectCard v-for="(project, index) in projects" :key="project.id" v-reveal="index * 60" :project="project" />
    </div>

    <EmptyState v-else :text="t('projects.empty')" :hint="t('common.emptyHint')" />
  </div>
</template>
