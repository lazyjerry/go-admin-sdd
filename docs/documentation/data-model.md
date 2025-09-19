# 資料模型設計

## 文件化實體模型

### 專案文件 (ProjectDocumentation)

**用途**: 定義專案層級的文件內容和結構
**屬性**:

- `name`: 專案名稱 (string) - "go-admin"
- `version`: 版本號 (string) - 基於 git tag 或 go.mod
- `description`: 專案描述 (string) - 功能概述
- `language`: 主要程式語言 (string) - "Go"
- `framework`: 使用框架 ([]string) - ["Gin", "GORM", "Casbin"]
- `lastUpdated`: 最後更新時間 (time.Time)

**關係**:

- 一對多：擁有多個 DocumentModule
