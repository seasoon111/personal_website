import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import type { SiteInfo } from '@/types'

const fallback: SiteInfo = {
  name: '个人网站',
  description: '记录文章、项目和日常的地方',
  icp: '',
  avatar_url: '',
  socials: [],
}

export const useSiteStore = defineStore('site', () => {
  const info = ref<SiteInfo>({ ...fallback })
  const loaded = ref(false)

  async function load() {
    if (loaded.value) return
    try {
      const data = await api.site()
      info.value = {
        ...fallback,
        ...data,
        socials: data.socials ?? [],
      }
      loaded.value = true
    } catch {
      // 后端未启动时保持占位信息，不阻塞页面渲染。
    }
  }

  return { info, loaded, load }
})
