import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from './i18n'
import './styles/main.css'
import 'highlight.js/styles/atom-one-light.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

// v-reveal：滚动进入视口时淡入上移。用法 v-reveal 或 v-reveal="120"（延迟毫秒）。
app.directive('reveal', {
  mounted(el: HTMLElement, binding) {
    el.classList.add('reveal')
    const delay = Number(binding.value ?? 0)
    if (delay > 0) el.style.transitionDelay = `${delay}ms`

    // 不支持 IntersectionObserver 时直接显示，避免内容永远透明。
    if (!('IntersectionObserver' in window)) {
      el.classList.add('reveal-in')
      return
    }

    let observerFired = false
    const observer = new IntersectionObserver(
      (entries) => {
        observerFired = true
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            el.classList.add('reveal-in')
            observer.unobserve(el)
          }
        })
      },
      { threshold: 0.12 },
    )
    observer.observe(el)

    // 兜底：个别环境（无头浏览器、被隐藏的容器）不会触发回调，2.5s 后强制显示。
    window.setTimeout(() => {
      if (!observerFired) {
        el.classList.add('reveal-in')
        observer.disconnect()
      }
    }, 2500)
  },
})

app.mount('#app')
