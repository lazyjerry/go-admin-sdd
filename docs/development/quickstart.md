# Go-Admin 快速開始指南

本指南將帶您快速體驗 Go-Admin 後台管理系統的核心功能，從零開始到擁有完整的管理後台。

## 5 分鐘快速體驗

### 第一步：環境確認

確保您的環境滿足基本需求：

```bash
# 檢查 Go 版本
go version
# 應輸出：go version go1.21.x 或更高

# 檢查 Git
git --version
```

### 第二步：下載並啟動

```bash
# 下載專案
git clone https://github.com/go-admin-team/go-admin.git
cd go-admin

# 安裝依賴
go mod tidy

# 快速啟動 (使用 SQLite)
./go-admin migrate -c config/settings.sqlite.yml
./go-admin server -c config/settings.sqlite.yml
```

### 第三步：訪問後台

1. 開啟瀏覽器，訪問：[http://localhost:8000](http://localhost:8000)
2. 使用預設帳號登入：
   - **使用者名稱**: `admin`
   - **密碼**: `123456`

🎉 **恭喜！** 您已成功啟動 Go-Admin 系統。

---

## 詳細入門指南

### 專案結構導覽

下載完成後，您會看到以下專案結構：

```
go-admin/
├── cmd/                    # 命令列工具
│   ├── api/               # API 服務命令
│   ├── migrate/           # 資料庫遷移命令
│   └── version/           # 版本命令
├── app/                    # 業務邏輯層
│   ├── admin/             # 管理員功能模組
│   │   ├── apis/          # API 處理器
│   │   ├── models/        # 資料模型
│   │   ├── router/        # 路由定義
│   │   └── service/       # 業務服務
│   ├── jobs/              # 背景任務模組
│   └── other/             # 其他業務模組
├── common/                 # 公用元件
│   ├── actions/           # 通用 CRUD 操作
│   ├── database/          # 資料庫連接
│   ├── middleware/        # HTTP 中間件
│   ├── models/            # 基礎資料模型
│   └── storage/           # 儲存相關
├── config/                 # 設定檔案
│   ├── settings.yml       # 主要設定檔
│   ├── settings.dev.yml   # 開發環境設定
│   └── settings.sqlite.yml # SQLite 設定
└── docs/                  # 專案文件
```

### 設定檔案說明

#### 基本設定結構

Go-Admin 使用 YAML 格式的設定檔案，主要包含以下部分：

```yaml
# config/settings.yml
settings:
  application:
    mode: dev # 運行模式: dev, test, prod
    name: go-admin # 應用程式名稱
    port: 8000 # 服務埠號
    host: 0.0.0.0 # 綁定主機

  database:
    dbtype: sqlite3 # 資料庫類型
    source: go-admin.db # 連接字串或檔案路徑

  jwt:
    secret: go-admin-secret # JWT 密鑰
    timeout: 24 # 權杖過期時間 (小時)

  log:
    level: info # 日誌級別
    path: storage/logs # 日誌路徑
```

### 資料庫設定選項

#### SQLite (推薦新手使用)

```yaml
settings:
  database:
    dbtype: sqlite3
    source: go-admin.db
```

#### MySQL

```yaml
settings:
  database:
    dbtype: mysql
    host: 127.0.0.1
    port: 3306
    username: root
    password: your_password
    dbname: go_admin
    config: charset=utf8mb4&parseTime=True&loc=Local
```

#### PostgreSQL

```yaml
settings:
  database:
    dbtype: postgres
    host: 127.0.0.1
    port: 5432
    username: postgres
    password: your_password
    dbname: go_admin
```

## 核心功能體驗

### 1. 使用者管理

登入後台後，您可以體驗以下功能：

#### 查看使用者列表

```bash
# API 方式查看
curl -X GET "http://localhost:8000/api/v1/sys-user" \
  -H "Authorization: Bearer YOUR_TOKEN"
```

#### 新增使用者

1. 進入「系統管理」→「使用者管理」
2. 點選「新增」按鈕
3. 填寫使用者資訊：
   - 使用者名稱
   - 真實姓名
   - 電子郵件
   - 手機號碼
   - 所屬部門
   - 角色分配

#### 編輯使用者權限

1. 在使用者列表中點選「編輯」
2. 調整角色權限
3. 儲存變更

### 2. 角色權限管理

#### 建立新角色

```bash
# 進入角色管理頁面
# 系統管理 → 角色管理 → 新增角色
```

設定角色資訊：

- **角色名稱**: 如「部門主管」
- **角色代碼**: 如「dept_manager」
- **角色描述**: 角色職責說明
- **選單權限**: 勾選可存取的選單
- **資料權限**: 設定資料存取範圍

#### 權限範圍說明

- **全部資料權限**: 可存取所有資料
- **自訂資料權限**: 指定部門資料
- **部門資料權限**: 本部門及子部門資料
- **本部門資料**: 僅本部門資料
- **僅本人資料**: 只能存取自己的資料

### 3. 選單管理

#### 新增選單項目

1. 進入「系統管理」→「選單管理」
2. 點選「新增」建立新選單
3. 設定選單屬性：

```yaml
# 選單配置範例
選單名稱: "產品管理"
選單圖示: "product"
排序: 100
路由地址: "/product"
元件路徑: "views/product/index"
選單類型: 目錄/選單/按鈕
顯示狀態: 顯示/隱藏
```

#### 按鈕權限設定

為選單項目新增操作按鈕：

- **新增**: `product:add`
- **編輯**: `product:edit`
- **刪除**: `product:delete`
- **查看**: `product:view`
- **匯出**: `product:export`

### 4. 部門管理

#### 建立組織架構

```
公司總部
├── 技術部
│   ├── 前端組
│   ├── 後端組
│   └── 測試組
├── 營運部
│   ├── 產品組
│   └── 營銷組
└── 行政部
    ├── 人資組
    └── 財務組
```

建立步驟：

1. 新增根部門：「公司總部」
2. 為總部新增子部門：「技術部」、「營運部」、「行政部」
3. 繼續新增更細層的部門結構

## 進階功能探索

### 1. 程式碼生成器

Go-Admin 提供強大的程式碼生成功能，可以根據資料庫表結構自動生成 CRUD 程式碼。

#### 建立測試資料表

```sql
-- 建立產品表
CREATE TABLE `products` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL COMMENT '產品名稱',
  `description` text COMMENT '產品描述',
  `price` decimal(10,2) NOT NULL COMMENT '價格',
  `category_id` int(11) COMMENT '分類ID',
  `status` tinyint(1) DEFAULT 1 COMMENT '狀態: 1-啟用 0-停用',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='產品表';
```

#### 使用程式碼生成器

```bash
# 生成產品管理模組
./go-admin gen -c config/settings.yml -t products

# 生成的檔案包括:
# - app/admin/apis/products.go (API 處理器)
# - app/admin/models/products.go (資料模型)
# - app/admin/router/products.go (路由定義)
# - app/admin/service/dto/products.go (DTO 物件)
```

#### 自訂生成模板

您可以修改 `template/` 目錄下的模板檔案來自訂生成的程式碼格式。

### 2. 表單建構器

使用視覺化表單建構器快速建立複雜表單：

#### 存取表單建構器

1. 進入「系統工具」→「表單建構」
2. 拖拉元件建立表單：
   - 基礎元件：輸入框、下拉選單、日期選擇器
   - 高級元件：檔案上傳、富文字編輯器、地址選擇器
   - 布局元件：分欄、標籤頁、摺疊面板

#### 表單配置

```json
{
	"formRef": "elForm",
	"formModel": "formData",
	"size": "medium",
	"labelPosition": "right",
	"labelWidth": 100,
	"formRules": "rules",
	"gutter": 15,
	"disabled": false,
	"span": 24,
	"formBtns": true
}
```

#### 匯出表單程式碼

建立完成後可以匯出：

- **Vue 元件程式碼**
- **JSON 配置**
- **HTML 程式碼**

### 3. API 介面測試

#### Swagger 文件

Go-Admin 自動生成 Swagger API 文件：

```bash
# 生成 API 文件
go generate

# 存取 Swagger UI
# http://localhost:8000/swagger/index.html
```

#### 常用 API 端點

```bash
# 使用者登入
POST /api/v1/login
{
  "username": "admin",
  "password": "123456"
}

# 取得使用者清單
GET /api/v1/sys-user

# 取得選單樹
GET /api/v1/menuTreeselect

# 取得部門樹
GET /api/v1/deptTree

# 上傳檔案
POST /api/v1/public/uploadFile
```

## 開發工作流程

### 1. 新功能開發流程

以開發「庫存管理」功能為例：

#### 第一步：設計資料結構

```sql
CREATE TABLE `inventory` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `product_id` int(11) NOT NULL COMMENT '產品ID',
  `warehouse_id` int(11) NOT NULL COMMENT '倉庫ID',
  `quantity` int(11) NOT NULL DEFAULT 0 COMMENT '庫存數量',
  `reserved_qty` int(11) NOT NULL DEFAULT 0 COMMENT '預留數量',
  `available_qty` int(11) GENERATED ALWAYS AS (quantity - reserved_qty) COMMENT '可用數量',
  `last_updated` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_product_warehouse` (`product_id`, `warehouse_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='庫存表';
```

#### 第二步：生成基礎程式碼

```bash
./go-admin gen -c config/settings.yml -t inventory
```

#### 第三步：自訂業務邏輯

編輯生成的檔案，新增業務邏輯：

```go
// app/admin/service/dto/inventory.go
type InventoryGetPageReq struct {
    dto.Pagination `search:"-"`
    ProductId      int    `form:"productId" search:"type:exact;column:product_id;table:inventory"`
    WarehouseId    int    `form:"warehouseId" search:"type:exact;column:warehouse_id;table:inventory"`
    dto.ControlBy  `search:"-"`
}

// 新增業務方法
func (s *InventoryService) UpdateStock(productId, warehouseId, quantity int) error {
    // 更新庫存邏輯
    return nil
}
```

#### 第四步：設定路由和權限

```go
// app/admin/router/inventory.go
func init() {
    routerNoCheckRole = append(routerNoCheckRole, registerInventoryRouter)
}

func registerInventoryRouter(v1 *gin.RouterGroup) {
    api := apis.Inventory{}
    r := v1.Group("/inventory").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
    {
        r.GET("", api.GetPage)
        r.POST("", api.Insert)
        r.PUT("/:id", api.Update)
        r.DELETE("", api.Delete)
        r.POST("/update-stock", api.UpdateStock) // 自訂 API
    }
}
```

#### 第五步：建立前端選單

在選單管理中新增：

- 父選單：「庫存管理」
- 子選單：「庫存清單」、「庫存調整」、「庫存報表」

### 2. 除錯和測試

#### 開啟除錯模式

```yaml
# config/settings.dev.yml
settings:
  application:
    mode: dev
  log:
    level: debug
```

#### 使用日誌除錯

```go
import "go-admin/common/log"

// 在程式碼中新增日誌
log.Info("開始處理庫存更新", "productId", productId)
log.Debug("查詢結果", "inventory", inventory)
log.Error("更新失敗", "error", err)
```

#### API 測試

```bash
# 測試庫存查詢 API
curl -X GET "http://localhost:8000/api/v1/inventory?productId=1" \
  -H "Authorization: Bearer YOUR_TOKEN"

# 測試庫存更新 API
curl -X POST "http://localhost:8000/api/v1/inventory/update-stock" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "productId": 1,
    "warehouseId": 1,
    "quantity": 100
  }'
```

## 部署和發布

### 1. 編譯生產版本

```bash
# 編譯 Linux 版本
env GOOS=linux GOARCH=amd64 go build -o go-admin-linux main.go

# 編譯 Windows 版本
env GOOS=windows GOARCH=amd64 go build -o go-admin-windows.exe main.go

# 本地編譯
go build -o go-admin main.go
```

### 2. 生產環境設定

```yaml
# config/settings.prod.yml
settings:
  application:
    mode: prod
    port: 8000
  database:
    dbtype: mysql
    host: your-db-host
    port: 3306
    username: your-username
    password: your-password
    dbname: go_admin_prod
  jwt:
    secret: your-production-secret-key
  log:
    level: info
    path: /var/log/go-admin
```

### 3. 使用 systemd 管理服務

建立服務檔案 `/etc/systemd/system/go-admin.service`：

```ini
[Unit]
Description=Go-Admin Service
After=network.target

[Service]
Type=simple
User=go-admin
WorkingDirectory=/opt/go-admin
ExecStart=/opt/go-admin/go-admin server -c /opt/go-admin/config/settings.prod.yml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

啟動服務：

```bash
sudo systemctl daemon-reload
sudo systemctl enable go-admin
sudo systemctl start go-admin
sudo systemctl status go-admin
```

## 效能監控

### 1. 內建監控端點

```bash
# 健康檢查
curl http://localhost:8000/api/v1/health

# 系統資訊
curl http://localhost:8000/api/v1/sys/server

# 記憶體使用情況
curl http://localhost:8000/debug/pprof/heap
```

### 2. 資料庫效能監控

```yaml
# 開啟 SQL 查詢日誌
settings:
  database:
    log-mode: info
    log-zap: true
    max-idle-conns: 10
    max-open-conns: 100
    conn-max-lifetime: 3600
```

### 3. 日誌分析

```bash
# 檢視錯誤日誌
tail -f storage/logs/go-admin.log | grep ERROR

# 分析慢查詢
grep "slow query" storage/logs/go-admin.log

# 統計 API 存取量
grep "GET\|POST\|PUT\|DELETE" storage/logs/go-admin.log | cut -d' ' -f4 | sort | uniq -c
```

## 下一步

完成快速開始後，建議您：

1. **閱讀詳細文件**：

   - [架構設計說明](./architecture.md)
   - [API 開發指南](./api-guide.md)
   - [權限系統說明](./permission-system.md)

2. **體驗進階功能**：

   - [程式碼生成工具](./code-generation.md)
   - [Docker 部署](../docker/deployment.md)
   - [生產環境配置](../deployment/production.md)

3. **參與社群**：

   - [GitHub Issues](https://github.com/go-admin-team/go-admin/issues)
   - [討論區](https://github.com/go-admin-team/go-admin/discussions)
   - [貢獻指南](./CONTRIBUTING.md)

4. **學習資源**：
   - [官方文件](https://www.go-admin.dev)
   - [影片教學](https://space.bilibili.com/565616721)
   - [範例專案](https://github.com/go-admin-team/go-admin-examples)

---

如果在使用過程中遇到任何問題，請參考 [常見問題解答](../faq.md) 或在 GitHub 上提交 Issue。
