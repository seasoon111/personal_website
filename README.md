# 个人网站

一个自己写、自己维护的个人站点：文章、项目、动态、关于我，外加一个够用的后台。

- 前台：米色纸感 + 陶土红的手账风，中文排版，支持中英文界面切换
- 后台：登录后管理文章、分类、标签、项目、动态、图片、留言和「关于」页
- 数据：MySQL，图片存本地目录（存储层已抽象，之后可换 OSS）

## 技术栈

| 层 | 选型 |
| --- | --- |
| 后端 | Go 1.26 · Gin · GORM · MySQL 8 · JWT · bcrypt · goldmark |
| 前端 | Vue 3 · TypeScript · Vite 6 · Tailwind 3 · vue-router · Pinia · vue-i18n · axios |
| 样式 | 自定义纸感主题 + `@tailwindcss/typography`，标题字体 LXGW WenKai（CDN） |

## 目录结构

```
personal_website/
├─ backend/                  Go 后端
│  ├─ cmd/server/            入口
│  ├─ internal/
│  │  ├─ config/             配置加载（.env）
│  │  ├─ database/           连接、迁移、初始化管理员
│  │  ├─ dto/                请求/响应结构
│  │  ├─ handler/            HTTP 层
│  │  ├─ middleware/         CORS、JWT、限流
│  │  ├─ model/              GORM 模型
│  │  ├─ pkg/                响应、密码、JWT、slug、markdown
│  │  ├─ router/             路由表
│  │  ├─ service/            业务逻辑（含数据访问）
│  │  └─ storage/            文件存储接口 + 本地实现
│  ├─ uploads/               上传的图片（按年/月分目录）
│  ├─ .env                   本地配置（不提交）
│  └─ .env.example           配置模板
├─ frontend/                 Vue 3 前端
│  └─ src/
│     ├─ api/                接口封装
│     ├─ components/         通用组件
│     ├─ composables/        toast 等
│     ├─ i18n/               中英文文案
│     ├─ layouts/            前台 / 后台布局
│     ├─ pages/              页面（含 pages/admin）
│     ├─ router/             路由与登录守卫
│     ├─ stores/             site / auth
│     ├─ styles/             主题与组件样式
│     ├─ types/              类型定义
│     └─ utils/              日期、markdown
├─ docs/                     需求方案、接口、部署文档
├─ scripts/                  启动、备份、冒烟测试、示例数据
└─ backups/                  备份输出（不提交）
```

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 18+（开发时用的是 22）
- MySQL 8（本地 3306）
- Windows PowerShell 5.1 或 PowerShell 7

### 1. 配置

```powershell
copy backend\.env.example backend\.env
```

编辑 `backend\.env`，至少填对 `DB_PASSWORD`，并把 `JWT_SECRET` 换成一串随机字符（至少 16 位）。

### 2. 一键启动

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts\dev.ps1
```

脚本会：检查环境 → 按 `.env` 建库 → 编译并启动后端（8080）→ 安装前端依赖（首次）→ 启动前端（5173）→ 打开浏览器。按 `Ctrl+C` 会连同后端一起停止。

- 前台：http://127.0.0.1:5173/
- 后台：http://127.0.0.1:5173/admin/login

### 3. 手动启动

```powershell
# 后端
cd backend
go run .\cmd\server        # 监听 http://127.0.0.1:8080

# 前端（另开一个终端）
cd frontend
npm install
npm run dev                # http://127.0.0.1:5173
```

### 4. 默认账号

首次启动且 `users` 表为空时，会用 `.env` 里的 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 建一个管理员：

```
admin / admin123
```

**登录后请立刻在「仪表盘 → 修改密码」里改掉。** 改密码后 `.env` 里的值不会自动同步，它只在首次建号时使用。

### 5. 写入示例内容（可选）

后端和前端都起来之后：

```powershell
powershell -NoProfile -ExecutionPolicy Bypass -File scripts\seed.ps1
```

会创建 3 个分类、5 个标签、3 篇文章、3 个项目、5 条动态，并生成并上传几张封面图，填好「关于」页。想从空白开始就别运行。

## 常用脚本

| 脚本 | 作用 |
| --- | --- |
| `scripts\dev.ps1` | 一键启动（建库、后端、前端） |
| `scripts\seed.ps1` | 写入示例内容 |
| `scripts\smoke-test.ps1` | 后端接口端到端测试（30 项） |
| `scripts\backup.ps1` | 备份数据库和 uploads，默认保留 14 份 |
| `scripts\admin-smoke.mjs` | 用无头 Chrome 检查后台各页面是否正常渲染 |

前端命令：

```powershell
npm run dev        # 开发
npm run build      # 类型检查 + 生产构建（输出 dist/）
npm run type-check # 只做类型检查
npm run preview    # 预览构建结果
```

## 配置项（backend\.env）

| 键 | 说明 | 默认 |
| --- | --- | --- |
| `APP_ENV` | `development` / `production` | development |
| `SERVER_PORT` | 后端端口 | 8080 |
| `DB_HOST` `DB_PORT` `DB_USER` `DB_PASSWORD` `DB_NAME` | MySQL 连接 | 127.0.0.1 / 3306 / root / — / personal_website |
| `JWT_SECRET` | 签名密钥，至少 16 位 | — |
| `JWT_TTL_HOURS` | 登录有效期（小时） | 168 |
| `ADMIN_USERNAME` `ADMIN_PASSWORD` | 首次建管理员用 | admin / admin123 |
| `SITE_NAME` `SITE_DESCRIPTION` `SITE_ICP` | 站点信息 | 个人网站 / … / 空 |
| `UPLOAD_DIR` | 图片存放目录 | ./uploads |
| `UPLOAD_MAX_MB` | 单张图片大小上限 | 10 |
| `CORS_ORIGINS` | 允许的前端来源，逗号分隔 | http://localhost:5173,http://127.0.0.1:5173 |

## 数据库表

`users`、`categories`、`tags`、`articles`、`article_tags`、`projects`、`moments`、`moment_images`、`images`、`about_profiles`、`contact_messages`。

启动时自动迁移（`AutoMigrate`），不需要手写建表 SQL。

## 已实现的功能

**前台**

- 首页：头像卡、最新文章、精选项目、最近动态
- 文章列表：分类筛选、标签筛选、关键词搜索、分页
- 文章详情：Markdown 渲染、代码高亮、自动目录、上一篇/下一篇、标签跳转
- 归档：按年份分组
- 项目列表 / 详情：技术栈、源码与演示链接
- 动态：图文、时间轴样式、图片灯箱
- 关于：头像、简介、技能、经历时间线、社交链接、留言表单
- 中英文界面切换（内容保持中文）、404 页、移动端适配

**后台**

- 登录（5 次/分钟限流）、改密码
- 仪表盘：文章/草稿/项目/动态/图片/未读留言统计，最近草稿
- 文章：列表筛选、新建/编辑（Markdown 编辑器带预览、插入图片、封面选择）、发布/转草稿、置顶、删除
- 分类、标签：增删改
- 项目：增删改、排序、上下架
- 动态：图文发布、编辑、删除、可见性
- 图片库：批量上传、复制链接、删除（被引用时禁止删除）
- 留言：未读筛选、标记已读、删除
- 关于页：头像、简介、技能、经历、社交链接

**其他**

- 统一响应 `{ code, message, data }`，前端统一拦截
- 公开接口无需登录；后台接口 JWT 校验；留言 3 次/小时限流
- 上传图片校验类型与大小，随机文件名，防目录穿越
- 图片被文章封面 / 项目封面 / 动态引用时不允许删除

## 和最初方案的差异

| 方案里写的 | 实际做的 | 原因 |
| --- | --- | --- |
| handler → service → repository 三层 | 合并为 handler → service | 项目规模下多一层只是转发，service 直接持有 `*gorm.DB`，需要时再拆 |
| Swagger 文档 | `docs/api.md` 手写文档 | 少一个依赖，维护成本更低 |
| md-editor-v3 | 自研轻量 Markdown 编辑器 | 体积小、样式统一、无额外依赖；功能覆盖日常写作 |
| 动态的英文标识 | 统一用 `moments` | 与中文「动态」对应，接口/表名一致 |

## 部署

见 [docs/deploy.md](docs/deploy.md)。当前是本地开发环境，尚未部署到服务器。

## 备份与恢复

```powershell
# 备份（数据库 + 图片）
powershell -NoProfile -ExecutionPolicy Bypass -File scripts\backup.ps1

# 恢复数据库
mysql --host=127.0.0.1 --user=root -p personal_website < backups\20260101-120000\personal_website.sql

# 恢复图片：把 uploads.zip 解压回 backend\uploads\
```

## 常见问题

**端口被占用**：后端改 `SERVER_PORT`，前端用 `npm run dev -- --port 5174`，同时把 `frontend\vite.config.ts` 的代理目标改成新端口。

**登录后立刻掉线**：检查浏览器是否禁用了 localStorage，以及 `JWT_SECRET` 是否在重启后变了。

**图片打不开**：确认 `UPLOAD_DIR` 存在且后端有写权限；开发环境由 Vite 把 `/uploads` 代理到 8080。

**中文乱码（脚本报错）**：`.ps1` 文件需要保存为 UTF-8 with BOM，PowerShell 5.1 才能正确识别中文。

**`npm install` 报 `EPERM spawn`**：某些沙箱环境禁止子进程管道，用 `npm install --ignore-scripts` 跳过 postinstall，esbuild 的可执行文件仍然可用。

## 接下来可以做的

- 部署到阿里云 + 域名 + HTTPS
- 文章草稿自动保存、图片粘贴上传
- 全文搜索（MySQL 全文索引或 Meilisearch）
- RSS、站点地图
- 留言回复邮件通知
