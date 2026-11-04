import MarkdownIt from 'markdown-it'

const md = new MarkdownIt({
  html: true,
  linkify: true,
  breaks: true,
})

// 渲染后台写入的 Markdown（仅站长本人可写，允许内嵌 HTML）。
export function renderMarkdown(source: string | null | undefined): string {
  if (!source) return ''
  return md.render(source)
}
