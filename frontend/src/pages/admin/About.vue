<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import ImagePicker from '@/components/ImagePicker.vue'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import { useToast } from '@/composables/useToast'
import type { ImageItem, SocialItem, TimelineItem } from '@/types'

const { t } = useI18n()
const toast = useToast()

const form = reactive({
  avatar_url: '',
  intro_md: '',
  skillsText: '',
  timelineText: '',
  socialsText: '',
})

const loading = ref(true)
const saving = ref(false)
const pickerOpen = ref(false)

function onAvatarPicked(items: ImageItem[]) {
  if (items[0]) form.avatar_url = items[0].url
}

function parseTimeline(text: string): TimelineItem[] {
  return text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [period = '', title = '', description = ''] = line.split('|').map((part) => part.trim())
      return { period, title, description }
    })
}

function parseSocials(text: string): SocialItem[] {
  return text
    .split('\n')
    .map((line) => line.trim())
    .filter(Boolean)
    .map((line) => {
      const [name = '', url = '', icon = ''] = line.split('|').map((part) => part.trim())
      return { name, url, icon }
    })
}

async function save() {
  saving.value = true
  try {
    await api.updateAbout({
      avatar_url: form.avatar_url,
      intro_md: form.intro_md,
      skills: form.skillsText
        .split('\n')
        .map((line) => line.trim())
        .filter(Boolean),
      timeline: parseTimeline(form.timelineText),
      socials: parseSocials(form.socialsText),
    })
    toast.success(t('admin.saved'))
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    saving.value = false
  }
}

onMounted(async () => {
  try {
    const profile = await api.adminAbout()
    form.avatar_url = profile.avatar_url
    form.intro_md = profile.intro_md
    form.skillsText = (profile.skills ?? []).join('\n')
    form.timelineText = (profile.timeline ?? [])
      .map((item) => `${item.period} | ${item.title} | ${item.description}`)
      .join('\n')
    form.socialsText = (profile.socials ?? [])
      .map((item) => `${item.name} | ${item.url} | ${item.icon}`)
      .join('\n')
  } catch (e) {
    toast.error((e as Error).message)
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="mx-auto max-w-4xl">
    <div class="mb-6 flex items-center justify-between">
      <h1 class="font-hand text-2xl">{{ t('admin.about') }}</h1>
      <button type="button" class="btn-primary" :disabled="saving" @click="save">
        {{ saving ? t('common.saving') : t('common.save') }}
      </button>
    </div>

    <div v-if="loading" class="h-64 animate-pulse rounded-2xl bg-white/60" />

    <div v-else class="grid gap-6 lg:grid-cols-[1fr_280px]">
      <div class="space-y-5">
        <div>
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.aboutIntro') }}</label>
          <MarkdownEditor v-model="form.intro_md" :rows="8" />
        </div>

        <div>
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.aboutSkills') }}</label>
          <textarea v-model="form.skillsText" rows="4" class="input-hand resize-y font-mono text-xs" placeholder="Go&#10;Vue3&#10;MySQL" />
        </div>

        <div>
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.aboutTimeline') }}</label>
          <textarea
            v-model="form.timelineText"
            rows="4"
            class="input-hand resize-y font-mono text-xs"
            placeholder="2024 | 开始写代码 | 从第一行 Hello World 开始&#10;2025 | 做了一个网站 | 用 Go + Vue"
          />
        </div>

        <div>
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.aboutSocials') }}</label>
          <textarea
            v-model="form.socialsText"
            rows="3"
            class="input-hand resize-y font-mono text-xs"
            placeholder="GitHub | https://github.com/you | github&#10;Email | mailto:you@example.com | mail"
          />
        </div>
      </div>

      <aside class="space-y-4">
        <div class="rounded-xl border border-line bg-white/70 p-4">
          <label class="mb-2 block text-xs text-ink-faint">{{ t('admin.avatar') }}</label>
          <img
            v-if="form.avatar_url"
            :src="form.avatar_url"
            class="aspect-square w-full rounded-xl border border-line object-cover"
          />
          <div v-else class="flex aspect-square w-full items-center justify-center rounded-xl bg-paper-deep font-hand text-5xl text-terracotta/50">
            我
          </div>
          <button type="button" class="btn-hand mt-2 w-full justify-center text-xs" @click="pickerOpen = true">
            选择头像
          </button>
          <p class="mt-2 text-xs leading-5 text-ink-faint">{{ t('admin.avatarHint') }}</p>
        </div>
      </aside>
    </div>

    <ImagePicker v-model:open="pickerOpen" @confirm="onAvatarPicked" />
  </div>
</template>
