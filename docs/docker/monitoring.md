# 容器監控指南

本文件詳細介紹如何監控 Go-Admin 容器化環境，包含效能監控、日誌管理、警報設定等方面。

## 概要

有效的監控系統是確保 Go-Admin 在生產環境穩定運行的關鍵。本指南涵蓋了從基礎監控到進階分析的完整解決方案。

## 監控架構

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Go-Admin  │───▶│ Prometheus  │───▶│   Grafana   │
│  Container  │    │   Metrics   │    │ Dashboard   │
└─────────────┘    └─────────────┘    └─────────────┘
       │                    │                 │
       ▼                    ▼                 ▼
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│    Logs     │───▶│ ELK Stack   │    │  AlertMgr   │
│  Container  │    │(E+L+K+F)    │    │ Notification│
└─────────────┘    └─────────────┘    └─────────────┘
```

## Docker Compose 監控配置

### 完整監控堆疊

```yaml
# docker-compose.monitoring.yml
version: "3.8"

services:
  # 主應用程式 (已存在)
  go-admin:
    build: .
    ports:
      - "8000:8000"
    labels:
      - "monitoring=enabled"
      - "service=go-admin"
    logging:
      driver: "fluentd"
      options:
        fluentd-address: localhost:24224
        tag: go-admin.logs
    networks:
      - monitoring

  # Prometheus - 指標收集
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./config/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml:ro
      - ./config/prometheus/rules/:/etc/prometheus/rules/:ro
      - prometheus-data:/prometheus
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--storage.tsdb.path=/prometheus"
      - "--web.console.libraries=/etc/prometheus/console_libraries"
      - "--web.console.templates=/etc/prometheus/consoles"
      - "--storage.tsdb.retention.time=200h"
      - "--web.enable-lifecycle"
    networks:
      - monitoring

  # Grafana - 視覺化儀表板
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_USER=admin
      - GF_SECURITY_ADMIN_PASSWORD=admin123
      - GF_USERS_ALLOW_SIGN_UP=false
    volumes:
      - grafana-data:/var/lib/grafana
      - ./config/grafana/provisioning/:/etc/grafana/provisioning/
      - ./config/grafana/dashboards/:/var/lib/grafana/dashboards/
    networks:
      - monitoring

  # Node Exporter - 系統指標
  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    restart: unless-stopped
    volumes:
      - /proc:/host/proc:ro
      - /sys:/host/sys:ro
      - /:/rootfs:ro
    command:
      - "--path.procfs=/host/proc"
      - "--path.rootfs=/rootfs"
      - "--path.sysfs=/host/sys"
      - "--collector.filesystem.mount-points-exclude=^/(sys|proc|dev|host|etc)($$|/)"
    ports:
      - "9100:9100"
    networks:
      - monitoring

  # cAdvisor - 容器指標
  cadvisor:
    image: gcr.io/cadvisor/cadvisor:latest
    container_name: cadvisor
    ports:
      - "8080:8080"
    volumes:
      - /:/rootfs:ro
      - /var/run:/var/run:rw
      - /sys:/sys:ro
      - /var/lib/docker/:/var/lib/docker:ro
      - /dev/disk/:/dev/disk:ro
    privileged: true
    devices:
      - /dev/kmsg
    networks:
      - monitoring

volumes:
  prometheus-data:
  grafana-data:

networks:
  monitoring:
    driver: bridge
```

## Prometheus 設定

### 基本配置

```yaml
# config/prometheus/prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "rules/*.yml"

scrape_configs:
  # Go-Admin 應用程式指標
  - job_name: "go-admin"
    static_configs:
      - targets: ["go-admin:8000"]
    metrics_path: "/metrics"
    scrape_interval: 10s

  # Node Exporter - 系統指標
  - job_name: "node-exporter"
    static_configs:
      - targets: ["node-exporter:9100"]

  # cAdvisor - 容器指標
  - job_name: "cadvisor"
    static_configs:
      - targets: ["cadvisor:8080"]

  # MySQL Exporter
  - job_name: "mysql-exporter"
    static_configs:
      - targets: ["mysql-exporter:9104"]

  # Redis Exporter
  - job_name: "redis-exporter"
    static_configs:
      - targets: ["redis-exporter:9121"]

alerting:
  alertmanagers:
    - static_configs:
        - targets:
            - alertmanager:9093
```

### 警報規則

```yaml
# config/prometheus/rules/go-admin.yml
groups:
  - name: go-admin.rules
    rules:
      # 應用程式可用性
      - alert: GoAdminDown
        expr: up{job="go-admin"} == 0
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Go-Admin 服務無法存取"
          description: "Go-Admin 服務已停止回應超過 1 分鐘"

      # 高回應時間
      - alert: HighResponseTime
        expr: http_request_duration_seconds{quantile="0.95"} > 2
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Go-Admin 回應時間過高"
          description: "95% 請求回應時間超過 2 秒"

      # 高錯誤率
      - alert: HighErrorRate
        expr: rate(http_requests_total{status=~"5.."}[5m]) > 0.1
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Go-Admin 錯誤率過高"
          description: "5xx 錯誤率超過 10%"

      # 記憶體使用率
      - alert: HighMemoryUsage
        expr: (container_memory_usage_bytes{name="go-admin"} / container_spec_memory_limit_bytes{name="go-admin"}) > 0.8
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Go-Admin 記憶體使用率過高"
          description: "記憶體使用率超過 80%"

      # CPU 使用率
      - alert: HighCPUUsage
        expr: (rate(container_cpu_usage_seconds_total{name="go-admin"}[5m]) * 100) > 80
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "Go-Admin CPU 使用率過高"
          description: "CPU 使用率超過 80%"
```

## Grafana 儀表板

### Go-Admin 應用程式儀表板

```json
{
	"dashboard": {
		"title": "Go-Admin 應用程式監控",
		"panels": [
			{
				"title": "請求總數",
				"type": "stat",
				"targets": [
					{
						"expr": "sum(rate(http_requests_total{job=\"go-admin\"}[5m]))",
						"legendFormat": "RPS"
					}
				]
			},
			{
				"title": "回應時間分布",
				"type": "graph",
				"targets": [
					{
						"expr": "histogram_quantile(0.50, rate(http_request_duration_seconds_bucket{job=\"go-admin\"}[5m]))",
						"legendFormat": "50th percentile"
					},
					{
						"expr": "histogram_quantile(0.95, rate(http_request_duration_seconds_bucket{job=\"go-admin\"}[5m]))",
						"legendFormat": "95th percentile"
					},
					{
						"expr": "histogram_quantile(0.99, rate(http_request_duration_seconds_bucket{job=\"go-admin\"}[5m]))",
						"legendFormat": "99th percentile"
					}
				]
			},
			{
				"title": "HTTP 狀態碼分布",
				"type": "piechart",
				"targets": [
					{
						"expr": "sum by (status) (rate(http_requests_total{job=\"go-admin\"}[5m]))",
						"legendFormat": "{{status}}"
					}
				]
			}
		]
	}
}
```

### 系統資源儀表板

```json
{
	"dashboard": {
		"title": "Go-Admin 系統資源",
		"panels": [
			{
				"title": "CPU 使用率",
				"type": "graph",
				"targets": [
					{
						"expr": "rate(container_cpu_usage_seconds_total{name=\"go-admin\"}[5m]) * 100",
						"legendFormat": "CPU Usage %"
					}
				]
			},
			{
				"title": "記憶體使用情況",
				"type": "graph",
				"targets": [
					{
						"expr": "container_memory_usage_bytes{name=\"go-admin\"} / 1024 / 1024",
						"legendFormat": "Memory Usage (MB)"
					},
					{
						"expr": "container_spec_memory_limit_bytes{name=\"go-admin\"} / 1024 / 1024",
						"legendFormat": "Memory Limit (MB)"
					}
				]
			},
			{
				"title": "網路 I/O",
				"type": "graph",
				"targets": [
					{
						"expr": "rate(container_network_receive_bytes_total{name=\"go-admin\"}[5m])",
						"legendFormat": "Received"
					},
					{
						"expr": "rate(container_network_transmit_bytes_total{name=\"go-admin\"}[5m])",
						"legendFormat": "Transmitted"
					}
				]
			},
			{
				"title": "磁碟 I/O",
				"type": "graph",
				"targets": [
					{
						"expr": "rate(container_fs_reads_bytes_total{name=\"go-admin\"}[5m])",
						"legendFormat": "Read"
					},
					{
						"expr": "rate(container_fs_writes_bytes_total{name=\"go-admin\"}[5m])",
						"legendFormat": "Write"
					}
				]
			}
		]
	}
}
```

## 日誌管理

### ELK Stack 配置

```yaml
# docker-compose.logging.yml
version: "3.8"

services:
  # Elasticsearch - 日誌儲存
  elasticsearch:
    image: elasticsearch:7.17.0
    container_name: elasticsearch
    environment:
      - discovery.type=single-node
      - cluster.name=go-admin-logs
      - node.name=es01
      - bootstrap.memory_lock=true
      - "ES_JAVA_OPTS=-Xms1g -Xmx1g"
    ulimits:
      memlock:
        soft: -1
        hard: -1
    volumes:
      - elasticsearch-data:/usr/share/elasticsearch/data
    ports:
      - "9200:9200"
    networks:
      - logging

  # Logstash - 日誌處理
  logstash:
    image: logstash:7.17.0
    container_name: logstash
    volumes:
      - ./config/logstash/logstash.conf:/usr/share/logstash/pipeline/logstash.conf:ro
    ports:
      - "5000:5000"
      - "9600:9600"
    environment:
      LS_JAVA_OPTS: "-Xmx256m -Xms256m"
    networks:
      - logging
    depends_on:
      - elasticsearch

  # Kibana - 日誌視覺化
  kibana:
    image: kibana:7.17.0
    container_name: kibana
    ports:
      - "5601:5601"
    environment:
      ELASTICSEARCH_URL: http://elasticsearch:9200
      ELASTICSEARCH_HOSTS: '["http://elasticsearch:9200"]'
    networks:
      - logging
    depends_on:
      - elasticsearch

  # Fluentd - 日誌收集
  fluentd:
    build: ./fluentd
    container_name: fluentd
    volumes:
      - ./config/fluentd/fluent.conf:/fluentd/etc/fluent.conf:ro
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
    ports:
      - "24224:24224"
      - "24224:24224/udp"
    networks:
      - logging
    depends_on:
      - elasticsearch

volumes:
  elasticsearch-data:

networks:
  logging:
    driver: bridge
```

### Logstash 配置

```ruby
# config/logstash/logstash.conf
input {
  beats {
    port => 5044
  }

  tcp {
    port => 5000
    codec => json
  }
}

filter {
  if [fields][service] == "go-admin" {
    # 解析 Go-Admin 日誌格式
    grok {
      match => {
        "message" => "%{TIMESTAMP_ISO8601:timestamp} %{LOGLEVEL:level} %{GREEDYDATA:msg}"
      }
    }

    # 解析時間戳記
    date {
      match => [ "timestamp", "yyyy-MM-dd HH:mm:ss.SSS" ]
    }

    # 添加標籤
    mutate {
      add_tag => [ "go-admin", "application" ]
    }
  }

  # 處理 HTTP 存取日誌
  if [fields][logtype] == "access" {
    grok {
      match => {
        "message" => "%{COMBINEDAPACHELOG}"
      }
    }
  }
}

output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "go-admin-%{+YYYY.MM.dd}"
  }

  stdout {
    codec => rubydebug
  }
}
```

### Fluentd 配置

```xml
# config/fluentd/fluent.conf
<source>
  @type forward
  port 24224
  bind 0.0.0.0
</source>

<filter go-admin.**>
  @type parser
  format json
  key_name log
  reserve_data true
</filter>

<match go-admin.**>
  @type elasticsearch
  host elasticsearch
  port 9200
  index_name go-admin
  type_name _doc

  <buffer>
    @type file
    path /var/log/fluentd-buffers/go-admin.buffer
    flush_mode interval
    retry_type exponential_backoff
    flush_thread_count 2
    flush_interval 5s
    retry_forever
    retry_max_interval 30
    chunk_limit_size 2M
    queue_limit_length 8
    overflow_action block
  </buffer>
</match>
```

## 應用程式指標暴露

### Go 應用程式整合

```go
// pkg/metrics/metrics.go
package metrics

import (
    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "strconv"
    "time"
)

var (
    // HTTP 請求總數
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )

    // HTTP 請求持續時間
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "http_request_duration_seconds",
            Help: "Duration of HTTP requests in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )

    // 活躍連接數
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )

    // 資料庫查詢指標
    dbQueriesTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "db_queries_total",
            Help: "Total number of database queries",
        },
        []string{"operation", "table"},
    )
)

// PrometheusMiddleware 監控中間件
func PrometheusMiddleware() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()

        c.Next()

        duration := time.Since(start)
        status := strconv.Itoa(c.Writer.Status())

        httpRequestsTotal.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            status,
        ).Inc()

        httpRequestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration.Seconds())
    })
}

// MetricsHandler 指標端點處理器
func MetricsHandler() gin.HandlerFunc {
    h := promhttp.Handler()
    return func(c *gin.Context) {
        h.ServeHTTP(c.Writer, c.Request)
    }
}
```

### 在 main.go 中整合

```go
// main.go
import (
    "your-project/pkg/metrics"
)

func main() {
    r := gin.Default()

    // 添加監控中間件
    r.Use(metrics.PrometheusMiddleware())

    // 暴露指標端點
    r.GET("/metrics", metrics.MetricsHandler())

    // 其他路由...

    r.Run(":8000")
}
```

## 警報管理

### AlertManager 配置

```yaml
# config/alertmanager/alertmanager.yml
global:
  smtp_smarthost: "localhost:587"
  smtp_from: "alerts@yourdomain.com"

route:
  group_by: ["alertname"]
  group_wait: 10s
  group_interval: 10s
  repeat_interval: 1h
  receiver: "web.hook"

receivers:
  - name: "web.hook"
    email_configs:
      - to: "admin@yourdomain.com"
        subject: "Go-Admin Alert: {{ .GroupLabels.alertname }}"
        body: |
          {{ range .Alerts }}
          Alert: {{ .Annotations.summary }}
          Description: {{ .Annotations.description }}
          Labels:
          {{ range .Labels.SortedPairs }} - {{ .Name }} = {{ .Value }}
          {{ end }}
          {{ end }}

    slack_configs:
      - api_url: "YOUR_SLACK_WEBHOOK_URL"
        channel: "#alerts"
        title: "Go-Admin Alert"
        text: "{{ range .Alerts }}{{ .Annotations.description }}{{ end }}"

inhibit_rules:
  - source_match:
      severity: "critical"
    target_match:
      severity: "warning"
    equal: ["alertname", "dev", "instance"]
```

## 健康檢查

### 應用程式健康檢查

```go
// pkg/health/health.go
package health

import (
    "database/sql"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

type HealthChecker struct {
    db    *sql.DB
    redis RedisClient
}

type HealthStatus struct {
    Status    string            `json:"status"`
    Timestamp time.Time         `json:"timestamp"`
    Services  map[string]string `json:"services"`
}

func (h *HealthChecker) HealthHandler(c *gin.Context) {
    status := &HealthStatus{
        Timestamp: time.Now(),
        Services:  make(map[string]string),
    }

    // 檢查資料庫連接
    if err := h.db.Ping(); err != nil {
        status.Services["database"] = "down"
        status.Status = "unhealthy"
    } else {
        status.Services["database"] = "up"
    }

    // 檢查 Redis 連接
    if err := h.redis.Ping(); err != nil {
        status.Services["redis"] = "down"
        status.Status = "unhealthy"
    } else {
        status.Services["redis"] = "up"
    }

    // 設定整體狀態
    if status.Status == "" {
        status.Status = "healthy"
    }

    if status.Status == "healthy" {
        c.JSON(http.StatusOK, status)
    } else {
        c.JSON(http.StatusServiceUnavailable, status)
    }
}
```

## 效能調優監控

### 關鍵效能指標 (KPIs)

1. **應用程式效能**

   - 回應時間：P50, P95, P99
   - 吞吐量：每秒請求數 (RPS)
   - 錯誤率：4xx/5xx 錯誤百分比

2. **系統資源**

   - CPU 使用率
   - 記憶體使用率
   - 磁碟 I/O
   - 網路頻寬

3. **資料庫效能**

   - 查詢回應時間
   - 連接池使用率
   - 慢查詢數量

4. **快取效能**
   - 快取命中率
   - 快取大小
   - 過期鍵數量

### 效能基準線設定

```yaml
# 效能基準線 (SLA)
response_time:
  p50: < 200ms
  p95: < 500ms
  p99: < 1000ms

throughput:
  min_rps: 100
  target_rps: 500

error_rate:
  max_4xx: 5%
  max_5xx: 1%

resource_usage:
  max_cpu: 70%
  max_memory: 80%
  max_disk_io: 80%
```

## 監控最佳實踐

### 1. 監控分層

```
Business Metrics (業務指標)
    ↓
Application Metrics (應用指標)
    ↓
Infrastructure Metrics (基礎設施指標)
```

### 2. 警報策略

- **關鍵警報**：立即通知 (SMS/電話)
- **重要警報**：5 分鐘內通知 (Email/Slack)
- **一般警報**：每小時彙總通知

### 3. 儀表板設計

- **高階概覽**：業務 KPIs
- **運維儀表板**：系統健康狀況
- **除錯儀表板**：詳細指標

### 4. 資料保留政策

```yaml
retention:
  high_resolution: 7d # 高解析度資料
  medium_resolution: 30d # 中等解析度資料
  low_resolution: 1y # 低解析度資料
```

## 常見問題排解

### 監控系統問題

1. **Prometheus 無法抓取指標**

   - 檢查目標服務是否正常運行
   - 驗證網路連通性
   - 確認 `/metrics` 端點可存取

2. **Grafana 儀表板顯示異常**

   - 檢查資料來源配置
   - 驗證查詢語法
   - 確認時間範圍設定

3. **日誌收集中斷**
   - 檢查 Fluentd 容器狀態
   - 驗證 Elasticsearch 連接
   - 確認日誌格式正確

### 除錯指令

```bash
# 檢查監控服務狀態
docker-compose -f docker-compose.monitoring.yml ps

# 查看 Prometheus 目標狀態
curl http://localhost:9090/api/v1/targets

# 檢查 Elasticsearch 健康狀況
curl http://localhost:9200/_cluster/health

# 查看容器資源使用情況
docker stats
```

## 安全考量

### 1. 存取控制

- 設定 Grafana 使用者權限
- 限制 Prometheus 查詢權限
- 保護敏感指標資料

### 2. 資料保護

- 加密傳輸資料
- 脫敏日誌中的敏感資訊
- 定期輪換認證憑證

### 3. 網路隔離

- 使用內部網路隔離監控服務
- 設定防火牆規則
- 啟用 TLS 加密

## 相關文件

- [Docker Compose 配置說明](./docker-compose.md)
- [Docker 部署指南](./deployment.md)
- [Go-Admin 架構說明](../development/architecture.md)
- [生產環境最佳實踐](../deployment/production.md)

## 版本記錄

| 版本  | 日期       | 更新內容               |
| ----- | ---------- | ---------------------- |
| 1.0.0 | 2025-09-19 | 初始版本，完整監控指南 |

---

_最後更新：2025 年 9 月 19 日_
