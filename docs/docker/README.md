# Go-Admin Docker 容器化指南

本文件提供 Go-Admin 後台管理系統的 Docker 容器化部署方案，包含單容器部署、Docker Compose 編排和生產環境最佳實踐。

## Docker 方案概覽

Go-Admin 提供多種 Docker 部署方案，適用於不同的使用場景：

### 部署方案對比

| 方案               | 適用場景   | 複雜度 | 特點                 |
| ------------------ | ---------- | ------ | -------------------- |
| **單容器部署**     | 開發測試   | 簡單   | 快速啟動，適合演示   |
| **Docker Compose** | 本地開發   | 中等   | 包含資料庫，完整環境 |
| **Kubernetes**     | 生產環境   | 複雜   | 高可用，自動擴展     |
| **Docker Swarm**   | 中小型生產 | 中等   | 集群管理，負載均衡   |

### 容器架構

```
┌─────────────────────────────────────────────┐
│                Load Balancer                │
│              (Nginx/Traefik)               │
└─────────────────────────────────────────────┘
                    │
    ┌───────────────┼───────────────┐
    │               │               │
┌─────────┐   ┌─────────┐   ┌─────────┐
│Go-Admin │   │Go-Admin │   │Go-Admin │
│Instance1│   │Instance2│   │Instance3│
└─────────┘   └─────────┘   └─────────┘
    │               │               │
    └───────────────┼───────────────┘
                    │
┌─────────────────────────────────────────────┐
│              Shared Services                │
├─────────────┬─────────────┬─────────────────┤
│   MySQL     │    Redis    │   File Storage  │
│  Database   │    Cache    │   (MinIO/NFS)   │
└─────────────┴─────────────┴─────────────────┘
```

## 快速開始

### 1. 使用預建映像 (推薦)

```bash
# 拉取官方映像
docker pull go-admin/go-admin:latest

# 快速啟動 (使用 SQLite)
docker run -d \
  --name go-admin \
  -p 8000:8000 \
  -v $(pwd)/data:/app/data \
  go-admin/go-admin:latest
```

### 2. 使用 Docker Compose (完整環境)

```bash
# 下載 docker-compose.yml
curl -O https://raw.githubusercontent.com/go-admin-team/go-admin/master/docker-compose.yml

# 啟動完整環境
docker-compose up -d

# 檢查服務狀態
docker-compose ps
```

### 3. 自行建構映像

```bash
# 複製專案
git clone https://github.com/go-admin-team/go-admin.git
cd go-admin

# 建構映像
docker build -t go-admin:local .

# 執行容器
docker run -d \
  --name go-admin-local \
  -p 8000:8000 \
  go-admin:local
```

## Dockerfile 解析

### 多階段建構

Go-Admin 使用多階段建構來最佳化映像大小：

```dockerfile
# Stage 1: 建構階段
FROM golang:1.21-alpine AS builder

# 設定工作目錄
WORKDIR /app

# 複製依賴檔案
COPY go.mod go.sum ./

# 下載依賴
RUN go mod download

# 複製原始碼
COPY . .

# 建構二進位檔案
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -o go-admin main.go

# Stage 2: 執行階段
FROM alpine:latest

# 安裝必要套件
RUN apk --no-cache add ca-certificates tzdata

# 設定時區
ENV TZ=Asia/Taipei

# 建立非 root 使用者
RUN addgroup -g 1001 -S go-admin && \
    adduser -u 1001 -S go-admin -G go-admin

# 設定工作目錄
WORKDIR /app

# 從建構階段複製執行檔
COPY --from=builder /app/go-admin .

# 複製設定檔案
COPY --from=builder /app/config ./config

# 建立必要目錄
RUN mkdir -p storage/logs storage/uploads && \
    chown -R go-admin:go-admin /app

# 切換到非 root 使用者
USER go-admin

# 暴露埠號
EXPOSE 8000

# 健康檢查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8000/api/v1/health || exit 1

# 啟動命令
CMD ["./go-admin", "server", "-c", "config/settings.yml"]
```

### 建構最佳化技巧

#### 1. 使用 .dockerignore

```bash
# .dockerignore
.git
.gitignore
README.md
Dockerfile
.dockerignore
node_modules
npm-debug.log
storage/logs/*
storage/uploads/*
docs/
test/
*.md
```

#### 2. 多階段建構優勢

- **減少映像大小**: 最終映像只包含執行檔案
- **安全性提升**: 不包含建構工具和原始碼
- **快取最佳化**: 分層建構提升重建速度

#### 3. 映像大小對比

```bash
# 單階段建構 (包含完整 Go 環境)
golang:1.21        ~800MB

# 多階段建構 (僅包含執行檔)
go-admin:latest    ~20MB
```

## Docker Compose 部署

### 基本 Compose 配置

```yaml
# docker-compose.yml
version: "3.8"

services:
  # Go-Admin 應用服務
  go-admin:
    image: go-admin/go-admin:latest
    container_name: go-admin-app
    restart: unless-stopped
    ports:
      - "8000:8000"
    environment:
      - DB_HOST=mysql
      - DB_PORT=3306
      - DB_USERNAME=go_admin
      - DB_PASSWORD=your_password
      - DB_NAME=go_admin
      - REDIS_HOST=redis
      - REDIS_PORT=6379
    depends_on:
      mysql:
        condition: service_healthy
      redis:
        condition: service_started
    volumes:
      - ./config:/app/config
      - app_storage:/app/storage
    networks:
      - go-admin-network

  # MySQL 資料庫服務
  mysql:
    image: mysql:8.0
    container_name: go-admin-mysql
    restart: unless-stopped
    environment:
      MYSQL_ROOT_PASSWORD: root_password
      MYSQL_DATABASE: go_admin
      MYSQL_USER: go_admin
      MYSQL_PASSWORD: your_password
    ports:
      - "3306:3306"
    volumes:
      - mysql_data:/var/lib/mysql
      - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
      timeout: 20s
      retries: 10
    networks:
      - go-admin-network

  # Redis 快取服務
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
      - go-admin-network

  # Nginx 反向代理
  nginx:
    image: nginx:alpine
    container_name: go-admin-nginx
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - app_storage:/var/www/storage
    depends_on:
      - go-admin
    networks:
      - go-admin-network

volumes:
  mysql_data:
  redis_data:
  app_storage:

networks:
  go-admin-network:
    driver: bridge
```

### 環境變數配置

建立 `.env` 檔案：

```bash
# .env
# 應用設定
APP_ENV=production
APP_DEBUG=false
APP_URL=https://your-domain.com

# 資料庫設定
DB_CONNECTION=mysql
DB_HOST=mysql
DB_PORT=3306
DB_DATABASE=go_admin
DB_USERNAME=go_admin
DB_PASSWORD=secure_password_here

# Redis 設定
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT 設定
JWT_SECRET=your-jwt-secret-key-here

# 檔案儲存
STORAGE_TYPE=local
STORAGE_PATH=/app/storage/uploads

# 郵件設定
MAIL_MAILER=smtp
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-app-password
```

### 進階 Compose 配置

```yaml
# docker-compose.prod.yml
version: "3.8"

services:
  go-admin:
    image: go-admin/go-admin:latest
    deploy:
      replicas: 3
      resources:
        limits:
          cpus: "0.5"
          memory: 512M
        reservations:
          cpus: "0.25"
          memory: 256M
      restart_policy:
        condition: on-failure
        delay: 5s
        max_attempts: 3
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s

  mysql:
    image: mysql:8.0
    command: --default-authentication-plugin=mysql_native_password
    environment:
      MYSQL_ROOT_PASSWORD_FILE: /run/secrets/mysql_root_password
      MYSQL_PASSWORD_FILE: /run/secrets/mysql_password
    secrets:
      - mysql_root_password
      - mysql_password
    volumes:
      - type: volume
        source: mysql_data
        target: /var/lib/mysql
      - type: bind
        source: ./mysql/my.cnf
        target: /etc/mysql/conf.d/my.cnf

secrets:
  mysql_root_password:
    external: true
  mysql_password:
    external: true
```

## 生產環境部署

### 1. 安全性配置

#### 密鑰管理

```bash
# 建立 Docker secrets
echo "secure_root_password" | docker secret create mysql_root_password -
echo "secure_app_password" | docker secret create mysql_password -
echo "jwt_secret_key_256_bits" | docker secret create jwt_secret -

# 在 Compose 中使用 secrets
services:
  go-admin:
    secrets:
      - jwt_secret
    environment:
      - JWT_SECRET_FILE=/run/secrets/jwt_secret
```

#### 網路隔離

```yaml
networks:
  frontend:
    driver: overlay
    external: true
  backend:
    driver: overlay
    internal: true # 內部網路，不對外

services:
  nginx:
    networks:
      - frontend

  go-admin:
    networks:
      - frontend
      - backend

  mysql:
    networks:
      - backend # 只能從後端存取
```

### 2. 效能最佳化

#### 資源限制

```yaml
services:
  go-admin:
    deploy:
      resources:
        limits:
          cpus: "1.0"
          memory: 1G
        reservations:
          cpus: "0.5"
          memory: 512M
```

#### 快取配置

```yaml
redis:
  image: redis:7-alpine
  command: >
    redis-server
    --maxmemory 256mb
    --maxmemory-policy allkeys-lru
    --save 900 1
    --save 300 10
    --save 60 10000
```

### 3. 監控和日誌

#### 日誌聚合

```yaml
version: "3.8"

services:
  go-admin:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"

  # ELK Stack 日誌收集
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.5.0
    environment:
      - discovery.type=single-node
      - "ES_JAVA_OPTS=-Xms512m -Xmx512m"

  logstash:
    image: docker.elastic.co/logstash/logstash:8.5.0
    volumes:
      - ./logstash/config:/usr/share/logstash/pipeline

  kibana:
    image: docker.elastic.co/kibana/kibana:8.5.0
    ports:
      - "5601:5601"
```

#### 監控指標

```yaml
# Prometheus 監控
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./prometheus/prometheus.yml:/etc/prometheus/prometheus.yml

# Grafana 儀表板
grafana:
  image: grafana/grafana
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
  volumes:
    - grafana_data:/var/lib/grafana
```

## 常用 Docker 指令

### 容器管理

```bash
# 檢視執行中的容器
docker ps

# 檢視所有容器 (包含停止的)
docker ps -a

# 檢視容器日誌
docker logs go-admin-app

# 即時跟蹤日誌
docker logs -f go-admin-app

# 進入容器 Shell
docker exec -it go-admin-app /bin/sh

# 檢視容器資源使用情況
docker stats go-admin-app
```

### 映像管理

```bash
# 列出本地映像
docker images

# 刪除映像
docker rmi go-admin:latest

# 清理未使用的映像
docker image prune

# 建構並標記映像
docker build -t go-admin:v1.0.0 .

# 推送映像到倉庫
docker push go-admin/go-admin:v1.0.0
```

### Docker Compose 指令

```bash
# 啟動所有服務
docker-compose up -d

# 停止所有服務
docker-compose down

# 重新建構並啟動
docker-compose up --build -d

# 檢視服務狀態
docker-compose ps

# 檢視服務日誌
docker-compose logs go-admin

# 執行一次性指令
docker-compose exec go-admin ./go-admin migrate

# 擴展服務實例
docker-compose up --scale go-admin=3 -d
```

## 故障排除

### 常見問題

#### 1. 容器無法啟動

```bash
# 檢查容器日誌
docker logs go-admin-app

# 檢查容器配置
docker inspect go-admin-app

# 檢查映像是否存在
docker images | grep go-admin
```

#### 2. 資料庫連接失敗

```bash
# 檢查網路連通性
docker exec go-admin-app ping mysql

# 檢查資料庫服務
docker exec mysql-container mysqladmin ping

# 測試資料庫連接
docker exec mysql-container mysql -u go_admin -p -e "SHOW DATABASES;"
```

#### 3. 埠號衝突

```bash
# 檢查埠號使用情況
netstat -tulpn | grep 8000

# 修改埠號對映
docker run -p 8080:8000 go-admin:latest
```

#### 4. 效能問題

```bash
# 檢視容器資源使用
docker stats

# 檢查系統資源
free -h
df -h

# 最佳化配置
docker update --memory=1g --cpus="1.5" go-admin-app
```

### 除錯技巧

#### 1. 使用除錯模式

```bash
# 以除錯模式啟動
docker run -it --rm \
  -e APP_DEBUG=true \
  -e LOG_LEVEL=debug \
  go-admin:latest
```

#### 2. 掛載本地程式碼

```bash
# 開發模式掛載
docker run -it --rm \
  -p 8000:8000 \
  -v $(pwd):/app \
  golang:1.21 \
  bash -c "cd /app && go run main.go"
```

#### 3. 網路診斷

```bash
# 檢查容器間網路
docker network ls
docker network inspect bridge

# 測試容器間連通性
docker exec container1 ping container2
docker exec container1 curl http://container2:8000/health
```

## 備份和恢復

### 資料備份

```bash
# 備份 MySQL 資料
docker exec mysql-container mysqldump \
  -u root -p go_admin > backup_$(date +%Y%m%d_%H%M%S).sql

# 備份 Volume 資料
docker run --rm \
  -v mysql_data:/source \
  -v $(pwd)/backups:/backup \
  alpine tar czf /backup/mysql_backup.tar.gz -C /source .

# 備份應用檔案
docker cp go-admin-app:/app/storage ./backups/storage
```

### 資料恢復

```bash
# 恢復 MySQL 資料
docker exec -i mysql-container mysql -u root -p go_admin < backup.sql

# 恢復 Volume 資料
docker run --rm \
  -v mysql_data:/target \
  -v $(pwd)/backups:/backup \
  alpine tar xzf /backup/mysql_backup.tar.gz -C /target

# 恢復應用檔案
docker cp ./backups/storage go-admin-app:/app/
```

## 升級策略

### 滾動升級

```bash
# 1. 拉取新版映像
docker pull go-admin/go-admin:v2.0.0

# 2. 備份資料
docker-compose exec mysql mysqldump go_admin > backup.sql

# 3. 更新 Compose 檔案
sed -i 's/go-admin:latest/go-admin:v2.0.0/g' docker-compose.yml

# 4. 滾動更新
docker-compose up -d --no-deps go-admin

# 5. 驗證服務
curl http://localhost:8000/api/v1/health
```

### 藍綠部署

```bash
# 1. 啟動新版本 (綠色環境)
docker-compose -f docker-compose.green.yml up -d

# 2. 健康檢查
./scripts/health-check.sh http://localhost:8001

# 3. 切換流量 (更新 Nginx 配置)
docker exec nginx nginx -s reload

# 4. 停止舊版本 (藍色環境)
docker-compose -f docker-compose.blue.yml down
```

---

Docker 容器化為 Go-Admin 提供了靈活、可擴展的部署方案，適合從開發測試到生產環境的各種需求。

相關文件：

- [Docker 部署詳細指南](./deployment.md)
- [Docker Compose 配置說明](./docker-compose.md)
- [容器監控指南](./monitoring.md)
- [生產環境最佳實踐](../deployment/production.md)
