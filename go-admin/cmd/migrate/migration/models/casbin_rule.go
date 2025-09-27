package models

// CasbinRule 對應資料表 `sys_casbin_rule`，用於存放 Casbin 權限規則。
// 欄位說明：
// - ID: 主鍵，自動遞增
// - Ptype: 規則類型（如 p, g 等）
// - V0..V5: 規則參數（Casbin 支援的多個欄位），用於存放 subject/obj/action 等資訊
// 注意：Ptype 及 V0..V5 使用 unique_index 組合以避免重複規則。
type CasbinRule struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Ptype string `gorm:"size:512;uniqueIndex:unique_index"`
	V0    string `gorm:"size:512;uniqueIndex:unique_index"`
	V1    string `gorm:"size:512;uniqueIndex:unique_index"`
	V2    string `gorm:"size:512;uniqueIndex:unique_index"`
	V3    string `gorm:"size:512;uniqueIndex:unique_index"`
	V4    string `gorm:"size:512;uniqueIndex:unique_index"`
	V5    string `gorm:"size:512;uniqueIndex:unique_index"`
}

func (CasbinRule) TableName() string {
	return "sys_casbin_rule"
}
