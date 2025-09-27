// Package models 系統使用者模型定義
// 定義系統使用者的資料結構和相關操作方法
package models

import (
	"go-admin/common/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SysUser 系統使用者模型
// 用於管理系統中的所有使用者帳號資訊
//
// 功能特性：
// - 支援密碼自動加密
// - 整合角色權限系統
// - 支援部門和崗位關聯
// - 提供完整的使用者生命週期管理
//
// 安全特性：
// - 密碼欄位自動排除於 JSON 序列化
// - 使用 bcrypt 演算法加密密碼
// - 支援鹽值增強安全性
type SysUser struct {
	// 基本識別資訊
	UserId   int      `gorm:"primaryKey;autoIncrement;comment:使用者唯一識別碼"  json:"userId"`
	Username string   `json:"username" gorm:"size:64;comment:登入使用者名稱，系統內唯一"`
	Password string   `json:"-" gorm:"size:128;comment:使用者密碼，bcrypt 加密儲存"`
	NickName string   `json:"nickName" gorm:"size:128;comment:使用者暱稱或顯示名稱"`
	
	// 聯絡資訊
	Phone    string   `json:"phone" gorm:"size:11;comment:手機號碼"`
	Email    string   `json:"email" gorm:"size:128;comment:電子郵件地址"`
	
	// 系統角色和組織架構
	RoleId   int      `json:"roleId" gorm:"size:20;comment:關聯的角色ID"`
	DeptId   int      `json:"deptId" gorm:"size:20;comment:所屬部門ID"`
	PostId   int      `json:"postId" gorm:"size:20;comment:擔任崗位ID"`
	
	// 安全和個人化設定
	Salt     string   `json:"-" gorm:"size:255;comment:密碼加鹽值，增強安全性"`
	Avatar   string   `json:"avatar" gorm:"size:255;comment:使用者頭像URL"`
	Sex      string   `json:"sex" gorm:"size:255;comment:性別"`
	Remark   string   `json:"remark" gorm:"size:255;comment:備註資訊"`
	Status   string   `json:"status" gorm:"size:4;comment:帳號狀態（1-啟用 2-停用）"`
	
	// 前端需要的輔助欄位（不儲存到資料庫）
	DeptIds  []int    `json:"deptIds" gorm:"-" comment:"部門ID陣列，用於前端多選"`
	PostIds  []int    `json:"postIds" gorm:"-" comment:"崗位ID陣列，用於前端多選"`
	RoleIds  []int    `json:"roleIds" gorm:"-" comment:"角色ID陣列，用於前端多選"`
	
	// 關聯模型（延遲載入）
	Dept     *SysDept `json:"dept" comment:"關聯的部門物件"`
	
	// 嵌入通用欄位
	models.ControlBy  // 建立人、更新人等控制欄位
	models.ModelTime  // 建立時間、更新時間等時間戳欄位
}

// TableName 指定資料表名稱
// 實作 GORM 的 Tabler 介面，自訂資料表名稱
func (*SysUser) TableName() string {
	return "sys_user"
}

// Generate 程式碼生成器支援方法
// 實作 ActiveRecord 介面，用於程式碼生成工具
// 返回當前物件的副本，避免指標問題
func (e *SysUser) Generate() models.ActiveRecord {
	o := *e
	return &o
}

// GetId 取得主鍵值
// 實作 ActiveRecord 介面，返回使用者的主鍵 UserId
// 用於通用的 CRUD 操作中識別記錄
func (e *SysUser) GetId() interface{} {
	return e.UserId
}

// Encrypt 密碼加密方法
// 使用 bcrypt 演算法對使用者密碼進行安全加密
// 
// 安全特性：
// - 使用 bcrypt.DefaultCost 預設強度
// - 自動生成隨機鹽值
// - 抗彩虹表攻擊
// - 計算複雜度可調整
// 
// 返回值：
//   - error: 加密過程中的錯誤，nil 表示成功
func (e *SysUser) Encrypt() (err error) {
	// 如果密碼為空，跳過加密
	if e.Password == "" {
		return
	}

	// 使用 bcrypt 生成密碼雜湊值
	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(e.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		// 將雜湊結果儲存回密碼欄位
		e.Password = string(hash)
		return
	}
}

// BeforeCreate GORM 建立前鉤子方法
// 在資料庫建立記錄之前自動執行
// 主要用於密碼加密處理
func (e *SysUser) BeforeCreate(_ *gorm.DB) error {
	return e.Encrypt()
}

// BeforeUpdate GORM 更新前鉤子方法  
// 在資料庫更新記錄之前自動執行
// 只有當密碼欄位不為空時才進行加密
// 避免空密碼覆蓋現有密碼的情況
func (e *SysUser) BeforeUpdate(_ *gorm.DB) error {
	var err error
	if e.Password != "" {
		err = e.Encrypt()
	}
	return err
}

// AfterFind GORM 查詢後鉤子方法
// 在從資料庫查詢到記錄後自動執行
// 用於設定前端需要的輔助欄位
// 
// 功能：
// - 將單一值轉換為陣列格式
// - 便於前端表單處理
// - 保持資料一致性
func (e *SysUser) AfterFind(_ *gorm.DB) error {
	// 將單一部門ID轉換為陣列（前端多選元件需要）
	e.DeptIds = []int{e.DeptId}
	
	// 將單一崗位ID轉換為陣列（前端多選元件需要）
	e.PostIds = []int{e.PostId}
	
	// 將單一角色ID轉換為陣列（前端多選元件需要）
	e.RoleIds = []int{e.RoleId}
	
	return nil
}
