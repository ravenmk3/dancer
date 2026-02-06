# Dancer

> 轻量级 DNS 管理工具，专为 CoreDNS 设计

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Echo](https://img.shields.io/badge/Echo-v4-00ADD8?style=flat)](https://echo.labstack.com)
[![etcd](https://img.shields.io/badge/etcd-v3-419EDA?style=flat)](https://etcd.io)

**Dancer** 是一个基于 Go + Echo 构建的 DNS 记录管理系统，使用 etcd 作为后端存储，天然适配 CoreDNS 的 etcd 插件。提供 RESTful API 和用户友好的 Web 界面。

---

## ✨ 特性

- 🔐 **JWT 认证** - HS256 签名，支持 Token 刷新
- 👥 **RBAC 权限** - Admin / Normal 角色分离
- 📝 **Zone/Domain 管理** - 清晰的二级域名和子域名管理
- 🔄 **自动 CoreDNS 同步** - 修改记录自动同步到 CoreDNS etcd 格式
- 🗄️ **etcd 存储** - 分布式高可用，双写机制确保数据一致性
- ⚙️ **可配置前缀** - CoreDNS etcd key 前缀可自定义（默认 `/skydns`）
- 🎨 **优雅日志** - logrus + lumberjack，支持轮转
- ⚡ **高性能** - Echo 框架，极简内存占用

---

## 🏗️ 架构

```
HTTP Request
    ↓
Handler (Echo) → 请求解析/响应封装
    ↓
Service → 业务逻辑/事务处理
    ↓
Storage (etcd) → 数据持久化 + CoreDNS 同步
```

---

## 🚀 快速开始

### 1. 配置

```toml
# config.toml
[app]
host = "0.0.0.0"
port = 8080
env = "development"

[etcd]
endpoints = ["http://localhost:2379"]
# CoreDNS etcd 插件的 key 前缀，默认 /skydns
# coredns_prefix = "/skydns"

[jwt]
secret = "your-256-bit-secret"
expiry = 86400

[logger]
level = "info"
file_path = "logs/dancer.log"
```

### 2. 启动

```bash
# 编译
go build -o dancer ./cmd/server

# 运行
./dancer -config config.toml
```

---

## 📡 API 概览

| 端点 | 描述 | 权限 |
|------|------|------|
| `POST /api/auth/login` | 用户登录 | 公开 |
| `POST /api/auth/refresh` | 刷新 Token | JWT |
| `POST /api/me` | 当前用户信息 | JWT |
| `POST /api/me/change-password` | 修改密码 | JWT |
| `POST /api/user/*` | 用户管理 | Admin |
| `POST /api/dns/zones/*` | Zone (二级域名) 管理 | Admin |
| `POST /api/dns/domains/*` | Domain (子域名) 管理 | JWT |

### 认证方式

```http
Authorization: Bearer <jwt-token>
```

---

## 📁 目录结构

```
dancer/
├── cmd/server/           # 程序入口
├── internal/
│   ├── auth/            # JWT / 密码 / 中间件
│   ├── config/          # TOML 配置
│   ├── errors/          # 业务错误
│   ├── handlers/        # HTTP 处理器
│   ├── logger/          # 日志系统
│   ├── models/          # 实体与 DTO
│   ├── router/          # 路由定义
│   ├── services/        # 业务逻辑层
│   └── storage/etcd/    # etcd 客户端
├── assets/              # 前端静态资源
└── config.toml          # 配置文件
```

---

## 🔧 CoreDNS 集成

Dancer 使用双写机制确保 CoreDNS 兼容性：

### 存储结构

```
# Dancer 管理数据
/dancer/zones/example.com              → Zone 元数据
/dancer/domains/example.com/www        → Domain 元数据（含 IP 列表）

# CoreDNS 使用数据（可配置前缀，默认 /skydns）
/skydns/com/example/www/x1             → {"host":"1.1.1.1","ttl":300}
/skydns/com/example/www/x2             → {"host":"1.1.1.2","ttl":300}
```

### CoreDNS 配置示例

```
example.com {
    etcd {
        path /skydns              # 与 dancer 配置一致
        endpoint http://localhost:2379
    }
    cache
}
```

### 工作流程

1. **创建/更新 Domain**：系统自动对比新旧 IP 列表，同步到 CoreDNS
2. **删除 Domain**：级联删除 CoreDNS 记录
3. **删除 Zone**：级联删除所有 Domain 和 CoreDNS 记录

---

## 🛡️ 安全

- 密码使用 **bcrypt** 加密存储
- JWT 支持过期时间配置
- API 全链路 HTTPS 友好
- Admin 操作权限隔离

---

## 📝 默认账号

启动后自动生成：
- **Username**: `admin`
- **Password**: `admin123`

⚠️ **生产环境请立即修改！**

---

## 📖 使用示例

### 1. 创建 Zone (需 Admin)

```bash
curl -X POST http://localhost:8080/api/dns/zones/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"zone":"example.com"}'
```

### 2. 创建 Domain

```bash
curl -X POST http://localhost:8080/api/dns/domains/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "zone": "example.com",
    "domain": "www",
    "ips": ["192.168.1.1", "192.168.1.2"],
    "ttl": 300
  }'
```

### 3. 更新 Domain IP 列表

```bash
curl -X POST http://localhost:8080/api/dns/domains/update \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "zone": "example.com",
    "domain": "www",
    "ips": ["192.168.1.3"],
    "ttl": 600
  }'
```

---

## 📚 文档

- [API 文档](docs/backend-api.md) - 详细的 API 说明
- [设计文档](docs/backend-design.md) - 架构和设计细节
