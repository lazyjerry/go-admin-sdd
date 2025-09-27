package models

import (
	"time"

	"gorm.io/gorm"
)

// ControlBy 表示記錄建立者與更新者的欄位，用於追蹤資料的建立與更新人員。
// 欄位：
// - CreateBy: 建立者的使用者 ID
// - UpdateBy: 更新者的使用者 ID
type ControlBy struct {
	CreateBy int `json:"createBy" gorm:"index;comment:創建者"`
	UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

// Model 為基礎模型，包含主鍵 Id。
// Id: 主鍵編號，自動遞增。
type Model struct {
	Id int `json:"id" gorm:"primaryKey;autoIncrement;comment:主鍵編碼"`
}

// ModelTime 包含常用的時間欄位，用於追蹤記錄的建立/更新/刪除時間。
// - CreatedAt: 建立時間
// - UpdatedAt: 最後更新時間
// - DeletedAt: 軟刪除時間（GORM 的 DeletedAt）
type ModelTime struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"comment:創建時間"`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"comment:最後更新时间"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index;comment:刪除時間"`
}
