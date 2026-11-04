// 后台页面冒烟：用 Chrome DevTools Protocol 打开每个后台路由，检查渲染文本与运行时错误。
// 用法（需先启动后端 8080 与前端 5173）：node scripts/admin-smoke.mjs
import { spawn } from 'node:child_process'
import { mkdirSync } from 'node:fs'
import { join } from 'node:path'

const CHROME = 'C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe'
const PORT = 9333
const APP = 'http://127.0.0.1:5173'
const API = 'http://127.0.0.1:8080/api/v1'
const profile = join(process.cwd(), '.cache', 'cdp-profile')
mkdirSync(profile, { recursive: true })

const sleep = (ms) => new Promise((r) => setTimeout(r, ms))

async function login() {
  const res = await fetch(`${API}/admin/login`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ username: 'admin', password: 'admin123' }),
  })
  const json = await res.json()
  if (!json?.data?.token) throw new Error('登录失败: ' + JSON.stringify(json))
  return json.data.token
}

const chrome = spawn(
  CHROME,
  [
    '--headless=new',
    '--no-sandbox',
    '--disable-gpu',
    '--disable-crash-reporter',
    '--disable-breakpad',
    '--no-first-run',
    '--window-size=1440,1000',
    `--remote-debugging-port=${PORT}`,
    `--user-data-dir=${profile}`,
    'about:blank',
  ],
  { stdio: 'ignore' },
)

let version = null
for (let i = 0; i < 40; i++) {
  try {
    const res = await fetch(`http://127.0.0.1:${PORT}/json/version`)
    version = await res.json()
    break
  } catch {
    await sleep(250)
  }
}
if (!version) {
  chrome.kill()
  throw new Error('Chrome 调试端口未就绪')
}

const ws = new WebSocket(version.webSocketDebuggerUrl)
await new Promise((resolve, reject) => {
  ws.addEventListener('open', resolve, { once: true })
  ws.addEventListener('error', reject, { once: true })
})

let seq = 0
const pending = new Map()
const listeners = new Map()
let runtimeErrors = []
let consoleErrors = []

ws.addEventListener('message', (event) => {
  const msg = JSON.parse(event.data)
  if (msg.id && pending.has(msg.id)) {
    const { resolve, reject } = pending.get(msg.id)
    pending.delete(msg.id)
    if (msg.error) reject(new Error(JSON.stringify(msg.error)))
    else resolve(msg.result)
    return
  }
  if (msg.method === 'Runtime.exceptionThrown') {
    const d = msg.params.exceptionDetails
    runtimeErrors.push(d.exception?.description ?? d.text)
  }
  if (msg.method === 'Runtime.consoleAPICalled' && msg.params.type === 'error') {
    consoleErrors.push(msg.params.args.map((a) => a.value ?? a.description ?? '').join(' '))
  }
  const key = `${msg.sessionId ?? ''}:${msg.method}`
  const waiting = listeners.get(key)
  if (waiting) {
    listeners.delete(key)
    waiting.forEach((fn) => fn(msg.params))
  }
})

function send(method, params = {}, sessionId) {
  return new Promise((resolve, reject) => {
    const id = ++seq
    pending.set(id, { resolve, reject })
    ws.send(JSON.stringify({ id, method, params, ...(sessionId ? { sessionId } : {}) }))
  })
}

function waitEvent(method, sessionId, timeout = 15000) {
  return new Promise((resolve, reject) => {
    const key = `${sessionId}:${method}`
    const timer = setTimeout(() => reject(new Error(`等待 ${method} 超时`)), timeout)
    const fn = (params) => {
      clearTimeout(timer)
      resolve(params)
    }
    if (!listeners.has(key)) listeners.set(key, [])
    listeners.get(key).push(fn)
  })
}

const { targetId } = await send('Target.createTarget', { url: 'about:blank' })
const { sessionId } = await send('Target.attachToTarget', { targetId, flatten: true })
await send('Runtime.enable', {}, sessionId)
await send('Page.enable', {}, sessionId)

const token = await login()

async function visit(path) {
  runtimeErrors = []
  consoleErrors = []
  const loaded = waitEvent('Page.loadEventFired', sessionId)
  await send('Page.navigate', { url: APP + path }, sessionId)
  await loaded
  await sleep(1600)
  const { result } = await send(
    'Runtime.evaluate',
    {
      expression:
        'JSON.stringify({ text: document.body.innerText, url: location.pathname, h1: document.querySelector("h1")?.innerText ?? "", inputs: Array.from(document.querySelectorAll("input")).map(i => i.value) })',
      returnByValue: true,
    },
    sessionId,
  )
  const data = JSON.parse(result.value)
  return { ...data, runtimeErrors: [...runtimeErrors], consoleErrors: [...consoleErrors] }
}

// 先到同源页面写入 token
await visit('/admin/login')
await send('Runtime.evaluate', { expression: `localStorage.setItem('pw_token', ${JSON.stringify(token)})` }, sessionId)

const firstArticle = await (
  await fetch(`${API}/admin/articles?page=1&page_size=1`, { headers: { Authorization: `Bearer ${token}` } })
).json()
const firstArticleId = firstArticle?.data?.items?.[0]?.id

const routes = [
  { path: '/admin', expect: '仪表盘' },
  { path: '/admin/articles', expect: '为什么我要做一个自己的网站' },
  { path: '/admin/articles/new', expect: '草稿' },
  { path: `/admin/articles/${firstArticleId}/edit`, expect: '发布', input: '为什么我要做一个自己的网站' },
  { path: '/admin/categories', expect: '技术' },
  { path: '/admin/tags', expect: 'Go' },
  { path: '/admin/projects', expect: '个人网站' },
  { path: '/admin/projects/new', expect: '保存' },
  { path: '/admin/moments', expect: '动态' },
  { path: '/admin/images', expect: '图片' },
  { path: '/admin/messages', expect: '留言' },
  { path: '/admin/about', expect: '关于' },
]

let failed = 0
for (const route of routes) {
  const page = await visit(route.path)
  const okPath = page.url === route.path
  const okText = page.text.includes(route.expect)
  const okInput = !route.input || page.inputs.some((value) => value.includes(route.input))
  const problems = []
  if (!okPath) problems.push(`跳转到 ${page.url}`)
  if (!okText) problems.push(`缺少「${route.expect}」`)
  if (!okInput) problems.push(`输入框缺少「${route.input}」`)
  if (page.runtimeErrors.length) problems.push('异常: ' + page.runtimeErrors.join(' | '))
  if (page.consoleErrors.length) problems.push('console.error: ' + page.consoleErrors.join(' | '))
  if (problems.length) failed++
  console.log(
    `${problems.length ? 'FAIL' : 'OK  '} ${route.path.padEnd(22)} len=${String(page.text.length).padEnd(6)} ${problems.join('; ')}`,
  )
}

ws.close()
chrome.kill()
console.log(failed === 0 ? '\n后台页面全部通过。' : `\n有 ${failed} 个页面存在问题。`)
process.exit(failed === 0 ? 0 : 1)
