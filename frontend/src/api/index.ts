import axios, { type AxiosRequestConfig } from 'axios'
import type {
  AboutProfile,
  ApiResponse,
  ArchiveGroup,
  Article,
  ArticleDetail,
  Category,
  ContactMessage,
  DashboardStats,
  ImageItem,
  Moment,
  PageResult,
  Project,
  SiteInfo,
  Tag,
  User,
} from '@/types'

const TOKEN_KEY = 'pw_token'

const http = axios.create({
  baseURL: '/api/v1',
  timeout: 20000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY)
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (res) => res.data,
  (error) => {
    const status = error.response?.status
    if (status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      if (location.pathname.startsWith('/admin') && location.pathname !== '/admin/login') {
        location.href = `/admin/login?redirect=${encodeURIComponent(location.pathname)}`
      }
    }
    const message = error.response?.data?.message || '网络异常，请稍后再试'
    return Promise.reject(new Error(message))
  },
)

// 后端统一返回 { code, message, data }，这里只把 data 交给调用方。
async function req<T>(config: AxiosRequestConfig): Promise<T> {
  const body = (await http.request(config)) as unknown as ApiResponse<T>
  return body.data
}

export const api = {
  // —— 公开 ——
  site: () => req<SiteInfo>({ url: '/site' }),
  about: () => req<AboutProfile>({ url: '/about' }),
  articles: (params: Record<string, unknown>) => req<PageResult<Article>>({ url: '/articles', params }),
  article: (slug: string) => req<ArticleDetail>({ url: `/articles/${encodeURIComponent(slug)}` }),
  archives: () => req<ArchiveGroup[]>({ url: '/archives' }),
  categories: () => req<Category[]>({ url: '/categories' }),
  tags: () => req<Tag[]>({ url: '/tags' }),
  projects: () => req<Project[]>({ url: '/projects' }),
  project: (slug: string) => req<Project>({ url: `/projects/${encodeURIComponent(slug)}` }),
  moments: (params: Record<string, unknown>) => req<PageResult<Moment>>({ url: '/moments', params }),
  contact: (body: { name: string; email: string; subject?: string; content: string }) =>
    req<void>({ url: '/contact', method: 'post', data: body }),

  // —— 认证 ——
  login: (body: { username: string; password: string }) =>
    req<{ token: string; expires_in: number; user: User }>({ url: '/admin/login', method: 'post', data: body }),
  me: () => req<User>({ url: '/admin/me' }),
  changePassword: (body: { old_password: string; new_password: string }) =>
    req<void>({ url: '/admin/password', method: 'put', data: body }),
  dashboard: () => req<DashboardStats>({ url: '/admin/dashboard' }),

  // —— 文章 ——
  adminArticles: (params: Record<string, unknown>) =>
    req<PageResult<Article>>({ url: '/admin/articles', params }),
  adminArticle: (id: number) => req<Article>({ url: `/admin/articles/${id}` }),
  createArticle: (body: Record<string, unknown>) => req<Article>({ url: '/admin/articles', method: 'post', data: body }),
  updateArticle: (id: number, body: Record<string, unknown>) =>
    req<Article>({ url: `/admin/articles/${id}`, method: 'put', data: body }),
  deleteArticle: (id: number) => req<void>({ url: `/admin/articles/${id}`, method: 'delete' }),
  setArticleStatus: (id: number, status: 'draft' | 'published') =>
    req<void>({ url: `/admin/articles/${id}/status`, method: 'patch', data: { status } }),
  setArticlePin: (id: number, is_pinned: boolean) =>
    req<void>({ url: `/admin/articles/${id}/pin`, method: 'patch', data: { is_pinned } }),

  // —— 分类 / 标签 ——
  createCategory: (body: { name: string; slug?: string; sort?: number }) =>
    req<Category>({ url: '/admin/categories', method: 'post', data: body }),
  updateCategory: (id: number, body: { name: string; slug?: string; sort?: number }) =>
    req<Category>({ url: `/admin/categories/${id}`, method: 'put', data: body }),
  deleteCategory: (id: number) => req<void>({ url: `/admin/categories/${id}`, method: 'delete' }),
  createTag: (body: { name: string; slug?: string }) =>
    req<Tag>({ url: '/admin/tags', method: 'post', data: body }),
  updateTag: (id: number, body: { name: string; slug?: string }) =>
    req<Tag>({ url: `/admin/tags/${id}`, method: 'put', data: body }),
  deleteTag: (id: number) => req<void>({ url: `/admin/tags/${id}`, method: 'delete' }),

  // —— 项目 ——
  adminProjects: (params: Record<string, unknown>) =>
    req<PageResult<Project>>({ url: '/admin/projects', params }),
  adminProject: (id: number) => req<Project>({ url: `/admin/projects/${id}` }),
  createProject: (body: Record<string, unknown>) => req<Project>({ url: '/admin/projects', method: 'post', data: body }),
  updateProject: (id: number, body: Record<string, unknown>) =>
    req<Project>({ url: `/admin/projects/${id}`, method: 'put', data: body }),
  deleteProject: (id: number) => req<void>({ url: `/admin/projects/${id}`, method: 'delete' }),

  // —— 动态 ——
  adminMoments: (params: Record<string, unknown>) => req<PageResult<Moment>>({ url: '/admin/moments', params }),
  adminMoment: (id: number) => req<Moment>({ url: `/admin/moments/${id}` }),
  createMoment: (body: { content: string; image_ids: number[]; is_public?: boolean }) =>
    req<Moment>({ url: '/admin/moments', method: 'post', data: body }),
  updateMoment: (id: number, body: { content: string; image_ids: number[]; is_public?: boolean }) =>
    req<Moment>({ url: `/admin/moments/${id}`, method: 'put', data: body }),
  deleteMoment: (id: number) => req<void>({ url: `/admin/moments/${id}`, method: 'delete' }),

  // —— 关于我 ——
  adminAbout: () => req<AboutProfile>({ url: '/admin/about' }),
  updateAbout: (body: Record<string, unknown>) => req<AboutProfile>({ url: '/admin/about', method: 'put', data: body }),

  // —— 图片库 ——
  images: (params: Record<string, unknown>) => req<PageResult<ImageItem>>({ url: '/admin/images', params }),
  uploadImages: (files: File[]) => {
    const form = new FormData()
    files.forEach((file) => form.append('files', file))
    return req<ImageItem[]>({ url: '/admin/images', method: 'post', data: form })
  },
  deleteImage: (id: number) => req<void>({ url: `/admin/images/${id}`, method: 'delete' }),

  // —— 留言 ——
  messages: (params: Record<string, unknown>) =>
    req<PageResult<ContactMessage>>({ url: '/admin/messages', params }),
  setMessageRead: (id: number, is_read: boolean) =>
    req<void>({ url: `/admin/messages/${id}/read`, method: 'patch', data: { is_read } }),
  deleteMessage: (id: number) => req<void>({ url: `/admin/messages/${id}`, method: 'delete' }),
}

export default api
