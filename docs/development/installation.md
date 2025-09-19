# Go-Admin 安裝指南

本文件提供 Go-Admin 後台管理系統的詳細安裝說明，包含環境準備、依賴安裝、設定配置和常見問題解決方案。

## 環境需求

### 基本需求

- **Go 語言**: 1.21 或更高版本
- **Node.js**: 16.0 或更高版本 (前端開發需要)
- **Git**: 用於版本控制和專案下載

### 資料庫支援

選擇以下任一資料庫：

- **MySQL**: 5.7 或更高版本 (建議 8.0+)
- **PostgreSQL**: 12 或更高版本
- **SQLite**: 3.x (開發環境推薦)
- **SQL Server**: 2017 或更高版本

### 作業系統支援

- **Linux**: Ubuntu 18.04+、CentOS 7+、Debian 10+
- **macOS**: 10.14 或更高版本
- **Windows**: 10 或 Windows Server 2019+

## 環境準備

### 1. Go 語言環境安裝

#### Linux/macOS

```bash
# 下載 Go 安裝包 (以 1.21.5 為例)
wget https://golang.org/dl/go1.21.5.linux-amd64.tar.gz

# 解壓縮到 /usr/local
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# 設定環境變數
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export GO111MODULE=on' >> ~/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bashrc
source ~/.bashrc

# 驗證安裝
go version
```

#### Windows

```powershell
# 使用 Chocolatey 安裝 (需要先安裝 Chocolatey)
choco install golang

# 或直接下載安裝包
# https://golang.org/dl/go1.21.5.windows-amd64.msi

# 設定環境變數
setx GOPATH "%USERPROFILE%\go"
setx GO111MODULE "on"
setx GOPROXY "https://goproxy.cn,direct"

# 驗證安裝
go version
```

### 2. 資料庫安裝配置

#### MySQL 安裝

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install mysql-server mysql-client

# CentOS/RHEL
sudo yum install mysql-server mysql

# 啟動服務
sudo systemctl start mysql
sudo systemctl enable mysql

# 安全設定
sudo mysql_secure_installation
```

#### PostgreSQL 安裝

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install postgresql postgresql-contrib

# CentOS/RHEL
sudo yum install postgresql-server postgresql-contrib

# 初始化和啟動
sudo postgresql-setup initdb
sudo systemctl start postgresql
sudo systemctl enable postgresql
```

#### SQLite 安裝

```bash
# Ubuntu/Debian
sudo apt install sqlite3

# macOS (使用 Homebrew)
brew install sqlite

# Windows (使用 Chocolatey)
choco install sqlite
```

### 3. 開發工具安裝

```bash
# Git 安裝
# Ubuntu/Debian
sudo apt install git

# CentOS/RHEL
sudo yum install git

# macOS
brew install git

# Windows
choco install git
```

## 專案安裝

### 1. 下載原始碼

```bash
# 建立工作目錄
mkdir ~/go-admin-project
cd ~/go-admin-project

# 下載後端專案
git clone https://github.com/go-admin-team/go-admin.git
cd go-admin

# 檢查最新標籤版本
git tag -l
git checkout v2.1.0  # 替換為最新穩定版本
```

### 2. 安裝依賴套件

```bash
# 進入專案目錄
cd go-admin

# 更新依賴套件
go mod tidy

# 下載相關依賴
go mod download

# 驗證依賴完整性
go mod verify
```

### 3. 編譯專案

```bash
# 編譯執行檔
go build -o go-admin main.go

# 或使用 Makefile (如果存在)
make build

# 驗證編譯結果
./go-admin version
```

## 設定配置

### 1. 建立設定檔案

```bash
# 複製範例設定檔
cp config/settings.yml.example config/settings.yml

# 或複製開發環境設定
cp config/settings.dev.yml config/settings.yml
```

### 2. 資料庫設定

#### MySQL 配置

編輯 `config/settings.yml`：

```yaml
database:
  driver: mysql
  source: "admin:admin123@tcp(127.0.0.1:3306)/go_admin?charset=utf8mb4&parseTime=True&loc=Local"

# 或使用分離式設定
settings:
  database:
    dbtype: mysql
    host: 127.0.0.1
    port: 3306
    username: admin
    password: admin123
    dbname: go_admin
    config: charset=utf8mb4&parseTime=True&loc=Local
    max-idle-conns: 10
    max-open-conns: 100
    log-mode: info
    log-zap: true
```

#### PostgreSQL 配置

```yaml
database:
  driver: postgres
  source: "host=127.0.0.1 port=5432 user=admin password=admin123 dbname=go_admin sslmode=disable TimeZone=Asia/Taipei"
```

#### SQLite 配置

```yaml
database:
  driver: sqlite3
  source: "go-admin.db"
```

### 3. 應用程式設定

```yaml
# 基本設定
settings:
  application:
    # 應用程式模式: dev, test, prod
    mode: dev
    # 應用程式名稱
    name: go-admin
    # 服務埠號
    port: 8000
    # 服務主機
    host: 0.0.0.0
    # 讀取超時時間
    readtimeout: 60
    # 寫入超時時間
    writetimeout: 60

  # JWT 設定
  jwt:
    # JWT 密鑰 (生產環境請更換為強密碼)
    secret: go-admin-secret-key
    # 權杖過期時間 (小時)
    timeout: 24

  # 日誌設定
  log:
    # 日誌路徑
    path: storage/logs
    # 日誌檔名
    filename: go-admin.log
    # 日誌層級: debug, info, warn, error
    level: info
    # 最大保存天數
    max-age: 30
    # 單檔案最大大小 (MB)
    max-size: 100
    # 最大保存檔案數
    max-backups: 10
```

### 4. 環境變數設定

建立 `.env` 檔案：

```bash
# 應用程式設定
APP_MODE=dev
APP_PORT=8000
APP_SECRET=your-secret-key-here

# 資料庫設定
DB_TYPE=mysql
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USERNAME=admin
DB_PASSWORD=admin123
DB_NAME=go_admin

# Redis 設定 (可選)
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT 設定
JWT_SECRET=go-admin-jwt-secret
JWT_TIMEOUT=24h

# 檔案上傳設定
UPLOAD_PATH=storage/uploads
MAX_UPLOAD_SIZE=10485760  # 10MB

# 郵件設定 (可選)
MAIL_HOST=smtp.gmail.com
MAIL_PORT=587
MAIL_USERNAME=your-email@gmail.com
MAIL_PASSWORD=your-password
MAIL_FROM=your-email@gmail.com
```

## 資料庫初始化

### 1. 建立資料庫

#### MySQL

```sql
-- 連接 MySQL
mysql -u root -p

-- 建立資料庫
CREATE DATABASE go_admin CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 建立使用者
CREATE USER 'admin'@'localhost' IDENTIFIED BY 'admin123';

-- 授予權限
GRANT ALL PRIVILEGES ON go_admin.* TO 'admin'@'localhost';
FLUSH PRIVILEGES;

-- 退出
EXIT;
```

#### PostgreSQL

```bash
# 切換到 postgres 使用者
sudo -u postgres psql

# 建立資料庫
CREATE DATABASE go_admin;

# 建立使用者
CREATE USER admin WITH PASSWORD 'admin123';

# 授予權限
GRANT ALL PRIVILEGES ON DATABASE go_admin TO admin;

# 退出
\q
```

### 2. 執行資料庫遷移

```bash
# 執行遷移指令建立表格
./go-admin migrate -c config/settings.yml

# 或使用詳細輸出
./go-admin migrate -c config/settings.yml --verbose

# 檢查遷移狀態
./go-admin migrate -c config/settings.yml --status
```

### 3. 初始化基礎資料

```bash
# 執行種子資料初始化
./go-admin seed -c config/settings.yml

# 建立管理員帳號
./go-admin user create -c config/settings.yml \
  --username admin \
  --password admin123 \
  --email admin@example.com \
  --role admin
```

## 啟動服務

### 1. 開發模式啟動

```bash
# 使用設定檔啟動
./go-admin server -c config/settings.yml

# 或使用開發設定
./go-admin server -c config/settings.dev.yml

# 指定埠號啟動
./go-admin server -c config/settings.yml --port 8080

# 背景執行
nohup ./go-admin server -c config/settings.yml > go-admin.log 2>&1 &
```

### 2. 生產模式啟動

```bash
# 使用生產設定啟動
./go-admin server -c config/settings.prod.yml

# 使用 systemd 服務
sudo systemctl start go-admin
sudo systemctl enable go-admin
```

### 3. 驗證安裝

```bash
# 檢查服務狀態
curl http://localhost:8000/api/v1/health

# 測試登入 API
curl -X POST http://localhost:8000/api/v1/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'

# 檢查版本資訊
./go-admin version
```

## 前端安裝 (可選)

### 1. 下載前端專案

```bash
# 回到專案根目錄
cd ~/go-admin-project

# 下載前端專案
git clone https://github.com/go-admin-team/go-admin-ui.git
cd go-admin-ui
```

### 2. 安裝前端依賴

```bash
# 使用 npm
npm install

# 或使用 yarn
yarn install

# 或使用 cnpm (中國地區)
npm install -g cnpm
cnpm install
```

### 3. 啟動前端開發服務

```bash
# 開發模式
npm run dev

# 或
yarn dev

# 前端服務通常運行在 http://localhost:9527
```

## 常見問題與解決方案

### 1. Go 編譯問題

#### 問題：CGO 編譯錯誤

```bash
# 錯誤訊息
# github.com/mattn/go-sqlite3
cgo: exec /missing-cc: exec: "/missing-cc": file does not exist
```

**解決方案：**

```bash
# Windows 安裝 MinGW
choco install mingw

# 或安裝 TDM-GCC
# 下載：https://jmeubank.github.io/tdm-gcc/

# Linux 安裝 build-essential
sudo apt install build-essential

# macOS 安裝 Xcode Command Line Tools
xcode-select --install
```

#### 問題：Go 模組下載失敗

```bash
# 設定 Go 代理
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 或使用其他代理
go env -w GOPROXY=https://goproxy.io,direct
```

### 2. 資料庫連接問題

#### 問題：MySQL 連接被拒絕

```bash
# 檢查 MySQL 服務狀態
sudo systemctl status mysql

# 重新啟動 MySQL
sudo systemctl restart mysql

# 檢查埠號是否開放
netstat -tlnp | grep 3306

# 測試連接
mysql -h 127.0.0.1 -u admin -p go_admin
```

#### 問題：權限不足

```sql
-- 重新授予權限
GRANT ALL PRIVILEGES ON go_admin.* TO 'admin'@'localhost';
GRANT ALL PRIVILEGES ON go_admin.* TO 'admin'@'%';
FLUSH PRIVILEGES;
```

### 3. 設定檔案問題

#### 問題：YAML 格式錯誤

```bash
# 檢查 YAML 語法
python -c "import yaml; yaml.safe_load(open('config/settings.yml'))"

# 或使用線上工具驗證
# https://yaml-online-parser.appspot.com/
```

#### 問題：設定項目缺失

```bash
# 比較範例設定檔
diff config/settings.yml config/settings.yml.example

# 補充缺失的設定項目
```

### 4. 效能最佳化

#### 資料庫連接池設定

```yaml
database:
  max-idle-conns: 10 # 最大空閒連接數
  max-open-conns: 100 # 最大開啟連接數
  conn-max-lifetime: 3600 # 連接最大生存時間 (秒)
```

#### 記憶體使用最佳化

```bash
# 設定 Go 記憶體限制
export GOGC=100
export GOMEMLIMIT=512MB

# 或在程式中設定
```

## 升級指南

### 版本升級步驟

```bash
# 備份資料庫
mysqldump -u admin -p go_admin > backup_$(date +%Y%m%d_%H%M%S).sql

# 備份設定檔
cp config/settings.yml config/settings.yml.backup

# 拉取最新程式碼
git fetch origin
git checkout v2.2.0  # 新版本標籤

# 更新依賴
go mod tidy

# 重新編譯
go build -o go-admin main.go

# 執行資料庫遷移
./go-admin migrate -c config/settings.yml

# 重新啟動服務
sudo systemctl restart go-admin
```

### 設定檔案升級

```bash
# 檢查新版本設定變更
git diff v2.1.0..v2.2.0 config/settings.yml.example

# 手動合併設定變更
# 或使用工具輔助合併
```

## 開發環境配置

### IDE 設定

#### VS Code

安裝推薦擴展：

```json
{
	"recommendations": ["golang.go", "ms-vscode.vscode-json", "redhat.vscode-yaml", "bradlc.vscode-tailwindcss", "esbenp.prettier-vscode"]
}
```

#### GoLand

設定 Go 模組：

- 開啟 Go Modules 整合
- 設定 GOPROXY
- 配置程式碼格式化

### 除錯設定

```bash
# 使用 delve 除錯器
go install github.com/go-delve/delve/cmd/dlv@latest

# 除錯模式啟動
dlv debug main.go

# 或透過 IDE 整合除錯
```

---

安裝完成後，請參閱 [快速開始指南](./quickstart.md) 了解基本使用方法。

如有安裝問題，請參考 [常見問題解答](../faq.md) 或提交 [GitHub Issue](https://github.com/go-admin-team/go-admin/issues)。
