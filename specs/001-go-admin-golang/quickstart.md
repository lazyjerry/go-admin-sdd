# Go-Admin 文檔化快速開始指南

## 概述

本指南將協助您快速體驗 go-admin 專案文檔化的成果，包含完整的繁體中文註解、詳細的架設說明以及二次開發教學。

## 前置需求

### 環境需求

- Go 1.24 或更新版本
- Git 版本控制系統
- 文字編輯器或 IDE (建議 VS Code)
- Docker (可選，用於容器化部署)

### 檢查環境

```bash
# 檢查 Go 版本
go version

# 檢查 Git 版本
git version

# 檢查 Docker 版本 (可選)
docker version
```

## 快速體驗流程

### 步驟 1: 克隆專案

```bash
# 克隆包含完整文檔的專案
git clone <repository-url>
cd go-admin-sdd

# 檢查專案結構
ls -la
```

**預期結果**: 看到以下目錄結構

```
├── README.md              # 完整專案說明 (繁體中文)
├── docs/                  # 文檔目錄
│   ├── docker/            # Docker 結構說明
│   ├── development/       # 二次開發教學
│   └── changes/           # 變更紀錄
└── go-admin/              # 主要專案 (含中文註解)
```

### 步驟 2: 閱讀專案文檔

```bash
# 查看主要 README
cat README.md

# 瀏覽 Docker 文檔
ls docs/docker/
cat docs/docker/README.md

# 查看開發指南
ls docs/development/
cat docs/development/README.md
```

**預期結果**: 所有文檔均使用繁體中文撰寫，內容完整且易於理解

### 步驟 3: 體驗中文註解

```bash
# 查看主要入口檔案
cd go-admin
cat main.go

# 查看核心模組註解
cat common/global/adm.go

# 查看 API 路由註解
cat app/admin/router/sys_menu.go
```

**預期結果**:

- 所有函數都有繁體中文註解說明用途
- 複雜邏輯包含流程說明
- 中英文之間有適當空格分隔

### 步驟 4: 本地環境設置

```bash
# 安裝相依套件
go mod tidy

# 檢查設定檔案
ls config/
cat config/settings.yml

# 初始化資料庫
go run main.go migrate -c=config/settings.yml
```

**預期結果**:

- 相依套件成功安裝
- 設定檔案有完整的中文說明註解
- 資料庫初始化成功

### 步驟 5: 啟動服務

```bash
# 啟動開發伺服器
go run main.go server -c=config/settings.yml

# 或使用 Docker (另開終端機)
docker-compose up -d
```

**預期結果**:

- 服務成功啟動在指定埠口
- 控制台顯示清晰的中文啟動日誌
- API 文檔可在瀏覽器中查看

### 步驟 6: 驗證 API 文檔

開啟瀏覽器訪問：

```
http://localhost:8000/swagger/admin/index.html
```

**預期結果**:

- Swagger UI 顯示完整的 API 文檔
- API 說明使用繁體中文
- 包含詳細的請求/回應範例

## 功能驗證清單

### ✅ 文檔完整性檢查

- [ ] README.md 包含完整專案說明
- [ ] 安裝指南清晰易懂
- [ ] Docker 部署說明完整
- [ ] 二次開發教學詳細
- [ ] API 文檔完整且最新

### ✅ 中文化品質檢查

- [ ] 所有文檔使用繁體中文
- [ ] 使用台灣慣用詞彙
- [ ] 中英文空格分隔正確
- [ ] 專業術語翻譯準確
- [ ] 語法通順易讀

### ✅ 程式碼註解檢查

- [ ] 主要函數有中文註解
- [ ] 複雜邏輯有流程說明
- [ ] API 端點說明完整
- [ ] 設定檔案註解詳細
- [ ] 資料模型欄位說明清楚

### ✅ 功能性檢查

- [ ] 專案可正常編譯
- [ ] 服務可正常啟動
- [ ] API 端點正常回應
- [ ] 資料庫操作正常
- [ ] 權限系統運作正常

## 常見問題排除

### Q1: 編譯錯誤

**症狀**: `go build` 或 `go run` 失敗
**解決**:

```bash
# 檢查 Go 版本
go version

# 清理模組快取
go clean -modcache
go mod download
go mod tidy
```

### Q2: 資料庫連線錯誤

**症狀**: 服務啟動時資料庫連線失敗
**解決**:

```bash
# 檢查設定檔案
cat config/settings.yml | grep -A 10 database

# 使用 SQLite (預設)
go run main.go migrate -c=config/settings.yml
```

### Q3: 埠口被占用

**症狀**: 服務啟動失敗，提示埠口被占用
**解決**:

```bash
# 檢查埠口使用情況
lsof -i :8000

# 修改設定檔案中的埠口號
vim config/settings.yml
```

### Q4: Docker 容器啟動失敗

**症狀**: `docker-compose up` 失敗
**解決**:

```bash
# 檢查 Docker 狀態
docker ps -a

# 重建映像檔
docker-compose build --no-cache
docker-compose up -d
```

## 下一步行動

### 🎯 深入學習

1. **閱讀架構文檔**: `docs/development/architecture.md`
2. **學習 API 開發**: `docs/development/api-guide.md`
3. **了解權限系統**: 查看 `common/middleware/permission.go` 中的註解

### 🛠️ 開始開發

1. **建立新模組**: 按照 `docs/development/module-creation.md`
2. **新增 API 端點**: 參考現有 `app/*/apis/` 目錄
3. **設計資料模型**: 查看 `app/*/models/` 範例

### 📚 參與貢獻

1. **查看貢獻指南**: `docs/development/CONTRIBUTING.md`
2. **了解變更流程**: `docs/changes/` 目錄說明
3. **參與文檔維護**: 持續改善中文註解和文檔品質

## 技術支援

如果在使用過程中遇到問題：

1. 查看相關文檔和 FAQ
2. 檢查 `docs/changes/` 中的已知問題
3. 查看程式碼中的中文註解和流程說明
4. 參考 Swagger API 文檔的詳細說明

**記住**: 現在所有的程式碼都有詳細的中文註解，可以直接閱讀程式碼來理解系統運作方式！
