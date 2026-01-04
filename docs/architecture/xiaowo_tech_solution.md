# 小窝项目技术方案与架构设计

## 执行摘要

基于对Syncplay和SyncTV开源项目的深入技术调研，结合小窝项目的PRD需求和MVP快速上线的要求，本文档设计了一套完整的技术方案。该方案采用Go语言后端配合现代前端框架的架构，通过WebSocket实现实时同步，SQLite提供数据存储，Docker容器化部署，确保MVP阶段快速开发和后续稳定运营的双重目标。

**核心设计原则**：
- **MVP优先**：快速迭代，最小可行产品先行
- **技术成熟**：采用经过验证的技术栈
- **可扩展性**：预留扩展空间，支持业务增长
- **开发效率**：平衡开发速度与代码质量

---

## 一、架构与选型

### 1.1 技术栈选型

**最终技术栈**：
```
后端: Go 1.21+ + Gin Framework
前端: Vue 3 + TypeScript + Vite
实时通信: WebSocket
消息协议: JSON (MVP) -> Protocol Buffers (后期)
数据库: SQLite (Embedded)
ORM: GORM
缓存: Go In-Memory (MVP)
部署: Docker + Docker Compose
日志: 结构化日志 (zap) + Docker Logs + 性能监控
错误处理: 统一错误码体系 + 中间件处理
```

**选型理由**：
- **Go语言**：高性能、并发能力强、生态成熟
- **Vue 3**：开发效率高、组件化、TypeScript支持
- **WebSocket**：实时双向通信、浏览器原生支持
- **SQLite**：单文件数据库，零配置，无需独立容器，完美契合单体架构与MVP需求。使用纯Go的SQLite驱动（github.com/glebarez/sqlite），避免CGO依赖，确保在Docker环境中稳定运行
- **Docker**：环境一致性、快速部署

### 1.2 系统整体架构

```
┌─────────────────────────────────────────────────────────────────┐
│                        小窝项目架构图                             │
├─────────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐               │
│  │   Web端     │  │   移动端    │  │   管理后台  │               │
│  │  Vue 3 App  │  │   (预留)    │  │  Admin UI   │               │
│  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘               │
│         │                │                │                      │
│         ▼                ▼                ▼                      │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                   Nginx 反向代理                            │ │
│  │            (SSL 终结 + 静态资源服务)                         │ │
│  └─────────────────┬───────────────────────────────────────────┘ │
│                    │                                           │
│  ┌─────────────────▼───────────────────────────────────────────┐ │
│  │                    小窝后端服务 (Go)                         │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐             │ │
│  │  │  HTTP API   │ │  WebSocket  │ │  业务逻辑   │             │ │
│  │  │   服务      │ │   Hub       │ │    层       │             │ │
│  │  └─────────────┘ └─────────────┘ └─────────────┘             │ │
│  │  ┌─────────────┐ ┌─────────────┐ ┌─────────────┐             │ │
│  │  │   房间管理   │ │  同步逻辑   │ │   消息处理  │             │ │
│  │  │    服务      │ │            │ │             │             │ │
│  │  └─────────────┘ └─────────────┘ └─────────────┘             │ │
│  └─────────────────┬───────────────────────────────────────────┘ │
│                    │                                           │
│  ┌─────────────────▼───────────────────────────────────────────┐ │
│  │                SQLite 数据库                                │ │
│  │     (rooms, room_members, user_sessions)                   │ │
│  └─────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────┘
```

**架构特点**：
- **单体架构**：所有服务在一个Go进程中，简化部署和运维
- **无状态设计**：应用层完全无状态，便于横向扩展
- **WebSocket优先**：核心功能通过WebSocket实时同步
- **SQLite单文件**：数据持久化，零运维成本

### 1.3 核心模块设计

**1. HTTP API服务**
- 提供房间管理、用户会话等RESTful接口
- 基于Gin框架，高性能路由
- JSON格式数据交换
- 统一错误处理中间件
- 性能监控中间件

**2. WebSocket Hub**
- 管理所有WebSocket连接
- 处理实时消息转发
- 房间广播机制
- 连接心跳检测
- 消息频率限制
- 连接数限制

**3. 业务逻辑层**
- 房间生命周期管理
- 用户会话管理
- 同步逻辑处理
- 业务规则验证

**4. 数据访问层**
- Repository模式封装
- GORM ORM映射
- SQLite事务管理
- 连接池优化

**5. 错误处理机制**
- 统一错误码体系（系统级、业务级、参数级）
- 结构化错误响应格式
- 错误日志记录
- 跨服务错误传播

**错误处理代码实现**：
```go
// 统一错误码定义
const (
    // 系统级错误 (1000-1999)
    ErrCodeInternalServer = 1001
    ErrCodeDatabaseError  = 1002
    ErrCodeNetworkError   = 1003
    
    // 业务级错误 (2000-2999)
    ErrCodeRoomNotFound   = 2001
    ErrCodeRoomFull       = 2002
    ErrCodeInvalidToken   = 2003
    ErrCodePermissionDenied = 2004
    
    // 参数级错误 (3000-3999)
    ErrCodeInvalidParam   = 3001
    ErrCodeMissingParam   = 3002
    ErrCodeParamTooLong   = 3003
)

// 统一错误响应格式
type ErrorResponse struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
    RequestID string `json:"request_id"`
}

// 错误处理中间件
func ErrorHandler() gin.HandlerFunc {
    return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
        log.Printf("[PANIC] %s\nStack: %s", recovered, string(debug.Stack()))
        
        requestID := c.GetHeader("X-Request-ID")
        if requestID == "" {
            requestID = generateRequestID()
        }
        
        c.JSON(http.StatusInternalServerError, ErrorResponse{
            Code:    ErrCodeInternalServer,
            Message: "内部服务器错误",
            RequestID: requestID,
        })
        c.Abort()
    })
}
```

**6. 日志系统**
- 结构化日志（zap框架）
- 分级日志（DEBUG、INFO、WARN、ERROR）
- 业务日志与技术日志分离
- 性能指标记录
- 分布式追踪支持（预留）

**日志系统代码实现**：
```go
// 结构化日志配置
type LogConfig struct {
    Level     string `json:"level"`
    Format    string `json:"format"`    // json/text
    Output    string `json:"output"`    // stdout/file
    Filename  string `json:"filename"`  // 文件输出时指定
    MaxSize   int    `json:"max_size"`  // MB
    MaxAge    int    `json:"max_age"`   // 天
    MaxBackup int    `json:"max_backup"`
}

// 业务日志记录器
type BusinessLogger struct {
    logger *zap.Logger
}

func (l *BusinessLogger) RoomCreated(roomID, sessionID string) {
    l.logger.Info("room_created",
        zap.String("room_id", roomID),
        zap.String("session_id", sessionID),
        zap.String("event", "business"),
    )
}

func (l *BusinessLogger) SyncDelay(roomID string, delayMs int64, userCount int) {
    l.logger.Info("sync_performance",
        zap.String("room_id", roomID),
        zap.Int64("delay_ms", delayMs),
        zap.Int("user_count", userCount),
        zap.String("event", "performance"),
    )
}

// 性能监控中间件
func PerformanceMonitor() gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        requestID := generateRequestID()
        
        c.Header("X-Request-ID", requestID)
        
        c.Next()
        
        duration := time.Since(start)
        
        // 记录性能指标
        if duration > 100*time.Millisecond {
            logger.Warn("slow_request",
                zap.String("path", c.Request.URL.Path),
                zap.Duration("duration", duration),
                zap.Int("status", c.Writer.Status()),
                zap.String("request_id", requestID),
            )
        }
    })
}
```

### 1.4 项目目录结构

```
xiaowo/
├── backend/                    # Go 后端
│   ├── cmd/
│   │   └── server/
│   │       └── main.go        # 应用入口
│   ├── internal/
│   │   ├── api/v1/            # API 接口层
│   │   │   ├── room.go        # 房间 API
│   │   │   ├── types.go       # 数据类型定义
│   │   │   ├── middleware.go  # 中间件
│   │   │   └── router.go      # 路由配置
│   │   ├── service/           # 业务逻辑层
│   │   │   ├── room_service.go
│   │   │   └── member_service.go
│   │   ├── repository/        # 数据访问层
│   │   │   ├── room_repo.go
│   │   │   └── member_repo.go
│   │   ├── model/             # 数据模型
│   │   │   ├── room.go
│   │   │   └── member.go
│   │   ├── websocket/         # WebSocket 处理
│   │   │   └── hub.go         # WebSocket Hub
│   │   └── database/          # 数据库配置
│   │       └── database.go    # SQLite 配置（WAL模式默认开启）
│   ├── pkg/
│   │   └── utils/             # 公共工具
│   ├── go.mod
│   └── Dockerfile
├── frontend/                   # Vue 前端 (预留)
│   ├── src/
│   ├── package.json
│   └── Dockerfile
├── nginx/                      # Nginx 配置
│   └── prod.conf
├── docker-compose.prod.yml     # 生产环境编排
├── docker-compose.yml          # 开发环境编排
├── README.md
└── docs/
    └── API.md                 # API 文档
```

#### 网络重连与状态恢复策略

**1. 重连后自动状态拉取机制**
- **触发时机**：WebSocket重连成功后立即执行
- **同步请求**：主动发送 `MSG_SYNC` 请求到服务器
- **状态响应**：服务器返回完整房间状态（播放/暂停状态、当前进度、播放倍速、缓冲状态）
- **技术实现**：
  ```javascript
  // WebSocket重连成功后的状态恢复
  function handleReconnectSuccess() {
    // 1. 发送同步请求
    socket.send(JSON.stringify({
      type: MSG_SYNC,
      room_id: currentRoomId,
      session_id: currentSessionId,
      timestamp: Date.now()
    }));
  }
  ```

**2. 平滑过渡策略 (基于四级同步策略)**
- **微小差异 (< 2秒)**：优先使用"温柔微调"策略（0.05s-3s范围）追帧
  - 计算公式：`追帧系数 = 1.0 + (time_diff / 10.0)`，限制在0.95x-1.05x范围
  - 特点：无感追平，避免画面突变
- **较大差异 (≥ 2秒)**：执行Seek操作，强制对齐
  - 显示"正在同步..."提示
  - 直接跳转到目标时间点
- **差异判定逻辑**：
  ```javascript
  function determineSyncStrategy(localTime, serverTime) {
    const timeDiff = Math.abs(localTime - serverTime);
    
    if (timeDiff < 2.0) {
      return 'gentle_adjust';  // 温柔微调
    } else {
      return 'seek';          // 强制跳转
    }
  }
  ```

**3. 状态恢复完整流程**
1. **检测重连成功** → WebSocket `onopen` 事件触发
2. **发送同步请求** → 主动请求最新房间状态
3. **接收状态响应** → 解析服务器返回的完整状态
4. **计算时间差异** → 对比本地与服务器时间
5. **选择恢复策略** → 基于差异大小选择微调或跳转
6. **执行状态应用** → 更新本地播放器状态
7. **恢复同步监听** → 重新开始正常同步流程

**4. 异常处理与降级**
- **服务器无响应**：等待3秒后重试，最多重试3次
- **状态不一致**：以服务器状态为准，强制覆盖本地状态
- **用户干预**：恢复过程中用户手动操作则中断自动恢复
- **网络波动**：短时间内多次重连时，合并恢复请求

---

## 二、数据模型与数据库架构

### 2.1 数据模型设计

**用户模型 (UserSession) [临时会话]**
```go
// 仅在内存或SQLite中维护，无需独立缓存服务
type UserSession struct {
    ID        string    `json:"id"`        // UUID
    Nickname  string    `json:"nickname"`  // 随机生成的趣味昵称 (e.g. "快乐的考拉")
    Avatar    string    `json:"avatar"`    // 随机头像URL
    RoomID    string    `json:"room_id"`   // 当前所在房间
    CreatedAt time.Time `json:"created_at"`
}
```

**房间模型 (Room)**
```go
type Room struct {
    ID               string    `json:"id"`                 // 6位数字+字母混合房间号
    Name             string    `json:"name"`               // 房间名称
    Description      string    `json:"description"`        // 房间描述
    CreatorSessionID string    `json:"creator_session_id"` // 创建者会话ID
    IsPrivate        bool      `json:"is_private"`         // 是否私密房间
    Password         string    `json:"password"`           // 房间密码 (如有)
    MaxUsers         int       `json:"max_users"`          // 最大用户数(固定为7)
    Status           string    `json:"status"`             // 房间状态: active/inactive
    MediaURL         string    `json:"media_url"`          // 媒体资源URL
    PlaybackState    string    `json:"playback_state"`     // 播放状态: playing/paused
    CurrentTime      float64   `json:"current_time"`       // 当前播放时间
    Duration         float64   `json:"duration"`           // 媒体总时长
    Settings         string    `json:"settings"`           // JSON格式的设置
    Version          int       `json:"version"`            // 版本号(乐观锁)
    LastActiveAt     time.Time `json:"last_active_at"`     // 最后活跃时间
    LastMemberLeftAt time.Time `json:"last_member_left_at"` // 最后一位成员离开时间
    CreatedAt        time.Time `json:"created_at"`
    UpdatedAt        time.Time `json:"updated_at"`
}
```

**房间号生成规则**:
- 格式: 6位字符，包含数字(0-9)和大写字母(A-Z)
- 唯一性: 确保生成的房间号唯一
- 生成算法: 使用随机数生成，配合查重机制
- 示例: 8XA92B, K7L3P2

**房间成员模型 (RoomMember)**
```go
type RoomMember struct {
    ID         string    `json:"id"`          // UUID
    RoomID     string    `json:"room_id"`     // 房间ID
    SessionID  string    `json:"session_id"`  // 用户会话ID
    Nickname   string    `json:"nickname"`    // 用户昵称
    Avatar     string    `json:"avatar"`      // 用户头像
    Role       string    `json:"role"`        // 角色: host/guest
    JoinedAt   time.Time `json:"joined_at"`   // 加入时间
    LastSeenAt time.Time `json:"last_seen_at"` // 最后在线时间
    CreatedAt  time.Time `json:"created_at"`
}
```

### 2.2 数据库表结构设计

**核心表结构 (SQLite Compatible DDL)**

```sql
-- 用户会话表 (临时会话)
CREATE TABLE user_sessions (
    id TEXT PRIMARY KEY,
    nickname TEXT NOT NULL,
    avatar TEXT NOT NULL,
    room_id TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 房间表
CREATE TABLE rooms (
    id TEXT PRIMARY KEY, -- 6位数字+字母混合房间号
    name TEXT NOT NULL,
    description TEXT,
    creator_session_id TEXT NOT NULL,
    is_private BOOLEAN DEFAULT 0,
    password TEXT,
    max_users INTEGER DEFAULT 7, -- 最大用户数固定为7
    status TEXT DEFAULT 'active',
    media_url TEXT NOT NULL,
    playback_state TEXT DEFAULT 'paused',
    current_time REAL DEFAULT 0,
    duration REAL DEFAULT 0,
    settings TEXT DEFAULT '{}',
    version INTEGER DEFAULT 0,
    last_active_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_member_left_at DATETIME, -- 最后一位成员离开时间
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- 房间自动销毁机制
-- 1. 当房间成员数为0时，更新last_member_left_at
-- 2. 定时任务每5分钟检查一次，若last_member_left_at超过10分钟，则将status改为inactive
-- 3. 另一个定时任务每小时清理inactive状态超过1小时的房间

-- 房间成员表
CREATE TABLE room_members (
    id TEXT PRIMARY KEY,
    room_id TEXT NOT NULL,
    session_id TEXT NOT NULL,
    nickname TEXT NOT NULL,
    avatar TEXT NOT NULL,
    role TEXT DEFAULT 'guest',
    joined_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_seen_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (room_id) REFERENCES rooms(id) ON DELETE CASCADE,
    FOREIGN KEY (session_id) REFERENCES user_sessions(id) ON DELETE CASCADE
);

-- 索引优化
CREATE INDEX idx_rooms_creator ON rooms(creator_session_id);
CREATE INDEX idx_rooms_status ON rooms(status);
CREATE INDEX idx_rooms_last_active ON rooms(last_active_at);
CREATE INDEX idx_room_members_room ON room_members(room_id);
CREATE INDEX idx_room_members_session ON room_members(session_id);
CREATE INDEX idx_user_sessions_room ON user_sessions(room_id);

-- 触发器: 自动更新 updated_at
CREATE TRIGGER rooms_update_updated_at 
    AFTER UPDATE ON rooms
    BEGIN
        UPDATE rooms SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
    END;
```

### 2.3 数据访问层设计

**数据库连接配置**：
```go
// database/database.go
package database

import (
    "os"
    "log"
    "time"
    "fmt"
    "github.com/glebarez/sqlite" // 纯Go驱动，无需CGO
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// InitDB 初始化SQLite数据库连接（优化版）
func InitDB() (*gorm.DB, error) {
    dbSource := os.Getenv("DB_SOURCE")
    if dbSource == "" {
        dbSource = "xiaowo.db"
    }
    
    // SQLite DSN配置优化
    dsn := fmt.Sprintf("%s?"+
        "cache=shared&"+
        "mode=rwc&"+
        "_journal_mode=WAL&"+
        "_busy_timeout=5000&"+
        "_synchronous=NORMAL&"+
        "_cache_size=10000&"+
        "_temp_store=memory&"+
        "_mmap_size=268435456", // 256MB mmap
        dbSource)
    
    // 配置日志级别 - 生产环境仅打印错误
    logLevel := logger.Info
    if os.Getenv("GIN_MODE") == "release" {
        logLevel = logger.Error
    }
    
    newLogger := logger.New(
        log.New(os.Stdout, "\r\n", log.LstdFlags),
        logger.Config{
            SlowThreshold:             time.Second,
            LogLevel:                  logLevel,
            IgnoreRecordNotFoundError: true,
            Colorful:                  true,
        },
    )
    
    db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
        Logger: newLogger,
    })
    if err != nil {
        return nil, err
    }
    
    // 获取底层 sql.DB 对象进行连接池配置
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // 连接池优化配置
    sqlDB.SetMaxOpenConns(25)                    // 最大连接数
    sqlDB.SetMaxIdleConns(10)                    // 最大空闲连接数
    sqlDB.SetConnMaxLifetime(5 * time.Minute)    // 连接最大生命周期
    sqlDB.SetConnMaxIdleTime(2 * time.Minute)    // 连接最大空闲时间
    
    // WAL模式优化
    _, err = sqlDB.Exec("PRAGMA wal_checkpoint(TRUNCATE)")
    if err != nil {
        log.Printf("WAL checkpoint failed: %v", err)
    }
    
    // 自动迁移数据库表结构
    // db.AutoMigrate(&model.Room{}, &model.RoomMember{}, &model.UserSession{})
    
    return db, nil
}

// 读写分离配置（为后续扩展预留）
type DatabaseConfig struct {
    Master *gorm.DB
    Slave  *gorm.DB
}

func InitReadWriteDB() (*DatabaseConfig, error) {
    master, err := InitDB()
    if err != nil {
        return nil, err
    }
    
    // MVP阶段主从使用同一数据库实例
    // 为后续扩展预留接口
    return &DatabaseConfig{
        Master: master,
        Slave:  master,
    }, nil
}
```

**Repository模式实现**：

```go
// 匿名用户服务接口
type AnonymousUserService interface {
    CreateSession() (*model.UserSession, string, error) // 返回 Session 和 Token
    ValidateToken(token string) (*model.UserSession, error)
}

// 房间仓库接口
type RoomRepository interface {
    Create(room *model.Room) error
    FindByID(id string) (*model.Room, error)
    Search(keyword string, page, size int) ([]*model.Room, int64)
    Update(room *model.Room) error
    Delete(id string) error
    FindActiveRooms() ([]*model.Room, error)
}

// 房间成员仓库接口
type RoomMemberRepository interface {
    Join(member *model.RoomMember) error
    Leave(roomID, sessionID string) error
    FindMembers(roomID string) ([]*model.RoomMember, error)
    Update(member *model.RoomMember) error
    FindBySession(sessionID string) (*model.RoomMember, error)
}
```

**数据写入优化策略 (Write Debounce)**

**1. 内存缓冲机制**
- **高频数据缓存**：`current_time`、`playback_state`、`last_active_at` 等高频变动数据，主要维护在Go内存结构体中
- **并发安全**：使用 `sync.RWMutex` 保证内存缓存的并发安全访问
- **缓存结构设计**：
  ```go
  type RoomStateCache struct {
      mu            sync.RWMutex
      lastWriteTime time.Time    // 上次写入SQLite时间
      dirty         bool         // 标记数据是否已修改
      currentTime   float64      // 当前播放时间
      playbackState string       // 播放状态 (playing/paused)
      roomID        string       // 房间ID
  }
  ```

**2. 异步落盘策略 (Debounce)**
- **定时同步**：每30秒将内存状态批量写入SQLite（通过定时器触发）
- **事件触发**：以下关键事件发生时立即写入：
  - 用户暂停播放 (`MSG_PLAYBACK` with `action: "pause"`)
  - 用户执行Seek (`MSG_SEEK`)
  - 用户离开房间 (`MSG_USER_LEFT`)
  - 房间销毁前（自动清理时）
- **批量写入**：将多个房间的状态变更合并为一次SQLite事务

**3. 性能收益分析**
- **磁盘IOPS降低**：从每秒多次写入降低到每30秒一次批量写入
- **SQLite锁竞争减少**：减少并发写入的锁竞争，提升系统并发能力
- **响应时间优化**：播放控制指令立即在内存中生效，无需等待磁盘写入
- **数据安全性**：关键事件（暂停、离开）仍保证立即持久化

**4. 技术实现示例**
```go
// RoomStateManager 房间状态管理器
type RoomStateManager struct {
    cache map[string]*RoomStateCache
    mu    sync.RWMutex
    ticker *time.Ticker
}

// UpdateCurrentTime 更新当前时间（高频调用）
func (m *RoomStateManager) UpdateCurrentTime(roomID string, currentTime float64) {
    m.mu.Lock()
    defer m.mu.Unlock()
    
    cache, exists := m.cache[roomID]
    if !exists {
        cache = &RoomStateCache{roomID: roomID}
        m.cache[roomID] = cache
    }
    
    cache.mu.Lock()
    cache.currentTime = currentTime
    cache.dirty = true
    cache.mu.Unlock()
}

// FlushToDatabase 定时刷写到数据库
func (m *RoomStateManager) FlushToDatabase() {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for _, cache := range m.cache {
        cache.mu.RLock()
        if cache.dirty && time.Since(cache.lastWriteTime) > 30*time.Second {
            // 异步写入数据库
            go m.persistToDatabase(cache)
        }
        cache.mu.RUnlock()
    }
}

// ForceFlush 强制将所有脏数据刷写到数据库（用于程序退出）
func (m *RoomStateManager) ForceFlush() {
    m.mu.RLock()
    defer m.mu.RUnlock()
    
    for _, cache := range m.cache {
        cache.mu.RLock()
        if cache.dirty {
            // 同步写入数据库，确保数据持久化
            m.persistToDatabase(cache)
        }
        cache.mu.RUnlock()
    }
}

// SetupSignalHandler 设置信号处理器，捕获退出信号
func (m *RoomStateManager) SetupSignalHandler() {
    c := make(chan os.Signal, 1)
    // 捕获 SIGINT (Ctrl+C) 和 SIGTERM (Docker容器停止)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
    
    go func() {
        <-c
        log.Println("Received shutdown signal, flushing all data...")
        m.ForceFlush()
        log.Println("All data flushed, exiting...")
        os.Exit(0)
    }()
}
```

**5. 防抖配置参数**
- **定时同步间隔**：30秒（可配置）
- **关键事件列表**：暂停、跳转、离开、销毁
- **最大脏数据时间**：60秒（超过此时间强制写入）
- **批量写入大小**：最多合并10个房间状态变更
- **强制刷写机制**：
  - 捕获SIGINT (Ctrl+C)和SIGTERM (Docker容器停止)信号
  - 程序退出前强制将所有脏数据同步写入数据库
  - 防止重启导致的数据回滚和丢失
  - 使用同步写入确保数据持久化

### 2.4 房间生命周期管理 (Room Lifecycle)

**自动销毁策略**：

1.  **活跃判定**：
    *   WebSocket 连接数 > 0
    *   或 `last_active_at` 在 10 分钟内（心跳更新）

2.  **清理机制**：
    *   **SQLite清理**：房间状态在数据库中维护，通过 `last_active_at` 字段判断活跃状态
    *   **Cron Job (核心)**：每分钟运行一次清理任务，物理删除 `last_active_at < NOW() - 10m` 的房间记录

3.  **数据清理**：
    *   由于 `messages` 不持久化，无需清理。
    *   `room_members` 通过外键级联删除 (ON DELETE CASCADE) 自动清理。

---

## 三、API接口规范设计

### 3.1 RESTful API设计

**基础信息**：
- **Base URL**: `http://localhost:8080/api/v1`
- **认证方式**: Bearer Token (Session Token)
- **Content-Type**: `application/json`

**核心API端点**：

```go
// 房间管理 API
POST   /api/v1/rooms              // 创建房间
GET    /api/v1/rooms              // 获取房间列表
GET    /api/v1/rooms/:room_id     // 获取房间详情
PUT    /api/v1/rooms/:room_id     // 更新房间信息
DELETE /api/v1/rooms/:room_id     // 删除房间

// 房间成员管理
POST   /api/v1/rooms/:room_id/join    // 加入房间
POST   /api/v1/rooms/:room_id/leave   // 离开房间
GET    /api/v1/rooms/:room_id/members // 获取房间成员

// 用户会话管理
POST   /api/v1/sessions           // 创建临时会话
GET    /api/v1/sessions/:token    // 验证会话
```

### 3.2 核心API接口

#### 房间生命周期管理逻辑

**1. 房间创建与返回机制**
- **LocalStorage检查**：当用户访问首页时，前端会检查LocalStorage中是否存在RoomID，且最后访问时间在15分钟内
- **动态提示**：若存在有效房间记录，显示"回到房间 XXXXXX"动态提示
- **分支处理**：
  - 点击"创建"：生成新房间
  - 点击"回到"：跳转旧房间（若已销毁则引导重建）

**2. 房间过期处理**
- **自动销毁**：无人在线10分钟后自动销毁
- **过期提示**：用户尝试回到已过期房间时，显示"房间已过期，是否创建新房间？"
- **一键重建**：提供快速重建房间的选项，保留用户习惯设置

**3. 房间状态管理**
- **Active状态**：房间内至少有1名在线用户
- **Inactive状态**：房间内无用户，超过10分钟后自动销毁
- **状态更新**：每次用户加入/离开时更新房间状态

**1. 创建房间接口**

```go
type CreateRoomRequest struct {
    Name        string  `json:"name" binding:"required" example:"我的观影房间"`
    Description string  `json:"description" example:"和朋友一起看电影"`
    MediaURL    string  `json:"media_url" binding:"required" example:"https://example.com/video.mp4"`
    IsPrivate   bool    `json:"is_private" example:"false"`
    Password    string  `json:"password" example:""`
    MaxUsers    int     `json:"max_users" binding:"min=1,max=100" example:"10"`
    Settings    RoomSettings `json:"settings"`
}

type RoomSettings struct {
    AutoPlay      bool    `json:"auto_play" example:"true"`      // 自动播放
    AllowControl  bool    `json:"allow_control" example:"true"`  // 允许控制
    SyncTolerance float64 `json:"sync_tolerance" example:"2.0"`   // 同步容忍度(秒)
}
```

**响应示例**：
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "我的观影房间",
        "media_url": "https://example.com/video.mp4",
        "status": "active",
        "playback_state": "paused",
        "current_time": 0,
        "duration": 3600,
        "created_at": "2024-01-01T12:00:00Z"
    }
}
```

**2. 加入房间接口**

```go
type JoinRoomRequest struct {
    Password string `json:"password" example:""`
}

type JoinRoomResponse struct {
    Room      *model.Room       `json:"room"`
    Members   []*model.RoomMember `json:"members"`
    Session   *model.UserSession `json:"session"`
    WebSocketURL string         `json:"websocket_url"`
    ServerTime int64           `json:"server_time"` // 云端时钟时间戳(毫秒)
    PlaybackRate float64       `json:"playback_rate"` // 当前播放速率（如1.0、1.5、2.0）
}
```

#### 访客加入流程优化

**1. 临时昵称分配机制**
- **自动分配**：用户加入房间时自动分配趣味昵称（如：潜水的章鱼、快乐的熊猫等）
- **存储位置**：存储在本地Cookie中，同一用户再次访问时保持一致
- **昵称列表**：预定义100+趣味昵称，避免重复

**2. 云端时钟机制**
- **设计原则**：服务器维护统一的"云端时钟"，所有同步基于服务器时间
- **时间获取**：
  - 进房时从JoinRoomResponse获取ServerTime
  - 定期通过WebSocket Calibration消息校准
- **播放速率同步**：
  - **设计原理**：播放速率是长期同步的关键因素，倍速差异会导致进度偏移随时间累积
  - **数学原理**：偏移量 = (rate_diff) × time_elapsed
    - 例如：1.5x vs 1.0x，10秒后累积5秒偏移
    - 偏移随时间线性增长，必须从源头消除
  - **实现要求**：
    - 服务器在返回时间信息时，**必须同时下发**当前的播放速率`playback_rate`
    - `JoinRoomResponse`和`Calibration`消息中均包含`playback_rate`字段
    - 新用户加入时立即应用房间倍速，避免进度偏离
    - 倍速变化时通过WebSocket实时广播，所有客户端同步更新
  - **同步保证**：确保所有客户端播放速率一致，维持长期同步精度
- **同步逻辑**：服务器基于云端时钟直接返回播放状态，不依赖房主客户端响应

**3. 进房同步体验保障**
- **精准对齐**：
  - 系统自动计算并抵消网络传输时延（RTT）
  - 访客获取进度时，服务器返回已校准的时间戳：`server_time + estimated_delay`
  - 前端基于校准后的时间戳直接Seek，实现毫米级对齐
- **拒绝等待房主**：
  - 访客进房获取进度时，**严禁**依赖房主客户端响应
  - 服务器独立维护房间状态，直接返回最新播放进度
  - 房主离线或卡顿时，新访客仍能正常进房同步
- **状态最终一致**：
  - 若在访客加载视频期间房间状态发生变化（如暂停、跳转）
  - 访客视频加载完毕后，**自动应用**最新的房间状态
  - 通过WebSocket实时推送状态更新，确保各端最终一致
- **网络时延计算**：
  - 使用多次ping-pong测量计算平均RTT
  - 考虑时钟漂移和网络抖动
  - 动态调整校准系数，适应不同网络环境

**4. 预加载与对齐**
- **静音加载**：播放器静音加载视频，避免浏览器自动播放限制
- **直接Seek**：缓冲就绪(onCanPlay)后，直接跳转(Seek)至服务器返回的目标时间点
- **状态处理**：画面已与房间同步，处于静音/暂停状态，等待用户激活

**5. 解除限制与播放**
- **遮罩引导**：显示全屏遮罩，提示"点击加入观影"
- **交互授权**：用户点击后获取播放权限，解除静音并应用房间状态
- **状态反馈**：信号灯变为🟢(绿灯)，表示已就绪

### 3.3 WebSocket消息协议

**连接建立**：
- **WebSocket URL**: `ws://localhost:8080/ws/room/:room_id?token=:session_token`
- **协议版本**: v1
- **心跳间隔**: 30秒

#### 多标签页冲突处理

**1. 冲突检测机制**
- **机制**：使用浏览器`localStorage`实现标签页间通信
- **实现**：
  - 连接WebSocket前，在localStorage中设置标记：`room_id:session_id:active`
  - 监听localStorage变化事件
  - 新标签页连接成功后，旧标签页检测到变化

**2. 强制断开旧连接**
- **逻辑**：
  - 新标签页连接WebSocket成功
  - 更新localStorage中的活跃状态
  - 旧标签页检测到状态变化
  - 旧标签页主动断开WebSocket连接
  - 旧标签页显示提示："已在别处打开，本页面已断开连接"

**3. 技术实现示例**
```javascript
// 监听localStorage变化
window.addEventListener('storage', (event) => {
  if (event.key === `xiaowo_active_session_${roomId}`) {
    const newSession = event.newValue;
    if (newSession && newSession !== currentSessionId) {
      // 断开旧连接
      ws.close();
      showNotification('已在别处打开，本页面已断开连接');
    }
  }
});

// 连接成功后更新状态
ws.onopen = () => {
  localStorage.setItem(`xiaowo_active_session_${roomId}`, currentSessionId);
};

// 关闭时清除状态
ws.onclose = () => {
  localStorage.removeItem(`xiaowo_active_session_${roomId}`);
};
```

**消息格式**：

```go
// 基础消息结构
type WSMessage struct {
    Type      string      `json:"type"`      // 消息类型
    RoomID    string      `json:"room_id"`   // 房间ID
    SessionID string      `json:"session_id"` // 会话ID
    Timestamp int64       `json:"timestamp"` // 时间戳
    Data      interface{} `json:"data"`      // 消息数据
}

// 消息类型定义
const (
    MSG_PING           = "ping"           // 心跳
    MSG_PONG           = "pong"           // 心跳响应
    MSG_CALIBRATION    = "calibration"    // 时间对时
    MSG_PLAYBACK       = "playback"       // 播放控制 (必须为绝对状态: play/pause)
    MSG_SEEK           = "seek"           // 跳转到指定时间
    MSG_SPEED_CHANGE   = "speed_change"   // 播放速度变更 (必须为绝对倍速值)
    MSG_SYNC           = "sync"           // 同步状态
    MSG_CHAT           = "chat"           // 聊天消息
    MSG_USER_JOINED    = "user_joined"    // 用户加入
    MSG_USER_LEFT      = "user_left"      // 用户离开
    MSG_BUFFER_START   = "buffer_start"   // 缓冲开始
    MSG_BUFFER_END     = "buffer_end"     // 缓冲结束
    MSG_ROOM_INFO      = "room_info"      // 房间信息更新
    MSG_READY          = "ready"          // 视频加载就绪
    MSG_ROOM_READY_STATUS = "room_ready_status" // 房间就绪状态通知
```

#### 冲突处理规则

**1. 防止指令回声 (Anti-Echo)**
- **机制**：收到同步指令执行动作后，严禁再次向外发送广播信号
- **实现**：
  - 每条指令包含SessionID，接收方检查SessionID
  - 若指令来自自身，直接执行，不广播
  - 若指令来自他人，执行后不广播
- **目的**：阻断无限循环，避免网络风暴

**2. 后发指令优先 (Latest Wins)**
- **机制**：多人操作时，系统依据动作发生的时间戳裁决
- **实现**：
  - 每条指令包含精确的Timestamp
  - 接收方比较指令时间戳与本地最新执行的指令时间戳
  - 仅执行时间戳更新的指令
- **目的**：确保最终状态一致，避免冲突

**3. 拖拽时的互不干扰 (Drag Isolation)**
- **机制**：在用户拖拽进度条时，不广播中间过程
- **实现**：
  - 拖拽过程中，仅在本地更新进度
  - 松手确认后，才发送最终的Seek指令
  - 其他端保持原状，直到收到最终指令
- **目的**：避免其他用户看到频繁跳变，提升体验

#### 缓冲保护与追帧协作逻辑

**1. 缓冲检测机制**
- **监控指标**：前端持续监控 `video.buffered` 状态和网络缓冲情况
- **阈值设定**：缓冲状态持续超过 **2秒** 触发缓冲事件
- **区分逻辑**：系统根据用户身份（房主/普通成员）执行不同策略

**2. 房主/多数人卡顿处理 (Global Pause)**
- **触发条件**：房主进入缓冲状态超过2秒
- **处理逻辑**：
  - 房主发送 `MSG_BUFFER_START` 消息
  - 服务器广播暂停指令给所有房间成员
  - 全员暂停播放，等待缓冲恢复
  - 房主缓冲恢复后发送 `MSG_BUFFER_END` 消息
  - 服务器广播继续播放指令
- **设计原则**：房主卡顿需全员暂停，保障观影体验

**3. 个别访客卡顿处理 (Local Pause)**
- **触发条件**：非房主普通成员进入缓冲状态
- **处理逻辑**：
  - 成员本地暂停播放，等待缓冲恢复
  - **严禁**发送缓冲消息影响其他用户
  - 缓冲恢复后，计算时间差 `time_diff`
  - 根据时间差自动触发追帧策略：
    - `time_diff < 3s`：温柔微调 (1.05x倍速)
    - `time_diff ≥ 3s`：幽灵模式 (1.5x-2.0x倍速)
- **设计原则**：个别用户问题不影响整体，快速追赶进度

**4. 恢复播放后的智能追帧**
- **触发时机**：缓冲结束后恢复播放瞬间
- **追帧逻辑**：
  - 系统自动计算各端间微小误差
  - 智能追帧立即介入，通过微调倍速无感抹平误差
  - 优先使用"温柔微调"策略，避免画面跳变
  - 若误差较大，自动升级为"快速追赶"策略
- **设计原则**：无感同步，保持观影沉浸感

**5. 技术实现要点**
- **前端监控**：使用 `video.readyState` 和 `video.buffered` API
- **消息传递**：房主缓冲消息必须广播，成员缓冲消息本地处理
- **状态恢复**：缓冲结束后需同步最新房间状态
- **追帧策略**：与四级同步策略共享追帧逻辑

#### 操作界限明确化

**1. 全局共享操作 (必须全员同步)**
- **播放/暂停**：`MSG_PLAYBACK` 消息，action必须为绝对状态 (`play`/`pause`)
- **进度跳转**：`MSG_SEEK` 消息，指定绝对时间戳
- **倍速设置**：`MSG_SPEED_CHANGE` 消息，指定绝对倍速值
- **设计原则**：影响画面内容与节奏的操作必须全员一致

**2. 本地独享操作 (严禁同步)**
- **音量大小**：本地设备音量控制
- **静音开关**：本地静音状态
- **全屏模式**：本地全屏显示
- **字幕开关**：本地字幕显示
- **画面比例**：本地播放器缩放
- **设计原则**：不影响他人观影体验的操作保持独立

**3. 提示防遮挡规则 (UI安全区域)**
- **禁止区域**：屏幕底部 **20%** 的区域
- **设计原则**：确保系统提示（如"User A暂停了"）不遮挡视频字幕
- **实现要求**：
  - 所有系统提示必须出现在屏幕顶部或中部
  - 鼠标悬停/点击信号灯时，延迟数值显示在侧边
  - 全局消息反馈严禁遮挡视频内容区域

**4. 消息广播范围定义**
- **全员广播**：`MSG_PLAYBACK`, `MSG_SEEK`, `MSG_SPEED_CHANGE`, `MSG_BUFFER_START` (房主), `MSG_BUFFER_END` (房主), `MSG_READY`, `MSG_ROOM_READY_STATUS`
- **本地处理**：`MSG_BUFFER_START` (成员), `MSG_BUFFER_END` (成员), 所有音量/静音操作
- **服务器推送**：`MSG_ROOM_INFO`, `MSG_USER_JOINED`, `MSG_USER_LEFT`

**5. 指令状态要求**
- **绝对状态**：所有同步指令必须是绝对状态（如"设为播放"、"跳转到120.5秒"、"设为1.5倍速"）
- **严禁相对动作**：禁止使用"切换状态"、"增加0.1倍速"等相对指令
- **幂等性保证**：多次收到相同指令必须产生相同结果

**核心消息示例**：

```go
// 1. 心跳消息 (每30秒)
{
    "type": "ping",
    "room_id": "room-123",
    "session_id": "session-456", 
    "timestamp": 1640995200000
}

// 2. 时间对时消息 (进房时 + 每5分钟)
{
    "type": "calibration",
    "room_id": "room-123",
    "session_id": "session-456",
    "timestamp": 1640995200000,
    "data": {
        "client_time": 1640995200000,
        "server_time": 1640995200050,
        "playback_rate": 1.0  // 当前播放速率，新用户加入时必须同步
    }
}

// 3. 播放控制消息
{
    "type": "playback",
    "room_id": "room-123", 
    "session_id": "session-456",
    "timestamp": 1640995200000,
    "data": {
        "action": "play", // play/pause
        "current_time": 120.5,
        "duration": 3600.0
    }
}

// 4. 同步状态消息
{
    "type": "sync",
    "room_id": "room-123",
    "session_id": "session-456", 
    "timestamp": 1640995200000,
    "data": {
        "current_time": 120.5,
        "playback_state": "playing",
        "playback_rate": 1.0
    }
}

// 5. 视频就绪消息
{
    "type": "ready",
    "room_id": "room-123",
    "session_id": "session-456",
    "timestamp": 1640995200000,
    "data": {
        "is_ready": true
    }
}

// 6. 房间就绪状态通知
{
    "type": "room_ready_status",
    "room_id": "room-123",
    "session_id": "",
    "timestamp": 1640995200000,
    "data": {
        "total_members": 3,
        "ready_members": 2,
        "member_statuses": {
            "session-456": true,  // 已就绪
            "session-789": true,  // 已就绪
            "session-012": false  // 未就绪
        }
    }
}
```

**WebSocket Hub 核心实现**：

```go
// WebSocketHub WebSocket连接管理
type WebSocketHub struct {
    rooms     map[string]*Room
    conns     map[*WebSocketConnection]bool
    broadcast chan []byte
    register  chan *WebSocketConnection
    unregister chan *WebSocketConnection
    
    // 连接统计与限制
    totalConnections    int32
    maxConnections     int32
    maxConnectionsPerIP int32
    
    // 消息频率限制
    rateLimiters map[string]*rate.Limiter
    
    // 定时任务
    cleanupTicker *time.Ticker
    heartbeatTicker *time.Ticker
    
    mu sync.RWMutex
}

// Room 房间结构
type Room struct {
    ID       string                `json:"id"`
    members  map[string]*WebSocketConnection
    mu       sync.RWMutex
    version  int                   // 版本号(乐观锁)
}

// WebSocketConnection WebSocket连接
type WebSocketConnection struct {
    roomID    string
    sessionID string
    conn      *websocket.Conn
    send      chan []byte
    ip        string              // 客户端IP地址
    mu        sync.RWMutex
    lastSeen  time.Time
    rtts      []int64              // 最近3次RTT记录
    closed    bool                 // 连接关闭标记
}

// NewWebSocketHub 创建WebSocket Hub
func NewWebSocketHub() *WebSocketHub {
    hub := &WebSocketHub{
        rooms:     make(map[string]*Room),
        conns:     make(map[*WebSocketConnection]bool),
        broadcast: make(chan []byte, 1024),
        register:  make(chan *WebSocketConnection),
        unregister: make(chan *WebSocketConnection),
        rateLimiters: make(map[string]*rate.Limiter),
        maxConnections: 2000,          // 最大连接数
        maxConnectionsPerIP: 10,       // 每IP最大连接数
    }
    
    // 初始化定时任务
    hub.cleanupTicker = time.NewTicker(1 * time.Minute)
    hub.heartbeatTicker = time.NewTicker(30 * time.Second)
    
    return hub
}

// Start 启动WebSocket Hub
func (h *WebSocketHub) Start() {
    // 启动定时任务
    go h.startCleanup()
    go h.startHeartbeat()
    
    for {
        select {
        case conn := <-h.register:
            if h.canAcceptConnection(conn) {
                h.addConnection(conn)
            } else {
                conn.conn.Close()
            }
        case conn := <-h.unregister:
            h.removeConnection(conn)
        case message := <-h.broadcast:
            h.broadcastMessage(message)
        }
    }
}

// 连接限制检查
func (h *WebSocketHub) canAcceptConnection(conn *WebSocketConnection) bool {
    h.mu.RLock()
    defer h.mu.RUnlock()
    
    // 检查总连接数
    if atomic.LoadInt32(&h.totalConnections) >= h.maxConnections {
        return false
    }
    
    // 检查IP级别的连接数
    ipConnections := h.countConnectionsByIP(conn.ip)
    if ipConnections >= h.maxConnectionsPerIP {
        return false
    }
    
    return true
}

// 按IP统计连接数
func (h *WebSocketHub) countConnectionsByIP(ip string) int {
    count := 0
    for conn := range h.conns {
        if conn.ip == ip {
            count++
        }
    }
    return count
}

// addConnection 添加连接
func (h *WebSocketHub) addConnection(conn *WebSocketConnection) {
    room := h.getOrCreateRoom(conn.roomID)
    
    room.mu.Lock()
    room.members[conn.sessionID] = conn
    room.version++
    room.mu.Unlock()
    
    h.mu.Lock()
    h.conns[conn] = true
    atomic.AddInt32(&h.totalConnections, 1)
    
    // 为每个连接创建速率限制器
    if _, exists := h.rateLimiters[conn.sessionID]; !exists {
        h.rateLimiters[conn.sessionID] = rate.NewLimiter(rate.Limit(10), 20) // 每秒10条消息，突发20条
    }
    h.mu.Unlock()
}

// removeConnection 移除连接  
func (h *WebSocketHub) removeConnection(conn *WebSocketConnection) {
    h.mu.Lock()
    if _, ok := h.conns[conn]; ok {
        delete(h.conns, conn)
        atomic.AddInt32(&h.totalConnections, -1)
        
        // 清理速率限制器
        delete(h.rateLimiters, conn.sessionID)
        
        // 标记连接为关闭
        conn.closed = true
        close(conn.send)
        
        room := h.getOrCreateRoom(conn.roomID)
        room.mu.Lock()
        delete(room.members, conn.sessionID)
        room.version++
        room.mu.Unlock()
    }
    h.mu.Unlock()
}

// 定期清理无效连接
func (h *WebSocketHub) startCleanup() {
    for range h.cleanupTicker.C {
        h.mu.Lock()
        for conn := range h.conns {
            if conn.closed {
                delete(h.conns, conn)
                atomic.AddInt32(&h.totalConnections, -1)
                delete(h.rateLimiters, conn.sessionID)
            }
        }
        h.mu.Unlock()
    }
}

// 心跳检测机制
func (h *WebSocketHub) startHeartbeat() {
    for range h.heartbeatTicker.C {
        h.mu.Lock()
        for conn := range h.conns {
            if time.Since(conn.lastSeen) > 90*time.Second {
                // 连接超时，关闭连接
                conn.closed = true
                conn.conn.Close()
                delete(h.conns, conn)
                atomic.AddInt32(&h.totalConnections, -1)
                delete(h.rateLimiters, conn.sessionID)
            } else {
                // 更新最后活跃时间
                conn.mu.Lock()
                conn.lastSeen = time.Now()
                conn.mu.Unlock()
                
                // 发送心跳
                select {
                case conn.send <- []byte("ping"):
                default:
                    // 发送失败，关闭连接
                    conn.closed = true
                    conn.conn.Close()
                    delete(h.conns, conn)
                    atomic.AddInt32(&h.totalConnections, -1)
                    delete(h.rateLimiters, conn.sessionID)
                }
            }
        }
        h.mu.Unlock()
    }
}

// 检查消息频率限制
func (h *WebSocketHub) checkRateLimit(sessionID string) bool {
    h.mu.RLock()
    limiter, exists := h.rateLimiters[sessionID]
    h.mu.RUnlock()
    
    if !exists {
        return true // 限制器不存在，允许通过
    }
    
    return limiter.Allow()
}

// calculateSmoothedRTT 计算平滑RTT
func (c *WebSocketConnection) calculateSmoothedRTT(newRTT int64) int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.rtts = append(c.rtts, newRTT)
    if len(c.rtts) > 3 {
        c.rtts = c.rtts[1:] // 保持最近3次
    }
    if len(c.rtts) < 3 {
        return newRTT // 不足3次直接返回
    }
    
    // 使用中位数平滑处理
    sortedRTTs := make([]int64, len(c.rtts))
    copy(sortedRTTs, c.rtts)
    sort.Slice(sortedRTTs, func(i, j int) bool { return sortedRTTs[i] < sortedRTTs[j] })
    return sortedRTTs[len(sortedRTTs)/2] // 中位数
}

// handleMessage 处理WebSocket消息
func (h *WebSocketHub) handleMessage(conn *WebSocketConnection, message []byte) {
    // 1. 检查消息频率限制
    if !h.checkRateLimit(conn.sessionID) {
        conn.send <- []byte(`{"type":"error","code":"rate_limited","message":"消息发送过于频繁"}`)
        return
    }
    
    // 2. 更新最后活跃时间
    conn.mu.Lock()
    conn.lastSeen = time.Now()
    conn.mu.Unlock()
    
    // 3. 根据消息类型处理
    // （此处省略具体消息处理逻辑）
}

// handleSync 处理同步消息
func (h *WebSocketHub) handleSync(conn *WebSocketConnection, message []byte) {
    var syncMsg v1.SyncMessage
    if err := json.Unmarshal(message, &syncMsg); err != nil {
        h.sendError(conn, "sync_failed", "同步消息格式错误")
        return
    }

    room := h.getOrCreateRoom(syncMsg.RoomID)
    room.mu.Lock()
    defer room.mu.Unlock()

    currentTime := syncMsg.Data.CurrentTime
    targetTime := h.calculateTargetTime(room)
    timeDiff := targetTime - currentTime

    var action SyncAction
    switch {
    case timeDiff > 2.0: // Level 2: 误差 > 2s，触发 Seek
        action = SyncAction{
            Type:       "seek",
            TargetTime: targetTime,
            Reason:     "large_difference",
        }
    case timeDiff > 0.5: // Level 3: 误差 0.5s - 2s，动态调整倍速
        speed := 1.0 + math.Min(timeDiff/10.0, 0.5) // 最大1.5倍速
        action = SyncAction{
            Type:          "playback_rate",
            PlaybackRate:  speed,
            TargetTime:    targetTime,
            Reason:        "medium_difference",
        }
    case timeDiff < -0.5: // 客户端超前，减速
        speed := 1.0 + math.Max(timeDiff/10.0, -0.3) // 最小0.7倍速
        action = SyncAction{
            Type:          "playback_rate", 
            PlaybackRate:  speed,
            TargetTime:    targetTime,
            Reason:        "client_ahead",
        }
    default: // Level 4: 误差 < 0.5s，忽略
        return
    }

    // 广播同步指令给房间内其他用户
    h.broadcastToRoom(room.ID, action)
}

// 四级同步策略说明 (完全对齐PRD):
// 完美同步 (🟢): |time_diff| < 0.05s → 忽略 (画面完全一致，无动作)
// 温柔微调 (🟡): 0.05s ≤ |time_diff| < 3s → 动态变速 (落后者1.05x，超前者0.95x，无感追平)
// 快速追赶 (🔴): 3s ≤ |time_diff| < 10s → 幽灵模式 (1.5x-2.0x倍速快速追赶)
// 强制对齐 (⚫): |time_diff| ≥ 10s → 强制跳转 (伴随Loading，直接跳至最新位置)

// 补充规则1: 长时自动校准
// 系统每隔30秒会自动进行一次静默对表，若发现偏差超过阈值，自动触发"温柔微调"

// 补充规则2: 倍速叠加逻辑
// 基准倍速: 用户手动设置的倍速为全局基准，全员同步
// 最终倍速: 系统的追帧微调是在基准倍速上进行乘法叠加
// 公式: 最终播放速度 = 用户基准倍速 × 追帧系数

// 前端播放器防坑指南：浏览器Autoplay Policy处理
// 问题背景：现代浏览器严格限制自动播放媒体内容，需要用户交互信任才能播放视频
// 解决方案：实现"点击加入"(Click to Join)遮罩层机制
// 
// 实施步骤：
// 1. 用户进入房间时显示遮罩层，提示"点击加入房间"
// 2. 用户点击"加入"按钮后，立即执行 video.play() 然后 video.pause()
// 3. 这一步是为了获取浏览器的"用户交互信任"
// 4. 之后WebSocket收到play指令时，JS才能正常自动播放视频
// 5. 遮罩层消失，开始正常的同步播放流程
//
// 技术实现示例：
// function handleUserJoin() {
//   const video = document.getElementById('video-player');
//   
//   // iOS Safari 音频自动播放限制处理：必须在用户交互时解除静音
//   video.muted = false; // 关键：解除静音状态，否则即使用户交互了也可能无声
//   
//   video.play().then(() => {
//     video.pause(); // 立即暂停，但已获取播放权限
//     hideJoinOverlay(); // 隐藏遮罩层
//     setWebSocketReady(true); // 标记WebSocket可以控制播放
//   }).catch(error => {
//     console.error('获取播放权限失败:', error);
//   });
// }
```

---

## 四、开发与部署策略

### 4.1 MVP功能范围

**核心功能 (必须)**：
- ✅ 房间创建与管理
- ✅ 临时用户会话 (无需注册)
- ✅ WebSocket实时同步
- ✅ 基础播放控制 (播放/暂停/跳转/倍速)
- ✅ 简单的聊天功能 (调试用)
- ✅ 房间自动清理
- ✅ 移动端适配

**功能预留 (后期)**：
- 🔄 播放历史记录
- 🔄 房间收藏夹
- 🔄 社交功能 (好友/群组)
- 🔄 高级同步设置

### 4.1.1 移动端适配设计

**布局逻辑**：
- **竖屏模式 (Portrait)**：采用"上视下控"结构
  - 上部 (View Zone)：视频播放器固定在顶部（保持16:9比例）
  - 中部 (Info Zone)：房间信息、成员头像、信号灯
  - 下部 (Control Zone)：大尺寸链接输入框、加载按钮、预设片源
  - 重点：输入框需位于屏幕下半部分，防止键盘弹起时遮挡视频

- **横屏模式 (Landscape)**：自动进入全屏沉浸模式
  - 隐藏所有非必要UI（输入框、成员列表）
  - 仅保留播放控制栏
  - 支持手势操作（如双击暂停/播放）

**交互优化**：
- 提供大尺寸的触控区域（如屏幕中央点击暂停/播放）
- 替代PC端的空格键快捷操作
- 优化移动端触摸手势

**响应式设计**：
- 使用CSS Grid和Flexbox实现响应式布局
- 媒体查询断点：
  - 移动端：< 768px
  - 平板：768px - 1024px
  - 桌面：> 1024px

**技术实现**：
- 使用Vue 3的响应式系统和计算属性
- 结合CSS媒体查询和现代CSS特性
- 利用Vite的CSS处理能力优化样式加载
- **iOS Safari全屏行为处理**：
  - **问题**：iOS Safari经常忽略`playsinline`属性（特别是旧版本或特殊设置下），强制全屏接管，导致"下控"区域不可见
  - **解决方案**：
    - 引入`iphone-inline-video` polyfill库确保内联播放
    - 在CSS中给`<video>`父容器设置`overflow: hidden`等黑魔法
    - 确保视频元素乖乖待在View Zone内，防止iOS Safari强制全屏
- **键盘弹出时的视口塌陷处理**：
  - **问题**：Android设备上软键盘弹出会导致视口高度缩小，使用100vh的布局会错乱
  - **CSS解决方案**：
    ```css
    /* 针对移动端的CSS优化 */
    @media (max-width: 768px) {
      .control-zone {
        /* 使用 svh (Small Viewport Height) 而不是 vh，自动适应键盘高度 */
        padding-bottom: env(safe-area-inset-bottom);
        min-height: 20svh;
      }
    }
    ```
  - **JavaScript解决方案**：
    ```javascript
    // 监听视觉视口变化，处理键盘弹起
    window.visualViewport.addEventListener('resize', () => {
      // 当键盘弹起，viewport高度变小，自动调整容器高度
      document.getElementById('app-container').style.height = 
          `${window.visualViewport.height}px`;
      
      // 可选：键盘弹起时临时隐藏非核心UI
      const isKeyboardVisible = window.innerHeight > window.visualViewport.height;
      if (isKeyboardVisible) {
        // 隐藏非核心UI，如成员列表
        document.getElementById('members-list').style.display = 'none';
      } else {
        // 恢复显示
        document.getElementById('members-list').style.display = 'block';
      }
    });
    ```

### 4.1.2 设备中断处理

**1. 页面不可见检测**
- **机制**：监听浏览器`visibilitychange`事件
- **触发条件**：
  - 用户切到其他标签页
  - 用户最小化浏览器窗口
  - 用户切换到其他应用
- **处理逻辑**：
  - 检测到页面不可见时，发送暂停指令给服务器
  - 服务器广播暂停指令给所有房间成员
  - 显示提示："User A 暂时离开，等待回来..."

**2. 系统级暂停处理**
- **机制**：监听设备级暂停事件
- **触发条件**：
  - 移动端设备锁屏
  - 移动端设备来电
  - 移动端设备闹钟响铃
- **处理逻辑**：
  - 检测到系统级暂停时，发送暂停指令
  - 恢复后发送提示："User A 已回来，继续播放？"
  - 等待用户确认后继续播放

**3. 网络恢复处理**
- **机制**：监听网络状态变化事件
- **触发条件**：
  - 网络从断开到恢复
  - 网络从弱到强
- **处理逻辑**：
  - 断网时发送暂停指令
  - 恢复后自动重连WebSocket
  - 同步最新进度并询问是否继续播放

**4. iOS Safari锁屏处理**
- **机制**：监听iOS Safari锁屏导致的WebSocket断开
- **触发条件**：
  - iOS 17.0+ 系统锁屏超过30秒
  - Safari进入后台模式导致WebSocket被系统强制断开
- **处理逻辑**：
  - 实现 `reconnect_on_visible` 机制：页面重新可见时自动重连
  - 结合 `visibilitychange` 事件，在页面可见时检查WebSocket连接状态
  - 重连成功后自动同步最新房间状态
  - 显示提示："网络已恢复，正在同步状态..."

### 4.1.3 预设演示片源

**1. 设计目的**
- 降低用户使用门槛
- 方便用户快速体验同步效果
- 提供无版权测试内容

**2. 内置片源列表**
- **片源1**：Big Buck Bunny (10分钟) - 经典开源动画
  - 链接：https://test-videos.co.uk/vids/bigbuckbunny/mp4/h264/1080/Big_Buck_Bunny_1080_10s_1MB.mp4
  - 特点：高清、短时长、快速加载

- **片源2**：Sintel (15分钟) - Blender基金会作品
  - 链接：https://test-videos.co.uk/vids/sintel/mp4/h264/720/Sintel_720_10s_1MB.mp4
  - 特点：高质量、情节完整

- **片源3**：Elephants Dream (10分钟) - 开源3D动画
  - 链接：https://test-videos.co.uk/vids/elephantsdream/mp4/h264/720/Elephants_Dream_720_10s_1MB.mp4
  - 特点：经典开源作品

- **片源4**：大闹天宫 (5分钟) - 中国经典动画片段
  - 链接：https://example.com/test-videos/danao.mp4
  - 特点：中国文化元素

- **片源5**：测试视频合集 (1分钟)
  - 链接：https://example.com/test-videos/short-collection.mp4
  - 特点：短时长、适合快速测试

**3. 技术实现**
- 在链接输入框下方添加"快速演示"标签页
- 点击片源自动填充到输入框并加载
- 提供片源预览和时长信息
- 支持一键切换不同片源

**4. 使用流程**
1. 用户进入房间
2. 点击"快速演示"标签
3. 选择一个片源
4. 系统自动填充并加载视频
5. 点击播放按钮开始同步体验

### 4.1.4 资源兼容性与稳定性

**1. HTTP/HTTPS 无感兼容**
- **问题背景**：现代浏览器限制在HTTPS页面中加载HTTP资源
- **解决方案**：
  - 实现后端视频链接代理服务
  - 所有视频请求通过后端代理转发
  - 自动将HTTP链接转换为HTTPS可访问的代理链接
- **技术实现**：
  ```go
  // 视频代理API
  GET /api/v1/proxy/video?url={encoded_video_url}
  ```
  - 后端验证链接安全性
  - 流式转发视频内容
  - 添加适当的CORS头部

**2. 防盗链适配**
- **检测机制**：
  - 尝试直接访问视频链接，检查HTTP状态码
  - 分析响应头中的`Content-Type`和`Content-Length`
  - 监控视频加载过程中的错误事件
- **友好提示**：
  - 检测到防盗链限制时，显示通俗提示："该视频链接受到保护，无法在小窝播放"
  - 提供"尝试其他链接"按钮，引导用户使用预设片源
  - 记录防盗链域名，优化后续兼容策略

**3. 前端视频源预检机制**

**设计目的**：防止无效链接或CORS受限链接污染房间状态，避免用户点击播放后导致全员黑屏/报错

**预检流程**：
1. **触发时机**：用户输入链接并点击"加载"按钮时
2. **隐藏测试**：前端创建一个隐藏的`<video>`标签尝试加载链接
3. **事件监听**：
   - 监听`canplay`事件：加载成功 → 发送UpdateRoom请求给服务器，更新房间视频链接
   - 监听`error`事件：加载失败 → 提示用户链接无效或存在跨域问题
4. **结果处理**：
   - **验证成功**：链接可用，广播给房间其他成员
   - **验证失败**：仅在当前用户端提示错误，不污染房间状态

**技术实现示例**：
```javascript
function precheckVideoSource(url) {
  return new Promise((resolve, reject) => {
    // 混合内容检查：HTTPS页面不允许加载HTTP视频
    if (window.location.protocol === 'https:' && url.startsWith('http:')) {
      reject({
        success: false,
        url: url,
        error: new Error('MIXED_CONTENT_ERROR'),
        message: 'HTTPS页面无法加载HTTP视频，请使用HTTPS链接或尝试通过服务器中转'
      });
      return;
    }
    
    const testVideo = document.createElement('video');
    testVideo.style.display = 'none';
    testVideo.preload = 'metadata'; // 只加载元数据，减少流量
    
    // 成功回调
    testVideo.addEventListener('canplay', () => {
      document.body.removeChild(testVideo);
      resolve({
        success: true,
        url: url,
        message: '视频链接验证通过'
      });
    });
    
    // 失败回调
    testVideo.addEventListener('error', (e) => {
      document.body.removeChild(testVideo);
      reject({
        success: false,
        url: url,
        error: e,
        message: '视频链接无效或存在跨域限制'
      });
    });
    
    // 超时处理（5秒）
    const timeoutId = setTimeout(() => {
      document.body.removeChild(testVideo);
      reject({
        success: false,
        url: url,
        message: '视频链接加载超时'
      });
    }, 5000);
    
    testVideo.addEventListener('canplay', () => clearTimeout(timeoutId));
    testVideo.addEventListener('error', () => clearTimeout(timeoutId));
    
    // 开始测试
    document.body.appendChild(testVideo);
    testVideo.src = url;
  });
}

// 使用示例
document.getElementById('load-btn').addEventListener('click', async () => {
  const url = document.getElementById('video-url').value;
  
  try {
    const result = await precheckVideoSource(url);
    // 验证成功，更新房间
    updateRoomVideo(url);
    showSuccess('视频链接有效，正在加载...');
  } catch (error) {
    // 验证失败，提示用户
    showError(`${error.message}，请检查链接后重试`);
  }
});
```

**价值与收益**：
- **提升用户体验**：避免无效链接导致全员观影中断
- **降低服务器负载**：减少因无效链接产生的错误请求
- **简化故障排查**：明确区分链接问题与系统问题
- **保护房间状态**：防止单个用户的错误输入影响其他成员

**4. 视频格式兼容性**
- **支持格式**：
  - MP4 (H.264 + AAC)：主流兼容格式
  - M3U8 (HLS)：流媒体协议，支持直播和长视频
  - WebM (VP9)：开源格式，Chrome/Firefox支持
- **格式检测**：
  - 通过文件扩展名和Content-Type初步判断
  - 前端尝试加载并监听`canplay`事件
  - 加载失败时自动尝试备用解码器

**5. 混合内容安全策略**
- **CORS配置**：
  - 后端服务器配置适当的CORS头部
  - 允许前端域名访问视频资源
  - 支持凭证传递（Cookies）
- **内容安全策略 (CSP)**：
  - 配置CSP策略，允许视频资源加载
  - 限制不安全的内容加载
  - 提供详细错误报告

**6. 播放稳定性保障**
- **重试机制**：
  - 视频加载失败时自动重试（最多3次）
  - 重试间隔指数退避：1s, 2s, 4s
  - 重试失败后显示错误提示
- **降级策略**：
  - 高清源失败时尝试标清源
  - MP4失败时尝试HLS源
  - 提供"切换清晰度"选项
- **监控报警**：
  - 记录视频加载成功率
  - 监控平均加载时间
  - 设置失败率阈值报警

**7. 技术实现要点**
- **代理服务（MVP阶段建议）**：使用Go标准库`httputil.ReverseProxy`实现高效流式转发，但**MVP阶段默认关闭或仅限预设片源使用**，避免服务器带宽和流量被1080P视频瞬间耗尽
- **链接验证**：检查域名白名单，防止滥用
- **缓存策略**：合理缓存视频元数据，减少重复请求
- **错误处理**：统一的错误处理中间件，提供友好提示
- **流量控制与成本考虑**：
  - **MVP策略**：用户自定义链接直接前端播放，不经过服务器代理
  - **预设片源**：仅对内置演示片源启用代理，保证基础体验
  - **成本控制**：代理服务作为备用方案，仅在万不得已时使用
  - **监控告警**：设置流量阈值告警，防止意外流量费用

### 4.1.5 前端播放器逻辑优化

#### 4.1.5.1 视觉体验建议

**1. 深色/暗黑模式**
- **强制要求**：播放页必须使用深色/暗黑模式主题
- **设计原则**：降低环境光干扰，提升观影沉浸感
- **实现方式**：
  - 使用CSS变量定义颜色主题
  - 自动检测系统主题偏好
  - 提供手动切换选项（预留）

**2. 动效与微交互**
- **淡入淡出切换**：所有UI状态切换使用淡入淡出动画
- **信号灯呼吸动效**：网络状态信号灯添加呼吸动画效果
- **加载动画**：视频加载时显示优雅的加载动画
- **过渡效果**：页面跳转、模态框展示等添加平滑过渡

**3. 视觉反馈**
- **操作反馈**：用户操作后提供明确的视觉反馈
- **状态指示**：清晰指示播放状态、网络状态、同步状态
- **进度提示**：拖动进度条时显示预览画面和时间点

#### 4.1.5.2 异常流程缺省页

**1. 解析失败页面**
- **触发条件**：视频链接解析失败或无法播放
- **界面设计**：
  - 显示友好的错误图标和提示文字
  - 提供"重试"按钮，重新尝试加载
  - 显示"使用预设片源"选项，快速切换
  - 保留原始链接输入框，方便修改

**2. 房间过期页面**
- **触发条件**：房间已过期或被销毁
- **界面设计**：
  - 显示"房间已过期"提示
  - 提供"创建新房间"按钮，快速重建
  - 保留原房间设置（如房间名称、密码）
  - 显示历史房间记录（如有）

**3. 网络异常页面**
- **触发条件**：网络连接失败或超时
- **界面设计**：
  - 显示网络异常图标
  - 提供"重新连接"按钮
  - 显示当前网络状态诊断信息
  - 提供离线使用指引

#### 4.1.5.3 网络健康度反馈

**1. 信号灯设计**
- **🟢 绿灯**：网络良好，延迟 < 100ms
- **🟡 黄灯**：网络一般，100ms ≤ 延迟 < 300ms
- **🔴 红灯**：网络较差，延迟 ≥ 300ms
- **⚫ 黑灯**：网络断开

**2. 延迟数值显示**
- **鼠标悬停**：悬停信号灯显示精确延迟数值
- **点击查看**：点击信号灯显示详细网络状态
- **实时更新**：延迟数值每5秒更新一次
- **历史记录**：显示最近10次延迟测量结果

**3. 网络状态面板**
- **延迟趋势**：显示延迟变化趋势图
- **丢包率**：显示网络丢包率
- **带宽信息**：显示当前带宽使用情况
- **建议操作**：根据网络状态提供优化建议

#### 4.1.5.4 全局消息反馈

**1. 消息位置规则**
- **严禁遮挡**：所有系统消息严禁遮挡视频底部20%区域
- **顶部显示**：重要系统消息显示在屏幕顶部
- **侧边提示**：次要提示信息显示在屏幕侧边
- **自动消失**：非关键消息3秒后自动消失

**2. 消息内容规范**
- **用户行为**："User A 暂停了视频"、"User B 加入了房间"
- **系统状态**："正在同步..."、"网络连接恢复"
- **错误提示**："视频加载失败，请重试"
- **操作确认**："已跳转到 01:30"

**3. 交互反馈**
- **可操作性**：重要消息提供操作按钮（如"重试"、"确定"）
- **关闭选项**：所有消息提供关闭按钮
- **持久性**：关键错误消息需用户确认后才消失

#### 4.1.5.5 播放器核心逻辑

**1. 状态管理**
- **播放状态**：playing, paused, buffering, ended
- **网络状态**：online, offline, reconnecting
- **同步状态**：synced, syncing, out_of_sync
- **房间状态**：joined, leaving, reconnecting
- **就绪状态**：ready, not_ready, checking

**2. 事件处理**
- **用户事件**：play, pause, seek, ratechange
- **视频事件**：canplay, waiting, seeking, ended
- **网络事件**：online, offline, visibilitychange
- **同步事件**：sync_start, sync_complete, sync_error

**3. 错误恢复**
- **自动重试**：可恢复错误自动重试机制
- **降级处理**：功能不可用时的降级方案
- **用户引导**：无法自动恢复时引导用户操作

**4. 全员就绪（Ready Check）机制**
- **设计目的**：避免房主快速播放导致访客加载不及时的问题
- **前端逻辑**：
  - 视频元素触发`onCanPlay`事件时，发送`MSG_READY`消息到服务器
  - 接收`MSG_ROOM_READY_STATUS`消息，更新成员列表就绪状态（显示绿点）
  - 房主点击播放按钮时，检查房间就绪状态
  - 如果存在未就绪成员，弹出提示："还有小伙伴没加载好，确定要开始吗？"
- **后端逻辑**：
  - 接收`MSG_READY`消息，更新成员就绪状态
  - 广播`MSG_ROOM_READY_STATUS`消息给所有房间成员
  - 维护每个成员的就绪状态
- **UI设计**：
  - 成员头像旁显示绿点（就绪）或红点（未就绪）
  - 房主播放按钮悬停时显示未就绪成员数量
  - 弹出提示框提供"继续播放"和"再等一下"选项

**5. iOS Safari后台保活机制**
- **问题背景**：iOS Safari在用户切后台（如回微信）时会冻结WebSocket连接和视频播放，导致同步中断
- **解决方案**：在用户点击"加入"按钮时，初始化并播放一个空的音频对象（Silent Audio），欺骗系统认为网页正在播放重要媒体
- **技术实现**：
  ```javascript
  // Hack for iOS Safari to keep connection alive in background
  const silentAudio = new Audio("data:audio/wav;base64,UklGRigAAABXQVZFZm10IBIAAAABAAEARKwAAIhYAQACABAAAABkYXRhAgAAAAEA");
  silentAudio.loop = true;
  silentAudio.play();
  
  // iOS Safari 音频自动播放限制处理
  const video = document.getElementById('video-player');
  video.muted = false; // 关键：必须在用户交互时解除静音，否则视频可能无声
  
  // 声调保护：确保变速不变调
  video.preservesPitch = true; // 关键：确保变速不变调
  video.webkitPreservesPitch = true; // 兼容旧版 Safari
  video.mozPreservesPitch = true; // 兼容旧版 Firefox
  ```
- **实施时机**：
  - 用户点击"加入"按钮触发`handleUserJoin()`时
  - 与`video.play()` + `video.pause()`操作同时执行
  - 确保音频在后台持续循环播放
- **效果**：保持WebSocket连接在后台活跃，防止iOS Safari冻结连接，提升移动端体验稳定性

#### 4.1.5.6 开发调试工具 (Developer Backdoor)

**设计目的**：为开发者和线上故障排查提供便捷的调试工具，快速定位同步问题（网络问题 vs. 逻辑Bug）

**1. 调试面板激活方式**
- **隐藏入口**：连续点击页面标题5次，唤出调试浮层面板
- **快捷键支持**：`Ctrl+Shift+D`（PC）或三指长按（移动端）切换显示
- **环境检测**：仅在生产环境的`debug=true`参数或开发环境下可用
- **激活逻辑**：
  ```javascript
  let debugClickCount = 0;
  let lastDebugClickTime = 0;
  
  document.getElementById('page-title').addEventListener('click', () => {
    const now = Date.now();
    if (now - lastDebugClickTime > 1000) {
      debugClickCount = 0; // 1秒超时重置
    }
    debugClickCount++;
    lastDebugClickTime = now;
    
    if (debugClickCount >= 5) {
      toggleDebugPanel();
      debugClickCount = 0;
    }
  });
  ```

**2. 实时监控信息展示**
- **网络状态**：
  - WebSocket连接状态（连接中/已连接/断开/重连中）
  - 实时RTT数值及历史趋势图（最近30秒）
  - 消息收发统计（发送/接收/丢失）
- **同步状态**：
  - 当前计算出的TimeDiff（本地时间-服务器时间）
  - 四级同步策略当前状态（完美/温柔微调/快速追赶/强制对齐）
  - 最近一次同步操作详情
- **房间信息**：
  - 当前房间ID、用户数、房主身份
  - 播放器状态（播放/暂停/缓冲）
  - 当前播放进度和倍速
- **消息追踪**：
  - 最后一条收到的WebSocket指令（类型、内容、时间戳）
  - 最后一条发送的指令
  - 指令处理延迟（收到到执行的时间差）

**3. 故障排查辅助功能**
- **状态快照**：一键导出当前完整状态（JSON格式，含时间戳）
- **指令重放**：手动发送测试指令到当前房间
- **网络模拟**：
  - 延迟模拟：添加50-1000ms网络延迟
  - 丢包模拟：模拟1%-20%的丢包率
  - 断开模拟：手动断开/重连WebSocket
- **性能分析**：
  - 帧率监测（FPS）
  - 内存占用统计
  - CPU使用率估算

**4. 调试面板UI设计**
```html
<div id="debug-panel" class="debug-panel" style="display: none;">
  <div class="debug-header">
    <h3>🔧 调试面板</h3>
    <button onclick="toggleDebugPanel()">×</button>
  </div>
  <div class="debug-tabs">
    <button data-tab="network">网络状态</button>
    <button data-tab="sync">同步状态</button>
    <button data-tab="room">房间信息</button>
    <button data-tab="messages">消息追踪</button>
    <button data-tab="tools">调试工具</button>
  </div>
  <div class="debug-content">
    <!-- 各选项卡内容 -->
  </div>
</div>
```

**5. 线上故障排查流程**
1. **用户报告问题** → 引导用户连续点击标题5次
2. **打开调试面板** → 查看RTT、TimeDiff、连接状态
3. **分析网络问题**：RTT > 300ms 或频繁断开 → 网络问题
4. **分析逻辑问题**：RTT正常但TimeDiff异常 → 同步逻辑问题
5. **导出状态快照** → 发给开发团队分析
6. **模拟用户环境** → 使用网络模拟工具复现问题

**6. 安全与隐私考虑**
- **生产环境默认隐藏**：仅通过特定方式激活
- **无敏感信息**：调试信息不包含用户隐私数据
- **自动关闭**：面板30分钟无操作自动隐藏
- **访问日志**：记录调试面板的开启和关闭时间

### 4.2 开发里程碑

**第1周**：
- [x] 项目初始化
- [x] 数据库设计
- [x] API框架搭建
- [x] 基础房间管理

**第2周**：
- [x] WebSocket基础连接
- [x] 用户会话管理
- [x] 房间加入/离开逻辑

**第3周**：
- [x] 播放同步逻辑
- [x] RTT补偿机制
- [x] 消息广播机制

**第4周**：
- [x] 功能测试和优化
- [x] 部署配置
- [x] 文档完善
- [x] MVP版本发布

### 4.3 技术债务管理

**MVP阶段允许的技术债务**：
1. 协议简化：使用JSON代替Protocol Buffers
2. 缓存简化：暂时不实现独立缓存服务（使用SQLite内置缓存）
3. 测试简化：单元测试优先，集成测试后续补充
4. 监控简化：基本日志记录，详细监控后续添加

**后续重构计划**：
1. 第5-6周：Protocol Buffers迁移
2. 第7-8周：缓存优化和性能调优
3. 第9-10周：完善测试覆盖
4. 第11-12周：监控和性能优化

### 4.6 测试与质量保证

**4.6.1 测试策略**

```go
// 测试配置
type TestConfig struct {
    UnitTest     bool `json:"unit_test"`
    IntegrationTest bool `json:"integration_test"`
    LoadTest     bool `json:"load_test"`
    CoverageThreshold float64 `json:"coverage_threshold"` // 80%
}
```

**4.6.2 测试层次**

**1. 单元测试**
- **测试范围**：单个函数、方法或模块
- **覆盖重点**：核心业务逻辑、数据访问层、工具函数
- **技术实现**：
  ```go
  // 单元测试示例
  func TestRoomService(t *testing.T) {
      db := setupTestDB()
      repo := NewRoomRepository(db)
      service := NewRoomService(repo)
      
      // 测试房间创建
      room := &model.Room{
          ID:   "TEST01",
          Name: "测试房间",
      }
      
      err := service.CreateRoom(room)
      assert.NoError(t, err)
      
      // 测试房间查找
      found, err := service.GetRoomByID("TEST01")
      assert.NoError(t, err)
      assert.Equal(t, "测试房间", found.Name)
  }
  ```

**2. 集成测试**
- **测试范围**：多个模块之间的协作
- **覆盖重点**：API接口、WebSocket连接、数据库交互
- **技术实现**：
  ```go
  // 集成测试示例
  func TestWebSocketIntegration(t *testing.T) {
      // 启动测试服务器
      ts := httptest.NewServer(nil)
      defer ts.Close()
      
      // 创建WebSocket连接
      wsURL := strings.Replace(ts.URL, "http", "ws", 1) + "/ws"
      conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
      assert.NoError(t, err)
      defer conn.Close()
      
      // 测试消息发送和接收
      message := `{"type":"sync","room_id":"TEST01","data":{"current_time":10.5}}`
      err = conn.WriteMessage(websocket.TextMessage, []byte(message))
      assert.NoError(t, err)
      
      // 验证响应
      _, response, err := conn.ReadMessage()
      assert.NoError(t, err)
      assert.Contains(t, string(response), "sync_response")
  }
  ```

**3. 性能测试**
- **测试范围**：系统在高负载下的性能表现
- **覆盖重点**：并发连接、消息处理、数据库性能
- **技术实现**：
  ```go
  // 性能测试示例
  func TestLoadPerformance(t *testing.T) {
      // 创建并发连接测试
      const concurrentUsers = 100
      const testDuration = 30 * time.Second
      
      var wg sync.WaitGroup
      start := time.Now()
      
      for i := 0; i < concurrentUsers; i++ {
          wg.Add(1)
          go func(userID int) {
              defer wg.Done()
              
              // 模拟用户行为
              simulateUserBehavior(userID, testDuration)
          }(i)
      }
      
      wg.Wait()
      duration := time.Since(start)
      
      // 性能断言
      avgResponseTime := duration / time.Duration(concurrentUsers)
      assert.Less(t, avgResponseTime, 100*time.Millisecond, "平均响应时间应小于100ms")
  }
  ```

**4.6.3 质量保证措施**

1. **代码审查**：
   - 所有代码变更必须经过代码审查
   - 重点审查：安全问题、性能瓶颈、代码风格

2. **静态代码分析**：
   - 使用 `golangci-lint` 进行静态代码检查
   - 检测：代码质量、安全漏洞、性能问题

3. **测试覆盖率要求**：
   - 核心业务逻辑：≥90%
   - 数据访问层：≥85%
   - API层：≥80%
   - 整体项目：≥80%

4. **CI/CD集成**：
   - 自动运行测试和代码检查
   - 测试失败不允许合并
   - 自动构建和部署

**4.6.4 性能测试指标**

| 指标名称 | 目标值 | 测试方式 |
| --- | --- | --- |
| 单房间支持用户数 | ≥20 | 并发连接测试 |
| 系统支持房间数 | ≥100 | 负载测试 |
| API响应时间 | <100ms | 性能测试 |
| WebSocket消息延迟 | <200ms | 实时性测试 |
| 数据库查询时间 | <50ms | 数据库性能测试 |

---

## 五、部署与扩展性

### 5.1 容器化部署

**Docker配置**：

```dockerfile
# 多阶段构建 Dockerfile
FROM golang:1.21-alpine AS builder

# 安装依赖
RUN apk add --no-cache git

# 设置工作目录
WORKDIR /app

# 复制依赖文件
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并构建
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/

# MVP阶段改用 Alpine 作为运行镜像，便于调试和权限控制
FROM alpine:latest
WORKDIR /app
# 安装基础依赖（CA证书用于HTTPS请求，Timezone数据）
RUN apk add --no-cache ca-certificates tzdata

# 复制二进制文件
COPY --from=builder /app/main .

# 暴露端口
EXPOSE 8080 8090

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# 启动应用 - 使用绝对路径，确保信号正确传递
CMD ["/app/main"]
```

**开发环境 docker-compose.yml**：

```yaml
version: '3.8'

services:
  app:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8080:8080"  # HTTP API
      - "8090:8090"  # WebSocket
    volumes:
      - ./backend:/app
      - ./data:/data
    environment:
      - GIN_MODE=debug
      - DB_DRIVER=sqlite
      - DB_SOURCE=/data/xiaowo.db
    command: go run cmd/server/main.go

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/dev.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    restart: unless-stopped
```

**Nginx WebSocket配置要点**：
- **超时设置**：Nginx默认`proxy_read_timeout`为60秒，为防止误杀WebSocket长连接，需显式配置WebSocket路径的超时时间
- **配置示例**（在nginx.conf的location块中添加）：
  ```nginx
  location /ws {
      proxy_pass http://app:8090;
      proxy_http_version 1.1;
      proxy_set_header Upgrade $http_upgrade;
      proxy_set_header Connection "upgrade";
      proxy_set_header Host $host;
      
      # 关键：设置长连接超时为300秒（5分钟）
      proxy_read_timeout 300s;
      proxy_connect_timeout 75s;
  }
  ```
- **心跳协调**：WebSocket心跳间隔30s与Nginx超时300s协调，确保连接稳定性

### 5.2 生产环境部署 (单机 Docker Compose)

**部署架构**：
- **单机部署**：一台 2C4G 云服务器即可满足初期需求
- **反向代理**：Nginx 处理 SSL 终结和静态资源
- **服务编排**：Docker Compose 管理 App 服务
- **数据存储**：SQLite 数据库文件挂载到宿主机，确保数据持久化（虽然是阅后即焚，但重启不应丢失活跃房间）

**生产环境 docker-compose.prod.yml**：

```yaml
version: '3.8'

services:
  app:
    image: xiaowo/app:latest
    restart: always
    environment:
      - GIN_MODE=release
      - DB_DRIVER=sqlite
      - DB_SOURCE=/data/xiaowo.db
    volumes:
      - ./data:/data
    networks:
      - xiaowo-net

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/prod.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/nginx/ssl
    depends_on:
      - app
    networks:
      - xiaowo-net

networks:
  xiaowo-net:
    driver: bridge
```

### 5.3 日志和监控（极简版）

**MVP阶段日志策略**：
- **应用日志**：直接输出到 Standard Output (stdout/stderr)
- **查看方式**：生产环境通过 `docker logs xiaowo-app` 查看
- **日志格式**：JSON格式化输出，便于后续分析
- **日志轮转**：Docker自动处理日志文件轮转

```yaml
# docker-compose.prod.yml 中的日志配置
services:
  app:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### 5.3.2 埋点体系设计

**1. 核心埋点指标**

| 指标名称 | 统计维度 | 说明 | 采集方式 |
| --- | --- | --- | --- |
| **房间创建数** | 日/小时 | 每日/每小时创建的房间数量 | 后端API调用日志 |
| **有效观影时长** | 用户/房间 | 用户在房间内实际观看的时长（北极星指标） | 前端埋点+后端统计 |
| **链接解析成功率** | 总成功率/按来源 | 视频链接成功解析播放的比例 | 后端API调用日志 |
| **异常中断率** | 总中断率/按原因 | 播放过程中异常中断的比例 | 前端埋点+后端日志 |
| **热门资源Top 10** | 播放次数/观看时长 | 最受欢迎的10个视频资源 | 后端数据库统计 |

**2. 前端埋点**

**关键事件埋点**：
- `room_enter`：用户进入房间
- `room_leave`：用户离开房间
- `video_load`：视频加载成功
- `video_play`：视频开始播放
- `video_pause`：视频暂停
- `video_complete`：视频播放完成
- `sync_action`：同步操作执行
- `error_occur`：发生错误

**埋点数据结构**：
```json
{
  "event_name": "video_play",
  "room_id": "8XA92B",
  "session_id": "session-123",
  "timestamp": 1640995200000,
  "duration": 0,
  "url": "https://example.com/video.mp4",
  "user_agent": "Mozilla/5.0...",
  "device_type": "mobile",
  "referrer": "https://example.com"
}
```

**3. 后端埋点**

**API调用日志**：
- 所有API调用均记录详细日志
- 包含请求参数、响应状态、处理时间

**WebSocket消息日志**：
- 关键WebSocket消息类型记录
- 包含房间ID、会话ID、消息内容

**4. 数据存储与分析**

**存储方式**：
- 实时日志：直接写入stdout，Docker收集
- 统计数据：定期汇总到SQLite数据库

**分析方式**：
- 实时监控：通过Docker日志查看实时情况
- 离线分析：定期导出数据进行离线分析
- 报表生成：生成每日/每周报表

**5. 技术实现**

**前端**：
```javascript
// 埋点工具函数
function trackEvent(eventName, data) {
  const event = {
    event_name: eventName,
    room_id: currentRoomId,
    session_id: currentSessionId,
    timestamp: Date.now(),
    ...data
  };
  
  // 发送到后端
  fetch('/api/v1/track', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(event)
  }).catch(err => {
    console.error('埋点发送失败:', err);
  });
}

// 使用示例
trackEvent('video_play', {
  url: currentVideoUrl,
  duration: video.duration
});
```

**后端**：
```go
// 埋点API处理
func TrackEvent(c *gin.Context) {
  var event TrackEventRequest
  if err := c.ShouldBindJSON(&event); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
    return
  }
  
  // 记录到日志
  log.Printf("[TrackEvent] %s %s %s", event.EventName, event.RoomID, event.SessionID)
  
  // 定期汇总到数据库
  go trackService.RecordEvent(event)
  
  c.JSON(http.StatusOK, gin.H{"code": 0, "message": "success"})
}
```

### 5.3.1 SQLite 数据备份策略

**重要提醒**：虽然是阅后即焚的观影应用，但万一服务器损坏或误删 data 目录，正在看电影的用户会掉线。

**简单备份方案**：
在 docker-compose.prod.yml 中添加定时备份机制，确保数据安全。

```yaml
# 在 docker-compose.prod.yml 中添加备份服务
services:
  # ... 其他服务 ...
  
  backup:
    image: alpine:latest
    container_name: xiaowo-backup
    restart: unless-stopped
    volumes:
      - ./data:/data
      - ./backup:/data/backup
    command: |
      sh -c "
        apk add --no-cache sqlite
        while true; do
          mkdir -p /data/backup
          sqlite3 /data/xiaowo.db '.backup /data/backup/xiaowo_$$(date +%H).db'
          find /data/backup -name 'xiaowo_*.db' -mtime +1 -delete
          sleep 3600
        done
      "
    depends_on:
      - app
```

**备份策略说明**：
- **备份频率**：每小时执行一次
- **备份保留**：保留24小时内的备份文件（自动清理过期文件）
- **备份路径**：宿主机 `./backup` 目录
- **命令说明**：
  - `sqlite3 /data/xiaowo.db ".backup '/data/backup/xiaowo_$(date +%H).db'"` - SQLite 在线热备命令
  - `-mtime +1` - 清理超过1天的备份文件

**恢复方法**：
```bash
# 停止应用
docker-compose down

# 恢复最新的备份文件
cp ./backup/xiaowo_latest.db ./data/xiaowo.db

# 重启应用
docker-compose up -d
```

这是一个**5分钟就能配好的"保险"**，简单但有效！

### 5.4 扩展性规划

**水平扩展策略**：

1. **无状态服务设计**
   - 用户认证服务
   - 房间管理服务
   - API网关

2. **有状态服务优化**
   - WebSocket连接管理
   - 房间状态同步

3. **数据层扩展**
   - SQLite → PostgreSQL迁移
   - 读写分离架构
   - 分库分表策略

### 5.5 非功能需求

#### 5.5.1 性能指标要求

**1. 响应时间要求**
- **进房速度**：< 1秒（从点击链接到看到视频画面）
- **操作延迟**：< 100ms（用户操作到视觉反馈）
- **同步性能**：端到端指令传输延迟 < 200ms
- **播放器加载时间**：播放器初始化与资源准备时间 < 500ms（视频起播时间取决于源站响应速度）

**2. 并发能力要求**
- **单房间容量**支持：至少20个并发连接
- **系统容量**：支持100个房间并发运行
- **WebSocket连接**：支持2000个同时在线连接
- **API吞吐量**：QPS（每秒查询率）≥ 100

**3. 资源使用要求**
- **内存占用**：单个实例内存使用 < 256MB
- **CPU使用**：正常负载下CPU使用率 < 30%
- **网络带宽**：支持100Mbps网络吞吐

#### 5.5.2 安全性要求

**1. 房间安全**
- **房间号防爆破**：限制单位时间内房间号尝试次数
- **访问频率限制**：IP级别和用户级别的请求频率限制
- **黑名单机制**：支持域名黑名单，自动拦截恶意资源

**2. 数据安全**
- **日志脱敏**：用户敏感信息在日志中自动脱敏
- **数据加密**：敏感数据在传输和存储时加密
- **访问控制**：严格的API访问权限控制

**3. 内容安全**
- **资源验证**：所有视频链接进行安全性验证
- **恶意内容过滤**：检测并拦截恶意视频内容
- **版权保护**：尊重内容版权，提供版权提示

#### 5.5.3 兼容性要求

**1. 浏览器兼容性**
- **桌面浏览器**：
  - Chrome 90+ ✅
  - Safari 14+ ✅
  - Firefox 88+ ✅
  - Edge 90+ ✅
- **移动浏览器**：
  - iOS Safari ✅
  - Android Chrome ✅
  - 微信内置浏览器 ✅

**2. 视频格式兼容性**
- **容器格式**：MP4, WebM, M3U8 (HLS)
- **视频编码**：H.264, VP9, AV1
- **音频编码**：AAC, MP3, Opus
- **分辨率支持**：480p, 720p, 1080p, 2K, 4K

**3. 网络环境兼容性**
- **网络协议**：HTTP/1.1, HTTP/2, WebSocket
- **代理支持**：支持正向代理和反向代理
- **CDN兼容**：支持主流CDN服务商

#### 5.5.4 可用性与可靠性

**1. 可用性要求**
- **服务可用性**：99.9% 正常运行时间
- **故障恢复**：单点故障恢复时间 < 5分钟
- **数据持久性**：数据备份恢复成功率 > 99.99%

**2. 容错能力**
- **网络抖动**：自动重连机制，支持网络中断恢复
- **服务降级**：核心功能降级策略，保证基本可用
- **过载保护**：流量控制和限流机制

**3. 监控告警**

**3.1 增强型健康检查**
- **多维度检查**：数据库连接、WebSocket状态、内存使用、磁盘空间、活跃房间数
- **状态分级**：健康(healthy)、降级(degraded)、不健康(unhealthy)
- **详细报告**：包含每个检查项的状态、持续时间和消息
- **HTTP端点**：`/api/v1/health` 返回完整健康状态

**3.2 性能指标监控**
- **请求指标**：请求总数、响应时间分布、状态码统计
- **WebSocket指标**：活跃连接数、消息速率、连接状态
- **业务指标**：活跃房间数、在线用户数、同步精度
- **数据库指标**：查询次数、慢查询、连接池状态
- **Prometheus集成**：支持Prometheus抓取指标

**3.3 异常告警**
- **告警配置**：
  - 高CPU使用率 (≥80%)
  - 高内存使用率 (≥85%)
  - 连接数过高 (≥1000)
  - 健康检查失败
  - 业务异常 (活跃房间数骤降)
- **告警方式**：Docker日志聚合、预留告警集成接口

**3.4 技术实现**

```go
// 增强的健康检查
type HealthChecker struct {
    db           *gorm.DB
    wsHub        *WebSocketHub
    config       *Config
}

type HealthStatus struct {
    Status    string            `json:"status"`    // "healthy", "unhealthy", "degraded"
    Timestamp time.Time         `json:"timestamp"`
    Checks    map[string]Check  `json:"checks"`
}

type Check struct {
    Status    string        `json:"status"`    // "pass", "fail", "warn"
    Duration  time.Duration `json:"duration"`
    Message   string        `json:"message"`
}

func (hc *HealthChecker) CheckAll() *HealthStatus {
    checks := make(map[string]Check)
    
    // 1. 数据库连接检查
    checks["database"] = hc.checkDatabase()
    
    // 2. WebSocket连接检查
    checks["websocket"] = hc.checkWebSocketHub()
    
    // 3. 内存使用检查
    checks["memory"] = hc.checkMemoryUsage()
    
    // 4. 磁盘空间检查
    checks["disk"] = hc.checkDiskSpace()
    
    // 5. 活跃房间数检查
    checks["rooms"] = hc.checkActiveRooms()
    
    // 计算整体状态
    overallStatus := hc.calculateOverallStatus(checks)
    
    return &HealthStatus{
        Status:    overallStatus,
        Timestamp: time.Now(),
        Checks:    checks,
    }
}

// 性能指标收集器
type MetricsCollector struct {
    requestCount    *prometheus.CounterVec
    requestDuration *prometheus.HistogramVec
    activeRooms     prometheus.Gauge
    activeConnections prometheus.Gauge
    dbConnections   prometheus.Gauge
}

func NewMetricsCollector() *MetricsCollector {
    return &MetricsCollector{
        requestCount: prometheus.NewCounterVec(
            prometheus.CounterOpts{
                Name: "xiaowo_requests_total",
                Help: "Total number of HTTP requests",
            },
            []string{"method", "endpoint", "status"},
        ),
        requestDuration: prometheus.NewHistogramVec(
            prometheus.HistogramOpts{
                Name: "xiaowo_request_duration_seconds",
                Help: "Duration of HTTP requests",
                Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1.0, 2.0, 5.0},
            },
            []string{"method", "endpoint"},
        ),
        activeRooms: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "xiaowo_active_rooms",
                Help: "Number of active rooms",
            },
        ),
        activeConnections: prometheus.NewGauge(
            prometheus.GaugeOpts{
                Name: "xiaowo_active_connections",
                Help: "Number of active WebSocket connections",
            },
        ),
    }
}

// 指标上报中间件
func MetricsMiddleware(metrics *MetricsCollector) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        // 记录请求指标
        metrics.requestCount.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
        
        metrics.requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration.Seconds())
    })
}

// 告警配置
type AlertConfig struct {
    HighCPUThreshold       float64 `json:"high_cpu_threshold"`       // 80%
    HighMemoryThreshold    float64 `json:"high_memory_threshold"`    // 85%
    HighConnectionThreshold int     `json:"high_connection_threshold"` // 1000
    LowRoomThreshold       int     `json:"low_room_threshold"`       // 0 (业务异常)
}
```

---

## 六、安全策略与最佳实践

### 6.1 认证与授权

**临时会话机制**：
- 无需注册登录
- 基于Token的会话管理
- 会话有效期控制

```go
// 会话创建
func CreateSession() (*UserSession, string, error) {
    session := &UserSession{
        ID:        generateUUID(),
        Nickname:  generateRandomNickname(),
        Avatar:    generateRandomAvatar(),
        CreatedAt: time.Now(),
    }
    
    token := generateToken(session.ID)
    return session, token, nil
}

// Token验证
func ValidateToken(token string) (*UserSession, error) {
    sessionID := parseToken(token)
    session, err := getSessionByID(sessionID)
    if err != nil {
        return nil, err
    }
    
    // 检查会话是否过期 (24小时)
    if time.Since(session.CreatedAt) > 24*time.Hour {
        return nil, errors.New("session expired")
    }
    
    return session, nil
}
```

### 6.2 数据安全

**输入验证**：
- 所有用户输入进行严格验证
- 防止SQL注入 (使用GORM参数化查询)
- XSS防护 (HTML编码)

**数据加密**：
- 敏感配置环境变量管理
- 数据库连接加密 (生产环境)

### 6.3 API安全

**6.3.1 多层级限流策略**

```go
// 多层级限流器
type MultiLayerRateLimiter struct {
    // IP级别限流
    ipLimiter *rate.Limiter
    
    // 用户级别限流
    userLimiter map[string]*rate.Limiter
    
    // 全局限流
    globalLimiter *rate.Limiter
    
    // 黑白名单
    blacklist   map[string]bool
    whitelist   map[string]bool
    
    mu sync.RWMutex
}

func NewMultiLayerRateLimiter() *MultiLayerRateLimiter {
    return &MultiLayerRateLimiter{
        ipLimiter:     rate.NewLimiter(rate.Every(time.Second/100), 200), // IP: 100 QPS, 突发200
        globalLimiter: rate.NewLimiter(rate.Every(time.Second/1000), 2000), // 全局: 1000 QPS
        userLimiter:   make(map[string]*rate.Limiter),
        blacklist:     make(map[string]bool),
        whitelist:     make(map[string]bool),
    }
}

func (r *MultiLayerRateLimiter) Allow(ip, sessionID string) (bool, string) {
    // 1. 白名单检查
    if r.whitelist[ip] {
        return true, ""
    }
    
    // 2. 黑名单检查
    if r.blacklist[ip] {
        return false, "IP已被封禁"
    }
    
    // 3. 全局限流
    if !r.globalLimiter.Allow() {
        return false, "服务器繁忙，请稍后重试"
    }
    
    // 4. IP级别限流
    if !r.ipLimiter.Allow() {
        return false, "请求过于频繁，请稍后重试"
    }
    
    // 5. 用户级别限流
    r.mu.Lock()
    userLimiter, exists := r.userLimiter[sessionID]
    if !exists {
        userLimiter = rate.NewLimiter(rate.Every(time.Second/50), 100) // 用户: 50 QPS
        r.userLimiter[sessionID] = userLimiter
    }
    r.mu.Unlock()
    
    if !userLimiter.Allow() {
        return false, "用户请求过于频繁"
    }
    
    return true, ""
}
```

**6.3.2 防DDoS保护**

```go
// 防DDoS保护
type DDoSProtection struct {
    connectionTracker map[string][]time.Time // IP -> 连接时间列表
    requestTracker    map[string][]RequestInfo // IP -> 请求信息
    thresholdConfig   *DDoSThreshold
    mu                sync.Mutex
}

type RequestInfo struct {
    Timestamp time.Time
    Endpoint  string
    UserAgent string
}

type DDoSThreshold struct {
    MaxConnectionsPerIP    int           // IP最大连接数
    ConnectionWindow       time.Duration // 连接统计窗口
    MaxRequestsPerIP       int           // IP最大请求数
    RequestWindow          time.Duration // 请求统计窗口
    SuspiciousUserAgents   []string      // 可疑User-Agent
}
```

**6.3.3 CORS配置**

```go
func CORSMiddleware() gin.HandlerFunc {
    config := cors.Config{
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }
    
    // 根据环境动态设置CORS规则
    if os.Getenv("GIN_MODE") == "debug" {
        config.AllowAllOrigins = true // 开发环境允许所有来源
    } else {
        config.AllowOrigins = []string{"https://yourdomain.com"} // 生产环境只允许指定域名
    }
    
    return cors.New(config)
}
```

### 6.4 WebSocket安全

**连接验证**：
- Token验证 (每个WebSocket连接)
- 房间权限检查
- 连接数限制

**消息过滤**：
- 消息大小限制 (1MB)
- 消息类型白名单
- 频率限制 (每连接每秒10条消息)

---

## 七、总结与建议

### 7.1 技术方案优势

**开发效率优势**：
1. **SQLite**：零运维成本，单文件数据库，开发调试极其方便
2. **单体架构**：快速开发，无需微服务复杂性
3. **Go语言**：高性能、并发处理能力强，部署简单
4. **WebSocket**：原生浏览器支持，实时同步体验优秀

**运维成本优势**：
1. **极简部署**：单Go二进制 + SQLite文件，Docker Compose一键部署
2. **资源占用低**：2C4G服务器即可满足初期需求
3. **故障排查简单**：Docker logs直接查看，无复杂依赖

**用户体验优势**：
1. **零门槛使用**：无需注册，生成趣味昵称即可使用
2. **实时同步**：RTT补偿机制确保精确同步
3. **稳定可靠**：乐观锁版本控制，避免竞态条件

### 7.2 风险评估与应对

**技术风险**：
1. **SQLite并发限制**：应对措施 - WAL模式 + 连接池配置
2. **WebSocket扩展性**：应对措施 - 后续可迁移至Redis集群
3. **同步精度问题**：应对措施 - RTT平滑处理 + 四级追帧策略

**运营风险**：
1. **用户增长超预期**：应对措施 - 架构预留扩展空间
2. **恶意用户攻击**：应对措施 - 限流 + 连接验证
3. **数据丢失风险**：应对措施 - SQLite自动备份机制

### 7.3 实施建议

**开发阶段**：
1. **优先级排序**：房间管理 → WebSocket同步 → 播放控制 → 聊天功能
2. **测试驱动**：每个核心功能都要有完整的测试用例
3. **性能监控**：定期检查WebSocket连接数和响应时间

**部署阶段**：
1. **灰度发布**：先在内网环境测试，再逐步上线
2. **监控告警**：设置基本的连接数异常告警
3. **备份策略**：定期备份SQLite数据库文件

**迭代优化**：
1. **用户反馈**：收集真实用户的使用反馈
2. **性能优化**：基于实际使用情况优化同步算法
3. **功能扩展**：按需添加高级功能

### 7.4 成功指标

**技术指标**：
- WebSocket连接成功率>99%
- 同步延迟<100ms
- API响应时间<200ms
- 系统可用性>99.5%

**业务指标**：
- 房间创建成功率>95%
- 用户留存率>50%
- 系统可用性>99.5%

**团队指标**：
- 代码质量评分>8.0
- 测试覆盖率>80%
- 技术债务控制在合理范围
- 团队技能持续提升

---


    }
}

// 指标上报中间件
func MetricsMiddleware(metrics *MetricsCollector) gin.HandlerFunc {
    return gin.HandlerFunc(func(c *gin.Context) {
        start := time.Now()
        
        c.Next()
        
        duration := time.Since(start)
        
        // 记录请求指标
        metrics.requestCount.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
            strconv.Itoa(c.Writer.Status()),
        ).Inc()
        
        metrics.requestDuration.WithLabelValues(
            c.Request.Method,
            c.FullPath(),
        ).Observe(duration.Seconds())
    })
}

// 告警配置
type AlertConfig struct {
    HighCPUThreshold       float64 `json:"high_cpu_threshold"`       // 80%
    HighMemoryThreshold    float64 `json:"high_memory_threshold"`    // 85%
    HighConnectionThreshold int     `json:"high_connection_threshold"` // 1000
    LowRoomThreshold       int     `json:"low_room_threshold"`       // 0 (业务异常)
}
```

### 8.5 安全防护增强

#### 8.5.1 请求限流优化

**问题分析**：当前限流策略实现细节不明确，缺少防DDoS机制。

**解决方案**：
```go
// 多层级限流器
type MultiLayerRateLimiter struct {
    // IP级别限流
    ipLimiter *rate.Limiter
    
    // 用户级别限流
    userLimiter map[string]*rate.Limiter
    
    // 全局限流
    globalLimiter *rate.Limiter
    
    // 黑白名单
    blacklist   map[string]bool
    whitelist   map[string]bool
    
    mu sync.RWMutex
}

func NewMultiLayerRateLimiter() *MultiLayerRateLimiter {
    return &MultiLayerRateLimiter{
        ipLimiter:     rate.NewLimiter(rate.Every(time.Second/100), 200), // IP: 100 QPS, 突发200
        globalLimiter: rate.NewLimiter(rate.Every(time.Second/1000), 2000), // 全局: 1000 QPS
        userLimiter:   make(map[string]*rate.Limiter),
        blacklist:     make(map[string]bool),
        whitelist:     make(map[string]bool),
    }
}

func (r *MultiLayerRateLimiter) Allow(ip, sessionID string) (bool, string) {
    // 1. 白名单检查
    if r.whitelist[ip] {
        return true, ""
    }
    
    // 2. 黑名单检查
    if r.blacklist[ip] {
        return false, "IP已被封禁"
    }
    
    // 3. 全局限流
    if !r.globalLimiter.Allow() {
        return false, "服务器繁忙，请稍后重试"
    }
    
    // 4. IP级别限流
    if !r.ipLimiter.Allow() {
        return false, "请求过于频繁，请稍后重试"
    }
    
    // 5. 用户级别限流
    r.mu.Lock()
    userLimiter, exists := r.userLimiter[sessionID]
    if !exists {
        userLimiter = rate.NewLimiter(rate.Every(time.Second/50), 100) // 用户: 50 QPS
        r.userLimiter[sessionID] = userLimiter
    }
    r.mu.Unlock()
    
    if !userLimiter.Allow() {
        return false, "用户请求过于频繁"
    }
    
    return true, ""
}

// 防DDoS保护
type DDoSProtection struct {
    connectionTracker map[string][]time.Time // IP -> 连接时间列表
    requestTracker    map[string][]RequestInfo // IP -> 请求信息
    thresholdConfig   *DDoSThreshold
    mu                sync.Mutex
}

type RequestInfo struct {
    Timestamp time.Time
    Endpoint  string
    UserAgent string
}

type DDoSThreshold struct {
    MaxConnectionsPerIP    int           // IP最大连接数
    ConnectionWindow       time.Duration // 连接统计窗口
    MaxRequestsPerIP       int           // IP最大请求数
    RequestWindow          time.Duration // 请求统计窗口
    SuspiciousUserAgents   []string      // 可疑User-Agent
}
```

### 8.6 测试和质量保证

#### 8.6.1 测试策略

**问题分析**：文档中没有提到测试策略和质量保证机制。

**解决方案**：
```go
// 测试配置
type TestConfig struct {
    UnitTest     bool `json:"unit_test"`
    IntegrationTest bool `json:"integration_test"`
    LoadTest     bool `json:"load_test"`
    CoverageThreshold float64 `json:"coverage_threshold"` // 80%
}

// 单元测试示例
func TestRoomService(t *testing.T) {
    db := setupTestDB()
    repo := NewRoomRepository(db)
    service := NewRoomService(repo)
    
    // 测试房间创建
    room := &model.Room{
        ID:   "TEST01",
        Name: "测试房间",
    }
    
    err := service.CreateRoom(room)
    assert.NoError(t, err)
    
    // 测试房间查找
    found, err := service.GetRoomByID("TEST01")
    assert.NoError(t, err)
    assert.Equal(t, "测试房间", found.Name)
}

// 集成测试示例
func TestWebSocketIntegration(t *testing.T) {
    // 启动测试服务器
    ts := httptest.NewServer(nil)
    defer ts.Close()
    
    // 创建WebSocket连接
    wsURL := strings.Replace(ts.URL, "http", "ws", 1) + "/ws"
    conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    assert.NoError(t, err)
    defer conn.Close()
    
    // 测试消息发送和接收
    message := `{"type":"sync","room_id":"TEST01","data":{"current_time":10.5}}`
    err = conn.WriteMessage(websocket.TextMessage, []byte(message))
    assert.NoError(t, err)
    
    // 验证响应
    _, response, err := conn.ReadMessage()
    assert.NoError(t, err)
    assert.Contains(t, string(response), "sync_response")
}

// 性能测试示例
func TestLoadPerformance(t *testing.T) {
    // 创建并发连接测试
    const concurrentUsers = 100
    const testDuration = 30 * time.Second
    
    var wg sync.WaitGroup
    start := time.Now()
    
    for i := 0; i < concurrentUsers; i++ {
        wg.Add(1)
        go func(userID int) {
            defer wg.Done()
            
            // 模拟用户行为
            simulateUserBehavior(userID, testDuration)
        }(i)
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    // 性能断言
    avgResponseTime := duration / time.Duration(concurrentUsers)
    assert.Less(t, avgResponseTime, 100*time.Millisecond, "平均响应时间应小于100ms")
}
```

---

**文档版本**：v1.2  
**更新时间**：2025-12-30  
**负责人**：后盾（系统架构师）  
**审核状态**：已审核  
**实施状态**：准备启动
