export interface ApiResponse<T> {
  code: number
  message: string
  data: T
}

export interface PageResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export interface Category {
  id: number
  name: string
  slug: string
  sort: number
  article_count: number
}

export interface Tag {
  id: number
  name: string
  slug: string
  article_count: number
}

export interface Article {
  id: number
  title: string
  slug: string
  summary: string
  cover_url: string
  category_id: number | null
  status: 'draft' | 'published'
  is_pinned: boolean
  published_at: string | null
  created_at: string
  updated_at: string
  content_md?: string
  content_html?: string
  category?: Category | null
  tags?: Tag[]
}

export interface ArticleDetail {
  article: Article
  prev: Article | null
  next: Article | null
}

export interface ArchiveGroup {
  year: number
  articles: Article[]
}

export interface Project {
  id: number
  title: string
  slug: string
  summary: string
  content_md?: string
  content_html?: string
  cover_url: string
  tech_stack: string[] | null
  repo_url: string
  demo_url: string
  sort: number
  status: 'draft' | 'published'
  created_at: string
  updated_at: string
}

export interface ImageItem {
  id: number
  filename: string
  original_name: string
  path: string
  url: string
  mime: string
  size: number
  width: number
  height: number
  created_at: string
}

export interface MomentImage {
  id: number
  moment_id: number
  image_id: number
  sort: number
  image?: ImageItem
}

export interface Moment {
  id: number
  content: string
  is_public: boolean
  created_at: string
  updated_at: string
  images: MomentImage[]
}

export interface TimelineItem {
  period: string
  title: string
  description: string
}

export interface SocialItem {
  name: string
  url: string
  icon: string
}

export interface AboutProfile {
  id: number
  avatar_url: string
  intro_md: string
  skills: string[] | null
  timeline: TimelineItem[] | null
  socials: SocialItem[] | null
}

export interface SiteInfo {
  name: string
  description: string
  icp: string
  avatar_url: string
  socials: SocialItem[] | null
}

export interface User {
  id: number
  username: string
}

export interface DashboardStats {
  articles: number
  drafts: number
  projects: number
  moments: number
  images: number
  unread_messages: number
  recent_drafts: Article[]
}

export interface ContactMessage {
  id: number
  name: string
  email: string
  subject: string
  content: string
  is_read: boolean
  created_at: string
}
