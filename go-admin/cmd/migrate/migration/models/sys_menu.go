package models

// SysMenu 表示系統功能選單（或權限項），對應資料表 `sys_menu`。
// 主要欄位說明：
// - MenuId: 選單編號（主鍵）
// - MenuName/Title: 選單名稱與標題
// - Path/Component: 前端路由與對應組件
// - Permission: 權限字串，用於後端授權檢查
// - ParentId: 父節點 ID 用於樹狀結構
// - NoCache/Breadcrumb/Visible/IsFrame: 前端顯示與快取設定
// - SysApi: 與 API 規則關聯（多對多）
// 同時包含 ControlBy 與 ModelTime 用於記錄建立/更新者與時間。

type SysMenu struct {
	MenuId     int       `json:"menuId" gorm:"primaryKey;autoIncrement"`
	MenuName   string    `json:"menuName" gorm:"size:128;"`
	Title      string    `json:"title" gorm:"size:128;"`
	Icon       string    `json:"icon" gorm:"size:128;"`
	Path       string    `json:"path" gorm:"size:128;"`
	Paths      string    `json:"paths" gorm:"size:128;"`
	MenuType   string    `json:"menuType" gorm:"size:1;"`
	Action     string    `json:"action" gorm:"size:16;"`
	Permission string    `json:"permission" gorm:"size:255;"`
	ParentId   int       `json:"parentId" gorm:"size:11;"`
	NoCache    bool      `json:"noCache" gorm:"size:8;"`
	Breadcrumb string    `json:"breadcrumb" gorm:"size:255;"`
	Component  string    `json:"component" gorm:"size:255;"`
	Sort       int       `json:"sort" gorm:"size:4;"`
	Visible    string    `json:"visible" gorm:"size:1;"`
	IsFrame    string    `json:"isFrame" gorm:"size:1;DEFAULT:0;"`
	SysApi     []SysApi  `json:"sysApi" gorm:"many2many:sys_menu_api_rule"`
	ControlBy
	ModelTime
}

func (SysMenu) TableName() string {
	return "sys_menu"
}