# Zone 管理模块测试用例（Admin）

## 测试概述

测试 Dancer DNS 系统的 Zone 管理功能（仅 Admin 权限），包括 Zone 列表查询、详情获取、创建、更新和删除功能。

---

## 测试用例清单

### TC-ZONE-001: Admin 列出所有 Zone

**测试模块**: Zone 管理模块  
**测试场景**: Admin 用户成功获取所有 Zone 列表  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动并运行
- Admin 用户已登录并获取有效 Token
- 系统中存在多个 Zone

**测试步骤**:
1. Admin 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/list`
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
    "zones": [
      {
        "zone": "example.com",
        "record_count": 5,
        "created_at": 1704067200,
        "updated_at": 1704067200
      },
      {
        "zone": "test.com",
        "record_count": 3,
        "created_at": 1704067200,
        "updated_at": 1704067200
      }
    ]
  }
}
```
- 返回所有 Zone 的列表
- record_count 显示该 Zone 下的 Domain 数量

---

### TC-ZONE-002: Normal 用户访问 Zone 列表

**测试模块**: Zone 管理模块  
**测试场景**: Normal 用户尝试获取 Zone 列表  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录并获取有效 Token

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/list`
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
- Normal 用户无法查看 Zone 列表

---

### TC-ZONE-003: 未登录访问 Zone 列表

**测试模块**: Zone 管理模块  
**测试场景**: 未携带 Token 访问 Zone 列表  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/list`
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

### TC-ZONE-004: Admin 获取 Zone 详情成功

**测试模块**: Zone 管理模块  
**测试场景**: Admin 成功获取指定 Zone 的详细信息  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `example.com` 已存在

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/get`
3. 请求体包含 zone 名称

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "zone": "example.com"
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
    "zone": {
      "zone": "example.com",
      "record_count": 5,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```

---

### TC-ZONE-005: 获取不存在的 Zone 详情

**测试模块**: Zone 管理模块  
**测试场景**: 尝试获取不存在的 Zone  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/get`
2. 请求体使用不存在的 zone 名称

**输入数据**:
```json
{
  "zone": "nonexistent.com"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "zone_not_found",
  "message": "Zone 不存在"
}
```

---

### TC-ZONE-006: Normal 用户获取 Zone 详情

**测试模块**: Zone 管理模块  
**测试场景**: Normal 用户尝试获取 Zone 详情  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录
- Zone 已存在

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/get`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "zone": "example.com"
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

---

### TC-ZONE-007: 获取 Zone 详情缺少 zone 参数

**测试模块**: Zone 管理模块  
**测试场景**: 请求体缺少 zone 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/get`
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

### TC-ZONE-008: Admin 创建 Zone 成功

**测试模块**: Zone 管理模块  
**测试场景**: Admin 成功创建新的 Zone  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `newzone.com` 不存在

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/create`
3. 请求体包含有效的 zone 名称

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "zone": "newzone.com"
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
    "zone": {
      "zone": "newzone.com",
      "record_count": 0,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```
- Zone 已创建，record_count 初始为 0
- 可以在 Zone 列表中看到新创建的 Zone

---

### TC-ZONE-009: 创建已存在的 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 尝试创建已存在的 Zone  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `existing.com` 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体使用已存在的 zone 名称

**输入数据**:
```json
{
  "zone": "existing.com"
}
```

**预期结果**:
- HTTP 状态码: 409
- 响应 JSON:
```json
{
  "code": "zone_exists",
  "message": "Zone 已存在"
}
```
- 不应创建重复的 Zone

---

### TC-ZONE-010: Normal 用户尝试创建 Zone

**测试模块**: Zone 管理模块  
**测试场景**: Normal 用户尝试创建 Zone  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/create`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "zone": "hackedzone.com"
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
- Zone 不应被创建

---

### TC-ZONE-011: 创建 Zone 缺少 zone 参数

**测试模块**: Zone 管理模块  
**测试场景**: 请求体缺少 zone 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
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

### TC-ZONE-012: 创建无效的 Zone 格式

**测试模块**: Zone 管理模块  
**测试场景**: zone 名称不是有效的 FQDN 格式  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体包含无效的 zone 格式

**输入数据**:
```json
{
  "zone": "not a valid domain"
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

### TC-ZONE-013: 创建带协议的 Zone

**测试模块**: Zone 管理模块  
**测试场景**: zone 包含 http:// 或 https:// 前缀  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体包含带协议的 zone

**输入数据**:
```json
{
  "zone": "https://example.com"
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

### TC-ZONE-014: 创建带路径的 Zone

**测试模块**: Zone 管理模块  
**测试场景**: zone 包含路径  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体包含带路径的 zone

**输入数据**:
```json
{
  "zone": "example.com/path"
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

### TC-ZONE-015: 创建空字符串 Zone

**测试模块**: Zone 管理模块  
**测试场景**: zone 为空字符串  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体 zone 为空字符串

**输入数据**:
```json
{
  "zone": ""
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

### TC-ZONE-016: 创建带空格 Zone

**测试模块**: Zone 管理模块  
**测试场景**: zone 包含空格  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体 zone 包含空格

**输入数据**:
```json
{
  "zone": "example zone.com"
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

### TC-ZONE-017: 创建多级子域名 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 创建多级子域名作为 Zone  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体 zone 为三级域名

**输入数据**:
```json
{
  "zone": "sub.example.com"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持多级子域名
    - HTTP 状态码: 200
    - Zone 创建成功
  - **情况 B**: 系统仅支持二级域名
    - HTTP 状态码: 400
    - 返回格式错误
- 需要确认系统设计规范

---

### TC-ZONE-018: 创建国际化域名（IDN）

**测试模块**: Zone 管理模块  
**测试场景**: 创建国际化域名 Zone  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/create`
2. 请求体 zone 为中文域名

**输入数据**:
```json
{
  "zone": "例子.测试"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持 IDN
    - HTTP 状态码: 200
    - Zone 创建成功
  - **情况 B**: 系统不支持 IDN
    - HTTP 状态码: 400
    - 返回格式错误

---

### TC-ZONE-019: Admin 更新 Zone 成功

**测试模块**: Zone 管理模块  
**测试场景**: Admin 成功更新 Zone（更新统计信息）  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `updatezone.com` 已存在

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/update`
3. 请求体包含 zone 名称

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "zone": "updatezone.com"
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
    "zone": {
      "zone": "updatezone.com",
      "record_count": 5,
      "created_at": 1704067200,
      "updated_at": 1704153600
    }
  }
}
```
- updated_at 时间戳已更新

---

### TC-ZONE-020: 更新不存在的 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 尝试更新不存在的 Zone  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/update`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "zone_not_found",
  "message": "Zone 不存在"
}
```

---

### TC-ZONE-021: Normal 用户尝试更新 Zone

**测试模块**: Zone 管理模块  
**测试场景**: Normal 用户尝试更新 Zone  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录
- Zone 已存在

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/update`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "zone": "example.com"
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

---

### TC-ZONE-022: 更新 Zone 缺少 zone 参数

**测试模块**: Zone 管理模块  
**测试场景**: 请求体缺少 zone 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/update`
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

### TC-ZONE-023: Admin 删除 Zone 成功

**测试模块**: Zone 管理模块  
**测试场景**: Admin 成功删除 Zone 及其所有 Domain  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `deleteme.com` 已存在
- `deleteme.com` 下有多个 Domain

**测试步骤**:
1. Admin 登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/delete`
3. 请求体包含 zone 名称
4. 验证 Zone 及其 Domain 都被删除

**输入数据**:
```http
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "zone": "deleteme.com"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "Zone deleted successfully",
  "data": null
}
```
- Zone 列表中不再包含 `deleteme.com`
- 该 Zone 下的所有 Domain 都被删除
- CoreDNS 中该 Zone 的所有记录都被清理

---

### TC-ZONE-024: 删除不存在的 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 尝试删除不存在的 Zone  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/delete`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "zone_not_found",
  "message": "Zone 不存在"
}
```

---

### TC-ZONE-025: Normal 用户尝试删除 Zone

**测试模块**: Zone 管理模块  
**测试场景**: Normal 用户尝试删除 Zone  
**优先级**: P1 (高)  
**测试类型**: 权限测试

**前置条件**:
- 系统已启动
- Normal 用户已登录
- Zone 已存在

**测试步骤**:
1. Normal 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/zones/delete`

**输入数据**:
```http
Authorization: Bearer <normal_token>
{
  "zone": "example.com"
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
- Zone 不应被删除

---

### TC-ZONE-026: 删除 Zone 缺少 zone 参数

**测试模块**: Zone 管理模块  
**测试场景**: 请求体缺少 zone 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/zones/delete`
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

### TC-ZONE-027: 删除 Zone 后创建同名 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 删除 Zone 后立即创建同名 Zone  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `recycle.com` 已存在

**测试步骤**:
1. Admin 删除 Zone `recycle.com`
2. 立即创建同名 Zone `recycle.com`
3. 在新 Zone 下创建 Domain

**输入数据**:
删除请求：
```json
{
  "zone": "recycle.com"
}
```

创建请求：
```json
{
  "zone": "recycle.com"
}
```

**预期结果**:
- 删除成功
- 创建同名 Zone 成功
- record_count 初始化为 0
- 新 Zone 与原 Zone 完全独立

---

### TC-ZONE-028: 创建 Zone 后时间戳验证

**测试模块**: Zone 管理模块  
**测试场景**: 验证创建 Zone 后的时间戳  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 记录当前时间戳
2. 创建新 Zone
3. 检查返回的 created_at 和 updated_at

**输入数据**:
```json
{
  "zone": "timetest.com"
}
```

**预期结果**:
- created_at 和 updated_at 为当前时间附近（±5 秒）
- created_at 等于 updated_at
- 时间戳为 Unix 时间戳（秒级）

---

### TC-ZONE-029: 更新 Zone 后时间戳验证

**测试模块**: Zone 管理模块  
**测试场景**: 验证更新 Zone 后 updated_at 更新  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone 已存在

**测试步骤**:
1. 获取 Zone 当前信息（记录 updated_at）
2. 等待 2 秒
3. 更新 Zone
4. 检查新的 updated_at

**输入数据**:
```json
{
  "zone": "example.com"
}
```

**预期结果**:
- 更新后的 updated_at > 更新前的 updated_at
- created_at 保持不变
- 时间戳为 Unix 时间戳（秒级）

---

### TC-ZONE-030: Zone record_count 正确性验证

**测试模块**: Zone 管理模块  
**测试场景**: 验证 record_count 统计准确  
**优先级**: P1 (高)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `counttest.com` 已存在

**测试步骤**:
1. 获取 Zone 当前 record_count
2. 在该 Zone 下创建 3 个 Domain
3. 再次获取 Zone 信息，验证 record_count
4. 删除 1 个 Domain
5. 再次验证 record_count

**输入数据**: 创建 Domain 请求

**预期结果**:
- 初始 record_count 为 N
- 创建 3 个 Domain 后，record_count = N + 3
- 删除 1 个 Domain 后，record_count = N + 2
- record_count 始终与实际 Domain 数量一致

---

### TC-ZONE-031: 并发创建相同 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 并发创建同名 Zone  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `concurrentzone.com` 不存在

**测试步骤**:
1. 同时发送 5 个创建同名 Zone 的请求

**输入数据**:
```json
{
  "zone": "concurrentzone.com"
}
```

**预期结果**:
- 只有 1 个请求成功（返回 200）
- 其他 4 个请求返回 409 (zone_exists)
- 系统中只创建一个 `concurrentzone.com`
- 无数据不一致问题

---

### TC-ZONE-032: 并发删除同一 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 并发删除同一 Zone  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `concurrentdelete.com` 已存在

**测试步骤**:
1. 同时发送 3 个删除同一 Zone 的请求

**输入数据**:
```json
{
  "zone": "concurrentdelete.com"
}
```

**预期结果**:
- 第 1 个请求成功（200）
- 其他请求返回 404 (zone_not_found) 或 200
- Zone 被删除（只删除一次）
- 无错误日志或异常

---

### TC-ZONE-033: 大写 Zone 名称处理

**测试模块**: Zone 管理模块  
**测试场景**: 使用大写字母创建 Zone  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建 Zone `EXAMPLE.COM`
2. 尝试创建 Zone `example.com`

**输入数据**:
```json
{
  "zone": "EXAMPLE.COM"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 大小写敏感
    - 两个 Zone 都创建成功（视为不同 Zone）
  - **情况 B**: 大小写不敏感
    - 第二个请求返回 409 (zone_exists)
- 需要确认系统设计规范

---

### TC-ZONE-034: 大小写混合 Zone 名称

**测试模块**: Zone 管理模块  
**测试场景**: 使用大小写混合创建 Zone  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建 Zone `Example.COM`

**输入数据**:
```json
{
  "zone": "Example.COM"
}
```

**预期结果**:
- HTTP 状态码: 200
- Zone 创建成功
- 查询时使用不同大小写应能获取（根据大小写敏感策略）

---

### TC-ZONE-035: Zone 名称前后空格处理

**测试模块**: Zone 管理模块  
**测试场景**: Zone 名称包含前后空格  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建 Zone `  example.com  `

**输入数据**:
```json
{
  "zone": "  example.com  "
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 自动 trim
    - HTTP 状态码: 200
    - Zone `example.com` 被创建
  - **情况 B**: 验证失败
    - HTTP 状态码: 400
    - 返回格式错误

---

### TC-ZONE-036: 超长的 Zone 名称

**测试模块**: Zone 管理模块  
**测试场景**: Zone 名称超过 253 字符限制  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建超长 Zone

**输入数据**:
```json
{
  "zone": "a.".repeat(100) + "com"
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

### TC-ZONE-037: 无效的 TLD 格式

**测试模块**: Zone 管理模块  
**测试场景**: Zone 使用无效的 TLD  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建 Zone `example.123`（数字 TLD）

**输入数据**:
```json
{
  "zone": "example.123"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 验证严格
    - HTTP 状态码: 400
  - **情况 B**: 验证宽松
    - HTTP 状态码: 200
- DNS 系统通常允许任意 TLD，但需要确认业务规则

---

### TC-ZONE-038: 单级域名 Zone

**测试模块**: Zone 管理模块  
**测试场景**: 尝试创建单级域名（无点号）  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- Admin 已登录

**测试步骤**:
1. 发送 POST 请求创建 Zone `localhost`

**输入数据**:
```json
{
  "zone": "localhost"
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 允许单级域名
    - HTTP 状态码: 200
  - **情况 B**: 要求至少二级域名
    - HTTP 状态码: 400

---

### TC-ZONE-039: 创建 Zone 后 etcd Key 验证

**测试模块**: Zone 管理模块  
**测试场景**: 验证 Zone 在 etcd 中的存储格式  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- Admin 已登录
- 有 etcd 直接访问权限

**测试步骤**:
1. 创建 Zone `testzone.com`
2. 直接查询 etcd 验证 Key 格式

**输入数据**:
```json
{
  "zone": "testzone.com"
}
```

**预期结果**:
- etcd 中存在 Key: `/dancer/zones/testzone.com`
- Value 包含 Zone 的完整信息
- 格式正确，可解析

---

### TC-ZONE-040: 删除 Zone 后 etcd Key 验证

**测试模块**: Zone 管理模块  
**测试场景**: 验证删除 Zone 后 etcd 清理完整  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- Admin 已登录
- Zone `cleanme.com` 已存在，下有多个 Domain
- 有 etcd 直接访问权限

**测试步骤**:
1. 删除 Zone `cleanme.com`
2. 查询 etcd 验证所有相关 Key 都被删除

**输入数据**:
```json
{
  "zone": "cleanme.com"
}
```

**预期结果**:
- `/dancer/zones/cleanme.com` 不存在
- `/dancer/domains/cleanme.com/` 下的所有 Key 不存在
- CoreDNS 前缀下的相关记录不存在
- etcd 中无残留数据

---

## 测试数据准备

### Admin 测试账户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| admin | admin123 | admin | Zone 管理测试 |

### Normal 测试账户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| normaluser | userpass123 | normal | 权限测试 |

### 测试用 Zone

| Zone 名称 | record_count | 用途 |
|-----------|--------------|------|
| example.com | 5 | 正常查询测试 |
| existing.com | 2 | 重复创建测试 |
| deleteme.com | 3 | 删除测试 |
| recycle.com | 1 | 删除后重建测试 |
| updatezone.com | 2 | 更新测试 |
| counttest.com | 0 | 计数测试 |

---

## 依赖和前置条件

1. Dancer 服务已启动并监听 8080 端口
2. etcd 服务正常运行
3. Admin 账号已创建
4. 测试数据库中预置了测试 Zone
5. 网络连接正常
6. 有方法直接查询 etcd 数据（用于集成测试）

---

## 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| 误删生产环境 Zone | 中 | 严重 | 测试环境与生产环境隔离 |
| 删除 Zone 后 Domain 残留 | 中 | 中 | 级联删除测试覆盖 |
| 并发操作导致数据不一致 | 中 | 中 | 并发测试覆盖 |
| CoreDNS 记录未清理 | 中 | 中 | etcd 直接验证 |

---

## 测试通过标准

- 所有 P0 和 P1 测试用例通过
- 权限控制正确，Normal 用户无法执行 Zone 管理操作
- Zone 删除时正确级联删除所有 Domain 和 CoreDNS 记录
- record_count 统计始终准确
- 并发操作不产生数据不一致
- etcd 中无残留数据
