# 小窝API测试用例设计

**版本**: v1.0  
**创建时间**: 2025-12-30  
**负责团队**: SRE & QA  

---

## 📋 测试覆盖范围

### 1. 房间管理测试
### 2. 成员操作测试  
### 3. 播放控制测试
### 4. WebSocket连接测试
### 5. 会话管理测试
### 6. 错误处理测试

---

## 1. 房间管理测试用例

### 1.1 创建房间测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_RM_001 | 创建公开房间成功 | P1 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. 提供name, media_url等必需参数<br>3. 设置is_private=false | 201 Created<br>返回房间信息和token |
| TC_RM_002 | 创建私密房间成功 | P1 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. 设置is_private=true<br>3. 提供password | 201 Created<br>返回房间信息和token |
| TC_RM_003 | 创建房间参数缺失 | P1 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. 缺少必需字段name或media_url | 400 Bad Request<br>错误信息指出缺失参数 |
| TC_RM_004 | 创建房间name长度超限 | P2 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. name字段超过100字符 | 400 Bad Request<br>错误信息指出字段限制 |
| TC_RM_005 | 创建房间max_users越界 | P2 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. max_users=0或max_users=1001 | 400 Bad Request<br>错误信息指出数值范围 |
| TC_RM_006 | 创建私密房间未提供密码 | P1 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. is_private=true但password为空 | 400 Bad Request<br>错误信息提示密码必需 |

### 1.2 获取房间信息测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_RM_007 | 获取房间详情成功 | P1 | 房间已创建 | 1. GET /api/v1/rooms/{room_id} | 200 OK<br>返回完整房间信息 |
| TC_RM_008 | 获取不存在房间 | P1 | 房间不存在 | 1. GET /api/v1/rooms/invalid_room_id | 404 Not Found<br>错误信息"房间不存在" |
| TC_RM_009 | 获取房间列表成功 | P1 | 多个房间已创建 | 1. GET /api/v1/rooms<br>2. 可选参数page, limit | 200 OK<br>返回房间列表和分页信息 |
| TC_RM_010 | 分页参数验证 | P2 | 多个房间已创建 | 1. GET /api/v1/rooms?page=-1&limit=1000 | 400 Bad Request<br>错误信息指出参数限制 |

### 1.3 更新房间测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_RM_011 | 房间创建者更新成功 | P1 | 房间已创建，当前用户为创建者 | 1. PUT /api/v1/rooms/{room_id}<br>2. 提供更新数据 | 200 OK<br>返回更新后的房间信息 |
| TC_RM_012 | 非创建者尝试更新 | P1 | 房间已创建，当前用户不是创建者 | 1. PUT /api/v1/rooms/{room_id}<br>2. 提供更新数据 | 403 Forbidden<br>错误信息"权限不足" |
| TC_RM_013 | 更新不存在房间 | P1 | 房间不存在 | 1. PUT /api/v1/rooms/invalid_room_id | 404 Not Found<br>错误信息"房间不存在" |

### 1.4 删除房间测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_RM_014 | 房间创建者删除成功 | P1 | 房间已创建，当前用户为创建者 | 1. DELETE /api/v1/rooms/{room_id} | 200 OK<br>成功删除信息 |
| TC_RM_015 | 非创建者尝试删除 | P1 | 房间已创建，当前用户不是创建者 | 1. DELETE /api/v1/rooms/{room_id} | 403 Forbidden<br>错误信息"权限不足" |
| TC_RM_016 | 删除不存在房间 | P1 | 房间不存在 | 1. DELETE /api/v1/rooms/invalid_room_id | 404 Not Found<br>错误信息"房间不存在" |

---

## 2. 成员操作测试用例

### 2.1 加入房间测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_MO_001 | 成功加入公开房间 | P1 | 房间已创建且为公开房间 | 1. POST /api/v1/rooms/{room_id}/join<br>2. 提供room_id和display_name | 200 OK<br>返回会话token和房间信息 |
| TC_MO_002 | 加入私密房间密码正确 | P1 | 房间已创建且为私密房间，密码已知 | 1. POST /api/v1/rooms/{room_id}/join<br>2. 提供正确的password | 200 OK<br>返回会话token和房间信息 |
| TC_MO_003 | 加入私密房间密码错误 | P1 | 房间已创建且为私密房间 | 1. POST /api/v1/rooms/{room_id}/join<br>2. 提供错误的password | 403 Forbidden<br>错误信息"密码错误" |
| TC_MO_004 | 加入不存在的房间 | P1 | 房间不存在 | 1. POST /api/v1/rooms/invalid_room_id/join | 404 Not Found<br>错误信息"房间不存在" |
| TC_MO_005 | 房间人数已满 | P2 | 房间人数已达max_users上限 | 1. POST /api/v1/rooms/{room_id}/join | 409 Conflict<br>错误信息"房间人数已满" |

### 2.2 退出房间测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_MO_006 | 正常退出房间 | P1 | 用户已加入房间 | 1. POST /api/v1/rooms/{room_id}/leave | 200 OK<br>成功退出信息 |
| TC_MO_007 | 退出未加入的房间 | P2 | 用户未加入该房间 | 1. POST /api/v1/rooms/{room_id}/leave | 400 Bad Request<br>错误信息"未加入该房间" |
| TC_MO_008 | 退出不存在的房间 | P1 | 房间不存在 | 1. POST /api/v1/rooms/invalid_room_id/leave | 404 Not Found<br>错误信息"房间不存在" |

### 2.3 获取房间成员测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_MO_009 | 获取房间成员列表成功 | P1 | 用户已加入房间或有访问权限 | 1. GET /api/v1/rooms/{room_id}/members | 200 OK<br>返回成员列表 |
| TC_MO_010 | 获取不存在房间成员 | P1 | 房间不存在 | 1. GET /api/v1/rooms/invalid_room_id/members | 404 Not Found<br>错误信息"房间不存在" |

---

## 3. 播放控制测试用例

### 3.1 播放控制API测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_PC_001 | 播放视频 | P1 | 用户已加入房间 | 1. POST /api/v1/rooms/{room_id}/play<br>2. 提供播放参数 | 200 OK<br>播放状态更新 |
| TC_PC_002 | 暂停视频 | P1 | 视频正在播放 | 1. POST /api/v1/rooms/{room_id}/pause | 200 OK<br>暂停状态更新 |
| TC_PC_003 | 跳转播放位置 | P1 | 视频已加载 | 1. POST /api/v1/rooms/{room_id}/seek<br>2. 提供time参数 | 200 OK<br>播放位置更新 |
| TC_PC_004 | 无效跳转位置 | P2 | 视频已加载 | 1. POST /api/v1/rooms/{room_id}/seek<br>2. time参数超出视频时长 | 400 Bad Request<br>错误信息"无效的播放位置" |
| TC_PC_005 | 获取播放状态 | P1 | 用户已加入房间 | 1. GET /api/v1/rooms/{room_id}/status | 200 OK<br>返回当前播放状态 |
| TC_PC_006 | 房间不存在时播放控制 | P1 | 房间不存在 | 1. POST /api/v1/rooms/invalid_room_id/play | 404 Not Found<br>错误信息"房间不存在" |

---

## 4. WebSocket连接测试用例

### 4.1 连接建立测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_WS_001 | 成功建立WebSocket连接 | P1 | 用户已加入房间，获得有效token | 1. 建立WebSocket连接<br>2. URL: /ws/room/{room_id}?token={valid_token} | 连接建立成功<br>收到欢迎消息 |
| TC_WS_002 | 缺少token的连接请求 | P1 | 无有效token | 1. 建立WebSocket连接<br>2. URL: /ws/room/{room_id} | 连接拒绝<br>返回401错误 |
| TC_WS_003 | 无效token的连接请求 | P1 | 持有无效token | 1. 建立WebSocket连接<br>2. URL: /ws/room/{room_id}?token=invalid_token | 连接拒绝<br>返回401错误 |
| TC_WS_004 | 房间不存在时的连接请求 | P1 | 房间不存在 | 1. 建立WebSocket连接<br>2. URL: /ws/room/invalid_room_id?token={valid_token} | 连接拒绝<br>返回404错误 |

### 4.2 消息通信测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_WS_005 | 播放状态同步消息 | P1 | WebSocket连接已建立 | 1. 发送播放控制消息<br>2. 验证广播到其他成员 | 收到播放状态更新消息 |
| TC_WS_006 | 用户加入/离开通知 | P1 | WebSocket连接已建立 | 1. 验证用户加入房间时收到通知<br>2. 验证用户离开房间时收到通知 | 收到相应的通知消息 |
| TC_WS_007 | 心跳保持连接 | P2 | WebSocket连接已建立 | 1. 定期发送ping消息<br>2. 验证pong响应 | 连接保持活跃状态 |

---

## 5. 会话管理测试用例

### 5.1 会话创建和验证测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_SM_001 | 创建会话成功 | P1 | 服务运行正常 | 1. POST /api/v1/sessions<br>2. 可选提供nickname | 201 Created<br>返回session_id和用户信息 |
| TC_SM_002 | 获取会话信息成功 | P1 | 会话已创建 | 1. GET /api/v1/sessions/{session_id} | 200 OK<br>返回会话详细信息 |
| TC_SM_003 | 验证会话有效性 | P1 | 会话已创建 | 1. GET /api/v1/sessions/{session_id}/validate | 200 OK<br>返回验证结果 |
| TC_SM_004 | 更新会话信息 | P1 | 会话已创建 | 1. PUT /api/v1/sessions/{session_id}<br>2. 提供nickname或avatar | 200 OK<br>返回更新后的会话信息 |
| TC_SM_005 | 会话心跳检测 | P1 | 会话已创建 | 1. POST /api/v1/sessions/{session_id}/heartbeat | 200 OK<br>心跳时间更新 |
| TC_SM_006 | 删除会话 | P1 | 会话已创建 | 1. DELETE /api/v1/sessions/{session_id} | 200 OK<br>会话删除成功 |
| TC_SM_007 | 获取不存在的会话 | P1 | 会话不存在 | 1. GET /api/v1/sessions/invalid_session_id | 404 Not Found<br>错误信息"会话不存在" |

---

## 6. 错误处理测试用例

### 6.1 HTTP状态码测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_ER_001 | 400错误 - 缺失必需参数 | P1 | 服务运行正常 | 1. POST /api/v1/rooms<br>2. 缺失必需字段 | 400 Bad Request<br>结构化错误响应 |
| TC_ER_002 | 401错误 - 未认证访问 | P1 | 需要认证的接口 | 1. 访问需要认证的接口<br>2. 不提供session_id | 401 Unauthorized<br>错误信息"未认证" |
| TC_ER_003 | 403错误 - 权限不足 | P1 | 用户权限不足 | 1. 尝试执行需要特定权限的操作 | 403 Forbidden<br>错误信息"权限不足" |
| TC_ER_004 | 404错误 - 资源不存在 | P1 | 访问不存在的资源 | 1. 访问不存在的房间或会话ID | 404 Not Found<br>错误信息"资源不存在" |
| TC_ER_005 | 409错误 - 资源冲突 | P2 | 资源状态冲突 | 1. 尝试创建重复资源<br>2. 尝试加入人数已满的房间 | 409 Conflict<br>错误信息说明冲突原因 |
| TC_ER_006 | 429错误 - 请求频率限制 | P3 | 短时间内大量请求 | 1. 在短时间内发送大量请求 | 429 Too Many Requests<br>错误信息"请求过于频繁" |
| TC_ER_007 | 500错误 - 服务器内部错误 | P2 | 服务器内部错误 | 1. 触发服务器内部错误情况 | 500 Internal Server Error<br>错误信息"服务器内部错误" |

### 6.2 错误响应格式测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_ER_008 | 错误响应格式正确性 | P1 | 发生错误情况 | 1. 触发各种错误情况<br>2. 验证错误响应格式 | 符合ErrorResponse结构<br>包含error, detail, code字段 |
| TC_ER_009 | 错误信息本地化 | P2 | 发生错误情况 | 1. 触发错误情况<br>2. 验证错误信息语言 | 错误信息为中文<br>描述清晰易懂 |

---

## 7. 性能和并发测试用例

### 7.1 性能基准测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_PF_001 | API响应时间基准 | P2 | 服务正常运行 | 1. 对每个API端点进行性能测试<br>2. 记录响应时间 | 95%的请求响应时间 < 500ms |
| TC_PF_002 | 并发用户创建房间 | P2 | 服务正常运行 | 1. 模拟50个并发用户同时创建房间 | 所有请求成功完成<br>无性能退化 |
| TC_PF_003 | WebSocket连接并发测试 | P2 | 服务正常运行 | 1. 模拟100个并发WebSocket连接 | 连接成功率 > 99%<br>无内存泄漏 |
| TC_PF_004 | 数据库连接池压力测试 | P2 | 服务正常运行 | 1. 持续高并发访问数据库的接口 | 数据库连接池稳定<br>无连接耗尽 |

---

## 8. 安全测试用例

### 8.1 认证和授权测试

| 测试用例ID | 用例名称 | 优先级 | 前置条件 | 测试步骤 | 期望结果 |
|-----------|----------|--------|----------|----------|----------|
| TC_SC_001 | Session Token验证 | P1 | 会话已创建 | 1. 使用有效的session_id访问受保护接口 | 正常访问 |
| TC_SC_002 | 无效Session Token | P1 | 持有无效token | 1. 使用过期的session_id访问受保护接口 | 401 Unauthorized |
| TC_SC_003 | 跨房间权限测试 | P1 | 用户A加入房间1，用户B加入房间2 | 1. 用户A尝试访问用户B的房间 | 403 Forbidden或404 Not Found |
| TC_SC_004 | 密码安全测试 | P1 | 私密房间已创建 | 1. 尝试暴力破解房间密码 | 失败率监控<br>触发频率限制 |

---

## 测试数据准备

### 测试用户数据
```json
{
  "test_users": [
    {
      "nickname": "测试用户001",
      "avatar": "https://example.com/avatar1.jpg"
    },
    {
      "nickname": "测试用户002", 
      "avatar": "https://example.com/avatar2.jpg"
    },
    {
      "nickname": "测试用户003",
      "avatar": "https://example.com/avatar3.jpg"
    }
  ]
}
```

### 测试房间数据
```json
{
  "test_rooms": [
    {
      "name": "测试公开房间",
      "description": "用于API测试的公开房间",
      "is_private": false,
      "max_users": 10,
      "media_url": "https://test.com/video1.mp4",
      "media_title": "测试视频1"
    },
    {
      "name": "测试私密房间",
      "description": "用于API测试的私密房间",
      "is_private": true,
      "password": "test123",
      "max_users": 5,
      "media_url": "https://test.com/video2.mp4", 
      "media_title": "测试视频2"
    }
  ]
}
```

---

## 测试执行优先级

### �� P1 - 必须通过 (Critical)
- 所有房间管理基本功能
- 成员加入/退出操作
- 播放控制核心功能
- WebSocket基础连接
- 会话创建和验证
- 基础错误处理

### �� P2 - 重要功能 (Major)  
- 参数边界值测试
- 并发访问测试
- 性能基准测试
- 权限控制测试
- 错误响应格式验证

### 🟢 P3 - 增强功能 (Minor)
- 性能压力测试
- 安全渗透测试
- 长连接稳定性测试
- 极端情况测试

---

## 自动化测试要求

1. **测试覆盖率**: 核心功能100%，边界条件90%
2. **执行时间**: 单次测试执行时间 < 5分钟
3. **环境隔离**: 每个测试用例独立执行，互不干扰
4. **数据清理**: 测试完成后自动清理测试数据
5. **报告生成**: 自动生成详细的测试执行报告
