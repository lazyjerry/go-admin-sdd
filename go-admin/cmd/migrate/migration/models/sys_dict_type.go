package models

// DictType (字典類型) 定義了一組字典分類，例如性別、狀態等分組。
// 欄位：
// - DictId: 字典類型主鍵
// - DictName: 顯示名稱
// - DictType: 程式中引用的類型代碼
// - Status/Remark: 狀態與備註
type DictType struct {
	DictId   int    `gorm:"primaryKey;autoIncrement;" json:"dictId"`
	DictName string `gorm:"size:128;" json:"dictName"` //字典名稱
	DictType string `gorm:"size:128;" json:"dictType"` //字典類型
	Status   int    `gorm:"size:4;" json:"status"`     //狀態
	Remark   string `gorm:"size:255;" json:"remark"`   //備註
	ControlBy
	ModelTime
}

func (DictType) TableName() string {
	return "sys_dict_type"
}
