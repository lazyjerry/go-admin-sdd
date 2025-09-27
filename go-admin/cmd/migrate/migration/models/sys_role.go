package models

// SysRole 表示系統角色，控制使用者的權限集合。
// 主要欄位說明：
// - RoleId: 角色編號（主鍵）
// - RoleName: 角色名稱
// - RoleKey: 角色代碼（程式中用於權限判斷）
// - Admin: 是否為管理員角色（特殊處理）
// - DataScope: 資料範圍設定，用於資料權限檢查（例如全域、部門、個人等）
// - SysMenu: 與選單的多對多關聯，代表該角色可存取的選單
// 同時包含 ControlBy 與 ModelTime。
type SysRole struct {
	RoleId    int       `json:"roleId" gorm:"primaryKey;autoIncrement"` // 角色編號
	RoleName  string    `json:"roleName" gorm:"size:128;"`              // 角色名稱
	Status    string    `json:"status" gorm:"size:4;"`                  //
	RoleKey   string    `json:"roleKey" gorm:"size:128;"`               //角色代碼
	RoleSort  int       `json:"roleSort" gorm:""`                       //角色排序
	Flag      string    `json:"flag" gorm:"size:128;"`                  //
	Remark    string    `json:"remark" gorm:"size:255;"`                //備註
	Admin     bool      `json:"admin" gorm:"size:4;"`
	DataScope string    `json:"dataScope" gorm:"size:128;"`
	SysMenu   []SysMenu `json:"sysMenu" gorm:"many2many:sys_role_menu;foreignKey:RoleId;joinForeignKey:role_id;references:MenuId;joinReferences:menu_id;"`
	ControlBy
	ModelTime
}

func (SysRole) TableName() string {
	return "sys_role"
}