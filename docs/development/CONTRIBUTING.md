# 貢獻指南

歡迎參與 Go-Admin 專案！本文件將指導您如何有效地為專案做出貢獻。

## 貢獻方式

### 我們歡迎以下類型的貢獻

- 🐛 **錯誤回報**：發現並回報 bug
- 🚀 **功能建議**：提出新功能想法
- 📝 **文檔改進**：完善文檔和註釋
- 🔧 **程式碼修復**：修復已知問題
- ✨ **功能開發**：開發新功能
- 🧪 **測試增強**：增加測試覆蓋率
- 🎨 **UI/UX 改進**：改善使用者體驗

## 開始之前

### 1. 熟悉專案

在開始貢獻之前，請：

- 閱讀 [README.md](../README.md) 了解專案概況
- 瀏覽 [架構說明](./architecture.md) 理解系統設計
- 查看 [程式碼標準](./code-standards.md) 了解編碼規範
- 閱讀 [API 指南](./api-guide.md) 熟悉 API 設計

### 2. 環境準備

```bash
# 1. Fork 專案到你的 GitHub 帳戶

# 2. 複製你的 Fork
git clone https://github.com/YOUR_USERNAME/go-admin-sdd.git
cd go-admin-sdd

# 3. 添加上游倉庫
git remote add upstream https://github.com/go-admin-team/go-admin-sdd.git

# 4. 安裝依賴
go mod download

# 5. 安裝開發工具
make install-tools
```

### 3. 開發環境設定

```bash
# 複製配置檔案
cp config/settings.example.yml config/settings.yml

# 啟動開發環境
docker-compose -f docker-compose.dev.yml up -d

# 執行資料庫遷移
make migrate

# 啟動開發伺服器
make dev
```

## 貢獻流程

### 1. 建立 Issue

在開始工作之前，請先建立或找到相關的 Issue：

#### 錯誤回報範本

```markdown
**問題描述**
簡要描述遇到的問題

**重現步驟**

1. 前往 '...'
2. 點擊 '....'
3. 滾動到 '....'
4. 看到錯誤

**期望行為**
描述您期望會發生什麼

**實際行為**
描述實際發生了什麼

**環境資訊**

- OS: [例如 macOS, Windows, Linux]
- Go 版本: [例如 1.21]
- 瀏覽器: [例如 Chrome, Safari]
- 版本: [例如 v2.0.0]

**額外內容**
添加任何其他相關的截圖或內容
```

#### 功能請求範本

```markdown
**功能描述**
清楚簡潔地描述您想要的功能

**問題背景**
這個功能請求是為了解決什麼問題？

**建議解決方案**
描述您希望如何實現這個功能

**替代方案**
描述您考慮過的其他替代解決方案

**額外內容**
添加任何其他相關的截圖、流程圖或資料
```

### 2. 分支策略

```bash
# 確保本地 main 分支是最新的
git checkout main
git pull upstream main

# 建立功能分支
git checkout -b feature/your-feature-name
# 或者修復分支
git checkout -b fix/issue-number-description
```

#### 分支命名規範

- `feature/功能名稱` - 新功能開發
- `fix/問題描述` - 錯誤修復
- `docs/文檔類型` - 文檔更新
- `refactor/重構範圍` - 程式碼重構
- `test/測試範圍` - 測試相關

### 3. 程式碼開發

#### 開發規範

```go
// 1. 遵循程式碼標準
// 參見 docs/development/code-standards.md

// 2. 添加適當的註釋
// CreateUser 建立新的使用者帳戶
//
// 參數：
//   req: 包含使用者資訊的創建請求
//
// 返回值：
//   *User: 建立的使用者物件
//   error: 如果建立失敗則返回錯誤
func CreateUser(req *CreateUserRequest) (*User, error) {
    // 實作邏輯...
}

// 3. 添加測試
func TestCreateUser_Success(t *testing.T) {
    // 測試實作...
}
```

#### 提交規範

```bash
# 遵循 Conventional Commits 規範
git commit -m "feat(user): add user profile management

- Add user profile CRUD operations
- Implement profile image upload
- Add profile validation rules

Closes #123"
```

### 4. 測試要求

在提交 PR 之前，請確保：

```bash
# 執行所有測試
make test

# 檢查測試覆蓋率
make test-coverage

# 執行程式碼檢查
make lint

# 執行安全檢查
make security-check

# 格式化程式碼
make fmt
```

#### 測試覆蓋率要求

- 新增功能的測試覆蓋率應達到 **80%** 以上
- 修復的錯誤應包含重現測試案例
- 關鍵業務邏輯的測試覆蓋率應達到 **90%** 以上

### 5. 提交 Pull Request

#### PR 標題格式

```
<type>(<scope>): <description>

例如：
feat(user): add user profile management
fix(auth): resolve JWT token validation issue
docs(api): update authentication documentation
```

#### PR 描述範本

```markdown
## 變更摘要

<!-- 簡要描述這個 PR 的目的和主要變更 -->

## 變更類型

<!-- 在適當的選項前打勾 -->

- [ ] 🐛 錯誤修復 (非破壞性變更，修復問題)
- [ ] ✨ 新功能 (非破壞性變更，增加功能)
- [ ] 💥 破壞性變更 (修復或功能會導致現有功能無法正常運作)
- [ ] 📝 文檔更新 (改進或新增文檔)
- [ ] 🧪 測試 (新增測試或修正現有測試)
- [ ] 🔧 重構 (程式碼改進，不影響外部行為)

## 相關 Issue

<!-- 如果這個 PR 解決了某個 Issue，請使用關鍵字連結 -->

- Closes #123
- Fixes #456
- Resolves #789

## 變更詳情

### 新增

- 新增功能 A
- 新增功能 B

### 修改

- 改進現有功能 C
- 最佳化效能

### 移除

- 移除已廢棄的功能 D

## 測試

<!-- 描述測試的類型和涵蓋範圍 -->

- [ ] 單元測試已通過
- [ ] 整合測試已通過
- [ ] 手動測試已完成

### 測試場景

1. 測試場景 A
2. 測試場景 B
3. 測試場景 C

## 螢幕截圖/影片

<!-- 如果有 UI 變更，請提供截圖或 GIF -->

## 核對清單

<!-- 在完成項目前打勾 -->

- [ ] 我已閱讀並遵循貢獻指南
- [ ] 程式碼遵循專案的程式碼標準
- [ ] 我已經進行了自我審查
- [ ] 我已經為程式碼添加了適當的註釋
- [ ] 我已經進行了相應的文檔更新
- [ ] 我的變更不會生成新的警告
- [ ] 我已添加了證明修復有效或功能正常的測試
- [ ] 新的和現有的單元測試都通過了

## 額外資訊

<!-- 任何其他相關資訊 -->
```

## 程式碼審查

### 審查標準

我們的程式碼審查關注以下方面：

#### 1. 功能性

- ✅ 程式碼是否實現了預期功能
- ✅ 是否處理了邊界條件
- ✅ 錯誤處理是否適當

#### 2. 程式碼品質

- ✅ 程式碼是否清晰易讀
- ✅ 是否遵循命名規範
- ✅ 是否有適當的註釋

#### 3. 效能考量

- ✅ 是否存在效能瓶頸
- ✅ 資源使用是否合理
- ✅ 是否有記憶體洩漏

#### 4. 安全性

- ✅ 是否存在安全漏洞
- ✅ 輸入驗證是否充分
- ✅ 權限控制是否正確

### 回應審查意見

當收到審查意見時：

1. **認真考慮每個建議**
2. **提出問題或不同觀點**
3. **及時修改程式碼**
4. **更新測試和文檔**

```bash
# 修改後提交變更
git add .
git commit -m "fix: address code review comments"
git push origin feature/your-feature-name
```

## 社群互動

### 行為準則

我們致力於為每個人提供友善、安全和歡迎的環境。請遵循以下原則：

- 🤝 **尊重他人**：善待每個貢獻者
- 💭 **建設性溝通**：提供有用的反饋
- 🎯 **專注主題**：保持討論相關
- 🚫 **避免歧視**：不容忍任何形式的歧視

### 溝通管道

- **GitHub Issues**: 報告 bug 和功能請求
- **GitHub Discussions**: 一般討論和問題
- **Pull Request**: 程式碼相關討論
- **Discord/Slack**: 即時溝通（如有）

### 獲得協助

如果您需要協助：

1. 查看現有的 [Issues](https://github.com/go-admin-team/go-admin-sdd/issues)
2. 搜尋 [Discussions](https://github.com/go-admin-team/go-admin-sdd/discussions)
3. 閱讀相關文檔
4. 建立新的 Issue 或 Discussion

## 開發工具

### 推薦的開發環境

#### IDE/編輯器

- **VS Code**: 搭配 Go 擴展
- **GoLand**: JetBrains 的 Go IDE
- **Vim/Neovim**: 搭配 vim-go 外掛

#### 有用的 VS Code 擴展

```json
{
	"recommendations": ["golang.go", "ms-vscode.vscode-json", "redhat.vscode-yaml", "ms-python.python", "bradlc.vscode-tailwindcss"]
}
```

### 開發腳本

我們提供了便利的 Makefile 指令：

```makefile
# 開發相關
make dev          # 啟動開發伺服器
make test         # 執行測試
make lint         # 程式碼檢查
make fmt          # 程式碼格式化

# 資料庫相關
make migrate      # 執行資料庫遷移
make seed         # 填入測試資料

# 建置相關
make build        # 建置應用程式
make docker       # 建置 Docker 映像

# 清理
make clean        # 清理建置檔案
```

### Git Hooks

建議設定以下 Git Hooks：

```bash
# 設定 pre-commit hook
cp scripts/pre-commit.sh .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit

# 設定 pre-push hook
cp scripts/pre-push.sh .git/hooks/pre-push
chmod +x .git/hooks/pre-push
```

## 發佈流程

### 版本號規範

我們使用 [Semantic Versioning](https://semver.org/)：

- **主版本號 (MAJOR)**：不相容的 API 變更
- **次版本號 (MINOR)**：向下相容的新功能
- **修訂號 (PATCH)**：向下相容的問題修正

### 發佈檢核清單

發佈新版本時：

- [ ] 更新 CHANGELOG.md
- [ ] 更新版本號
- [ ] 建立 Git tag
- [ ] 建置和測試
- [ ] 發佈 GitHub Release
- [ ] 更新文檔

## 常見問題

### Q: 我是新手，能夠貢獻嗎？

A: 當然可以！我們歡迎所有層級的貢獻者。建議從以下開始：

- 修復文檔中的錯字
- 改進測試覆蓋率
- 處理標記為 `good first issue` 的 Issue

### Q: 如何決定要處理哪個 Issue？

A: 建議順序：

1. 標記為 `good first issue` 的問題
2. 您感興趣的功能領域
3. 您遇到的問題
4. 高優先級但無人處理的 Issue

### Q: 我的 PR 被拒絕了，怎麼辦？

A: 不要沮喪！這很正常：

1. 仔細閱讀審查意見
2. 詢問不清楚的地方
3. 根據建議修改程式碼
4. 重新提交審查

### Q: 如何跟上專案的最新動態？

A: 建議方式：

- Watch 這個倉庫獲得通知
- 定期閱讀 CHANGELOG
- 參與 Discussions
- 關注 Release notes

## 貢獻者認可

我們感謝每一位貢獻者的努力！

### 貢獻者列表

貢獻者會自動列入：

- [Contributors](https://github.com/go-admin-team/go-admin-sdd/graphs/contributors) 頁面
- README.md 的貢獻者區塊
- 發佈說明中的致謝

### 成為維護者

優秀的貢獻者可能會被邀請成為維護者：

**條件：**

- 持續的高質量貢獻
- 積極參與社群互動
- 熟悉專案架構和目標
- 展現領導能力

**責任：**

- 審查 Pull Requests
- 維護程式碼品質
- 協助新貢獻者
- 參與專案規劃

## 資源連結

### 專案文檔

- [專案架構](./architecture.md)
- [API 指南](./api-guide.md)
- [測試指南](./testing.md)
- [程式碼標準](./code-standards.md)

### 外部資源

- [Go 官方文檔](https://golang.org/doc/)
- [Gin 框架文檔](https://gin-gonic.com/docs/)
- [GORM 文檔](https://gorm.io/docs/)
- [Conventional Commits](https://www.conventionalcommits.org/)

## 聯絡我們

如果您有任何問題或建議：

- 📧 Email: admin@go-admin.dev
- 🐙 GitHub: [@go-admin-team](https://github.com/go-admin-team)
- 💬 Discussions: [專案討論區](https://github.com/go-admin-team/go-admin-sdd/discussions)

---

感謝您考慮為 Go-Admin 做出貢獻！您的每一個貢獻都讓這個專案變得更好。

_最後更新：2025 年 9 月 19 日_
