<script setup lang="ts">
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
</script>

<template>
  <div class="container-page flex min-h-[60vh] flex-col items-center justify-center text-center">
    <div class="text-6xl">🍃</div>
    <h1 class="mt-6 font-hand text-3xl text-ink">{{ t('common.notFound') }}</h1>
    <p class="mt-3 text-sm text-ink-soft">{{ t('common.notFoundHint') }}</p>
    <router-link to="/" class="btn-primary mt-8">{{ t('common.backHome') }}</router-link>
  </div>
</template>
