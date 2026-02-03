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
- 📝 **DNS 管理** - CRUD 操作，完美适配 CoreDNS etcd 格式
- 🗄️ **etcd 存储** - 分布式高可用，域名存储
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
Storage (etcd) → 数据持久化
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
| `POST /api/dns/records/*` | DNS 记录管理 | JWT |

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

Dancer 使用与 CoreDNS etcd 插件兼容的 Key 格式：

```
/coredns/{反转域名}/{记录名}

示例：
  github.com    → /coredns/com/github
  api.github.com → /coredns/com/github/api
```

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
