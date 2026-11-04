# 备份数据库和上传的图片
# 用法：powershell -NoProfile -ExecutionPolicy Bypass -File scripts\backup.ps1 [-Keep 14]
param(
    [int]$Keep = 14
)

$ErrorActionPreference = 'Stop'
$ProgressPreference = 'SilentlyContinue'

$root = Split-Path -Parent $PSScriptRoot
$envFile = Join-Path $root 'backend\.env'
if (-not (Test-Path $envFile)) { throw "找不到 backend\.env" }

$cfg = @{}
Get-Content $envFile | ForEach-Object {
    if ($_ -match '^\s*([A-Z_]+)\s*=\s*(.*)$') { $cfg[$matches[1]] = $matches[2].Trim() }
}
$dbHost = if ($cfg['DB_HOST']) { $cfg['DB_HOST'] } else { '127.0.0.1' }
$dbPort = if ($cfg['DB_PORT']) { $cfg['DB_PORT'] } else { '3306' }
$dbUser = if ($cfg['DB_USER']) { $cfg['DB_USER'] } else { 'root' }
$dbName = if ($cfg['DB_NAME']) { $cfg['DB_NAME'] } else { 'personal_website' }

$stamp = Get-Date -Format 'yyyyMMdd-HHmmss'
$target = Join-Path $root "backups\$stamp"
New-Item -ItemType Directory -Force -Path $target | Out-Null

Write-Host "== 导出数据库 $dbName ==" -ForegroundColor Cyan
$env:MYSQL_PWD = $cfg['DB_PASSWORD']
$dumpFile = Join-Path $target "$dbName.sql"
& mysqldump --host=$dbHost --port=$dbPort --user=$dbUser --default-character-set=utf8mb4 --single-transaction --routines --result-file=$dumpFile $dbName
if ($LASTEXITCODE -ne 0) { throw 'mysqldump 失败' }

Write-Host '== 打包上传目录 ==' -ForegroundColor Cyan
$uploads = Join-Path $root 'backend\uploads'
if (Test-Path $uploads) {
    Compress-Archive -Path (Join-Path $uploads '*') -DestinationPath (Join-Path $target 'uploads.zip') -Force
}

$size = [Math]::Round(((Get-ChildItem $target -Recurse | Measure-Object Length -Sum).Sum / 1MB), 2)
Write-Host "备份完成：$target（$size MB）" -ForegroundColor Green

Write-Host '== 清理旧备份 ==' -ForegroundColor Cyan
$all = Get-ChildItem (Join-Path $root 'backups') -Directory | Sort-Object Name -Descending
if ($all.Count -gt $Keep) {
    $all | Select-Object -Skip $Keep | ForEach-Object {
        Remove-Item $_.FullName -Recurse -Force
        Write-Host "已删除旧备份 $($_.Name)"
    }
}
