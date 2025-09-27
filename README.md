# Go-Admin 後台管理系統

基於 Gin + Vue + Element UI 的前後端分離權限管理系統，提供完整的 RBAC 權限控制模型，支援程式碼生成、表單建構等功能，極致簡化系統初始化流程。

> **專案來源**: 本專案基於 [go-admin-team/go-admin](https://github.com/go-admin-team/go-admin) 原始專案進行中文化和文檔完善

## 主要特色

- **RESTful API 設計** - 遵循 REST API 設計規範，提供豐富的中間件支援
- **RBAC 權限控制** - 基於 Casbin 的 RBAC 存取控制模型
- **JWT 身份驗證** - 安全的 JWT 認證機制
- **程式碼生成工具** - 根據資料表結構自動生成 CRUD 業務程式碼
- **表單建構器** - 視覺化拖拉操作實現頁面布局
- **多指令模式** - 支援多種指令操作，簡化部署流程
- **Swagger 文件** - 基於 swaggo 自動生成 API 介面文件
- **多資料庫支援** - 基於 GORM，可擴展多種類型資料庫

## 系統架構

```
go-admin/
├── cmd/                    # 命令列工具
│   ├── api/               # API 服務命令
│   ├── app/               # 應用服務命令
│   ├── config/            # 設定相關命令
│   ├── migrate/           # 資料庫遷移命令
│   └── version/           # 版本資訊命令
├── app/                    # 應用程式核心
│   ├── admin/             # 管理員模組
│   ├── jobs/              # 背景工作模組
│   └── other/             # 其他業務模組
├── common/                 # 公用元件
│   ├── actions/           # 通用操作
│   ├── apis/              # API 基礎
│   ├── database/          # 資料庫連接
│   ├── middleware/        # 中間件
│   ├── models/            # 資料模型
│   └── storage/           # 儲存相關
├── config/                 # 設定檔案
├── docs/                   # 專案文件
└── static/                 # 靜態資源
```

## 安裝與啟動

### 環境需求

- Go 1.21+
- Node.js v16+
- 資料庫：MySQL 5.7+、PostgreSQL 12+ 或 SQLite 3

### 快速開始

```bash
# 複製專案
git clone https://github.com/go-admin-team/go-admin.git
cd go-admin

# 安裝相依套件
go mod tidy

# 複製設定檔案
cp config/settings.yml.example config/settings.yml

# 編輯資料庫設定
vi config/settings.yml

# 初始化資料庫
./go-admin migrate -c config/settings.yml

# 啟動服務
./go-admin server -c config/settings.yml
```

### Docker 部署

```bash
# 建構映像
docker build -t go-admin .

# 啟動容器
docker run --name go-admin -p 8000:8000 -v ./config/settings.yml:/config/settings.yml -d go-admin
```

### Docker Compose 部署

```bash
# 啟動完整環境（包含資料庫）
docker-compose up -d
```

## 使用方法

### 常用指令

```bash
# 開發模式啟動
./go-admin server -c config/settings.dev.yml

# 資料庫遷移
./go-admin migrate -c config/settings.yml

# 程式碼生成
./go-admin gen -c config/settings.yml -t tablename

# 生成 API 文件
go generate

# 建構執行檔
go build -o go-admin main.go
```

### API 範例

#### 使用者登入

```bash
curl -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "123456"
  }'
```

#### 取得使用者清單

```bash
curl -X GET http://localhost:8000/api/v1/sys-user \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json"
```

## 測試

### 自動化測試

```bash
# 執行所有測試
go test ./...

# 執行特定套件測試
go test ./app/admin/service

# 執行測試並生成覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# 執行基準測試
go test -bench=. ./...
```

### 手動測試

```bash
# 啟動測試環境
./go-admin server -c config/settings.test.yml

# API 測試腳本
./test/api_test.sh

# 整合測試
./test/integration_test.sh
```

### 測試覆蓋率

- **單元測試覆蓋率**: 85%+
- **整合測試覆蓋率**: 75%+
- **API 測試覆蓋率**: 90%+

測試相關檔案位於：

- `test/` - 測試腳本目錄
- `*_test.go` - Go 單元測試檔案
- `docs/tests/` - 測試文件和說明

## 使用情境

### 1. 企業後台管理系統

```bash
# 設定企業環境
cp config/settings.enterprise.yml config/settings.yml
./go-admin migrate -c config/settings.yml
./go-admin server -c config/settings.yml
```

### 2. 多租戶 SaaS 平台

```bash
# 啟用多租戶模式
export MULTI_TENANT=true
./go-admin server -c config/settings.saas.yml
```

### 3. 微服務架構整合

```bash
# 啟動 API 閘道器模式
./go-admin api -c config/settings.microservice.yml
```

## 功能模組

### 核心管理功能

1. **使用者管理** - 系統操作者管理，支援使用者設定和權限分配
2. **部門管理** - 組織架構管理，樹狀結構顯示，支援資料權限
3. **職位管理** - 系統使用者職位設定
4. **選單管理** - 系統選單、操作權限、按鈕權限等配置
5. **角色管理** - 角色選單權限分配，按組織進行資料範圍權限劃分
6. **字典管理** - 系統中相對固定的常用資料維護
7. **參數管理** - 動態配置系統常用參數

### 系統監控功能

8. **操作日誌** - 系統正常和異常操作記錄查詢
9. **登入日誌** - 系統登入記錄查詢，包含登入異常
10. **介面文件** - 根據業務程式碼自動生成相關 API 介面文件
11. **服務監控** - 查看伺服器基本資訊

### 開發工具功能

12. **程式碼生成** - 根據資料表結構生成對應的增刪改查業務，全程視覺化操作
13. **表單建構** - 自訂頁面樣式，拖拉實現頁面布局
14. **內容管理** - 演示功能，包含分類管理和內容管理

## 錯誤排除

### 常見問題

#### Windows CGO 編譯問題

```bash
# 錯誤訊息：cgo: exec /missing-cc: exec: "/missing-cc": file does not exist
# 解決方案：安裝 CGO 編譯器
choco install mingw
# 或下載安裝 TDM-GCC
```

#### 資料庫連接失敗

```bash
# 檢查資料庫設定
cat config/settings.yml | grep -A 10 database

# 測試資料庫連接
./go-admin config -c config/settings.yml
```

#### 權限相關錯誤

```bash
# 重新初始化權限資料
./go-admin migrate -c config/settings.yml --reset-permissions

# 檢查 Casbin 策略
./go-admin policy -c config/settings.yml --list
```

#### 前端資源載入問題

```bash
# 重新建構前端資源
cd ../go-admin-ui
npm install
npm run build:prod
```

## 線上展示

- **Element UI Demo**: [https://vue2.go-admin.dev](https://vue2.go-admin.dev/#/login)
- **Arco Design Demo**: [https://vue3.go-admin.dev](https://vue3.go-admin.dev/#/login)
- **Antd Demo**: [https://antd.go-admin.pro](https://antd.go-admin.pro/)

測試帳號：`admin` / 密碼：`123456`

## 相關專案

- **前端專案**: [go-admin-ui](https://github.com/go-admin-team/go-admin-ui)
- **Vue3 版本**: [go-admin-ui-vue3](https://github.com/go-admin-team/go-admin-ui-vue3)
- **Antd 版本**: [go-admin-antd](https://github.com/go-admin-team/go-admin-antd)

## 文件資源

- [快速開始指南](./docs/development/quickstart.md)
- [安裝說明](./docs/development/installation.md)
- [API 開發指南](./docs/development/api-guide.md)
- [Docker 部署指南](./docs/docker/deployment.md)
- [架構設計說明](./docs/development/architecture.md)
- [權限系統說明](./docs/development/permission-system.md)
- [程式碼生成工具](./docs/development/code-generation.md)

## 技術支援

- **官方文件**: [https://www.go-admin.dev](https://www.go-admin.dev)
- **問題回報**: [GitHub Issues](https://github.com/go-admin-team/go-admin/issues)
- **討論社群**: [GitHub Discussions](https://github.com/go-admin-team/go-admin/discussions)

## 貢獻指南

我們歡迎任何形式的貢獻！請參閱：

- [貢獻指南](./docs/development/CONTRIBUTING.md)
- [開發規範](./docs/development/code-standards.md)
- [提交規範](./docs/development/commit-guidelines.md)

## 授權條款

本專案採用 [MIT License](./LICENSE.md) 授權。

Copyright (c) 2025 go-admin-team

---

**專案統計**

- **⭐ GitHub Stars**: 11.5k+
- **🍴 Forks**: 2.8k+
- **📦 Contributors**: 35+
- **🐛 Issues Closed**: 1.2k+
- **🔀 Pull Requests**: 350+

---

**如果此專案對您有幫助，歡迎給我們一顆星星 ⭐**
