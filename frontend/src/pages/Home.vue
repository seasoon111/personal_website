<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import api from '@/api'
import ArticleCard from '@/components/ArticleCard.vue'
import MomentCard from '@/components/MomentCard.vue'
import ProjectCard from '@/components/ProjectCard.vue'
import SectionTitle from '@/components/SectionTitle.vue'
import EmptyState from '@/components/EmptyState.vue'
import { useSiteStore } from '@/stores/site'
import type { AboutProfile, Article, Moment, Project } from '@/types'
import { renderMarkdown } from '@/utils/markdown'

const { t } = useI18n()
const site = useSiteStore()

const about = ref<AboutProfile | null>(null)
const articles = ref<Article[]>([])
const projects = ref<Project[]>([])
const moments = ref<Moment[]>([])
const loading = ref(true)

const introHtml = computed(() => renderMarkdown(about.value?.intro_md))
const displayName = computed(() => site.info.name)
const socials = computed(() => about.value?.socials ?? site.info.socials ?? [])

onMounted(async () => {
  try {
    const [a, ar, pr, mo] = await Promise.all([
      api.about(),
      api.articles({ page: 1, page_size: 5 }),
      api.projects(),
      api.moments({ page: 1, page_size: 4 }),
    ])
    about.value = a
    articles.value = ar.items ?? []
    projects.value = (pr ?? []).slice(0, 3)
    moments.value = mo.items ?? []
  } catch {
    // 首屏接口失败时展示空状态，不阻断页面。
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div>
    <!-- 首屏：自我介绍 -->
    <section class="container-page pt-12 sm:pt-20">
      <div class="grid items-center gap-10 md:grid-cols-[1fr_280px]">
        <div>
          <p class="text-sm tracking-widest text-ink-faint">{{ t('home.hello') }}</p>
          <h1 class="mt-2 font-hand text-4xl leading-tight text-ink sm:text-5xl">
            {{ displayName }}
            <span class="ml-1 inline-block -rotate-6 text-2xl text-terracotta">✳</span>
          </h1>

          <div class="mt-5 max-w-prose text-sm leading-7 text-ink-soft">
            <div v-if="introHtml" class="intro-preview" v-html="introHtml" />
            <p v-else>{{ site.info.description }}</p>
          </div>

          <div class="mt-6 flex flex-wrap items-center gap-3">
            <router-link to="/about" class="btn-primary">
              {{ t('home.aboutLink') }}
              <span aria-hidden="true">→</span>
            </router-link>
            <a v-if="socials.length" :href="socials[0].url" target="_blank" rel="noopener noreferrer" class="btn-hand">
              {{ socials[0].name }}
            </a>
          </div>

          <div v-if="socials.length > 1" class="mt-4 flex flex-wrap gap-3 text-xs text-ink-faint">
            <a
              v-for="social in socials.slice(1)"
              :key="social.name"
              :href="social.url"
              target="_blank"
              rel="noopener noreferrer"
              class="underline decoration-line decoration-dashed underline-offset-4 hover:text-terracotta"
            >
              {{ social.name }}
            </a>
          </div>
        </div>

        <div class="relative mx-auto w-52 sm:w-64">
          <div class="tape relative rotate-1 rounded-2xl border border-line bg-white/70 p-3 shadow-hand-lg">
            <img
              v-if="about?.avatar_url"
              :src="about.avatar_url"
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
          <p class="mt-4 text-center font-hand text-sm text-ink-faint">—— 一个在写东西的人</p>
        </div>
      </div>
    </section>

    <!-- 最近文章 -->
    <section class="container-page mt-20">
      <SectionTitle :title="t('home.latestArticles')" :to="'/articles'" :link-label="t('home.allArticles')" />
      <div v-if="loading" class="space-y-3">
        <div v-for="i in 3" :key="i" class="h-28 animate-pulse rounded-2xl bg-white/50" />
      </div>
      <div v-else-if="articles.length" class="space-y-3">
        <ArticleCard v-for="(article, index) in articles" :key="article.id" v-reveal="index * 60" :article="article" />
      </div>
      <EmptyState v-else :hint="t('common.emptyHint')" />
    </section>

    <!-- 精选项目 -->
    <section class="container-page mt-20">
      <SectionTitle :title="t('home.featuredProjects')" :to="'/projects'" :link-label="t('home.allProjects')" />
      <div v-if="projects.length" class="grid gap-5 sm:grid-cols-2 lg:grid-cols-3">
        <ProjectCard v-for="(project, index) in projects" :key="project.id" v-reveal="index * 80" :project="project" />
      </div>
      <EmptyState v-else :hint="t('common.emptyHint')" />
    </section>

    <!-- 最近动态 -->
    <section class="container-page mt-20">
      <SectionTitle :title="t('home.latestMoments')" :to="'/moments'" :link-label="t('home.allMoments')" />
      <div v-if="moments.length" class="space-y-6">
        <MomentCard v-for="(moment, index) in moments" :key="moment.id" v-reveal="index * 60" :moment="moment" />
      </div>
      <EmptyState v-else :hint="t('common.emptyHint')" />
    </section>
  </div>
</template>

<style scoped>
/* 首页只露前 4 行，完整介绍留在关于页 */
.intro-preview {
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
.intro-preview :deep(p) {
  margin: 0 0 0.5rem;
}
</style>
