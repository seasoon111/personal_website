# 部署到服务器（阿里云 ECS）

目前项目只在本地跑过。这份文档是准备上服务器时的操作清单，按顺序做即可。示例用 Ubuntu 22.04 + 域名 `example.com`，实际替换成自己的。

## 0. 需要准备

- 一台 2 核 2G 以上的 ECS，装 Ubuntu 22.04
- 一个域名，解析到服务器 IP（备案按国内要求办）
- 服务器上装好：MySQL 8、Nginx
- 本地装好：Go（交叉编译用）、Node.js

## 1. 服务器初始化

```bash
sudo apt update && sudo apt upgrade -y
sudo apt install -y nginx mysql-server unzip
sudo systemctl enable --now nginx mysql

# 建目录
sudo mkdir -p /var/www/personal-website/{backend,uploads,dist}
sudo chown -R $USER:$USER /var/www/personal-website
```

MySQL 建库和账号（不要用 root 跑应用）：

```bash
sudo mysql
```

```sql
CREATE DATABASE personal_website DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
CREATE USER 'pw'@'localhost' IDENTIFIED BY '换成强密码';
GRANT ALL PRIVILEGES ON personal_website.* TO 'pw'@'localhost';
FLUSH PRIVILEGES;
EXIT;
```

## 2. 本地构建

后端交叉编译成 Linux 可执行文件：

```powershell
cd backend
$env:GOOS='linux'; $env:GOARCH='amd64'
go build -trimpath -ldflags="-s -w" -o server .\cmd\server
```

前端构建：

```powershell
cd frontend
npm run build      # 产物在 frontend\dist
```

## 3. 上传

```powershell
scp backend\server root@服务器IP:/var/www/personal-website/backend/server
scp -r frontend\dist\* root@服务器IP:/var/www/personal-website/dist/
```

服务器上放配置文件 `/var/www/personal-website/backend/.env`（**不要提交到 Git**）：

```ini
APP_ENV=production
SERVER_PORT=8080
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=pw
DB_PASSWORD=换成强密码
DB_NAME=personal_website
JWT_SECRET=换成一串足够长的随机字符
JWT_TTL_HOURS=168
ADMIN_USERNAME=admin
ADMIN_PASSWORD=换成强密码
SITE_NAME=你的站点名
SITE_DESCRIPTION=一句话介绍
SITE_ICP=你的备案号
UPLOAD_DIR=/var/www/personal-website/uploads
UPLOAD_MAX_MB=10
CORS_ORIGINS=https://example.com
```

```bash
chmod +x /var/www/personal-website/backend/server
```

首次启动会自动建表，并用 `ADMIN_USERNAME` / `ADMIN_PASSWORD` 创建管理员。

## 4. systemd 托管后端

`/etc/systemd/system/personal-website.service`：

```ini
[Unit]
Description=Personal Website API
After=network.target mysql.service

[Service]
Type=simple
User=www-data
Group=www-data
WorkingDirectory=/var/www/personal-website/backend
ExecStart=/var/www/personal-website/backend/server
Restart=always
RestartSec=3
EnvironmentFile=/var/www/personal-website/backend/.env

[Install]
WantedBy=multi-user.target
```

```bash
sudo chown -R www-data:www-data /var/www/personal-website
sudo systemctl daemon-reload
sudo systemctl enable --now personal-website
sudo systemctl status personal-website
curl http://127.0.0.1:8080/api/v1/health
```

## 5. Nginx 反向代理 + 静态站点

`/etc/nginx/sites-available/personal-website`：

```nginx
server {
    listen 80;
    server_name example.com www.example.com;

    root /var/www/personal-website/dist;
    index index.html;

    client_max_body_size 20m;

    # 接口
    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 60s;
    }

    # 上传的图片
    location /uploads/ {
        alias /var/www/personal-website/uploads/;
        expires 30d;
        add_header Cache-Control "public";
    }

    # 前端静态资源长缓存
    location /assets/ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # SPA 路由回退
    location / {
        try_files $uri $uri/ /index.html;
    }
}
```

```bash
sudo ln -s /etc/nginx/sites-available/personal-website /etc/nginx/sites-enabled/
sudo rm -f /etc/nginx/sites-enabled/default
sudo nginx -t && sudo systemctl reload nginx
```

## 6. HTTPS

```bash
sudo apt install -y certbot python3-certbot-nginx
sudo certbot --nginx -d example.com -d www.example.com
sudo systemctl status certbot.timer   # 自动续期
```

启用 HTTPS 后，记得把 `.env` 里的 `CORS_ORIGINS` 改成 `https://example.com`，然后重启后端：

```bash
sudo systemctl restart personal-website
```

## 7. 定时备份

```bash
sudo mkdir -p /var/backups/personal-website
sudo crontab -e
```

```cron
# 每天 3:10 备份数据库
10 3 * * * mysqldump --single-transaction --default-character-set=utf8mb4 personal_website | gzip > /var/backups/personal-website/db-$(date +\%F).sql.gz

# 每天 3:20 打包图片
20 3 * * * cd /var/www/personal-website && tar czf /var/backups/personal-website/uploads-$(date +\%F).tar.gz uploads

# 每周日 4:00 删除 30 天前的备份
0 4 * * 0 find /var/backups/personal-website -type f -mtime +30 -delete
```

数据库备份要带密码时，建议写 `~/.my.cnf`（权限 600），不要把密码写进 crontab：

```ini
[mysqldump]
user=pw
password=换成强密码
```

## 8. 更新流程

```powershell
# 本地
cd backend; $env:GOOS='linux'; $env:GOARCH='amd64'; go build -trimpath -o server .\cmd\server
cd ..\frontend; npm run build

scp backend\server root@服务器IP:/var/www/personal-website/backend/server.new
scp -r frontend\dist\* root@服务器IP:/var/www/personal-website/dist/
```

```bash
# 服务器
cd /var/www/personal-website/backend
mv server.new server && chmod +x server
sudo systemctl restart personal-website
```

数据库结构变化由 `AutoMigrate` 在启动时处理。**注意**：AutoMigrate 只加字段/索引，不会删列或改类型；破坏性变更要自己写 SQL 并先备份。

## 9. 上线检查清单

- [ ] `JWT_SECRET` 已换成随机值，且没有提交到 Git
- [ ] `ADMIN_PASSWORD` 是强密码，登录后立即在后台改一次密码
- [ ] `CORS_ORIGINS` 只包含正式域名
- [ ] `UPLOAD_DIR` 存在且 `www-data` 可写
- [ ] `curl https://example.com/api/v1/health` 返回 `code: 0`
- [ ] 上传一张图片，确认 `/uploads/...` 能访问
- [ ] 前台各页面、后台登录和增删改都点一遍
- [ ] 备份 cron 生效，手动跑一次确认文件生成
- [ ] 阿里云安全组只开放 22 / 80 / 443，3306 不对公网开放
