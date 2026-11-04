// 日期格式化：中文用「2026 年 11 月 4 日」，英文用「Nov 4, 2026」。
export function formatDate(iso: string | null | undefined, locale: string, withDay = true): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''

  if (locale === 'en-US') {
    return d.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: withDay ? 'numeric' : undefined,
    })
  }
  const month = d.getMonth() + 1
  return withDay ? `${d.getFullYear()} 年 ${month} 月 ${d.getDate()} 日` : `${d.getFullYear()} 年 ${month} 月`
}

// 动态时间线用的短日期：11.04 / Nov 4
export function formatShort(iso: string | null | undefined, locale: string): string {
  if (!iso) return ''
  const d = new Date(iso)
  if (Number.isNaN(d.getTime())) return ''
  if (locale === 'en-US') {
    return d.toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
  }
  return `${String(d.getMonth() + 1).padStart(2, '0')}.${String(d.getDate()).padStart(2, '0')}`
}

export function yearOf(iso: string | null | undefined): string {
  if (!iso) return ''
  const d = new Date(iso)
  return Number.isNaN(d.getTime()) ? '' : String(d.getFullYear())
}
