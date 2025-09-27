package common

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// GetClientIP 取得客戶端 IP 的工具函式。
// 它按照以下優先順序檢查請求標頭/上下文：
// 1. X-Forwarded-For (通常由反向代理設置)
// 2. X-Real-IP
// 3. RemoteIP (Gin 提供)
// 4. ClientIP (Gin 提供，會嘗試多種來源)
// 若都無法取得，則回傳 127.0.0.1 作為兜底值。
func GetClientIP(c *gin.Context) string {
	// 優先從 X-Forwarded-For 獲取 IP
	ip := c.Request.Header.Get("X-Forwarded-For")
	if ip == "" || strings.Contains(ip, "127.0.0.1") {
		// 如果為空或為本地地址，則嘗試從 X-Real-IP 獲取
		ip = c.Request.Header.Get("X-real-ip")
	}
	if ip == "" {
		// 如果仍然為空，則使用 RemoteIP
		ip = c.RemoteIP()
	}
	if ip == "" || ip == "127.0.0.1" {
		// 如果仍然為空或為本地地址，則使用 ClientIP
		ip = c.ClientIP()
	}
	if ip == "" {
		// 最後兜底為本地地址
		ip = "127.0.0.1"
	}
	return ip
}
