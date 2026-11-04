<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const auth = useAuthStore()

const form = ref({ username: 'admin', password: '' })
const loading = ref(false)
const error = ref('')

async function submit() {
  error.value = ''
  if (!form.value.username.trim() || !form.value.password) {
    error.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  try {
    await auth.login(form.value.username.trim(), form.value.password)
    router.replace((route.query.redirect as string) || '/admin')
  } catch (e) {
    error.value = (e as Error).message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="flex min-h-screen items-center justify-center bg-paper px-4">
    <div class="w-full max-w-sm">
      <div class="rounded-2xl border border-line bg-white/80 p-7 shadow-hand-lg">
        <div class="text-center">
          <div class="text-3xl">📔</div>
          <h1 class="mt-3 font-hand text-2xl text-ink">{{ t('admin.login') }}</h1>
          <p class="mt-1 text-xs text-ink-faint">{{ t('admin.loginHint') }}</p>
        </div>

        <form class="mt-7 space-y-4" @submit.prevent="submit">
          <div>
            <label class="mb-1 block text-xs text-ink-faint">{{ t('admin.username') }}</label>
            <input v-model="form.username" class="input-hand" autocomplete="username" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">{{ t('admin.password') }}</label>
            <input
              v-model="form.password"
              type="password"
              class="input-hand"
              autocomplete="current-password"
              placeholder="••••••"
            />
          </div>

          <p v-if="error" class="text-sm text-terracotta-deep">{{ error }}</p>

          <button type="submit" class="btn-primary w-full" :disabled="loading">
            {{ loading ? t('admin.submitting') : t('admin.submit') }}
          </button>
        </form>
      </div>

      <div class="mt-5 text-center text-sm">
        <router-link to="/" class="text-ink-faint underline decoration-line decoration-dashed underline-offset-4 hover:text-terracotta">
          ← {{ t('common.backHome') }}
        </router-link>
      </div>
    </div>
  </div>
</template>
