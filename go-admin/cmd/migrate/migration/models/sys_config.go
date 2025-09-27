package models

// SysConfig 表示系統配置項目，對應資料表 `sys_config`。
// 欄位說明：
// - ConfigName: 配置名稱
// - ConfigKey: 配置鍵
// - ConfigValue: 配置值
// - ConfigType: 配置類型（例如：string, number 等）
// - IsFrontend: 是否為前台可見配置（0/1）
// - Remark: 備註
// 同時繼承 ControlBy 與 ModelTime 用於記錄建立/更新者與時間。
type SysConfig struct {
	Model
	ConfigName  string `json:"configName" gorm:"type:varchar(128);comment:ConfigName"`
	ConfigKey   string `json:"configKey" gorm:"type:varchar(128);comment:ConfigKey"`
	ConfigValue string `json:"configValue" gorm:"type:varchar(255);comment:ConfigValue"`
	ConfigType  string `json:"configType" gorm:"type:varchar(64);comment:ConfigType"`
	IsFrontend  int    `json:"isFrontend" gorm:"type:varchar(64);comment:是否前台"`
	Remark      string `json:"remark" gorm:"type:varchar(128);comment:Remark"`
	ControlBy
	ModelTime
}

func (SysConfig) TableName() string {
	return "sys_config"
}
