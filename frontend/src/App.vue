<script setup lang="ts">
import { onMounted } from 'vue'
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()

onMounted(() => {
  site.load()
})
</script>

<template>
  <router-view />
</template>
