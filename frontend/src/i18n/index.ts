import { createI18n } from 'vue-i18n'
import zhCN from '@/locales/zh-CN'
import enUS from '@/locales/en-US'

export const SUPPORTED_LOCALES = ['zh-CN', 'en-US'] as const
export type Locale = (typeof SUPPORTED_LOCALES)[number]

const saved = (localStorage.getItem('pw_locale') as Locale | null) ?? 'zh-CN'

const i18n = createI18n({
  legacy: false,
  locale: SUPPORTED_LOCALES.includes(saved) ? saved : 'zh-CN',
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS,
  },
})

export default i18n
