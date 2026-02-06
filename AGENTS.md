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
- `Zone`: Zone(二级域名), RecordCount, CreatedAt, UpdatedAt
- `Domain`: Zone, Domain, Name, IPs[], TTL, RecordCount, CreatedAt, UpdatedAt
- `Config`: App(host/port/env), Etcd(..., coredns_prefix), JWT(secret/expiry), Logger

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

# Zone 管理 (Admin)
POST /api/dns/zones/list
POST /api/dns/zones/get
POST /api/dns/zones/create
POST /api/dns/zones/update
POST /api/dns/zones/delete

# Domain 管理 (JWT)
POST /api/dns/domains/list
POST /api/dns/domains/get
POST /api/dns/domains/create
POST /api/dns/domains/update
POST /api/dns/domains/delete
```

## etcd Key 规范

- 用户: `/dance/users/{user-id}`
- Zone: `/dancer/zones/{zone}` (例: `example.com`)
- Domain: `/dancer/domains/{zone}/{domain}` (例: `www.example.com`)
- CoreDNS: `{prefix}/{反转zone}/{domain}/x{n}` (例: `/skydns/com/example/www/x1`)

## 错误类型

ErrUserNotFound, ErrUserExists, ErrInvalidCredentials, ErrWrongPassword, ErrRecordNotFound, ErrRecordExists, ErrZoneNotFound, ErrZoneExists, ErrDomainNotFound, ErrDomainExists, ErrInvalidToken, ErrTokenExpired, ErrUnauthorized, ErrForbidden, ErrInvalidInput, ErrEtcdUnavailable

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
