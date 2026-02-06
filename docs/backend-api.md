# Dancer DNS Management API

## 概述

Dancer 是一个基于 etcd 的 DNS 管理工具，专为 CoreDNS 设计，提供 RESTful API 接口。

**基础 URL**: `http://{host}:{port}/api`

**默认端口**: `8080`

## 认证机制

所有 API (除登录、刷新 Token、健康检查外) 都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 权限级别

1. **普通用户 (normal)**: 可以管理 DNS Domain
2. **管理员 (admin)**: 可以管理用户、Zone 和 Domain

## 响应格式

### 成功响应

```json
{
  "code": "success",
  "message": "success",
  "data": { ... }
}
```

### 错误响应

```json
{
  "code": "error_code",
  "message": "error message"
}
```

## 错误码

| 错误码 | HTTP 状态码 | 说明 |
|--------|-------------|------|
| `success` | 200 | 操作成功 |
| `invalid_input` | 400 | 请求参数无效 |
| `invalid_credentials` | 401 | 用户名或密码错误 |
| `unauthorized` | 401 | 未认证或 Token 无效 |
| `forbidden` | 403 | 权限不足 |
| `user_not_found` | 404 | 用户不存在 |
| `user_exists` | 409 | 用户已存在 |
| `zone_not_found` | 404 | Zone (二级域名) 不存在 |
| `zone_exists` | 409 | Zone 已存在 |
| `domain_not_found` | 404 | Domain 不存在 |
| `domain_exists` | 409 | Domain 已存在 |
| `service_unavailable` | 503 | etcd 服务不可用 |
| `internal_error` | 500 | 服务器内部错误 |

---

## API 端点

### 认证模块

#### 1. 用户登录

**请求**

```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "admin123"
}
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user-uuid",
      "username": "admin",
      "user_type": "admin",
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `invalid_credentials` (401): 用户名或密码错误
- `invalid_input` (400): 请求参数缺失或格式错误

---

#### 2. 刷新 Token

**请求**

```http
POST /api/auth/refresh
Authorization: Bearer <current_token>
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**错误场景**

- `unauthorized` (401): Token 无效或过期

---

### 用户个人信息模块

#### 3. 获取当前用户信息

**请求**

```http
POST /api/me
Authorization: Bearer <token>
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "id": "user-uuid",
    "username": "admin",
    "user_type": "admin",
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```

**错误场景**

- `user_not_found` (404): 用户不存在
- `unauthorized` (401): Token 无效或过期

---

#### 4. 修改当前用户密码

**请求**

```http
POST /api/me/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**字段约束**

- `new_password`: 最小长度 6 个字符

**响应**

```json
{
  "code": "success",
  "message": "password changed successfully",
  "data": null
}
```

**错误场景**

- `invalid_input` (400): 旧密码错误或新密码格式不符合要求
- `user_not_found` (404): 用户不存在
- `unauthorized` (401): Token 无效或过期

---

### 用户管理模块 (Admin)

#### 5. 列出所有用户

**请求**

```http
POST /api/user/list
Authorization: Bearer <token> (需 Admin 权限)
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "users": [
      {
        "id": "user-uuid-1",
        "username": "admin",
        "user_type": "admin",
        "created_at": 1704067200,
        "updated_at": 1704067200
      },
      {
        "id": "user-uuid-2",
        "username": "user1",
        "user_type": "normal",
        "created_at": 1704067200,
        "updated_at": 1704067200
      }
    ]
  }
}
```

**错误场景**

- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 6. 创建用户

**请求**

```http
POST /api/user/create
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "username": "newuser",
  "password": "password123",
  "user_type": "normal"
}
```

**字段约束**

- `username`: 3-32 个字符，必填
- `password`: 最少 6 个字符，必填
- `user_type`: `admin` 或 `normal`，必填

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "id": "user-uuid",
    "username": "newuser",
    "user_type": "normal",
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```

**错误场景**

- `user_exists` (409): 用户名已存在
- `invalid_input` (400): 请求参数不符合约束
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 7. 更新用户

**请求**

```http
POST /api/user/update
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "id": "user-uuid",
  "username": "updateduser",
  "password": "newpassword123",
  "user_type": "normal"
}
```

**字段约束**

- `id`: 必填
- `username`: 3-32 个字符，可选
- `password`: 最少 6 个字符，可选
- `user_type`: `admin` 或 `normal`，可选

**响应**

```json
{
  "code": "success",
  "message": "user updated successfully",
  "data": null
}
```

**错误场景**

- `user_not_found` (404): 用户不存在
- `user_exists` (409): 更新后的用户名已存在
- `invalid_input` (400): 请求参数不符合约束
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 8. 删除用户

**请求**

```http
POST /api/user/delete
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "id": "user-uuid"
}
```

**字段约束**

- `id`: 必填

**响应**

```json
{
  "code": "success",
  "message": "user deleted successfully",
  "data": null
}
```

**错误场景**

- `user_not_found` (404): 用户不存在
- `forbidden` (403): 不能删除默认管理员或权限不足
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

### Zone 管理模块 (Admin)

Zone 代表二级域名，如 `example.com`。

#### 9. 列出所有 Zone

**请求**

```http
POST /api/dns/zones/list
Authorization: Bearer <token> (需 Admin 权限)
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "zones": [
      {
        "zone": "example.com",
        "record_count": 5,
        "created_at": 1704067200,
        "updated_at": 1704067200
      }
    ]
  }
}
```

**错误场景**

- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 10. 获取 Zone 详情

**请求**

```http
POST /api/dns/zones/get
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "zone": "example.com"
}
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "zone": {
      "zone": "example.com",
      "record_count": 5,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 11. 创建 Zone

**请求**

```http
POST /api/dns/zones/create
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "zone": "example.com"
}
```

**字段约束**

- `zone`: 有效的二级域名（FQDN），必填

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "zone": {
      "zone": "example.com",
      "record_count": 0,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `zone_exists` (409): Zone 已存在
- `invalid_input` (400): 请求参数不符合约束
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 12. 更新 Zone

**请求**

```http
POST /api/dns/zones/update
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "zone": "example.com"
}
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "zone": {
      "zone": "example.com",
      "record_count": 5,
      "created_at": 1704067200,
      "updated_at": 1704153600
    }
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

#### 13. 删除 Zone

**请求**

```http
POST /api/dns/zones/delete
Authorization: Bearer <token> (需 Admin 权限)
Content-Type: application/json

{
  "zone": "example.com"
}
```

**说明**

- 删除 Zone 会**级联删除**该 Zone 下的所有 Domain 及其 CoreDNS 记录

**响应**

```json
{
  "code": "success",
  "message": "Zone deleted successfully",
  "data": null
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `forbidden` (403): 非 Admin 用户
- `unauthorized` (401): Token 无效或过期

---

### Domain 管理模块 (JWT)

Domain 代表完整域名（子域名），如 Zone `example.com` 下的 `www` 或 `@`（根）。

#### 14. 列出 Zone 下所有 Domain

**请求**

```http
POST /api/dns/domains/list
Authorization: Bearer <token>
Content-Type: application/json

{
  "zone": "example.com"
}
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "domains": [
      {
        "zone": "example.com",
        "domain": "www",
        "name": "www.example.com",
        "ips": ["192.168.1.1", "192.168.1.2"],
        "ttl": 300,
        "record_count": 2,
        "created_at": 1704067200,
        "updated_at": 1704067200
      },
      {
        "zone": "example.com",
        "domain": "@",
        "name": "example.com",
        "ips": ["192.168.1.10"],
        "ttl": 600,
        "record_count": 1,
        "created_at": 1704067200,
        "updated_at": 1704067200
      }
    ]
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `unauthorized` (401): Token 无效或过期

---

#### 15. 获取 Domain 详情

**请求**

```http
POST /api/dns/domains/get
Authorization: Bearer <token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www"
}
```

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "domain": {
      "zone": "example.com",
      "domain": "www",
      "name": "www.example.com",
      "ips": ["192.168.1.1", "192.168.1.2"],
      "ttl": 300,
      "record_count": 2,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `domain_not_found` (404): Domain 不存在
- `unauthorized` (401): Token 无效或过期

---

#### 16. 创建 Domain

**请求**

```http
POST /api/dns/domains/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.1", "192.168.1.2"],
  "ttl": 300
}
```

**字段约束**

- `zone`: 已存在的 Zone 名称，必填
- `domain`: 子域名部分（如 `www` 或 `@` 代表根），必填
- `ips`: IP 地址数组，必填，每个 IP 必须是有效格式
- `ttl`: TTL (秒)，必填，最小值 1

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "domain": {
      "zone": "example.com",
      "domain": "www",
      "name": "www.example.com",
      "ips": ["192.168.1.1", "192.168.1.2"],
      "ttl": 300,
      "record_count": 2,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在，需要先创建 Zone
- `domain_exists` (409): Domain 已存在
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

#### 17. 更新 Domain

**请求**

```http
POST /api/dns/domains/update
Authorization: Bearer <token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.3", "192.168.1.4"],
  "ttl": 600
}
```

**字段约束**

- `zone`: 必填
- `domain`: 必填
- `ips`: IP 地址数组，必填，会**替换**现有的所有 IP
- `ttl`: 可选，不填则保持原值

**说明**

- 系统会自动比较新旧 IP 列表，添加新 IP、删除不再使用的 IP，保持 CoreDNS 记录与请求一致

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "domain": {
      "zone": "example.com",
      "domain": "www",
      "name": "www.example.com",
      "ips": ["192.168.1.3", "192.168.1.4"],
      "ttl": 600,
      "record_count": 2,
      "created_at": 1704067200,
      "updated_at": 1704153600
    }
  }
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `domain_not_found` (404): Domain 不存在
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

#### 18. 删除 Domain

**请求**

```http
POST /api/dns/domains/delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www"
}
```

**说明**

- 删除 Domain 会**级联删除**该 Domain 的所有 CoreDNS 记录

**响应**

```json
{
  "code": "success",
  "message": "Domain deleted successfully",
  "data": null
}
```

**错误场景**

- `zone_not_found` (404): Zone 不存在
- `domain_not_found` (404): Domain 不存在
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

## 健康检查

### 端点

```http
GET /api/health
POST /api/health
```

**说明**: 公开端点，无需认证。反映 etcd 连接状态。

**响应（所有 components 为 up）**:

```json
{
  "status": "up",
  "components": {
    "etcd": "up"
  }
}
```

HTTP 状态码: 200

**响应（任一 component 为 down）**:

```json
{
  "status": "down",
  "components": {
    "etcd": "down"
  }
}
```

HTTP 状态码: 503

---

## 数据模型

### User

| 字段 | 类型 | 说明 |
|------|------|------|
| `id` | string | 用户唯一标识 |
| `username` | string | 用户名 |
| `user_type` | string | 用户类型: `admin` / `normal` |
| `created_at` | int64 | 创建时间 (Unix 时间戳) |
| `updated_at` | int64 | 更新时间 (Unix 时间戳) |

### Zone

| 字段 | 类型 | 说明 |
|------|------|------|
| `zone` | string | 二级域名 (如 `example.com`) |
| `record_count` | int | 该 Zone 下的 Domain 数量 |
| `created_at` | int64 | 创建时间 (Unix 时间戳) |
| `updated_at` | int64 | 更新时间 (Unix 时间戳) |

### Domain

| 字段 | 类型 | 说明 |
|------|------|------|
| `zone` | string | 所属 Zone (如 `example.com`) |
| `domain` | string | 子域名部分 (如 `www` 或 `@`) |
| `name` | string | 完整域名 (如 `www.example.com`) |
| `ips` | []string | IP 地址列表 |
| `ttl` | int | TTL (秒) |
| `record_count` | int | IP 记录数量 |
| `created_at` | int64 | 创建时间 (Unix 时间戳) |
| `updated_at` | int64 | 更新时间 (Unix 时间戳) |

---

## etcd Key 规范

### 用户数据

```
/dancer/users/{user-id}
```

### Dancer 管理数据

#### Zone

```
/dancer/zones/{zone}
```

示例: `/dancer/zones/example.com`

#### Domain

```
/dancer/domains/{zone}/{domain}
```

示例: 
- `/dancer/domains/example.com/www` (www.example.com)
- `/dancer/domains/example.com/@` (example.com 根域名)

### CoreDNS 记录数据

```
/{prefix}/{反转zone}/{domain}/x{n}
```

- `{prefix}`: CoreDNS etcd 前缀，默认 `/skydns`，可配置
- `{反转zone}`: Zone 的反转格式，如 `example.com` → `com/example`
- `{domain}`: 子域名
- `x{n}`: 记录索引，如 `x1`, `x2`...

示例 (prefix=/skydns):
- `www.example.com` → `/skydns/com/example/www/x1`, `/skydns/com/example/www/x2`...
- `example.com` (根) → `/skydns/com/example/x1`...

---

## 使用示例

### 1. 登录并获取 Token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 2. 创建 Zone (需 Admin)

```bash
curl -X POST http://localhost:8080/api/dns/zones/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"zone":"example.com"}'
```

### 3. 创建 Domain

```bash
curl -X POST http://localhost:8080/api/dns/domains/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"zone":"example.com","domain":"www","ips":["192.168.1.1","192.168.1.2"],"ttl":300}'
```

### 4. 列出 Zone 下所有 Domain

```bash
curl -X POST http://localhost:8080/api/dns/domains/list \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"zone":"example.com"}'
```

### 5. 更新 Domain IP 列表

```bash
curl -X POST http://localhost:8080/api/dns/domains/update \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"zone":"example.com","domain":"www","ips":["192.168.1.3"],"ttl":600}'
```

---

## 注意事项

1. 所有 API 使用 POST 方法 (除 /api/health 外)
2. 时间戳使用 Unix 时间戳 (int64)
3. 密码使用 bcrypt 加密存储
4. JWT Token 过期时间可配置
5. 默认管理员账号在系统启动时自动创建
6. 健康检查端点 /api/health 同时支持 GET 和 POST 方法
7. 创建 Domain 前必须先创建对应的 Zone
8. 删除 Zone 会级联删除其下所有 Domain 和 CoreDNS 记录
9. Domain 的 `ips` 字段在更新时会**完全替换**原有 IP 列表
