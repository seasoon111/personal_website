# 接口文档

## 约定

- 基础地址：`http://127.0.0.1:8080/api/v1`（开发环境也可走前端代理 `http://127.0.0.1:5173/api/v1`）
- 请求体与响应体均为 JSON，编码 UTF-8
- 需要登录的接口带请求头：`Authorization: Bearer <token>`
- 分页参数：`page`（默认 1）、`page_size`（默认 10，最大 100）

### 响应结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {}
}
```

`code` 为 `0` 表示成功，其余为业务错误码：

| HTTP | code | 含义 |
| --- | --- | --- |
| 400 | 40000 | 参数错误 |
| 401 | 40100 | 未登录或 token 失效 |
| 403 | 40300 | 无权限 |
| 404 | 40400 | 资源不存在 |
| 409 | 40900 | 冲突（slug 重复、图片被引用等） |
| 429 | 42900 | 请求过于频繁 |
| 500 | 50000 | 服务器错误 |

### 分页结构

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [],
    "total": 42,
    "page": 1,
    "page_size": 10
  }
}
```

---

## 公开接口

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 健康检查 |
| GET | `/site` | 站点信息 |
| GET | `/about` | 关于我页内容 |
| GET | `/articles` | 文章列表（仅已发布） |
| GET | `/articles/:slug` | 文章详情 |
| GET | `/archives` | 归档（按年分组） |
| GET | `/categories` | 分类列表（含文章数） |
| GET | `/tags` | 标签列表（含文章数） |
| GET | `/projects` | 项目列表（仅已发布） |
| GET | `/projects/:slug` | 项目详情 |
| GET | `/moments` | 动态列表（仅公开） |
| POST | `/contact` | 提交留言（3 次/小时） |

### GET /site

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "name": "个人网站",
    "description": "记录文章、项目和日常的地方",
    "icp": "",
    "avatar_url": "/uploads/2026/11/xxx.png",
    "socials": [{ "name": "GitHub", "url": "https://github.com/", "icon": "github" }]
  }
}
```

### GET /about

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "avatar_url": "/uploads/2026/11/xxx.png",
    "intro_md": "你好，我是……",
    "intro_html": "<p>你好，我是……</p>",
    "skills": ["Go", "Vue3"],
    "timeline": [{ "period": "2026", "title": "搭了这个网站", "description": "……" }],
    "socials": [{ "name": "Email", "url": "mailto:hello@example.com", "icon": "mail" }]
  }
}
```

### GET /articles

查询参数：`page`、`page_size`、`category`（分类 slug）、`tag`（标签 slug）、`keyword`（标题/摘要模糊匹配）。

返回 `PageResult`，`items` 为文章数组（列表不返回 `content_md` / `content_html`）：

```json
{
  "id": 1,
  "title": "为什么我要做一个自己的网站",
  "slug": "why-own-website",
  "summary": "……",
  "cover_url": "/uploads/2026/11/xxx.png",
  "category_id": 3,
  "category": { "id": 3, "name": "随笔", "slug": "notes", "sort": 3 },
  "tags": [{ "id": 4, "name": "思考", "slug": "thinking" }],
  "status": "published",
  "is_pinned": true,
  "published_at": "2026-11-04T21:30:00+08:00",
  "created_at": "2026-11-04T21:30:00+08:00",
  "updated_at": "2026-11-04T21:30:00+08:00"
}
```

### GET /articles/:slug

返回：

```json
{
  "article": { "...": "文章字段，含 content_md 与 content_html" },
  "prev": { "title": "上一篇标题", "slug": "prev-slug" },
  "next": { "title": "下一篇标题", "slug": "next-slug" }
}
```

`prev` 为更早的一篇，`next` 为更新的一篇，没有则为 `null`。

### GET /archives

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    {
      "year": 2026,
      "count": 3,
      "articles": [{ "id": 1, "title": "……", "slug": "……", "published_at": "……" }]
    }
  ]
}
```

### GET /projects

查询参数：`page`、`page_size`。按 `sort` 倒序、`id` 倒序。

```json
{
  "id": 1,
  "title": "个人网站",
  "slug": "personal-website",
  "summary": "……",
  "cover_url": "/uploads/2026/11/xxx.png",
  "tech_stack": ["Go", "Vue3"],
  "repo_url": "https://github.com/",
  "demo_url": "",
  "sort": 30,
  "status": "published",
  "created_at": "……",
  "updated_at": "……"
}
```

### GET /projects/:slug

返回项目详情，含 `content_md` / `content_html`。

### GET /moments

查询参数：`page`、`page_size`。

```json
{
  "id": 5,
  "content": "接下来想写一写这个网站是怎么搭的。",
  "is_public": true,
  "images": [
    {
      "id": 1,
      "sort": 0,
      "image": { "id": 9, "url": "/uploads/2026/11/xxx.png", "width": 1200, "height": 800 }
    }
  ],
  "created_at": "……"
}
```

### POST /contact

```json
{ "name": "小明", "email": "hi@example.com", "subject": "你好", "content": "想请教一个问题……" }
```

成功返回 `{ "code": 0, "message": "留言已收到" }`。限流：同一 IP 每小时 3 次。

---

## 后台接口

除 `/admin/login` 外均需 `Authorization: Bearer <token>`。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| POST | `/admin/login` | 登录（5 次/分钟） |
| GET | `/admin/me` | 当前用户 |
| PUT | `/admin/password` | 修改密码 |
| GET | `/admin/dashboard` | 统计数据 |
| GET | `/admin/articles` | 文章列表（含草稿） |
| POST | `/admin/articles` | 新建文章 |
| GET | `/admin/articles/:id` | 文章详情 |
| PUT | `/admin/articles/:id` | 更新文章 |
| DELETE | `/admin/articles/:id` | 删除文章（软删除） |
| PATCH | `/admin/articles/:id/status` | 发布 / 转草稿 |
| PATCH | `/admin/articles/:id/pin` | 置顶 / 取消置顶 |
| POST | `/admin/categories` | 新建分类 |
| PUT | `/admin/categories/:id` | 更新分类 |
| DELETE | `/admin/categories/:id` | 删除分类（文章转为未分类） |
| POST | `/admin/tags` | 新建标签 |
| PUT | `/admin/tags/:id` | 更新标签 |
| DELETE | `/admin/tags/:id` | 删除标签（解除文章关联） |
| GET | `/admin/projects` | 项目列表（含草稿） |
| POST | `/admin/projects` | 新建项目 |
| GET | `/admin/projects/:id` | 项目详情 |
| PUT | `/admin/projects/:id` | 更新项目 |
| DELETE | `/admin/projects/:id` | 删除项目 |
| GET | `/admin/moments` | 动态列表（含不公开） |
| POST | `/admin/moments` | 新建动态 |
| GET | `/admin/moments/:id` | 动态详情 |
| PUT | `/admin/moments/:id` | 更新动态 |
| DELETE | `/admin/moments/:id` | 删除动态 |
| GET | `/admin/about` | 关于我（编辑用） |
| PUT | `/admin/about` | 保存关于我 |
| GET | `/admin/images` | 图片列表 |
| POST | `/admin/images` | 上传图片（multipart，字段名 `files`，可多张） |
| DELETE | `/admin/images/:id` | 删除图片（被引用时 409） |
| GET | `/admin/messages` | 留言列表（`unread=1` 只看未读） |
| PATCH | `/admin/messages/:id/read` | 标记已读 / 未读 |
| DELETE | `/admin/messages/:id` | 删除留言 |

### POST /admin/login

```json
{ "username": "admin", "password": "admin123" }
```

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "token": "eyJhbGciOi...",
    "expires_at": "2026-11-11T21:30:00+08:00",
    "user": { "id": 1, "username": "admin" }
  }
}
```

### POST /admin/articles

```json
{
  "title": "文章标题",
  "slug": "optional-slug",
  "summary": "摘要",
  "content_md": "## 小标题\n\n正文",
  "cover_url": "/uploads/2026/11/xxx.png",
  "category_id": 1,
  "tag_ids": [1, 2],
  "status": "published",
  "is_pinned": false
}
```

- `slug` 留空会按标题自动生成（保留中英文、数字），重复时自动加后缀
- `status` 只接受 `draft` / `published`；首次发布时写入 `published_at`
- `content_html` 由后端用 goldmark 渲染，前端直接展示

### PATCH /admin/articles/:id/status

```json
{ "status": "published" }
```

### PATCH /admin/articles/:id/pin

```json
{ "is_pinned": true }
```

### PUT /admin/password

```json
{ "old_password": "admin123", "new_password": "至少8位的新密码" }
```

### POST /admin/images

`multipart/form-data`，字段名 `files`，可一次传多张。限制：jpg / png / gif / webp，单张 ≤ `UPLOAD_MAX_MB`（默认 10MB）。

```json
{
  "code": 0,
  "message": "ok",
  "data": [
    { "id": 12, "url": "/uploads/2026/11/xxx.png", "original_name": "cover.png", "size": 102400, "width": 1200, "height": 800 }
  ]
}
```

### PUT /admin/about

```json
{
  "avatar_url": "/uploads/2026/11/avatar.png",
  "intro_md": "你好，我是……",
  "skills": ["Go", "Vue3"],
  "timeline": [{ "period": "2026", "title": "搭了这个网站", "description": "……" }],
  "socials": [{ "name": "GitHub", "url": "https://github.com/", "icon": "github" }]
}
```

### GET /admin/dashboard

```json
{
  "articles": 3,
  "drafts": 0,
  "projects": 3,
  "moments": 5,
  "images": 5,
  "unread_messages": 1,
  "recent_drafts": [{ "id": 4, "title": "草稿标题", "updated_at": "……" }]
}
```

---

## 静态资源

上传的图片通过 `GET /uploads/yyyy/mm/<随机名>.<ext>` 访问，由后端直接提供（生产环境建议交给 Nginx）。
