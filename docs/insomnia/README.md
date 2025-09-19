# Go-Admin API 測試集合

本目錄包含 Go-Admin 後台管理系統的完整 API 測試集合，使用 Insomnia REST Client 進行 API 開發和測試。

## 📋 集合內容

### 核心功能測試

1. **系統健康檢查**

   - 系統狀態檢查
   - 服務可用性驗證

2. **身份驗證**

   - 使用者登入
   - 使用者登出
   - JWT Token 管理

3. **使用者管理**

   - 使用者列表查詢（支援分頁和搜尋）
   - 使用者詳情查詢
   - 使用者建立
   - 使用者資訊更新
   - 使用者刪除

4. **角色管理**

   - 角色列表查詢
   - 角色建立
   - 角色權限配置

5. **部門管理**

   - 部門樹狀結構查詢
   - 組織架構管理

6. **選單管理**

   - 系統選單查詢
   - 選單權限控制

7. **字典管理**
   - 系統字典資料查詢
   - 常用資料維護

### 錯誤處理測試

- **404 錯誤**：查詢不存在的資源
- **401 錯誤**：未授權存取
- **422 錯誤**：輸入驗證錯誤
- **500 錯誤**：伺服器內部錯誤

## 🚀 快速開始

### 1. 匯入集合

1. 開啟 Insomnia REST Client
2. 點選「Import/Export」→「Import Data」
3. 選擇「From File」
4. 匯入 `go-admin-api-collection.yaml` 檔案

### 2. 環境配置

匯入後會自動建立三個環境：

- **Local Development**：`http://localhost:8000`
- **Docker Development**：`http://localhost:8080`
- **Production**：`https://api.go-admin.dev`

### 3. 使用流程

#### 基本測試流程

```
1. 選擇適當的環境 (Local/Docker/Production)
2. 執行「系統健康檢查」確認服務可用
3. 執行「使用者登入」取得 JWT Token
4. 更新環境變數中的 jwt_token
5. 執行其他 API 測試
```

#### 完整 CRUD 測試流程

```
使用者管理測試流程：
1. 登入系統 → 取得 Token
2. 獲取使用者列表 → 查看現有資料
3. 建立新使用者 → 記錄回傳的 user_id
4. 更新環境變數 user_id
5. 獲取指定使用者詳情 → 驗證建立結果
6. 更新使用者資訊 → 測試更新功能
7. 刪除指定使用者 → 測試刪除功能
```

## 🔧 環境變數說明

### 必要變數

| 變數名稱    | 說明                   | 範例值                    |
| ----------- | ---------------------- | ------------------------- |
| `base_url`  | API 基礎網址           | `http://localhost:8000`   |
| `jwt_token` | 登入後取得的 JWT Token | `eyJhbGciOiJIUzI1NiIs...` |
| `user_id`   | 測試用的使用者 ID      | `1`                       |

### 設定步驟

1. **自動取得 JWT Token**

   - 執行「使用者登入」請求
   - 從回應中複製 `token` 值
   - 手動更新環境變數 `jwt_token`

2. **動態使用者 ID**
   - 建立新使用者後從回應取得 `userId`
   - 更新環境變數 `user_id` 用於後續測試

## 📝 測試資料

### 預設登入帳號

```json
{
	"username": "admin",
	"password": "123456",
	"code": "",
	"uuid": ""
}
```

### 測試使用者資料

```json
{
	"username": "testuser",
	"nickName": "測試使用者",
	"phone": "0912345678",
	"email": "test@example.com",
	"sex": "1",
	"deptId": 1,
	"postId": 1,
	"roleId": 2,
	"status": "2",
	"password": "password123"
}
```

### 測試角色資料

```json
{
	"roleName": "測試角色",
	"roleKey": "test_role",
	"roleSort": 1,
	"remark": "測試角色說明",
	"admin": false,
	"status": "2"
}
```

## 🔍 API 回應格式

### 成功回應

```json
{
	"code": 200,
	"message": "操作成功",
	"data": {
		// 實際資料內容
	},
	"timestamp": "2025-09-19T10:30:00Z",
	"requestId": "req_12345"
}
```

### 分頁回應

```json
{
  "code": 200,
  "message": "查詢成功",
  "data": {
    "list": [...],
    "count": 100,
    "pageIndex": 1,
    "pageSize": 10
  }
}
```

### 錯誤回應

```json
{
	"code": 400,
	"message": "請求參數錯誤",
	"data": null,
	"timestamp": "2025-09-19T10:30:00Z",
	"requestId": "req_12345"
}
```

## 🛠️ 故障排除

### 常見問題

#### 1. 401 未授權錯誤

**問題**：API 回應 401 Unauthorized

**解決方法**：

- 確認已執行登入請求
- 檢查 `jwt_token` 環境變數是否正確設定
- 確認 Token 未過期

#### 2. 404 資源不存在

**問題**：API 回應 404 Not Found

**解決方法**：

- 檢查 API 端點網址是否正確
- 確認資源 ID 是否存在
- 檢查 `base_url` 環境變數設定

#### 3. 連接錯誤

**問題**：無法連接到伺服器

**解決方法**：

- 確認 Go-Admin 服務已啟動
- 檢查 `base_url` 設定
- 確認網路連接正常
- 檢查防火牆設定

#### 4. 驗證錯誤

**問題**：422 Unprocessable Entity

**解決方法**：

- 檢查請求資料格式
- 確認必要欄位已提供
- 驗證資料類型和格式

### 除錯技巧

1. **檢查請求標頭**

   - 確認 `Content-Type: application/json`
   - 確認 `Authorization: Bearer {token}` 格式正確

2. **查看完整回應**

   - 檢查 HTTP 狀態碼
   - 查看錯誤訊息詳情
   - 檢查回應標頭資訊

3. **使用系統健康檢查**
   - 優先執行健康檢查確認服務狀態
   - 確認基礎服務功能正常

## 📚 進階使用

### 批次測試

1. 使用 Insomnia 的「Collection Runner」功能
2. 設定測試順序和資料流
3. 自動化回歸測試

### 自動化測試腳本

```javascript
// 在 Insomnia 中使用 Pre-request Script
const baseUrl = insomnia.environment.base_url;
const token = insomnia.environment.jwt_token;

// 在 After Response Script 中處理回應
const response = insomnia.response.json();
if (response.data && response.data.token) {
	insomnia.environment.jwt_token = response.data.token;
}
```

### 測試資料管理

- 使用環境變數儲存測試資料
- 建立測試資料重設流程
- 實作測試資料清理機制

## 📋 檢查清單

### 測試前檢查

- [ ] Go-Admin 服務已啟動
- [ ] 資料庫連接正常
- [ ] 環境變數已正確設定
- [ ] 測試帳號可正常登入

### 測試執行檢查

- [ ] 系統健康檢查通過
- [ ] 使用者登入成功並取得 Token
- [ ] 所有 CRUD 操作測試通過
- [ ] 錯誤處理測試符合預期
- [ ] 權限控制測試正常

### 測試後檢查

- [ ] 測試資料已清理
- [ ] 無殘留測試帳號
- [ ] 系統狀態恢復正常
- [ ] 測試結果已記錄

## 🤝 貢獻指南

歡迎提交改進建議和錯誤回報：

1. **新增測試案例**：擴展 API 覆蓋範圍
2. **改善測試資料**：提供更真實的測試場景
3. **優化文檔**：改善使用說明和故障排除
4. **回報問題**：提交發現的 API 問題和建議

## 📄 授權條款

本測試集合遵循與 Go-Admin 專案相同的 MIT License 授權條款。

---

**最後更新：2025 年 9 月 19 日**
