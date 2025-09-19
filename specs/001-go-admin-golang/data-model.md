# 資料模型設計

## 文檔化實體模型

### 專案文檔 (ProjectDocumentation)

**用途**: 定義專案層級的文檔內容和結構
**屬性**:

- `name`: 專案名稱 (string) - "go-admin"
- `version`: 版本號 (string) - 基於 git tag 或 go.mod
- `description`: 專案描述 (string) - 功能概述
- `language`: 主要程式語言 (string) - "Go"
- `framework`: 使用框架 ([]string) - ["Gin", "GORM", "Casbin"]
- `lastUpdated`: 最後更新時間 (time.Time)

**關係**:

- 一對多：擁有多個 DocumentModule

### 文檔模組 (DocumentModule)

**用途**: 定義不同類型文檔的模組化內容
**屬性**:

- `id`: 唯一識別碼 (string)
- `type`: 文檔類型 (DocumentType) - 列舉值
- `title`: 模組標題 (string)
- `path`: 檔案路徑 (string)
- `priority`: 優先順序 (int) - 用於排序
- `status`: 完成狀態 (DocumentStatus)

**關係**:

- 屬於：ProjectDocumentation
- 一對多：包含多個 DocumentSection

### 文檔章節 (DocumentSection)

**用途**: 文檔內部的具體章節內容
**屬性**:

- `id`: 章節識別碼 (string)
- `title`: 章節標題 (string)
- `content`: 章節內容 (string) - Markdown 格式
- `order`: 章節順序 (int)
- `level`: 標題層級 (int) - 1-6 對應 H1-H6

**關係**:

- 屬於：DocumentModule
- 可選：包含多個 CodeExample

### 程式碼範例 (CodeExample)

**用途**: 文檔中的程式碼示範和說明
**屬性**:

- `id`: 範例識別碼 (string)
- `language`: 程式語言 (string) - "go", "bash", "yaml"
- `code`: 程式碼內容 (string)
- `description`: 說明文字 (string)
- `isRunnable`: 是否可執行 (bool)

### 程式碼註解 (CodeAnnotation)

**用途**: 定義程式碼檔案中的中文註解內容
**屬性**:

- `filePath`: 檔案路徑 (string)
- `lineNumber`: 行號 (int)
- `annotationType`: 註解類型 (AnnotationType)
- `originalText`: 原始內容 (string) - 可能是英文或簡體中文
- `translatedText`: 翻譯後內容 (string) - 繁體中文
- `context`: 上下文說明 (string) - 業務邏輯背景

**關係**:

- 屬於：特定的程式碼檔案

## 列舉類型定義

### DocumentType (文檔類型)

```go
type DocumentType string

const (
    TypeREADME      DocumentType = "readme"      // 專案說明
    TypeInstall     DocumentType = "install"     // 安裝指南
    TypeDeploy      DocumentType = "deploy"      // 部署文檔
    TypeAPI         DocumentType = "api"         // API 文檔
    TypeArchitecture DocumentType = "architecture" // 架構說明
    TypeDocker      DocumentType = "docker"      // Docker 文檔
    TypeDevelopment DocumentType = "development" // 開發指南
    TypeChangelog   DocumentType = "changelog"   // 變更紀錄
    TypeFAQ         DocumentType = "faq"         // 常見問題
)
```

### DocumentStatus (文檔狀態)

```go
type DocumentStatus string

const (
    StatusPlanned    DocumentStatus = "planned"    // 規劃中
    StatusInProgress DocumentStatus = "in_progress" // 進行中
    StatusReview     DocumentStatus = "review"     // 審核中
    StatusCompleted  DocumentStatus = "completed"  // 已完成
    StatusObsolete   DocumentStatus = "obsolete"   // 已過時
)
```

### AnnotationType (註解類型)

```go
type AnnotationType string

const (
    TypeFunction    AnnotationType = "function"    // 函數註解
    TypeStruct      AnnotationType = "struct"      // 結構體註解
    TypeInterface   AnnotationType = "interface"   // 介面註解
    TypePackage     AnnotationType = "package"     // 套件註解
    TypeVariable    AnnotationType = "variable"    // 變數註解
    TypeConstant    AnnotationType = "constant"    // 常數註解
    TypeLogic       AnnotationType = "logic"       // 業務邏輯註解
    TypeConfig      AnnotationType = "config"      // 設定註解
)
```

## 驗證規則

### ProjectDocumentation 驗證

- `name`: 非空，長度 1-100 字元
- `version`: 符合語義版本格式 (semver)
- `description`: 非空，長度 10-500 字元

### DocumentModule 驗證

- `type`: 必須為有效的 DocumentType 列舉值
- `title`: 非空，長度 1-200 字元
- `path`: 有效的檔案路徑格式
- `priority`: 正整數

### DocumentSection 驗證

- `title`: 非空，長度 1-200 字元
- `content`: Markdown 格式驗證
- `order`: 正整數，同一模組內唯一
- `level`: 1-6 之間的整數

### CodeAnnotation 驗證

- `filePath`: 檔案必須存在於專案中
- `lineNumber`: 正整數，不超過檔案總行數
- `translatedText`: 非空，繁體中文字元檢查
- `annotationType`: 必須為有效的 AnnotationType 列舉值

## 狀態轉換

### DocumentStatus 狀態流程

```
planned → in_progress → review → completed
   ↓           ↓          ↓
obsolete ← obsolete ← obsolete
```

**轉換規則**:

- planned → in_progress: 開始撰寫文檔
- in_progress → review: 文檔內容初稿完成
- review → completed: 審核通過
- review → in_progress: 需要修改
- completed → obsolete: 內容過時需要更新
- - → obsolete: 任何狀態都可標記為過時

## 資料持久化考慮

### 檔案系統結構

- 文檔以 Markdown 檔案形式存儲
- 元資料可選用 YAML Front Matter
- 程式碼註解直接修改原始檔案

### 版本控制整合

- 所有文檔變更通過 Git 追蹤
- 使用 Git hooks 自動更新 lastUpdated 時間戳
- 支援分支開發，文檔與程式碼同步版本管理
