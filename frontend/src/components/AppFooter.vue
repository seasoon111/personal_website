<script setup lang="ts">
import { useSiteStore } from '@/stores/site'

const site = useSiteStore()
const year = new Date().getFullYear()
</script>

<template>
  <footer class="mt-20 border-t border-dashed border-line py-10">
    <div class="container-page flex flex-col items-center gap-4 text-center text-sm text-ink-faint">
      <div v-if="site.info.socials?.length" class="flex flex-wrap justify-center gap-4">
        <a
          v-for="social in site.info.socials"
          :key="social.name"
          :href="social.url"
          target="_blank"
          rel="noopener noreferrer"
          class="text-ink-soft underline decoration-line decoration-dashed underline-offset-4 transition hover:text-terracotta"
        >
          {{ social.name }}
        </a>
      </div>

      <p>© {{ year }} {{ site.info.name }} · 用 Go + Vue 慢慢搭起来的</p>

      <p v-if="site.info.icp" class="text-xs">
        <a href="https://beian.miit.gov.cn/" target="_blank" rel="noopener noreferrer" class="hover:text-terracotta">
          {{ site.info.icp }}
        </a>
      </p>
    </div>
  </footer>
</template>
