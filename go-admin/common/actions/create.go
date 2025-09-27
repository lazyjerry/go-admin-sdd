package actions

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-admin-team/go-admin-core/sdk/api"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"

	"go-admin/common/dto"
	"go-admin/common/models"
)

// CreateAction 為通用的新增處理器。
// 流程：
// 1. 從上下文取得 logger 與 DB 連線
// 2. 生成 Request DTO 並綁定請求參數
// 3. 生成對應的模型實例並設定建立者
// 4. 執行資料庫新增操作，處理錯誤並回傳結果
func CreateAction(control dto.Control) gin.HandlerFunc {
	return func(c *gin.Context) {
		log := api.GetRequestLogger(c)
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		//新增操作
		req := control.Generate()
		err = req.Bind(c)
		if err != nil {
			response.Error(c, http.StatusUnprocessableEntity, err, err.Error())
			return
		}
		var object models.ActiveRecord
		object, err = req.GenerateM()
		if err != nil {
			response.Error(c, 500, err, "模型生成失敗")
			return
		}
		object.SetCreateBy(user.GetUserId(c))
		err = db.WithContext(c).Create(object).Error
		if err != nil {
			log.Errorf("Create error: %s", err)
			response.Error(c, 500, err, "創建失敗")
			return
		}
		response.OK(c, object.GetId(), "創建成功")
		c.Next()
	}
}
