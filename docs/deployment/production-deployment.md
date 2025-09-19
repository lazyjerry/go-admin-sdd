# Go-Admin 生產環境部署指南

本文件提供 Go-Admin 後台管理系統在生產環境中的完整部署指南，包含系統需求、部署架構、安全配置、效能最佳化等內容。

## 目錄

1. [生產環境架構](#生產環境架構)
2. [系統需求](#系統需求)
3. [部署準備](#部署準備)
4. [應用程式部署](#應用程式部署)
5. [資料庫配置](#資料庫配置)
6. [反向代理配置](#反向代理配置)
7. [SSL/TLS 配置](#ssl-tls-配置)
8. [監控與日誌](#監控與日誌)
9. [備份策略](#備份策略)
10. [安全配置](#安全配置)
11. [效能最佳化](#效能最佳化)
12. [故障排除](#故障排除)

## 生產環境架構

### 推薦架構圖

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Load Balancer │    │     Firewall    │    │      CDN        │
│    (Nginx/HAP)  │────│                 │────│   (CloudFlare)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
          │
          │
┌─────────▼─────────┐
│  Web Server       │
│  (Nginx)          │
│  Port: 80/443     │
└─────────┬─────────┘
          │
          │
┌─────────▼─────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Go-Admin App     │    │     Redis       │    │   File Storage  │
│  Port: 8000       │────│   (Cache)       │    │   (MinIO/OSS)   │
│  Instances: 2+    │    │                 │    │                 │
└─────────┬─────────┘    └─────────────────┘    └─────────────────┘
          │
          │
┌─────────▼─────────┐    ┌─────────────────┐    ┌─────────────────┐
│    Database       │    │   Monitoring    │    │     Logging     │
│  (MySQL Cluster)  │    │  (Prometheus)   │    │   (ELK Stack)   │
│  Master/Slave     │    │                 │    │                 │
└───────────────────┘    └─────────────────┘    └─────────────────┘
```

### 環境分層

- **生產環境 (Production)**: 對外服務的正式環境
- **預生產環境 (Staging)**: 生產環境的完整複製，用於最終測試
- **測試環境 (Testing)**: 功能測試和整合測試環境
- **開發環境 (Development)**: 開發者本地環境

## 系統需求

### 硬體需求

#### 最小配置

- **CPU**: 2 核心 2.0GHz
- **記憶體**: 4GB RAM
- **硬碟**: 50GB SSD
- **網路**: 100Mbps

#### 推薦配置

- **CPU**: 4 核心 2.4GHz+
- **記憶體**: 8GB+ RAM
- **硬碟**: 100GB+ SSD
- **網路**: 1Gbps

#### 高可用配置

- **應用伺服器**: 2+ 台（負載均衡）
- **資料庫伺服器**: 主從複製 + 讀寫分離
- **快取伺服器**: Redis 集群
- **檔案伺服器**: 分散式儲存

### 軟體需求

#### 作業系統

- **Ubuntu**: 20.04 LTS+ (推薦)
- **CentOS**: 8+ / AlmaLinux 8+
- **RHEL**: 8+
- **Debian**: 11+

#### 基礎軟體

- **Go**: 1.21+ (編譯用)
- **Nginx**: 1.18+
- **MySQL**: 8.0+ 或 MariaDB 10.6+
- **Redis**: 6.0+
- **Git**: 2.25+

#### 監控軟體

- **Prometheus**: 2.30+
- **Grafana**: 8.0+
- **Node Exporter**: 1.2+
- **Alertmanager**: 0.23+

## 部署準備

### 1. 系統初始化

```bash
#!/bin/bash
# 系統初始化腳本

# 更新系統
sudo apt update && sudo apt upgrade -y

# 安裝基礎套件
sudo apt install -y curl wget git vim htop tree unzip

# 設定時區
sudo timedatectl set-timezone Asia/Taipei

# 設定 NTP 同步
sudo apt install -y ntp
sudo systemctl enable ntp
sudo systemctl start ntp

# 建立系統使用者
sudo adduser --system --group --home /opt/go-admin go-admin
sudo mkdir -p /opt/go-admin/{app,logs,data,backup,scripts}
sudo chown -R go-admin:go-admin /opt/go-admin
```

### 2. 防火牆配置

```bash
# 安裝 UFW
sudo apt install -y ufw

# 預設策略
sudo ufw default deny incoming
sudo ufw default allow outgoing

# 開放必要端口
sudo ufw allow 22/tcp    # SSH
sudo ufw allow 80/tcp    # HTTP
sudo ufw allow 443/tcp   # HTTPS
sudo ufw allow 8000/tcp  # Go-Admin (內部)
sudo ufw allow 3306/tcp  # MySQL (內部)
sudo ufw allow 6379/tcp  # Redis (內部)

# 啟用防火牆
sudo ufw enable

# 檢查狀態
sudo ufw status verbose
```

### 3. 效能調優

```bash
# 建立系統最佳化配置
sudo tee /etc/sysctl.d/99-go-admin.conf << EOF
# 網路效能最佳化
net.core.rmem_max = 16777216
net.core.wmem_max = 16777216
net.ipv4.tcp_rmem = 4096 65536 16777216
net.ipv4.tcp_wmem = 4096 65536 16777216
net.core.netdev_max_backlog = 5000
net.ipv4.tcp_congestion_control = bbr

# 檔案描述符限制
fs.file-max = 1000000

# 記憶體管理
vm.swappiness = 10
vm.dirty_ratio = 15
vm.dirty_background_ratio = 5
EOF

# 套用設定
sudo sysctl -p /etc/sysctl.d/99-go-admin.conf

# 設定使用者限制
sudo tee /etc/security/limits.d/go-admin.conf << EOF
go-admin soft nofile 65536
go-admin hard nofile 65536
go-admin soft nproc 32768
go-admin hard nproc 32768
EOF
```

## 應用程式部署

### 1. 編譯和打包

```bash
#!/bin/bash
# 建置腳本 build.sh

set -e

echo "🔨 開始編譯 Go-Admin..."

# 設定編譯參數
export CGO_ENABLED=1
export GOOS=linux
export GOARCH=amd64

# 版本資訊
VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u '+%Y-%m-%d %H:%M:%S UTC')
COMMIT_SHA=$(git rev-parse --short HEAD)

# 編譯參數
LDFLAGS="-s -w -X 'main.Version=$VERSION' -X 'main.BuildTime=$BUILD_TIME' -X 'main.CommitSHA=$COMMIT_SHA'"

# 清理之前的建置
rm -rf dist/
mkdir -p dist/

# 編譯主程式
echo "📦 編譯主程式..."
go build -ldflags "$LDFLAGS" -o dist/go-admin main.go

# 複製配置檔案
echo "📋 複製配置檔案..."
cp -r config dist/
cp -r static dist/
cp -r template dist/

# 建立版本檔案
echo "$VERSION" > dist/VERSION
echo "$BUILD_TIME" > dist/BUILD_TIME

# 打包
echo "🗜️  建立發布包..."
cd dist
tar -czf go-admin-${VERSION}-linux-amd64.tar.gz *
cd ..

echo "✅ 編譯完成: dist/go-admin-${VERSION}-linux-amd64.tar.gz"
```

### 2. 部署腳本

```bash
#!/bin/bash
# 部署腳本 deploy.sh

set -e

# 配置參數
APP_NAME="go-admin"
APP_USER="go-admin"
APP_DIR="/opt/go-admin"
SERVICE_NAME="go-admin"
BACKUP_DIR="/opt/go-admin/backup"

# 顏色輸出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# 檢查權限
if [[ $EUID -ne 0 ]]; then
   log_error "此腳本需要 root 權限執行"
   exit 1
fi

# 檢查參數
if [ $# -eq 0 ]; then
    log_error "使用方法: $0 <release-package.tar.gz>"
    exit 1
fi

RELEASE_PACKAGE=$1

if [ ! -f "$RELEASE_PACKAGE" ]; then
    log_error "發布包不存在: $RELEASE_PACKAGE"
    exit 1
fi

log_info "開始部署 $APP_NAME..."

# 停止服務
log_info "停止現有服務..."
if systemctl is-active --quiet $SERVICE_NAME; then
    systemctl stop $SERVICE_NAME
    log_info "服務已停止"
fi

# 備份現有版本
if [ -d "$APP_DIR/app" ]; then
    log_info "備份現有版本..."
    BACKUP_NAME="backup-$(date +%Y%m%d-%H%M%S)"
    mkdir -p $BACKUP_DIR
    tar -czf $BACKUP_DIR/$BACKUP_NAME.tar.gz -C $APP_DIR app
    log_info "備份完成: $BACKUP_DIR/$BACKUP_NAME.tar.gz"
fi

# 解壓新版本
log_info "部署新版本..."
mkdir -p $APP_DIR/app
tar -xzf $RELEASE_PACKAGE -C $APP_DIR/app
chown -R $APP_USER:$APP_USER $APP_DIR

# 設定執行權限
chmod +x $APP_DIR/app/go-admin

# 啟動服務
log_info "啟動服務..."
systemctl start $SERVICE_NAME
systemctl enable $SERVICE_NAME

# 檢查服務狀態
sleep 3
if systemctl is-active --quiet $SERVICE_NAME; then
    log_info "✅ 部署成功！服務正在運行"

    # 顯示版本資訊
    if [ -f "$APP_DIR/app/VERSION" ]; then
        VERSION=$(cat $APP_DIR/app/VERSION)
        log_info "部署版本: $VERSION"
    fi
else
    log_error "❌ 服務啟動失敗"
    journalctl -u $SERVICE_NAME --no-pager -l
    exit 1
fi
```

### 3. Systemd 服務配置

```ini
# /etc/systemd/system/go-admin.service
[Unit]
Description=Go-Admin Backend Service
Documentation=https://github.com/go-admin-team/go-admin
After=network.target mysql.service redis.service

[Service]
Type=simple
User=go-admin
Group=go-admin
WorkingDirectory=/opt/go-admin/app
ExecStart=/opt/go-admin/app/go-admin server -c /opt/go-admin/app/config/settings.yml
ExecReload=/bin/kill -USR2 $MAINPID

# 重啟策略
Restart=always
RestartSec=5
StartLimitInterval=60s
StartLimitBurst=3

# 資源限制
LimitNOFILE=65536
LimitNPROC=32768

# 環境變數
Environment=GIN_MODE=release
Environment=GO_ENV=production

# 日誌
StandardOutput=journal
StandardError=journal
SyslogIdentifier=go-admin

# 安全設定
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ReadWritePaths=/opt/go-admin/logs /opt/go-admin/data

[Install]
WantedBy=multi-user.target
```

## 資料庫配置

### 1. MySQL 主從配置

#### 主伺服器配置 (master)

```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf

[mysqld]
# 基本配置
server-id = 1
log-bin = mysql-bin
binlog-format = ROW
gtid-mode = ON
enforce-gtid-consistency = ON

# 效能最佳化
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
innodb_flush_log_at_trx_commit = 1
sync_binlog = 1

# 連接配置
max_connections = 1000
max_connect_errors = 100000

# 字符集配置
character_set_server = utf8mb4
collation_server = utf8mb4_unicode_ci
```

#### 從伺服器配置 (slave)

```ini
# /etc/mysql/mysql.conf.d/mysqld.cnf

[mysqld]
# 基本配置
server-id = 2
relay-log = mysql-relay-bin
log-slave-updates = ON
read-only = ON

# 復製配置
gtid-mode = ON
enforce-gtid-consistency = ON
slave-skip-errors = 1062,1053,1146

# 效能最佳化（同主伺服器）
innodb_buffer_pool_size = 2G
innodb_log_file_size = 256M
```

### 2. 資料庫初始化腳本

```bash
#!/bin/bash
# 資料庫初始化腳本

DB_HOST="localhost"
DB_PORT="3306"
DB_ROOT_PASSWORD="your_root_password"
DB_NAME="go_admin"
DB_USER="go_admin_user"
DB_PASSWORD="secure_password"

# 建立資料庫和使用者
mysql -h $DB_HOST -P $DB_PORT -u root -p$DB_ROOT_PASSWORD << EOF
CREATE DATABASE IF NOT EXISTS $DB_NAME CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER IF NOT EXISTS '$DB_USER'@'localhost' IDENTIFIED BY '$DB_PASSWORD';
CREATE USER IF NOT EXISTS '$DB_USER'@'%' IDENTIFIED BY '$DB_PASSWORD';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER
ON $DB_NAME.* TO '$DB_USER'@'localhost';

GRANT SELECT, INSERT, UPDATE, DELETE, CREATE, DROP, INDEX, ALTER
ON $DB_NAME.* TO '$DB_USER'@'%';

FLUSH PRIVILEGES;
EOF

echo "資料庫初始化完成"
```

## 反向代理配置

### Nginx 配置

```nginx
# /etc/nginx/sites-available/go-admin
upstream go_admin_backend {
    least_conn;
    server 127.0.0.1:8000 weight=1 max_fails=3 fail_timeout=30s;
    server 127.0.0.1:8001 weight=1 max_fails=3 fail_timeout=30s;
    keepalive 32;
}

# HTTP 重定向到 HTTPS
server {
    listen 80;
    server_name your-domain.com www.your-domain.com;
    return 301 https://$server_name$request_uri;
}

# HTTPS 主配置
server {
    listen 443 ssl http2;
    server_name your-domain.com www.your-domain.com;

    # SSL 配置
    ssl_certificate /etc/ssl/certs/your-domain.crt;
    ssl_certificate_key /etc/ssl/private/your-domain.key;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # 安全標頭
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # 日誌配置
    access_log /var/log/nginx/go-admin.access.log;
    error_log /var/log/nginx/go-admin.error.log;

    # Gzip 壓縮
    gzip on;
    gzip_comp_level 6;
    gzip_types
        text/plain
        text/css
        text/xml
        text/javascript
        application/json
        application/javascript
        application/xml+rss
        image/svg+xml;

    # 靜態檔案配置
    location /static/ {
        alias /opt/go-admin/app/static/;
        expires 30d;
        add_header Cache-Control "public, no-transform";
        access_log off;
    }

    # 檔案上傳配置
    location /upload/ {
        alias /opt/go-admin/data/uploads/;
        expires 7d;
        add_header Cache-Control "public";
        access_log off;
    }

    # API 代理配置
    location /api/ {
        proxy_pass http://go_admin_backend;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_set_header X-Forwarded-Host $server_name;

        # 超時配置
        proxy_connect_timeout 60s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;

        # 緩衝配置
        proxy_buffering on;
        proxy_buffer_size 4k;
        proxy_buffers 8 4k;

        # 限制上傳大小
        client_max_body_size 10M;
    }

    # 前端代理配置
    location / {
        proxy_pass http://go_admin_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # 健康檢查
    location /health {
        proxy_pass http://go_admin_backend/api/v1/health;
        access_log off;
    }
}
```

## SSL/TLS 配置

### 1. Let's Encrypt 自動證書

```bash
#!/bin/bash
# SSL 證書自動配置腳本

DOMAIN="your-domain.com"
EMAIL="admin@your-domain.com"

# 安裝 Certbot
sudo apt update
sudo apt install -y certbot python3-certbot-nginx

# 獲取證書
sudo certbot --nginx -d $DOMAIN -d www.$DOMAIN --email $EMAIL --agree-tos --no-eff-email

# 設定自動續約
sudo crontab -l > /tmp/crontab.tmp
echo "0 12 * * * /usr/bin/certbot renew --quiet && systemctl reload nginx" >> /tmp/crontab.tmp
sudo crontab /tmp/crontab.tmp
rm /tmp/crontab.tmp

echo "SSL 證書配置完成"
```

### 2. 證書管理腳本

```bash
#!/bin/bash
# 證書檢查和更新腳本

check_certificates() {
    echo "檢查 SSL 證書狀態..."

    for domain in your-domain.com www.your-domain.com; do
        expiry_date=$(echo | openssl s_client -servername $domain -connect $domain:443 2>/dev/null | openssl x509 -noout -dates | grep notAfter | cut -d= -f2)
        expiry_epoch=$(date -d "$expiry_date" +%s)
        current_epoch=$(date +%s)
        days_until_expiry=$(( (expiry_epoch - current_epoch) / 86400 ))

        echo "域名: $domain"
        echo "到期日期: $expiry_date"
        echo "剩餘天數: $days_until_expiry"

        if [ $days_until_expiry -lt 30 ]; then
            echo "⚠️  證書即將到期，需要更新"
        else
            echo "✅ 證書有效"
        fi
        echo "---"
    done
}

renew_certificates() {
    echo "更新 SSL 證書..."
    certbot renew --quiet

    if [ $? -eq 0 ]; then
        echo "✅ 證書更新成功"
        systemctl reload nginx
    else
        echo "❌ 證書更新失敗"
        exit 1
    fi
}

case "$1" in
    check)
        check_certificates
        ;;
    renew)
        renew_certificates
        ;;
    *)
        echo "使用方法: $0 {check|renew}"
        exit 1
        ;;
esac
```

## 監控與日誌

### 1. Prometheus 配置

```yaml
# /etc/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "go-admin-rules.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093

scrape_configs:
  # Go-Admin 應用監控
  - job_name: "go-admin"
    static_configs:
      - targets: ["localhost:8000", "localhost:8001"]
    metrics_path: "/metrics"
    scrape_interval: 30s

  # 系統監控
  - job_name: "node"
    static_configs:
      - targets: ["localhost:9100"]

  # MySQL 監控
  - job_name: "mysql"
    static_configs:
      - targets: ["localhost:9104"]

  # Redis 監控
  - job_name: "redis"
    static_configs:
      - targets: ["localhost:9121"]

  # Nginx 監控
  - job_name: "nginx"
    static_configs:
      - targets: ["localhost:9113"]
```

### 2. 告警規則

```yaml
# /etc/prometheus/go-admin-rules.yml
groups:
  - name: go-admin.rules
    rules:
      # 應用程式告警
      - alert: GoAdminDown
        expr: up{job="go-admin"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Go-Admin 服務停止"
          description: "Go-Admin 實例 {{ $labels.instance }} 已停止運行超過 1 分鐘"

      - alert: GoAdminHighResponseTime
        expr: histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Go-Admin 回應時間過長"
          description: "Go-Admin 95% 回應時間超過 1 秒"

      - alert: GoAdminHighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "Go-Admin 錯誤率過高"
          description: "Go-Admin 5xx 錯誤率超過 10%"

      # 系統資源告警
      - alert: HighCPUUsage
        expr: 100 - (avg by(instance) (rate(node_cpu_seconds_total{mode="idle"}[2m])) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "CPU 使用率過高"
          description: "主機 {{ $labels.instance }} CPU 使用率超過 80%"

      - alert: HighMemoryUsage
        expr: (node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100 > 85
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "記憶體使用率過高"
          description: "主機 {{ $labels.instance }} 記憶體使用率超過 85%"

      - alert: DiskSpaceLow
        expr: (node_filesystem_avail_bytes / node_filesystem_size_bytes) * 100 < 10
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "磁碟空間不足"
          description: "主機 {{ $labels.instance }} 磁碟空間剩餘不足 10%"

      # 資料庫告警
      - alert: MySQLDown
        expr: mysql_up == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "MySQL 資料庫停止"
          description: "MySQL 實例 {{ $labels.instance }} 已停止運行"

      - alert: MySQLSlowQueries
        expr: rate(mysql_global_status_slow_queries[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "MySQL 慢查詢過多"
          description: "MySQL 慢查詢數量超過每分鐘 10 個"
```

### 3. 日誌配置

```yaml
# /etc/filebeat/filebeat.yml
filebeat.inputs:
  - type: log
    enabled: true
    paths:
      - /opt/go-admin/logs/*.log
    fields:
      service: go-admin
      environment: production
    multiline.pattern: '^\d{4}-\d{2}-\d{2}'
    multiline.negate: true
    multiline.match: after

  - type: log
    enabled: true
    paths:
      - /var/log/nginx/go-admin.access.log
    fields:
      service: nginx
      log_type: access

  - type: log
    enabled: true
    paths:
      - /var/log/nginx/go-admin.error.log
    fields:
      service: nginx
      log_type: error

output.elasticsearch:
  hosts: ["localhost:9200"]
  index: "go-admin-%{+yyyy.MM.dd}"

logging.level: info
logging.to_files: true
logging.files:
  path: /var/log/filebeat
  name: filebeat
  keepfiles: 7
  permissions: 0644
```

### 4. Grafana 儀表板

```json
{
	"dashboard": {
		"id": null,
		"title": "Go-Admin 監控儀表板",
		"tags": ["go-admin", "production"],
		"timezone": "browser",
		"panels": [
			{
				"title": "應用程式狀態",
				"type": "stat",
				"targets": [
					{
						"expr": "up{job=\"go-admin\"}",
						"legendFormat": "{{instance}}"
					}
				]
			},
			{
				"title": "QPS (每秒請求數)",
				"type": "graph",
				"targets": [
					{
						"expr": "rate(http_requests_total[5m])",
						"legendFormat": "{{instance}}"
					}
				]
			},
			{
				"title": "回應時間",
				"type": "graph",
				"targets": [
					{
						"expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))",
						"legendFormat": "P95"
					},
					{
						"expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))",
						"legendFormat": "P99"
					}
				]
			},
			{
				"title": "錯誤率",
				"type": "graph",
				"targets": [
					{
						"expr": "rate(http_requests_total{status=~\"5..\"}[5m])",
						"legendFormat": "5xx 錯誤"
					}
				]
			},
			{
				"title": "系統資源",
				"type": "graph",
				"targets": [
					{
						"expr": "100 - (avg(rate(node_cpu_seconds_total{mode=\"idle\"}[2m])) * 100)",
						"legendFormat": "CPU 使用率"
					},
					{
						"expr": "(node_memory_MemTotal_bytes - node_memory_MemAvailable_bytes) / node_memory_MemTotal_bytes * 100",
						"legendFormat": "記憶體使用率"
					}
				]
			}
		]
	}
}
```

## 備份策略

### 1. 資料庫備份腳本

```bash
#!/bin/bash
# 資料庫備份腳本

# 配置參數
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="go_admin"
DB_USER="backup_user"
DB_PASSWORD="backup_password"
BACKUP_DIR="/opt/go-admin/backup/mysql"
RETENTION_DAYS=7

# 建立備份目錄
mkdir -p $BACKUP_DIR

# 備份檔案名稱
BACKUP_FILE="go_admin_$(date +%Y%m%d_%H%M%S).sql.gz"

# 執行備份
echo "開始備份資料庫..."
mysqldump -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD \
    --single-transaction \
    --routines \
    --triggers \
    --events \
    --set-gtid-purged=OFF \
    $DB_NAME | gzip > $BACKUP_DIR/$BACKUP_FILE

if [ $? -eq 0 ]; then
    echo "✅ 資料庫備份成功: $BACKUP_FILE"

    # 清理舊備份
    find $BACKUP_DIR -name "go_admin_*.sql.gz" -mtime +$RETENTION_DAYS -delete
    echo "🧹 清理 $RETENTION_DAYS 天前的備份檔案"
else
    echo "❌ 資料庫備份失敗"
    exit 1
fi
```

### 2. 應用程式備份腳本

```bash
#!/bin/bash
# 應用程式備份腳本

APP_DIR="/opt/go-admin"
BACKUP_DIR="/opt/go-admin/backup/app"
RETENTION_DAYS=30

# 建立備份目錄
mkdir -p $BACKUP_DIR

# 備份檔案名稱
BACKUP_FILE="go_admin_app_$(date +%Y%m%d_%H%M%S).tar.gz"

echo "開始備份應用程式..."

# 建立備份
tar -czf $BACKUP_DIR/$BACKUP_FILE \
    --exclude="$APP_DIR/backup" \
    --exclude="$APP_DIR/logs/*.log" \
    --exclude="$APP_DIR/data/tmp" \
    -C $APP_DIR \
    app config data

if [ $? -eq 0 ]; then
    echo "✅ 應用程式備份成功: $BACKUP_FILE"

    # 清理舊備份
    find $BACKUP_DIR -name "go_admin_app_*.tar.gz" -mtime +$RETENTION_DAYS -delete
    echo "🧹 清理 $RETENTION_DAYS 天前的備份檔案"
else
    echo "❌ 應用程式備份失敗"
    exit 1
fi
```

### 3. 自動備份 Crontab

```bash
# 編輯 crontab
sudo crontab -e

# 添加以下內容
# 每日凌晨 2 點備份資料庫
0 2 * * * /opt/go-admin/scripts/backup_database.sh >> /opt/go-admin/logs/backup.log 2>&1

# 每週日凌晨 3 點備份應用程式
0 3 * * 0 /opt/go-admin/scripts/backup_app.sh >> /opt/go-admin/logs/backup.log 2>&1

# 每小時檢查服務狀態
0 * * * * /opt/go-admin/scripts/health_check.sh >> /opt/go-admin/logs/health.log 2>&1
```

### 4. 恢復腳本

```bash
#!/bin/bash
# 資料庫恢復腳本

if [ $# -eq 0 ]; then
    echo "使用方法: $0 <backup_file.sql.gz>"
    exit 1
fi

BACKUP_FILE=$1
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="go_admin"
DB_USER="restore_user"
DB_PASSWORD="restore_password"

if [ ! -f "$BACKUP_FILE" ]; then
    echo "❌ 備份檔案不存在: $BACKUP_FILE"
    exit 1
fi

echo "⚠️  警告：此操作將覆蓋現有資料庫！"
read -p "請輸入 'YES' 確認繼續: " confirmation

if [ "$confirmation" != "YES" ]; then
    echo "取消恢復操作"
    exit 1
fi

echo "開始恢復資料庫..."

# 停止應用程式服務
systemctl stop go-admin

# 恢復資料庫
zcat $BACKUP_FILE | mysql -h $DB_HOST -P $DB_PORT -u $DB_USER -p$DB_PASSWORD $DB_NAME

if [ $? -eq 0 ]; then
    echo "✅ 資料庫恢復成功"

    # 重啟應用程式服務
    systemctl start go-admin

    echo "✅ 服務已重啟"
else
    echo "❌ 資料庫恢復失敗"
    exit 1
fi
```

## 安全配置

### 1. 應用程式安全配置

```yaml
# config/settings.prod.yml 安全配置
settings:
  application:
    mode: release
    # 禁用調試模式
    debug: false

  # JWT 安全配置
  jwt:
    # 使用強密鑰（至少 32 字元）
    secret: "your-very-strong-jwt-secret-key-here-32-chars+"
    # 較短的過期時間
    timeout: 2
    # 啟用刷新權杖
    refresh: true

  # 資料庫安全配置
  database:
    # 使用專用使用者，非 root
    username: "go_admin_user"
    password: "strong-database-password"
    # 啟用 SSL 連接
    config: "charset=utf8mb4&parseTime=True&loc=Local&tls=true"

  # Redis 安全配置
  redis:
    # 設定密碼
    password: "strong-redis-password"

  # 日誌安全配置
  log:
    # 不要記錄敏感資訊
    level: "info"
    # 定期清理日誌
    max-age: 30

  # 檔案上傳安全配置
  upload:
    # 限制檔案類型
    allowed_types: ["jpg", "jpeg", "png", "gif", "pdf", "doc", "docx"]
    # 限制檔案大小
    max_size: 10485760 # 10MB
    # 安全的上傳路徑
    path: "/opt/go-admin/data/uploads"
```

### 2. 系統安全加固

```bash
#!/bin/bash
# 系統安全加固腳本

# 禁用不必要的服務
sudo systemctl disable avahi-daemon
sudo systemctl disable cups
sudo systemctl disable bluetooth

# 配置 SSH 安全
sudo tee /etc/ssh/sshd_config.d/99-security.conf << EOF
# 禁用 root 登入
PermitRootLogin no

# 禁用密碼登入，強制使用金鑰
PasswordAuthentication no
PubkeyAuthentication yes

# 限制登入嘗試
MaxAuthTries 3
MaxStartups 10:30:100

# 其他安全配置
Protocol 2
X11Forwarding no
AllowTcpForwarding no
ClientAliveInterval 300
ClientAliveCountMax 2
EOF

# 重啟 SSH 服務
sudo systemctl restart sshd

# 安裝 fail2ban
sudo apt install -y fail2ban

# 配置 fail2ban
sudo tee /etc/fail2ban/jail.d/go-admin.conf << EOF
[sshd]
enabled = true
port = ssh
filter = sshd
logpath = /var/log/auth.log
maxretry = 3
bantime = 3600

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/go-admin.error.log
maxretry = 5
bantime = 1800
EOF

sudo systemctl enable fail2ban
sudo systemctl start fail2ban

echo "✅ 安全配置完成"
```

### 3. 網路安全配置

```bash
# IPTables 規則配置
sudo iptables -F

# 允許本地回環
sudo iptables -A INPUT -i lo -j ACCEPT

# 允許已建立的連接
sudo iptables -A INPUT -m conntrack --ctstate ESTABLISHED,RELATED -j ACCEPT

# 允許 SSH (限制來源 IP)
sudo iptables -A INPUT -p tcp --dport 22 -s YOUR_OFFICE_IP -j ACCEPT

# 允許 HTTP/HTTPS
sudo iptables -A INPUT -p tcp --dport 80 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# 允許內部通信（應用程式端口）
sudo iptables -A INPUT -p tcp --dport 8000 -s 127.0.0.1 -j ACCEPT
sudo iptables -A INPUT -p tcp --dport 8001 -s 127.0.0.1 -j ACCEPT

# 預設拒絕所有其他連接
sudo iptables -A INPUT -j DROP

# 保存規則
sudo iptables-save > /etc/iptables/rules.v4

# 開機自動載入
echo 'iptables-restore < /etc/iptables/rules.v4' >> /etc/rc.local
```

## 效能最佳化

### 1. 應用程式效能調優

```yaml
# config/settings.prod.yml 效能配置
settings:
  application:
    # 生產模式
    mode: release
    # 適當的讀寫超時
    readtimeout: 60
    writetimeout: 60

  database:
    # 連接池配置
    max-idle-conns: 50
    max-open-conns: 200
    conn-max-lifetime: 1800

    # 查詢超時
    query-timeout: 30

  redis:
    # 連接池配置
    pool-size: 100
    min-idle-conns: 10
    idle-timeout: 300

  # 快取配置
  cache:
    # 啟用多層快取
    enable: true
    ttl: 3600

    # Redis 快取
    redis:
      enable: true
      prefix: "go-admin:"

    # 記憶體快取
    memory:
      enable: true
      size: 1000
```

### 2. 資料庫效能調優

```sql
-- MySQL 效能調優 SQL

-- 建立必要的索引
CREATE INDEX idx_user_username ON sys_users(username);
CREATE INDEX idx_user_email ON sys_users(email);
CREATE INDEX idx_user_status ON sys_users(status);
CREATE INDEX idx_user_created_at ON sys_users(created_at);

CREATE INDEX idx_role_status ON sys_roles(status);
CREATE INDEX idx_api_path ON sys_apis(path);

-- 查詢最佳化
-- 啟用查詢快取
SET GLOBAL query_cache_size = 268435456;
SET GLOBAL query_cache_type = ON;

-- 慢查詢日誌
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 2;
SET GLOBAL log_queries_not_using_indexes = 'ON';
```

### 3. 系統效能監控

```bash
#!/bin/bash
# 效能監控腳本

check_performance() {
    echo "=== 系統效能檢查 ==="

    # CPU 使用率
    echo "CPU 使用率:"
    top -bn1 | grep "Cpu(s)" | awk '{print $2}' | cut -d'%' -f1

    # 記憶體使用
    echo "記憶體使用:"
    free -h

    # 磁碟 I/O
    echo "磁碟 I/O:"
    iostat -x 1 1

    # 網路連接
    echo "網路連接數:"
    netstat -an | grep :80 | wc -l

    # 應用程式狀態
    echo "Go-Admin 進程:"
    ps aux | grep go-admin | grep -v grep

    # 資料庫連接
    echo "MySQL 連接數:"
    mysql -e "SHOW STATUS LIKE 'Threads_connected';" 2>/dev/null
}

check_performance
```

### 4. 效能調優建議

```markdown
## Go-Admin 效能最佳化建議

### 應用層面

1. **連接池配置**: 根據併發量調整資料庫和 Redis 連接池大小
2. **快取策略**: 實施多層快取（記憶體 + Redis）
3. **異步處理**: 使用訊息佇列處理耗時操作
4. **靜態檔案**: 使用 CDN 分發靜態資源

### 資料庫層面

1. **索引最佳化**: 為常用查詢欄位建立適當索引
2. **查詢最佳化**: 避免 N+1 查詢，使用聯表查詢
3. **讀寫分離**: 配置主從複製，讀取操作使用從庫
4. **分庫分表**: 對於大數據量場景考慮分庫分表

### 系統層面

1. **硬體配置**: SSD 硬碟，充足的記憶體
2. **網路最佳化**: 使用 HTTP/2，啟用 Gzip 壓縮
3. **負載均衡**: 多實例部署，使用 Nginx 負載均衡
4. **監控告警**: 完善的監控體系，及時發現問題
```

## 故障排除

### 1. 常見問題診斷

```bash
#!/bin/bash
# 故障診斷腳本

diagnose_issue() {
    echo "=== Go-Admin 故障診斷 ==="

    # 檢查服務狀態
    echo "1. 檢查服務狀態:"
    systemctl status go-admin

    # 檢查端口監聽
    echo "2. 檢查端口監聽:"
    netstat -tlnp | grep :8000

    # 檢查日誌錯誤
    echo "3. 檢查應用程式日誌:"
    tail -n 50 /opt/go-admin/logs/go-admin.log | grep -i error

    # 檢查系統日誌
    echo "4. 檢查系統日誌:"
    journalctl -u go-admin -n 20 --no-pager

    # 檢查資源使用
    echo "5. 檢查資源使用:"
    ps aux | grep go-admin | head -1

    # 檢查磁碟空間
    echo "6. 檢查磁碟空間:"
    df -h /opt/go-admin

    # 檢查資料庫連接
    echo "7. 檢查資料庫連接:"
    mysql -h localhost -u go_admin_user -p -e "SELECT 1;" 2>&1

    # 檢查 Redis 連接
    echo "8. 檢查 Redis 連接:"
    redis-cli ping 2>&1
}

diagnose_issue
```

### 2. 緊急恢復程序

```bash
#!/bin/bash
# 緊急恢復腳本

emergency_recovery() {
    echo "🚨 執行緊急恢復程序..."

    # 停止應用程式
    echo "停止應用程式服務..."
    systemctl stop go-admin

    # 檢查磁碟空間
    echo "檢查磁碟空間..."
    df -h /opt/go-admin

    # 清理臨時檔案
    echo "清理臨時檔案..."
    find /opt/go-admin/logs -name "*.log" -mtime +7 -delete
    find /opt/go-admin/data/tmp -type f -mtime +1 -delete

    # 重置設定檔案
    echo "檢查設定檔案..."
    if [ ! -f /opt/go-admin/app/config/settings.yml ]; then
        echo "❌ 設定檔案遺失，從備份恢復..."
        cp /opt/go-admin/backup/config/settings.yml /opt/go-admin/app/config/
    fi

    # 檢查資料庫連接
    echo "檢查資料庫連接..."
    mysql -h localhost -u go_admin_user -p -e "SELECT 1;" > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        echo "❌ 資料庫連接失敗，請檢查資料庫服務"
        return 1
    fi

    # 重啟應用程式
    echo "重啟應用程式服務..."
    systemctl start go-admin

    # 等待服務啟動
    sleep 5

    # 檢查服務狀態
    if systemctl is-active --quiet go-admin; then
        echo "✅ 緊急恢復成功，服務已正常運行"
    else
        echo "❌ 緊急恢復失敗，請手動檢查"
        return 1
    fi
}

emergency_recovery
```

### 3. 效能問題排除

```bash
#!/bin/bash
# 效能問題排除腳本

troubleshoot_performance() {
    echo "=== 效能問題排除 ==="

    # 檢查 CPU 使用率
    echo "1. 檢查 CPU 使用率:"
    top -bn1 | head -20

    # 檢查記憶體使用
    echo "2. 檢查記憶體使用:"
    free -h
    cat /proc/meminfo | grep -E "(MemTotal|MemFree|MemAvailable|Cached)"

    # 檢查磁碟 I/O
    echo "3. 檢查磁碟 I/O:"
    iostat -x 1 3

    # 檢查網路連接
    echo "4. 檢查網路連接:"
    ss -tuln | grep -E ":80|:443|:8000|:3306|:6379"

    # 檢查應用程式效能
    echo "5. 檢查應用程式回應時間:"
    curl -w "@curl-format.txt" -o /dev/null -s http://localhost:8000/api/v1/health

    # 檢查資料庫效能
    echo "6. 檢查 MySQL 慢查詢:"
    mysql -e "SHOW STATUS LIKE 'Slow_queries';"

    # 檢查 Redis 效能
    echo "7. 檢查 Redis 統計:"
    redis-cli info stats | grep -E "(total_commands_processed|expired_keys|evicted_keys)"
}

# 建立 curl 格式檔案
cat > curl-format.txt << 'EOF'
     time_namelookup:  %{time_namelookup}\n
        time_connect:  %{time_connect}\n
     time_appconnect:  %{time_appconnect}\n
    time_pretransfer:  %{time_pretransfer}\n
       time_redirect:  %{time_redirect}\n
  time_starttransfer:  %{time_starttransfer}\n
                     ----------\n
          time_total:  %{time_total}\n
EOF

troubleshoot_performance
```

---

## 部署清單

### 上線前檢查清單

- [ ] 系統環境準備完成
- [ ] 資料庫配置和初始化完成
- [ ] 應用程式編譯和部署完成
- [ ] Nginx 反向代理配置完成
- [ ] SSL 證書安裝和配置完成
- [ ] 防火牆和安全配置完成
- [ ] 監控和告警配置完成
- [ ] 日誌收集配置完成
- [ ] 備份策略實施完成
- [ ] 效能測試通過
- [ ] 安全檢測通過
- [ ] 文檔更新完成

### 上線後驗證清單

- [ ] 應用程式服務正常運行
- [ ] 資料庫連接正常
- [ ] Redis 快取正常
- [ ] 前端頁面載入正常
- [ ] API 介面回應正常
- [ ] 使用者登入功能正常
- [ ] 檔案上傳功能正常
- [ ] SSL 證書有效
- [ ] 監控指標正常
- [ ] 告警通知正常
- [ ] 日誌記錄正常
- [ ] 備份任務正常

## 維護建議

### 日常維護

1. **每日檢查**

   - 服務狀態監控
   - 系統資源使用率
   - 應用程式日誌檢查
   - 備份任務執行狀態

2. **每週檢查**

   - 安全更新檢查
   - 效能指標分析
   - 磁碟空間清理
   - 資料庫效能最佳化

3. **每月檢查**
   - 系統安全掃描
   - 備份恢復測試
   - 災難恢復演練
   - 容量規劃評估

### 升級建議

1. **應用程式升級**

   - 在預生產環境充分測試
   - 制定回滾計劃
   - 選擇低峰時段執行
   - 監控升級後的系統狀態

2. **系統升級**
   - 定期更新作業系統補丁
   - 升級中介軟體版本
   - 更新監控工具
   - 保持安全工具最新

通過遵循本部署指南，您可以建立一個穩定、安全、高效能的 Go-Admin 生產環境。建議在實際部署前在測試環境中充分驗證所有配置和腳本。
