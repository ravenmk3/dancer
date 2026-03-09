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
│   │   │   ├── zone.go            # Zone CRUD 操作
│   │   │   └── domain.go          # Domain CRUD + CoreDNS 同步
│   │   └── key_prefix.go          # etcd key 前缀定义
│   ├── models/                     # 数据模型
│   │   ├── user.go                # 用户模型
│   │   ├── zone.go                # Zone 模型
│   │   ├── domain.go              # Domain 模型
│   │   └── dto.go                 # 请求/响应 DTO
│   ├── handlers/                   # HTTP 处理器
│   │   ├── base.go                # 基础响应结构
│   │   ├── user.go                # 用户管理处理器
│   │   ├── zone.go                # Zone 管理处理器
│   │   ├── domain.go              # Domain 管理处理器
│   │   └── health.go              # 健康检查处理器
│   ├── services/                   # 业务逻辑层
│   │   ├── user_service.go        # 用户业务逻辑
│   │   ├── zone_service.go        # Zone 业务逻辑
│   │   └── domain_service.go      # Domain 业务逻辑
│   └── router/                     # 路由配置
│       ├── router.go              # Echo 路由定义
│       └── logger.go              # 自定义访问日志中间件
├── assets/                         # 前端构建产物 (Go embed)
├── config.toml                     # 配置文件
go.mod
go.sum
Dockerfile
```

## 2. 模块划分

| 模块 | 职责 |
|------|------|
| **cmd/server** | 程序入口，负责初始化配置、日志、存储，启动 HTTP 服务 |
| **config** | 加载并解析 config.toml，提供配置访问接口 |
| **errors** | 定义业务错误类型 |
| **auth** | JWT Token 生成/验证、密码加密、认证中间件、RBAC |
| **logger** | 彩色控制台输出 + 轮转文件日志 (使用 logrus) |
| **storage/etcd** | 封装 etcd v3 客户端，提供用户、Zone、Domain 的 CRUD 操作及 CoreDNS 同步 |
| **models** | 定义用户、Zone、Domain 实体及请求/响应 DTO |
| **handlers** | HTTP 处理器，解析请求、调用服务层，返回响应 |
| **services** | 业务逻辑层，封装用户管理、Zone 管理和 Domain 管理的核心逻辑 |
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
        ReconnectInterval    int      `toml:"reconnect_interval"`     // 初始重连间隔(秒)
        MaxReconnectInterval int      `toml:"max_reconnect_interval"` // 最大重连间隔(秒)
        HealthCheckInterval  int      `toml:"health_check_interval"`  // 健康检查间隔(秒)
        DialTimeout          int      `toml:"dial_timeout"`           // 连接超时(秒)
        CorednsPrefix        string   `toml:"coredns_prefix"`         // CoreDNS etcd key 前缀
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
    ID        string   `json:"id"`
    Username  string   `json:"username"`
    Password  string   `json:"-"`          // 不序列化
    UserType  UserType `json:"user_type"`
    CreatedAt int64    `json:"created_at"`
    UpdatedAt int64    `json:"updated_at"`
}
```

### 3.3 Zone 模型 (internal/models/zone.go)

```go
type Zone struct {
    Zone        string `json:"zone"`          // 二级域名，如 example.com
    RecordCount int    `json:"record_count"`  // 该 zone 下的域名数量
    CreatedAt   int64  `json:"created_at"`    // 创建时间戳
    UpdatedAt   int64  `json:"updated_at"`    // 更新时间戳
}
```

### 3.4 Domain 模型 (internal/models/domain.go)

```go
type Domain struct {
    Zone        string   `json:"zone"`         // 所属 zone，如 example.com
    Domain      string   `json:"domain"`       // 子域名部分，如 www
    Name        string   `json:"name"`         // 完整域名，如 www.example.com
    IPs         []string `json:"ips"`          // IP 地址列表
    TTL         int      `json:"ttl"`          // TTL (秒)
    RecordCount int      `json:"record_count"` // IP 记录数量
    CreatedAt   int64    `json:"created_at"`   // 创建时间戳
    UpdatedAt   int64    `json:"updated_at"`   // 更新时间戳
}
```

## 4. API 路由设计

```
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

# Zone 管理 (Admin 权限)
POST   /api/dns/zones/list          # 列举所有 Zone
POST   /api/dns/zones/get           # 获取 Zone 详情
POST   /api/dns/zones/create        # 创建 Zone
POST   /api/dns/zones/update        # 更新 Zone（当前仅刷新 updated_at 时间戳）
POST   /api/dns/zones/delete        # 删除 Zone（级联删除）

# Domain 管理 (JWT 认证)
POST   /api/dns/domains/list        # 列举 Zone 下所有 Domain
POST   /api/dns/domains/get         # 获取 Domain 详情
POST   /api/dns/domains/create      # 创建 Domain
POST   /api/dns/domains/update      # 更新 Domain（IP 列表替换）
POST   /api/dns/domains/delete      # 删除 Domain（级联删除）
```

## 5. etcd Key 规划

| 数据类型 | Key 格式 | 示例 |
|---------|---------|------|
| 用户记录 | `/dancer/users/{user-id}` | `/dancer/users/1701234567890` |
| Zone | `/dancer/zones/{zone}` | `/dancer/zones/example.com` |
| Domain | `/dancer/domains/{zone}/{domain}` | `/dancer/domains/example.com/www` |
| CoreDNS | `{prefix}/{反转zone}/{domain}/x{n}` | `/skydns/com/example/www/x1` |

### 5.1 etcd 客户端自动重连

#### 连接状态管理

```
┌─────────────────────────────────────────────────────┐
│                 EtcdClientManager                    │
├─────────────────────────────────────────────────────┤
│  状态: disconnected / connecting / connected          │
│  后台 goroutine 自动重连                              │
│  健康检查定时器                                       │
└─────────────────────────────────────────────────────┘
```

#### 重连策略

- **首次连接**: 异步尝试连接，失败则后台重试
- **断线检测**: 每 30 秒健康检查一次
- **指数退避**: 5s → 10s → 20s → 30s (上限)
- **等待超时**: 存储操作默认等待 5 秒

#### 配置项

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `reconnect_interval` | 5 | 初始重连间隔(秒) |
| `max_reconnect_interval` | 30 | 最大重连间隔(秒) |
| `health_check_interval` | 30 | 健康检查间隔(秒) |
| `dial_timeout` | 5 | 连接超时(秒) |
| `coredns_prefix` | /skydns | CoreDNS etcd key 前缀 |

### 5.2 CoreDNS 同步机制

Domain 的增删改操作会自动同步到 CoreDNS 的 etcd key：

```go
// 同步流程:
1. Domain Create/Update/Delete 操作
2. 比较新旧 IP 列表差异
3. 删除多余的 CoreDNS 记录
4. 添加新增的 CoreDNS 记录（SkyDNS 格式）
5. 更新 Domain 元数据

// CoreDNS 记录格式 (SkyDNS):
{
  "host": "192.168.1.1",
  "ttl": 300
}
```

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
reconnect_interval = 5          # 初始重连间隔(秒)
max_reconnect_interval = 30     # 最大重连间隔(秒)
health_check_interval = 30      # 健康检查间隔(秒)
dial_timeout = 5               # 连接超时(秒)
coredns_prefix = "/skydns"     # CoreDNS etcd key 前缀

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

## 9. 业务错误处理规范 (internal/errors/errors.go)

### 9.1 BusinessError 接口

所有业务错误统一通过 `BusinessError` 接口管理：

```go
// BusinessError 业务错误接口
type BusinessError interface {
    error
    Code() string   // 业务错误码，如 "user_not_found"
    Status() int    // HTTP 状态码，如 404
}
```

实现结构：

```go
type businessError struct {
    code    string  // 错误码
    status  int     // HTTP 状态码
    message string  // 错误消息
}

func (e *businessError) Error() string { return e.message }
func (e *businessError) Code() string  { return e.code }
func (e *businessError) Status() int   { return e.status }

// Is 实现 error.Is 接口，按错误码比较
func (e *businessError) Is(target error) bool {
    t, ok := target.(*businessError)
    if !ok {
        return false
    }
    return e.code == t.code
}
```

### 9.2 创建业务错误

使用统一工厂函数：

```go
func NewBusinessError(code string, status int, message string) BusinessError

// 示例
var ErrUserNotFound = NewBusinessError("user_not_found", 404, "user not found")
```

### 9.3 预定义业务错误

```go
var (
    ErrUserNotFound       = NewBusinessError("user_not_found", 404, "user not found")
    ErrUserExists         = NewBusinessError("user_exists", 409, "user already exists")
    ErrInvalidCredentials = NewBusinessError("invalid_credentials", 401, "invalid username or password")
    ErrWrongPassword      = NewBusinessError("wrong_password", 400, "wrong password")
    ErrSamePassword       = NewBusinessError("same_password", 400, "new password cannot be the same as old password")
    ErrCannotDeleteDefaultAdmin = NewBusinessError("cannot_delete_default_admin", 403, "cannot delete default admin user")
    ErrRecordNotFound     = NewBusinessError("record_not_found", 404, "DNS record not found")
    ErrRecordExists       = NewBusinessError("record_exists", 409, "DNS record already exists")
    ErrZoneNotFound       = NewBusinessError("zone_not_found", 404, "zone not found")
    ErrZoneExists         = NewBusinessError("zone_exists", 409, "zone already exists")
    ErrDomainNotFound     = NewBusinessError("domain_not_found", 404, "domain not found")
    ErrDomainExists       = NewBusinessError("domain_exists", 409, "domain already exists")
    ErrInvalidToken       = NewBusinessError("invalid_token", 401, "invalid token")
    ErrTokenExpired       = NewBusinessError("token_expired", 401, "token expired")
    ErrUnauthorized       = NewBusinessError("unauthorized", 401, "unauthorized")
    ErrForbidden          = NewBusinessError("forbidden", 403, "forbidden")
    ErrInvalidInput       = NewBusinessError("invalid_input", 400, "invalid input")
    ErrPasswordTooLong    = NewBusinessError("invalid_input", 400, "password exceeds maximum length of 72 bytes")
    ErrEtcdUnavailable    = NewBusinessError("service_unavailable", 503, "etcd service temporarily unavailable, please retry later")
)
```

### 9.4 错误判断方式

```go
// 方式 1: errors.Is (判断错误类型)
if errors.Is(err, ErrUserNotFound) {
    // 处理用户不存在
}

// 方式 2: errors.As (获取错误详情)
var bizErr BusinessError
if errors.As(err, &bizErr) {
    status := bizErr.Status()  // 获取 HTTP 状态码
    code := bizErr.Code()      // 获取错误码
    message := bizErr.Error()  // 获取错误消息
}
```

### 9.5 全局错误处理器

在 `internal/router/router.go` 中注册自定义错误处理器：

```go
func New(...) *echo.Echo {
    e := echo.New()
    // 设置全局错误处理器
    e.HTTPErrorHandler = customHTTPErrorHandler
    // ... 路由配置
}

func customHTTPErrorHandler(err error, c echo.Context) {
    // 如果响应已经写入，直接返回
    if c.Response().Committed {
        return
    }

    // 尝试从 echo.HTTPError 中提取内部错误
    var httpErr *echo.HTTPError
    if errors.As(err, &httpErr) {
        if innerErr, ok := httpErr.Message.(error); ok {
            err = innerErr
        } else {
            // Message 不是 error 类型，返回 HTTP 错误
            c.JSON(httpErr.Code, Response{
                Code:    "http_error",
                Message: httpErr.Error(),
            })
            return
        }
    }

    // 检查是否是业务错误，自动映射 HTTP 状态码
    var bizErr BusinessError
    if errors.As(err, &bizErr) {
        c.JSON(bizErr.Status(), Response{
            Code:    bizErr.Code(),
            Message: bizErr.Error(),
        })
        return
    }

    // 未知错误
    logger.Log.WithError(err).Error("Unhandled error")
    c.JSON(http.StatusInternalServerError, Response{
        Code:    "internal_error",
        Message: "internal server error",
    })
}
```

### 9.6 开发规范

1. **所有业务错误必须使用 `BusinessError` 接口**，不要直接使用 `errors.New()`
2. **新增错误码时**：
   - 在 `internal/errors/errors.go` 中定义
   - 选择合适的 HTTP 状态码
   - 使用小写下划线的英文错误码（如 `user_not_found`）
3. **不要在 handler 层手动设置 HTTP 状态码**，由全局错误处理器自动映射
4. **错误响应格式**：`{"code": "xxx_error", "message": "..."}`

## 10. 三层架构

```
HTTP Request
    ↓
Handler (解析请求, 调用服务层)
    ↓
Service (业务逻辑)
    ↓
Storage (etcd CRUD + CoreDNS 同步)
```

## 11. 关键实现文件

1. `internal/storage/etcd/client.go` - etcd 客户端封装（支持自动重连）
2. `internal/storage/etcd/domain.go` - Domain CRUD + CoreDNS 同步
3. `internal/auth/middleware.go` - JWT 认证中间件
4. `internal/models/dto.go` - 请求/响应数据结构
5. `internal/services/domain_service.go` - Domain 业务逻辑（IP 列表对比同步）
6. `internal/errors/errors.go` - 错误定义
7. `internal/handlers/health.go` - 健康检查处理器
8. `internal/router/logger.go` - 自定义访问日志中间件
9. `cmd/server/main.go` - 程序入口

## 12. 错误响应格式

```json
{
  "code": "xxx_error",
  "message": "<错误信息>"
}
```
