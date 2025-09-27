package models

// SysPost 表示系統中的職位（崗位）資訊。
// 欄位：
// - PostId: 崗位編號（主鍵）
// - PostName: 崗位名稱
// - PostCode: 崗位代碼
// - Sort: 顯示排序
// - Status: 啟用/停用狀態
// - Remark: 備註
// 同時包含 ControlBy 與 ModelTime。
type SysPost struct {
	PostId   int    `gorm:"primaryKey;autoIncrement" json:"postId"` //崗位編號
	PostName string `gorm:"size:128;" json:"postName"`              //崗位名稱
	PostCode string `gorm:"size:128;" json:"postCode"`              //崗位代碼
	Sort     int    `gorm:"size:4;" json:"sort"`                    //顯示排序
	Status   int    `gorm:"size:4;" json:"status"`                  //啟用/停用狀態
	Remark   string `gorm:"size:255;" json:"remark"`                //備註
	ControlBy
	ModelTime
}

func (SysPost) TableName() string {
	return "sys_post"
}