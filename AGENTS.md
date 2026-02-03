# AGENTS.md

## 项目概述

Dancer - DNS 管理工具

## 技术栈

- **Web 框架**: Echo v4
- **数据存储**: etcd v3
- **认证**: JWT (HS256)
- **密码加密**: bcrypt
- **日志**: logrus + lumberjack
- **配置**: TOML

## 目录结构

```
cmd/server/main.go              # 程序入口
internal/
  config/                       # TOML 配置加载
  errors/                       # 业务错误定义
  auth/                         # JWT、密码加密、中间件
  logger/                       # 日志初始化与轮转
  storage/etcd/                 # etcd 客户端及 CRUD 操作
  models/                       # 实体模型与 DTO
  handlers/                     # HTTP 处理器
  services/                     # 业务逻辑层
  router/                       # 路由配置
```

## 架构分层

HTTP 请求 → Handler → Service → Storage(etcd)

## 关键模型

- `User`: ID, Username, Password, UserType(admin/normal), CreatedAt, UpdatedAt
- `DNSRecord`: ID, Domain, IP, TTL, CreatedAt, UpdatedAt
- `Config`: App(host/port/env), Etcd(endpoints/credentials/reconnect_interval/max_reconnect_interval/health_check_interval/dial_timeout), JWT(secret/expiry), Logger

## API 路由

GET/POST /api/health            # 健康检查

POST /api/auth/login
POST /api/auth/refresh

POST /api/me                    # 获取当前用户
POST /api/me/change-password    # 修改密码

POST /api/user/list             # Admin 权限
POST /api/user/create
POST /api/user/update
POST /api/user/delete

POST /api/dns/records/list
POST /api/dns/records/create
POST /api/dns/records/update
POST /api/dns/records/delete
```

## etcd Key 规范

- 用户: `/dance/users/{user-id}`
- DNS: `/coredns/{反转域名}/{子域名}` (例: `github.com` → `/coredns/com/github`)

## 错误类型

ErrUserNotFound, ErrUserExists, ErrInvalidCredentials, ErrWrongPassword, ErrRecordNotFound, ErrRecordExists, ErrInvalidToken, ErrTokenExpired, ErrUnauthorized, ErrForbidden, ErrInvalidInput, ErrEtcdUnavailable

## 开发规范

1. 所有 API 使用 POST 方法
2. JWT 从 Header 获取: `Authorization: Bearer <token>`
3. 管理员权限检查中间件: `RequireAdmin()`
4. 错误响应格式: `{"code": "xxx_error", "message": "..."}`
5. 时间戳使用 int64(Unix timestamp)

## 关键文件

- `internal/auth/middleware.go` - JWT 认证中间件
- `internal/storage/etcd/client.go` - etcd 客户端封装（支持自动重连）
- `internal/handlers/health.go` - 健康检查处理器
- `internal/router/logger.go` - 自定义访问日志中间件
- `internal/services/user_service.go` - 用户业务逻辑
- `cmd/server/main.go` - 程序入口
