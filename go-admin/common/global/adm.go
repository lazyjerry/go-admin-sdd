// Package global 全域設定和常數定義
// 包含系統版本資訊、資料庫驅動配置等全域變數
package global

const (
	// Version Go-Admin 系統版本號
	// 用於：
	// - API 版本標識
	// - 系統升級檢查
	// - 相容性驗證
	// - 日誌記錄和除錯
	Version = "2.2.0"
)

var (
	// Driver 資料庫驅動類型
	// 支援的驅動包括：
	// - mysql: MySQL/MariaDB 資料庫
	// - postgres: PostgreSQL 資料庫  
	// - sqlite3: SQLite 資料庫（開發和測試用）
	// - sqlserver: Microsoft SQL Server 資料庫
	// - oracle: Oracle 資料庫（企業版）
	//
	// 該變數在系統初始化時設定，影響：
	// - 資料庫連接配置
	// - SQL 語法差異處理
	// - 資料類型轉換
	// - 效能最佳化策略
	Driver string
)
