// Package models API 回應格式定義
// 提供統一的 RESTful API 回應結構，確保前後端介面一致性
package models

// Response 統一的 API 回應結構
// 所有 API 端點都使用此結構回傳資料，確保回應格式的一致性
//
// 回應格式遵循以下原則：
// - 統一的狀態碼系統（HTTP 狀態碼 + 業務邏輯碼）
// - 標準化的訊息格式
// - 可追蹤的請求識別碼
// - 靈活的資料承載能力
//
// 使用場景：
// - RESTful API 回應封裝
// - 錯誤訊息標準化處理
// - 成功操作結果回傳
// - 分頁查詢結果封裝
// - 日誌追蹤和除錯支援
type Response struct {
	Code      int         `json:"code" example:"200"`       // HTTP 狀態碼，200表示成功，其他表示各種錯誤類型
	Data      interface{} `json:"data"`                     // 回應資料內容，可以是物件、陣列或基本類型
	Msg       string      `json:"msg"`                      // 回應訊息，提供操作結果的文字說明
	RequestId string      `json:"requestId"`                // 請求追蹤識別碼，用於日誌關聯和問題除錯
}

// Page 分頁查詢結果結構
// 用於包裝分頁查詢的結果資料，提供完整的分頁資訊
//
// 分頁特性：
// - 支援大資料量的分頁處理
// - 提供總記錄數統計
// - 前端分頁元件整合
// - 效能最佳化的查詢策略
//
// 使用場景：
// - 使用者列表分頁查詢
// - 訂單記錄分頁顯示
// - 日誌記錄分頁瀏覽
// - 任何需要分頁的資料查詢
type Page struct {
	List      interface{} `json:"list"`      // 當前頁面的資料列表，通常是物件陣列
	Count     int         `json:"count"`     // 符合查詢條件的總記錄數
	PageIndex int         `json:"pageIndex"` // 目前頁碼，從1開始計算
	PageSize  int         `json:"pageSize"`  // 每頁顯示的記錄數量
}

// ReturnOK 設定成功回應狀態
// 將回應標記為成功狀態（HTTP 200）
//
// 回傳值：
//   - *Response: 回應物件指標，支援鏈式呼叫
//
// 使用範例：
//   response := &Response{Data: userData, Msg: "查詢成功"}
//   return response.ReturnOK()
//
// 適用場景：
// - API 操作成功完成
// - 資料查詢成功回傳
// - 建立、更新、刪除操作成功
func (res *Response) ReturnOK() *Response {
	res.Code = 200
	return res
}

// ReturnError 設定錯誤回應狀態
// 根據提供的錯誤碼設定回應狀態
//
// 參數：
//   - code: HTTP 錯誤狀態碼（400、401、404、500等）
//
// 回傳值：
//   - *Response: 回應物件指標，支援鏈式呼叫
//
// 常用錯誤碼：
// - 400: 請求參數錯誤
// - 401: 未授權存取
// - 403: 權限不足
// - 404: 資源不存在
// - 422: 資料驗證失敗
// - 500: 伺服器內部錯誤
//
// 使用範例：
//   response := &Response{Msg: "使用者不存在"}
//   return response.ReturnError(404)
//
// 適用場景：
// - 資料驗證失敗
// - 權限檢查不通過
// - 資源查詢失敗
// - 系統異常處理
func (res *Response) ReturnError(code int) *Response {
	res.Code = code
	return res
}
