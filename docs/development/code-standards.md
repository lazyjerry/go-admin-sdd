# 程式碼品質標準

本文件定義 Go-Admin 專案的程式碼品質標準，包含程式碼風格、最佳實踐和品質檢查工具。

## 程式碼風格

### Go 語言規範

遵循官方 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) 和 [Effective Go](https://golang.org/doc/effective_go.html) 指南。

### 格式化工具

```bash
# 使用 gofmt 格式化程式碼
gofmt -w .

# 使用 goimports 處理 import
goimports -w .

# 使用 gci 排序 import
gci write --skip-generated -s standard -s default -s "prefix(go-admin)" .
```

### 命名規範

#### 1. 變數命名

```go
// ✅ 好的命名
var userID int
var maxRetryCount int
var httpClient *http.Client

// ❌ 避免的命名
var uid int
var max int
var client *http.Client
```

#### 2. 函數命名

```go
// ✅ 公開函數使用大寫開頭
func GetUserByID(id int) (*User, error) {}
func CreateUser(req *CreateUserRequest) error {}

// ✅ 私有函數使用小寫開頭
func validateEmail(email string) bool {}
func hashPassword(password string) string {}

// ✅ 介面命名使用 -er 後綴
type UserRepository interface {
    FindByID(id int) (*User, error)
}
```

#### 3. 結構體命名

```go
// ✅ 結構體使用大駝峰命名
type UserProfile struct {
    ID       int    `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
}

// ✅ 方法接收器使用簡短名稱
func (u *User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}
```

### 註釋規範

#### 1. 套件註釋

```go
// Package admin 提供管理員相關的業務邏輯處理
//
// 主要功能包括：
// - 使用者管理
// - 角色權限管理
// - 系統設定管理
//
// 範例使用：
//   service := admin.NewUserService(db)
//   user, err := service.GetByID(1)
package admin
```

#### 2. 函數註釋

```go
// GetUserByID 根據使用者ID獲取使用者資訊
//
// 參數：
//   id: 使用者ID，必須大於0
//
// 返回值：
//   *User: 使用者資訊，如果不存在則為nil
//   error: 錯誤資訊，如果沒有錯誤則為nil
//
// 範例：
//   user, err := GetUserByID(123)
//   if err != nil {
//       log.Printf("獲取使用者失敗: %v", err)
//       return
//   }
func GetUserByID(id int) (*User, error) {
    // 實作邏輯...
}
```

#### 3. 結構體註釋

```go
// User 代表系統中的使用者實體
//
// 包含使用者的基本資訊和狀態，支援軟刪除。
// 密碼欄位經過 bcrypt 加密儲存。
type User struct {
    ID       int       `json:"id" gorm:"primaryKey"`           // 使用者唯一識別碼
    Username string    `json:"username" gorm:"uniqueIndex"`     // 使用者名稱，必須唯一
    Email    string    `json:"email" gorm:"uniqueIndex"`        // 電子信箱，必須唯一
    Password string    `json:"-" gorm:"not null"`               // 加密後的密碼
    Status   int       `json:"status" gorm:"default:1"`         // 狀態：1啟用，0停用

    CreatedAt time.Time      `json:"created_at"`                 // 建立時間
    UpdatedAt time.Time      `json:"updated_at"`                 // 更新時間
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`             // 刪除時間，軟刪除
}
```

## 程式碼組織

### 檔案結構

```
app/admin/
├── apis/           # API 層
│   ├── user.go     # 使用者 API
│   └── role.go     # 角色 API
├── models/         # 資料模型
│   ├── user.go     # 使用者模型
│   └── role.go     # 角色模型
├── services/       # 業務邏輯
│   ├── user.go     # 使用者服務
│   └── role.go     # 角色服務
├── dto/           # 資料傳輸物件
│   ├── user.go     # 使用者 DTO
│   └── role.go     # 角色 DTO
└── router/        # 路由配置
    └── router.go   # 路由註冊
```

### Import 規範

```go
import (
    // 標準庫
    "context"
    "fmt"
    "time"

    // 第三方套件
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    // 專案內部套件
    "go-admin/common/models"
    "go-admin/common/response"
)
```

### 常數定義

```go
// 使用 iota 定義相關常數
type UserStatus int

const (
    UserStatusInactive UserStatus = iota // 0: 停用
    UserStatusActive                     // 1: 啟用
    UserStatusBlocked                    // 2: 封鎖
)

// 字串常數使用有意義的名稱
const (
    DefaultPageSize     = 20
    MaxPageSize        = 100
    DefaultTimeout     = 30 * time.Second

    // 快取鍵前綴
    CacheKeyUserPrefix = "user:"
    CacheKeyRolePrefix = "role:"
)
```

## 錯誤處理

### 錯誤定義

```go
// 使用自訂錯誤類型
type BusinessError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

func (e BusinessError) Error() string {
    return e.Message
}

// 預定義錯誤
var (
    ErrUserNotFound     = BusinessError{Code: 1001, Message: "使用者不存在"}
    ErrInvalidPassword  = BusinessError{Code: 1002, Message: "密碼錯誤"}
    ErrDuplicateUser    = BusinessError{Code: 1003, Message: "使用者已存在"}
)
```

### 錯誤處理模式

```go
// ✅ 好的錯誤處理
func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid user id: %d", id)
    }

    user, err := userRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    if user == nil {
        return nil, ErrUserNotFound
    }

    return user, nil
}

// ❌ 避免忽略錯誤
func BadExample() {
    user, _ := GetUser(123) // 不要忽略錯誤
    fmt.Println(user.Name)
}
```

## 效能最佳化

### 1. 記憶體管理

```go
// ✅ 預先分配 slice 容量
func ProcessUsers(count int) []User {
    users := make([]User, 0, count) // 預先分配容量
    // 處理邏輯...
    return users
}

// ✅ 使用物件池減少 GC 壓力
var userPool = sync.Pool{
    New: func() interface{} {
        return &User{}
    },
}

func GetUserFromPool() *User {
    return userPool.Get().(*User)
}

func PutUserToPool(user *User) {
    user.Reset() // 重置物件狀態
    userPool.Put(user)
}
```

### 2. 資料庫最佳化

```go
// ✅ 使用批量操作
func CreateUsers(users []User) error {
    // 使用 CreateInBatches 而不是逐個建立
    return db.CreateInBatches(&users, 100).Error
}

// ✅ 使用 Select 指定欄位
func GetUserBasicInfo(id int) (*User, error) {
    var user User
    err := db.Select("id", "username", "email").
        Where("id = ?", id).
        First(&user).Error
    return &user, err
}

// ✅ 使用預載入避免 N+1 問題
func GetUsersWithRoles() ([]User, error) {
    var users []User
    err := db.Preload("Role").Find(&users).Error
    return users, err
}
```

### 3. 並發控制

```go
// ✅ 使用 context 控制超時
func GetUserWithTimeout(ctx context.Context, id int) (*User, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()

    type result struct {
        user *User
        err  error
    }

    ch := make(chan result, 1)
    go func() {
        user, err := userRepo.FindByID(id)
        ch <- result{user: user, err: err}
    }()

    select {
    case res := <-ch:
        return res.user, res.err
    case <-ctx.Done():
        return nil, ctx.Err()
    }
}

// ✅ 使用 sync.Once 確保初始化一次
var (
    dbInstance *gorm.DB
    dbOnce     sync.Once
)

func GetDB() *gorm.DB {
    dbOnce.Do(func() {
        dbInstance = initDB()
    })
    return dbInstance
}
```

## 安全最佳實踐

### 1. 輸入驗證

```go
// ✅ 嚴格的輸入驗證
func ValidateCreateUserRequest(req *CreateUserRequest) error {
    if req == nil {
        return errors.New("request is nil")
    }

    // 使用者名稱驗證
    if len(req.Username) < 3 || len(req.Username) > 32 {
        return errors.New("使用者名稱長度必須在3-32字元之間")
    }

    if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(req.Username) {
        return errors.New("使用者名稱只能包含字母、數字和底線")
    }

    // 電子信箱驗證
    if !isValidEmail(req.Email) {
        return errors.New("無效的電子信箱格式")
    }

    // 密碼強度驗證
    if !isStrongPassword(req.Password) {
        return errors.New("密碼不符合安全要求")
    }

    return nil
}
```

### 2. SQL 注入防護

```go
// ✅ 使用參數化查詢
func GetUserByUsername(username string) (*User, error) {
    var user User
    // GORM 自動處理 SQL 注入防護
    err := db.Where("username = ?", username).First(&user).Error
    return &user, err
}

// ❌ 避免字串拼接
func BadGetUser(username string) (*User, error) {
    var user User
    // 危險！容易 SQL 注入
    query := fmt.Sprintf("SELECT * FROM users WHERE username = '%s'", username)
    err := db.Raw(query).Scan(&user).Error
    return &user, err
}
```

### 3. 敏感資訊處理

```go
// ✅ 密碼加密
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hash), nil
}

// ✅ 敏感資料脫敏
func (u *User) MaskSensitiveData() {
    if len(u.Phone) >= 7 {
        u.Phone = u.Phone[:3] + "****" + u.Phone[len(u.Phone)-4:]
    }

    if len(u.Email) > 0 {
        parts := strings.Split(u.Email, "@")
        if len(parts) == 2 && len(parts[0]) > 2 {
            masked := parts[0][:2] + "***" + "@" + parts[1]
            u.Email = masked
        }
    }
}
```

## 測試規範

### 測試覆蓋率要求

- **單元測試覆蓋率**: >= 80%
- **整合測試覆蓋率**: >= 60%
- **關鍵業務邏輯**: >= 90%

### 測試命名規範

```go
// 測試函數命名: Test[FunctionName]_[Scenario]_[Expected]
func TestCreateUser_ValidInput_Success(t *testing.T) {}
func TestCreateUser_DuplicateUsername_ReturnsError(t *testing.T) {}
func TestCreateUser_InvalidEmail_ReturnsValidationError(t *testing.T) {}

// 基準測試命名: Benchmark[FunctionName]_[Scenario]
func BenchmarkCreateUser_SingleUser(b *testing.B) {}
func BenchmarkCreateUser_BatchUsers(b *testing.B) {}
```

## 品質檢查工具

### 靜態分析工具

```bash
# golangci-lint - 綜合 linter
golangci-lint run ./...

# gosec - 安全檢查
gosec ./...

# staticcheck - 靜態分析
staticcheck ./...

# go vet - Go 官方檢查工具
go vet ./...
```

### golangci-lint 配置

```yaml
# .golangci.yml
linters-settings:
  govet:
    check-shadowing: true
  gocyclo:
    min-complexity: 15
  maligned:
    suggest-new: true
  dupl:
    threshold: 100
  goconst:
    min-len: 2
    min-occurrences: 2
  misspell:
    locale: US
  lll:
    line-length: 120

linters:
  enable:
    - bodyclose
    - deadcode
    - depguard
    - dogsled
    - dupl
    - errcheck
    - gochecknoinits
    - goconst
    - gocyclo
    - gofmt
    - goimports
    - golint
    - gosec
    - gosimple
    - govet
    - ineffassign
    - interfacer
    - lll
    - maligned
    - misspell
    - nakedret
    - scopelint
    - staticcheck
    - structcheck
    - stylecheck
    - typecheck
    - unconvert
    - unparam
    - unused
    - varcheck
    - whitespace

issues:
  exclude-rules:
    - path: _test\.go
      linters:
        - gocyclo
        - errcheck
        - dupl
        - gosec
```

### Pre-commit Hook

```bash
#!/bin/sh
# .git/hooks/pre-commit

echo "Running pre-commit checks..."

# 格式化檢查
if ! gofmt -l . | grep -q '^$'; then
    echo "Error: Code is not formatted. Run 'gofmt -w .'"
    exit 1
fi

# 靜態分析
if ! golangci-lint run ./...; then
    echo "Error: Linting failed"
    exit 1
fi

# 測試執行
if ! go test -race -short ./...; then
    echo "Error: Tests failed"
    exit 1
fi

echo "All checks passed!"
```

## 效能監控

### 效能基準

```go
// 建立效能基準測試
func BenchmarkUserService_CreateUser(b *testing.B) {
    db := setupTestDB()
    service := NewUserService(db)

    req := &CreateUserRequest{
        Username: "benchmark_user",
        Email:    "benchmark@example.com",
        Password: "password123",
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        req.Username = fmt.Sprintf("user_%d", i)
        req.Email = fmt.Sprintf("user_%d@example.com", i)

        err := service.CreateUser(req)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 效能監控指標

- **API 回應時間**: P95 < 200ms
- **資料庫查詢**: 平均 < 50ms
- **記憶體使用**: 穩定無洩漏
- **CPU 使用率**: 正常負載 < 70%

## 文檔要求

### 1. README 檔案

每個套件都應包含 README.md 說明：

````markdown
# Package Admin

## 概述

管理員模組提供使用者、角色和權限管理功能。

## 主要功能

- 使用者 CRUD 操作
- 角色權限管理
- 登入認證

## 使用範例

```go
service := admin.NewUserService(db)
user, err := service.CreateUser(&CreateUserRequest{
    Username: "admin",
    Email:    "admin@example.com",
})
```
````

## API 文檔

詳見 [API 文檔](./docs/api.md)

````

### 2. API 文檔

使用 Swagger 註釋生成 API 文檔：

```go
// CreateUser 建立新使用者
// @Summary 建立使用者
// @Description 建立新的系統使用者
// @Tags 使用者管理
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "使用者資訊"
// @Success 200 {object} response.Response{data=User}
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /api/v1/users [post]
func (api *UserAPI) CreateUser(c *gin.Context) {
    // 實作...
}
````

## 版本控制

### Git 提交規範

```bash
# 提交訊息格式
<type>(<scope>): <subject>

<body>

<footer>
```

### 提交類型

- `feat`: 新功能
- `fix`: 錯誤修復
- `docs`: 文檔更新
- `style`: 程式碼格式化
- `refactor`: 程式碼重構
- `test`: 測試相關
- `chore`: 建置工具或輔助工具的變動

### 範例提交

```bash
git commit -m "feat(user): add user profile management

- Add user profile CRUD operations
- Implement profile image upload
- Add profile validation rules

Closes #123"
```

## 相關文件

- [測試指南](./testing.md)
- [Go-Admin 架構說明](./architecture.md)
- [API 開發指南](./api-guide.md)
- [貢獻指南](./CONTRIBUTING.md)

## 版本記錄

| 版本  | 日期       | 更新內容                     |
| ----- | ---------- | ---------------------------- |
| 1.0.0 | 2025-09-19 | 初始版本，完整程式碼品質標準 |

---

_最後更新：2025 年 9 月 19 日_
