package actions

import (
	"errors"

	"github.com/gin-gonic/gin"
	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/jwtauth/user"
	"github.com/go-admin-team/go-admin-core/sdk/pkg/response"
	"gorm.io/gorm"
)

type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

// DataPermission 結構代表使用者的資料權限資訊，用於在查詢時套用資料範圍限制。
// 欄位：
// - DataScope: 權限範圍代碼（例如：全局、部門、個人等）
// - UserId, DeptId, RoleId: 用於生成範圍查詢條件的識別資訊
//
// PermissionAction 中介軟體：
// - 於請求進入時從 DB 查詢並構造 DataPermission，存放於 Context（Key: PermissionKey）
// - 後續查詢會呼叫 Permission 函式將該權限轉為 GORM scopes
func PermissionAction() gin.HandlerFunc {
	return func(c *gin.Context) {
		db, err := pkg.GetOrm(c)
		if err != nil {
			log.Error(err)
			return
		}

		msgID := pkg.GenerateMsgIDFromContext(c)
		var p = new(DataPermission)
		if userId := user.GetUserIdStr(c); userId != "" {
			p, err = newDataPermission(db, userId)
			if err != nil {
				log.Errorf("MsgID[%s] PermissionAction error: %s", msgID, err)
				response.Error(c, 500, err, "權限範圍鑑定錯誤")
				c.Abort()
				return
			}
		}
		c.Set(PermissionKey, p)
		c.Next()
	}
}

func newDataPermission(tx *gorm.DB, userId interface{}) (*DataPermission, error) {
	var err error
	p := &DataPermission{}

	err = tx.Table("sys_user").
		Select("sys_user.user_id", "sys_role.role_id", "sys_user.dept_id", "sys_role.data_scope").
		Joins("left join sys_role on sys_role.role_id = sys_user.role_id").
		Where("sys_user.user_id = ?", userId).
		Scan(p).Error
	if err != nil {
		err = errors.New("獲取使用者資料出錯 msg:" + err.Error())
		return nil, err
	}
	return p, nil
}

// Permission 為 GORM scope 生成器，根據 DataPermission 的 DataScope 生成對應的 WHERE 條件。
func Permission(tableName string, p *DataPermission) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if !config.ApplicationConfig.EnableDP {
			return db
		}
		switch p.DataScope {
		case "2":
			return db.Where(tableName+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", p.RoleId)
		case "3":
			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where dept_id = ? )", p.DeptId)
		case "4":
			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%/"+pkg.IntToString(p.DeptId)+"/%")
		case "5":
			return db.Where(tableName+".create_by = ?", p.UserId)
		default:
			return db
		}
	}
}

func getPermissionFromContext(c *gin.Context) *DataPermission {
	p := new(DataPermission)
	if pm, ok := c.Get(PermissionKey); ok {
		switch pm.(type) {
		case *DataPermission:
			p = pm.(*DataPermission)
		}
	}
	return p
}

// GetPermissionFromContext 提供非action寫法資料範圍約束
func GetPermissionFromContext(c *gin.Context) *DataPermission {
	return getPermissionFromContext(c)
}
