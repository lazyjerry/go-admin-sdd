package global

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/api"
)

// LoadPolicy 方法用於加載 Casbin 的策略。
// 它會從當前請求的 Host 中獲取策略並進行初始化。
func LoadPolicy(c *gin.Context) (*casbin.SyncedEnforcer, error) {
	log := api.GetRequestLogger(c)
	if err := sdk.Runtime.GetCasbinKey(c.Request.Host).LoadPolicy(); err == nil {
		return sdk.Runtime.GetCasbinKey(c.Request.Host), err
	} else {
		log.Errorf("casbin rbac_model or policy init error, %s ", err.Error())
		return nil, err
	}
}
