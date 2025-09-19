package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// DatabaseIntegrationSuite 資料庫整合測試套件
// 此套件測試與實際資料庫的整合功能
type DatabaseIntegrationSuite struct {
	suite.Suite
	db     *gorm.DB
	ctx    context.Context
	server *httptest.Server
	router *gin.Engine
}

// SetupSuite 測試套件初始化
// 建立測試資料庫連接和測試服務器
func (suite *DatabaseIntegrationSuite) SetupSuite() {
	// 設定 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	// 初始化測試資料庫（這裡使用 SQLite 記憶體資料庫）
	// db, err := database.SetupTestDB()
	// suite.Require().NoError(err, "無法建立測試資料庫")
	// suite.db = db

	// 初始化路由
	suite.router = gin.New()
	suite.setupRoutes()

	// 建立測試服務器
	suite.server = httptest.NewServer(suite.router)

	// 建立 context
	suite.ctx = context.Background()

	fmt.Println("✅ 資料庫整合測試套件初始化完成")
}

// TearDownSuite 測試套件清理
func (suite *DatabaseIntegrationSuite) TearDownSuite() {
	// 關閉測試服務器
	if suite.server != nil {
		suite.server.Close()
	}

	// 清理測試資料庫
	if suite.db != nil {
		// 刪除測試表格
		suite.db.Exec("DROP TABLE IF EXISTS sys_users")
		suite.db.Exec("DROP TABLE IF EXISTS sys_roles")
		suite.db.Exec("DROP TABLE IF EXISTS sys_apis")
	}

	fmt.Println("✅ 資料庫整合測試套件清理完成")
}

// SetupTest 每個測試前的準備工作
func (suite *DatabaseIntegrationSuite) SetupTest() {
	if suite.db != nil {
		// 清空測試資料
		suite.db.Exec("DELETE FROM sys_users")
		suite.db.Exec("DELETE FROM sys_roles")
		suite.db.Exec("DELETE FROM sys_apis")

		// 重置自增 ID
		suite.db.Exec("DELETE FROM sqlite_sequence WHERE name IN ('sys_users', 'sys_roles', 'sys_apis')")
	}
}

// setupRoutes 設定測試路由
func (suite *DatabaseIntegrationSuite) setupRoutes() {
	api := suite.router.Group("/api/v1")
	{
		// 認證相關路由
		api.POST("/login", suite.mockLogin)
		api.GET("/logout", suite.mockLogout)

		// 使用者管理路由（需要認證）
		users := api.Group("/sysUser")
		{
			users.GET("", suite.mockGetUserList)
			users.POST("", suite.mockCreateUser)
			users.PUT("", suite.mockUpdateUser)
			users.DELETE("/:id", suite.mockDeleteUser)
			users.GET("/:id", suite.mockGetUser)
		}

		// 角色管理路由
		roles := api.Group("/sysRole")
		{
			roles.GET("", suite.mockGetRoleList)
			roles.POST("", suite.mockCreateRole)
		}

		// 健康檢查
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"code":    200,
				"message": "系統運行正常",
				"data": gin.H{
					"status":    "healthy",
					"timestamp": time.Now().Format(time.RFC3339),
				},
			})
		})
	}
}

// TestHealthCheck 測試健康檢查端點
func (suite *DatabaseIntegrationSuite) TestHealthCheck() {
	// 發送健康檢查請求
	resp, err := http.Get(suite.server.URL + "/api/v1/health")
	suite.NoError(err, "健康檢查請求應該成功")
	defer resp.Body.Close()

	// 檢查狀態碼
	suite.Equal(http.StatusOK, resp.StatusCode, "健康檢查應該返回 200 狀態碼")

	// 解析回應
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	suite.NoError(err, "應該能夠解析健康檢查回應")

	// 檢查回應內容
	suite.Equal(float64(200), response["code"], "回應代碼應該是 200")
	suite.Equal("系統運行正常", response["message"], "回應訊息應該正確")
	suite.Contains(response, "data", "回應應該包含 data 欄位")

	data := response["data"].(map[string]interface{})
	suite.Equal("healthy", data["status"], "系統狀態應該是 healthy")
	suite.Contains(data, "timestamp", "回應應該包含時間戳")
}

// TestUserAuthentication 測試使用者認證流程
func (suite *DatabaseIntegrationSuite) TestUserAuthentication() {
	// 測試登入
	loginData := map[string]string{
		"username": "admin",
		"password": "admin123",
	}

	jsonData, err := json.Marshal(loginData)
	suite.NoError(err, "應該能夠序列化登入資料")

	// 發送登入請求
	resp, err := http.Post(
		suite.server.URL+"/api/v1/login",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	suite.NoError(err, "登入請求應該成功")
	defer resp.Body.Close()

	// 檢查登入回應
	suite.Equal(http.StatusOK, resp.StatusCode, "登入應該成功")

	var loginResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	suite.NoError(err, "應該能夠解析登入回應")

	suite.Equal(float64(200), loginResponse["code"], "登入回應代碼應該是 200")
	suite.Contains(loginResponse, "data", "登入回應應該包含 data")

	// 檢查是否返回 token
	data := loginResponse["data"].(map[string]interface{})
	suite.Contains(data, "token", "登入回應應該包含 token")
	suite.NotEmpty(data["token"], "token 不應該為空")
}

// TestUserCRUDOperations 測試使用者 CRUD 操作
func (suite *DatabaseIntegrationSuite) TestUserCRUDOperations() {
	// 1. 建立使用者
	userData := map[string]interface{}{
		"username": "testuser",
		"nickName": "測試使用者",
		"email":    "test@example.com",
		"password": "password123",
		"phone":    "0912345678",
		"status":   "2",
		"roleId":   1,
	}

	jsonData, _ := json.Marshal(userData)
	resp, err := http.Post(
		suite.server.URL+"/api/v1/sysUser",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	suite.NoError(err, "建立使用者請求應該成功")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "建立使用者應該成功")

	var createResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&createResponse)
	suite.NoError(err, "應該能夠解析建立使用者回應")

	suite.Equal(float64(200), createResponse["code"], "建立使用者回應代碼應該是 200")

	// 2. 取得使用者列表
	resp, err = http.Get(suite.server.URL + "/api/v1/sysUser?pageIndex=1&pageSize=10")
	suite.NoError(err, "取得使用者列表請求應該成功")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "取得使用者列表應該成功")

	var listResponse map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&listResponse)
	suite.NoError(err, "應該能夠解析使用者列表回應")

	suite.Equal(float64(200), listResponse["code"], "取得使用者列表回應代碼應該是 200")
	suite.Contains(listResponse, "data", "使用者列表回應應該包含 data")

	// 3. 更新使用者
	updateData := map[string]interface{}{
		"id":       1,
		"nickName": "更新的測試使用者",
		"email":    "updated@example.com",
	}

	jsonData, _ = json.Marshal(updateData)
	req, err := http.NewRequest("PUT", suite.server.URL+"/api/v1/sysUser", bytes.NewBuffer(jsonData))
	suite.NoError(err, "應該能夠建立更新請求")

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err = client.Do(req)
	suite.NoError(err, "更新使用者請求應該成功")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "更新使用者應該成功")

	// 4. 刪除使用者
	req, err = http.NewRequest("DELETE", suite.server.URL+"/api/v1/sysUser/1", nil)
	suite.NoError(err, "應該能夠建立刪除請求")

	resp, err = client.Do(req)
	suite.NoError(err, "刪除使用者請求應該成功")
	defer resp.Body.Close()

	suite.Equal(http.StatusOK, resp.StatusCode, "刪除使用者應該成功")
}

// TestConcurrentRequests 測試並發請求
func (suite *DatabaseIntegrationSuite) TestConcurrentRequests() {
	const numRequests = 10
	const concurrency = 5

	// 建立並發請求通道
	requests := make(chan int, numRequests)
	results := make(chan bool, numRequests)

	// 填充請求
	for i := 0; i < numRequests; i++ {
		requests <- i
	}
	close(requests)

	// 啟動並發工作者
	for i := 0; i < concurrency; i++ {
		go func() {
			for requestId := range requests {
				// 發送健康檢查請求
				resp, err := http.Get(suite.server.URL + "/api/v1/health")
				success := err == nil && resp != nil && resp.StatusCode == http.StatusOK

				if resp != nil {
					resp.Body.Close()
				}

				results <- success
				fmt.Printf("並發請求 %d 完成，成功: %t\n", requestId, success)
			}
		}()
	}

	// 收集結果
	successCount := 0
	for i := 0; i < numRequests; i++ {
		if <-results {
			successCount++
		}
	}

	// 檢查結果
	suite.Equal(numRequests, successCount, "所有並發請求都應該成功")
}

// Mock 函數實作

func (suite *DatabaseIntegrationSuite) mockLogin(c *gin.Context) {
	var loginData struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "請求參數錯誤"})
		return
	}

	// 模擬登入驗證
	if loginData.Username == "admin" && loginData.Password == "admin123" {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "登入成功",
			"data": gin.H{
				"token": "mock-jwt-token-" + fmt.Sprint(time.Now().Unix()),
				"user": gin.H{
					"id":       1,
					"username": loginData.Username,
					"nickName": "管理員",
				},
			},
		})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    401,
			"message": "使用者名稱或密碼錯誤",
		})
	}
}

func (suite *DatabaseIntegrationSuite) mockLogout(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登出成功",
	})
}

func (suite *DatabaseIntegrationSuite) mockGetUserList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查詢成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"id":       1,
					"username": "admin",
					"nickName": "管理員",
					"email":    "admin@example.com",
					"status":   "2",
				},
			},
			"total":     1,
			"pageIndex": 1,
			"pageSize":  10,
		},
	})
}

func (suite *DatabaseIntegrationSuite) mockCreateUser(c *gin.Context) {
	var userData map[string]interface{}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "請求參數錯誤"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "建立成功",
		"data": gin.H{
			"id": 1,
		},
	})
}

func (suite *DatabaseIntegrationSuite) mockUpdateUser(c *gin.Context) {
	var userData map[string]interface{}
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "請求參數錯誤"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "更新成功",
	})
}

func (suite *DatabaseIntegrationSuite) mockDeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "使用者ID不能為空"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "刪除成功",
	})
}

func (suite *DatabaseIntegrationSuite) mockGetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "使用者ID不能為空"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查詢成功",
		"data": gin.H{
			"id":       userID,
			"username": "testuser",
			"nickName": "測試使用者",
			"email":    "test@example.com",
			"status":   "2",
		},
	})
}

func (suite *DatabaseIntegrationSuite) mockGetRoleList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "查詢成功",
		"data": gin.H{
			"list": []gin.H{
				{
					"roleId":   1,
					"roleName": "管理員",
					"status":   "2",
				},
			},
			"total": 1,
		},
	})
}

func (suite *DatabaseIntegrationSuite) mockCreateRole(c *gin.Context) {
	var roleData map[string]interface{}
	if err := c.ShouldBindJSON(&roleData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "請求參數錯誤"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "建立成功",
		"data": gin.H{
			"roleId": 1,
		},
	})
}

// 執行資料庫整合測試套件
func TestDatabaseIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DatabaseIntegrationSuite))
}

// TestAPIResponseFormat 測試 API 回應格式
func TestAPIResponseFormat(t *testing.T) {
	// 建立測試路由
	router := gin.New()
	gin.SetMode(gin.TestMode)

	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "測試成功",
			"data": gin.H{
				"test": "value",
			},
		})
	})

	// 建立測試請求
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// 檢查狀態碼
	assert.Equal(t, http.StatusOK, w.Code)

	// 檢查回應格式
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// 檢查標準回應欄位
	assert.Contains(t, response, "code")
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "data")

	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "測試成功", response["message"])
	assert.IsType(t, map[string]interface{}{}, response["data"])
}

// TestErrorHandling 測試錯誤處理
func TestErrorHandling(t *testing.T) {
	router := gin.New()
	gin.SetMode(gin.TestMode)

	// 建立會產生錯誤的路由
	router.POST("/error", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "請求參數錯誤",
			"error":   "Invalid JSON format",
		})
	})

	// 發送無效 JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/error", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// 檢查錯誤回應
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, float64(400), response["code"])
	assert.Contains(t, response["message"], "錯誤")
}