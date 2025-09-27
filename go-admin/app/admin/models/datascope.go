package models

import (
	"errors"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"
	"gorm.io/gorm"

	log "github.com/go-admin-team/go-admin-core/logger"
	"github.com/go-admin-team/go-admin-core/sdk/config"
)

// DataPermission 表示使用者的資料權限範圍資訊，用於在查詢時套用資料範圍過濾。
// 欄位：
// - DataScope: 權限範圍代碼（例如：2=自定義角色範圍, 3=部門, 4=部門及其子部門, 5=個人）
// - UserId / DeptId / RoleId: 用於生成查詢條件的識別資訊
type DataPermission struct {
	DataScope string
	UserId    int
	DeptId    int
	RoleId    int
}

// GetDataScope 根據當前 DataPermission 與資料表名稱生成 GORM 查詢條件（db），
// 並返回修改後的 db。若應用配置關閉資料權限（EnableDP=false），會直接回傳原始 db。
func (e *DataPermission) GetDataScope(tableName string, db *gorm.DB) (*gorm.DB, error) {

	if !config.ApplicationConfig.EnableDP {
		usageStr := `数据权限已经为您` + pkg.Green(`关闭`) + `，如需开启请参考配置文件字段说明`
		log.Debug("%s\n", usageStr)
		return db, nil
	}
	user := new(SysUser)
	role := new(SysRole)
	err := db.Find(user, e.UserId).Error
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	err = db.Find(role, user.RoleId).Error
	if err != nil {
		return nil, errors.New("获取用户数据出错 msg:" + err.Error())
	}
	if role.DataScope == "2" {
		db = db.Where(tableName+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
	}
	if role.DataScope == "3" {
		db = db.Where(tableName+".create_by in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
	}
	if role.DataScope == "4" {
		db = db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+pkg.IntToString(user.DeptId)+"%")
	}
	if role.DataScope == "5" || role.DataScope == "" {
		db = db.Where(tableName+".create_by = ?", e.UserId)
	}

	return db, nil
}

//func DataScopes(tableName string, userId int) func(db *gorm.DB) *gorm.DB {
//	return func(db *gorm.DB) *gorm.DB {
//		user := new(SysUser)
//		role := new(SysRole)
//		user.UserId = userId
//		err := db.Find(user, userId).Error
//		if err != nil {
//			db.Error = errors.New("获取用户数据出错 msg:" + err.Error())
//			return db
//		}
//		err = db.Find(role, user.RoleId).Error
//		if err != nil {
//			db.Error = errors.New("获取用户数据出错 msg:" + err.Error())
//			return db
//		}
//		if role.DataScope == "2" {
//			return db.Where(tableName+".create_by in (select sys_user.user_id from sys_role_dept left join sys_user on sys_user.dept_id=sys_role_dept.dept_id where sys_role_dept.role_id = ?)", user.RoleId)
//		}
//		if role.DataScope == "3" {
//			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where dept_id = ? )", user.DeptId)
//		}
//		if role.DataScope == "4" {
//			return db.Where(tableName+".create_by in (SELECT user_id from sys_user where sys_user.dept_id in(select dept_id from sys_dept where dept_path like ? ))", "%"+pkg.IntToString(user.DeptId)+"%")
//		}
//		if role.DataScope == "5" || role.DataScope == "" {
//			return db.Where(tableName+".create_by = ?", userId)
//		}
//		return db
//	}
//}
