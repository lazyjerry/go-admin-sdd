# Docker Compose 配置說明

本文件詳細說明 Go-Admin 專案的 Docker Compose 配置，包含服務組態、網路設定、儲存管理等方面。

## 概要

Go-Admin 提供完整的 Docker Compose 配置，可快速啟動包含前端、後端和資料庫在內的完整開發環境。

## 配置檔案結構

```yaml
# docker-compose.yml - 主要配置檔案
version: "3.8"

services:
  # Go-Admin 後端服務
  go-admin:
    build: .
    ports:
      - "8000:8000"
    environment:
      - GIN_MODE=release
      - DB_HOST=mysql
    depends_on:
      - mysql
    networks:
      - go-admin-network

  # MySQL 資料庫服務
  mysql:
    image: mysql:8.0
    ports:
      - "3306:3306"
    environment:
      MYSQL_ROOT_PASSWORD: 123456
      MYSQL_DATABASE: go-admin
    volumes:
      - mysql-data:/var/lib/mysql
    networks:
      - go-admin-network

  # Redis 快取服務
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    networks:
      - go-admin-network

volumes:
  mysql-data:

networks:
  go-admin-network:
    driver: bridge
```

## 服務詳細配置

### Go-Admin 後端服務

#### 基本配置

```yaml
go-admin:
  build:
    context: .
    dockerfile: Dockerfile
  container_name: go-admin-backend
  restart: unless-stopped
  ports:
    - "8000:8000" # HTTP 服務埠
  environment:
    # 應用程式配置
    - GIN_MODE=release
    - APP_ENV=production
    - APP_DEBUG=false

    # 資料庫配置
    - DB_DRIVER=mysql
    - DB_HOST=mysql
    - DB_PORT=3306
    - DB_NAME=go-admin
    - DB_USER=root
    - DB_PASSWORD=123456

    # Redis 配置
    - REDIS_HOST=redis
    - REDIS_PORT=6379
    - REDIS_PASSWORD=

    # JWT 配置
    - JWT_SECRET=your-secret-key
    - JWT_TIMEOUT=7200
```

#### 進階配置選項

```yaml
# 健康檢查
healthcheck:
  test: ["CMD", "curl", "-f", "http://localhost:8000/api/health"]
  interval: 30s
  timeout: 10s
  retries: 3
  start_period: 40s

# 資源限制
deploy:
  resources:
    limits:
      cpus: "1.0"
      memory: 512M
    reservations:
      cpus: "0.5"
      memory: 256M

# 日誌配置
logging:
  driver: "json-file"
  options:
    max-size: "10m"
    max-file: "3"
```

### MySQL 資料庫服務

#### 基本配置

```yaml
mysql:
  image: mysql:8.0
  container_name: go-admin-mysql
  restart: unless-stopped
  ports:
    - "3306:3306"
  environment:
    # 必要的環境變數
    MYSQL_ROOT_PASSWORD: 123456
    MYSQL_DATABASE: go-admin
    MYSQL_USER: admin
    MYSQL_PASSWORD: admin123

    # MySQL 配置參數
    MYSQL_CHARACTER_SET_SERVER: utf8mb4
    MYSQL_COLLATION_SERVER: utf8mb4_unicode_ci
```

#### 進階 MySQL 配置

```yaml
# 自訂 MySQL 配置
command: >
  --default-authentication-plugin=mysql_native_password
  --character-set-server=utf8mb4
  --collation-server=utf8mb4_unicode_ci
  --sql_mode=STRICT_TRANS_TABLES,NO_ZERO_DATE,NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO
  --max_connections=1000
  --innodb_buffer_pool_size=256M

# 資料持久化
volumes:
  - mysql-data:/var/lib/mysql
  - ./config/mysql/my.cnf:/etc/mysql/conf.d/my.cnf:ro
  - ./scripts/init.sql:/docker-entrypoint-initdb.d/init.sql:ro

# 健康檢查
healthcheck:
  test: ["CMD", "mysqladmin", "ping", "-h", "localhost"]
  timeout: 20s
  retries: 10
```

### Redis 快取服務

#### 基本配置

```yaml
redis:
  image: redis:7-alpine
  container_name: go-admin-redis
  restart: unless-stopped
  ports:
    - "6379:6379"
  command: redis-server --appendonly yes
```

#### Redis 進階配置

```yaml
# 自訂 Redis 配置
command: >
  redis-server
  --maxmemory 256mb
  --maxmemory-policy allkeys-lru
  --appendonly yes
  --appendfsync everysec

# 資料持久化
volumes:
  - redis-data:/data
  - ./config/redis/redis.conf:/usr/local/etc/redis/redis.conf:ro

# 健康檢查
healthcheck:
  test: ["CMD", "redis-cli", "ping"]
  interval: 30s
  timeout: 3s
  retries: 3
```

## 網路配置

### 預設網路設定

```yaml
networks:
  go-admin-network:
    driver: bridge
    ipam:
      config:
        - subnet: 172.20.0.0/16
```

### 自訂網路配置

```yaml
networks:
  # 前端網路
  frontend:
    driver: bridge
    internal: false

  # 後端網路
  backend:
    driver: bridge
    internal: true

  # 資料庫網路
  database:
    driver: bridge
    internal: true
```

## 儲存管理

### 資料卷配置

```yaml
volumes:
  # MySQL 資料持久化
  mysql-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./data/mysql

  # Redis 資料持久化
  redis-data:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./data/redis

  # 應用程式日誌
  app-logs:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./logs

  # 檔案上傳儲存
  uploads:
    driver: local
    driver_opts:
      type: none
      o: bind
      device: ./uploads
```

### 備份配置

```yaml
# 資料庫備份服務
mysql-backup:
  image: mysql:8.0
  depends_on:
    - mysql
  volumes:
    - ./backup:/backup
    - ./scripts/backup.sh:/backup.sh:ro
  command: /backup.sh
  networks:
    - go-admin-network
```

## 環境變數設定

### 開發環境

```bash
# .env.development
COMPOSE_PROJECT_NAME=go-admin-dev
APP_ENV=development
APP_DEBUG=true
DB_PASSWORD=dev_password
LOG_LEVEL=debug
```

### 測試環境

```bash
# .env.testing
COMPOSE_PROJECT_NAME=go-admin-test
APP_ENV=testing
APP_DEBUG=false
DB_PASSWORD=test_password
LOG_LEVEL=info
```

### 生產環境

```bash
# .env.production
COMPOSE_PROJECT_NAME=go-admin-prod
APP_ENV=production
APP_DEBUG=false
DB_PASSWORD=strong_production_password
JWT_SECRET=random_strong_secret_key
LOG_LEVEL=warn
```

## 常用指令

### 啟動服務

```bash
# 啟動所有服務
docker-compose up -d

# 啟動特定服務
docker-compose up -d go-admin mysql

# 重新建置並啟動
docker-compose up --build -d
```

### 停止服務

```bash
# 停止所有服務
docker-compose down

# 停止並清理資料卷
docker-compose down -v

# 停止並清理映像
docker-compose down --rmi all
```

### 查看狀態

```bash
# 查看服務狀態
docker-compose ps

# 查看服務日誌
docker-compose logs -f go-admin

# 查看資源使用情況
docker-compose top
```

### 管理服務

```bash
# 重新啟動服務
docker-compose restart go-admin

# 擴展服務副本
docker-compose up --scale go-admin=3 -d

# 進入容器
docker-compose exec go-admin bash
```

## 資料庫初始化

### 初始化腳本

```sql
-- scripts/init.sql
USE go-admin;

-- 建立初始管理員帳戶
INSERT INTO sys_users (username, password, role_id, status)
VALUES ('admin', '$2a$10$hashed_password', 1, 1);

-- 建立預設角色
INSERT INTO sys_roles (name, code, description)
VALUES ('管理員', 'admin', '系統管理員角色');

-- 建立預設選單
INSERT INTO sys_menus (name, path, component, icon, sort)
VALUES ('儀表板', '/dashboard', 'Dashboard', 'dashboard', 1);
```

### 遷移執行

```bash
# 在容器中執行資料庫遷移
docker-compose exec go-admin go run cmd/migrate/main.go

# 或使用 Makefile
docker-compose exec go-admin make migrate
```

## 監控與日誌

### 日誌收集

```yaml
# 增加 ELK 堆疊
elasticsearch:
  image: elasticsearch:7.17.0
  environment:
    - discovery.type=single-node

logstash:
  image: logstash:7.17.0
  volumes:
    - ./config/logstash:/usr/share/logstash/pipeline

kibana:
  image: kibana:7.17.0
  ports:
    - "5601:5601"
```

### 效能監控

```yaml
# 增加 Prometheus + Grafana
prometheus:
  image: prom/prometheus
  ports:
    - "9090:9090"
  volumes:
    - ./config/prometheus.yml:/etc/prometheus/prometheus.yml

grafana:
  image: grafana/grafana
  ports:
    - "3000:3000"
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
```

## 疑難排解

### 常見問題

1. **容器無法啟動**

   - 檢查連接埠是否衝突
   - 確認 Docker 服務運行正常
   - 查看容器日誌：`docker-compose logs [服務名稱]`

2. **資料庫連接失敗**

   - 確認資料庫服務已啟動
   - 檢查網路連通性
   - 驗證資料庫認證資訊

3. **效能問題**
   - 調整資源限制
   - 最佳化資料庫配置
   - 檢查磁碟空間

### 除錯工具

```bash
# 進入容器除錯
docker-compose exec go-admin sh

# 查看網路配置
docker network ls
docker network inspect go-admin_go-admin-network

# 查看資料卷
docker volume ls
docker volume inspect go-admin_mysql-data
```

## 最佳實踐

### 安全建議

1. **密碼管理**

   - 使用 Docker Secrets 管理敏感資訊
   - 定期更換預設密碼
   - 避免在配置檔案中明文儲存密碼

2. **網路安全**

   - 使用內部網路隔離服務
   - 僅暴露必要的連接埠
   - 設定適當的防火牆規則

3. **資料安全**
   - 定期備份資料
   - 加密敏感資料儲存
   - 設定適當的檔案權限

### 效能最佳化

1. **資源配置**

   - 根據負載調整記憶體限制
   - 合理分配 CPU 資源
   - 最佳化磁碟 I/O

2. **快取策略**
   - 配置 Redis 快取
   - 使用 CDN 加速靜態資源
   - 啟用 Gzip 壓縮

## 相關文件

- [Docker 部署指南](./deployment.md)
- [容器監控指南](./monitoring.md)
- [Go-Admin 架構說明](../development/architecture.md)
- [生產環境最佳實踐](../deployment/production.md)

## 版本記錄

| 版本  | 日期       | 更新內容               |
| ----- | ---------- | ---------------------- |
| 1.0.0 | 2025-09-19 | 初始版本，基本配置說明 |

---

_最後更新：2025 年 9 月 19 日_
