# 认证模块测试用例

## 测试概述

测试 Dancer DNS 系统的认证功能，包括登录和 Token 刷新功能。

---

## 测试用例清单

### TC-AUTH-001: 正常登录成功

**测试模块**: 认证模块  
**测试场景**: 使用正确的用户名和密码登录  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动并运行
- 存在测试用户: `testuser` / `password123`
- 用户状态正常

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含正确的用户名和密码

**输入数据**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": "user-uuid",
      "username": "testuser",
      "user_type": "normal",
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```
- Token 格式为有效的 JWT
- 用户信息包含所有必需字段

---

### TC-AUTH-002: 错误密码登录失败

**测试模块**: 认证模块  
**测试场景**: 使用错误的密码登录  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 存在测试用户: `testuser` / `password123`

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含正确的用户名和错误的密码

**输入数据**:
```json
{
  "username": "testuser",
  "password": "wrongpassword"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```
- 不返回任何 Token

---

### TC-AUTH-003: 不存在的用户登录

**测试模块**: 认证模块  
**测试场景**: 使用不存在的用户名登录  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户 `nonexistentuser` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含不存在的用户名

**输入数据**:
```json
{
  "username": "nonexistentuser",
  "password": "somepassword"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```
- 错误信息不应泄露用户存在性（与错误密码返回相同信息）

---

### TC-AUTH-004: 缺失用户名参数

**测试模块**: 认证模块  
**测试场景**: 请求体缺少 username 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体只包含 password 字段

**输入数据**:
```json
{
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-005: 缺失密码参数

**测试模块**: 认证模块  
**测试场景**: 请求体缺少 password 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体只包含 username 字段

**输入数据**:
```json
{
  "username": "testuser"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-006: 空请求体

**测试模块**: 认证模块  
**测试场景**: 发送空 JSON 对象  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体为空对象 `{}`

**输入数据**:
```json
{}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-007: 用户名和密码均为空字符串

**测试模块**: 认证模块  
**测试场景**: username 和 password 都是空字符串  
**优先级**: P2 (中)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含空字符串用户名和密码

**输入数据**:
```json
{
  "username": "",
  "password": ""
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-008: 用户名只包含空白字符

**测试模块**: 认证模块  
**测试场景**: username 为空白字符（空格、Tab等）  
**优先级**: P2 (中)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含空白用户名

**输入数据**:
```json
{
  "username": "   ",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-009: SQL 注入攻击测试

**测试模块**: 认证模块  
**测试场景**: 尝试 SQL 注入攻击  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- 存在测试用户: `testuser` / `password123`

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含 SQL 注入字符串

**输入数据**:
```json
{
  "username": "testuser' OR '1'='1",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```
- 系统应正确处理特殊字符，防止注入攻击

---

### TC-AUTH-010: XSS 攻击测试

**测试模块**: 认证模块  
**测试场景**: 尝试 XSS 攻击  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含 XSS 攻击字符串

**输入数据**:
```json
{
  "username": "<script>alert('xss')</script>",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```
- 系统应对输入进行安全处理

---

### TC-AUTH-011: 超长用户名测试

**测试模块**: 认证模块  
**测试场景**: 用户名超过最大长度限制  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含超长用户名（200个字符）

**输入数据**:
```json
{
  "username": "a".repeat(200),
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数缺失或格式错误"
}
```

---

### TC-AUTH-012: 超长密码测试

**测试模块**: 认证模块  
**测试场景**: 密码超过最大长度限制  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含超长密码（1000个字符）

**输入数据**:
```json
{
  "username": "testuser",
  "password": "a".repeat(1000)
}
```

**预期结果**:
- HTTP 状态码: 401 或 400
- 系统应正常处理，不应崩溃

---

### TC-AUTH-013: 特殊字符用户名测试

**测试模块**: 认证模块  
**测试场景**: 用户名包含特殊字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含特殊字符用户名

**输入数据**:
```json
{
  "username": "user@#$%^&*()",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```

---

### TC-AUTH-014: 中文用户名测试

**测试模块**: 认证模块  
**测试场景**: 用户名包含中文字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/login`
2. 请求体包含中文用户名

**输入数据**:
```json
{
  "username": "测试用户",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "invalid_credentials",
  "message": "用户名或密码错误"
}
```
- 系统应正确处理 Unicode 字符

---

### TC-AUTH-015: 正常刷新 Token

**测试模块**: 认证模块  
**测试场景**: 使用有效的 Token 刷新  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录并获取有效 Token

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/refresh`
2. 在 Header 中携带有效的 Authorization Token

**输入数据**:
```http
Authorization: Bearer <valid_token>
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```
- 新 Token 与旧 Token 不同
- 新 Token 应有效且可正常使用

---

### TC-AUTH-016: 过期 Token 刷新失败

**测试模块**: 认证模块  
**测试场景**: 使用过期的 Token 刷新  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录但 Token 已过期

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/refresh`
2. 在 Header 中携带已过期的 Token

**输入数据**:
```http
Authorization: Bearer <expired_token>
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "unauthorized",
  "message": "Token 无效或过期"
}
```
- 不返回新 Token

---

### TC-AUTH-017: 无效 Token 格式刷新

**测试模块**: 认证模块  
**测试场景**: 使用格式错误的 Token 刷新  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/refresh`
2. 在 Header 中携带格式错误的 Token

**输入数据**:
```http
Authorization: Bearer invalid.token.here
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "unauthorized",
  "message": "Token 无效或过期"
}
```

---

### TC-AUTH-018: 缺少 Authorization Header

**测试模块**: 认证模块  
**测试场景**: 请求头中不包含 Authorization  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/refresh`
2. 不携带 Authorization Header

**输入数据**: 无

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "unauthorized",
  "message": "Token 无效或过期"
}
```

---

### TC-AUTH-019: 错误的 Authorization 格式

**测试模块**: 认证模块  
**测试场景**: Authorization 头格式不正确（缺少 Bearer）  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 存在有效 Token

**测试步骤**:
1. 发送 POST 请求到 `/api/auth/refresh`
2. Authorization Header 缺少 "Bearer " 前缀

**输入数据**:
```http
Authorization: <valid_token>
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "unauthorized",
  "message": "Token 无效或过期"
}
```

---

### TC-AUTH-020: 篡改的 Token 刷新

**测试模块**: 认证模块  
**测试场景**: Token 被篡改（签名无效）  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- 存在有效 Token

**测试步骤**:
1. 修改有效 Token 的某一位字符
2. 发送 POST 请求到 `/api/auth/refresh`
3. 使用篡改后的 Token

**输入数据**:
```http
Authorization: Bearer <tampered_token>
```

**预期结果**:
- HTTP 状态码: 401
- 响应 JSON:
```json
{
  "code": "unauthorized",
  "message": "Token 无效或过期"
}
```
- 系统应检测出 Token 被篡改

---

### TC-AUTH-021: 并发登录测试

**测试模块**: 认证模块  
**测试场景**: 同一用户并发多次登录  
**优先级**: P2 (中)  
**测试类型**: 性能/并发测试

**前置条件**:
- 系统已启动
- 存在测试用户: `testuser` / `password123`

**测试步骤**:
1. 同时发送 10 个登录请求
2. 所有请求使用相同的用户名和密码

**输入数据**:
```json
{
  "username": "testuser",
  "password": "password123"
}
```

**预期结果**:
- 所有请求都应返回 200
- 每个请求返回不同的 Token
- 系统不应出现竞态条件

---

### TC-AUTH-022: 多次错误登录后锁定（如支持）

**测试模块**: 认证模块  
**测试场景**: 连续多次输入错误密码  
**优先级**: P2 (中)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- 存在测试用户: `testuser` / `password123`

**测试步骤**:
1. 连续发送 5 次错误密码登录请求
2. 第 6 次发送正确密码

**输入数据**:
```json
{
  "username": "testuser",
  "password": "wrongpassword"
}
```
(前5次)

```json
{
  "username": "testuser",
  "password": "password123"
}
```
(第6次)

**预期结果**:
- 前 5 次返回 401 (invalid_credentials)
- 第 6 次：
  - 如果支持锁定：返回 403 或类似状态码，提示账户已锁定
  - 如果不支持锁定：返回 200，登录成功

---

### TC-AUTH-023: 登录后 Token 验证

**测试模块**: 认证模块  
**测试场景**: 验证登录返回的 Token 是否可用  
**优先级**: P0 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- 存在测试用户

**测试步骤**:
1. 调用登录接口获取 Token
2. 使用获取的 Token 调用需要认证的接口（如 /api/me）

**输入数据**: 登录请求

**预期结果**:
- 登录成功，返回 Token
- 使用 Token 调用其他接口成功（返回 200）
- Token 应包含正确的用户信息

---

### TC-AUTH-024: Token 过期时间验证

**测试模块**: 认证模块  
**测试场景**: 验证 Token 过期时间设置正确  
**优先级**: P1 (高)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- JWT Token 过期时间配置为较短（如 1 分钟用于测试）

**测试步骤**:
1. 登录获取 Token
2. 立即使用 Token 调用接口（应成功）
3. 等待 Token 过期
4. 再次使用 Token 调用接口

**输入数据**: 登录请求

**预期结果**:
- 步骤 2: 返回 200，请求成功
- 步骤 4: 返回 401 (unauthorized)，Token 已过期

---

## 测试数据准备

### 测试用户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| admin | admin123 | admin | 管理员用户测试 |
| testuser | password123 | normal | 普通用户测试 |
| normaluser | userpass456 | normal | 普通用户测试 |

### 环境配置

- 测试环境: http://localhost:8080
- JWT 过期时间: 建议设置为 15 分钟（或根据配置）
- etcd 状态: 正常运行

---

## 依赖和前置条件

1. Dancer 服务已启动并监听 8080 端口
2. etcd 服务正常运行
3. 测试数据库中包含预定义的测试用户
4. 网络连接正常

---

## 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| etcd 不可用 | 低 | 高 | 测试前检查 etcd 状态 |
| 并发测试导致数据不一致 | 中 | 中 | 使用独立测试数据 |
| Token 过期时间配置过长 | 低 | 中 | 测试环境缩短过期时间 |

---

## 测试通过标准

- 所有 P0 和 P1 测试用例通过
- 安全测试用例（SQL注入、XSS）正确处理
- 系统在各种边界条件下保持稳定
- 无内存泄漏或性能问题
