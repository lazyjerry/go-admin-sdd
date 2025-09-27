// Package models 系統 API 模型定義
// 用於管理系統中所有 API 端點的資訊和權限控制
package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/bitly/go-simplejson"
	"github.com/go-admin-team/go-admin-core/sdk"
	"github.com/go-admin-team/go-admin-core/sdk/runtime"
	"github.com/go-admin-team/go-admin-core/storage"

	"go-admin/common/models"
)

// SysApi 系統 API 模型
// 用於記錄和管理系統中所有 RESTful API 端點資訊
//
// 主要功能：
// - API 權限控制的基礎資料
// - 支援動態 API 發現和註冊
// - 與角色權限系統整合
// - 提供 API 文件化支援
//
// 權限整合：
// - 透過 Casbin 實作細粒度權限控制
// - 支援角色與 API 的多對多關聯
// - 動態權限檢查和驗證
type SysApi struct {
	// 基本識別資訊
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement;comment:API唯一識別碼"`
	Handle string `json:"handle" gorm:"size:128;comment:API處理器名稱，用於程式識別"`
	Title  string `json:"title" gorm:"size:128;comment:API顯示標題，從Swagger註解獲取"`
	
	// 路由資訊
	Path   string `json:"path" gorm:"size:128;comment:API路徑，如 /api/v1/users"`
	Action string `json:"action" gorm:"size:16;comment:HTTP請求方法（GET/POST/PUT/DELETE等）"`
	Type   string `json:"type" gorm:"size:16;comment:API類型（SYS-系統API BUS-業務API）"`
	
	// 嵌入通用欄位
	models.ModelTime  // 建立時間、更新時間等時間戳欄位
	models.ControlBy  // 建立人、更新人等控制欄位
}

// TableName 指定資料表名稱
// 實作 GORM 的 Tabler 介面，自訂資料表名稱
func (*SysApi) TableName() string {
	return "sys_api"
}

// Generate 程式碼生成器支援方法
// 實作 ActiveRecord 介面，用於程式碼生成工具
// 返回當前 API 物件的副本，避免指標問題
func (e *SysApi) Generate() models.ActiveRecord {
	o := *e
	return &o
}

// GetId 取得主鍵值
// 實作 ActiveRecord 介面，返回 API 的主鍵 Id
// 用於通用的 CRUD 操作中識別 API 記錄
func (e *SysApi) GetId() interface{} {
	return e.Id
}

// SaveSysApi 自動發現並儲存系統 API 資訊
// 掃描應用程式中的所有路由，自動建立或更新 API 記錄到資料庫
//
// 功能特性：
// - 自動路由發現：掃描所有註冊的路由端點
// - Swagger 整合：從 Swagger 文件中提取 API 標題
// - 路徑正規化：處理動態路徑參數格式轉換
// - 智慧過濾：排除系統內建和靜態資源路由
//
// 參數：
//   - message: 包含路由資訊的訊息物件
//
// 返回值：
//   - error: 處理過程中的錯誤，nil 表示成功
//
// 使用場景：
// - 系統啟動時自動掃描 API
// - 開發環境動態更新 API 列表
// - 權限系統初始化
func SaveSysApi(message storage.Messager) (err error) {
	// 序列化訊息內容為 JSON
	var rb []byte
	rb, err = json.Marshal(message.GetValues())
	if err != nil {
		err = fmt.Errorf("JSON 序列化錯誤: %v", err.Error())
		return err
	}

	// 反序列化為路由列表結構
	var l runtime.Routers
	err = json.Unmarshal(rb, &l)
	if err != nil {
		err = fmt.Errorf("JSON 反序列化錯誤: %s", err.Error())
		return err
	}
	
	// 取得所有資料庫連接實例
	dbList := sdk.Runtime.GetDb()
	
	// 遍歷每個資料庫連接
	for _, d := range dbList {
		// 遍歷所有路由端點
		for _, v := range l.List {
			// 過濾系統內建和靜態資源路由
			// 排除 HEAD 請求、Swagger 文件、靜態檔案、表單生成器、系統表格等
			if v.HttpMethod != "HEAD" ||
				strings.Contains(v.RelativePath, "/swagger/") ||
				strings.Contains(v.RelativePath, "/static/") ||
				strings.Contains(v.RelativePath, "/form-generator/") ||
				strings.Contains(v.RelativePath, "/sys/tables") {

				// 從 Swagger 文件中提取 API 標題
				// 根據介面方法註解中的 @Summary 填充介面名稱
				// 主要用於程式碼生成器產生的 API
				jsonFile, _ := ioutil.ReadFile("docs/swagger.json")
				jsonData, _ := simplejson.NewFromReader(bytes.NewReader(jsonFile))
				
				// 處理動態路由參數格式轉換
				// 將 Gin 格式的 :id 轉換為 Swagger 格式的 {id}
				urlPath := v.RelativePath
				idPatten := "(.*)/:(\\w+)"  // 正則表達式：匹配 :參數 格式
				reg, _ := regexp.Compile(idPatten)
				if reg.MatchString(urlPath) {
					// 替換 :id 為 {id} 格式
					urlPath = reg.ReplaceAllString(v.RelativePath, "${1}/{${2}}")
				}
				
				// 從 Swagger JSON 中提取 API 標題
				apiTitle, _ := jsonData.Get("paths").Get(urlPath).Get(strings.ToLower(v.HttpMethod)).Get("summary").String()

				// 建立或更新 API 記錄
				// 使用 FirstOrCreate 避免重複記錄
				err := d.Debug().Where(SysApi{Path: v.RelativePath, Action: v.HttpMethod}).
					Attrs(SysApi{Handle: v.Handler, Title: apiTitle}).
					FirstOrCreate(&SysApi{}).
					Error
				if err != nil {
					err := fmt.Errorf("儲存系統 API 時發生錯誤: %s", err.Error())
					return err
				}
			}
		}
	}
	return nil
}
