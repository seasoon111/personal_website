<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import EmptyState from '@/components/EmptyState.vue'
import type { AboutProfile } from '@/types'
import { renderMarkdown } from '@/utils/markdown'

const { t } = useI18n()

const profile = ref<AboutProfile | null>(null)
const loading = ref(true)

const introHtml = computed(() => renderMarkdown(profile.value?.intro_md))
const skills = computed(() => profile.value?.skills ?? [])
const timeline = computed(() => profile.value?.timeline ?? [])
const socials = computed(() => profile.value?.socials ?? [])

const form = reactive({ name: '', email: '', subject: '', content: '' })
const sending = ref(false)
const sent = ref(false)
const formError = ref('')

async function submit() {
  formError.value = ''
  if (!form.name.trim() || !form.email.trim() || form.content.trim().length < 5) {
    formError.value = '请填写称呼、邮箱和至少 5 个字的留言'
    return
  }
  sending.value = true
  try {
    await api.contact({ ...form })
    sent.value = true
    form.name = ''
    form.email = ''
    form.subject = ''
    form.content = ''
  } catch (e) {
    formError.value = (e as Error).message
  } finally {
    sending.value = false
  }
}

onMounted(async () => {
  try {
    profile.value = await api.about()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="container-page">
    <div v-if="loading" class="space-y-4 pt-14">
      <div class="h-10 w-40 animate-pulse rounded bg-white/60" />
      <div class="h-40 animate-pulse rounded-2xl bg-white/50" />
    </div>

    <template v-else>
      <!-- 头部 -->
      <section class="grid gap-10 pt-12 sm:pt-16 md:grid-cols-[220px_1fr]">
        <div class="relative mx-auto w-44 sm:w-52">
          <div class="tape relative -rotate-1 rounded-2xl border border-line bg-white/70 p-3 shadow-hand-lg">
            <img
              v-if="profile?.avatar_url"
              :src="profile.avatar_url"
              alt="avatar"
              class="aspect-square w-full rounded-xl object-cover"
            />
            <div
              v-else
              class="flex aspect-square w-full items-center justify-center rounded-xl bg-paper-deep font-hand text-6xl text-terracotta/60"
            >
              我
            </div>
          </div>
        </div>

        <div>
          <h1 class="font-hand text-3xl text-ink sm:text-4xl">{{ t('about.title') }}</h1>
          <p class="mt-1 text-sm text-ink-faint">{{ t('about.subtitle') }}</p>
          <div class="article-body mt-5 text-sm" v-html="introHtml" />
          <div v-if="socials.length" class="mt-5 flex flex-wrap gap-3 text-sm">
            <a
              v-for="social in socials"
              :key="social.name"
              :href="social.url"
              target="_blank"
              rel="noopener noreferrer"
              class="btn-hand"
            >
              {{ social.name }} ↗
            </a>
          </div>
        </div>
      </section>

      <!-- 技能 -->
      <section v-if="skills.length" class="mt-16">
        <h2 class="dashed-title">{{ t('about.skills') }}</h2>
        <div class="mt-5 flex flex-wrap gap-2">
          <span
            v-for="skill in skills"
            :key="skill"
            class="rounded-full border border-line bg-white/70 px-3.5 py-1.5 text-sm text-ink-soft"
          >
            {{ skill }}
          </span>
        </div>
      </section>

      <!-- 经历 -->
      <section v-if="timeline.length" class="mt-16">
        <h2 class="dashed-title">{{ t('about.timeline') }}</h2>
        <ol class="mt-6 space-y-6 border-l border-dashed border-line pl-6">
          <li v-for="(item, index) in timeline" :key="index" class="relative">
            <span class="absolute -left-[31px] top-1.5 h-2.5 w-2.5 rounded-full bg-terracotta" />
            <p class="font-hand text-lg text-terracotta-deep">{{ item.period }}</p>
            <p class="mt-1 text-sm font-medium text-ink">{{ item.title }}</p>
            <p v-if="item.description" class="mt-1 text-sm leading-6 text-ink-soft">{{ item.description }}</p>
          </li>
        </ol>
      </section>

      <!-- 联系表单 -->
      <section class="mt-16 max-w-xl">
        <h2 class="dashed-title">{{ t('about.contact') }}</h2>
        <p class="mt-3 text-sm text-ink-faint">{{ t('about.contactHint') }}</p>

        <div v-if="sent" class="mt-6 rounded-2xl border border-sage-soft bg-sage/10 p-5 text-sm text-ink">
          {{ t('about.sent') }}
        </div>

        <form v-else class="mt-6 space-y-4" @submit.prevent="submit">
          <div class="grid gap-4 sm:grid-cols-2">
            <div>
              <label class="mb-1 block text-xs text-ink-faint">{{ t('about.name') }}</label>
              <input v-model="form.name" class="input-hand" :placeholder="t('about.namePlaceholder')" maxlength="64" />
            </div>
            <div>
              <label class="mb-1 block text-xs text-ink-faint">{{ t('about.email') }}</label>
              <input
                v-model="form.email"
                type="email"
                class="input-hand"
                :placeholder="t('about.emailPlaceholder')"
                maxlength="128"
              />
            </div>
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">{{ t('about.subject') }}</label>
            <input v-model="form.subject" class="input-hand" :placeholder="t('about.subjectPlaceholder')" maxlength="200" />
          </div>
          <div>
            <label class="mb-1 block text-xs text-ink-faint">{{ t('about.message') }}</label>
            <textarea
              v-model="form.content"
              rows="5"
              class="input-hand resize-y"
              :placeholder="t('about.messagePlaceholder')"
              maxlength="5000"
            />
          </div>
          <p v-if="formError" class="text-sm text-terracotta-deep">{{ formError }}</p>
          <button type="submit" class="btn-primary" :disabled="sending">
            {{ sending ? t('about.sending') : t('about.send') }}
          </button>
        </form>
      </section>
    </template>

    <EmptyState v-if="!loading && !profile" class="mt-16" :hint="t('common.emptyHint')" />
  </div>
</template>
