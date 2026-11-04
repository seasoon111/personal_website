<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import MomentCard from '@/components/MomentCard.vue'
import PageHeader from '@/components/PageHeader.vue'
import Pagination from '@/components/Pagination.vue'
import type { Moment, PageResult } from '@/types'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()

const result = ref<PageResult<Moment> | null>(null)
const loading = ref(true)
const page = ref(Number(route.query.page ?? 1) || 1)

async function load() {
  loading.value = true
  try {
    result.value = await api.moments({ page: page.value, page_size: 10 })
  } finally {
    loading.value = false
  }
}

function changePage(p: number) {
  page.value = p
  router.push({ query: { page: String(p) } })
  load()
}

onMounted(load)
</script>

<template>
  <div class="container-page">
    <PageHeader :title="t('moments.title')" :subtitle="t('moments.subtitle')" emoji="💬" />

    <div v-if="loading" class="space-y-6">
      <div v-for="i in 3" :key="i" class="h-32 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <div v-else-if="result?.items?.length" class="space-y-7">
      <MomentCard v-for="(moment, index) in result.items" :key="moment.id" v-reveal="index * 50" :moment="moment" />
      <Pagination :page="result.page" :page-size="result.page_size" :total="result.total" @update:page="changePage" />
    </div>

    <EmptyState v-else :text="t('moments.empty')" :hint="t('common.emptyHint')" />
  </div>
</template>
