# 用户个人信息模块测试用例

## 测试概述

测试 Dancer DNS 系统的用户个人信息管理功能，包括获取当前用户信息和修改密码功能。

---

## 测试用例清单

### TC-PROFILE-001: 获取当前用户信息成功

**测试模块**: 用户个人信息模块  
**测试场景**: 使用有效 Token 获取当前用户信息  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动并运行
- 用户已登录并获取有效 Token
- 用户存在于系统中

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
2. 在 Header 中携带有效的 Authorization Token
3. 请求体可为空

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
    "id": "user-uuid",
    "username": "testuser",
    "user_type": "normal",
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```
- 返回的用户信息与 Token 对应的用户一致
- 不包含密码字段

---

### TC-PROFILE-002: 未携带 Token 获取用户信息

**测试模块**: 用户个人信息模块  
**测试场景**: 请求中不包含 Authorization Header  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
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

### TC-PROFILE-003: 使用无效 Token 获取用户信息

**测试模块**: 用户个人信息模块  
**测试场景**: 使用格式错误的 Token  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
2. 在 Header 中携带无效的 Token

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

### TC-PROFILE-004: 使用过期 Token 获取用户信息

**测试模块**: 用户个人信息模块  
**测试场景**: Token 已过期  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 存在已过期的 Token

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
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

---

### TC-PROFILE-005: Token 对应用户不存在

**测试模块**: 用户个人信息模块  
**测试场景**: Token 有效但用户已被删除  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户 A 已登录并获取 Token
- 用户 A 已被其他管理员删除

**测试步骤**:
1. 用户 A 登录获取 Token
2. 管理员删除用户 A
3. 使用用户 A 的 Token 发送 POST 请求到 `/api/me`

**输入数据**:
```http
Authorization: Bearer <user_a_token>
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "user_not_found",
  "message": "用户不存在"
}
```

---

### TC-PROFILE-006: 错误的 Authorization 格式

**测试模块**: 用户个人信息模块  
**测试场景**: Authorization 头缺少 "Bearer " 前缀  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 存在有效 Token

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
2. Authorization Header 格式不正确

**输入数据**:
```http
Authorization: <valid_token>
```
或
```http
Authorization: Basic <valid_token>
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

### TC-PROFILE-007: 正常修改密码成功

**测试模块**: 用户个人信息模块  
**测试场景**: 使用正确的旧密码修改为新密码  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录并获取有效 Token
- 用户当前密码: `oldpassword123`

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 在 Header 中携带有效的 Authorization Token
3. 请求体包含正确的旧密码和符合要求的新密码

**输入数据**:
```http
Authorization: Bearer <valid_token>
Content-Type: application/json

{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "password changed successfully",
  "data": null
}
```
- 使用旧密码登录应失败
- 使用新密码登录应成功

---

### TC-PROFILE-008: 错误的旧密码修改失败

**测试模块**: 用户个人信息模块  
**测试场景**: 提供错误的旧密码  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- 用户当前密码: `currentpassword`

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体包含错误的旧密码

**输入数据**:
```json
{
  "old_password": "wrongpassword",
  "new_password": "newpassword456"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```
- 密码不应被修改
- 使用原密码登录仍应成功

---

### TC-PROFILE-009: 新密码长度不足

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码少于 6 个字符  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体包含少于 6 个字符的新密码

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "12345"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```

---

### TC-PROFILE-010: 新密码为空字符串

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码为空  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体新密码为空字符串

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": ""
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```

---

### TC-PROFILE-011: 缺少 old_password 字段

**测试模块**: 用户个人信息模块  
**测试场景**: 请求体缺少旧密码字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体只包含 new_password

**输入数据**:
```json
{
  "new_password": "newpassword456"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```

---

### TC-PROFILE-012: 缺少 new_password 字段

**测试模块**: 用户个人信息模块  
**测试场景**: 请求体缺少新密码字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体只包含 old_password

**输入数据**:
```json
{
  "old_password": "oldpassword123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```

---

### TC-PROFILE-013: 空请求体修改密码

**测试模块**: 用户个人信息模块  
**测试场景**: 发送空 JSON 对象  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体为空对象

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
  "message": "旧密码错误或新密码格式不符合要求"
}
```

---

### TC-PROFILE-014: 新密码与旧密码相同

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码与旧密码完全一致  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体中 old_password 和 new_password 相同

**输入数据**:
```json
{
  "old_password": "password123",
  "new_password": "password123"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统拒绝
    - HTTP 状态码: 400
    - 响应包含新密码不能与旧密码相同的提示
  - **情况 B**: 系统允许
    - HTTP 状态码: 200
    - 密码修改成功（允许相同密码）
- 需要确认系统设计

---

### TC-PROFILE-015: 超长新密码测试

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码超过最大长度限制  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体包含超长新密码（如 200 个字符）

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "a".repeat(200)
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "旧密码错误或新密码格式不符合要求"
}
```
- 系统应正常处理，不应崩溃

---

### TC-PROFILE-016: 特殊字符新密码

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码包含特殊字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体包含带特殊字符的新密码

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "P@ssw0rd!#$%^&*()"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "password changed successfully",
  "data": null
}
```
- 使用新密码登录应成功

---

### TC-PROFILE-017: 中文密码测试

**测试模块**: 用户个人信息模块  
**测试场景**: 新密码包含中文字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 请求体包含中文密码

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "中文密码测试123"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持中文密码
    - HTTP 状态码: 200
    - 密码修改成功
  - **情况 B**: 系统不支持中文密码
    - HTTP 状态码: 400
    - 返回格式错误信息

---

### TC-PROFILE-018: 修改密码后 Token 失效验证

**测试模块**: 用户个人信息模块  
**测试场景**: 修改密码后验证旧 Token 是否仍有效  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录并获取 Token

**测试步骤**:
1. 获取当前 Token
2. 调用修改密码接口
3. 使用原 Token 访问需要认证的接口

**输入数据**: 修改密码请求

**预期结果**:
- 修改密码成功
- 原 Token 仍然有效（因为 JWT 是无状态的）
- 或者系统应使所有旧 Token 失效（如果实现了 Token 黑名单）
- 需要确认系统设计

---

### TC-PROFILE-019: 未登录用户修改密码

**测试模块**: 用户个人信息模块  
**测试场景**: 未携带 Token 尝试修改密码  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/me/change-password`
2. 不携带 Authorization Header
3. 请求体包含密码信息

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
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

### TC-PROFILE-020: 用户被删除后修改密码

**测试模块**: 用户个人信息模块  
**测试场景**: Token 有效但用户已被删除  
**优先级**: P2 (中)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户 A 已登录并获取 Token
- 用户 A 已被管理员删除

**测试步骤**:
1. 用户 A 登录获取 Token
2. 管理员删除用户 A
3. 使用用户 A 的 Token 发送修改密码请求

**输入数据**:
```http
Authorization: Bearer <user_a_token>
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "user_not_found",
  "message": "用户不存在"
}
```

---

### TC-PROFILE-021: 密码哈希安全性验证

**测试模块**: 用户个人信息模块  
**测试场景**: 验证密码是否正确哈希存储  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- 用户已登录
- 有 etcd 直接访问权限

**测试步骤**:
1. 修改密码
2. 直接查询 etcd 中的用户数据
3. 检查密码字段

**输入数据**:
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**预期结果**:
- etcd 中存储的密码应为 bcrypt 哈希值
- 不应存储明文密码
- 哈希值长度应符合 bcrypt 标准（通常 60 字符）

---

### TC-PROFILE-022: 并发修改密码测试

**测试模块**: 用户个人信息模块  
**测试场景**: 同一用户并发多次修改密码  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- 用户已登录并获取 Token

**测试步骤**:
1. 同时发送 5 个修改密码请求
2. 每个请求使用不同的新密码
3. 使用最后一个修改的密码登录

**输入数据**:
```json
{
  "old_password": "password123",
  "new_password": "newpass{1-5}"
}
```

**预期结果**:
- 只有最后一次成功的修改有效
- 使用最后一次修改的密码登录成功
- 使用之前的密码登录失败
- 系统不应出现数据不一致

---

### TC-PROFILE-023: 获取不同用户类型信息

**测试模块**: 用户个人信息模块  
**测试场景**: 验证不同用户类型返回的信息正确  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 存在 admin 和 normal 两种类型的用户

**测试步骤**:
1. admin 用户登录并获取信息
2. normal 用户登录并获取信息
3. 对比返回的 user_type 字段

**输入数据**: 无

**预期结果**:
- admin 用户返回: `"user_type": "admin"`
- normal 用户返回: `"user_type": "normal"`
- 用户信息完整且正确

---

### TC-PROFILE-024: 时间戳格式验证

**测试模块**: 用户个人信息模块  
**测试场景**: 验证返回的时间戳格式正确  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/me`
2. 检查返回的 created_at 和 updated_at 字段

**输入数据**: 无

**预期结果**:
- created_at 和 updated_at 为 int64 类型
- 值为 Unix 时间戳（秒级）
- 值应在合理范围内（如 2020-2030 年之间）
- updated_at >= created_at

---

## 测试数据准备

### 测试用户

| 用户名 | 当前密码 | 用户类型 | 用途 |
|--------|----------|----------|------|
| testuser | oldpassword123 | normal | 修改密码测试 |
| admin | admin123 | admin | 管理员用户测试 |
| tobedeleted | password123 | normal | 删除用户场景测试 |

### 密码测试数据

| 测试场景 | 旧密码 | 新密码 | 预期结果 |
|----------|--------|--------|----------|
| 正常修改 | oldpassword123 | newpassword456 | 成功 |
| 错误旧密码 | wrongpassword | newpassword456 | 失败 |
| 密码太短 | oldpassword123 | 12345 | 失败 |
| 空新密码 | oldpassword123 | "" | 失败 |

---

## 依赖和前置条件

1. Dancer 服务已启动并监听 8080 端口
2. etcd 服务正常运行
3. 测试数据库中包含预定义的测试用户
4. 网络连接正常
5. 有方法直接查询 etcd 数据（用于安全测试）

---

## 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 密码明文存储 | 低 | 严重 | 验证 etcd 中的密码哈希 |
| 并发修改导致数据不一致 | 低 | 中 | 并发测试覆盖 |
| 旧 Token 在密码修改后仍可用 | 中 | 中 | 验证 Token 失效策略 |

---

## 测试通过标准

- 所有 P0 和 P1 测试用例通过
- 密码以哈希形式存储，永不存储明文
- 修改密码后原 Token 行为符合设计
- 系统在各种边界条件下保持稳定
