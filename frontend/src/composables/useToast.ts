import { reactive } from 'vue'

export interface ToastItem {
  id: number
  text: string
  type: 'success' | 'error'
}

const state = reactive<{ items: ToastItem[] }>({ items: [] })
let seq = 0

export function useToast() {
  function push(text: string, type: ToastItem['type']) {
    const id = ++seq
    state.items.push({ id, text, type })
    window.setTimeout(() => {
      const index = state.items.findIndex((item) => item.id === id)
      if (index >= 0) state.items.splice(index, 1)
    }, 2600)
  }

  return {
    items: state.items,
    success: (text: string) => push(text, 'success'),
    error: (text: string) => push(text, 'error'),
  }
}
