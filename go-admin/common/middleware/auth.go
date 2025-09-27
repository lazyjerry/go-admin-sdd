// Package middleware JWT 身份認證中間件
// 提供基於 JWT (JSON Web Token) 的使用者身份驗證和授權功能
package middleware

import (
	"time"

	"go-admin/common/middleware/handler"

	"github.com/go-admin-team/go-admin-core/sdk/config"
	jwt "github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
)

// AuthInit 初始化 JWT 身份認證中間件
// 建立並配置 JWT 中間件實例，用於保護需要身份驗證的 API 端點
//
// JWT 功能特性：
// - 無狀態身份驗證（不依賴 Session）
// - 支援 Token 自動續期
// - 多種 Token 傳遞方式（Header、Query、Cookie）
// - 開發環境特殊配置（長期有效 Token）
// - 生產環境安全配置（短期有效 Token）
//
// 安全考量：
// - 使用 HMAC-SHA256 演算法簽名
// - 可配置的 Token 過期時間
// - 支援 Token 黑名單機制
// - 防 CSRF 攻擊保護
//
// 回傳值：
//   - *jwt.GinJWTMiddleware: 已配置的 JWT 中間件實例
//   - error: 初始化過程中的錯誤（如配置無效等）
//
// 使用範例：
//   authMiddleware, err := AuthInit()
//   if err != nil {
//     log.Fatal("JWT初始化失敗:", err)
//   }
//   router.Use(authMiddleware.MiddlewareFunc())
func AuthInit() (*jwt.GinJWTMiddleware, error) {
	// 設定 Token 過期時間
	// 開發環境：設定極長時間避免頻繁重新登入影響開發效率
	// 生產環境：根據配置設定安全的過期時間
	timeout := time.Hour
	if config.ApplicationConfig.Mode == "dev" {
		// 開發模式：設定約100年的超長過期時間（876010小時）
		// 注意：這僅用於開發便利性，生產環境絕不可使用
		timeout = time.Duration(876010) * time.Hour
	} else {
		// 生產模式：使用配置檔案中指定的過期時間
		// 如果未配置則使用預設的1小時
		if config.JwtConfig.Timeout != 0 {
			timeout = time.Duration(config.JwtConfig.Timeout) * time.Second
		}
	}

	// 建立並配置 JWT 中間件
	return jwt.New(&jwt.GinJWTMiddleware{
		Realm:           "Go-Admin System",                                    // JWT 領域標識，用於描述受保護的範圍
		Key:             []byte(config.JwtConfig.Secret),                      // JWT 簽名密鑰，從配置檔案讀取
		Timeout:         timeout,                                             // Token 有效期限
		MaxRefresh:      time.Hour,                                          // Token 最大續期時間
		PayloadFunc:     handler.PayloadFunc,                                // Token 載荷構建函數
		IdentityHandler: handler.IdentityHandler,                            // 身份識別處理函數
		Authenticator:   handler.Authenticator,                              // 身份認證處理函數（登入驗證）
		Authorizator:    handler.Authorizator,                               // 授權驗證處理函數（權限檢查）
		Unauthorized:    handler.Unauthorized,                               // 未授權處理函數（錯誤回應）
		TokenLookup:     "header: Authorization, query: token, cookie: jwt", // Token 查找方式：優先從 Header，其次 Query，最後 Cookie
		TokenHeadName:   "Bearer",                                           // Header 中 Token 的前綴名稱
		TimeFunc:        time.Now,                                           // 時間函數，用於 Token 時效性檢查
	})
}