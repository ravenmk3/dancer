# Dancer DNS 管理工具 - 后端设计文档

**Go 版本**: 1.25

## 技术栈

| 类别 | 技术选型 |
|------|----------|
| Web 框架 | Echo |
| 数据存储 | etcd v3 |
| 日志库 | logrus + lumberjack |
| JWT | golang-jwt/jwt/v5 |
| 密码加密 | bcrypt |
| 配置格式 | TOML |

## 1. 项目目录结构

```
dancer/
├── cmd/
│   └── server/
│       └── main.go                 # 程序入口
├── internal/
│   ├── config/                     # 配置模块
│   │   ├── config.go              # 配置加载逻辑
│   │   └── model.go               # 配置数据结构
│   ├── errors/                     # 错误定义
│   │   └── errors.go              # 业务错误类型
│   ├── auth/                       # 认证授权模块
│   │   ├── jwt.go                 # JWT token 生成与验证
│   │   ├── middleware.go          # Echo 中间件
│   │   └── password.go            # 密码加密/验证
│   ├── logger/                     # 日志模块
│   │   ├── logger.go              # 日志器初始化
│   │   └── rotate.go              # 文件轮转逻辑
│   ├── storage/                    # 存储模块
│   │   ├── etcd/
│   │   │   ├── client.go          # etcd 客户端封装
│   │   │   ├── user.go            # 用户 CRUD 操作
│   │   │   └── dns.go             # DNS 记录 CRUD 操作
│   │   └── key_prefix.go          # etcd key 前缀定义
│   ├── models/                     # 数据模型
│   │   ├── user.go                # 用户模型
│   │   ├── dns.go                 # DNS 记录模型
│   │   └── dto.go                 # 请求/响应 DTO
│   ├── handlers/                   # HTTP 处理器
│   │   ├── base.go                # 基础响应结构
│   │   ├── user.go                # 用户管理处理器
│   │   ├── dns.go                 # DNS 记录处理器
│   │   └── health.go              # 健康检查处理器
│   ├── services/                   # 业务逻辑层
│   │   ├── user_service.go        # 用户业务逻辑
│   │   └── dns_service.go         # DNS 业务逻辑
│   └── router/                     # 路由配置
│       ├── router.go              # Echo 路由定义
│       └── logger.go              # 自定义访问日志中间件
├── assets/                         # 前端构建产物 (Go embed)
├── config.toml                     # 配置文件
├── go.mod
├── go.sum
└── Dockerfile
```

## 2. 模块划分

| 模块 | 职责 |
|------|------|
| **cmd/server** | 程序入口，负责初始化配置、日志、存储，启动 HTTP 服务 |
| **config** | 加载并解析 config.toml，提供配置访问接口 |
| **errors** | 定义业务错误类型 |
| **auth** | JWT Token 生成/验证、密码加密、认证中间件、RBAC |
| **logger** | 彩色控制台输出 + 轮转文件日志 (使用 logrus) |
| **storage/etcd** | 封装 etcd v3 客户端，提供用户和 DNS 记录的 CRUD 操作 |
| **models** | 定义用户、DNS 记录实体及请求/响应 DTO |
| **handlers** | HTTP 处理器，解析请求、调用服务层，返回响应 |
| **services** | 业务逻辑层，封装用户管理和 DNS 管理的核心逻辑 |
| **router** | 定义 Echo 路由组和中间件注册 |

## 3. 数据结构定义

### 3.1 配置模型 (internal/config/model.go)

```go
type Config struct {
    App struct {
        Host string `toml:"host"`
        Port int    `toml:"port"`
        Env  string `toml:"env"`
    } `toml:"app"`

    Etcd struct {
        Endpoints            []string `toml:"endpoints"`
        Username             string   `toml:"username"`
        Password             string   `toml:"password"`
        ReconnectInterval    int      `toml:"reconnect_interval"`
        MaxReconnectInterval int      `toml:"max_reconnect_interval"`
        HealthCheckInterval  int      `toml:"health_check_interval"`
        DialTimeout          int      `toml:"dial_timeout"`
    } `toml:"etcd"`

    JWT struct {
        Secret string `toml:"secret"`
        Expiry int64  `toml:"expiry"`
    } `toml:"jwt"`

    Logger struct {
        Level     string `toml:"level"`
        FilePath  string `toml:"file_path"`
        MaxSize   int    `toml:"max_size"`
        MaxBackup int    `toml:"max_backup"`
        MaxAge    int    `toml:"max_age"`
    } `toml:"logger"`
}
```

### 3.2 用户模型 (internal/models/user.go)

```go
type UserType string

const (
    UserTypeAdmin  UserType = "admin"
    UserTypeNormal UserType = "normal"
)

type User struct {
    ID        string    `json:"id"`
    Username  string    `json:"username"`
    Password  string    `json:"-"`          // 不序列化
    UserType  UserType  `json:"user_type"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 3.3 DNS 记录模型 (internal/models/dns.go)

```go
type DNSRecord struct {
    ID        string    `json:"id"`
    Domain    string    `json:"domain"`
    IP        string    `json:"ip"`
    TTL       int64     `json:"ttl"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

## 4. API 路由设计

GET/POST /api/health                # 健康检查

POST   /api/auth/login              # 用户登录
POST   /api/auth/refresh            # 刷新 Token

# 当前用户 (JWT 认证)
POST   /api/me                      # 获取当前登录用户信息
POST   /api/me/change-password      # 修改当前用户密码

# 用户管理 (Admin 权限)
POST   /api/user/list               # 列举用户
POST   /api/user/create             # 创建用户
POST   /api/user/update             # 更新用户
POST   /api/user/delete             # 删除用户

# DNS 记录管理 (JWT 认证)
POST   /api/dns/records/list        # 列举记录
POST   /api/dns/records/create      # 创建记录
POST   /api/dns/records/update      # 更新记录
POST   /api/dns/records/delete      # 删除记录
```

## 5. etcd Key 规划

| 数据类型 | Key 格式 | 示例 |
|---------|---------|------|
| 用户记录 | `/dance/users/{user-id}` | `/dance/users/1701234567890` |
| DNS 记录 | `/coredns/{域名层级}` | `/coredns/com/github` |

### DNS 记录 Key 说明

CoreDNS 使用按域名层级组织的 key 格式：

```
/coredns/{域名反转}/{子域名}/{记录名}
```

示例：
- `github.com` → `/coredns/com/github/`
- `api.github.com` → `/coredns/com/github/api/`
- 多个 A 记录: `/coredns/com/github/x1`, `/coredns/com/github/x2` ...

## 5.1 etcd 客户端自动重连

### 连接状态管理

```
┌─────────────────────────────────────────────────────┐
│                 EtcdClientManager                    │
├─────────────────────────────────────────────────────┤
│  状态: disconnected / connecting / connected          │
│  后台 goroutine 自动重连                              │
│  健康检查定时器                                       │
└─────────────────────────────────────────────────────┘
```

### 重连策略

- **首次连接**: 异步尝试连接，失败则后台重试
- **断线检测**: 每 30 秒健康检查一次
- **指数退避**: 5s → 10s → 20s → 30s (上限)
- **等待超时**: 存储操作默认等待 5 秒

### 配置项

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `reconnect_interval` | 5 | 初始重连间隔(秒) |
| `max_reconnect_interval` | 30 | 最大重连间隔(秒) |
| `health_check_interval` | 30 | 健康检查间隔(秒) |
| `dial_timeout` | 5 | 连接超时(秒) |

## 6. 认证授权

- JWT (HS256 算法)
- 从 Header 获取: `Authorization: Bearer <token>`
- 管理员权限检查中间件: `RequireAdmin()`

## 7. 日志系统

- 库: logrus + lumberjack
- 控制台: 彩色输出 (开发环境)
- 文件: 支持轮转 (max_size, max_backup, max_age)
- 访问日志: 自定义中间件 (DEBUG 级别)
  - 格式: `DEBU[2026-02-03 23:26:42] 127.0.0.1 | GET /api/health | 200 | 0ms | 0B/43B`

## 8. 配置文件 (config.toml)

```toml
[app]
host = "0.0.0.0"
port = 8080
env = "development"

[etcd]
endpoints = ["http://localhost:2379"]
# username = ""
# password = ""
# reconnect_interval = 5          # 初始重连间隔(秒)
# max_reconnect_interval = 30     # 最大重连间隔(秒)
# health_check_interval = 30      # 健康检查间隔(秒)
# dial_timeout = 5               # 连接超时(秒)

[jwt]
secret = "your-256-bit-secret"
expiry = 86400

[logger]
level = "info"
file_path = "logs/dancer.log"
max_size = 100
max_backup = 7
max_age = 7
```

## 9. 错误定义 (internal/errors/errors.go)

```go
var (
    ErrUserNotFound        = errors.New("user not found")
    ErrUserExists          = errors.New("user already exists")
    ErrInvalidCredentials  = errors.New("invalid username or password")
    ErrWrongPassword       = errors.New("wrong password")
    ErrRecordNotFound      = errors.New("DNS record not found")
    ErrRecordExists        = errors.New("DNS record already exists")
    ErrInvalidToken        = errors.New("invalid token")
    ErrTokenExpired        = errors.New("token expired")
    ErrUnauthorized        = errors.New("unauthorized")
    ErrForbidden           = errors.New("forbidden")
    ErrInvalidInput        = errors.New("invalid input")
    ErrEtcdUnavailable     = errors.New("etcd service temporarily unavailable")
)
```

## 10. 三层架构

```
HTTP Request
    ↓
Handler (解析请求, 调用服务层)
    ↓
Service (业务逻辑)
    ↓
Storage (etcd CRUD)
```

## 11. 关键实现文件

1. `internal/storage/etcd/client.go` - etcd 客户端封装（支持自动重连）
2. `internal/auth/middleware.go` - JWT 认证中间件
3. `internal/models/dto.go` - 请求/响应数据结构
4. `internal/services/user_service.go` - 用户业务逻辑
5. `internal/errors/errors.go` - 错误定义
6. `internal/handlers/health.go` - 健康检查处理器
7. `internal/router/logger.go` - 自定义访问日志中间件
8. `cmd/server/main.go` - 程序入口

## 12. 错误响应格式

```json
{
  "code": "xxx_error",
  "message": "<错误信息>"
}
```
