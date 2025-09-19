# Go-Admin 專案術語對照表

此文件定義 Go-Admin 專案中使用的技術術語和翻譯標準，確保文件和程式碼註解的一致性。

---

## 核心技術術語

### Web 框架相關

| 英文術語   | 繁體中文 | 簡體中文 | 說明                |
| ---------- | -------- | -------- | ------------------- |
| Framework  | 框架     | 框架     | Web 開發框架        |
| Middleware | 中間件   | 中间件   | HTTP 請求處理中間件 |
| Controller | 控制器   | 控制器   | MVC 架構中的控制器  |
| Model      | 模型     | 模型     | 數據模型            |
| View       | 視圖     | 视图     | 用戶界面視圖        |
| Router     | 路由器   | 路由器   | HTTP 路由處理器     |
| Handler    | 處理器   | 处理器   | 請求處理函數        |
| Service    | 服務     | 服务     | 業務邏輯服務        |

### 資料庫相關

| 英文術語    | 繁體中文     | 簡體中文     | 說明                      |
| ----------- | ------------ | ------------ | ------------------------- |
| Database    | 資料庫       | 数据库       | 數據存儲系統              |
| Table       | 資料表       | 数据表       | 資料庫表格                |
| Column      | 欄位         | 字段         | 表格欄位                  |
| Primary Key | 主鍵         | 主键         | 主要鍵值                  |
| Foreign Key | 外鍵         | 外键         | 外部參照鍵                |
| Index       | 索引         | 索引         | 資料庫索引                |
| Migration   | 遷移         | 迁移         | 資料庫結構遷移            |
| Seed        | 種子資料     | 种子数据     | 初始化資料                |
| ORM         | 物件關聯對映 | 对象关系映射 | Object-Relational Mapping |
| Transaction | 交易         | 事务         | 資料庫交易                |

### 認證授權相關

| 英文術語       | 繁體中文 | 簡體中文 | 說明         |
| -------------- | -------- | -------- | ------------ |
| Authentication | 身份驗證 | 身份验证 | 用戶身份確認 |
| Authorization  | 授權     | 授权     | 權限控制     |
| Permission     | 權限     | 权限     | 操作權限     |
| Role           | 角色     | 角色     | 用戶角色     |
| User           | 使用者   | 用户     | 系統用戶     |
| Admin          | 管理員   | 管理员   | 系統管理員   |
| Token          | 權杖     | 令牌     | 認證權杖     |
| Session        | 會話     | 会话     | 用戶會話     |
| Login          | 登入     | 登录     | 用戶登入     |
| Logout         | 登出     | 登出     | 用戶登出     |

### 系統架構相關

| 英文術語      | 繁體中文     | 簡體中文     | 說明                              |
| ------------- | ------------ | ------------ | --------------------------------- |
| API           | 應用程式介面 | 应用程序接口 | Application Programming Interface |
| REST          | REST         | REST         | RESTful API 架構                  |
| HTTP          | HTTP         | HTTP         | 超文本傳輸協議                    |
| JSON          | JSON         | JSON         | JavaScript Object Notation        |
| Configuration | 設定         | 配置         | 系統設定                          |
| Environment   | 環境         | 环境         | 執行環境                          |
| Development   | 開發         | 开发         | 開發環境                          |
| Production    | 正式環境     | 生产环境     | 生產環境                          |
| Testing       | 測試         | 测试         | 測試環境                          |
| Deployment    | 部署         | 部署         | 系統部署                          |

### Go 語言特定術語

| 英文術語  | 繁體中文 | 簡體中文 | 說明        |
| --------- | -------- | -------- | ----------- |
| Package   | 套件     | 包       | Go 套件     |
| Module    | 模組     | 模块     | Go 模組     |
| Import    | 匯入     | 导入     | 套件匯入    |
| Interface | 介面     | 接口     | Go 介面     |
| Struct    | 結構體   | 结构体   | 資料結構    |
| Method    | 方法     | 方法     | 結構體方法  |
| Function  | 函數     | 函数     | Go 函數     |
| Goroutine | Go 協程  | Go 协程  | Go 並發機制 |
| Channel   | 通道     | 通道     | Go 通道     |
| Context   | 上下文   | 上下文   | Go Context  |

### 業務邏輯術語

| 英文術語   | 繁體中文 | 簡體中文 | 說明         |
| ---------- | -------- | -------- | ------------ |
| Dashboard  | 儀表板   | 仪表板   | 管理後台首頁 |
| Menu       | 選單     | 菜单     | 系統選單     |
| Department | 部門     | 部门     | 組織部門     |
| Position   | 職位     | 职位     | 工作職位     |
| Log        | 日誌     | 日志     | 系統日誌     |
| Audit      | 稽核     | 审计     | 操作稽核     |
| Backup     | 備份     | 备份     | 資料備份     |
| Cache      | 快取     | 缓存     | 資料快取     |
| Queue      | 佇列     | 队列     | 任務佇列     |
| Job        | 工作     | 作业     | 背景工作     |

---

## 文件寫作規範

### 1. 語言使用原則

- **主要語言**: 使用繁體中文 (台灣用語)
- **技術術語**: 優先使用中文，必要時保留英文原文
- **程式碼**: 變數和函數名稱保持英文，註解使用中文
- **文件標題**: 中英文混合時需要適當的空格分隔

### 2. 標點符號規範

- 中文句子使用中文標點符號 (，。；：！？)
- 英文句子使用英文標點符號 (, . ; : ! ?)
- 中英文混合時，以主要語言的標點符號為準
- 程式碼區塊內使用英文標點符號

### 3. 數字和單位

- 中文語境中的數字優先使用阿拉伯數字
- 單位和度量使用標準符號 (MB, GB, ms, etc.)
- 百分比使用 % 符號
- 時間格式: 24 小時制 (14:30) 或 ISO 格式 (2025-09-19)

### 4. 專有名詞處理

- **產品名稱**: 保持原文 (Go-Admin, Docker, Kubernetes)
- **公司名稱**: 保持原文 (Google, Microsoft, GitHub)
- **技術品牌**: 保持原文但可附中文說明
- **檔案名稱**: 保持原始格式，使用反引號標記

---

## 程式碼註解規範

### 1. 函數註解格式

```go
// CreateUser 建立新使用者帳號
// 參數:
//   - req: 使用者建立請求資料
//   - ctx: 請求上下文
// 返回:
//   - User: 建立的使用者物件
//   - error: 錯誤資訊 (如果有)
func CreateUser(req CreateUserRequest, ctx context.Context) (*User, error) {
    // 驗證使用者輸入資料
    if err := validateUserInput(req); err != nil {
        return nil, fmt.Errorf("輸入資料驗證失敗: %w", err)
    }

    // 建立使用者記錄
    user := &User{
        Name:  req.Name,
        Email: req.Email,
    }

    return user, nil
}
```

### 2. 結構體註解格式

```go
// User 表示系統使用者
type User struct {
    ID       uint   `json:"id" gorm:"primaryKey" comment:"使用者唯一識別碼"`
    Name     string `json:"name" gorm:"size:100" comment:"使用者姓名"`
    Email    string `json:"email" gorm:"uniqueIndex" comment:"電子郵件地址"`
    Password string `json:"-" gorm:"size:255" comment:"密碼雜湊值"`
    RoleID   uint   `json:"role_id" comment:"角色ID"`

    // 關聯關係
    Role *Role `json:"role,omitempty" gorm:"foreignKey:RoleID" comment:"使用者角色"`
}
```

### 3. 常量和變數註解

```go
const (
    // DefaultPageSize 預設分頁大小
    DefaultPageSize = 20

    // MaxPageSize 最大分頁大小
    MaxPageSize = 100

    // TokenExpireDuration 權杖過期時間 (24小時)
    TokenExpireDuration = 24 * time.Hour
)

var (
    // ErrUserNotFound 使用者不存在錯誤
    ErrUserNotFound = errors.New("使用者不存在")

    // ErrInvalidPassword 密碼錯誤
    ErrInvalidPassword = errors.New("密碼錯誤")
)
```

---

## 文件連結規範

### 1. 內部連結格式

- 相對路徑: `[安裝指南](./development/installation.md)`
- 章節連結: `[API 說明](#api-說明)`
- 跨文件連結: `[Docker 部署](../docker/deployment.md)`

### 2. 外部連結格式

- 完整 URL: `[Go 官方網站](https://golang.org/)`
- GitHub 連結: `[專案儲存庫](https://github.com/go-admin-team/go-admin)`

### 3. 圖片連結格式

```markdown
![架構圖](./images/architecture-diagram.png "Go-Admin 系統架構圖")
```

---

## 版本控制規範

### 1. Commit 訊息格式

```
類型(範圍): 簡短描述

詳細說明 (可選)

相關 Issue: #123
```

**類型選項:**

- `docs`: 文件更新
- `feat`: 新功能
- `fix`: 錯誤修復
- `refactor`: 程式碼重構
- `style`: 程式碼格式調整
- `test`: 測試相關

**範例:**

```
docs(readme): 更新安裝指南和快速開始說明

- 添加 Docker 安裝步驟
- 更新環境需求說明
- 修正程式碼範例錯誤

相關 Issue: #456
```

### 2. 分支命名規範

- `main`: 主分支
- `develop`: 開發分支
- `feature/功能名稱`: 功能開發分支
- `hotfix/修復說明`: 緊急修復分支
- `docs/文件說明`: 文件更新分支

---

此術語對照表將持續更新，以確保專案文件的一致性和專業性。如有新的術語需要添加，請參考此規範進行定義。
