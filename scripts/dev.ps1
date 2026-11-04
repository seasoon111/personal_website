# 一键启动：检查环境 -> 建库 -> 启动后端 -> 启动前端
# 兼容 Windows PowerShell 5.1，用法：powershell -NoProfile -ExecutionPolicy Bypass -File scripts\dev.ps1
param(
    [switch]$NoBrowser,
    [switch]$SkipInstall
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$root = Split-Path -Parent $PSScriptRoot
$backendDir = Join-Path $root 'backend'
$frontendDir = Join-Path $root 'frontend'
$envFile = Join-Path $backendDir '.env'
$exe = Join-Path $backendDir 'bin\server.exe'

function Write-Step($text) { Write-Host "== $text ==" -ForegroundColor Cyan }
function Fail($text) { Write-Host "!! $text" -ForegroundColor Red; exit 1 }

Write-Step '检查环境'
foreach ($cmd in @('go', 'node', 'npm', 'mysql')) {
    if (-not (Get-Command $cmd -ErrorAction SilentlyContinue)) { Fail "找不到 $cmd，请先安装并加入 PATH" }
}
Write-Host ("go     " + (go version))
Write-Host ("node   " + (node -v))

if (-not (Test-Path $envFile)) {
    Copy-Item (Join-Path $backendDir '.env.example') $envFile
    Write-Host '已从 .env.example 生成 backend\.env，请检查数据库密码后重新运行。' -ForegroundColor Yellow
    exit 0
}

# 读取 .env
$cfg = @{}
Get-Content $envFile | ForEach-Object {
    if ($_ -match '^\s*([A-Z_]+)\s*=\s*(.*)$') { $cfg[$matches[1]] = $matches[2].Trim() }
}
$dbHost = if ($cfg['DB_HOST']) { $cfg['DB_HOST'] } else { '127.0.0.1' }
$dbPort = if ($cfg['DB_PORT']) { $cfg['DB_PORT'] } else { '3306' }
$dbUser = if ($cfg['DB_USER']) { $cfg['DB_USER'] } else { 'root' }
$dbPass = $cfg['DB_PASSWORD']
$dbName = if ($cfg['DB_NAME']) { $cfg['DB_NAME'] } else { 'personal_website' }

Write-Step "确认数据库 $dbName"
$env:MYSQL_PWD = $dbPass
$createSql = "CREATE DATABASE IF NOT EXISTS ``$dbName`` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;"
mysql --host=$dbHost --port=$dbPort --user=$dbUser -e $createSql
if ($LASTEXITCODE -ne 0) { Fail '数据库连接失败，请检查 backend\.env 中的账号密码，并确认 MySQL 已启动' }
Write-Host '数据库就绪'

Write-Step '编译后端'
$env:GOCACHE = Join-Path $root '.cache\go-build'
$env:GOMODCACHE = Join-Path $root '.cache\go-mod'
$env:GOTMPDIR = Join-Path $root '.cache\go-tmp'
New-Item -ItemType Directory -Force -Path (Split-Path $exe) | Out-Null
Push-Location $backendDir
try {
    go build -o $exe .\cmd\server
    if ($LASTEXITCODE -ne 0) { Fail '后端编译失败' }
}
finally { Pop-Location }

Write-Step '启动后端 (http://127.0.0.1:8080)'
$server = Start-Process -FilePath $exe -WorkingDirectory $backendDir -PassThru -WindowStyle Hidden
$healthy = $false
for ($i = 0; $i -lt 30; $i++) {
    Start-Sleep -Milliseconds 500
    try {
        $res = Invoke-RestMethod -Uri 'http://127.0.0.1:8080/api/v1/health' -TimeoutSec 2
        if ($res.code -eq 0) { $healthy = $true; break }
    }
    catch { }
    if ($server.HasExited) { Fail '后端进程已退出，请直接运行 backend\bin\server.exe 查看报错' }
}
if (-not $healthy) { Fail '后端启动超时' }
Write-Host ("后端已启动，PID " + $server.Id)

Write-Step '准备前端依赖'
Push-Location $frontendDir
try {
    if (-not (Test-Path (Join-Path $frontendDir 'node_modules')) -and -not $SkipInstall) {
        Write-Host '首次运行，安装前端依赖（可能需要几分钟）…'
        npm install --no-audit --no-fund
        if ($LASTEXITCODE -ne 0) { Fail '前端依赖安装失败' }
    }

    if (-not $NoBrowser) { Start-Process 'http://127.0.0.1:5173/' }

    Write-Host ''
    Write-Host '前台 http://127.0.0.1:5173/' -ForegroundColor Green
    Write-Host '后台 http://127.0.0.1:5173/admin/login' -ForegroundColor Green
    Write-Host '按 Ctrl+C 停止（后端会一起关闭）' -ForegroundColor DarkGray
    Write-Host ''

    try {
        npm run dev
    }
    finally {
        if (-not $server.HasExited) { Stop-Process -Id $server.Id -Force -ErrorAction SilentlyContinue }
        Write-Host '已停止后端。' -ForegroundColor DarkGray
    }
}
finally { Pop-Location }
