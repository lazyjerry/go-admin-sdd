# 測試指南

本文件詳細說明 Go-Admin 專案的測試策略、測試框架使用和最佳實踐。

## 測試架構

### 測試層級

```
測試金字塔
    ↗ E2E 測試 (End-to-End)
   ↗  整合測試 (Integration)
  ↗   單元測試 (Unit Tests)
 ↗    靜態分析 (Static Analysis)
```

### 測試類型

| 測試類型     | 範圍          | 工具                | 執行頻率 |
| ------------ | ------------- | ------------------- | -------- |
| **單元測試** | 函數/方法級別 | Go testing, testify | 每次提交 |
| **整合測試** | 模組間互動    | testcontainers      | 每日構建 |
| **API 測試** | HTTP 介面     | httptest, gin       | 每次發布 |
| **效能測試** | 系統效能      | go test -bench      | 每週執行 |
| **安全測試** | 安全掃描      | gosec, nancy        | 持續整合 |

## 測試環境設定

### 依賴安裝

```bash
# 測試框架
go get github.com/stretchr/testify
go get github.com/gin-gonic/gin
go get github.com/testcontainers/testcontainers-go

# 模擬工具
go get github.com/golang/mock/gomock
go get github.com/golang/mock/mockgen

# 資料庫測試
go get github.com/DATA-DOG/go-sqlmock
go get gorm.io/driver/sqlite
```

### 測試配置

```yaml
# config/settings.test.yml
app:
  env: testing
  debug: true

database:
  driver: sqlite
  source: ":memory:"

redis:
  host: localhost
  port: 6379
  db: 1

log:
  level: debug
  format: text
```

## 單元測試

### 測試結構

```go
// 標準測試檔案結構
package service

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

// 測試套件
type UserServiceTestSuite struct {
    suite.Suite
    service *UserService
    db      *gorm.DB
}

// 設定測試環境
func (s *UserServiceTestSuite) SetupTest() {
    // 初始化測試資料庫
    s.db = setupTestDB()
    s.service = &UserService{DB: s.db}
}

// 清理測試環境
func (s *UserServiceTestSuite) TearDownTest() {
    cleanupTestDB(s.db)
}

// 執行測試套件
func TestUserServiceSuite(t *testing.T) {
    suite.Run(t, new(UserServiceTestSuite))
}
```

### 測試範例

```go
// service/user_test.go
func (s *UserServiceTestSuite) TestCreateUser() {
    // Arrange (準備)
    req := &dto.UserCreateReq{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "password123",
    }

    // Act (執行)
    err := s.service.CreateUser(req)

    // Assert (驗證)
    s.NoError(err)
    s.NotNil(req.GetId())

    // 驗證資料庫記錄
    var user models.User
    err = s.db.First(&user, req.GetId()).Error
    s.NoError(err)
    s.Equal("testuser", user.Username)
}

func (s *UserServiceTestSuite) TestCreateUser_DuplicateUsername() {
    // 測試重複使用者名稱的錯誤處理
    req1 := &dto.UserCreateReq{
        Username: "duplicate",
        Email:    "user1@example.com",
    }
    req2 := &dto.UserCreateReq{
        Username: "duplicate", // 相同使用者名稱
        Email:    "user2@example.com",
    }

    // 建立第一個使用者
    err := s.service.CreateUser(req1)
    s.NoError(err)

    // 建立重複使用者應該失敗
    err = s.service.CreateUser(req2)
    s.Error(err)
    s.Contains(err.Error(), "使用者名稱已存在")
}
```

## API 測試

### HTTP 測試

```go
// apis/user_test.go
package apis

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestUserApi_CreateUser(t *testing.T) {
    // 設定 Gin 測試模式
    gin.SetMode(gin.TestMode)

    // 建立測試路由
    router := gin.New()
    api := &UserApi{}
    router.POST("/users", api.CreateUser)

    // 準備測試資料
    reqData := map[string]interface{}{
        "username": "testuser",
        "email":    "test@example.com",
        "password": "password123",
    }
    jsonData, _ := json.Marshal(reqData)

    // 建立測試請求
    req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    // 執行請求
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)

    // 驗證回應
    assert.Equal(t, http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, "建立成功", response["msg"])
}
```

### 認證測試

```go
func TestUserApi_WithAuth(t *testing.T) {
    // 測試需要認證的端點
    router := setupAuthRouter()

    // 無認證請求
    req, _ := http.NewRequest("GET", "/api/v1/users", nil)
    w := httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusUnauthorized, w.Code)

    // 有效認證請求
    token := generateTestToken("testuser")
    req.Header.Set("Authorization", "Bearer "+token)
    w = httptest.NewRecorder()
    router.ServeHTTP(w, req)
    assert.Equal(t, http.StatusOK, w.Code)
}
```

## 資料庫測試

### 記憶體資料庫

```go
// test/database.go
func SetupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Silent),
    })
    if err != nil {
        panic("failed to connect test database")
    }

    // 自動遷移測試表
    err = db.AutoMigrate(
        &models.User{},
        &models.Role{},
        &models.Department{},
    )
    if err != nil {
        panic("failed to migrate test database")
    }

    return db
}
```

### 測試資料工廠

```go
// test/factories/user_factory.go
package factories

import (
    "go-admin/app/admin/models"
    "github.com/bxcodec/faker/v3"
)

func CreateUser(overrides ...func(*models.User)) *models.User {
    user := &models.User{
        Username: faker.Username(),
        Email:    faker.Email(),
        Nickname: faker.Name(),
        Status:   "1",
    }

    // 應用自訂覆蓋
    for _, override := range overrides {
        override(user)
    }

    return user
}

func CreateUserWithRole(roleId int) *models.User {
    return CreateUser(func(u *models.User) {
        u.RoleId = roleId
    })
}
```

## Mock 測試

### 生成 Mock

```bash
# 安裝 mockgen
go install github.com/golang/mock/mockgen@latest

# 生成 Mock 介面
mockgen -source=service/user.go -destination=mocks/user_service_mock.go
```

### 使用 Mock

```go
// service/user_test.go
import (
    "go-admin/mocks"
    "github.com/golang/mock/gomock"
)

func TestUserService_WithMock(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 建立 Mock 物件
    mockRepo := mocks.NewMockUserRepository(ctrl)
    service := &UserService{Repo: mockRepo}

    // 設定 Mock 期望
    mockRepo.EXPECT().
        FindByUsername("testuser").
        Return(nil, nil).  // 返回值
        Times(1)           // 呼叫次數

    // 執行測試
    user, err := service.GetByUsername("testuser")

    // 驗證結果
    assert.NoError(t, err)
    assert.Nil(t, user)
}
```

## 整合測試

### Testcontainers 測試

```go
// integration/database_test.go
package integration

import (
    "context"
    "testing"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/modules/mysql"
)

func TestDatabaseIntegration(t *testing.T) {
    ctx := context.Background()

    // 啟動 MySQL 容器
    mysqlContainer, err := mysql.RunContainer(ctx,
        testcontainers.WithImage("mysql:8.0"),
        mysql.WithDatabase("testdb"),
        mysql.WithUsername("root"),
        mysql.WithPassword("password"),
    )
    if err != nil {
        t.Fatal(err)
    }
    defer mysqlContainer.Terminate(ctx)

    // 獲取連接資訊
    host, _ := mysqlContainer.Host(ctx)
    port, _ := mysqlContainer.MappedPort(ctx, "3306")

    // 建立資料庫連接並測試
    dsn := fmt.Sprintf("root:password@tcp(%s:%s)/testdb", host, port.Port())
    db := setupDatabase(dsn)

    // 執行整合測試
    testUserCRUD(t, db)
}
```

## 效能測試

### 基準測試

```go
// service/user_benchmark_test.go
func BenchmarkUserService_CreateUser(b *testing.B) {
    db := SetupTestDB()
    service := &UserService{DB: db}

    b.ResetTimer() // 重置計時器

    for i := 0; i < b.N; i++ {
        req := &dto.UserCreateReq{
            Username: fmt.Sprintf("user%d", i),
            Email:    fmt.Sprintf("user%d@example.com", i),
        }
        service.CreateUser(req)
    }
}

func BenchmarkUserService_GetUser(b *testing.B) {
    // 效能測試：查詢使用者
    db := SetupTestDB()
    service := &UserService{DB: db}

    // 預先建立測試資料
    user := createTestUser(db)

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        _, err := service.GetUser(user.Id)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 並發測試

```go
func TestUserService_Concurrent(t *testing.T) {
    db := SetupTestDB()
    service := &UserService{DB: db}

    // 並發建立使用者
    var wg sync.WaitGroup
    errors := make(chan error, 100)

    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()

            req := &dto.UserCreateReq{
                Username: fmt.Sprintf("concurrent_user_%d", index),
                Email:    fmt.Sprintf("user%d@concurrent.com", index),
            }

            if err := service.CreateUser(req); err != nil {
                errors <- err
            }
        }(i)
    }

    wg.Wait()
    close(errors)

    // 檢查是否有錯誤
    for err := range errors {
        t.Errorf("Concurrent create error: %v", err)
    }
}
```

## 測試執行

### 基本指令

```bash
# 執行所有測試
go test ./...

# 執行特定套件測試
go test ./app/admin/service

# 執行特定測試函數
go test -run TestUserService_CreateUser

# 顯示詳細輸出
go test -v ./...

# 測試覆蓋率
go test -cover ./...

# 生成覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 進階選項

```bash
# 效能測試
go test -bench=. -benchmem ./...

# 競態條件檢測
go test -race ./...

# 超時設定
go test -timeout 30s ./...

# 並行執行
go test -parallel 4 ./...

# 快取禁用
go test -count=1 ./...
```

## CI/CD 整合

### GitHub Actions

```yaml
# .github/workflows/test.yml
name: Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: password
          MYSQL_DATABASE: test_db
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
      - uses: actions/checkout@v3

      - name: Set up Go
        uses: actions/setup-go@v3
        with:
          go-version: 1.21

      - name: Cache dependencies
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}

      - name: Download dependencies
        run: go mod download

      - name: Run tests
        run: |
          go test -v -race -coverprofile=coverage.out ./...
          go tool cover -func=coverage.out

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

## 測試最佳實踐

### 1. 測試命名

```go
// 好的測試命名
func TestUserService_CreateUser_Success(t *testing.T) {}
func TestUserService_CreateUser_DuplicateUsername_ReturnsError(t *testing.T) {}
func TestUserService_CreateUser_InvalidEmail_ReturnsValidationError(t *testing.T) {}

// 避免的命名
func TestCreateUser(t *testing.T) {}
func TestUser(t *testing.T) {}
```

### 2. 測試組織

```go
// 使用子測試組織相關測試
func TestUserService_CreateUser(t *testing.T) {
    t.Run("Success", func(t *testing.T) {
        // 成功案例測試
    })

    t.Run("DuplicateUsername", func(t *testing.T) {
        // 重複使用者名稱測試
    })

    t.Run("InvalidInput", func(t *testing.T) {
        // 無效輸入測試
    })
}
```

### 3. 測試資料管理

```go
// 使用表格驅動測試
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name     string
        email    string
        expected bool
    }{
        {"Valid email", "user@example.com", true},
        {"Invalid email", "invalid-email", false},
        {"Empty email", "", false},
        {"Email without domain", "user@", false},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := validateEmail(tt.email)
            assert.Equal(t, tt.expected, result)
        })
    }
}
```

### 4. 測試隔離

```go
// 每個測試都應該是獨立的
func (s *UserServiceTestSuite) SetupTest() {
    // 每個測試前重置狀態
    s.db.Exec("DELETE FROM users")
    s.db.Exec("DELETE FROM roles")
}
```

## 測試報告

### 覆蓋率報告

```bash
# 生成 HTML 覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# 查看覆蓋率統計
go tool cover -func=coverage.out
```

### 測試報告生成

```go
// 使用 testify 生成測試報告
func TestMain(m *testing.M) {
    // 設定測試環境
    setup()

    // 執行測試
    code := m.Run()

    // 清理測試環境
    teardown()

    os.Exit(code)
}
```

## 相關文件

- [程式碼品質標準](./code-standards.md)
- [Go-Admin 架構說明](./architecture.md)
- [API 開發指南](./api-guide.md)
- [貢獻指南](./CONTRIBUTING.md)

## 版本記錄

| 版本  | 日期       | 更新內容               |
| ----- | ---------- | ---------------------- |
| 1.0.0 | 2025-09-19 | 初始版本，完整測試指南 |

---

_最後更新：2025 年 9 月 19 日_
