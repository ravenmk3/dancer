# Domain 管理模块测试用例

## 测试概述

测试 Dancer DNS 系统的 Domain 管理功能（所有 JWT 认证用户），包括 Domain 列表查询、详情获取、创建、更新和删除功能。

---

## 测试用例清单

### TC-DOMAIN-001: 列出 Zone 下所有 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 成功获取指定 Zone 下的所有 Domain  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动并运行
- 用户已登录并获取有效 Token
- Zone `example.com` 存在且有多个 Domain

**测试步骤**:
1. 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/domains/list`
3. 请求体包含 zone 名称

**输入数据**:
```http
Authorization: Bearer <valid_token>
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
- 返回该 Zone 下所有 Domain 列表
- Domain 包含完整的 DNS 记录信息

---

### TC-DOMAIN-002: 列出不存在 Zone 的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 尝试获取不存在 Zone 下的 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/list`
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

### TC-DOMAIN-003: 未登录访问 Domain 列表

**测试模块**: Domain 管理模块  
**测试场景**: 未携带 Token 访问 Domain 列表  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/list`
2. 不携带 Authorization Header

**输入数据**:
```json
{
  "zone": "example.com"
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

### TC-DOMAIN-004: 列出 Domain 缺少 zone 参数

**测试模块**: Domain 管理模块  
**测试场景**: 请求体缺少 zone 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/list`
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

### TC-DOMAIN-005: 获取 Domain 详情成功

**测试模块**: Domain 管理模块  
**测试场景**: 成功获取指定 Domain 的详细信息  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `www` 存在于该 Zone

**测试步骤**:
1. 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/domains/get`
3. 请求体包含 zone 和 domain

**输入数据**:
```http
Authorization: Bearer <valid_token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www"
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

---

### TC-DOMAIN-006: 获取不存在的 Domain 详情

**测试模块**: Domain 管理模块  
**测试场景**: 尝试获取不存在的 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `nonexistent` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/get`
2. 请求体使用不存在的 domain

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "nonexistent"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "domain_not_found",
  "message": "Domain 不存在"
}
```

---

### TC-DOMAIN-007: 获取不存在的 Zone 下的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: Zone 不存在时获取 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/get`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com",
  "domain": "www"
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

### TC-DOMAIN-008: 获取 Domain 详情缺少参数

**测试模块**: Domain 管理模块  
**测试场景**: 请求体缺少 zone 或 domain 字段  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/get`
2. 请求体只包含 zone

**输入数据**:
```json
{
  "zone": "example.com"
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

### TC-DOMAIN-009: 创建 Domain 成功

**测试模块**: Domain 管理模块  
**测试场景**: 成功创建新的 Domain  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `newdomain` 不存在于该 Zone

**测试步骤**:
1. 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/domains/create`
3. 请求体包含完整的 Domain 信息

**输入数据**:
```http
Authorization: Bearer <valid_token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "newdomain",
  "ips": ["192.168.1.100", "192.168.1.101"],
  "ttl": 300
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
    "domain": {
      "zone": "example.com",
      "domain": "newdomain",
      "name": "newdomain.example.com",
      "ips": ["192.168.1.100", "192.168.1.101"],
      "ttl": 300,
      "record_count": 2,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```
- Domain 创建成功
- CoreDNS 中自动创建相应的 DNS 记录
- Zone 的 record_count 增加 1

---

### TC-DOMAIN-010: 创建已存在的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 尝试创建已存在的 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `www` 已存在于该 Zone

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体使用已存在的 domain

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.200"],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 409
- 响应 JSON:
```json
{
  "code": "domain_exists",
  "message": "Domain 已存在"
}
```
- 不应创建重复的 Domain

---

### TC-DOMAIN-011: 在不存在的 Zone 下创建 Domain

**测试模块**: Domain 管理模块  
**测试场景**: Zone 不存在时创建 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com",
  "domain": "www",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "zone_not_found",
  "message": "Zone 不存在，需要先创建 Zone"
}
```
- 不应创建 Domain

---

### TC-DOMAIN-012: 创建 Domain 缺少必填字段

**测试模块**: Domain 管理模块  
**测试场景**: 请求体缺少 zone、domain、ips 或 ttl  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体缺少 ips 字段

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test",
  "ttl": 300
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

### TC-DOMAIN-013: 创建 Domain 使用无效的 IP 地址

**测试模块**: Domain 管理模块  
**测试场景**: ips 数组中包含无效的 IP 格式  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ips 包含无效格式

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test",
  "ips": ["192.168.1", "invalid", "256.1.1.1"],
  "ttl": 300
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

### TC-DOMAIN-014: 创建 Domain 使用无效的 TTL

**测试模块**: Domain 管理模块  
**测试场景**: ttl 值小于 1  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ttl 为 0

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test",
  "ips": ["192.168.1.100"],
  "ttl": 0
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

### TC-DOMAIN-015: 创建 Domain 使用负数 TTL

**测试模块**: Domain 管理模块  
**测试场景**: ttl 为负数  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ttl 为 -1

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test",
  "ips": ["192.168.1.100"],
  "ttl": -1
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

### TC-DOMAIN-016: 创建 Domain 使用空 IP 数组

**测试模块**: Domain 管理模块  
**测试场景**: ips 为空数组  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ips 为空数组

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test",
  "ips": [],
  "ttl": 300
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

### TC-DOMAIN-017: 创建 Domain 使用空字符串 domain

**测试模块**: Domain 管理模块  
**测试场景**: domain 字段为空字符串  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 为空字符串

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "",
  "ips": ["192.168.1.100"],
  "ttl": 300
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

### TC-DOMAIN-018: 创建 Domain 使用 @ 作为根域名

**测试模块**: Domain 管理模块  
**测试场景**: 使用 @ 符号代表 Zone 根域名  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- 根域名 `@` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 为 "@"

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "@",
  "ips": ["192.168.1.10"],
  "ttl": 600
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
    "domain": {
      "zone": "example.com",
      "domain": "@",
      "name": "example.com",
      "ips": ["192.168.1.10"],
      "ttl": 600,
      "record_count": 1,
      "created_at": 1704067200,
      "updated_at": 1704067200
    }
  }
}
```
- name 字段应为 `example.com`（Zone 名本身）
- CoreDNS 记录正确创建

---

### TC-DOMAIN-019: 创建 Domain 使用多级子域名

**测试模块**: Domain 管理模块  
**测试场景**: 创建多级子域名  
**优先级**: P1 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 为多级子域名

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "sub.www",
  "ips": ["192.168.1.50"],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 200
- Domain 创建成功
- name 字段为 `sub.www.example.com`
- CoreDNS 记录路径正确（含点的子域名）

---

### TC-DOMAIN-020: 创建 Domain 使用通配符 *

**测试模块**: Domain 管理模块  
**测试场景**: 创建通配符 Domain  
**优先级**: P1 (高)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 为 "*"

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "*",
  "ips": ["192.168.1.200"],
  "ttl": 300
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持通配符
    - HTTP 状态码: 200
    - Domain 创建成功
    - CoreDNS 支持通配符解析
  - **情况 B**: 系统不支持通配符
    - HTTP 状态码: 400
    - 返回格式错误

---

### TC-DOMAIN-021: 创建 Domain 使用特殊字符

**测试模块**: Domain 管理模块  
**测试场景**: domain 包含 DNS 不允许的特殊字符  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 包含特殊字符

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "test@#$%",
  "ips": ["192.168.1.100"],
  "ttl": 300
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
- DNS 标签只允许字母、数字和连字符

---

### TC-DOMAIN-022: 创建 Domain 使用超长域名标签

**测试模块**: Domain 管理模块  
**测试场景**: domain 标签超过 63 字符限制  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 domain 为 64 个字符

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "a".repeat(64),
  "ips": ["192.168.1.100"],
  "ttl": 300
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

### TC-DOMAIN-023: 创建 Domain 使用 IPv6 地址

**测试模块**: Domain 管理模块  
**测试场景**: ips 数组包含 IPv6 地址  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ips 包含 IPv6 地址

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "ipv6test",
  "ips": ["2001:db8::1", "192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 系统支持 IPv6
    - HTTP 状态码: 200
    - Domain 创建成功
  - **情况 B**: 系统仅支持 IPv4
    - HTTP 状态码: 400
    - 返回格式错误

---

### TC-DOMAIN-024: 创建大量 IP 的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: ips 数组包含大量 IP 地址  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/create`
2. 请求体 ips 包含 50 个 IP 地址

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "manyips",
  "ips": ["192.168.1."+i for i in range(1, 51)],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 200
- Domain 创建成功
- record_count 为 50
- CoreDNS 中创建 50 条记录
- 系统性能无显著下降

---

### TC-DOMAIN-025: 更新 Domain 成功

**测试模块**: Domain 管理模块  
**测试场景**: 成功更新 Domain 的 IP 列表和 TTL  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `www` 存在，当前 IPs: ["192.168.1.1", "192.168.1.2"]

**测试步骤**:
1. 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/domains/update`
3. 请求体包含新的 IP 列表

**输入数据**:
```http
Authorization: Bearer <valid_token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.3", "192.168.1.4", "192.168.1.5"],
  "ttl": 600
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
    "domain": {
      "zone": "example.com",
      "domain": "www",
      "name": "www.example.com",
      "ips": ["192.168.1.3", "192.168.1.4", "192.168.1.5"],
      "ttl": 600,
      "record_count": 3,
      "created_at": 1704067200,
      "updated_at": 1704153600
    }
  }
}
```
- IP 列表完全替换为新列表
- record_count 更新为 3
- CoreDNS 记录同步更新（删除旧 IP，添加新 IP）
- updated_at 时间戳更新

---

### TC-DOMAIN-026: 更新 Domain 不更新 TTL

**测试模块**: Domain 管理模块  
**测试场景**: 只更新 IPs，不更新 TTL  
**优先级**: P1 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 和 Domain 已存在
- Domain 当前 TTL 为 300

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体不包含 ttl 字段

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.100"]
}
```

**预期结果**:
- HTTP 状态码: 200
- IP 列表已更新
- TTL 保持原值 300
- updated_at 已更新

---

### TC-DOMAIN-027: 更新不存在的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 尝试更新不存在的 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `nonexistent` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体使用不存在的 domain

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "nonexistent",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "domain_not_found",
  "message": "Domain 不存在"
}
```

---

### TC-DOMAIN-028: 在不存在的 Zone 下更新 Domain

**测试模块**: Domain 管理模块  
**测试场景**: Zone 不存在时更新 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com",
  "domain": "www",
  "ips": ["192.168.1.100"],
  "ttl": 300
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

### TC-DOMAIN-029: 更新 Domain 缺少必填字段

**测试模块**: Domain 管理模块  
**测试场景**: 请求体缺少 zone、domain 或 ips  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体缺少 ips

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www"
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

### TC-DOMAIN-030: 更新 Domain 使用无效的 IP

**测试模块**: Domain 管理模块  
**测试场景**: 更新时使用无效的 IP 格式  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 和 Domain 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体 ips 包含无效 IP

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["invalid.ip"],
  "ttl": 300
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
- Domain 保持原值不变

---

### TC-DOMAIN-031: 更新 Domain 清空所有 IP

**测试模块**: Domain 管理模块  
**测试场景**: 更新时传入空 IP 数组  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 和 Domain 已存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/update`
2. 请求体 ips 为空数组

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": [],
  "ttl": 300
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
- Domain 保持原值不变

---

### TC-DOMAIN-032: 删除 Domain 成功

**测试模块**: Domain 管理模块  
**测试场景**: 成功删除 Domain 及其 CoreDNS 记录  
**优先级**: P0 (高)  
**测试类型**: 正向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `deleteme` 存在且有多个 IP

**测试步骤**:
1. 用户登录获取 Token
2. 发送 POST 请求到 `/api/dns/domains/delete`
3. 请求体包含 zone 和 domain
4. 验证 Domain 和 CoreDNS 记录都被删除

**输入数据**:
```http
Authorization: Bearer <valid_token>
Content-Type: application/json

{
  "zone": "example.com",
  "domain": "deleteme"
}
```

**预期结果**:
- HTTP 状态码: 200
- 响应 JSON:
```json
{
  "code": "success",
  "message": "Domain deleted successfully",
  "data": null
}
```
- Domain 列表中不再包含 `deleteme`
- CoreDNS 中该 Domain 的所有记录都被删除
- Zone 的 record_count 减少 1

---

### TC-DOMAIN-033: 删除不存在的 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 尝试删除不存在的 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `nonexistent` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/delete`
2. 请求体使用不存在的 domain

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "nonexistent"
}
```

**预期结果**:
- HTTP 状态码: 404
- 响应 JSON:
```json
{
  "code": "domain_not_found",
  "message": "Domain 不存在"
}
```

---

### TC-DOMAIN-034: 在不存在的 Zone 下删除 Domain

**测试模块**: Domain 管理模块  
**测试场景**: Zone 不存在时删除 Domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `nonexistent.com` 不存在

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/delete`
2. 请求体使用不存在的 zone

**输入数据**:
```json
{
  "zone": "nonexistent.com",
  "domain": "www"
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

### TC-DOMAIN-035: 删除 Domain 缺少必填字段

**测试模块**: Domain 管理模块  
**测试场景**: 请求体缺少 zone 或 domain  
**优先级**: P1 (高)  
**测试类型**: 负向测试

**前置条件**:
- 系统已启动
- 用户已登录

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/delete`
2. 请求体缺少 domain

**输入数据**:
```json
{
  "zone": "example.com"
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

### TC-DOMAIN-036: 删除 Domain 后创建同名 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 删除 Domain 后立即创建同名 Domain  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- Domain `recycle` 存在

**测试步骤**:
1. 删除 Domain `recycle`
2. 立即创建同名 Domain `recycle`，使用不同的 IP

**输入数据**:
删除请求：
```json
{
  "zone": "example.com",
  "domain": "recycle"
}
```

创建请求：
```json
{
  "zone": "example.com",
  "domain": "recycle",
  "ips": ["192.168.2.100"],
  "ttl": 300
}
```

**预期结果**:
- 删除成功
- 创建同名 Domain 成功
- 新 Domain 与旧 Domain 完全独立
- record_count 从 0 变为 1

---

### TC-DOMAIN-037: Zone 被删除后访问其 Domain

**测试模块**: Domain 管理模块  
**测试场景**: Zone 被删除后尝试访问其 Domain  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `tobedeleted.com` 存在且有 Domain
- 用户有权限删除 Zone（Admin）

**测试步骤**:
1. 获取 Zone `tobedeleted.com` 的 Domain 列表（记录成功）
2. Admin 删除 Zone `tobedeleted.com`
3. 再次尝试获取该 Zone 的 Domain 列表

**输入数据**:
```json
{
  "zone": "tobedeleted.com"
}
```

**预期结果**:
- 步骤 1: 返回 Domain 列表（200）
- 步骤 3: 返回 404 (zone_not_found)
- 级联删除正确执行

---

### TC-DOMAIN-038: 创建 Domain 后 CoreDNS 记录验证

**测试模块**: Domain 管理模块  
**测试场景**: 验证 Domain 创建后 CoreDNS 记录格式正确  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone `example.com` 存在
- 有 etcd 直接访问权限
- CoreDNS 前缀配置为 `/skydns`

**测试步骤**:
1. 创建 Domain `www`，IPs: ["192.168.1.1", "192.168.1.2"]
2. 直接查询 etcd 验证 CoreDNS 记录

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.1", "192.168.1.2"],
  "ttl": 300
}
```

**预期结果**:
- etcd 中存在 Keys:
  - `/skydns/com/example/www/x1` - 对应 IP 192.168.1.1
  - `/skydns/com/example/www/x2` - 对应 IP 192.168.1.2
- Value 格式符合 CoreDNS 规范
- TTL 值正确设置

---

### TC-DOMAIN-039: 更新 Domain 后 CoreDNS 记录同步

**测试模块**: Domain 管理模块  
**测试场景**: 验证更新 Domain 后 CoreDNS 记录同步更新  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- 用户已登录
- Domain `www.example.com` 存在，IPs: ["192.168.1.1", "192.168.1.2"]
- 有 etcd 直接访问权限

**测试步骤**:
1. 查询 etcd 记录当前状态（应有 x1, x2）
2. 更新 Domain，IPs: ["192.168.1.3"]
3. 再次查询 etcd 验证

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.3"],
  "ttl": 600
}
```

**预期结果**:
- 更新前：存在 x1 (192.168.1.1), x2 (192.168.1.2)
- 更新后：
  - 删除 x1, x2
  - 新建 x1 (192.168.1.3)
- 不存在 x2
- TTL 更新为 600

---

### TC-DOMAIN-040: 删除 Domain 后 CoreDNS 记录清理

**测试模块**: Domain 管理模块  
**测试场景**: 验证删除 Domain 后 CoreDNS 记录完全清理  
**优先级**: P1 (高)  
**测试类型**: 集成测试

**前置条件**:
- 系统已启动
- 用户已登录
- Domain `cleanme.example.com` 存在且有多个 IP
- 有 etcd 直接访问权限

**测试步骤**:
1. 查询 etcd 记录当前状态
2. 删除 Domain `cleanme`
3. 再次查询 etcd 验证

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "cleanme"
}
```

**预期结果**:
- 删除前：存在 `/skydns/com/example/cleanme/x1`, `x2`...
- 删除后：不存在 `/skydns/com/example/cleanme/` 下的任何 Key
- 清理完整，无残留记录

---

### TC-DOMAIN-041: 并发创建同一 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 并发创建同名 Domain  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- 多个用户已登录
- Zone `example.com` 存在
- Domain `concurrent` 不存在

**测试步骤**:
1. 同时发送 5 个创建同名 Domain 的请求

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "concurrent",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- 只有 1 个请求成功（返回 200）
- 其他请求返回 409 (domain_exists)
- 系统中只创建一个 Domain
- 无数据不一致问题

---

### TC-DOMAIN-042: 并发更新同一 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 并发更新同一 Domain  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- 多个用户已登录
- Domain `example.com` 存在

**测试步骤**:
1. 同时发送 3 个更新请求，使用不同的 IP

**输入数据**:
```json
// 请求 1
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.10"],
  "ttl": 300
}

// 请求 2
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.20"],
  "ttl": 300
}

// 请求 3
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.30"],
  "ttl": 300
}
```

**预期结果**:
- 所有请求都可能成功（200）
- 最终 Domain 的 IP 为最后一次成功的更新
- CoreDNS 记录与最终状态一致
- 无数据损坏

---

### TC-DOMAIN-043: 并发删除同一 Domain

**测试模块**: Domain 管理模块  
**测试场景**: 并发删除同一 Domain  
**优先级**: P2 (中)  
**测试类型**: 并发测试

**前置条件**:
- 系统已启动
- 多个用户已登录
- Domain `example.com` 存在

**测试步骤**:
1. 同时发送 3 个删除同一 Domain 的请求

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "concurrentdelete"
}
```

**预期结果**:
- 第 1 个请求成功（200）
- 其他请求返回 404 (domain_not_found) 或 200
- Domain 被删除（只删除一次）
- CoreDNS 记录被清理

---

### TC-DOMAIN-044: Domain 列表按名称排序

**测试模块**: Domain 管理模块  
**测试场景**: 验证 Domain 列表返回顺序  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 存在且有多个 Domain

**测试步骤**:
1. 发送 POST 请求到 `/api/dns/domains/list`
2. 检查返回列表的顺序

**输入数据**:
```json
{
  "zone": "example.com"
}
```

**预期结果**:
- HTTP 状态码: 200
- Domain 列表按 domain 名称排序（字母顺序）
- 或使用创建时间排序（根据设计）
- 顺序应一致且可预测

---

### TC-DOMAIN-045: 创建 Domain 后时间戳验证

**测试模块**: Domain 管理模块  
**测试场景**: 验证创建 Domain 后的时间戳  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 记录当前时间戳
2. 创建新 Domain
3. 检查返回的 created_at 和 updated_at

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "timetest",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- created_at 和 updated_at 为当前时间附近（±5 秒）
- created_at 等于 updated_at
- 时间戳为 Unix 时间戳（秒级）

---

### TC-DOMAIN-046: 更新 Domain 后时间戳验证

**测试模块**: Domain 管理模块  
**测试场景**: 验证更新 Domain 后 updated_at 更新  
**优先级**: P2 (中)  
**测试类型**: 功能测试

**前置条件**:
- 系统已启动
- 用户已登录
- Domain 已存在

**测试步骤**:
1. 获取 Domain 当前信息（记录 updated_at）
2. 等待 2 秒
3. 更新 Domain
4. 检查新的 updated_at

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "www",
  "ips": ["192.168.1.200"],
  "ttl": 300
}
```

**预期结果**:
- 更新后的 updated_at > 更新前的 updated_at
- created_at 保持不变
- 时间戳为 Unix 时间戳（秒级）

---

### TC-DOMAIN-047: 特殊 Domain 名称测试

**测试模块**: Domain 管理模块  
**测试场景**: 测试 DNS 特殊名称  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 尝试创建以下 Domain：
   - `_dmarc` (DMARC 记录)
   - `_domainkey` (DKIM)
   - `_spf` (SPF)
   - `_acme-challenge` (Let's Encrypt)

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "_dmarc",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- HTTP 状态码: 200
- Domain 创建成功（下划线在 DNS 中是允许的）
- CoreDNS 记录正确创建

---

### TC-DOMAIN-048: 大写 Domain 名称处理

**测试模块**: Domain 管理模块  
**测试场景**: 使用大写字母创建 Domain  
**优先级**: P2 (中)  
**测试类型**: 边界测试

**前置条件**:
- 系统已启动
- 用户已登录
- Zone 已存在

**测试步骤**:
1. 创建 Domain `WWW`
2. 尝试创建 Domain `www`

**输入数据**:
```json
{
  "zone": "example.com",
  "domain": "WWW",
  "ips": ["192.168.1.100"],
  "ttl": 300
}
```

**预期结果**:
- 以下两种情况之一：
  - **情况 A**: 大小写敏感
    - 两个 Domain 都创建成功
  - **情况 B**: 大小写不敏感
    - 第二个请求返回 409 (domain_exists)
- DNS 标准规定域名不区分大小写

---

## 测试数据准备

### 测试账户

| 用户名 | 密码 | 用户类型 | 用途 |
|--------|------|----------|------|
| admin | admin123 | admin | Admin 用户测试 |
| normaluser | userpass123 | normal | Normal 用户测试 |

### 测试用 Zone

| Zone 名称 | 状态 | 用途 |
|-----------|------|------|
| example.com | 存在，有 Domain | Domain 操作测试 |
| empty.com | 存在，无 Domain | 空列表测试 |
| nonexistent.com | 不存在 | 错误场景测试 |

### 测试用 Domain

| Zone | Domain | IPs | TTL | 用途 |
|------|--------|-----|-----|------|
| example.com | www | ["192.168.1.1", "192.168.1.2"] | 300 | 正常测试 |
| example.com | @ | ["192.168.1.10"] | 600 | 根域名测试 |
| example.com | deleteme | ["192.168.1.50"] | 300 | 删除测试 |
| example.com | recycle | ["192.168.1.60"] | 300 | 重建测试 |
| example.com | cleanme | ["192.168.1.70", "192.168.1.71"] | 300 | CoreDNS 清理测试 |

---

## 依赖和前置条件

1. Dancer 服务已启动并监听 8080 端口
2. etcd 服务正常运行
3. CoreDNS 配置正确，使用 etcd 作为后端
4. 测试数据库中预置了测试 Zone 和 Domain
5. 网络连接正常
6. 有方法直接查询 etcd 数据（用于集成测试）

---

## 风险评估

| 风险 | 可能性 | 影响 | 缓解措施 |
|------|--------|------|----------|
| CoreDNS 记录未同步 | 中 | 高 | CoreDNS 集成测试覆盖 |
| Zone 删除后 Domain 残留 | 低 | 高 | 级联删除验证 |
| IP 格式验证不严 | 中 | 中 | IP 格式边界测试 |
| 并发更新导致数据不一致 | 中 | 中 | 并发测试覆盖 |

---

## 测试通过标准

- 所有 P0 和 P1 测试用例通过
- Domain CRUD 操作正常工作
- CoreDNS 记录与 Domain 数据始终同步
- 更新 Domain 时 IP 列表完全替换
- 删除 Domain 后 CoreDNS 记录完全清理
- record_count 统计始终准确
- 并发操作不产生数据不一致
