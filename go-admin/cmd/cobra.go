// Package cmd 命令列介面實作
// 基於 Cobra 框架實作的多功能命令列工具
// 支援服務器啟動、資料庫遷移、程式碼生成等多種操作模式
package cmd

import (
	"errors"
	"fmt"
	"go-admin/cmd/app"
	"go-admin/common/global"
	"os"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"github.com/spf13/cobra"

	"go-admin/cmd/api"
	"go-admin/cmd/config"
	"go-admin/cmd/migrate"
	"go-admin/cmd/version"
)

// rootCmd 根命令定義
// 作為所有子命令的入口點，提供統一的命令列介面
var rootCmd = &cobra.Command{
	Use:          "go-admin",                              // 命令名稱
	Short:        "Go-Admin 企業級管理系統",                    // 簡短描述
	SilenceUsage: true,                                   // 發生錯誤時不顯示使用說明
	Long: `Go-Admin 是一個基於 Gin + Vue + Element UI 的企業級管理系統

主要功能：
  • 完整的 RBAC 權限管理系統
  • 支援多種資料庫（MySQL、PostgreSQL、SQLite）
  • RESTful API 設計
  • JWT 身份驗證
  • Docker 容器化部署
  • 自動程式碼生成
  • 豐富的中間件支援
  
技術架構：
  • 後端：Go + Gin + GORM + Casbin
  • 前端：Vue + Element UI + TypeScript
  • 資料庫：支援主流關係型資料庫
  • 快取：Redis
  • 訊息佇列：可選支援`,
	
	// 參數驗證函數
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("請至少提供一個命令參數"))
		}
		return nil
	},
	
	// 全域前置處理函數（在所有子命令執行前）
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	
	// 根命令執行函數（當沒有指定子命令時執行）
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

// tip 顯示使用提示訊息
// 包含版本資訊、基本使用說明和文件連結
func tip() {
	usageStr := `歡迎使用 ` + pkg.Green(`Go-Admin `+global.Version) + ` 企業級管理系統！`
	usageStr1 := `使用 ` + pkg.Red(`-h`) + ` 查看可用命令，或使用 ` + pkg.Red(`[命令] -h`) + ` 查看具體命令說明`
	usageStr2 := `更多資訊請參考：https://doc.go-admin.dev/guide/ksks`
	
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1) 
	fmt.Printf("%s\n", usageStr2)
}

// init 初始化函數，註冊所有可用的子命令
func init() {
	// 註冊 API 服務器命令（啟動 HTTP 服務）
	rootCmd.AddCommand(api.StartCmd)
	
	// 註冊資料庫遷移命令（執行資料庫結構變更）
	rootCmd.AddCommand(migrate.StartCmd)
	
	// 註冊版本查詢命令（顯示版本資訊）
	rootCmd.AddCommand(version.StartCmd)
	
	// 註冊配置管理命令（配置檔案相關操作）
	rootCmd.AddCommand(config.StartCmd)
	
	// 註冊應用程式命令（程式碼生成等工具）
	rootCmd.AddCommand(app.StartCmd)
}

// Execute 執行命令列程式
// 作為整個應用程式的入口點，處理命令解析和執行
// 如果命令執行失敗，程式將以錯誤狀態碼退出
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// 命令執行失敗時，以非零狀態碼退出
		// 這遵循 Unix 慣例，便於腳本和 CI/CD 系統處理
		os.Exit(-1)
	}
}