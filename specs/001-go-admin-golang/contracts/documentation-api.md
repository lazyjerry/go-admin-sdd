# 文檔化 API 契約

## 概述

本文檔定義文檔化過程中需要的操作介面和流程契約。雖然這是文檔化專案而非 API 開發，但我們定義標準化的操作流程來確保一致性。

## 文檔操作契約

### 1. 文檔建立操作

**操作**: CreateDocument
**輸入**:

```yaml
type: DocumentType # 文檔類型
title: string # 文檔標題
path: string # 檔案路徑
content: string # 文檔內容 (Markdown)
language: "zh-TW" # 語言標識
```

**輸出**:

```yaml
success: boolean # 操作成功與否
message: string # 結果訊息
filePath: string # 建立的檔案路徑
lastModified: timestamp # 最後修改時間
```

**驗證規則**:

- title 不可為空
- path 必須為有效路徑且不存在衝突
- content 必須為有效 Markdown 格式
- 繁體中文內容檢查

### 2. 程式碼註解操作

**操作**: AddCodeAnnotation
**輸入**:

```yaml
filePath: string # 目標檔案路徑
lineNumber: integer # 行號
annotationType: AnnotationType # 註解類型
originalContent: string # 原始內容
chineseAnnotation: string # 中文註解內容
preserveOriginal: boolean # 是否保留原始註解
```

**輸出**:

```yaml
success: boolean # 操作成功與否
message: string # 結果訊息
modifiedLines: integer # 修改的行數
backupPath: string # 備份檔案路徑
```

**驗證規則**:

- filePath 檔案必須存在
- lineNumber 必須在有效範圍內
- chineseAnnotation 必須為繁體中文
- 語法檢查確保註解不破壞程式碼

### 3. 文檔更新操作

**操作**: UpdateDocument
**輸入**:

```yaml
filePath: string          # 現有檔案路徑
sections: []Section       # 要更新的章節
updateType: UpdateType    # 更新類型 (append/replace/insert)
changeReason: string      # 變更原因
```

**Section 結構**:

```yaml
title: string # 章節標題
content: string # 章節內容
position: integer # 插入位置 (僅 insert 時使用)
```

**輸出**:

```yaml
success: boolean # 操作成功與否
message: string # 結果訊息
changelogEntry: string # 變更紀錄項目
oldVersion: string # 舊版本內容雜湊
newVersion: string # 新版本內容雜湊
```

### 4. 翻譯操作

**操作**: TranslateContent
**輸入**:

```yaml
sourceText: string       # 原始文字 (簡體中文或英文)
sourceLanguage: string   # 原始語言 ("en", "zh-CN")
targetLanguage: "zh-TW"  # 目標語言 (繁體中文)
context: string          # 上下文 (技術文檔、註解等)
terminology: []Term      # 專業術語對照
```

**Term 結構**:

```yaml
source: string # 原始術語
target: string # 對應譯詞
category: string # 術語分類
```

**輸出**:

```yaml
translatedText: string   # 翻譯結果
confidence: float       # 信心度 (0-1)
suggestions: []string   # 替代翻譯建議
usedTerms: []Term      # 使用的術語對照
```

## 品質保證契約

### 5. 文檔驗證操作

**操作**: ValidateDocument
**輸入**:

```yaml
filePath: string         # 要驗證的檔案
checkTypes: []CheckType  # 檢查類型
strictMode: boolean      # 嚴格模式
```

**CheckType 列舉**:

- markdown_syntax # Markdown 語法檢查
- chinese_traditional # 繁體中文檢查
- terminology_consistency # 術語一致性
- link_validation # 連結有效性
- structure_compliance # 結構規範性

**輸出**:

```yaml
isValid: boolean        # 整體是否有效
errors: []ValidationError # 錯誤清單
warnings: []ValidationWarning # 警告清單
score: integer         # 品質評分 (0-100)
```

### 6. 變更追蹤操作

**操作**: TrackChange
**輸入**:

```yaml
changeType: ChangeType # 變更類型
filePath: string # 影響檔案
description: string # 變更描述
author: string # 變更者
impact: ImpactLevel # 影響程度
```

**ChangeType 列舉**:

- document_created # 新建文檔
- document_updated # 更新文檔
- annotation_added # 新增註解
- translation_updated # 更新翻譯
- structure_changed # 結構變更

**ImpactLevel 列舉**:

- minor # 輕微變更 (錯字修正)
- moderate # 中等變更 (內容補充)
- major # 重大變更 (結構調整)

**輸出**:

```yaml
changeId: string # 變更識別碼
timestamp: timestamp # 變更時間
changelogPath: string # 變更紀錄檔案路徑
```

## 批次操作契約

### 7. 批次註解操作

**操作**: BatchAnnotateFiles
**輸入**:

```yaml
filePatterns: []string  # 檔案模式 (支援 glob)
excludePatterns: []string # 排除模式
annotationRules: []AnnotationRule # 註解規則
dryRun: boolean        # 僅模擬不實際修改
```

**AnnotationRule 結構**:

```yaml
targetPattern: string # 目標模式 (函數名、類型等)
template: string # 註解模板
priority: integer # 優先順序
```

**輸出**:

```yaml
processedFiles: integer # 處理檔案數
addedAnnotations: integer # 新增註解數
skippedFiles: []string # 跳過的檔案清單
errors: []ProcessingError # 處理錯誤
```

## 錯誤處理契約

### 標準錯誤格式

```yaml
code: string # 錯誤代碼
message: string # 錯誤訊息
details: object # 詳細資訊
suggestion: string # 建議解決方案
```

### 常見錯誤代碼

- `INVALID_PATH`: 無效檔案路徑
- `SYNTAX_ERROR`: 語法錯誤
- `ENCODING_ERROR`: 編碼問題
- `TRANSLATION_FAILED`: 翻譯失敗
- `VALIDATION_FAILED`: 驗證失敗
- `BACKUP_FAILED`: 備份失敗

## 輸出格式標準

所有操作支援以下輸出格式：

- JSON: 程式化處理
- YAML: 人類可讀
- Plain Text: 終端顯示

格式透過 `Accept` 或 `--format` 參數指定。
