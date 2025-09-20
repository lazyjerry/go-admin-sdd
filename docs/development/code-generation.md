# 程式碼生成工具說明

本文件詳細介紹 Go-Admin 的程式碼生成工具，包含使用方法、模板設計、自訂配置等內容。

## 概要

Go-Admin 內建強大的程式碼生成工具，可根據資料表結構自動生成 CRUD 業務程式碼，大幅提升開發效率，減少重複性工作。

## 功能特色

### 主要功能

- **自動 CRUD 生成**：根據資料表自動生成增刪改查介面
- **RESTful API**：遵循 REST 設計規範的 API 端點
- **前端頁面生成**：可選擇生成對應的前端管理頁面
- **權限整合**：自動整合 RBAC 權限控制
- **參數驗證**：自動生成資料驗證規則
- **API 文檔**：自動生成 Swagger 文檔

### 支援的生成內容

```
業務模組生成
├── Model (資料模型)
│   ├── 結構體定義
│   ├── 資料庫關聯
│   └── 驗證規則
├── API (介面層)
│   ├── RESTful 端點
│   ├── 請求處理
│   └── 回應格式
├── Service (業務邏輯)
│   ├── CRUD 操作
│   ├── 業務規則
│   └── 資料處理
├── Router (路由配置)
│   ├── 路由註冊
│   ├── 中間件配置
│   └── 權限控制
└── 前端頁面 (可選)
    ├── 列表頁面
    ├── 表單頁面
    └── 詳情頁面
```

## 工具使用

### 指令行介面

```bash
# 基本語法
go run cmd/gen/main.go [選項] 表名

# 常用指令
go run cmd/gen/main.go -table=user_profiles          # 生成單表
go run cmd/gen/main.go -table=orders,products       # 生成多表
go run cmd/gen/main.go -all                          # 生成所有表
go run cmd/gen/main.go -table=users -with-frontend  # 包含前端
```

### 指令參數

| 參數              | 簡寫 | 說明                 | 範例                         |
| ----------------- | ---- | -------------------- | ---------------------------- |
| `--table`         | `-t` | 指定要生成的資料表名 | `-t users`                   |
| `--output`        | `-o` | 輸出目錄             | `-o ./generated`             |
| `--package`       | `-p` | 套件名稱             | `-p admin`                   |
| `--author`        | `-a` | 作者資訊             | `-a "John Doe"`              |
| `--prefix`        |      | 表名前綴             | `--prefix sys_`              |
| `--exclude`       | `-e` | 排除欄位             | `-e id,created_at`           |
| `--with-frontend` | `-f` | 生成前端頁面         | `-f`                         |
| `--template-dir`  |      | 自訂模板目錄         | `--template-dir ./templates` |
| `--config`        | `-c` | 配置檔案             | `-c gen.yaml`                |
| `--dry-run`       |      | 預覽生成內容         | `--dry-run`                  |

### 配置檔案

```yaml
# gen.yaml - 生成工具配置
database:
  host: localhost
  port: 3306
  username: root
  password: ""
  database: go_admin
  charset: utf8mb4

generator:
  # 輸出設定
  output_dir: "./app/generated"
  package_name: "generated"

  # 作者資訊
  author: "Go-Admin Generator"
  email: "admin@example.com"

  # 模板設定
  template_dir: "./template"
  custom_templates: true

  # 生成選項
  generate_frontend: false
  generate_swagger: true
  generate_tests: true

  # 表名設定
  table_prefix: "sys_"
  exclude_tables:
    - "migrations"
    - "sys_oper_logs"

  # 欄位設定
  exclude_fields:
    - "created_at"
    - "updated_at"
    - "deleted_at"

  # 關聯設定
  foreign_keys:
    enabled: true
    suffix: "_id"

naming:
  # 命名規則
  model_suffix: ""
  api_suffix: "Api"
  service_suffix: "Service"

  # 檔案命名
  snake_case: true
  lower_camel_case: false
```

## 生成範例

### 範例 1：使用者管理模組

假設有一個 `user_profiles` 表：

```sql
CREATE TABLE `user_profiles` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `real_name` varchar(64) DEFAULT NULL,
  `avatar` varchar(255) DEFAULT NULL,
  `phone` varchar(20) DEFAULT NULL,
  `email` varchar(128) DEFAULT NULL,
  `birthday` date DEFAULT NULL,
  `gender` tinyint(1) DEFAULT '0',
  `address` varchar(255) DEFAULT NULL,
  `bio` text,
  `status` tinyint(1) DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`)
);
```

執行生成指令：

```bash
go run cmd/gen/main.go -table=user_profiles -p admin
```

### 生成的檔案結構

```
app/admin/
├── models/
│   └── user_profile.go          # 資料模型
├── apis/
│   └── user_profile.go          # API 介面
├── services/
│   └── user_profile.go          # 業務邏輯
└── router/
    └── user_profile.go          # 路由配置
```

### 生成的 Model

```go
// app/admin/models/user_profile.go
package models

import (
    "time"
    "go-admin/common/models"
)

// UserProfile 使用者設定檔
type UserProfile struct {
    models.Model

    UserId   int       `json:"userId" gorm:"column:user_id;comment:使用者ID"`
    RealName string    `json:"realName" gorm:"column:real_name;size:64;comment:真實姓名"`
    Avatar   string    `json:"avatar" gorm:"column:avatar;size:255;comment:頭像"`
    Phone    string    `json:"phone" gorm:"column:phone;size:20;comment:手機號碼"`
    Email    string    `json:"email" gorm:"column:email;size:128;comment:電子信箱"`
    Birthday time.Time `json:"birthday" gorm:"column:birthday;comment:生日"`
    Gender   int       `json:"gender" gorm:"column:gender;comment:性別"`
    Address  string    `json:"address" gorm:"column:address;size:255;comment:地址"`
    Bio      string    `json:"bio" gorm:"column:bio;type:text;comment:個人簡介"`
    Status   int       `json:"status" gorm:"column:status;comment:狀態"`

    models.ModelTime
    models.ControlBy
}

// TableName 設定表名
func (UserProfile) TableName() string {
    return "user_profiles"
}

// Generate 模型生成
func (e *UserProfile) Generate() models.ActiveRecord {
    o := *e
    return &o
}

// GetId 獲取主鍵
func (e *UserProfile) GetId() interface{} {
    return e.Id
}
```

### 生成的 API

```go
// app/admin/apis/user_profile.go
package apis

import (
    "strconv"

    "github.com/gin-gonic/gin"
    "go-admin/app/admin/models"
    "go-admin/app/admin/service"
    "go-admin/common/apis"
    "go-admin/common/dto"
)

type UserProfileApi struct {
    apis.Api
}

// GetUserProfileList 獲取使用者設定檔列表
// @Summary 獲取使用者設定檔列表
// @Description 獲取使用者設定檔列表
// @Tags 使用者設定檔管理
// @Accept json
// @Product json
// @Param data query dto.UserProfileGetPageReq true "查詢參數"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.UserProfile}} "成功"
// @Router /api/v1/user-profiles [get]
// @Security Bearer
func (e UserProfileApi) GetUserProfileList(c *gin.Context) {
    req := dto.UserProfileGetPageReq{}
    s := service.UserProfile{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    list := make([]models.UserProfile, 0)
    var count int64

    err = s.GetPage(&req, &list, &count)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("獲取使用者設定檔列表失敗，\\r\\n失敗資訊 %s", err.Error()))
        return
    }

    e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查詢成功")
}

// GetUserProfile 獲取使用者設定檔詳情
// @Summary 獲取使用者設定檔詳情
// @Description 獲取使用者設定檔詳情
// @Tags 使用者設定檔管理
// @Accept json
// @Product json
// @Param id path int true "使用者設定檔ID"
// @Success 200 {object} response.Response{data=models.UserProfile} "成功"
// @Router /api/v1/user-profiles/{id} [get]
// @Security Bearer
func (e UserProfileApi) GetUserProfile(c *gin.Context) {
    req := dto.UserProfileGetReq{}
    s := service.UserProfile{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req, nil).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
    var object models.UserProfile

    err = s.Get(&req, &object)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("獲取使用者設定檔失敗，\\r\\n失敗資訊 %s", err.Error()))
        return
    }

    e.OK(object, "查詢成功")
}

// InsertUserProfile 建立使用者設定檔
// @Summary 建立使用者設定檔
// @Description 建立使用者設定檔
// @Tags 使用者設定檔管理
// @Accept json
// @Product json
// @Param data body dto.UserProfileInsertReq true "使用者設定檔資訊"
// @Success 200 {object} response.Response{data=models.UserProfile} "成功"
// @Router /api/v1/user-profiles [post]
// @Security Bearer
func (e UserProfileApi) InsertUserProfile(c *gin.Context) {
    req := dto.UserProfileInsertReq{}
    s := service.UserProfile{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
    // 設定建立人
    req.SetCreateBy(user.GetUserId(c))

    err = s.Insert(&req)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, fmt.Sprintf("建立使用者設定檔失敗，\\r\\n失敗資訊 %s", err.Error()))
        return
    }

    e.OK(req.GetId(), "建立成功")
}

// UpdateUserProfile 更新使用者設定檔
// @Summary 更新使用者設定檔
// @Description 更新使用者設定檔
// @Tags 使用者設定檔管理
// @Accept json
// @Product json
// @Param id path int true "使用者設定檔ID"
// @Param data body dto.UserProfileUpdateReq true "使用者設定檔資訊"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/user-profiles/{id} [put]
// @Security Bearer
func (e UserProfileApi) UpdateUserProfile(c *gin.Context) {
    req := dto.UserProfileUpdateReq{}
    s := service.UserProfile{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }
    req.SetUpdateBy(user.GetUserId(c))

    err = s.Update(&req)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, fmt.Sprintf("更新使用者設定檔失敗，\\r\\n失敗資訊 %s", err.Error()))
        return
    }
    e.OK(req.GetId(), "更新成功")
}

// DeleteUserProfile 刪除使用者設定檔
// @Summary 刪除使用者設定檔
// @Description 刪除使用者設定檔
// @Tags 使用者設定檔管理
// @Accept json
// @Product json
// @Param data body dto.UserProfileDeleteReq true "刪除資料"
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/user-profiles [delete]
// @Security Bearer
func (e UserProfileApi) DeleteUserProfile(c *gin.Context) {
    s := service.UserProfile{}
    req := dto.UserProfileDeleteReq{}
    err := e.MakeContext(c).
        MakeOrm().
        Bind(&req).
        MakeService(&s.Service).
        Errors
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, err.Error())
        return
    }

    err = s.Remove(&req)
    if err != nil {
        e.Logger.Error(err)
        e.Error(500, err, fmt.Sprintf("刪除使用者設定檔失敗，\\r\\n失敗資訊 %s", err.Error()))
        return
    }
    e.OK(req.GetId(), "刪除成功")
}
```

### 生成的 Service

```go
// app/admin/service/user_profile.go
package service

import (
    "errors"

    "gorm.io/gorm"

    "go-admin/app/admin/models"
    "go-admin/common/dto"
    "go-admin/common/service"
    cDto "go-admin/common/dto"
)

type UserProfile struct {
    service.Service
}

// GetPage 獲取使用者設定檔列表
func (e *UserProfile) GetPage(c *dto.UserProfileGetPageReq, p *[]models.UserProfile, count *int64) error {
    var err error
    var data models.UserProfile

    err = e.Orm.Model(&data).
        Scopes(
            cDto.MakeCondition(c.GetNeedSearch()),
            cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
        ).
        Find(p).Limit(-1).Offset(-1).
        Count(count).Error
    if err != nil {
        e.Log.Errorf("UserProfile GetPage error:%s\\r\\n", err)
        return err
    }
    return nil
}

// Get 獲取使用者設定檔詳情
func (e *UserProfile) Get(d *dto.UserProfileGetReq, p *models.UserProfile) error {
    var data models.UserProfile

    err := e.Orm.Model(&data).
        First(p, d.GetId()).Error
    if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
        err = errors.New("查看對象不存在或無權限")
        e.Log.Errorf("UserProfile Get error:%s\\r\\n", err)
        return err
    }
    if err != nil {
        e.Log.Errorf("db error:%s", err)
        return err
    }
    return nil
}

// Insert 建立使用者設定檔
func (e *UserProfile) Insert(c *dto.UserProfileInsertReq) error {
    var err error
    var data models.UserProfile
    c.Generate(&data)
    err = e.Orm.Create(&data).Error
    if err != nil {
        e.Log.Errorf("UserProfile Insert error:%s\\r\\n", err)
        return err
    }
    c.SetId(data.GetId())
    return nil
}

// Update 更新使用者設定檔
func (e *UserProfile) Update(c *dto.UserProfileUpdateReq) error {
    var err error
    var data = models.UserProfile{}
    e.Orm.First(&data, c.GetId())
    c.Generate(&data)

    db := e.Orm.Save(&data)
    if err = db.Error; err != nil {
        e.Log.Errorf("UserProfile Update error:%s\\r\\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("無權限更新該資料")
    }
    return nil
}

// Remove 刪除使用者設定檔
func (e *UserProfile) Remove(d *dto.UserProfileDeleteReq) error {
    var data models.UserProfile

    db := e.Orm.Delete(&data, d.GetId())
    if err := db.Error; err != nil {
        e.Log.Errorf("UserProfile Remove error:%s\\r\\n", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("無權限刪除該資料")
    }
    return nil
}
```

### 生成的路由

```go
// app/admin/router/user_profile.go
package router

import (
    "github.com/gin-gonic/gin"
    jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

    "go-admin/app/admin/apis"
    "go-admin/common/middleware"
)

func init() {
    routerCheckRole = append(routerCheckRole, registerUserProfileRouter)
}

// registerUserProfileRouter 註冊使用者設定檔路由
func registerUserProfileRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
    api := apis.UserProfileApi{}
    r := v1.Group("/user-profiles").Use(authMiddleware.MiddlewareFunc()).Use(middleware.AuthCheckRole())
    {
        r.GET("", api.GetUserProfileList)
        r.GET("/:id", api.GetUserProfile)
        r.POST("", api.InsertUserProfile)
        r.PUT("/:id", api.UpdateUserProfile)
        r.DELETE("", api.DeleteUserProfile)
    }
}
```

## 模板系統

### 模板結構

```
template/
├── api.template              # API 介面模板
├── service.template          # 業務邏輯模板
├── model.template           # 資料模型模板
├── router.template          # 路由配置模板
├── dto.template             # 資料傳輸物件模板
├── frontend/                # 前端模板
│   ├── list.vue.template    # 列表頁面
│   ├── form.vue.template    # 表單頁面
│   └── api.js.template      # 前端 API
└── test/                    # 測試模板
    ├── api_test.template    # API 測試
    └── service_test.template# 業務邏輯測試
```

### 自訂模板

```go
// template/api.template
package apis

import (
    "github.com/gin-gonic/gin"
    "{{.PackageName}}/app/{{.ModuleName}}/models"
    "{{.PackageName}}/app/{{.ModuleName}}/service"
    "{{.PackageName}}/common/apis"
)

type {{.ClassName}}Api struct {
    apis.Api
}

// Get{{.ClassName}}List 獲取{{.Comment}}列表
// @Summary 獲取{{.Comment}}列表
// @Description 獲取{{.Comment}}列表
// @Tags {{.Comment}}管理
// @Accept json
// @Product json
// @Success 200 {object} response.Response "成功"
// @Router /api/v1/{{.ModulePath}} [get]
// @Security Bearer
func (e {{.ClassName}}Api) Get{{.ClassName}}List(c *gin.Context) {
    // 實作邏輯...
}

{{range .Methods}}
// {{.Name}} {{.Comment}}
func (e {{$.ClassName}}Api) {{.Name}}(c *gin.Context) {
    // {{.Comment}}實作
}
{{end}}
```

### 模板變數

| 變數               | 說明     | 範例             |
| ------------------ | -------- | ---------------- |
| `{{.PackageName}}` | 套件名稱 | `go-admin`       |
| `{{.ModuleName}}`  | 模組名稱 | `admin`          |
| `{{.ClassName}}`   | 類別名稱 | `UserProfile`    |
| `{{.TableName}}`   | 資料表名 | `user_profiles`  |
| `{{.Comment}}`     | 註釋說明 | `使用者設定檔`   |
| `{{.Author}}`      | 作者資訊 | `Go-Admin Team`  |
| `{{.CreateTime}}`  | 建立時間 | `2025-09-19`     |
| `{{.Fields}}`      | 欄位陣列 | 包含所有欄位資訊 |
| `{{.PrimaryKey}}`  | 主鍵欄位 | `Id`             |
| `{{.Methods}}`     | 方法陣列 | CRUD 方法列表    |

## 進階配置

### 欄位類型映射

```yaml
# 資料庫類型 -> Go 類型映射
type_mapping:
  # 整數類型
  int: "int"
  bigint: "int64"
  tinyint: "int8"

  # 字串類型
  varchar: "string"
  char: "string"
  text: "string"
  longtext: "string"

  # 浮點類型
  float: "float32"
  double: "float64"
  decimal: "decimal.Decimal"

  # 時間類型
  datetime: "time.Time"
  date: "time.Time"
  timestamp: "time.Time"

  # JSON 類型
  json: "datatypes.JSON"

  # 布林類型
  boolean: "bool"
```

### 驗證規則生成

```go
// 自動生成驗證標籤
type UserProfile struct {
    RealName string `json:"realName" validate:"required,min=2,max=64" comment:"真實姓名"`
    Email    string `json:"email" validate:"email,max=128" comment:"電子信箱"`
    Phone    string `json:"phone" validate:"phone,len=11" comment:"手機號碼"`
    Age      int    `json:"age" validate:"min=0,max=150" comment:"年齡"`
}
```

### 關聯關係處理

```go
// 一對多關係
type User struct {
    Id       int           `json:"id"`
    Username string        `json:"username"`
    Profiles []UserProfile `json:"profiles" gorm:"foreignKey:UserId"`
}

// 多對多關係
type Role struct {
    Id    int    `json:"id"`
    Name  string `json:"name"`
    Users []User `json:"users" gorm:"many2many:user_roles;"`
}
```

## 擴展功能

### 自訂生成器

```go
// generators/custom_generator.go
package generators

import (
    "go-admin/cmd/gen/generator"
)

type CustomGenerator struct {
    generator.BaseGenerator
}

func (g *CustomGenerator) Generate(table *generator.TableInfo) error {
    // 自訂生成邏輯
    return g.generateCustomCode(table)
}

func (g *CustomGenerator) generateCustomCode(table *generator.TableInfo) error {
    // 實作自訂生成邏輯
    templateData := map[string]interface{}{
        "TableName": table.Name,
        "Fields":    table.Fields,
        // 其他模板變數
    }

    return g.RenderTemplate("custom.template", templateData, "output.go")
}
```

### 外掛系統

```go
// plugins/swagger_plugin.go
package plugins

type SwaggerPlugin struct{}

func (p *SwaggerPlugin) Name() string {
    return "swagger"
}

func (p *SwaggerPlugin) Execute(context *generator.Context) error {
    // 生成 Swagger 文檔
    return p.generateSwaggerDocs(context)
}

func (p *SwaggerPlugin) generateSwaggerDocs(context *generator.Context) error {
    // 實作 Swagger 文檔生成
    return nil
}
```

## 最佳實踐

### 1. 命名規範

```yaml
# 推薦的命名規範
naming_conventions:
  # 資料表命名 (蛇形命名)
  table_name: "user_profiles"

  # Go 結構體 (大駝峰)
  struct_name: "UserProfile"

  # JSON 欄位 (小駝峰)
  json_field: "realName"

  # 資料庫欄位 (蛇形命名)
  db_field: "real_name"
```

### 2. 模板組織

```
templates/
├── base/                    # 基礎模板
│   ├── model.go.tmpl
│   ├── api.go.tmpl
│   └── service.go.tmpl
├── modules/                 # 模組特定模板
│   ├── admin/
│   └── user/
└── custom/                  # 自訂模板
    └── special.go.tmpl
```

### 3. 程式碼品質

```go
// 生成的程式碼應包含
// 1. 完整的註釋
// 2. 錯誤處理
// 3. 參數驗證
// 4. 日誌記錄
// 5. 權限檢查

func (e UserProfileApi) InsertUserProfile(c *gin.Context) {
    // 參數綁定和驗證
    req := dto.UserProfileInsertReq{}
    if err := e.MakeContext(c).Bind(&req).Errors; err != nil {
        e.Logger.Error(err)
        e.Error(400, err, "參數驗證失敗")
        return
    }

    // 權限檢查
    if !e.checkPermission(c, "user:profile:create") {
        e.Error(403, nil, "無建立權限")
        return
    }

    // 業務邏輯處理
    s := service.UserProfile{}
    if err := s.Insert(&req); err != nil {
        e.Logger.Error(err)
        e.Error(500, err, "建立失敗")
        return
    }

    e.OK(req.GetId(), "建立成功")
}
```

## 疑難排解

### 常見問題

1. **生成失敗**

   ```bash
   # 檢查資料庫連接
   go run cmd/gen/main.go -test-db

   # 檢查表結構
   go run cmd/gen/main.go -describe user_profiles
   ```

2. **模板錯誤**

   ```bash
   # 驗證模板語法
   go run cmd/gen/main.go -validate-template api.template

   # 預覽生成內容
   go run cmd/gen/main.go -table users --dry-run
   ```

3. **權限問題**

   ```bash
   # 檢查檔案權限
   chmod +x cmd/gen/main.go

   # 檢查輸出目錄權限
   chmod 755 app/admin/
   ```

### 除錯模式

```bash
# 啟用詳細日誌
go run cmd/gen/main.go -table users -v

# 啟用除錯模式
go run cmd/gen/main.go -table users --debug

# 輸出詳細資訊
go run cmd/gen/main.go -table users --verbose
```

## 相關文件

- [API 開發指南](./api-guide.md)
- [資料庫設計文檔](./database-design.md)
- [Go-Admin 架構說明](./architecture.md)
- [權限系統說明](./permission-system.md)

## 版本記錄

| 版本  | 日期       | 更新內容                         |
| ----- | ---------- | -------------------------------- |
| 1.0.0 | 2025-09-19 | 初始版本，完整程式碼生成工具說明 |

---

_最後更新：2025 年 9 月 19 日_
