# Go-Admin 測試指南

本指南提供 Go-Admin 專案的完整測試策略、工具使用和最佳實務，幫助開發者建立高品質的測試套件。

## 目錄

1. [測試概述](#測試概述)
2. [測試環境設定](#測試環境設定)
3. [單元測試](#單元測試)
4. [整合測試](#整合測試)
5. [API 測試](#api-測試)
6. [效能測試](#效能測試)
7. [測試工具與框架](#測試工具與框架)
8. [測試最佳實務](#測試最佳實務)
9. [持續整合測試](#持續整合測試)

## 測試概述

### 測試策略

Go-Admin 採用多層次測試策略：

```
┌─────────────────┐
│   E2E 測試      │  ← 完整功能流程測試
├─────────────────┤
│   整合測試      │  ← API、資料庫整合
├─────────────────┤
│   單元測試      │  ← 個別函數、模組
└─────────────────┘
```

### 測試類型

- **單元測試 (Unit Tests)**: 測試個別函數和方法
- **整合測試 (Integration Tests)**: 測試元件間的互動
- **API 測試**: 測試 REST API 端點
- **效能測試 (Performance Tests)**: 測試系統效能和負載
- **端對端測試 (E2E Tests)**: 測試完整使用者流程

## 測試環境設定

### 1. 安裝測試依賴

```bash
# 進入專案目錄
cd go-admin

# 安裝測試相關工具
go install github.com/stretchr/testify@latest
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega@latest
go install github.com/golang/mock/mockgen@latest
go install github.com/vektra/mockery/v2@latest
```

### 2. 測試設定檔案

建立測試專用的設定檔案 `config/settings.test.yml`：

```yaml
# 測試環境設定
settings:
  application:
    mode: test
    name: go-admin-test
    port: 8080
    host: 127.0.0.1

  database:
    dbtype: sqlite3
    path: test.db
    config: cache=shared&mode=memory
    log-mode: silent
    log-zap: false

  jwt:
    secret: test-jwt-secret-key
    timeout: 1

  log:
    level: error
    path: ./test_logs
    filename: test.log
```

### 3. 測試資料庫設定

```bash
# 建立測試資料庫初始化腳本
cat > scripts/test_db_setup.sql << EOF
-- 測試資料庫初始化
CREATE DATABASE IF NOT EXISTS go_admin_test;
USE go_admin_test;

-- 建立測試使用者
CREATE USER IF NOT EXISTS 'test_user'@'localhost' IDENTIFIED BY 'test_pass';
GRANT ALL PRIVILEGES ON go_admin_test.* TO 'test_user'@'localhost';
FLUSH PRIVILEGES;
EOF
```

### 4. 環境變數設定

建立測試專用的 `.env.test` 檔案：

```bash
# 測試環境變數
APP_MODE=test
DB_TYPE=sqlite3
DB_PATH=test.db
LOG_LEVEL=error
JWT_SECRET=test-secret
```

## 單元測試

### 1. 基本測試結構

```go
// test/unit/user_test.go
package unit

import (
    "testing"
    "go-admin/app/admin/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
)

// UserTestSuite 使用者相關測試套件
type UserTestSuite struct {
    suite.Suite
    user *models.SysUser
}

// SetupTest 測試前置設定
func (suite *UserTestSuite) SetupTest() {
    suite.user = &models.SysUser{
        Username: "testuser",
        NickName: "測試使用者",
        Email:    "test@example.com",
    }
}

// TestUserValidation 測試使用者驗證
func (suite *UserTestSuite) TestUserValidation() {
    // 測試有效使用者
    err := suite.user.Validate()
    assert.NoError(suite.T(), err)

    // 測試無效使用者名稱
    suite.user.Username = ""
    err = suite.user.Validate()
    assert.Error(suite.T(), err)
}

// TestPasswordHashing 測試密碼雜湊
func (suite *UserTestSuite) TestPasswordHashing() {
    password := "testpassword"

    // 測試密碼雜湊
    hashedPassword, err := suite.user.HashPassword(password)
    assert.NoError(suite.T(), err)
    assert.NotEmpty(suite.T(), hashedPassword)
    assert.NotEqual(suite.T(), password, hashedPassword)

    // 測試密碼驗證
    isValid := suite.user.CheckPassword(password, hashedPassword)
    assert.True(suite.T(), isValid)
}

// 執行測試套件
func TestUserTestSuite(t *testing.T) {
    suite.Run(t, new(UserTestSuite))
}
```

### 2. 服務層測試

```go
// test/unit/user_service_test.go
package unit

import (
    "context"
    "testing"
    "go-admin/app/admin/service"
    "go-admin/app/admin/models"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// MockUserRepository 模擬使用者儲存庫
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) GetByID(id int) (*models.SysUser, error) {
    args := m.Called(id)
    return args.Get(0).(*models.SysUser), args.Error(1)
}

func (m *MockUserRepository) Create(user *models.SysUser) error {
    args := m.Called(user)
    return args.Error(0)
}

// TestUserService 使用者服務測試
func TestUserService(t *testing.T) {
    // 建立模擬儲存庫
    mockRepo := new(MockUserRepository)
    userService := service.NewUserService(mockRepo)

    t.Run("建立使用者成功", func(t *testing.T) {
        // 準備測試資料
        user := &models.SysUser{
            Username: "newuser",
            Email:    "newuser@example.com",
        }

        // 設定模擬行為
        mockRepo.On("Create", user).Return(nil)

        // 執行測試
        err := userService.CreateUser(context.Background(), user)

        // 驗證結果
        assert.NoError(t, err)
        mockRepo.AssertExpectations(t)
    })

    t.Run("取得使用者成功", func(t *testing.T) {
        // 準備測試資料
        expectedUser := &models.SysUser{
            Model:    models.Model{Id: 1},
            Username: "testuser",
        }

        // 設定模擬行為
        mockRepo.On("GetByID", 1).Return(expectedUser, nil)

        // 執行測試
        user, err := userService.GetUser(context.Background(), 1)

        // 驗證結果
        assert.NoError(t, err)
        assert.Equal(t, expectedUser.Username, user.Username)
        mockRepo.AssertExpectations(t)
    })
}
```

### 3. 工具函數測試

```go
// test/unit/utils_test.go
package unit

import (
    "testing"
    "go-admin/common/utils"
    "github.com/stretchr/testify/assert"
)

// TestPasswordUtils 密碼工具測試
func TestPasswordUtils(t *testing.T) {
    tests := []struct {
        name     string
        password string
        wantErr  bool
    }{
        {
            name:     "有效密碼",
            password: "ValidPassword123",
            wantErr:  false,
        },
        {
            name:     "密碼太短",
            password: "123",
            wantErr:  true,
        },
        {
            name:     "空密碼",
            password: "",
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := utils.ValidatePassword(tt.password)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}

// TestEmailValidation 電子郵件驗證測試
func TestEmailValidation(t *testing.T) {
    validEmails := []string{
        "test@example.com",
        "user.name@domain.co.uk",
        "test123@test-domain.org",
    }

    invalidEmails := []string{
        "invalid-email",
        "@example.com",
        "test@",
        "",
    }

    for _, email := range validEmails {
        t.Run("有效郵件: "+email, func(t *testing.T) {
            assert.True(t, utils.IsValidEmail(email))
        })
    }

    for _, email := range invalidEmails {
        t.Run("無效郵件: "+email, func(t *testing.T) {
            assert.False(t, utils.IsValidEmail(email))
        })
    }
}
```

## 整合測試

### 1. 資料庫整合測試

```go
// test/integration/database_test.go
package integration

import (
    "testing"
    "go-admin/common/database"
    "go-admin/app/admin/models"
    "github.com/stretchr/testify/suite"
    "gorm.io/gorm"
)

// DatabaseTestSuite 資料庫整合測試套件
type DatabaseTestSuite struct {
    suite.Suite
    db *gorm.DB
}

// SetupSuite 測試套件前置設定
func (suite *DatabaseTestSuite) SetupSuite() {
    // 初始化測試資料庫
    db, err := database.Setup("config/settings.test.yml")
    suite.Require().NoError(err)
    suite.db = db

    // 自動遷移測試表格
    err = suite.db.AutoMigrate(
        &models.SysUser{},
        &models.SysRole{},
        &models.SysApi{},
    )
    suite.Require().NoError(err)
}

// TearDownSuite 測試套件後置清理
func (suite *DatabaseTestSuite) TearDownSuite() {
    // 清理測試資料
    suite.db.Exec("DROP TABLE IF EXISTS sys_users")
    suite.db.Exec("DROP TABLE IF EXISTS sys_roles")
    suite.db.Exec("DROP TABLE IF EXISTS sys_apis")
}

// SetupTest 每個測試前置設定
func (suite *DatabaseTestSuite) SetupTest() {
    // 清空測試表格
    suite.db.Exec("DELETE FROM sys_users")
    suite.db.Exec("DELETE FROM sys_roles")
    suite.db.Exec("DELETE FROM sys_apis")
}

// TestUserCRUD 測試使用者 CRUD 操作
func (suite *DatabaseTestSuite) TestUserCRUD() {
    // Create - 建立使用者
    user := &models.SysUser{
        Username: "testuser",
        NickName: "測試使用者",
        Email:    "test@example.com",
    }

    result := suite.db.Create(user)
    suite.NoError(result.Error)
    suite.NotZero(user.Id)

    // Read - 讀取使用者
    var foundUser models.SysUser
    result = suite.db.First(&foundUser, user.Id)
    suite.NoError(result.Error)
    suite.Equal(user.Username, foundUser.Username)

    // Update - 更新使用者
    foundUser.NickName = "更新的使用者"
    result = suite.db.Save(&foundUser)
    suite.NoError(result.Error)

    // 驗證更新
    var updatedUser models.SysUser
    suite.db.First(&updatedUser, user.Id)
    suite.Equal("更新的使用者", updatedUser.NickName)

    // Delete - 刪除使用者
    result = suite.db.Delete(&foundUser)
    suite.NoError(result.Error)

    // 驗證刪除
    var deletedUser models.SysUser
    result = suite.db.First(&deletedUser, user.Id)
    suite.Error(result.Error)
    suite.True(errors.Is(result.Error, gorm.ErrRecordNotFound))
}

// TestUserRoleRelation 測試使用者角色關聯
func (suite *DatabaseTestSuite) TestUserRoleRelation() {
    // 建立角色
    role := &models.SysRole{
        RoleName: "測試角色",
        Status:   "2", // 啟用
    }
    suite.db.Create(role)

    // 建立使用者
    user := &models.SysUser{
        Username: "testuser",
        RoleId:   int(role.RoleId),
    }
    suite.db.Create(user)

    // 預載入關聯
    var userWithRole models.SysUser
    result := suite.db.Preload("Role").First(&userWithRole, user.Id)
    suite.NoError(result.Error)
    suite.NotNil(userWithRole.Role)
    suite.Equal(role.RoleName, userWithRole.Role.RoleName)
}

// 執行資料庫測試套件
func TestDatabaseTestSuite(t *testing.T) {
    suite.Run(t, new(DatabaseTestSuite))
}
```

### 2. 服務整合測試

```go
// test/integration/service_integration_test.go
package integration

import (
    "context"
    "testing"
    "go-admin/app/admin/service"
    "go-admin/common/database"
    "github.com/stretchr/testify/suite"
)

// ServiceIntegrationSuite 服務整合測試套件
type ServiceIntegrationSuite struct {
    suite.Suite
    userService *service.SysUser
    ctx         context.Context
}

// SetupSuite 設定測試套件
func (suite *ServiceIntegrationSuite) SetupSuite() {
    // 初始化資料庫
    db, err := database.Setup("config/settings.test.yml")
    suite.Require().NoError(err)

    // 初始化服務
    suite.userService = service.NewSysUserService(db)
    suite.ctx = context.Background()
}

// TestCreateAndGetUser 測試建立和取得使用者
func (suite *ServiceIntegrationSuite) TestCreateAndGetUser() {
    // 建立使用者資料
    userData := &service.CreateUserRequest{
        Username: "integrationtest",
        NickName: "整合測試使用者",
        Email:    "integration@test.com",
        Password: "testpassword123",
    }

    // 建立使用者
    userID, err := suite.userService.CreateUser(suite.ctx, userData)
    suite.NoError(err)
    suite.NotZero(userID)

    // 取得使用者
    user, err := suite.userService.GetUser(suite.ctx, userID)
    suite.NoError(err)
    suite.Equal(userData.Username, user.Username)
    suite.Equal(userData.NickName, user.NickName)
}

// 執行服務整合測試套件
func TestServiceIntegrationSuite(t *testing.T) {
    suite.Run(t, new(ServiceIntegrationSuite))
}
```

## API 測試

### 1. HTTP API 測試

```go
// test/api/user_api_test.go
package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "go-admin/cmd/app"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/suite"
)

// UserAPITestSuite 使用者 API 測試套件
type UserAPITestSuite struct {
    suite.Suite
    router   *gin.Engine
    server   *httptest.Server
    authToken string
}

// SetupSuite API 測試套件設定
func (suite *UserAPITestSuite) SetupSuite() {
    // 設定測試模式
    gin.SetMode(gin.TestMode)

    // 初始化路由
    suite.router = app.InitRouter("config/settings.test.yml")
    suite.server = httptest.NewServer(suite.router)

    // 取得認證權杖
    suite.authToken = suite.getAuthToken()
}

// TearDownSuite API 測試套件清理
func (suite *UserAPITestSuite) TearDownSuite() {
    suite.server.Close()
}

// getAuthToken 取得認證權杖
func (suite *UserAPITestSuite) getAuthToken() string {
    loginData := map[string]string{
        "username": "admin",
        "password": "admin123",
    }

    jsonData, _ := json.Marshal(loginData)
    req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    var response struct {
        Code int `json:"code"`
        Data struct {
            Token string `json:"token"`
        } `json:"data"`
    }

    json.Unmarshal(w.Body.Bytes(), &response)
    return response.Data.Token
}

// TestUserLogin 測試使用者登入
func (suite *UserAPITestSuite) TestUserLogin() {
    loginData := map[string]string{
        "username": "admin",
        "password": "admin123",
    }

    jsonData, _ := json.Marshal(loginData)
    req := httptest.NewRequest("POST", "/api/v1/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    suite.NoError(err)
    suite.Equal(float64(200), response["code"])
}

// TestCreateUser 測試建立使用者
func (suite *UserAPITestSuite) TestCreateUser() {
    userData := map[string]interface{}{
        "username": "newapiuser",
        "nickName": "新 API 使用者",
        "email":    "newapiuser@example.com",
        "password": "password123",
        "roleId":   1,
    }

    jsonData, _ := json.Marshal(userData)
    req := httptest.NewRequest("POST", "/api/v1/sysUser", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+suite.authToken)

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    suite.NoError(err)
    suite.Equal(float64(200), response["code"])
}

// TestGetUserList 測試取得使用者列表
func (suite *UserAPITestSuite) TestGetUserList() {
    req := httptest.NewRequest("GET", "/api/v1/sysUser?pageIndex=1&pageSize=10", nil)
    req.Header.Set("Authorization", "Bearer "+suite.authToken)

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    suite.NoError(err)
    suite.Equal(float64(200), response["code"])

    // 檢查分頁資料
    data := response["data"].(map[string]interface{})
    suite.Contains(data, "list")
    suite.Contains(data, "total")
}

// TestUpdateUser 測試更新使用者
func (suite *UserAPITestSuite) TestUpdateUser() {
    // 首先建立一個使用者
    createData := map[string]interface{}{
        "username": "updatetest",
        "nickName": "待更新使用者",
        "email":    "update@test.com",
        "password": "password123",
        "roleId":   1,
    }

    jsonData, _ := json.Marshal(createData)
    req := httptest.NewRequest("POST", "/api/v1/sysUser", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+suite.authToken)

    w := httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    // 取得建立的使用者 ID
    var createResponse map[string]interface{}
    json.Unmarshal(w.Body.Bytes(), &createResponse)
    userID := createResponse["data"].(map[string]interface{})["id"]

    // 更新使用者
    updateData := map[string]interface{}{
        "id":       userID,
        "nickName": "已更新使用者",
        "email":    "updated@test.com",
    }

    jsonData, _ = json.Marshal(updateData)
    req = httptest.NewRequest("PUT", "/api/v1/sysUser", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+suite.authToken)

    w = httptest.NewRecorder()
    suite.router.ServeHTTP(w, req)

    suite.Equal(http.StatusOK, w.Code)
}

// 執行 API 測試套件
func TestUserAPITestSuite(t *testing.T) {
    suite.Run(t, new(UserAPITestSuite))
}
```

## 效能測試

### 1. 基準測試

```go
// test/benchmark/user_benchmark_test.go
package benchmark

import (
    "context"
    "testing"
    "go-admin/app/admin/service"
    "go-admin/common/database"
)

var userService *service.SysUser

// 初始化基準測試
func init() {
    db, _ := database.Setup("config/settings.test.yml")
    userService = service.NewSysUserService(db)
}

// BenchmarkUserCreate 建立使用者基準測試
func BenchmarkUserCreate(b *testing.B) {
    ctx := context.Background()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        userData := &service.CreateUserRequest{
            Username: fmt.Sprintf("benchuser%d", i),
            NickName: "基準測試使用者",
            Email:    fmt.Sprintf("bench%d@test.com", i),
            Password: "password123",
        }

        _, err := userService.CreateUser(ctx, userData)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// BenchmarkUserGet 取得使用者基準測試
func BenchmarkUserGet(b *testing.B) {
    ctx := context.Background()

    // 準備測試資料
    userData := &service.CreateUserRequest{
        Username: "benchgetuser",
        NickName: "基準測試使用者",
        Email:    "benchget@test.com",
        Password: "password123",
    }

    userID, _ := userService.CreateUser(ctx, userData)

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := userService.GetUser(ctx, userID)
        if err != nil {
            b.Fatal(err)
        }
    }
}

// BenchmarkPasswordHash 密碼雜湊基準測試
func BenchmarkPasswordHash(b *testing.B) {
    password := "testpassword123"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
        if err != nil {
            b.Fatal(err)
        }
    }
}
```

### 2. 負載測試

```go
// test/load/api_load_test.go
package load

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "sync"
    "testing"
    "time"
)

// LoadTestConfig 負載測試配置
type LoadTestConfig struct {
    BaseURL     string
    Concurrency int
    Duration    time.Duration
    AuthToken   string
}

// APILoadTest API 負載測試
func TestAPILoadTest(t *testing.T) {
    config := &LoadTestConfig{
        BaseURL:     "http://localhost:8080",
        Concurrency: 10,
        Duration:    30 * time.Second,
        AuthToken:   getTestToken(),
    }

    // 執行負載測試
    results := runLoadTest(config, testUserListAPI)

    // 分析結果
    t.Logf("總請求數: %d", results.TotalRequests)
    t.Logf("成功請求: %d", results.SuccessRequests)
    t.Logf("失敗請求: %d", results.FailedRequests)
    t.Logf("平均回應時間: %v", results.AvgResponseTime)
    t.Logf("QPS: %.2f", results.QPS)

    // 檢查效能指標
    if results.AvgResponseTime > 500*time.Millisecond {
        t.Errorf("平均回應時間過長: %v", results.AvgResponseTime)
    }

    if results.QPS < 50 {
        t.Errorf("QPS 過低: %.2f", results.QPS)
    }
}

// LoadTestResult 負載測試結果
type LoadTestResult struct {
    TotalRequests     int
    SuccessRequests   int
    FailedRequests    int
    AvgResponseTime   time.Duration
    QPS              float64
}

// runLoadTest 執行負載測試
func runLoadTest(config *LoadTestConfig, testFunc func(string, string) error) *LoadTestResult {
    var wg sync.WaitGroup
    var mu sync.Mutex

    result := &LoadTestResult{}
    startTime := time.Now()
    endTime := startTime.Add(config.Duration)

    // 啟動並發測試
    for i := 0; i < config.Concurrency; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()

            for time.Now().Before(endTime) {
                reqStartTime := time.Now()
                err := testFunc(config.BaseURL, config.AuthToken)
                reqDuration := time.Since(reqStartTime)

                mu.Lock()
                result.TotalRequests++
                if err != nil {
                    result.FailedRequests++
                } else {
                    result.SuccessRequests++
                }
                result.AvgResponseTime = (result.AvgResponseTime*time.Duration(result.TotalRequests-1) + reqDuration) / time.Duration(result.TotalRequests)
                mu.Unlock()
            }
        }()
    }

    wg.Wait()

    actualDuration := time.Since(startTime)
    result.QPS = float64(result.TotalRequests) / actualDuration.Seconds()

    return result
}

// testUserListAPI 測試使用者列表 API
func testUserListAPI(baseURL, authToken string) error {
    url := fmt.Sprintf("%s/api/v1/sysUser?pageIndex=1&pageSize=10", baseURL)

    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return err
    }

    req.Header.Set("Authorization", "Bearer "+authToken)

    client := &http.Client{Timeout: 5 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    }

    return nil
}

// getTestToken 取得測試權杖
func getTestToken() string {
    // 實作取得測試權杖的邏輯
    // 這裡簡化為直接回傳固定權杖
    return "test-auth-token"
}
```

## 測試工具與框架

### 1. 推薦測試框架

- **Testify**: 豐富的斷言和模擬功能
- **Ginkgo & Gomega**: BDD 風格測試框架
- **GoMock**: 模擬物件產生器
- **GoConvey**: Web UI 測試框架

### 2. 安裝和設定

```bash
# 安裝 Testify
go get github.com/stretchr/testify

# 安裝 Ginkgo
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go get github.com/onsi/gomega/...

# 安裝 GoMock
go install github.com/golang/mock/mockgen@latest

# 安裝 GoConvey
go get github.com/smartystreets/goconvey/convey
```

### 3. 測試輔助工具

```go
// test/helpers/test_helpers.go
package helpers

import (
    "testing"
    "go-admin/common/database"
    "gorm.io/gorm"
)

// TestDatabase 測試資料庫輔助
type TestDatabase struct {
    DB *gorm.DB
}

// NewTestDatabase 建立測試資料庫
func NewTestDatabase(t *testing.T) *TestDatabase {
    db, err := database.Setup("config/settings.test.yml")
    if err != nil {
        t.Fatalf("無法建立測試資料庫: %v", err)
    }

    return &TestDatabase{DB: db}
}

// Cleanup 清理測試資料
func (td *TestDatabase) Cleanup() {
    td.DB.Exec("DELETE FROM sys_users")
    td.DB.Exec("DELETE FROM sys_roles")
    td.DB.Exec("DELETE FROM sys_apis")
}

// SeedTestData 填入測試資料
func (td *TestDatabase) SeedTestData() {
    // 建立測試角色
    role := &models.SysRole{
        RoleName: "測試角色",
        Status:   "2",
    }
    td.DB.Create(role)

    // 建立測試使用者
    user := &models.SysUser{
        Username: "testuser",
        NickName: "測試使用者",
        Email:    "test@example.com",
        RoleId:   int(role.RoleId),
    }
    td.DB.Create(user)
}

// AssertUserExists 斷言使用者存在
func (td *TestDatabase) AssertUserExists(t *testing.T, username string) {
    var user models.SysUser
    result := td.DB.Where("username = ?", username).First(&user)
    if result.Error != nil {
        t.Errorf("使用者 %s 不存在: %v", username, result.Error)
    }
}
```

## 測試最佳實務

### 1. 測試命名規範

```go
// 好的測試命名
func TestUserService_CreateUser_Success(t *testing.T) {}
func TestUserService_CreateUser_WithInvalidEmail_ReturnsError(t *testing.T) {}
func TestUserAPI_Login_WithValidCredentials_ReturnsToken(t *testing.T) {}

// 不好的測試命名
func TestUser(t *testing.T) {}
func Test1(t *testing.T) {}
func TestCreateUser(t *testing.T) {}
```

### 2. 測試結構組織

```
test/
├── unit/           # 單元測試
│   ├── models/     # 模型測試
│   ├── services/   # 服務測試
│   └── utils/      # 工具測試
├── integration/    # 整合測試
│   ├── database/   # 資料庫測試
│   └── services/   # 服務整合測試
├── api/           # API 測試
│   └── v1/        # API v1 測試
├── benchmark/     # 效能測試
├── load/          # 負載測試
├── helpers/       # 測試輔助工具
└── fixtures/      # 測試資料固件
```

### 3. 測試資料管理

```go
// test/fixtures/user_fixtures.go
package fixtures

import "go-admin/app/admin/models"

// UserFixtures 使用者測試資料固件
var UserFixtures = struct {
    ValidUser       *models.SysUser
    AdminUser       *models.SysUser
    InactiveUser    *models.SysUser
}{
    ValidUser: &models.SysUser{
        Username: "validuser",
        NickName: "有效使用者",
        Email:    "valid@example.com",
        Status:   "2",
    },
    AdminUser: &models.SysUser{
        Username: "admin",
        NickName: "管理員",
        Email:    "admin@example.com",
        Status:   "2",
        RoleId:   1,
    },
    InactiveUser: &models.SysUser{
        Username: "inactive",
        NickName: "停用使用者",
        Email:    "inactive@example.com",
        Status:   "1",
    },
}
```

### 4. 錯誤測試

```go
// 測試錯誤情況
func TestUserService_CreateUser_WithDuplicateUsername_ReturnsError(t *testing.T) {
    // 準備: 建立第一個使用者
    user1 := &models.SysUser{Username: "duplicate"}
    err := userService.CreateUser(ctx, user1)
    assert.NoError(t, err)

    // 執行: 嘗試建立相同使用者名稱的使用者
    user2 := &models.SysUser{Username: "duplicate"}
    err = userService.CreateUser(ctx, user2)

    // 驗證: 應該回傳錯誤
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "使用者名稱已存在")
}
```

## 持續整合測試

### 1. GitHub Actions 設定

建立 `.github/workflows/test.yml`：

```yaml
name: 測試

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: root
          MYSQL_DATABASE: go_admin_test
        ports:
          - 3306:3306
        options: >-
          --health-cmd="mysqladmin ping"
          --health-interval=10s
          --health-timeout=5s
          --health-retries=3

    steps:
      - uses: actions/checkout@v3

      - name: 設定 Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.21

      - name: 快取模組
        uses: actions/cache@v3
        with:
          path: ~/go/pkg/mod
          key: ${{ runner.os }}-go-${{ hashFiles('**/go.sum') }}
          restore-keys: |
            ${{ runner.os }}-go-

      - name: 下載依賴
        run: go mod download

      - name: 執行單元測試
        run: go test -v ./test/unit/...

      - name: 執行整合測試
        run: go test -v ./test/integration/...
        env:
          DB_TYPE: mysql
          DB_HOST: 127.0.0.1
          DB_PORT: 3306
          DB_USERNAME: root
          DB_PASSWORD: root
          DB_NAME: go_admin_test

      - name: 執行 API 測試
        run: go test -v ./test/api/...

      - name: 產生測試覆蓋率報告
        run: go test -coverprofile=coverage.out ./...

      - name: 上傳覆蓋率到 Codecov
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
```

### 2. 測試指令腳本

建立 `scripts/test.sh`：

```bash
#!/bin/bash

# 測試執行腳本
set -e

echo "🧪 開始執行測試..."

# 設定測試環境
export APP_MODE=test
export DB_TYPE=sqlite3
export DB_PATH=test.db
export LOG_LEVEL=error

# 清理之前的測試資料
rm -f test.db
rm -rf test_logs

echo "📋 執行單元測試..."
go test -v -timeout 30s ./test/unit/...

echo "🔗 執行整合測試..."
go test -v -timeout 60s ./test/integration/...

echo "🌐 執行 API 測試..."
go test -v -timeout 120s ./test/api/...

echo "⚡ 執行效能測試..."
go test -v -bench=. -benchmem ./test/benchmark/...

echo "📊 產生測試覆蓋率報告..."
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

echo "✅ 所有測試執行完成!"
echo "📈 測試覆蓋率報告已產生: coverage.html"
```

### 3. 測試品質檢查

```bash
#!/bin/bash
# scripts/test_quality.sh

# 檢查測試覆蓋率
COVERAGE=$(go test -coverprofile=coverage.out ./... | grep "coverage:" | awk '{print $3}' | sed 's/%//')

if (( $(echo "$COVERAGE < 80" | bc -l) )); then
    echo "❌ 測試覆蓋率過低: $COVERAGE% (最低要求: 80%)"
    exit 1
else
    echo "✅ 測試覆蓋率達標: $COVERAGE%"
fi

# 檢查測試執行時間
SLOW_TESTS=$(go test -v ./... | grep -E "PASS.*[0-9]+\.[0-9]+s" | awk '$3 > 1.0 {print $2, $3}')

if [ ! -z "$SLOW_TESTS" ]; then
    echo "⚠️  發現執行時間過長的測試:"
    echo "$SLOW_TESTS"
fi

echo "✅ 測試品質檢查完成"
```

---

## 執行測試

### 基本測試執行

```bash
# 執行所有測試
go test ./...

# 執行特定套件測試
go test ./test/unit/...
go test ./test/integration/...
go test ./test/api/...

# 執行單個測試檔案
go test ./test/unit/user_test.go

# 執行特定測試函數
go test -run TestUserValidation ./test/unit/...

# 詳細輸出
go test -v ./...

# 產生覆蓋率報告
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 效能測試執行

```bash
# 執行基準測試
go test -bench=. ./test/benchmark/...

# 執行記憶體分析
go test -bench=. -benchmem ./test/benchmark/...

# 指定執行時間
go test -bench=. -benchtime=10s ./test/benchmark/...
```

通過完整的測試策略，Go-Admin 專案能夠確保程式碼品質、功能正確性和系統穩定性。建議開發者在開發過程中持續撰寫和執行測試，以維持高品質的程式碼標準。
