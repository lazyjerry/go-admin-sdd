package actions

import (
	"errors"
	"net/http"

	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"

	"go-admin/common/dto"
	"go-admin/common/models"
)

// ViewAction 是一個通用的資料詳情查詢處理器。
// 它的主要流程：
// 1. 從上下文取得 DB 連線
// 2. 生成請求 DTO 並綁定參數
// 3. 生成模型實例並檢查資料權限
// 4. 執行查詢並回傳結果
func ViewAction(control dto.Control, f func() interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		//查看詳情
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, "參數驗證失敗")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			response.Error(c, 500, err, "模型生成失敗")
			return
		}

		var rsp interface{}
		if f != nil {
			rsp = f()
		} else {
			rsp, _ = req.GenerateM()
		}

		//資料權限檢查
		p := GetPermissionFromContext(c)

		err = db.Model(object).WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).First(rsp).Error

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, http.StatusNotFound, nil, "查看物件不存在或無權查看")
			return
		}
		if err != nil {
			log.Errorf("MsgID[%s] View error: %s", msgID, err)
			response.Error(c, 500, err, "查看失敗")
			return
		}
		response.OK(c, rsp, "查詢成功")
		c.Next()
	}
}
