package fixtures
package fixtures

import "time"

// UserFixtures 使用者測試資料固件
// 提供各種測試場景的使用者資料
var UserFixtures = struct {
	ValidUser       map[string]interface{}
	AdminUser       map[string]interface{}
	InactiveUser    map[string]interface{}
	InvalidUser     map[string]interface{}
	MinimalUser     map[string]interface{}
}{
	// 有效的一般使用者
	ValidUser: map[string]interface{}{
		"id":        1,
		"username":  "validuser",
		"nickName":  "有效使用者",
		"email":     "valid@example.com",
		"phone":     "0912345678",
		"status":    "2", // 啟用
		"roleId":    2,   // 一般使用者角色
		"avatar":    "",
		"sex":       "1", // 男性
		"remark":    "測試用的有效使用者",
		"createdAt": time.Now().Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 管理員使用者
	AdminUser: map[string]interface{}{
		"id":        2,
		"username":  "admin",
		"nickName":  "系統管理員",
		"email":     "admin@example.com",
		"phone":     "0987654321",
		"status":    "2", // 啟用
		"roleId":    1,   // 管理員角色
		"avatar":    "",
		"sex":       "1", // 男性
		"remark":    "系統管理員帳號",
		"createdAt": time.Now().AddDate(0, 0, -30).Format(time.RFC3339), // 30天前建立
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 停用的使用者
	InactiveUser: map[string]interface{}{
		"id":        3,
		"username":  "inactiveuser",
		"nickName":  "停用使用者",
		"email":     "inactive@example.com",
		"phone":     "0911111111",
		"status":    "1", // 停用
		"roleId":    2,   // 一般使用者角色
		"avatar":    "",
		"sex":       "2", // 女性
		"remark":    "已停用的使用者帳號",
		"createdAt": time.Now().AddDate(0, 0, -15).Format(time.RFC3339),
		"updatedAt": time.Now().AddDate(0, 0, -5).Format(time.RFC3339),
	},

	// 無效的使用者（用於測試驗證）
	InvalidUser: map[string]interface{}{
		"id":       4,
		"username": "", // 空的使用者名稱（無效）
		"nickName": "無效使用者",
		"email":    "invalid-email", // 無效的電子郵件格式
		"phone":    "123",           // 無效的電話號碼
		"status":   "3",             // 無效的狀態值
		"roleId":   999,             // 不存在的角色ID
	},

	// 最小必要欄位的使用者
	MinimalUser: map[string]interface{}{
		"username": "minimaluser",
		"email":    "minimal@example.com",
		"status":   "2",
		"roleId":   2,
	},
}

// RoleFixtures 角色測試資料固件
var RoleFixtures = struct {
	AdminRole    map[string]interface{}
	UserRole     map[string]interface{}
	ManagerRole  map[string]interface{}
	GuestRole    map[string]interface{}
}{
	// 管理員角色
	AdminRole: map[string]interface{}{
		"roleId":     1,
		"roleName":   "管理員",
		"status":     "2", // 啟用
		"roleKey":    "admin",
		"roleSort":   1,
		"flag":       "",
		"remark":     "系統管理員角色，擁有所有權限",
		"admin":      true,
		"dataScope":  "1", // 全部資料權限
		"createdAt":  time.Now().AddDate(0, 0, -365).Format(time.RFC3339), // 一年前建立
		"updatedAt":  time.Now().Format(time.RFC3339),
	},

	// 一般使用者角色
	UserRole: map[string]interface{}{
		"roleId":     2,
		"roleName":   "一般使用者",
		"status":     "2", // 啟用
		"roleKey":    "user",
		"roleSort":   3,
		"flag":       "",
		"remark":     "一般使用者角色，基本權限",
		"admin":      false,
		"dataScope":  "5", // 僅本人資料權限
		"createdAt":  time.Now().AddDate(0, 0, -300).Format(time.RFC3339),
		"updatedAt":  time.Now().Format(time.RFC3339),
	},

	// 管理者角色
	ManagerRole: map[string]interface{}{
		"roleId":     3,
		"roleName":   "部門主管",
		"status":     "2", // 啟用
		"roleKey":    "manager",
		"roleSort":   2,
		"flag":       "",
		"remark":     "部門主管角色，管理部門資料",
		"admin":      false,
		"dataScope":  "4", // 部門資料權限
		"createdAt":  time.Now().AddDate(0, 0, -200).Format(time.RFC3339),
		"updatedAt":  time.Now().Format(time.RFC3339),
	},

	// 訪客角色
	GuestRole: map[string]interface{}{
		"roleId":     4,
		"roleName":   "訪客",
		"status":     "1", // 停用
		"roleKey":    "guest",
		"roleSort":   4,
		"flag":       "",
		"remark":     "訪客角色，唯讀權限",
		"admin":      false,
		"dataScope":  "5", // 僅本人資料權限
		"createdAt":  time.Now().AddDate(0, 0, -100).Format(time.RFC3339),
		"updatedAt":  time.Now().AddDate(0, 0, -30).Format(time.RFC3339),
	},
}

// APIFixtures API 測試資料固件
var APIFixtures = struct {
	LoginAPI       map[string]interface{}
	UserListAPI    map[string]interface{}
	UserCreateAPI  map[string]interface{}
	UserUpdateAPI  map[string]interface{}
	UserDeleteAPI  map[string]interface{}
}{
	// 登入 API
	LoginAPI: map[string]interface{}{
		"id":       1,
		"handle":   "登入",
		"title":    "使用者登入",
		"path":     "/api/v1/login",
		"type":     "",
		"action":   "POST",
		"createdAt": time.Now().AddDate(0, 0, -365).Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 使用者列表 API
	UserListAPI: map[string]interface{}{
		"id":       2,
		"handle":   "GetSysUserList",
		"title":    "取得使用者列表",
		"path":     "/api/v1/sysUser",
		"type":     "",
		"action":   "GET",
		"createdAt": time.Now().AddDate(0, 0, -365).Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 建立使用者 API
	UserCreateAPI: map[string]interface{}{
		"id":       3,
		"handle":   "CreateSysUser",
		"title":    "建立使用者",
		"path":     "/api/v1/sysUser",
		"type":     "",
		"action":   "POST",
		"createdAt": time.Now().AddDate(0, 0, -365).Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 更新使用者 API
	UserUpdateAPI: map[string]interface{}{
		"id":       4,
		"handle":   "UpdateSysUser",
		"title":    "更新使用者",
		"path":     "/api/v1/sysUser",
		"type":     "",
		"action":   "PUT",
		"createdAt": time.Now().AddDate(0, 0, -365).Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},

	// 刪除使用者 API
	UserDeleteAPI: map[string]interface{}{
		"id":       5,
		"handle":   "DeleteSysUser",
		"title":    "刪除使用者",
		"path":     "/api/v1/sysUser",
		"type":     "",
		"action":   "DELETE",
		"createdAt": time.Now().AddDate(0, 0, -365).Format(time.RFC3339),
		"updatedAt": time.Now().Format(time.RFC3339),
	},
}

// AuthFixtures 認證測試資料固件
var AuthFixtures = struct {
	ValidLogin    map[string]interface{}
	InvalidLogin  map[string]interface{}
	EmptyLogin    map[string]interface{}
	AdminLogin    map[string]interface{}
}{
	// 有效登入資料
	ValidLogin: map[string]interface{}{
		"username": "validuser",
		"password": "password123",
	},

	// 無效登入資料
	InvalidLogin: map[string]interface{}{
		"username": "invaliduser",
		"password": "wrongpassword",
	},

	// 空的登入資料
	EmptyLogin: map[string]interface{}{
		"username": "",
		"password": "",
	},

	// 管理員登入資料
	AdminLogin: map[string]interface{}{
		"username": "admin",
		"password": "admin123",
	},
}

// ResponseFixtures 回應測試資料固件
var ResponseFixtures = struct {
	SuccessResponse map[string]interface{}
	ErrorResponse   map[string]interface{}
	ValidationError map[string]interface{}
	NotFoundError   map[string]interface{}
}{
	// 成功回應
	SuccessResponse: map[string]interface{}{
		"code":    200,
		"message": "操作成功",
		"data": map[string]interface{}{
			"id": 1,
		},
	},

	// 錯誤回應
	ErrorResponse: map[string]interface{}{
		"code":    500,
		"message": "系統內部錯誤",
		"error":   "Database connection failed",
	},

	// 驗證錯誤
	ValidationError: map[string]interface{}{
		"code":    400,
		"message": "請求參數錯誤",
		"errors": map[string]interface{}{
			"username": []string{"使用者名稱不能為空"},
			"email":    []string{"電子郵件格式不正確"},
		},
	},

	// 找不到資源
	NotFoundError: map[string]interface{}{
		"code":    404,
		"message": "資源不存在",
	},
}

// TestCases 測試案例集合
var TestCases = struct {
	UserValidation []TestCase
	AuthTests      []TestCase
	APITests       []TestCase
}{
	// 使用者驗證測試案例
	UserValidation: []TestCase{
		{
			Name:        "有效使用者資料",
			Input:       UserFixtures.ValidUser,
			Expected:    true,
			Description: "測試有效的使用者資料應該通過驗證",
		},
		{
			Name:        "無效使用者資料",
			Input:       UserFixtures.InvalidUser,
			Expected:    false,
			Description: "測試無效的使用者資料應該驗證失敗",
		},
		{
			Name:        "最小必要欄位",
			Input:       UserFixtures.MinimalUser,
			Expected:    true,
			Description: "測試只包含最小必要欄位的使用者資料",
		},
	},

	// 認證測試案例
	AuthTests: []TestCase{
		{
			Name:        "有效登入",
			Input:       AuthFixtures.ValidLogin,
			Expected:    true,
			Description: "測試有效的登入憑證",
		},
		{
			Name:        "無效登入",
			Input:       AuthFixtures.InvalidLogin,
			Expected:    false,
			Description: "測試無效的登入憑證",
		},
		{
			Name:        "空的登入資料",
			Input:       AuthFixtures.EmptyLogin,
			Expected:    false,
			Description: "測試空的登入資料",
		},
	},

	// API 測試案例
	APITests: []TestCase{
		{
			Name:        "取得使用者列表",
			Input:       map[string]interface{}{"pageIndex": 1, "pageSize": 10},
			Expected:    ResponseFixtures.SuccessResponse,
			Description: "測試取得使用者列表 API",
		},
		{
			Name:        "建立使用者",
			Input:       UserFixtures.ValidUser,
			Expected:    ResponseFixtures.SuccessResponse,
			Description: "測試建立使用者 API",
		},
	},
}

// TestCase 測試案例結構
type TestCase struct {
	Name        string      `json:"name"`
	Input       interface{} `json:"input"`
	Expected    interface{} `json:"expected"`
	Description string      `json:"description"`
}

// DatabaseFixtures 資料庫測試資料固件
var DatabaseFixtures = struct {
	CreateTables []string
	SeedData     []map[string]interface{}
	CleanupSQL   []string
}{
	// 建立測試表格 SQL
	CreateTables: []string{
		`CREATE TABLE IF NOT EXISTS sys_users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username VARCHAR(64) NOT NULL UNIQUE,
			nick_name VARCHAR(128),
			email VARCHAR(128),
			phone VARCHAR(11),
			status VARCHAR(4) DEFAULT '2',
			avatar VARCHAR(255),
			sex VARCHAR(255),
			role_id INTEGER,
			remark VARCHAR(255),
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS sys_roles (
			role_id INTEGER PRIMARY KEY AUTOINCREMENT,
			role_name VARCHAR(128) NOT NULL,
			status VARCHAR(4) DEFAULT '2',
			role_key VARCHAR(128),
			role_sort INTEGER,
			flag VARCHAR(128),
			remark VARCHAR(255),
			admin BOOLEAN DEFAULT FALSE,
			data_scope VARCHAR(4),
			created_at DATETIME,
			updated_at DATETIME
		)`,
		`CREATE TABLE IF NOT EXISTS sys_apis (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			handle VARCHAR(128),
			title VARCHAR(128),
			path VARCHAR(128),
			type VARCHAR(16),
			action VARCHAR(16),
			created_at DATETIME,
			updated_at DATETIME
		)`,
	},

	// 種子資料
	SeedData: []map[string]interface{}{
		UserFixtures.AdminUser,
		UserFixtures.ValidUser,
		RoleFixtures.AdminRole,
		RoleFixtures.UserRole,
		APIFixtures.LoginAPI,
		APIFixtures.UserListAPI,
	},

	// 清理 SQL
	CleanupSQL: []string{
		"DELETE FROM sys_users",
		"DELETE FROM sys_roles",
		"DELETE FROM sys_apis",
		"DELETE FROM sqlite_sequence WHERE name IN ('sys_users', 'sys_roles', 'sys_apis')",
	},
}