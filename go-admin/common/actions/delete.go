package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/common/dto"
	"go-admin/common/models"
)

// DeleteAction 為通用的刪除處理器。
// 流程：
// 1. 取得 DB 並生成請求 DTO
// 2. 綁定參數並生成模型實例，設定更新者
// 3. 執行資料權限檢查並呼叫 Delete
// 4. 根據 RowsAffected 判斷是否有刪除權限並回傳結果
func DeleteAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		//刪除操作
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			log.Errorf("MsgID[%s] Bind error: %s", msgID, err)
			response.Error(c, http.StatusUnprocessableEntity, err, "參數驗證失敗")
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			response.Error(c, 500, err, "模型生成失敗")
			return
		}

		object.SetUpdateBy(user.GetUserId(c))

		//資料權限檢查
		p := GetPermissionFromContext(c)

		db = db.WithContext(c).Scopes(
			Permission(object.TableName(), p),
		).Where(req.GetId()).Delete(object)
		if err = db.Error; err != nil {
			log.Errorf("MsgID[%s] Delete error: %s", msgID, err)
			response.Error(c, 500, err, "刪除失敗")
			return
		}
		if db.RowsAffected == 0 {
			response.Error(c, http.StatusForbidden, nil, "無權刪除該資料")
			return
		}
		response.OK(c, object.GetId(), "刪除成功")
		c.Next()
	}
}
