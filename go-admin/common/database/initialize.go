// Package database 資料庫初始化和配置管理
// 負責系統資料庫連接的建立、配置和管理，支援多資料庫架構
package database

import (
	"time"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk"
	toolsConfig "github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	mycasbin "github.com/go-admin-team/go-admin-core/sdk/pkg/casbin"
	toolsDB "github.com/go-admin-team/go-admin-core/tools/database"
	. "github.com/go-admin-team/go-admin-core/tools/gorm/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"

	"go-admin/common/global"
)

// Setup 初始化資料庫配置
// 功能說明：
// - 遍歷所有配置的資料庫
// - 建立多個資料庫連接
// - 支援主從架構和讀寫分離
// - 初始化 Casbin 權限控制
//
// 支援特性：
// - 多資料庫實例管理
// - 連接池配置和最佳化
// - 資料庫健康檢查
// - 自動重連機制
func Setup() {
	// 遍歷配置檔案中定義的所有資料庫配置
	// 每個配置可能包含不同的資料庫實例（如主庫、從庫等）
	for k := range toolsConfig.DatabasesConfig {
		setupSimpleDatabase(k, toolsConfig.DatabasesConfig[k])
	}
}

// setupSimpleDatabase 建立單個資料庫連接
// 參數說明：
//   - host: 資料庫主機識別碼，用於多資料庫管理
//   - c: 資料庫配置物件，包含連接參數和池配置
//
// 功能流程：
// 1. 設定全域資料庫驅動類型
// 2. 配置資料庫連接池參數
// 3. 建立 GORM 資料庫實例
// 4. 初始化 Casbin 權限控制器
// 5. 註冊到 SDK 執行時環境
func setupSimpleDatabase(host string, c *toolsConfig.Database) {
	// 設定全域資料庫驅動（如果尚未設定）
	// 用於後續的資料庫操作和 SQL 語法適配
	if global.Driver == "" {
		global.Driver = c.Driver
	}
	
	// 記錄資料庫連接資訊（不包含敏感資訊）
	log.Infof("正在初始化資料庫: %s => %s", host, pkg.Green(c.Source))
	
	// 配置資料庫解析器（支援讀寫分離）
	// registers 用於配置主從庫的路由策略
	registers := make([]toolsDB.ResolverConfigure, len(c.Registers))
	for i := range c.Registers {
		registers[i] = toolsDB.NewResolverConfigure(
			c.Registers[i].Sources,  // 主庫連接配置
			c.Registers[i].Replicas, // 從庫連接配置
			c.Registers[i].Policy,   // 讀寫分離策略
			c.Registers[i].Tables)   // 適用的資料表
	}
	
	// 建立資料庫連接配置
	// 包含連接池參數：最大空閒連接數、最大開啟連接數、連接生命週期等
	resolverConfig := toolsDB.NewConfigure(
		c.Source,           // 資料庫連接字串
		c.MaxIdleConns,     // 最大空閒連接數
		c.MaxOpenConns,     // 最大開啟連接數  
		c.ConnMaxIdleTime,  // 連接最大空閒時間
		c.ConnMaxLifeTime,  // 連接最大生命週期
		registers)          // 讀寫分離配置
	
	// 初始化 GORM 資料庫實例
	db, err := resolverConfig.Init(&gorm.Config{
		// 命名策略配置
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用單數形式的資料表名稱
		},
		// 日誌配置
		Logger: New(
			logger.Config{
				SlowThreshold: time.Second,                    // 慢查詢閾值（1秒）
				Colorful:      true,                          // 啟用彩色日誌輸出
				LogLevel: logger.LogLevel(                    // 設定日誌等級
					log.DefaultLogger.Options().Level.LevelForGorm()),
			},
		),
	}, opens[c.Driver]) // 使用對應的資料庫驅動

	// 檢查資料庫連接結果
	if err != nil {
		// 連接失敗，記錄錯誤並終止程式
		log.Fatal(pkg.Red(c.Driver+" 資料庫連接失敗: "), err)
	} else {
		// 連接成功，記錄成功訊息
		log.Info(pkg.Green(c.Driver + " 資料庫連接成功!"))
	}

	// 初始化 Casbin 權限控制器
	// 用於實作基於角色的存取控制（RBAC）
	e := mycasbin.Setup(db, "")

	// 將資料庫實例和 Casbin 執行器註冊到 SDK 執行時環境
	// 供其他模組使用
	sdk.Runtime.SetDb(host, db)      // 註冊資料庫實例
	sdk.Runtime.SetCasbin(host, e)   // 註冊權限控制器
}
