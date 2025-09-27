// Package models 系統角色模型定義
// 定義 RBAC 權限系統中的角色資料結構和操作方法
package models

import "go-admin/common/models"

// SysRole 系統角色模型
// 基於 RBAC（基於角色的存取控制）設計的核心角色實體
//
// 角色功能：
// - 權限分組和管理的基本單位
// - 支援階層式角色設計
// - 資料權限範圍控制
// - 選單和功能權限分配
//
// 權限特性：
// - 支援多對多關聯（角色-選單、角色-部門）
// - 資料權限範圍可配置
// - 支援管理員標識
// - 靈活的權限繼承機制
type SysRole struct {
	// 基本識別資訊
	RoleId    int        `json:"roleId" gorm:"primaryKey;autoIncrement;comment:角色唯一識別碼"` 
	RoleName  string     `json:"roleName" gorm:"size:128;comment:角色顯示名稱"`
	RoleKey   string     `json:"roleKey" gorm:"size:128;comment:角色代碼，用於程式識別"`
	RoleSort  int        `json:"roleSort" gorm:"comment:角色排序，數字越小優先級越高"`
	
	// 狀態和標識
	Status    string     `json:"status" gorm:"size:4;comment:角色狀態（1-停用 2-啟用）"`
	Flag      string     `json:"flag" gorm:"size:128;comment:角色標記，用於特殊標識"`
	Admin     bool       `json:"admin" gorm:"size:4;comment:是否為管理員角色"`
	
	// 權限控制
	DataScope string     `json:"dataScope" gorm:"size:128;comment:資料權限範圍（1-全部 2-自訂 3-部門 4-部門及子部門 5-僅本人）"`
	
	// 描述資訊
	Remark    string     `json:"remark" gorm:"size:255;comment:角色描述和備註"`
	
	// 前端需要的輔助欄位（不儲存到資料庫）
	Params    string     `json:"params" gorm:"-" comment:"前端傳遞的額外參數"`
	MenuIds   []int      `json:"menuIds" gorm:"-" comment:"關聯的選單ID陣列"`
	DeptIds   []int      `json:"deptIds" gorm:"-" comment:"資料權限部門ID陣列"`
	
	// 多對多關聯關係
	SysDept   []SysDept  `json:"sysDept" gorm:"many2many:sys_role_dept;foreignKey:RoleId;joinForeignKey:role_id;references:DeptId;joinReferences:dept_id;" comment:"角色關聯的部門列表"`
	SysMenu   *[]SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;" comment:"角色關聯的選單權限"`
	
	// 嵌入通用欄位
	models.ControlBy  // 建立人、更新人等控制欄位
	models.ModelTime  // 建立時間、更新時間等時間戳欄位
}

// TableName 指定資料表名稱
// 實作 GORM 的 Tabler 介面，自訂資料表名稱
func (*SysRole) TableName() string {
	return "sys_role"
}

// Generate 程式碼生成器支援方法
// 實作 ActiveRecord 介面，用於程式碼生成工具
// 返回當前角色物件的副本，避免指標問題
func (e *SysRole) Generate() models.ActiveRecord {
	o := *e
	return &o
}

// GetId 取得主鍵值
// 實作 ActiveRecord 介面，返回角色的主鍵 RoleId
// 用於通用的 CRUD 操作中識別角色記錄
func (e *SysRole) GetId() interface{} {
	return e.RoleId
}
