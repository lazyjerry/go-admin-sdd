# Go-Admin API 開發指南

本指南詳細說明如何在 Go-Admin 框架中開發和擴展 RESTful API，包含 API 設計原則、實作範例和最佳實踐。

## API 設計原則

### RESTful API 標準

Go-Admin 遵循 RESTful API 設計標準，使用標準 HTTP 方法和狀態碼：

| HTTP 方法 | 用途     | 範例 URL                 | 說明           |
| --------- | -------- | ------------------------ | -------------- |
| GET       | 查詢資源 | `GET /api/v1/users`      | 取得使用者列表 |
| POST      | 建立資源 | `POST /api/v1/users`     | 建立新使用者   |
| PUT       | 完整更新 | `PUT /api/v1/users/1`    | 完整更新使用者 |
| PATCH     | 部分更新 | `PATCH /api/v1/users/1`  | 部分更新使用者 |
| DELETE    | 刪除資源 | `DELETE /api/v1/users/1` | 刪除使用者     |

### 統一回應格式

```json
{
	"code": 200,
	"message": "操作成功",
	"data": {
		// 實際資料內容
	},
	"timestamp": "2025-09-19T10:30:00Z",
	"requestId": "req_12345"
}
```

### 狀態碼規範

- **200 OK**: 請求成功
- **201 Created**: 資源建立成功
- **400 Bad Request**: 請求參數錯誤
- **401 Unauthorized**: 未授權存取
- **403 Forbidden**: 權限不足
- **404 Not Found**: 資源不存在
- **500 Internal Server Error**: 伺服器內部錯誤

## API 開發流程

### 1. 定義資料模型

首先建立資料模型，定義資料結構和資料庫關聯：

```go
// app/admin/models/product.go
package models

import (
    "go-admin/common/models"
)

// Product 產品模型
type Product struct {
    models.Model
    Name        string  `json:"name" gorm:"size:100;not null;comment:產品名稱"`
    Description string  `json:"description" gorm:"type:text;comment:產品描述"`
    Price       float64 `json:"price" gorm:"type:decimal(10,2);not null;comment:價格"`
    CategoryID  int     `json:"categoryId" gorm:"not null;comment:分類ID"`
    Status      int     `json:"status" gorm:"default:1;comment:狀態 1-啟用 0-停用"`

    // 關聯關係
    Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID"`

    models.ControlBy
    models.ModelTime
}

// TableName 指定資料表名稱
func (Product) TableName() string {
    return "products"
}

// Generate 程式碼生成器接口
func (e *Product) Generate() models.ActiveRecord {
    o := *e
    return &o
}

// GetId 取得主鍵
func (e *Product) GetId() interface{} {
    return e.Id
}
```

### 2. 建立 DTO (資料傳輸物件)

定義 API 請求和回應的資料結構：

```go
// app/admin/service/dto/product.go
package dto

import (
    "go-admin/app/admin/models"
    "go-admin/common/dto"
    common "go-admin/common/models"
)

// ProductGetPageReq 分頁查詢請求
type ProductGetPageReq struct {
    dto.Pagination `search:"-"`
    Name           string `form:"name" search:"type:contains;column:name;table:product" comment:"產品名稱"`
    CategoryID     int    `form:"categoryId" search:"type:exact;column:category_id;table:product" comment:"分類ID"`
    Status         int    `form:"status" search:"type:exact;column:status;table:product" comment:"狀態"`
    dto.ControlBy  `search:"-"`
}

func (m *ProductGetPageReq) GetNeedSearch() interface{} {
    return *m
}

// ProductInsertReq 建立請求
type ProductInsertReq struct {
    Id          int     `json:"-" comment:"主鍵"`
    Name        string  `json:"name" comment:"產品名稱" validate:"required,max=100"`
    Description string  `json:"description" comment:"產品描述"`
    Price       float64 `json:"price" comment:"價格" validate:"required,min=0"`
    CategoryID  int     `json:"categoryId" comment:"分類ID" validate:"required"`
    Status      int     `json:"status" comment:"狀態" validate:"oneof=0 1"`
    common.ControlBy
}

func (s *ProductInsertReq) Generate(model *models.Product) {
    if s.Id == 0 {
        model.Model = common.Model{Id: s.Id}
    }
    model.Name = s.Name
    model.Description = s.Description
    model.Price = s.Price
    model.CategoryID = s.CategoryID
    model.Status = s.Status
}

func (s *ProductInsertReq) GetId() interface{} {
    return s.Id
}

// ProductUpdateReq 更新請求
type ProductUpdateReq struct {
    Id          int     `uri:"id" comment:"主鍵"`
    Name        string  `json:"name" comment:"產品名稱" validate:"required,max=100"`
    Description string  `json:"description" comment:"產品描述"`
    Price       float64 `json:"price" comment:"價格" validate:"required,min=0"`
    CategoryID  int     `json:"categoryId" comment:"分類ID" validate:"required"`
    Status      int     `json:"status" comment:"狀態" validate:"oneof=0 1"`
    common.ControlBy
}

func (s *ProductUpdateReq) Generate(model *models.Product) {
    model.Name = s.Name
    model.Description = s.Description
    model.Price = s.Price
    model.CategoryID = s.CategoryID
    model.Status = s.Status
}

func (s *ProductUpdateReq) GetId() interface{} {
    return s.Id
}

// ProductGetReq 單個查詢請求
type ProductGetReq struct {
    Id int `uri:"id"`
}

func (s *ProductGetReq) GetId() interface{} {
    return s.Id
}

// ProductDeleteReq 刪除請求
type ProductDeleteReq struct {
    Ids []int `json:"ids" comment:"主鍵列表"`
    common.ControlBy
}

func (s *ProductDeleteReq) GetId() interface{} {
    return s.Ids
}
```

### 3. 實作業務邏輯層

建立 Service 層處理業務邏輯：

```go
// app/admin/service/product.go
package service

import (
    "errors"
    "gorm.io/gorm"

    "go-admin/app/admin/models"
    "go-admin/app/admin/service/dto"
    "go-admin/common/actions"
    cDto "go-admin/common/dto"
    "go-admin/common/log"
    "go-admin/common/service"
)

type Product struct {
    service.Service
}

// GetPage 分頁查詢產品
func (e *Product) GetPage(c *dto.ProductGetPageReq, p *actions.DataPermission, list *[]models.Product, count *int64) error {
    var err error
    var data models.Product

    err = e.Orm.Model(&data).
        Scopes(
            cDto.MakeCondition(c.GetNeedSearch()),
            cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
            actions.Permission(data.TableName(), p),
        ).
        Preload("Category").  // 預載入分類資訊
        Find(list).Limit(-1).Offset(-1).
        Count(count).Error
    if err != nil {
        log.Errorf("ProductService GetPage error: %s", err)
        return err
    }
    return nil
}

// Get 查詢單個產品
func (e *Product) Get(d *dto.ProductGetReq, p *actions.DataPermission, model *models.Product) error {
    var data models.Product

    err := e.Orm.Model(&data).
        Scopes(
            actions.Permission(data.TableName(), p),
        ).
        Preload("Category").
        First(model, d.GetId()).Error
    if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
        err = errors.New("查看物件不存在或無權查看")
        log.Errorf("ProductService Get error: %s", err)
        return err
    }
    if err != nil {
        log.Errorf("ProductService Get error: %s", err)
        return err
    }
    return nil
}

// Insert 建立產品
func (e *Product) Insert(c *dto.ProductInsertReq) error {
    var err error
    var data models.Product
    c.Generate(&data)

    // 業務邏輯驗證
    if err = e.validateProduct(&data); err != nil {
        return err
    }

    err = e.Orm.Create(&data).Error
    if err != nil {
        log.Errorf("ProductService Insert error: %s", err)
        return err
    }
    return nil
}

// Update 更新產品
func (e *Product) Update(c *dto.ProductUpdateReq, p *actions.DataPermission) error {
    var err error
    var data = models.Product{}

    // 檢查記錄是否存在
    err = e.Orm.Scopes(
        actions.Permission(data.TableName(), p),
    ).First(&data, c.GetId()).Error
    if err != nil {
        log.Errorf("ProductService Update error: %s", err)
        return err
    }

    // 更新資料
    c.Generate(&data)

    // 業務邏輯驗證
    if err = e.validateProduct(&data); err != nil {
        return err
    }

    err = e.Orm.Save(&data).Error
    if err != nil {
        log.Errorf("ProductService Update error: %s", err)
        return err
    }
    return nil
}

// Remove 刪除產品
func (e *Product) Remove(c *dto.ProductDeleteReq, p *actions.DataPermission) error {
    var data models.Product

    db := e.Orm.Model(&data).
        Scopes(
            actions.Permission(data.TableName(), p),
        ).Delete(&data, c.GetId())
    if err := db.Error; err != nil {
        log.Errorf("ProductService Remove error: %s", err)
        return err
    }
    if db.RowsAffected == 0 {
        return errors.New("無權刪除該資料")
    }
    return nil
}

// validateProduct 驗證產品資料
func (e *Product) validateProduct(product *models.Product) error {
    // 檢查分類是否存在
    var category models.Category
    if err := e.Orm.First(&category, product.CategoryID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return errors.New("指定的分類不存在")
        }
        return err
    }

    // 檢查產品名稱是否重複
    var count int64
    err := e.Orm.Model(&models.Product{}).
        Where("name = ? AND id != ?", product.Name, product.Id).
        Count(&count).Error
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("產品名稱已存在")
    }

    return nil
}
```

### 4. 建立 API 控制器

實作 API 端點處理器：

```go
// app/admin/apis/product.go
package apis

import (
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/go-admin-team/go-admin-core/sdk/api"
    "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
    _ "github.com/go-admin-team/go-admin-core/sdk/pkg/response"

    "go-admin/app/admin/models"
    "go-admin/app/admin/service"
    "go-admin/app/admin/service/dto"
    "go-admin/common/actions"
)

type Product struct {
    api.Api
}

// GetPage 分頁查詢產品
// @Summary 分頁查詢產品列表
// @Description 取得產品分頁列表
// @Tags 產品管理
// @Param name query string false "產品名稱"
// @Param categoryId query int false "分類ID"
// @Param status query int false "狀態"
// @Param pageSize query int false "頁面大小"
// @Param pageIndex query int false "頁面索引"
// @Success 200 {object} response.Response{data=response.Page{list=[]models.Product}} "成功"
// @Router /api/v1/product [get]
// @Security Bearer
func (e Product) GetPage(c *gin.Context) {
    req := dto.ProductGetPageReq{}
    s := service.Product{}
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

    p := actions.GetPermissionFromContext(c)
    list := make([]models.Product, 0)
    var count int64

    err = s.GetPage(&req, p, &list, &count)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("取得產品列表失敗，%s", err.Error()))
        return
    }

    e.PageOK(list, int(count), req.GetPageIndex(), req.GetPageSize(), "查詢成功")
}

// Get 查詢單個產品
// @Summary 查詢產品詳情
// @Description 根據ID查詢產品詳情
// @Tags 產品管理
// @Param id path int true "產品ID"
// @Success 200 {object} response.Response{data=models.Product} "成功"
// @Router /api/v1/product/{id} [get]
// @Security Bearer
func (e Product) Get(c *gin.Context) {
    req := dto.ProductGetReq{}
    s := service.Product{}
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
    var object models.Product

    p := actions.GetPermissionFromContext(c)
    err = s.Get(&req, p, &object)
    if err != nil {
        e.Error(http.StatusUnprocessableEntity, err, fmt.Sprintf("取得產品失敗，%s", err.Error()))
        return
    }

    e.OK( object, "查詢成功")
}

// Insert 建立產品
// @Summary 建立產品
// @Description 建立新的產品
// @Tags 產品管理
// @Accept application/json
// @Product application/json
// @Param data body dto.ProductInsertReq true "產品資料"
// @Success 200 {object} response.Response	"成功"
// @Router /api/v1/product [post]
// @Security Bearer
func (e Product) Insert(c *gin.Context) {
    req := dto.ProductInsertReq{}
    s := service.Product{}
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
    // 設定建立人員
    req.SetCreateBy(user.GetUserId(c))

    err = s.Insert(&req)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("建立產品失敗，%s", err.Error()))
        return
    }

    e.OK(req.GetId(), "建立成功")
}

// Update 更新產品
// @Summary 更新產品
// @Description 更新產品資訊
// @Tags 產品管理
// @Accept application/json
// @Product application/json
// @Param id path int true "產品ID"
// @Param data body dto.ProductUpdateReq true "產品資料"
// @Success 200 {object} response.Response	"成功"
// @Router /api/v1/product/{id} [put]
// @Security Bearer
func (e Product) Update(c *gin.Context) {
    req := dto.ProductUpdateReq{}
    s := service.Product{}
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
    p := actions.GetPermissionFromContext(c)

    err = s.Update(&req, p)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("更新產品失敗，%s", err.Error()))
        return
    }
    e.OK( req.GetId(), "更新成功")
}

// Delete 刪除產品
// @Summary 刪除產品
// @Description 刪除產品（支援批量刪除）
// @Tags 產品管理
// @Param data body dto.ProductDeleteReq true "產品ID列表"
// @Success 200 {object} response.Response	"成功"
// @Router /api/v1/product [delete]
// @Security Bearer
func (e Product) Delete(c *gin.Context) {
    s := service.Product{}
    req := dto.ProductDeleteReq{}
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
    p := actions.GetPermissionFromContext(c)

    err = s.Remove(&req, p)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("刪除產品失敗，%s", err.Error()))
        return
    }
    e.OK( req.GetId(), "刪除成功")
}
```

### 5. 註冊路由

建立路由定義檔案：

```go
// app/admin/router/product.go
package router

import (
    "github.com/gin-gonic/gin"
    jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"

    "go-admin/app/admin/apis"
    "go-admin/common/middleware"
)

func init() {
    routerNoCheckRole = append(routerNoCheckRole, registerProductRouter)
}

// registerProductRouter 註冊產品路由
func registerProductRouter(v1 *gin.RouterGroup) {
    api := apis.Product{}
    r := v1.Group("/product").Use(middleware.Auth()).Use(middleware.AuthCheckRole())
    {
        r.GET("", api.GetPage)
        r.GET("/:id", api.Get)
        r.POST("", api.Insert)
        r.PUT("/:id", api.Update)
        r.DELETE("", api.Delete)
    }
}
```

## 進階 API 功能

### 1. 檔案上傳處理

```go
// UploadFile 檔案上傳
func (e Product) UploadFile(c *gin.Context) {
    file, header, err := c.Request.FormFile("file")
    if err != nil {
        e.Error(400, err, "檔案上傳失敗")
        return
    }
    defer file.Close()

    // 檔案類型驗證
    allowedTypes := []string{"image/jpeg", "image/png", "image/gif"}
    contentType := header.Header.Get("Content-Type")
    if !contains(allowedTypes, contentType) {
        e.Error(400, nil, "不支援的檔案類型")
        return
    }

    // 檔案大小驗證
    maxSize := int64(10 * 1024 * 1024) // 10MB
    if header.Size > maxSize {
        e.Error(400, nil, "檔案大小超過限制")
        return
    }

    // 儲存檔案
    filename := generateFilename(header.Filename)
    filepath := fmt.Sprintf("storage/uploads/%s", filename)

    if err := c.SaveUploadedFile(header, filepath); err != nil {
        e.Error(500, err, "檔案儲存失敗")
        return
    }

    e.OK(gin.H{
        "filename": filename,
        "url":      fmt.Sprintf("/uploads/%s", filename),
        "size":     header.Size,
    }, "上傳成功")
}
```

### 2. 資料匯出功能

```go
// ExportExcel 匯出 Excel
func (e Product) ExportExcel(c *gin.Context) {
    req := dto.ProductGetPageReq{}
    s := service.Product{}
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

    // 取得所有資料（不分頁）
    req.PageSize = -1
    p := actions.GetPermissionFromContext(c)
    list := make([]models.Product, 0)
    var count int64

    err = s.GetPage(&req, p, &list, &count)
    if err != nil {
        e.Error(500, err, fmt.Sprintf("取得資料失敗，%s", err.Error()))
        return
    }

    // 建立 Excel 檔案
    file := excelize.NewFile()
    sheet := "產品列表"

    // 設定標題
    headers := []string{"ID", "產品名稱", "描述", "價格", "分類", "狀態", "建立時間"}
    for i, header := range headers {
        cell := fmt.Sprintf("%s1", string(rune('A'+i)))
        file.SetCellValue(sheet, cell, header)
    }

    // 填入資料
    for i, product := range list {
        row := i + 2
        file.SetCellValue(sheet, fmt.Sprintf("A%d", row), product.Id)
        file.SetCellValue(sheet, fmt.Sprintf("B%d", row), product.Name)
        file.SetCellValue(sheet, fmt.Sprintf("C%d", row), product.Description)
        file.SetCellValue(sheet, fmt.Sprintf("D%d", row), product.Price)
        file.SetCellValue(sheet, fmt.Sprintf("E%d", row), product.Category.Name)
        file.SetCellValue(sheet, fmt.Sprintf("F%d", row), getStatusText(product.Status))
        file.SetCellValue(sheet, fmt.Sprintf("G%d", row), product.CreatedAt.Format("2006-01-02 15:04:05"))
    }

    // 設定回應標頭
    filename := fmt.Sprintf("products_%s.xlsx", time.Now().Format("20060102150405"))
    c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

    // 輸出檔案
    if err := file.Write(c.Writer); err != nil {
        e.Error(500, err, "匯出失敗")
        return
    }
}
```

### 3. 批次操作

```go
// BatchUpdate 批次更新
func (e Product) BatchUpdate(c *gin.Context) {
    req := struct {
        Ids    []int `json:"ids" validate:"required"`
        Status int   `json:"status" validate:"oneof=0 1"`
    }{}

    s := service.Product{}
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

    p := actions.GetPermissionFromContext(c)

    // 批次更新
    var data models.Product
    db := e.Orm.Model(&data).
        Scopes(actions.Permission(data.TableName(), p)).
        Where("id IN ?", req.Ids).
        Update("status", req.Status)

    if err := db.Error; err != nil {
        e.Error(500, err, "批次更新失敗")
        return
    }

    e.OK(gin.H{
        "updated": db.RowsAffected,
    }, fmt.Sprintf("成功更新 %d 筆記錄", db.RowsAffected))
}
```

## API 測試

### 1. 單元測試

```go
// app/admin/apis/product_test.go
package apis

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

    "go-admin/app/admin/service/dto"
    "go-admin/common/database"
)

func TestProduct_Insert(t *testing.T) {
    // 設定測試環境
    gin.SetMode(gin.TestMode)
    r := gin.New()

    // 初始化資料庫
    database.Setup()

    // 註冊路由
    api := Product{}
    r.POST("/api/v1/product", api.Insert)

    // 準備測試資料
    product := dto.ProductInsertReq{
        Name:        "測試產品",
        Description: "這是一個測試產品",
        Price:       99.99,
        CategoryID:  1,
        Status:      1,
    }

    jsonData, _ := json.Marshal(product)

    // 發送請求
    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/api/v1/product", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")

    r.ServeHTTP(w, req)

    // 驗證結果
    assert.Equal(t, 200, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, 200, int(response["code"].(float64)))
}

func TestProduct_GetPage(t *testing.T) {
    gin.SetMode(gin.TestMode)
    r := gin.New()

    api := Product{}
    r.GET("/api/v1/product", api.GetPage)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/api/v1/product?pageSize=10&pageIndex=1", nil)

    r.ServeHTTP(w, req)

    assert.Equal(t, 200, w.Code)

    var response map[string]interface{}
    err := json.Unmarshal(w.Body.Bytes(), &response)
    assert.NoError(t, err)
    assert.Equal(t, 200, int(response["code"].(float64)))
}
```

### 2. 整合測試

```bash
#!/bin/bash
# test_api.sh - API 整合測試腳本

BASE_URL="http://localhost:8000/api/v1"
TOKEN=""

# 登入取得 Token
login() {
    response=$(curl -s -X POST "$BASE_URL/login" \
        -H "Content-Type: application/json" \
        -d '{"username":"admin","password":"admin123"}')

    TOKEN=$(echo $response | jq -r '.data.token')
    echo "Login successful, token: $TOKEN"
}

# 測試產品 CRUD
test_product_crud() {
    echo "Testing Product CRUD..."

    # 1. 建立產品
    echo "1. Creating product..."
    create_response=$(curl -s -X POST "$BASE_URL/product" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "測試產品",
            "description": "API 測試產品",
            "price": 199.99,
            "categoryId": 1,
            "status": 1
        }')

    product_id=$(echo $create_response | jq -r '.data')
    echo "Created product ID: $product_id"

    # 2. 查詢產品
    echo "2. Getting product..."
    get_response=$(curl -s -X GET "$BASE_URL/product/$product_id" \
        -H "Authorization: Bearer $TOKEN")

    echo "Get response: $get_response"

    # 3. 更新產品
    echo "3. Updating product..."
    update_response=$(curl -s -X PUT "$BASE_URL/product/$product_id" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
            "name": "測試產品(更新)",
            "description": "API 測試產品(更新)",
            "price": 299.99,
            "categoryId": 1,
            "status": 1
        }')

    echo "Update response: $update_response"

    # 4. 刪除產品
    echo "4. Deleting product..."
    delete_response=$(curl -s -X DELETE "$BASE_URL/product" \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d "{\"ids\":[$product_id]}")

    echo "Delete response: $delete_response"
}

# 執行測試
main() {
    login
    if [ -n "$TOKEN" ]; then
        test_product_crud
        echo "API tests completed!"
    else
        echo "Login failed!"
        exit 1
    fi
}

main
```

## API 文件生成

### Swagger 註解規範

Go-Admin 使用 swaggo 自動生成 Swagger 文件，需要在 API 方法上添加註解：

```go
// @Summary 方法摘要
// @Description 詳細描述
// @Tags 標籤（用於分組）
// @Accept 接受的內容類型
// @Produce 輸出的內容類型
// @Param 參數名 參數位置 參數類型 是否必須 "參數描述"
// @Success 200 {object} 回應類型 "成功描述"
// @Failure 400 {object} response.Response "失敗描述"
// @Router /api路徑 [HTTP方法]
// @Security Bearer
```

### 生成文件

```bash
# 安裝 swag 工具
go install github.com/swaggo/swag/cmd/swag@latest

# 生成 Swagger 文件
swag init

# 或使用專案的 Makefile
make docs

# 存取文件 URL
# http://localhost:8000/swagger/index.html
```

## 效能最佳化

### 1. 資料庫查詢最佳化

```go
// 避免 N+1 問題，使用 Preload
func (e *Product) GetPage(c *dto.ProductGetPageReq, p *actions.DataPermission, list *[]models.Product, count *int64) error {
    var data models.Product

    err := e.Orm.Model(&data).
        Scopes(
            cDto.MakeCondition(c.GetNeedSearch()),
            cDto.Paginate(c.GetPageSize(), c.GetPageIndex()),
            actions.Permission(data.TableName(), p),
        ).
        Preload("Category").
        Preload("Tags").
        Find(list).Limit(-1).Offset(-1).
        Count(count).Error
    return err
}

// 使用 Select 只查詢需要的欄位
func (e *Product) GetList() ([]models.Product, error) {
    var products []models.Product
    err := e.Orm.Select("id, name, price, status").
        Where("status = ?", 1).
        Find(&products).Error
    return products, err
}
```

### 2. 快取策略

```go
// 使用 Redis 快取
func (e *Product) GetCachedProduct(id int) (*models.Product, error) {
    // 嘗試從快取取得
    cacheKey := fmt.Sprintf("product:%d", id)
    cached, err := redis.Get(cacheKey)
    if err == nil && cached != "" {
        var product models.Product
        json.Unmarshal([]byte(cached), &product)
        return &product, nil
    }

    // 從資料庫查詢
    var product models.Product
    err = e.Orm.First(&product, id).Error
    if err != nil {
        return nil, err
    }

    // 寫入快取
    productJSON, _ := json.Marshal(product)
    redis.Set(cacheKey, string(productJSON), time.Hour)

    return &product, nil
}
```

### 3. 請求限流

```go
// 使用 gin 中間件實作限流
func RateLimit(maxRequests int, window time.Duration) gin.HandlerFunc {
    limiter := make(map[string]*rate.Limiter)
    mutex := &sync.RWMutex{}

    return func(c *gin.Context) {
        ip := c.ClientIP()

        mutex.Lock()
        if _, exists := limiter[ip]; !exists {
            limiter[ip] = rate.NewLimiter(rate.Every(window/time.Duration(maxRequests)), maxRequests)
        }
        mutex.Unlock()

        if !limiter[ip].Allow() {
            c.JSON(429, gin.H{
                "code":    429,
                "message": "請求頻率過高，請稍後再試",
            })
            c.Abort()
            return
        }

        c.Next()
    }
}
```

## 安全最佳實踐

### 1. 輸入驗證

```go
// 使用 validator 進行參數驗證
type ProductInsertReq struct {
    Name        string  `json:"name" validate:"required,min=1,max=100"`
    Description string  `json:"description" validate:"max=1000"`
    Price       float64 `json:"price" validate:"required,min=0"`
    CategoryID  int     `json:"categoryId" validate:"required,gt=0"`
    Status      int     `json:"status" validate:"oneof=0 1"`
}

// 自訂驗證器
func validateProductName(fl validator.FieldLevel) bool {
    name := fl.Field().String()
    // 檢查是否包含非法字元
    matched, _ := regexp.MatchString(`^[a-zA-Z0-9\s\-_]+$`, name)
    return matched
}
```

### 2. SQL 注入防護

```go
// 正確使用參數化查詢
func (e *Product) SearchByName(name string) ([]models.Product, error) {
    var products []models.Product
    // 安全：使用參數化查詢
    err := e.Orm.Where("name LIKE ?", "%"+name+"%").Find(&products).Error
    return products, err
}

// 錯誤：字串拼接容易導致 SQL 注入
func (e *Product) SearchByNameUnsafe(name string) ([]models.Product, error) {
    var products []models.Product
    // 危險：直接拼接 SQL
    sql := fmt.Sprintf("SELECT * FROM products WHERE name LIKE '%%%s%%'", name)
    err := e.Orm.Raw(sql).Find(&products).Error
    return products, err
}
```

### 3. 權限控制

```go
// 資料權限檢查
func (e *Product) CheckDataPermission(c *gin.Context, productId int) error {
    userId := jwt.GetUserId(c)
    user, err := getUserById(userId)
    if err != nil {
        return err
    }

    // 檢查使用者是否有存取該產品的權限
    if user.Role != "admin" {
        var product models.Product
        err := e.Orm.Where("id = ? AND create_by = ?", productId, userId).
            First(&product).Error
        if err != nil {
            return errors.New("無權限存取該資源")
        }
    }

    return nil
}
```

---

遵循本指南的 API 開發規範，可以建立安全、高效且易於維護的 RESTful API。

相關文件：

- [模組建立教學](./module-creation.md)
- [權限系統說明](./permission-system.md)
- [資料庫設計文檔](./database-design.md)
- [程式碼生成工具](./code-generation.md)
