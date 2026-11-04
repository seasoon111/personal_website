<script setup lang="ts">
import { computed, nextTick, ref } from 'vue'
import { renderMarkdown } from '@/utils/markdown'

const props = withDefaults(
  defineProps<{
    modelValue: string
    rows?: number
    placeholder?: string
  }>(),
  { rows: 18, placeholder: '' },
)

const emit = defineEmits<{ (e: 'update:modelValue', value: string): void }>()

const tab = ref<'write' | 'preview'>('write')
const textarea = ref<HTMLTextAreaElement | null>(null)
const html = computed(() => renderMarkdown(props.modelValue))

function update(value: string) {
  emit('update:modelValue', value)
}

function surround(before: string, after = before) {
  const el = textarea.value
  if (!el) return
  const start = el.selectionStart
  const end = el.selectionEnd
  const selected = props.modelValue.slice(start, end)
  const next = props.modelValue.slice(0, start) + before + selected + after + props.modelValue.slice(end)
  update(next)
  nextTick(() => {
    el.focus()
    el.selectionStart = start + before.length
    el.selectionEnd = start + before.length + selected.length
  })
}

function prefixLine(prefix: string) {
  const el = textarea.value
  if (!el) return
  const start = el.selectionStart
  const lineStart = props.modelValue.lastIndexOf('\n', start - 1) + 1
  const next = props.modelValue.slice(0, lineStart) + prefix + props.modelValue.slice(lineStart)
  update(next)
  nextTick(() => el.focus())
}

// 供父组件在选图后插入 Markdown 图片语法
function insert(text: string) {
  const el = textarea.value
  const pos = el ? el.selectionStart : props.modelValue.length
  const next = props.modelValue.slice(0, pos) + text + props.modelValue.slice(pos)
  update(next)
  nextTick(() => el?.focus())
}

defineExpose({ insert })
</script>

<template>
  <div class="overflow-hidden rounded-xl border border-line bg-white/80">
    <div class="flex flex-wrap items-center gap-1 border-b border-line bg-paper-deep/50 px-2 py-1.5">
      <button type="button" class="editor-btn font-bold" title="加粗" @click="surround('**')">B</button>
      <button type="button" class="editor-btn italic" title="斜体" @click="surround('*')">I</button>
      <button type="button" class="editor-btn" title="二级标题" @click="prefixLine('## ')">H2</button>
      <button type="button" class="editor-btn" title="链接" @click="surround('[', '](https://)')">🔗</button>
      <button type="button" class="editor-btn font-mono" title="行内代码" @click="surround('`')">&lt;/&gt;</button>
      <button type="button" class="editor-btn font-mono" title="代码块" @click="surround('\n```\n', '\n```\n')">{ }</button>
      <button type="button" class="editor-btn" title="引用" @click="prefixLine('> ')">❝</button>
      <button type="button" class="editor-btn" title="分割线" @click="insert('\n\n---\n\n')">—</button>

      <div class="ml-auto flex rounded-full bg-white p-0.5 text-xs">
        <button
          type="button"
          class="rounded-full px-2.5 py-1"
          :class="tab === 'write' ? 'bg-terracotta text-white' : 'text-ink-soft'"
          @click="tab = 'write'"
        >
          编辑
        </button>
        <button
          type="button"
          class="rounded-full px-2.5 py-1"
          :class="tab === 'preview' ? 'bg-terracotta text-white' : 'text-ink-soft'"
          @click="tab = 'preview'"
        >
          预览
        </button>
      </div>
    </div>

    <textarea
      v-show="tab === 'write'"
      ref="textarea"
      :value="modelValue"
      :rows="rows"
      :placeholder="placeholder"
      class="w-full resize-y border-0 bg-transparent px-4 py-3 font-mono text-sm leading-6 text-ink outline-none"
      @input="update(($event.target as HTMLTextAreaElement).value)"
    />

    <div v-if="tab === 'preview'" class="max-h-[70vh] overflow-y-auto px-4 py-3">
      <div v-if="modelValue.trim()" class="article-body" v-html="html" />
      <p v-else class="text-sm text-ink-faint">还没有内容</p>
    </div>
  </div>
</template>

<style scoped>
.editor-btn {
  @apply rounded px-2 py-1 text-xs text-ink-soft transition hover:bg-white hover:text-terracotta-deep;
}
</style>
