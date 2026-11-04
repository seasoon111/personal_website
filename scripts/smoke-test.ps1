# 接口冒烟测试：覆盖登录、内容 CRUD、图片上传、动态、关于我、留言、仪表盘。
# 兼容 Windows PowerShell 5.1（不依赖 pwsh 7 的 -Form）。
# 用法：powershell -NoProfile -File scripts\smoke-test.ps1 [-BaseUrl http://127.0.0.1:8080] [-User admin] [-Pass admin123]
param(
    [string]$BaseUrl = 'http://127.0.0.1:8080',
    [string]$User = 'admin',
    [string]$Pass = 'admin123'
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'
$script:failed = 0
$script:passed = 0

function Invoke-Json {
    param(
        [string]$Method = 'GET',
        [string]$Url,
        [hashtable]$Headers = @{},
        $Body = $null
    )
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
        return $text | ConvertFrom-Json
    }
    finally {
        $client.Dispose()
    }
}

function Step([string]$name, [scriptblock]$action) {
    try {
        $result = & $action
        $script:passed++
        Write-Host "  [PASS] $name" -ForegroundColor Green
        return $result
    }
    catch {
        $script:failed++
        Write-Host "  [FAIL] $name -> $($_.Exception.Message)" -ForegroundColor Red
        if ($_.ErrorDetails.Message) { Write-Host "         $($_.ErrorDetails.Message)" -ForegroundColor DarkGray }
        return $null
    }
}

function ExpectStatus([string]$name, [int]$expected, [scriptblock]$action) {
    try {
        & $action | Out-Null
        if ($expected -eq 200) {
            $script:passed++
            Write-Host "  [PASS] $name" -ForegroundColor Green
        }
        else {
            $script:failed++
            Write-Host "  [FAIL] $name -> 期望 $expected，实际 200" -ForegroundColor Red
        }
    }
    catch {
        $actual = 0
        if ($_.Exception.Response) { $actual = [int]$_.Exception.Response.StatusCode }
        if ($actual -eq $expected) {
            $script:passed++
            Write-Host "  [PASS] $name ($actual)" -ForegroundColor Green
        }
        else {
            $script:failed++
            Write-Host "  [FAIL] $name -> 期望 $expected，实际 $actual" -ForegroundColor Red
        }
    }
}

Write-Host "== 冒烟测试 $BaseUrl ==" -ForegroundColor Cyan

Step '健康检查' { Invoke-Json -Url "$BaseUrl/api/v1/health" } | Out-Null

$login = Step '登录' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/login" -Body @{ username = $User; password = $Pass }
}
if (-not $login) { Write-Host '登录失败，终止测试' -ForegroundColor Red; exit 1 }
$headers = @{ Authorization = "Bearer $($login.data.token)" }

Step '当前用户' { Invoke-Json -Url "$BaseUrl/api/v1/admin/me" -Headers $headers } | Out-Null

$stamp = Get-Date -Format 'MMddHHmmss'
$category = Step '新建分类' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/categories" -Headers $headers `
        -Body @{ name = "测试分类$stamp"; slug = "test-cat-$stamp"; sort = 1 }
}
$tag = Step '新建标签' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/tags" -Headers $headers `
        -Body @{ name = "测试标签$stamp"; slug = "test-tag-$stamp" }
}

$article = Step '新建并发布文章' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/articles" -Headers $headers -Body @{
        title       = "冒烟测试文章 $stamp"
        slug        = "smoke-$stamp"
        summary     = '这是一篇由冒烟测试创建的文章'
        content_md  = "# 标题`n`n正文 **加粗**。`n`n| 列 | 值 |`n| --- | --- |`n| a | b |"
        category_id = $category.data.id
        tag_ids     = @($tag.data.id)
        status      = 'published'
    }
}

Step '前台文章列表' { Invoke-Json -Url "$BaseUrl/api/v1/articles?page=1&page_size=5" } | Out-Null
Step '前台文章详情' { Invoke-Json -Url "$BaseUrl/api/v1/articles/smoke-$stamp" } | Out-Null
Step '归档' { Invoke-Json -Url "$BaseUrl/api/v1/archives" } | Out-Null
Step '分类列表' { Invoke-Json -Url "$BaseUrl/api/v1/categories" } | Out-Null
Step '标签列表' { Invoke-Json -Url "$BaseUrl/api/v1/tags" } | Out-Null
Step '后台文章列表' { Invoke-Json -Url "$BaseUrl/api/v1/admin/articles?page=1" -Headers $headers } | Out-Null
Step '置顶文章' {
    Invoke-Json -Method Patch -Url "$BaseUrl/api/v1/admin/articles/$($article.data.id)/pin" -Headers $headers -Body @{ is_pinned = $true }
} | Out-Null

Step '新建项目' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/projects" -Headers $headers -Body @{
        title = "冒烟测试项目 $stamp"; slug = "smoke-project-$stamp"; summary = '测试项目'
        content_md = '项目详情'; tech_stack = @('Go', 'Vue3', 'MySQL'); status = 'published'; sort = 10
    }
} | Out-Null
Step '前台项目列表' { Invoke-Json -Url "$BaseUrl/api/v1/projects" } | Out-Null
Step '前台项目详情' { Invoke-Json -Url "$BaseUrl/api/v1/projects/smoke-project-$stamp" } | Out-Null

$pngPath = Join-Path $env:TEMP "smoke-$stamp.png"
[IO.File]::WriteAllBytes($pngPath, [Convert]::FromBase64String('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=='))
$upload = Step '上传图片' { Invoke-Upload -Url "$BaseUrl/api/v1/admin/images" -Headers $headers -FilePath $pngPath }
Remove-Item $pngPath -ErrorAction SilentlyContinue

$imageId = $null
if ($upload) {
    $items = @($upload.data)
    if ($items.Count -gt 0) { $imageId = $items[0].id }
}
Step '图片列表' { Invoke-Json -Url "$BaseUrl/api/v1/admin/images?page=1" -Headers $headers } | Out-Null

Step '新建动态' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/admin/moments" -Headers $headers `
        -Body @{ content = "冒烟测试动态 $stamp"; image_ids = @($imageId); is_public = $true }
} | Out-Null
Step '前台动态列表' { Invoke-Json -Url "$BaseUrl/api/v1/moments?page=1" } | Out-Null

if ($imageId) {
    ExpectStatus '删除被引用图片应 409' 409 {
        Invoke-Json -Method Delete -Url "$BaseUrl/api/v1/admin/images/$imageId" -Headers $headers
    }
}

Step '更新关于我' {
    Invoke-Json -Method Put -Url "$BaseUrl/api/v1/admin/about" -Headers $headers -Body @{
        avatar_url = ''
        intro_md   = '你好，这是冒烟测试写入的自我介绍。'
        skills     = @('Go', 'Vue3')
        timeline   = @(@{ period = '2024'; title = '开始写代码'; description = 'hello' })
        socials    = @(@{ name = 'GitHub'; url = 'https://github.com/'; icon = 'github' })
    }
} | Out-Null
Step '前台关于我' { Invoke-Json -Url "$BaseUrl/api/v1/about" } | Out-Null
Step '站点信息' { Invoke-Json -Url "$BaseUrl/api/v1/site" } | Out-Null

Step '提交联系表单' {
    Invoke-Json -Method Post -Url "$BaseUrl/api/v1/contact" `
        -Body @{ name = '测试访客'; email = 'test@example.com'; subject = '冒烟测试'; content = '这是一条冒烟测试留言。' }
} | Out-Null

$messages = Step '后台留言列表' { Invoke-Json -Url "$BaseUrl/api/v1/admin/messages?page=1" -Headers $headers }
if ($messages) {
    $msgItems = @($messages.data.items)
    if ($msgItems.Count -gt 0) {
        $msgId = $msgItems[0].id
        Step '标记留言已读' {
            Invoke-Json -Method Patch -Url "$BaseUrl/api/v1/admin/messages/$msgId/read" -Headers $headers -Body @{ is_read = $true }
        } | Out-Null
    }
}

Step '仪表盘统计' { Invoke-Json -Url "$BaseUrl/api/v1/admin/dashboard" -Headers $headers } | Out-Null

ExpectStatus '未登录访问后台应 401' 401 {
    Invoke-Json -Url "$BaseUrl/api/v1/admin/articles"
}

if ($article) {
    Step '删除文章' {
        Invoke-Json -Method Delete -Url "$BaseUrl/api/v1/admin/articles/$($article.data.id)" -Headers $headers
    } | Out-Null
}

Write-Host ''
Write-Host "结果：通过 $script:passed 项，失败 $script:failed 项" -ForegroundColor $(if ($script:failed -eq 0) { 'Green' } else { 'Red' })
exit $(if ($script:failed -eq 0) { 0 } else { 1 })
