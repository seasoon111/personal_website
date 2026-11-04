import { defineStore } from 'pinia'
import { ref } from 'vue'
import api from '@/api'
import type { User } from '@/types'

const TOKEN_KEY = 'pw_token'

export const useAuthStore = defineStore('auth', () => {
  const token = ref(localStorage.getItem(TOKEN_KEY) ?? '')
  const user = ref<User | null>(null)

  async function login(username: string, password: string) {
    const data = await api.login({ username, password })
    token.value = data.token
    user.value = data.user
    localStorage.setItem(TOKEN_KEY, data.token)
  }

  async function fetchMe() {
    if (!token.value) return
    try {
      user.value = await api.me()
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem(TOKEN_KEY)
  }

  return { token, user, login, fetchMe, logout }
})
