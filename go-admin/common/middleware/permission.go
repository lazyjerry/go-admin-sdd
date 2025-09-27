// Package middleware RBAC 權限檢查中間件
// 基於 Casbin 實現的細粒度權限控制系統
package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2/util"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
)

// AuthCheckRole RBAC 權限檢查中間件
// 基於 Casbin 框架實現精細的角色權限控制，支援：
// - 基於角色的存取控制（RBAC）
// - URL 路徑和 HTTP 方法的權限驗證
// - 管理員特殊權限處理
// - 白名單機制（排除特定 API 的權限檢查）
// - 詳細的權限檢查日誌記錄
//
// 權限檢查流程：
// 1. 從 JWT Token 中提取使用者角色資訊
// 2. 檢查是否為管理員角色（直接放行）
// 3. 檢查當前請求是否在白名單中
// 4. 使用 Casbin 執行權限策略驗證
// 5. 記錄權限檢查結果並決定放行或拒絕
//
// 回傳值：
//   - gin.HandlerFunc: Gin 框架的中間件處理函數
//
// 權限策略配置：
// - 支援角色繼承和權限傳遞
// - 可配置的權限策略檔案
// - 動態權限策略更新
// - 細粒度的 API 級別控制
func AuthCheckRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 取得請求日誌記錄器，用於記錄權限檢查過程
		log := api.GetRequestLogger(c)
		
		// 從 JWT Token 中提取使用者身份資訊
		// JWT 載荷包含 rolekey（角色識別碼）等重要身份資訊
		data, _ := c.Get(jwtauth.JwtPayloadKey)
		v := data.(jwtauth.MapClaims)
		
		// 取得 Casbin 權限執行器實例
		// 根據請求主機名取得對應的權限策略執行器
		e := sdk.Runtime.GetCasbinKey(c.Request.Host)
		var res, casbinExclude bool
		var err error
		
		// 管理員特殊權限處理
		// admin 角色擁有所有 API 的存取權限，直接放行
		if v["rolekey"] == "admin" {
			log.Infof("管理員角色存取，直接放行 - 角色: admin, 方法: %s, 路徑: %s", c.Request.Method, c.Request.URL.Path)
			res = true
			c.Next()
			return
		}
		
		// 白名單機制檢查
		// 某些 API 端點（如公開資源、健康檢查等）不需要權限驗證
		for _, i := range CasbinExclude {
			if util.KeyMatch2(c.Request.URL.Path, i.Url) && c.Request.Method == i.Method {
				casbinExclude = true
				break
			}
		}
		
		// 如果在白名單中，跳過權限檢查
		if casbinExclude {
			log.Infof("白名單排除，無需權限驗證 - 方法: %s, 路徑: %s", c.Request.Method, c.Request.URL.Path)
			c.Next()
			return
		}
		
		// 執行 Casbin 權限策略檢查
		// 驗證當前角色是否有權限存取指定的 API 端點和 HTTP 方法
		res, err = e.Enforce(v["rolekey"], c.Request.URL.Path, c.Request.Method)
		if err != nil {
			log.Errorf("權限檢查發生錯誤 - 錯誤: %s, 方法: %s, 路徑: %s", err, c.Request.Method, c.Request.URL.Path)
			response.Error(c, 500, err, "權限驗證系統錯誤")
			return
		}

		// 根據權限檢查結果決定是否放行
		if res {
			// 權限驗證通過，記錄成功日誌並繼續處理請求
			log.Infof("權限驗證通過 - 結果: %v, 角色: %s, 方法: %s, 路徑: %s", res, v["rolekey"], c.Request.Method, c.Request.URL.Path)
			c.Next()
		} else {
			// 權限驗證失敗，記錄警告日誌並回傳403錯誤
			log.Warnf("權限驗證失敗 - 結果: %v, 角色: %s, 方法: %s, 路徑: %s, 訊息: 當前請求無權限，請管理員確認", res, v["rolekey"], c.Request.Method, c.Request.URL.Path)
			c.JSON(http.StatusOK, gin.H{
				"code": 403,
				"msg":  "很抱歉，您沒有該介面的存取權限，請聯絡管理員",
			})
			c.Abort()
			return
		}
	}
}
