# 用户管理模块测试用例（Admin）

## 测试概述

测试 Dancer DNS 系统的用户管理功能（仅 Admin 权限），包括用户列表查询、创建、更新和删除功能。

---

## 测试用例清单

### TC-USER-001: Admin 列出所有用户

**测试模块**: 用户管理模块  
**测试场景**: Admin 用户成功获取所有用户列表  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动并运行
- Admin 用户已登录并获取有效 Token
- 系统中存在多个用户（admin 和 normal 类型）

**测试步骤**:
1. Admin 用户登录获取 Token
2. 发送 POST 请求到 `/api/user/list`
3. 在 Header 中携带 Admin 的 Authorization Token

**输入数据**:
```http
Authorization: Bearer <admin_token>
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
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
- 返回所有用户的列表
- 用户信息不包含密码字段
- 列表包含 admin 和 normal 两种类型的用户

---

### TC-USER-002: Normal 用户访问用户列表

**测试模块**: 用户管理模块  
**测试场景**: Normal 用户尝试获取用户列表  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录并获取有效 Token

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/user/list`
3. 在 Header 中携带 Normal 用户的 Token

**输入数据**:
```http
Authorization: Bearer <normal_token>
```

**预期结果**:
- HTTP 状态码: 403
- 响应 JSON:
```json
{
  "code": "forbidden",
  "message": "权限不足"
}
```
- Normal 用户无法查看用户列表

---

### TC-USER-003: 未登录访问用户列表

**测试模块**: 用户管理模块  
**测试场景**: 未携带 Token 访问用户列表  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/user/list`
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

### TC-USER-004: 使用过期 Token 访问用户列表

**测试模块**: 用户管理模块  
**测试场景**: 使用已过期的 Admin Token  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 存在已过期的 Admin Token

**测试步骤**:
1. 发送 POST 请求到 `/api/user/list`
2. 在 Header 中携带已过期的 Admin Token

**输入数据**:
```http
Authorization: Bearer <expired_admin_token>
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

### TC-USER-005: Admin 创建 Normal 用户成功

**测试模块**: 用户管理模块  
**测试场景**: Admin 成功创建普通用户  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户名 `newnormaluser` 不存在

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/user/create`
3. 请求体包含新用户信息

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "username": "newnormaluser",
  "password": "password123",
  "user_type": "normal"
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
    "id": "user-uuid",
    "username": "newnormaluser",
    "user_type": "normal",
    "created_at": 1704067200,
    "updated_at": 1704067200
  }
}
```
- 新用户可以使用 username/password 登录
- 创建的用户 type 为 normal
- 密码被正确哈希存储

---

### TC-USER-006: Admin 创建 Admin 用户成功

**测试模块**: 用户管理模块  
**测试场景**: Admin 成功创建另一个 Admin 用户  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户名 `newadminuser` 不存在

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/user/create`
3. 请求体包含 user_type 为 admin

**输入数据**:
```json
{
  "username": "newadminuser",
  "password": "adminpass123",
  "user_type": "admin"
}
```

**预期结果**:
- HTTP 状态码: 200
- 创建的用户 type 为 admin
- 新 Admin 可以执行 Admin 操作

---

### TC-USER-007: 创建已存在的用户名

**测试模块**: 用户管理模块  
**测试场景**: 尝试创建用户名已存在的用户  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `existinguser` 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体使用已存在的用户名

**输入数据**:
```json
{
  "username": "existinguser",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 409
- 响应 JSON:
```json
{
  "code": "user_exists",
  "message": "用户已存在"
}
```
- 不应创建重复用户
- 原有用户数据保持不变

---

### TC-USER-008: Normal 用户尝试创建用户

**测试模块**: 用户管理模块  
**测试场景**: Normal 用户尝试创建新用户  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/user/create`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "username": "anotheruser",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 403
- 响应 JSON:
```json
{
  "code": "forbidden",
  "message": "权限不足"
}
```
- 用户不应被创建

---

### TC-USER-009: 创建用户缺少 username 字段

**测试模块**: 用户管理模块  
**测试场景**: 请求体缺少必填的 username  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体不包含 username 字段

**输入数据**:
```json
{
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-010: 创建用户缺少 password 字段

**测试模块**: 用户管理模块  
**测试场景**: 请求体缺少必填的 password  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体不包含 password 字段

**输入数据**:
```json
{
  "username": "newuser",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-011: 创建用户缺少 user_type 字段

**测试模块**: 用户管理模块  
**测试场景**: 请求体缺少必填的 user_type  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体不包含 user_type 字段

**输入数据**:
```json
{
  "username": "newuser",
  "password": "password123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-012: 用户名长度小于 3 个字符

**测试模块**: 用户管理模块  
**测试场景**: 用户名只有 1-2 个字符  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 username 为 "ab"（2个字符）

**输入数据**:
```json
{
  "username": "ab",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-013: 用户名长度超过 32 个字符

**测试模块**: 用户管理模块  
**测试场景**: 用户名超过 32 个字符  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 username 为 33 个字符

**输入数据**:
```json
{
  "username": "thisisaverylongusernamethatexceedsthirtytwo",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-014: 密码长度少于 6 个字符

**测试模块**: 用户管理模块  
**测试场景**: 密码只有 5 个字符  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 password 为 5 个字符

**输入数据**:
```json
{
  "username": "newuser",
  "password": "12345",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-015: 无效的 user_type 值

**测试模块**: 用户管理模块  
**测试场景**: user_type 不是 admin 或 normal  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 user_type 为无效值

**输入数据**:
```json
{
  "username": "newuser",
  "password": "password123",
  "user_type": "superadmin"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-016: 空用户名

**测试模块**: 用户管理模块  
**测试场景**: username 为空字符串  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 username 为空字符串

**输入数据**:
```json
{
  "username": "",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-017: 用户名包含特殊字符

**测试模块**: 用户管理模块  
**测试场景**: username 包含特殊字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 username 包含特殊字符

**输入数据**:
```json
{
  "username": "user@#$%",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持特殊字符用户名
    - HTTP 状态码: 200
    - 用户创建成功
  - **情况 B**: 系统不支持
    - HTTP 状态码: 400
    - 返回格式错误
- 需要确认系统设计规范

---

### TC-USER-018: 中文用户名测试

**测试模块**: 用户管理模块  
**测试场景**: username 包含中文字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. 请求体 username 为中文

**输入数据**:
```json
{
  "username": "测试用户",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持中文用户名
    - HTTP 状态码: 200
    - 用户创建成功
  - **情况 B**: 系统不支持
    - HTTP 状态码: 400
    - 返回格式错误

---

### TC-USER-019: 创建用户后验证登录

**测试模块**: 用户管理模块  
**测试场景**: 创建用户后验证该用户可以正常登录  
**优先级**: P0 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. Admin 创建新用户
2. 使用新用户的凭据登录

**输入数据**:
```json
{
  "username": "testnewuser",
  "password": "testpass123",
  "user_type": "normal"
}
```

**预期结果**:
- 用户创建成功
- 新用户可以使用 username/password 登录成功
- 登录返回的 user_type 与创建时一致

---

### TC-USER-020: Admin 更新用户信息成功

**测试模块**: 用户管理模块  
**测试场景**: Admin 成功更新现有用户信息  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `updatableuser` 已存在，id 为 `user-uuid-123`

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/user/update`
3. 请求体包含用户 id 和更新的字段

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "id": "user-uuid-123",
  "username": "updatedusername",
  "password": "newpassword456",
  "user_type": "admin"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "user updated successfully",
  "data": null
}
```
- 用户信息已更新
- 使用新密码可以登录

---

### TC-USER-021: 只更新部分字段

**测试模块**: 用户管理模块  
**测试场景**: 只更新 username，不更新 password 和 user_type  
**优先级**: P1 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 目标用户已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 请求体只包含 id 和 username

**输入数据**:
```json
{
  "id": "user-uuid-123",
  "username": "newusername"
}
```

**预期结果**:
- HTTP 状态码: 200
- 用户名已更新
- 密码保持不变
- user_type 保持不变

---

### TC-USER-022: 更新不存在的用户

**测试模块**: 用户管理模块  
**测试场景**: 尝试更新不存在的用户  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 id `nonexistent-id` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 请求体使用不存在的用户 id

**输入数据**:
```json
{
  "id": "nonexistent-id",
  "username": "newname"
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

### TC-USER-023: 更新用户名为已存在的用户名

**测试模块**: 用户管理模块  
**测试场景**: 将用户名改为另一个已存在的用户名  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 A (id: user-a) 和用户 B (id: user-b) 都存在
- 用户 B 的用户名为 `existingname`

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 将用户 A 的 username 改为 `existingname`

**输入数据**:
```json
{
  "id": "user-a",
  "username": "existingname"
}
```

**预期结果**:
- HTTP 状态码: 409
- 响应 JSON:
```json
{
  "code": "user_exists",
  "message": "用户已存在"
}
```
- 用户 A 的用户名应保持不变

---

### TC-USER-024: Normal 用户尝试更新用户

**测试模块**: 用户管理模块  
**测试场景**: Normal 用户尝试更新其他用户信息  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/user/update`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "id": "some-user-id",
  "username": "hackedname"
}
```

**预期结果**:
- HTTP 状态码: 403
- 响应 JSON:
```json
{
  "code": "forbidden",
  "message": "权限不足"
}
```
- 用户信息不应被修改

---

### TC-USER-025: 更新用户缺少 id 字段

**测试模块**: 用户管理模块  
**测试场景**: 请求体缺少必填的 id  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 请求体不包含 id 字段

**输入数据**:
```json
{
  "username": "newname",
  "password": "newpass123"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-026: 更新密码少于 6 个字符

**测试模块**: 用户管理模块  
**测试场景**: 更新密码但新密码太短  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 目标用户已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 请求体 password 为 5 个字符

**输入数据**:
```json
{
  "id": "user-uuid",
  "password": "12345"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-027: 更新用户名长度不符合要求

**测试模块**: 用户管理模块  
**测试场景**: 更新 username 但长度不符合 3-32 字符要求  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 目标用户已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/update`
2. 请求体 username 为 2 个字符

**输入数据**:
```json
{
  "id": "user-uuid",
  "username": "ab"
}
```

**预期结果**:
- HTTP 状态码: 400
- 响应 JSON:
```json
{
  "code": "invalid_input",
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-028: Admin 删除用户成功

**测试模块**: 用户管理模块  
**测试场景**: Admin 成功删除普通用户  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `deletableuser` 已存在，id 为 `user-to-delete`
- 该用户不是默认管理员

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/user/delete`
3. 请求体包含用户 id

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "id": "user-to-delete"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "user deleted successfully",
  "data": null
}
```
- 该用户无法再登录
- 用户列表中不再包含该用户

---

### TC-USER-029: 删除不存在的用户

**测试模块**: 用户管理模块  
**测试场景**: 尝试删除不存在的用户  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 id `nonexistent-id` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/user/delete`
2. 请求体使用不存在的用户 id

**输入数据**:
```json
{
  "id": "nonexistent-id"
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

### TC-USER-030: Normal 用户尝试删除用户

**测试模块**: 用户管理模块  
**测试场景**: Normal 用户尝试删除其他用户  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/user/delete`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "id": "some-user-id"
}
```

**预期结果**:
- HTTP 状态码: 403
- 响应 JSON:
```json
{
  "code": "forbidden",
  "message": "权限不足"
}
```
- 用户不应被删除

---

### TC-USER-031: 删除用户缺少 id 字段

**测试模块**: 用户管理模块  
**测试场景**: 请求体缺少必填的 id  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/delete`
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
  "message": "请求参数不符合约束"
}
```

---

### TC-USER-032: 删除默认管理员

**测试模块**: 用户管理模块  
**测试场景**: 尝试删除系统默认管理员  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 存在默认管理员账号（如初始化创建的 admin）

**测试步骤**:
1. 发送 POST 请求到 `/api/user/delete`
2. 请求体 id 为默认管理员 id

**输入数据**:
```json
{
  "id": "default-admin-id"
}
```

**预期结果**:
- HTTP 状态码: 403
- 响应 JSON:
```json
{
  "code": "forbidden",
  "message": "不能删除默认管理员或权限不足"
}
```
- 默认管理员不应被删除

---

### TC-USER-033: 删除用户后验证

**测试模块**: 用户管理模块  
**测试场景**: 删除用户后验证该用户无法再登录  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `tobedeleted` 已存在，密码为 `password123`

**测试步骤**:
1. Admin 删除用户 `tobedeleted`
2. 尝试使用被删除用户的凭据登录

**输入数据**:
删除请求：
```json
{
  "id": "tobedeleted-user-id"
}
```

登录请求：
```json
{
  "username": "tobedeleted",
  "password": "password123"
}
```

**预期结果**:
- 删除成功
- 登录返回 401 (invalid_credentials)

---

### TC-USER-034: 删除用户后立即创建同名用户

**测试模块**: 用户管理模块  
**测试场景**: 删除用户后立即创建同名用户  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `recycleuser` 已存在

**测试步骤**:
1. Admin 删除用户 `recycleuser`
2. 立即创建同名用户 `recycleuser`

**输入数据**:
删除请求：
```json
{
  "id": "recycleuser-id"
}
```

创建请求：
```json
{
  "username": "recycleuser",
  "password": "newpassword123",
  "user_type": "normal"
}
```

**预期结果**:
- 删除成功
- 创建同名用户成功
- 新用户有全新的 id
- 与原用户完全独立

---

### TC-USER-035: 并发创建相同用户名

**测试模块**: 用户管理模块  
**测试场景**: 并发创建同名用户  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户名 `concurrentuser` 不存在

**测试步骤**:
1. 同时发送 5 个创建同名用户的请求

**输入数据**:
```json
{
  "username": "concurrentuser",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- 只有 1 个请求成功（返回 200）
- 其他 4 个请求返回 409 (user_exists)
- 系统中只创建一个 `concurrentuser` 用户
- 无数据不一致问题

---

### TC-USER-036: 并发更新同一用户

**测试模块**: 用户管理模块  
**测试场景**: 并发更新同一用户不同字段  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `concurrentupdate` 已存在

**测试步骤**:
1. 同时发送 3 个更新请求，分别更新 username、password、user_type

**输入数据**:
```json
// 请求 1
{
  "id": "concurrentupdate-id",
  "username": "newname1"
}

// 请求 2
{
  "id": "concurrentupdate-id",
  "password": "newpass1"
}

// 请求 3
{
  "id": "concurrentupdate-id",
  "user_type": "admin"
}
```

**预期结果**:
- 所有请求都应成功（200）或只有最后一个成功
- 最终用户状态应为最后一次更新的值
- 无数据损坏或竞态条件

---

### TC-USER-037: 并发删除同一用户

**测试模块**: 用户管理模块  
**测试场景**: 并发删除同一用户  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户 `concurrentdelete` 已存在

**测试步骤**:
1. 同时发送 3 个删除同一用户的请求

**输入数据**:
```json
{
  "id": "concurrentdelete-id"
}
```

**预期结果**:
- 第 1 个请求成功（200）
- 其他请求返回 404 (user_not_found) 或 200
- 用户被删除（只删除一次）
- 无错误日志或异常

---

### TC-USER-038: SQL 注入测试 - 创建用户

**测试模块**: 用户管理模块  
**测试场景**: 尝试在 username 中注入 SQL  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. username 包含 SQL 注入代码

**输入数据**:
```json
{
  "username": "user'; DROP TABLE users; --",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400 或 200
- 用户被创建（username 为字面量字符串）
- 不应执行任何 SQL 命令
- 不应影响其他用户数据

---

### TC-USER-039: XSS 攻击测试 - 创建用户

**测试模块**: 用户管理模块  
**测试场景**: 尝试在 username 中注入 XSS  
**优先级**: P1 (高)  
**测试类型**: 安全测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/user/create`
2. username 包含 XSS 代码

**输入数据**:
```json
{
  "username": "<script>alert('xss')</script>",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- HTTP 状态码: 400 或 200
- 用户被创建或创建失败
- 系统应对输入进行过滤或转义
- 不应执行 JavaScript 代码

---

### TC-USER-040: 用户列表包含用户数量验证

**测试模块**: 用户管理模块  
**测试场景**: 验证用户列表返回正确数量  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 系统中已知有 N 个用户

**测试步骤**:
1. 发送 POST 请求到 `/api/user/list`
2. 统计返回的用户数量

**输入数据**: 无

**预期结果**:
- HTTP 状态码: 200
- 返回的用户数组长度等于系统中的实际用户数
- 每个用户对象包含正确的字段

---

### TC-USER-041: 创建用户后时间戳验证

**测试模块**: 用户管理模块  
**测试场景**: 验证创建用户后的时间戳正确  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 记录当前时间戳
2. 创建新用户
3. 检查返回的 created_at 和 updated_at

**输入数据**:
```json
{
  "username": "timetestuser",
  "password": "password123",
  "user_type": "normal"
}
```

**预期结果**:
- created_at 和 updated_at 为当前时间附近（±5 秒）
- created_at 等于 updated_at
- 时间戳为 Unix 时间戳（秒级）

---

### TC-USER-042: 更新用户后时间戳验证

**测试模块**: 用户管理模块  
**测试场景**: 验证更新用户后 updated_at 更新  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 用户已存在

**测试步骤**:
1. 获取用户当前信息（记录 updated_at）
2. 等待 2 秒
3. 更新用户信息
4. 检查新的 updated_at

**输入数据**:
```json
{
  "id": "user-uuid",
  "username": "updatedname"
}
```

**预期结果**:
- 更新后的 updated_at > 更新前的 updated_at
- created_at 保持不变
- 时间戳为 Unix 时间戳（秒级）

---

## 测试数据准备

### Admin 测试账户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| admin | admin123 | admin | 管理员操作测试 |

### Normal 测试账户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| normaluser | userpass123 | normal | 权限测试 |

### 用于操作的用户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| updatableuser | pass123 | normal | 更新测试 |
| deletableuser | pass123 | normal | 删除测试 |
| tobedeleted | pass123 | normal | 删除后验证 |
| existinguser | pass123 | normal | 重复用户名测试 |

---

## 依赖和前置条件

1. Dancer 服务已启动并监听 8080 端口
2. etcd 服务正常运行
3. 默认管理员账号已创建
4. 测试数据库中预置了测试用户
5. 网络连接正常

---

## 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 误删生产环境用户 | 中 | 严重 | 测试环境与生产环境隔离 |
| 删除所有管理员 | 低 | 严重 | 禁止删除最后一个管理员 |
| 并发操作导致数据不一致 | 中 | 中 | 并发测试覆盖 |
| 权限绕过 | 低 | 严重 | 权限测试覆盖所有接口 |

---

## 测试通过标准

- 所有 P0 和 P1 测试用例通过
- 权限控制正确，Normal 用户无法执行 Admin 操作
- 无法删除默认管理员或最后一个管理员
- 并发操作不产生数据不一致
- 安全测试（SQL注入、XSS）正确处理
