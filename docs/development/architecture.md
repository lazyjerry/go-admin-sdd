# Go-Admin 系統架構設計

本文件詳細說明 Go-Admin 後台管理系統的整體架構設計、核心元件、資料流向和技術選型，幫助開發者深入理解系統的設計理念和實作方式。

## 總體架構

Go-Admin 採用 **分層架構** + **模組化設計**，遵循 **領域驅動設計 (DDD)** 原則，實現前後端分離的現代化 Web 應用架構。

### 架構層次

```
┌─────────────────────────────────────────┐
│              前端層 (Frontend)            │
│  Vue.js + Element UI / Arco Design     │
└─────────────────────────────────────────┘
                    │ HTTP/HTTPS
                    ▼
┌─────────────────────────────────────────┐
│               閘道層 (Gateway)            │
│     Nginx + SSL + 負載均衡               │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│              應用層 (Application)         │
│        Gin Framework + 中間件            │
├─────────────────────────────────────────┤
│              業務層 (Business)            │
│    APIs + Services + Models             │
├─────────────────────────────────────────┤
│              數據層 (Data)               │
│        GORM + Database                  │
└─────────────────────────────────────────┘
                    │
                    ▼
┌─────────────────────────────────────────┐
│             基礎設施層 (Infrastructure)    │
│   MySQL/PostgreSQL + Redis + 檔案儲存    │
└─────────────────────────────────────────┘
```

### 核心設計原則

1. **單一職責原則** - 每個模組職責明確，互不干擾
2. **開放封閉原則** - 對擴展開放，對修改封閉
3. **介面隔離原則** - 依賴抽象介面，不依賴具體實作
4. **相依倒置原則** - 高階模組不依賴低階模組
5. **模組化設計** - 功能模組可獨立開發和部署

## 目錄結構詳解

### 專案根目錄

```
go-admin/
├── cmd/                    # 命令列介面層
├── app/                    # 應用程式核心層
├── common/                 # 公用基礎元件層
├── config/                 # 設定檔案層
├── docs/                   # 文件資源層
├── static/                 # 靜態資源層
├── template/               # 程式碼模板層
├── test/                   # 測試程式碼層
├── main.go                 # 應用程式入口
├── go.mod                  # Go 模組定義
├── go.sum                  # 相依性版本鎖定
├── Dockerfile              # Docker 容器定義
├── docker-compose.yml      # Docker 編排配置
└── Makefile               # 建構腳本
```

### 命令列介面層 (cmd/)

```
cmd/
├── cobra.go               # Cobra 命令列框架設定
├── api/                   # API 服務相關命令
│   ├── jobs.go           # 背景任務服務
│   ├── other.go          # 其他服務
│   └── server.go         # HTTP 伺服器
├── app/                   # 應用程式服務命令
│   └── server.go         # 應用伺服器啟動
├── config/                # 設定相關命令
│   └── server.go         # 設定服務
├── migrate/               # 資料庫遷移命令
│   ├── server.go         # 遷移服務
│   └── migration/        # 遷移腳本
└── version/               # 版本資訊命令
    └── server.go         # 版本查詢服務
```

**職責說明**：

- 提供命令列操作介面
- 負責應用程式啟動和初始化
- 實現資料庫遷移和種子資料
- 支援多種服務模式啟動

### 應用程式核心層 (app/)

```
app/
├── admin/                 # 管理員功能模組
│   ├── apis/             # API 控制器層
│   ├── models/           # 資料模型層
│   ├── router/           # 路由定義層
│   └── service/          # 業務邏輯層
├── jobs/                 # 背景任務模組
│   ├── examples.go       # 任務範例
│   ├── jobbase.go        # 任務基礎類別
│   ├── type.go           # 任務類型定義
│   ├── apis/             # 任務 API
│   ├── models/           # 任務模型
│   ├── router/           # 任務路由
│   └── service/          # 任務服務
└── other/                # 其他業務模組
    ├── apis/             # 其他 API
    ├── models/           # 其他模型
    ├── router/           # 其他路由
    └── service/          # 其他服務
```

**模組化設計**：

- 每個模組包含完整的 MVC 結構
- 模組間透過介面進行溝通
- 支援熱插拔和獨立部署

### 公用基礎元件層 (common/)

```
common/
├── actions/              # 通用 CRUD 操作
│   ├── create.go        # 建立操作
│   ├── delete.go        # 刪除操作
│   ├── index.go         # 查詢操作
│   ├── update.go        # 更新操作
│   ├── view.go          # 檢視操作
│   ├── permission.go    # 權限檢查
│   └── type.go          # 類型定義
├── apis/                # API 基礎類別
│   └── api.go           # API 基礎介面
├── database/            # 資料庫連接管理
│   ├── initialize.go    # 資料庫初始化
│   ├── open.go          # 連接開啟
│   └── open_sqlite3.go  # SQLite 連接
├── dto/                 # 資料傳輸物件
│   ├── auto_form.go     # 自動表單
│   ├── generate.go      # 程式碼生成
│   ├── order.go         # 排序處理
│   ├── pagination.go    # 分頁處理
│   ├── search.go        # 搜尋處理
│   └── type.go          # DTO 類型
├── middleware/          # HTTP 中間件
│   ├── auth.go          # 身份驗證
│   ├── cors.go          # 跨域處理
│   ├── logger.go        # 日誌記錄
│   └── rate_limit.go    # 流量限制
├── models/              # 基礎資料模型
├── response/            # HTTP 回應處理
├── service/             # 基礎服務介面
└── storage/             # 儲存服務
```

## 核心元件設計

### 1. HTTP 路由系統

#### 路由結構

```go
// 路由註冊架構
type Router interface {
    Register(engine *gin.Engine)
}

// 模組路由定義
func InitRouter() *gin.Engine {
    r := gin.New()

    // 全域中間件
    r.Use(middleware.Logger())
    r.Use(middleware.Recovery())
    r.Use(middleware.CORS())

    // API 路由群組
    api := r.Group("/api/v1")
    {
        // 公開路由 (無需認證)
        public := api.Group("/public")
        registerPublicRoutes(public)

        // 認證路由 (需要登入)
        auth := api.Group("/")
        auth.Use(middleware.Auth())
        registerAuthRoutes(auth)

        // 管理員路由 (需要管理員權限)
        admin := api.Group("/")
        admin.Use(middleware.Auth(), middleware.AuthCheckRole())
        registerAdminRoutes(admin)
    }

    return r
}
```

#### 中間件管道

```
HTTP Request
     │
     ▼
┌─────────────┐
│   Logger    │ ← 記錄請求日誌
└─────────────┘
     │
     ▼
┌─────────────┐
│   Recovery  │ ← 錯誤恢復處理
└─────────────┘
     │
     ▼
┌─────────────┐
│    CORS     │ ← 跨域資源共享
└─────────────┘
     │
     ▼
┌─────────────┐
│    Auth     │ ← JWT 身份驗證
└─────────────┘
     │
     ▼
┌─────────────┐
│ Permission  │ ← RBAC 權限檢查
└─────────────┘
     │
     ▼
┌─────────────┐
│  Handler    │ ← 業務邏輯處理
└─────────────┘
```

### 2. 權限控制系統 (RBAC)

#### 權限模型設計

```
使用者 (User) ──┐
              │ N:N
              ├───→ 角色 (Role) ──┐
              │                 │ N:N
部門 (Dept) ───┘                 ├───→ 權限 (Permission)
                                │
選單 (Menu) ────────────────────┘
```

#### Casbin 策略配置

```ini
# 角色繼承關係
g, user1, role_admin
g, user2, role_manager
g, user3, role_user

# 權限策略
p, role_admin, /api/v1/*, *
p, role_manager, /api/v1/user, GET|POST|PUT
p, role_user, /api/v1/user, GET

# 資料權限
g2, user1, dept_all
g2, user2, dept_tech
g2, user3, dept_tech
```

#### 權限檢查流程

```go
func AuthCheckRole() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 取得使用者資訊
        userID := jwt.GetUserIdFromToken(c)
        user := getUserByID(userID)

        // 2. 取得請求資源
        resource := c.Request.URL.Path
        method := c.Request.Method

        // 3. Casbin 權限檢查
        enforcer := casbin.GetEnforcer()
        allowed := enforcer.Enforce(user.Role, resource, method)

        if !allowed {
            response.Error(c, 403, "權限不足")
            c.Abort()
            return
        }

        // 4. 資料權限檢查
        checkDataScope(c, user)
        c.Next()
    }
}
```

### 3. 資料存取層 (DAL)

#### GORM 模型定義

```go
// 基礎模型
type BaseModel struct {
    ID        uint           `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time      `json:"createdAt"`
    UpdatedAt time.Time      `json:"updatedAt"`
    DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// 控制欄位 (審計追蹤)
type ControlBy struct {
    CreateBy int `json:"createBy" gorm:"index;comment:建立者"`
    UpdateBy int `json:"updateBy" gorm:"index;comment:更新者"`
}

// 使用者模型
type SysUser struct {
    BaseModel
    Username string `json:"username" gorm:"size:64;uniqueIndex;comment:使用者名稱"`
    Password string `json:"-" gorm:"size:128;comment:密碼"`
    NickName string `json:"nickName" gorm:"size:128;comment:暱稱"`
    Phone    string `json:"phone" gorm:"size:11;comment:手機號"`
    RoleId   int    `json:"roleId" gorm:"size:20;comment:角色ID"`
    DeptId   int    `json:"deptId" gorm:"size:20;comment:部門ID"`
    Status   int    `json:"status" gorm:"size:4;default:1;comment:狀態"`
    ControlBy

    // 關聯關係
    Role *SysRole `json:"role,omitempty"`
    Dept *SysDept `json:"dept,omitempty"`
}
```

#### Repository 模式

```go
// 儲存庫介面
type Repository interface {
    Create(entity interface{}) error
    Update(entity interface{}) error
    Delete(id uint) error
    FindByID(id uint, entity interface{}) error
    FindAll(entities interface{}) error
    FindWithCondition(condition interface{}, entities interface{}) error
}

// GORM 實作
type GormRepository struct {
    db *gorm.DB
}

func (r *GormRepository) Create(entity interface{}) error {
    return r.db.Create(entity).Error
}

func (r *GormRepository) FindByID(id uint, entity interface{}) error {
    return r.db.First(entity, id).Error
}
```

### 4. 業務服務層

#### 服務介面設計

```go
// 服務介面
type UserService interface {
    GetPage(req *dto.UserGetPageReq) (*dto.UserGetPageResp, error)
    Get(req *dto.UserGetReq) (*dto.UserGetResp, error)
    Create(req *dto.UserCreateReq) error
    Update(req *dto.UserUpdateReq) error
    Delete(req *dto.UserDeleteReq) error
}

// 服務實作
type userService struct {
    repo Repository
    log  logger.Logger
}

func (s *userService) Create(req *dto.UserCreateReq) error {
    // 1. 參數驗證
    if err := s.validateCreateReq(req); err != nil {
        return err
    }

    // 2. 業務邏輯處理
    user := &models.SysUser{
        Username: req.Username,
        NickName: req.NickName,
        Password: s.hashPassword(req.Password),
    }

    // 3. 資料持久化
    return s.repo.Create(user)
}
```

### 5. API 控制器層

#### RESTful API 設計

```go
// API 控制器
type UserApi struct {
    service UserService
}

// GET /api/v1/users - 分頁查詢使用者
func (api *UserApi) GetPage(c *gin.Context) {
    var req dto.UserGetPageReq
    if err := c.ShouldBindQuery(&req); err != nil {
        response.Error(c, 400, err.Error())
        return
    }

    result, err := api.service.GetPage(&req)
    if err != nil {
        response.Error(c, 500, err.Error())
        return
    }

    response.OK(c, result, "查詢成功")
}

// POST /api/v1/users - 建立使用者
func (api *UserApi) Create(c *gin.Context) {
    var req dto.UserCreateReq
    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, 400, err.Error())
        return
    }

    // 設定操作人員
    req.SetCreateBy(jwt.GetUserId(c))

    if err := api.service.Create(&req); err != nil {
        response.Error(c, 500, err.Error())
        return
    }

    response.OK(c, nil, "建立成功")
}
```

## 資料流架構

### 請求處理流程

```
Client Request
     │
     ▼
┌─────────────────┐
│   Gin Router    │ ← HTTP 路由分發
└─────────────────┘
     │
     ▼
┌─────────────────┐
│   Middleware    │ ← 認證、權限、日誌
└─────────────────┘
     │
     ▼
┌─────────────────┐
│  API Handler    │ ← 參數綁定、驗證
└─────────────────┘
     │
     ▼
┌─────────────────┐
│  Service Layer  │ ← 業務邏輯處理
└─────────────────┘
     │
     ▼
┌─────────────────┐
│ Repository Layer│ ← 資料存取
└─────────────────┘
     │
     ▼
┌─────────────────┐
│    Database     │ ← 資料持久化
└─────────────────┘
```

### 回應處理流程

```
Database Result
     │
     ▼
┌─────────────────┐
│   Model/Entity  │ ← ORM 物件對映
└─────────────────┘
     │
     ▼
┌─────────────────┐
│   DTO/Response  │ ← 資料轉換
└─────────────────┘
     │
     ▼
┌─────────────────┐
│   JSON Marshal  │ ← 序列化
└─────────────────┘
     │
     ▼
┌─────────────────┐
│  HTTP Response  │ ← 回應客戶端
└─────────────────┘
```

## 設定管理架構

### 設定檔案層次

```yaml
# 基礎設定
settings:
  # 應用程式設定
  application:
    mode: dev
    name: go-admin
    port: 8000

  # 資料庫設定
  database:
    dbtype: mysql
    host: 127.0.0.1
    port: 3306

  # 快取設定
  redis:
    host: 127.0.0.1
    port: 6379

  # 日誌設定
  log:
    level: info
    path: storage/logs
```

### 設定載入機制

```go
// 設定管理器
type Config struct {
    Application ApplicationConfig `yaml:"application"`
    Database    DatabaseConfig    `yaml:"database"`
    Redis       RedisConfig       `yaml:"redis"`
    Log         LogConfig         `yaml:"log"`
    JWT         JWTConfig         `yaml:"jwt"`
}

// 設定載入
func LoadConfig(path string) (*Config, error) {
    var config Config

    // 載入 YAML 檔案
    data, err := ioutil.ReadFile(path)
    if err != nil {
        return nil, err
    }

    // 解析設定
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    // 環境變數覆蓋
    overrideWithEnv(&config)

    // 設定驗證
    if err := validateConfig(&config); err != nil {
        return nil, err
    }

    return &config, nil
}
```

## 擴展性設計

### 1. 模組化擴展

#### 新增業務模組

```go
// 1. 建立模組目錄結構
app/
└── inventory/              # 庫存管理模組
    ├── apis/
    │   └── inventory.go
    ├── models/
    │   └── inventory.go
    ├── router/
    │   └── inventory.go
    └── service/
        └── inventory.go

// 2. 實作模組介面
type InventoryModule struct{}

func (m *InventoryModule) Name() string {
    return "inventory"
}

func (m *InventoryModule) Router() gin.HandlerFunc {
    return inventory.RegisterRouter
}

func (m *InventoryModule) Models() []interface{} {
    return []interface{}{
        &models.Inventory{},
    }
}

// 3. 註冊模組
func init() {
    module.Register(&InventoryModule{})
}
```

### 2. 中間件擴展

```go
// 自訂中間件
func CustomMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 前置處理
        start := time.Now()

        // 執行後續處理
        c.Next()

        // 後置處理
        latency := time.Since(start)
        log.Info("請求處理完成", "latency", latency)
    }
}

// 中間件註冊
func RegisterMiddlewares(r *gin.Engine) {
    r.Use(CustomMiddleware())
}
```

### 3. 儲存擴展

```go
// 儲存介面
type Storage interface {
    Upload(file io.Reader, filename string) (string, error)
    Download(filename string) (io.Reader, error)
    Delete(filename string) error
}

// 本地儲存實作
type LocalStorage struct {
    basePath string
}

// 雲端儲存實作
type CloudStorage struct {
    bucket string
    client *s3.Client
}

// 儲存工廠
func NewStorage(config StorageConfig) Storage {
    switch config.Type {
    case "local":
        return &LocalStorage{basePath: config.Path}
    case "s3":
        return &CloudStorage{bucket: config.Bucket}
    default:
        return &LocalStorage{}
    }
}
```

## 效能最佳化策略

### 1. 資料庫最佳化

#### 連接池設定

```go
// 資料庫連接池配置
func configureDB(db *gorm.DB) {
    sqlDB, _ := db.DB()

    // 設定最大空閒連接數
    sqlDB.SetMaxIdleConns(10)

    // 設定最大開啟連接數
    sqlDB.SetMaxOpenConns(100)

    // 設定連接最大生存時間
    sqlDB.SetConnMaxLifetime(time.Hour)

    // 設定連接最大空閒時間
    sqlDB.SetConnMaxIdleTime(time.Minute * 30)
}
```

#### 查詢最佳化

```go
// 使用索引和預載入
func (s *userService) GetUserWithRole(id uint) (*User, error) {
    var user User
    err := s.db.
        Preload("Role").          // 預載入角色
        Preload("Dept").          // 預載入部門
        Where("id = ?", id).      // 使用主鍵查詢
        First(&user).Error
    return &user, err
}

// 分頁查詢最佳化
func (s *userService) GetPageOptimized(req *dto.GetPageReq) (*dto.PageResp, error) {
    var users []User
    var total int64

    query := s.db.Model(&User{})

    // 計算總數 (不載入資料)
    query.Count(&total)

    // 分頁查詢 (只選擇需要的欄位)
    err := query.
        Select("id, username, nick_name, status").
        Offset(req.GetOffset()).
        Limit(req.GetLimit()).
        Find(&users).Error

    return &dto.PageResp{
        List:  users,
        Total: total,
    }, err
}
```

### 2. 快取策略

#### Redis 快取實作

```go
// 快取服務介面
type CacheService interface {
    Get(key string) (string, error)
    Set(key string, value string, expiration time.Duration) error
    Delete(key string) error
}

// Redis 實作
type redisCache struct {
    client *redis.Client
}

// 快取裝飾器
func WithCache(service UserService, cache CacheService) UserService {
    return &cachedUserService{
        service: service,
        cache:   cache,
    }
}

type cachedUserService struct {
    service UserService
    cache   CacheService
}

func (s *cachedUserService) Get(id uint) (*User, error) {
    // 嘗試從快取取得
    cacheKey := fmt.Sprintf("user:%d", id)
    if cached, err := s.cache.Get(cacheKey); err == nil {
        var user User
        json.Unmarshal([]byte(cached), &user)
        return &user, nil
    }

    // 從資料庫取得
    user, err := s.service.Get(id)
    if err != nil {
        return nil, err
    }

    // 寫入快取
    userJSON, _ := json.Marshal(user)
    s.cache.Set(cacheKey, string(userJSON), time.Hour)

    return user, nil
}
```

### 3. 併發處理

#### Goroutine 池

```go
// 工作者池
type WorkerPool struct {
    workers   int
    taskQueue chan Task
    wg        sync.WaitGroup
}

func NewWorkerPool(workers int) *WorkerPool {
    return &WorkerPool{
        workers:   workers,
        taskQueue: make(chan Task, workers*2),
    }
}

func (p *WorkerPool) Start() {
    for i := 0; i < p.workers; i++ {
        go p.worker()
    }
}

func (p *WorkerPool) worker() {
    for task := range p.taskQueue {
        task.Execute()
        p.wg.Done()
    }
}

func (p *WorkerPool) Submit(task Task) {
    p.wg.Add(1)
    p.taskQueue <- task
}
```

## 安全性架構

### 1. 身份驗證

```go
// JWT 權杖管理
type JWTManager struct {
    secretKey   string
    tokenExpiry time.Duration
}

func (j *JWTManager) GenerateToken(user *User) (string, error) {
    claims := &jwt.MapClaims{
        "user_id":  user.ID,
        "username": user.Username,
        "role":     user.Role,
        "exp":      time.Now().Add(j.tokenExpiry).Unix(),
        "iat":      time.Now().Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(j.secretKey))
}
```

### 2. 輸入驗證

```go
// 參數驗證
type CreateUserReq struct {
    Username string `json:"username" validate:"required,min=3,max=20"`
    Password string `json:"password" validate:"required,min=6,max=20"`
    Email    string `json:"email" validate:"required,email"`
    Phone    string `json:"phone" validate:"required,phone"`
}

// 驗證中間件
func ValidateRequest() gin.HandlerFunc {
    return func(c *gin.Context) {
        if err := validator.New().Struct(req); err != nil {
            response.Error(c, 400, "參數驗證失敗: "+err.Error())
            c.Abort()
            return
        }
        c.Next()
    }
}
```

### 3. SQL 注入防護

```go
// 使用 GORM 預處理語句
func (s *userService) FindByUsername(username string) (*User, error) {
    var user User
    // GORM 自動處理 SQL 注入防護
    err := s.db.Where("username = ?", username).First(&user).Error
    return &user, err
}

// 動態查詢建構
func (s *userService) Search(req *SearchReq) ([]User, error) {
    query := s.db.Model(&User{})

    if req.Username != "" {
        query = query.Where("username LIKE ?", "%"+req.Username+"%")
    }

    if req.Status != 0 {
        query = query.Where("status = ?", req.Status)
    }

    var users []User
    err := query.Find(&users).Error
    return users, err
}
```

## 監控和日誌架構

### 1. 結構化日誌

```go
// 日誌配置
type Logger struct {
    *zap.Logger
}

func NewLogger(config LogConfig) *Logger {
    cfg := zap.NewProductionConfig()
    cfg.OutputPaths = []string{config.Path}
    cfg.Level = zap.NewAtomicLevelAt(getLogLevel(config.Level))

    logger, _ := cfg.Build()
    return &Logger{logger}
}

// 結構化日誌記錄
func (l *Logger) LogRequest(c *gin.Context, start time.Time) {
    latency := time.Since(start)

    l.Info("HTTP請求",
        zap.String("method", c.Request.Method),
        zap.String("path", c.Request.URL.Path),
        zap.Int("status", c.Writer.Status()),
        zap.Duration("latency", latency),
        zap.String("ip", c.ClientIP()),
        zap.String("user_agent", c.Request.UserAgent()),
    )
}
```

### 2. 效能監控

```go
// 效能指標收集
type MetricsCollector struct {
    requestCount    prometheus.CounterVec
    requestDuration prometheus.HistogramVec
    activeUsers     prometheus.Gauge
}

func (m *MetricsCollector) RecordRequest(method, path string, duration time.Duration) {
    m.requestCount.WithLabelValues(method, path).Inc()
    m.requestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// 健康檢查端點
func HealthCheck(c *gin.Context) {
    health := map[string]interface{}{
        "status":    "ok",
        "timestamp": time.Now().Unix(),
        "version":   version.Version,
        "database":  checkDatabaseHealth(),
        "redis":     checkRedisHealth(),
    }

    c.JSON(200, health)
}
```

---

Go-Admin 的架構設計充分考慮了可擴展性、可維護性和效能需求，為企業級應用提供了堅實的技術基礎。

如需了解具體實作細節，請參考：

- [API 開發指南](./api-guide.md)
- [程式碼生成工具](./code-generation.md)
- [權限系統說明](./permission-system.md)
- [效能最佳化指南](../deployment/performance.md)
