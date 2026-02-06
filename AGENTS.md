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
docs/                           # 项目文档
  backend-api.md                # 完整 API 文档
  backend-design.md             # 后端架构设计文档
  testcase/                     # 测试用例文档
cmd/server/main.go              # 程序入口
internal/
  config/                       # TOML 配置加载
  errors/                       # 业务错误定义
  auth/                         # JWT、密码加密、中间件
  logger/                       # 日志初始化与轮转
  storage/                      # 存储相关
  models/                       # 实体模型与 DTO
  handlers/                     # HTTP 处理器
  services/                     # 业务逻辑层
  router/                       # 路由配置
```

## 架构分层

HTTP 请求 → Handler → Service → Storage(etcd)

## 关键模型

- **User**: 用户信息，包含账号类型（管理员/普通用户）
- **Zone**: 域名区域（二级域名，如 example.com）
- **Domain**: 域名记录，包含完整域名、IP 列表和 TTL
- **Config**: 应用配置，包括服务端口、etcd 连接、JWT 密钥、日志等

## API 路由

- GET/POST /api/health          # 健康检查
- POST /api/auth/login
- POST /api/auth/refresh
- POST /api/me                  # 获取当前用户
- POST /api/me/change-password  # 修改密码

### 用户管理 (Admin)

- POST /api/user/list
- POST /api/user/create
- POST /api/user/update
- POST /api/user/delete

### Zone 管理 (Admin)

- POST /api/dns/zones/list
- POST /api/dns/zones/get
- POST /api/dns/zones/create
- POST /api/dns/zones/update
- POST /api/dns/zones/delete

### Domain 管理 (JWT)

- POST /api/dns/domains/list
- POST /api/dns/domains/get
- POST /api/dns/domains/create
- POST /api/dns/domains/update
- POST /api/dns/domains/delete

## etcd Key 规范

- 用户: `/dancer/users/{user-id}`
- Zone: `/dancer/zones/{zone}` (例: `example.com`)
- Domain: `/dancer/domains/{zone}/{domain}` (例: `www.example.com`)
- CoreDNS: `{prefix}/{反转zone}/{domain}/x{n}` (例: `/skydns/com/example/www/x1`)

## 错误类型

按业务场景分类：
- **用户相关**: 用户不存在、用户已存在
- **认证相关**: 无效凭据、密码错误、无效/过期 Token、未授权/禁止访问
- **DNS 相关**: Zone/Domain 不存在或已存在、记录不存在或已存在
- **系统相关**: 无效输入、etcd 不可用

## 开发规范

1. 所有 API 使用 POST 方法
2. JWT 从 Header 获取: `Authorization: Bearer <token>`
3. 管理员权限检查中间件: `RequireAdmin()`
4. 错误响应格式: `{"code": "xxx_error", "message": "..."}` (如 `wrong_password`, `user_not_found`)
5. 时间戳使用 int64(Unix timestamp)

## 关键文件

- `cmd/server/main.go` - 程序入口
- `internal/router/` - 路由配置和访问日志
- `internal/handlers/` - HTTP 处理器
- `internal/services/` - 业务逻辑层
- `internal/storage/etcd/` - etcd 存储操作
- `internal/auth/` - JWT 认证和密码加密
