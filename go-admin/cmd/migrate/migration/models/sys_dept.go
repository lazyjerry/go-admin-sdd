package models

// SysDept 表示組織結構中的部門節點。
// 主要欄位：
// - DeptId: 部門編號（主鍵）
// - ParentId: 上級部門 ID
// - DeptPath: 部門路徑，用於快速查詢子節點
// - DeptName: 部門名稱
// - Sort: 排序權重
// - Leader/Phone/Email: 聯絡資訊
// 同時包含 ControlBy 與 ModelTime。
type SysDept struct {
	DeptId   int    `json:"deptId" gorm:"primaryKey;autoIncrement;"` //部門編號
	ParentId int    `json:"parentId" gorm:""`                        //上級部門 ID
	DeptPath string `json:"deptPath" gorm:"size:255;"`               //部門路徑
	DeptName string `json:"deptName"  gorm:"size:128;"`              //部門名稱
	Sort     int    `json:"sort" gorm:"size:4;"`                     //排序權重
	Leader   string `json:"leader" gorm:"size:128;"`                 //負責人
	Phone    string `json:"phone" gorm:"size:11;"`                   //手機
	Email    string `json:"email" gorm:"size:64;"`                   //郵箱
	Status   int    `json:"status" gorm:"size:4;"`                   //狀態
	ControlBy
	ModelTime
}

func (SysDept) TableName() string {
	return "sys_dept"
}
