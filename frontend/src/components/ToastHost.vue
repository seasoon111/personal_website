<script setup lang="ts">
import { useToast } from '@/composables/useToast'

const { items } = useToast()
</script>

<template>
  <teleport to="body">
    <div class="pointer-events-none fixed inset-x-0 top-4 z-[60] flex flex-col items-center gap-2">
      <transition-group name="toast">
        <div
          v-for="item in items"
          :key="item.id"
          class="pointer-events-auto rounded-full border px-4 py-2 text-sm shadow-hand"
          :class="
            item.type === 'success'
              ? 'border-sage-soft bg-white text-ink'
              : 'border-terracotta-soft bg-terracotta text-white'
          "
        >
          {{ item.text }}
        </div>
      </transition-group>
    </div>
  </teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}
</style>
