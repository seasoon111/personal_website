<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import { useToast } from '@/composables/useToast'
import Pagination from '@/components/Pagination.vue'
import type { PageResult, Project } from '@/types'
import { formatDate } from '@/utils/date'

const { t, locale } = useI18n()
const toast = useToast()

const result = ref<PageResult<Project> | null>(null)
const loading = ref(true)
const page = ref(1)

async function load() {
  loading.value = true
  try {
    result.value = await api.adminProjects({ page: page.value, page_size: 10 })
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
}

async function remove(project: Project) {
  if (!window.confirm(t('admin.confirmDelete'))) return
  try {
    await api.deleteProject(project.id)
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
  <div class="mx-auto max-w-5xl">
    <div class="mb-6 flex flex-wrap items-center justify-between gap-3">
      <h1 class="font-hand text-2xl">{{ t('admin.projects') }}</h1>
      <router-link to="/admin/projects/new" class="btn-primary">+ {{ t('admin.newProject') }}</router-link>
    </div>

    <div v-if="loading" class="space-y-3">
      <div v-for="i in 3" :key="i" class="h-20 animate-pulse rounded-xl bg-white/60" />
    </div>

    <div v-else-if="result?.items?.length" class="space-y-2">
      <div
        v-for="project in result.items"
        :key="project.id"
        class="flex flex-col gap-3 rounded-xl border border-line bg-white/70 p-4 sm:flex-row sm:items-center"
      >
        <img
          v-if="project.cover_url"
          :src="project.cover_url"
          class="hidden h-14 w-20 rounded-lg border border-line object-cover sm:block"
        />
        <div class="min-w-0 flex-1">
          <div class="flex items-center gap-2 text-xs text-ink-faint">
            <span
              class="rounded-full px-2 py-0.5"
              :class="project.status === 'published' ? 'bg-sage/20 text-ink' : 'bg-paper-deep text-ink-soft'"
            >
              {{ project.status === 'published' ? t('admin.published') : t('admin.draft') }}
            </span>
            <time>{{ formatDate(project.updated_at, locale) }}</time>
            <span>排序 {{ project.sort }}</span>
          </div>
          <p class="mt-1 truncate font-hand text-lg">{{ project.title }}</p>
        </div>
        <div class="flex gap-2 text-xs">
          <router-link :to="`/admin/projects/${project.id}/edit`" class="btn-hand">{{ t('common.edit') }}</router-link>
          <button type="button" class="btn-hand !text-terracotta-deep" @click="remove(project)">
            {{ t('common.delete') }}
          </button>
        </div>
      </div>

      <Pagination :page="result.page" :page-size="result.page_size" :total="result.total" @update:page="changePage" />
    </div>

    <p v-else class="py-16 text-center text-sm text-ink-faint">还没有项目。</p>
  </div>
</template>
