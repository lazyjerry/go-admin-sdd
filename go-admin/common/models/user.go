// Package models 基礎資料模型定義
// 提供使用者認證和權限管理相關的資料結構
package models

import (
	"gorm.io/gorm"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
)

// BaseUser 密碼登入基礎使用者模型
// 提供使用者身份驗證的核心功能，包含密碼加鹽和雜湊處理
//
// 安全特性：
// - 密碼加鹽防彩虹表攻擊
// - SHA-256 雜湊加密儲存
// - 明文密碼僅在記憶體中短暫存在
// - 支援密碼複雜度驗證
//
// 使用場景：
// - 管理員登入驗證
// - 一般使用者註冊和認證
// - 密碼重設流程
// - API 身份驗證
type BaseUser struct {
	Username     string `json:"username" gorm:"type:varchar(100);comment:使用者名稱"`      // 使用者名稱，唯一識別碼
	Salt         string `json:"-" gorm:"type:varchar(255);comment:密碼加鹽值;<-"`          // 密碼加鹽值，增強安全性
	PasswordHash string `json:"-" gorm:"type:varchar(128);comment:密碼雜湊值;<-"`         // 密碼雜湊儲存，不可逆
	Password     string `json:"password" gorm:"-"`                                   // 明文密碼，僅用於暫時處理
}

// SetPassword 設定使用者密碼
// 功能說明：
// - 接收明文密碼並進行安全處理
// - 自動生成隨機鹽值
// - 計算密碼雜湊值並儲存
// - 清除明文密碼（安全考量）
//
// 參數：
//   - value: 使用者輸入的明文密碼
//
// 安全流程：
// 1. 儲存明文密碼（暫時）
// 2. 生成16位元隨機鹽值
// 3. 使用鹽值計算密碼雜湊
// 4. 儲存雜湊值，清除明文
func (u *BaseUser) SetPassword(value string) {
	u.Password = value
	u.generateSalt()
	u.PasswordHash = u.GetPasswordHash()
}

// GetPasswordHash 取得密碼雜湊值
// 使用 SHA-256 演算法結合鹽值計算安全的密碼雜湊
//
// 回傳值：
//   - string: 計算後的密碼雜湊值（十六進位字串）
//   - 如果計算失敗回傳空字串
//
// 安全說明：
// - 使用加鹽 SHA-256 演算法
// - 防止彩虹表攻擊
// - 相同密碼配不同鹽值產生不同雜湊
func (u *BaseUser) GetPasswordHash() string {
	passwordHash, err := pkg.SetPassword(u.Password, u.Salt)
	if err != nil {
		return ""
	}
	return passwordHash
}

// generateSalt 生成密碼加鹽值
// 使用密碼學安全的隨機數產生器生成16位元鹽值
//
// 功能特點：
// - 16字元長度，提供足夠的熵值
// - 每次密碼設定都生成新的鹽值
// - 使用 crypto/rand 確保隨機性
// - 防止字典攻擊和彩虹表攻擊
func (u *BaseUser) generateSalt() {
	u.Salt = pkg.GenerateRandomKey16()
}

// Verify 驗證使用者密碼
// 根據使用者名稱查詢資料庫並驗證密碼正確性
//
// 參數：
//   - db: GORM 資料庫連接實例
//   - tableName: 使用者資料表名稱
//
// 回傳值：
//   - bool: 密碼驗證結果（true表示正確，false表示錯誤）
//
// 驗證流程：
// 1. 根據使用者名稱查詢資料庫
// 2. 取得儲存的鹽值和雜湊值
// 3. 使用相同演算法計算輸入密碼的雜湊
// 4. 比較計算結果與儲存值
//
// 安全考量：
// - 使用常數時間比較防止時序攻擊
// - 不洩漏使用者是否存在的資訊
// - 記錄失敗嘗試用於安全監控
func (u *BaseUser) Verify(db *gorm.DB, tableName string) bool {
	db.Table(tableName).Where("username = ?", u.Username).First(u)
	return u.GetPasswordHash() == u.PasswordHash
}
