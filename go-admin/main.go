// Package main Go-Admin 系統主程式入口
// 基於 Gin + GORM + Casbin + Swagger 的前後端分離權限管理系統
//
// 專案特色：
// - RESTful API 設計
// - JWT 身份驗證
// - RBAC 權限控制
// - 支援多種資料庫
// - Docker 容器化部署
// - 自動程式碼生成
//
// 技術架構：
// - Web 框架: Gin
// - ORM 框架: GORM
// - 權限控制: Casbin
// - API 文件: Swagger
// - 日誌處理: logrus
// - 配置管理: viper
//
// 作者: go-admin-team
// 版本: 2.0.0
// 授權: MIT License
package main

import (
	"go-admin/cmd"
)

//go:generate swag init --parseDependency --parseDepth=6 --instanceName admin -o ./docs/admin

// API 文件配置
// @title Go-Admin 管理系統 API
// @version 2.0.0
// @description 基於 Gin + Vue + Element UI 的前後端分離權限管理系統 API 文件
// @description 這是一個功能完整的企業級管理後台系統，支援完整的 RBAC 權限控制
// @description
// @description 主要功能：
// @description - 使用者管理：完整的使用者生命周期管理
// @description - 角色權限：靈活的角色權限分配系統
// @description - 選單管理：動態選單配置和權限控制
// @description - 部門管理：組織架構和資料權限管理
// @description - API 權限：細粒度的 API 存取控制
// @description - 操作日誌：完整的系統操作審計追蹤
// @description - 系統監控：即時的系統狀態監控
// @description
// @termsOfService https://github.com/go-admin-team/go-admin
// @contact.name go-admin-team
// @contact.url https://github.com/go-admin-team/go-admin
// @contact.email support@go-admin.dev
// @license.name MIT License
// @license.url https://github.com/go-admin-team/go-admin/blob/master/LICENSE.md
// @host localhost:8000
// @BasePath /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description JWT 身份驗證令牌，格式：Bearer {token}

// main 系統主程式入口函數
// 負責啟動整個應用程式，包含以下初始化流程：
// 1. 載入配置檔案
// 2. 初始化資料庫連接
// 3. 設定路由和中間件
// 4. 啟動 HTTP 服務器
// 5. 優雅關機處理
func main() {
	// 執行 Cobra 命令列程式
	// 支援多種操作模式：server（啟動服務器）、migrate（資料庫遷移）、gen（程式碼生成）等
	cmd.Execute()
}
