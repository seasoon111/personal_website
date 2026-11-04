# 写入示例内容：分类、标签、文章、项目、动态、关于我、封面图。
# 兼容 Windows PowerShell 5.1。
# 用法：powershell -NoProfile -File scripts\seed.ps1
param(
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$User = 'admin',
    [string]$Pass = 'admin123'
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

function Invoke-Json {
    param([string]$Method = 'GET', [string]$Url, [hashtable]$Headers = @{}, $Body = $null)
    $params = @{ Uri = $Url; Method = $Method; Headers = $Headers; ErrorAction = 'Stop' }
    if ($null -ne $Body) {
        $json = if ($Body -is [string]) { $Body } else { $Body | ConvertTo-Json -Depth 10 -Compress }
        $params.Body = [Text.Encoding]::UTF8.GetBytes($json)
        $params.ContentType = 'application/json; charset=utf-8'
    }
    Invoke-RestMethod @params
}

function Invoke-Upload {
    param([string]$Url, [hashtable]$Headers, [string]$FilePath)
    Add-Type -AssemblyName System.Net.Http
    $client = New-Object System.Net.Http.HttpClient
    try {
        foreach ($k in $Headers.Keys) { $client.DefaultRequestHeaders.Add($k, $Headers[$k]) }
        $content = New-Object System.Net.Http.MultipartFormDataContent
        $bytes = [IO.File]::ReadAllBytes($FilePath)
        $fileContent = New-Object System.Net.Http.ByteArrayContent(, $bytes)
        $fileContent.Headers.ContentType = [System.Net.Http.Headers.MediaTypeHeaderValue]::Parse('image/png')
        $content.Add($fileContent, 'files', [IO.Path]::GetFileName($FilePath))
        $resp = $client.PostAsync($Url, $content).Result
        $text = $resp.Content.ReadAsStringAsync().Result
        if (-not $resp.IsSuccessStatusCode) { throw "HTTP $([int]$resp.StatusCode): $text" }
        return ($text | ConvertFrom-Json).data[0]
    }
    finally { $client.Dispose() }
}

function New-Cover {
    param([string]$Path, [string]$From, [string]$To)
    Add-Type -AssemblyName System.Drawing
    $w = 1200; $h = 800
    $bmp = New-Object System.Drawing.Bitmap($w, $h)
    $g = [System.Drawing.Graphics]::FromImage($bmp)
    $g.SmoothingMode = [System.Drawing.Drawing2D.SmoothingMode]::AntiAlias
    $rect = New-Object System.Drawing.Rectangle(0, 0, $w, $h)
    $brush = New-Object System.Drawing.Drawing2D.LinearGradientBrush(
        $rect,
        [System.Drawing.ColorTranslator]::FromHtml($From),
        [System.Drawing.ColorTranslator]::FromHtml($To),
        35)
    $g.FillRectangle($brush, $rect)

    $circles = @(
        @{ X = 880; Y = 120; S = 380; A = 46; C = '#C4633F' },
        @{ X = 120; Y = 520; S = 260; A = 38; C = '#8A9A7B' },
        @{ X = 700; Y = 560; S = 180; A = 30; C = '#A44E2E' }
    )
    foreach ($c in $circles) {
        $color = [System.Drawing.ColorTranslator]::FromHtml($c.C)
        $sb = New-Object System.Drawing.SolidBrush([System.Drawing.Color]::FromArgb($c.A, $color.R, $color.G, $color.B))
        $g.FillEllipse($sb, $c.X, $c.Y, $c.S, $c.S)
        $sb.Dispose()
    }

    $pen = New-Object System.Drawing.Pen([System.Drawing.Color]::FromArgb(90, 59, 51, 44), 2)
    $pen.DashStyle = [System.Drawing.Drawing2D.DashStyle]::Dash
    $g.DrawLine($pen, 80, 700, 1120, 700)
    $pen.Dispose()

    $bmp.Save($Path, [System.Drawing.Imaging.ImageFormat]::Png)
    $g.Dispose(); $bmp.Dispose(); $brush.Dispose()
}

Write-Host '== 登录 ==' -ForegroundColor Cyan
$login = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/login" -Body @{ username = $User; password = $Pass }
$headers = @{ Authorization = "Bearer $($login.data.token)" }

$tmp = Join-Path $env:TEMP "pw-seed-$(Get-Random)"
New-Item -ItemType Directory -Force -Path $tmp | Out-Null

Write-Host '== 生成并上传封面 ==' -ForegroundColor Cyan
New-Cover -Path (Join-Path $tmp 'cover-1.png') -From '#FAF5EC' -To '#F1E7D6'
New-Cover -Path (Join-Path $tmp 'cover-2.png') -From '#F1E7D6' -To '#E8B39C'
New-Cover -Path (Join-Path $tmp 'cover-3.png') -From '#EADFCB' -To '#C3CDB6'
New-Cover -Path (Join-Path $tmp 'avatar.png') -From '#E8B39C' -To '#FAF5EC'
New-Cover -Path (Join-Path $tmp 'moment-1.png') -From '#FAF5EC' -To '#C3CDB6'

$cover1 = Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath (Join-Path $tmp 'cover-1.png')
$cover2 = Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath (Join-Path $tmp 'cover-2.png')
$cover3 = Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath (Join-Path $tmp 'cover-3.png')
$avatar = Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath (Join-Path $tmp 'avatar.png')
$momentImage = Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath (Join-Path $tmp 'moment-1.png')

Write-Host '== 分类与标签 ==' -ForegroundColor Cyan
$tech = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/categories" -Headers $headers -Body @{ name = '技术'; slug = 'tech'; sort = 1 }
$life = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/categories" -Headers $headers -Body @{ name = '生活'; slug = 'life'; sort = 2 }
$notes = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/categories" -Headers $headers -Body @{ name = '随笔'; slug = 'notes'; sort = 3 }

$tagGo = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers -Body @{ name = 'Go'; slug = 'go' }
$tagVue = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers -Body @{ name = 'Vue3'; slug = 'vue3' }
$tagDB = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers -Body @{ name = 'MySQL'; slug = 'mysql' }
$tagThink = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers -Body @{ name = '思考'; slug = 'thinking' }
$tagTool = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers -Body @{ name = '工具'; slug = 'tools' }

Write-Host '== 文章 ==' -ForegroundColor Cyan
$article1 = @'
# 为什么我要做一个自己的网站

写东西的地方很多：社交平台、笔记软件、文档工具。但它们都有一个共同点——**不属于我**。

平台的规则会变，产品会关停，链接会失效。而一个自己的域名和一台自己的服务器，至少能让这些东西留得久一点。

## 我想要的

- 一个能安静写字的地方，不用考虑流量和算法
- 一个能放项目的地方，不用被压缩成几张截图
- 一个能记录日常的地方，像以前的朋友圈

## 技术上的选择

用 Go 写后端，Vue3 写前端，MySQL 存数据。都是很普通的技术，普通到几年后还能找到人维护。

> 工具应该服务于表达，而不是成为表达本身。

接下来会慢慢把这里填满。
'@

$article2 = @'
# 用 Go + Gin 搭一个能长期维护的 API

这个网站的后端很小，但我想让它**结构清楚**，因为要维护很多年。

## 分层

```
handler  ->  service  ->  model
```

- `handler` 只负责解析参数、调用 service、返回统一响应
- `service` 放业务逻辑和数据库查询
- `model` 是 GORM 模型，字段与表结构一一对应

## 统一响应

所有接口都返回同样的结构：

```json
{ "code": 0, "message": "ok", "data": {} }
```

前端只需要写一次拦截器，后面所有请求都省心。

## 错误处理

业务错误用 `errors.Is` 判断，handler 里统一映射成 HTTP 状态码：

```go
switch {
case errors.Is(err, service.ErrNotFound):
    response.NotFound(c, "内容不存在")
case errors.Is(err, service.ErrInUse):
    response.Conflict(c, "资源正在被使用")
}
```

这样新增业务时不用每次都写一遍错误分支。

## 关于依赖

后端只用了 Gin、GORM、JWT、bcrypt 和 goldmark。依赖越少，几年后升级越轻松。
'@

$article3 = @'
# 手账风界面的一些取舍

这个站点的视觉目标是「像一本摊开的笔记本」，而不是「一个后台管理系统」。

## 配色

主色只有三个：

| 角色 | 色值 | 用途 |
| --- | --- | --- |
| 纸 | `#FAF5EC` | 页面底色 |
| 墨 | `#3B332C` | 正文 |
| 陶土红 | `#C4633F` | 强调、链接、按钮 |

颜色少了，反而不容易出错。

## 手写体

标题用「霞鹜文楷」，正文用系统无衬线。手写体只出现在标题、日期和少量装饰上——**大面积手写体会影响阅读**。

## 装饰的边界

胶带、虚线、纸纹、轻微旋转，都只出现在卡片边缘和空白处。正文区域保持干净，行高 1.8 以上。

> 装饰是为了让内容更可爱，不是为了抢内容的风头。

## 响应式

三档：手机、平板、桌面。手机上隐藏侧边装饰，保留主要信息。
'@

$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/articles" -Headers $headers -Body @{
    title = '为什么我要做一个自己的网站'; slug = 'why-own-website'; summary = '平台会变，产品会关停，而一个自己的域名能让写下的东西留得久一点。'
    content_md = $article1; cover_url = $cover1.url; category_id = $notes.data.id
    tag_ids = @($tagThink.data.id); status = 'published'; is_pinned = $true
}
$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/articles" -Headers $headers -Body @{
    title = '用 Go + Gin 搭一个能长期维护的 API'; slug = 'go-gin-api'; summary = '分层、统一响应、错误映射，以及为什么我只留了五个依赖。'
    content_md = $article2; cover_url = $cover2.url; category_id = $tech.data.id
    tag_ids = @($tagGo.data.id, $tagDB.data.id); status = 'published'
}
$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/articles" -Headers $headers -Body @{
    title = '手账风界面的一些取舍'; slug = 'handbook-ui'; summary = '三个颜色、一款手写体，以及装饰应该停在哪里。'
    content_md = $article3; cover_url = $cover3.url; category_id = $tech.data.id
    tag_ids = @($tagVue.data.id); status = 'published'
}

Write-Host '== 项目 ==' -ForegroundColor Cyan
$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/projects" -Headers $headers -Body @{
    title = '个人网站'; slug = 'personal-website'; summary = 'Go + Vue3 + MySQL 的个人站点，包含博客、作品集、动态和后台管理。'
    content_md = "## 这个项目`n`n从零开始搭的一个个人网站，目标是**能自己维护很多年**。`n`n- 后端：Gin + GORM + MySQL`n- 前端：Vue3 + TypeScript + Tailwind`n- 部署：先本地，后续上阿里云"
    cover_url = $cover1.url; tech_stack = @('Go', 'Gin', 'Vue3', 'MySQL'); status = 'published'; sort = 30
}
$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/projects" -Headers $headers -Body @{
    title = '命令行书签工具'; slug = 'cli-bookmark'; summary = '一个把常用链接存在本地、用命令行快速打开的小工具。'
    content_md = "## 为什么做`n`n浏览器书签越存越多，真正要找的时候又找不到。`n`n于是写了一个只用键盘的工具：`bm add`、`bm open`、`bm list`。"
    cover_url = $cover2.url; tech_stack = @('Go', 'SQLite'); status = 'published'; sort = 20
}
$null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/projects" -Headers $headers -Body @{
    title = '图片压缩小工具'; slug = 'image-squeezer'; summary = '拖进去、压完、下载，不传服务器。'
    content_md = "## 特点`n`n纯浏览器端处理，图片不会上传到任何地方。`n`n用 Canvas 重新编码，支持批量。"
    cover_url = $cover3.url; tech_stack = @('Vue3', 'Canvas'); status = 'published'; sort = 10
}

Write-Host '== 动态 ==' -ForegroundColor Cyan
$moments = @(
    @{ content = '把网站的第一版跑起来了，虽然还很空，但看着它一点点长出来挺开心的。'; image_ids = @() },
    @{ content = '今天调了一下午的样式，最后发现只是少了一个 `border`。'; image_ids = @() },
    @{ content = '给文章详情页加了目录，长文终于不用一直往上滚了。'; image_ids = @($momentImage.id) },
    @{ content = '傍晚出去走了走，风已经有点凉了。'; image_ids = @() },
    @{ content = '接下来想写一写这个网站是怎么搭的，从数据库到部署。'; image_ids = @() }
)
foreach ($m in $moments) {
    $null = Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/moments" -Headers $headers -Body @{
        content = $m.content; image_ids = $m.image_ids; is_public = $true
    }
}

Write-Host '== 关于我 ==' -ForegroundColor Cyan
$null = Invoke-Json -Method Put -Url "$BaseUrl/api/v1/admin/about" -Headers $headers -Body @{
    avatar_url = $avatar.url
    intro_md = "你好，我是这个网站的主人。`n`n白天写代码，晚上写点东西。这里放我的文章、项目和日常，也算是一个长期的练习。`n`n如果你也想做一个自己的网站，欢迎从「关于」页给我留言。"
    skills = @('Go', 'Vue3', 'TypeScript', 'MySQL', 'Tailwind CSS')
    timeline = @(
        @{ period = '2024'; title = '开始认真写代码'; description = '从第一行 Hello World 开始，慢慢把兴趣变成了习惯。' },
        @{ period = '2025'; title = '做了第一个完整的项目'; description = '前后端连起来的那一刻，才真正理解「系统」是什么意思。' },
        @{ period = '2026'; title = '搭了这个网站'; description = '想有一个自己的地方，把学过的东西沉淀下来。' }
    )
    socials = @(
        @{ name = 'GitHub'; url = 'https://github.com/'; icon = 'github' },
        @{ name = 'Email'; url = 'mailto:hello@example.com'; icon = 'mail' }
    )
}

Remove-Item -Recurse -Force $tmp -ErrorAction SilentlyContinue
Write-Host '示例数据写入完成。' -ForegroundColor Green
