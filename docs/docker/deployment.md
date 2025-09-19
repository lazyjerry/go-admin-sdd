# Go-Admin Docker 部署指南

本指南提供 Go-Admin 系統的 Docker 容器化部署詳細步驟，包含單容器和多容器部署方案。

## 快速部署

### 單容器部署 (SQLite)

適合快速體驗和開發測試：

```bash
# 1. 拉取映像
docker pull go-admin/go-admin:latest

# 2. 建立資料目錄
mkdir -p ./go-admin-data/{config,storage,logs}

# 3. 建立設定檔案
cat > ./go-admin-data/config/settings.yml << EOF
settings:
  application:
    mode: prod
    name: go-admin
    port: 8000
    host: 0.0.0.0
  database:
    dbtype: sqlite3
    source: /app/storage/go-admin.db
  jwt:
    secret: your-jwt-secret-key
  log:
    level: info
    path: /app/logs
EOF

# 4. 啟動容器
docker run -d \
  --name go-admin \
  -p 8000:8000 \
  -v $(pwd)/go-admin-data/config:/app/config \
  -v $(pwd)/go-admin-data/storage:/app/storage \
  -v $(pwd)/go-admin-data/logs:/app/logs \
  go-admin/go-admin:latest

# 5. 初始化資料庫
docker exec go-admin ./go-admin migrate -c config/settings.yml

# 6. 驗證部署
curl http://localhost:8000/api/v1/health
```

### Docker Compose 部署 (推薦)

完整的生產環境部署：

```bash
# 1. 下載部署配置
curl -L https://github.com/go-admin-team/go-admin/raw/master/docker-compose.yml -o docker-compose.yml

# 2. 建立環境變數檔案
cat > .env << EOF
# 應用設定
APP_ENV=production
APP_PORT=8000

# 資料庫設定
DB_HOST=mysql
DB_PORT=3306
DB_NAME=go_admin
DB_USERNAME=go_admin
DB_PASSWORD=SecurePassword123

# Redis 設定
REDIS_HOST=redis
REDIS_PORT=6379

# JWT 設定
JWT_SECRET=your-production-jwt-secret-key-here
EOF

# 3. 啟動所有服務
docker-compose up -d

# 4. 檢查服務狀態
docker-compose ps

# 5. 初始化資料庫
docker-compose exec go-admin ./go-admin migrate -c config/settings.yml

# 6. 建立管理員帳號
docker-compose exec go-admin ./go-admin user create \
  --username admin \
  --password admin123 \
  --email admin@example.com
```

## Docker Compose 配置

### 基礎配置檔案

```yaml
# docker-compose.yml
version: "3.8"

services:
  # Go-Admin 應用
  go-admin:
    image: go-admin/go-admin:latest
    container_name: go-admin
    restart: unless-stopped
    ports:
      - "${APP_PORT:-8000}:8000"
    environment:
      - DB_HOST=${DB_HOST}
      - DB_PORT=${DB_PORT}
      - DB_NAME=${DB_NAME}
      - DB_USERNAME=${DB_USERNAME}
      - DB_PASSWORD=${DB_PASSWORD}
      - REDIS_HOST=${REDIS_HOST}
      - REDIS_PORT=${REDIS_PORT}
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./config:/app/config
      - app_storage:/app/storage
      - app_logs:/app/logs
    networks:
      - go-admin-net

  # MySQL 資料庫
  mysql:
    image: mysql:8.0
    container_name: go-admin-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_DATABASE: ${DB_NAME}
      MYSQL_USER: ${DB_USERNAME}
      MYSQL_PASSWORD: ${DB_PASSWORD}
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10
    networks:
      - go-admin-net

  # Redis 快取
  redis:
    image: redis:7-alpine
    container_name: go-admin-redis
    restart: unless-stopped
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    networks:
      - go-admin-net

volumes:
  mysql_data:
  redis_data:
  app_storage:
  app_logs:

networks:
  go-admin-net:
    driver: bridge
```

### 生產環境配置

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  go-admin:
    image: go-admin/go-admin:v2.1.0 # 指定穩定版本
    deploy:
      replicas: 2
      resources:
        limits:
          cpus: "1.0"
          memory: 1G
        reservations:
          cpus: "0.5"
          memory: 512M
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/api/v1/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  mysql:
    image: mysql:8.0
    command: >
      --default-authentication-plugin=mysql_native_password
      --innodb-buffer-pool-size=256M
      --max-connections=200
      --slow-query-log=1
      --slow-query-log-file=/var/log/mysql/slow.log
      --long-query-time=2
    volumes:
      - type: volume
        source: mysql_data
        target: /var/lib/mysql
      - type: bind
        source: ./mysql/my.cnf
        target: /etc/mysql/conf.d/my.cnf

  nginx:
    image: nginx:alpine
    container_name: go-admin-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
      - app_storage:/var/www/storage:ro
    depends_on:
      - go-admin
    networks:
      - go-admin-net
```

## 自定義 Dockerfile

### 最佳化 Dockerfile

```dockerfile
# Multi-stage build for Go-Admin
FROM golang:1.21-alpine AS builder

# 安裝必要的系統套件
RUN apk add --no-cache git ca-certificates tzdata gcc musl-dev

# 設定工作目錄
WORKDIR /app

# 複製 go mod 檔案並下載依賴
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# 複製原始碼
COPY . .

# 建構應用程式
RUN CGO_ENABLED=1 GOOS=linux go build \
    -ldflags="-w -s -X main.version=$(git describe --tags --always)" \
    -a -installsuffix cgo \
    -o go-admin main.go

# Final stage
FROM alpine:latest

# 安裝執行時依賴
RUN apk --no-cache add ca-certificates tzdata curl

# 設定時區
ENV TZ=Asia/Taipei
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# 建立應用使用者
RUN addgroup -g 1001 -S go-admin && \
    adduser -u 1001 -S go-admin -G go-admin

# 建立應用目錄
WORKDIR /app

# 從建構階段複製檔案
COPY --from=builder /app/go-admin .
COPY --from=builder /app/config ./config
COPY --from=builder /app/template ./template

# 建立必要目錄並設定權限
RUN mkdir -p storage/{logs,uploads} && \
    chown -R go-admin:go-admin /app

# 切換到非特權使用者
USER go-admin

# 暴露服務埠
EXPOSE 8000

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8000/api/v1/health || exit 1

# 預設啟動命令
CMD ["./go-admin", "server", "-c", "config/settings.yml"]
```

### 建構腳本

```bash
#!/bin/bash
# build.sh

set -e

# 設定變數
IMAGE_NAME="go-admin"
TAG=${1:-latest}
REGISTRY=${REGISTRY:-"go-admin"}

echo "開始建構 Docker 映像..."

# 建構映像
docker build \
    --build-arg BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') \
    --build-arg VCS_REF=$(git rev-parse --short HEAD) \
    --build-arg VERSION=${TAG} \
    -t ${REGISTRY}/${IMAGE_NAME}:${TAG} \
    -t ${REGISTRY}/${IMAGE_NAME}:latest \
    .

echo "映像建構完成: ${REGISTRY}/${IMAGE_NAME}:${TAG}"

# 檢視映像大小
docker images ${REGISTRY}/${IMAGE_NAME}:${TAG}

# 安全掃描 (可選)
if command -v trivy &> /dev/null; then
    echo "執行安全掃描..."
    trivy image ${REGISTRY}/${IMAGE_NAME}:${TAG}
fi
```

## 網路和安全配置

### Nginx 反向代理

```nginx
# nginx/nginx.conf
upstream go-admin-backend {
    least_conn;
    server go-admin:8000 max_fails=3 fail_timeout=30s;
    # 如果有多個實例
    # server go-admin-2:8000 max_fails=3 fail_timeout=30s;
}

server {
    listen 80;
    server_name your-domain.com;

    # 強制 HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    # SSL 配置
    ssl_certificate /etc/nginx/ssl/cert.pem;
    ssl_certificate_key /etc/nginx/ssl/key.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;

    # 安全標頭
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";

    # 檔案上傳限制
    client_max_body_size 10M;

    # API 代理
    location /api/ {
        proxy_pass http://go-admin-backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # 超時設定
        proxy_connect_timeout 30s;
        proxy_send_timeout 60s;
        proxy_read_timeout 60s;
    }

    # 靜態檔案
    location /uploads/ {
        alias /var/www/storage/uploads/;
        expires 30d;
        add_header Cache-Control "public, immutable";
    }

    # 健康檢查
    location /health {
        access_log off;
        return 200 "healthy\n";
        add_header Content-Type text/plain;
    }
}
```

### 防火牆規則

```bash
# UFW 防火牆設定
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp

# 僅允許內部網路存取資料庫
sudo ufw allow from 172.18.0.0/16 to any port 3306
sudo ufw allow from 172.18.0.0/16 to any port 6379

sudo ufw enable
```

## 監控和日誌

### 日誌配置

```yaml
# 在 docker-compose.yml 中添加日誌配置
services:
  go-admin:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
        labels: "service=go-admin,env=production"
```

### 監控指標

```bash
# 監控腳本
#!/bin/bash
# monitor.sh

echo "=== Docker 容器狀態 ==="
docker-compose ps

echo -e "\n=== 系統資源使用 ==="
docker stats --no-stream

echo -e "\n=== 應用健康檢查 ==="
curl -s http://localhost:8000/api/v1/health | jq .

echo -e "\n=== 資料庫連接 ==="
docker-compose exec mysql mysqladmin ping

echo -e "\n=== Redis 狀態 ==="
docker-compose exec redis redis-cli ping
```

## 備份策略

### 自動備份腳本

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backups/$(date +%Y%m%d)"
mkdir -p $BACKUP_DIR

# 備份資料庫
echo "備份 MySQL 資料庫..."
docker-compose exec mysql mysqldump \
    -u root -p$DB_PASSWORD go_admin \
    | gzip > $BACKUP_DIR/mysql_backup.sql.gz

# 備份應用檔案
echo "備份應用檔案..."
docker run --rm \
    -v go-admin_app_storage:/source \
    -v $BACKUP_DIR:/backup \
    alpine tar czf /backup/storage_backup.tar.gz -C /source .

# 清理舊備份 (保留 7 天)
find /backups -type d -mtime +7 -exec rm -rf {} +

echo "備份完成: $BACKUP_DIR"
```

### Cron 定時任務

```bash
# 設定定時備份
crontab -e

# 每日凌晨 2 點執行備份
0 2 * * * /path/to/backup.sh >> /var/log/backup.log 2>&1
```

## 故障排除

### 常見問題解決

```bash
# 1. 檢視容器日誌
docker-compose logs -f go-admin

# 2. 檢查容器健康狀態
docker inspect --format='{{.State.Health.Status}}' go-admin

# 3. 測試網路連通性
docker-compose exec go-admin ping mysql
docker-compose exec go-admin nc -zv mysql 3306

# 4. 檢查資料庫
docker-compose exec mysql mysql -u go_admin -p -e "SHOW TABLES;"

# 5. 重新啟動服務
docker-compose restart go-admin

# 6. 清理並重建
docker-compose down -v
docker-compose up --build -d
```

---

透過 Docker 部署 Go-Admin 可以實現快速、一致且可擴展的部署方案，適合各種環境需求。

相關文件：

- [Docker Compose 配置詳解](./docker-compose.md)
- [容器監控指南](./monitoring.md)
- [生產環境最佳實踐](../deployment/production.md)
