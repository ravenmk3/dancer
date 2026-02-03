# Dancer DNS Management API

## 概述

Dancer 是一个基于 etcd 的 DNS 管理工具，提供 RESTful API 接口。

**基础 URL**: `http://{host}:{port}/api`

**默认端口**: `8080`

## 认证机制

所有 API (除登录和刷新 Token 外) 都需要在请求头中携带 JWT Token：

```
Authorization: Bearer <token>
```

### 权限级别

1. **普通用户 (normal)**: 可以管理自己的 DNS 记录
2. **管理员 (admin)**: 可以管理用户和 DNS 记录

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
| `record_not_found` | 404 | DNS 记录不存在 |
| `record_exists` | 409 | DNS 记录已存在 |
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

### DNS 记录管理模块

#### 9. 列出 DNS 记录

**请求**

```http
POST /api/dns/records/list
Authorization: Bearer <token>
Content-Type: application/json

{
  "domain": "example.com"
}
```

**说明**

- `domain`: 可选，为空时列出所有 DNS 记录

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "records": [
      {
        "key": "/coredns/com/example",
        "id": "record-uuid",
        "domain": "example.com",
        "ip": "192.168.1.1",
        "ttl": 300,
        "created_at": 1704067200,
        "updated_at": 1704067200
      }
    ]
  }
}
```

**错误场景**

- `unauthorized` (401): Token 无效或过期

---

#### 10. 创建 DNS 记录

**请求**

```http
POST /api/dns/records/create
Authorization: Bearer <token>
Content-Type: application/json

{
  "domain": "example.com",
  "ip": "192.168.1.1",
  "ttl": 300
}
```

**字段约束**

- `domain`: 有效的 FQDN (完全限定域名)，必填
- `ip`: 有效的 IP 地址，必填
- `ttl`: 最小值 60 秒，必填

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "record": {
      "key": "/coredns/com/example",
      "id": "record-uuid",
      "domain": "example.com",
      "ip": "192.168.1.1",
      "ttl": 300,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

**错误场景**

- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

#### 11. 更新 DNS 记录

**请求**

```http
POST /api/dns/records/update
Authorization: Bearer <token>
Content-Type: application/json

{
  "key": "/coredns/com/example",
  "domain": "example.com",
  "ip": "192.168.1.2",
  "ttl": 600
}
```

**字段约束**

- `key`: 必填，DNS 记录的 etcd key
- `domain`: 有效的 FQDN，可选
- `ip`: 有效的 IP 地址，可选
- `ttl`: 最小值 60 秒，可选

**响应**

```json
{
  "code": "success",
  "message": "success",
  "data": {
    "record": {
      "key": "/coredns/com/example",
      "id": "record-uuid",
      "domain": "example.com",
      "ip": "192.168.1.2",
      "ttl": 600,
      "created_at": 1704067200,
      "updated_at": 1704153600
    }
  }
}
```

**错误场景**

- `record_not_found` (404): DNS 记录不存在
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

#### 12. 删除 DNS 记录

**请求**

```http
POST /api/dns/records/delete
Authorization: Bearer <token>
Content-Type: application/json

{
  "key": "/coredns/com/example"
}
```

**字段约束**

- `key`: 必填，DNS 记录的 etcd key

**响应**

```json
{
  "code": "success",
  "message": "DNS record deleted successfully",
  "data": null
}
```

**错误场景**

- `record_not_found` (404): DNS 记录不存在
- `invalid_input` (400): 请求参数不符合约束
- `unauthorized` (401): Token 无效或过期

---

## 健康检查

### 新增端点

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

### DNSRecord

| 字段 | 类型 | 说明 |
|------|------|------|
| `key` | string | etcd 存储 key (如 `/coredns/com/example`) |
| `id` | string | 记录唯一标识 |
| `domain` | string | 域名 (如 `example.com`) |
| `ip` | string | IP 地址 |
| `ttl` | int64 | TTL (秒)，最小 60 |
| `created_at` | int64 | 创建时间 (Unix 时间戳) |
| `updated_at` | int64 | 更新时间 (Unix 时间戳) |

---

## etcd Key 规范

### 用户数据

```
/dance/users/{user-id}
```

### DNS 记录

域名反转存储，例如 `github.com`:

```
/coredns/com/github
```

子域名 `www.github.com`:

```
/coredns/com/github/www
```

---

## 使用示例

### 登录并获取 Token

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}'
```

### 创建 DNS 记录

```bash
curl -X POST http://localhost:8080/api/dns/records/create \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"domain":"example.com","ip":"192.168.1.1","ttl":300}'
```

### 列出 DNS 记录

```bash
curl -X POST http://localhost:8080/api/dns/records/list \
  -H "Authorization: Bearer <token>"
```

---

## 注意事项

1. 所有 API 使用 POST 方法 (除 /api/health 外)
2. 时间戳使用 Unix 时间戳 (int64)
3. 密码使用 bcrypt 加密存储
4. JWT Token 过期时间可配置
5. 默认管理员账号在系统启动时自动创建
6. 健康检查端点 /api/health 同时支持 GET 和 POST 方法
7. 错误响应由统一的 HTTPErrorHandler 处理
