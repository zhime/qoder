# API接口文档

<cite>
**本文档引用的文件**  
- [router.go](file://backend/internal/api/router.go)
- [auth.go](file://backend/internal/api/auth.go)
- [types.go](file://backend/internal/api/types.go)
- [jwt.go](file://backend/internal/auth/jwt.go)
- [password.go](file://backend/internal/auth/password.go)
- [user.go](file://backend/internal/service/user.go)
- [user.go](file://backend/internal/model/user.go)
- [user.go](file://backend/internal/api/user.go)
- [monitor.go](file://backend/internal/api/monitor.go)
</cite>

## 目录
1. [简介](#简介)
2. [认证接口 (/auth)](#认证接口-auth)
3. [用户接口 (/user)](#用户接口-user)
4. [服务器接口 (/server)](#服务器接口-server)
5. [部署接口 (/deployment)](#部署接口-deployment)
6. [任务接口 (/task)](#任务接口-task)
7. [监控接口 (/monitor)](#监控接口-monitor)
8. [通用响应结构](#通用响应结构)
9. [错误处理与建议](#错误处理与建议)

## 简介
本接口文档详细描述了qoder系统的RESTful API设计，基于`router.go`中定义的路由结构，涵盖认证、用户管理、服务器、部署、任务及监控等核心模块。所有受保护的接口均需通过JWT身份验证，请求头中必须携带`Authorization: Bearer <token>`。文档包含每个端点的HTTP方法、路径、请求参数、响应格式及状态码，并提供实际调用示例。

## 认证接口 /auth

### POST /api/auth/login - 用户登录
- **功能**：用户凭用户名/邮箱和密码登录，返回JWT访问令牌和刷新令牌。
- **HTTP方法**：POST
- **URL路径**：`/api/auth/login`
- **请求头**：`Content-Type: application/json`
- **请求体Schema**：
  ```json
  {
    "username": "string (required)",
    "password": "string (required, min length 6)"
  }
  ```
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "登录成功",
    "data": {
      "access_token": "string",
      "refresh_token": "string",
      "user": {
        "id": 1,
        "username": "admin",
        "email": "admin@example.com",
        "role": "admin",
        "status": 1
      }
    }
  }
  ```
- **状态码**：
  - 200: 登录成功
  - 400: 请求参数错误
  - 401: 用户名或密码错误
  - 500: 令牌生成失败

**JWT令牌生成与密码验证流程分析**  
登录接口在`auth.go`中实现，调用`userService.Login`进行用户验证。该服务在`service/user.go`中通过用户名或邮箱查询用户，使用`auth.CheckPassword`（见`auth/password.go`）比对BCrypt哈希密码。验证通过后，`jwtManager.GenerateToken`和`GenerateRefreshToken`（见`auth/jwt.go`）生成有效期分别为1小时和7天的JWT令牌。

**curl示例**：
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"password123"}'
```

### POST /api/auth/refresh - 刷新令牌
- **功能**：使用刷新令牌获取新的访问令牌和刷新令牌。
- **HTTP方法**：POST
- **URL路径**：`/api/auth/refresh`
- **请求头**：`Content-Type: application/json`
- **请求体Schema**：
  ```json
  {
    "refresh_token": "string (required)"
  }
  ```
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "刷新成功",
    "data": {
      "access_token": "string",
      "refresh_token": "string"
    }
  }
  ```
- **状态码**：
  - 200: 刷新成功
  - 400: 请求参数错误
  - 401: 刷新令牌无效或已过期

**Section sources**
- [auth.go](file://backend/internal/api/auth.go#L70-L108)
- [jwt.go](file://backend/internal/auth/jwt.go#L65-L83)

### GET /api/auth/profile - 获取用户信息
- **功能**：获取当前认证用户的基本信息。
- **HTTP方法**：GET
- **URL路径**：`/api/auth/profile`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "status": 1
    }
  }
  ```
- **状态码**：
  - 200: 获取成功
  - 401: 用户未认证或令牌无效
  - 404: 用户不存在

**Section sources**
- [auth.go](file://backend/internal/api/auth.go#L110-L138)

### POST /api/auth/logout - 用户登出
- **功能**：用户登出，当前令牌失效（当前实现为简单返回成功，实际应加入Redis黑名单）。
- **HTTP方法**：POST
- **URL路径**：`/api/auth/logout`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "登出成功"
  }
  ```
- **状态码**：
  - 200: 登出成功
  - 401: 用户未认证

**Section sources**
- [auth.go](file://backend/internal/api/auth.go#L140-L159)

## 用户接口 /user

### GET /api/users - 获取用户列表
- **功能**：分页获取所有用户列表（需管理员权限）。
- **HTTP方法**：GET
- **URL路径**：`/api/users?page=1&page_size=10`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `page`: 页码（默认1）
  - `page_size`: 每页数量（默认10，最大100）
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "list": [
        {
          "id": 1,
          "username": "admin",
          "email": "admin@example.com",
          "role": "admin",
          "status": 1
        }
      ],
      "total": 1,
      "page": 1,
      "page_size": 10
    }
  }
  ```
- **状态码**：
  - 200: 获取成功
  - 400: 参数错误
  - 500: 查询失败

**Section sources**
- [user.go](file://backend/internal/api/user.go#L34-L65)
- [user.go](file://backend/internal/service/user.go#L128-L167)

### GET /api/users/:id - 获取用户详情
- **功能**：根据ID获取单个用户信息。
- **HTTP方法**：GET
- **URL路径**：`/api/users/1`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "id": 1,
      "username": "admin",
      "email": "admin@example.com",
      "role": "admin",
      "status": 1
    }
  }
  ```
- **状态码**：
  - 200: 获取成功
  - 400: ID无效
  - 404: 用户不存在

**Section sources**
- [user.go](file://backend/internal/api/user.go#L67-L93)

### POST /api/users - 创建用户
- **功能**：创建新用户（需管理员权限）。
- **HTTP方法**：POST
- **URL路径**：`/api/users`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **请求体Schema**：
  ```json
  {
    "username": "string (required)",
    "email": "string (required, valid email)",
    "password": "string (required, min length 6)",
    "role": "string (required, oneof: user, admin)"
  }
  ```
- **成功响应**（201 Created）：
  ```json
  {
    "code": 201,
    "message": "用户创建成功",
    "data": {
      "id": 2,
      "username": "newuser",
      "email": "newuser@example.com",
      "role": "user",
      "status": 1
    }
  }
  ```
- **状态码**：
  - 201: 创建成功
  - 400: 参数错误或用户名/邮箱已存在

**Section sources**
- [user.go](file://backend/internal/api/user.go#L10-L32)
- [user.go](file://backend/internal/service/user.go#L70-L108)

### PUT /api/users/:id - 更新用户
- **功能**：更新用户信息（需管理员权限）。
- **HTTP方法**：PUT
- **URL路径**：`/api/users/1`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **请求体Schema**（可选字段）：
  ```json
  {
    "username": "string",
    "email": "string",
    "role": "string (oneof: user, admin)",
    "status": 0 or 1
  }
  ```
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "更新成功"
  }
  ```
- **状态码**：
  - 200: 更新成功
  - 400: 参数错误或用户不存在
  - 404: 用户不存在

**Section sources**
- [user.go](file://backend/internal/api/user.go#L95-L147)

### DELETE /api/users/:id - 删除用户
- **功能**：删除用户（软删除，需管理员权限）。
- **HTTP方法**：DELETE
- **URL路径**：`/api/users/1`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "删除成功"
  }
  ```
- **状态码**：
  - 200: 删除成功
  - 400: ID无效或用户不存在

**Section sources**
- [user.go](file://backend/internal/api/user.go#L149-L176)

## 服务器接口 /server

### GET /api/servers - 获取服务器列表
- **功能**：获取所有服务器列表。
- **HTTP方法**：GET
- **URL路径**：`/api/servers`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "servers list"
  }
  ```
- **状态码**：
  - 200: 成功

### POST /api/servers - 创建服务器
- **功能**：创建新的服务器配置。
- **HTTP方法**：POST
- **URL路径**：`/api/servers`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "create server"
  }
  ```
- **状态码**：
  - 200: 成功

**服务器CRUD操作规范说明**  
根据`router.go`，服务器接口目前仅实现了GET和POST方法，PUT和DELETE方法尚未实现。完整的CRUD规范应包含：
- **创建**：POST /api/servers，请求体包含`CreateServerRequest`（见`types.go`）。
- **读取**：GET /api/servers（列表），GET /api/servers/:id（详情）。
- **更新**：PUT /api/servers/:id，请求体包含`UpdateServerRequest`。
- **删除**：DELETE /api/servers/:id。

**Section sources**
- [router.go](file://backend/internal/api/router.go#L70-L77)

## 部署接口 /deployment

### GET /api/deployments - 获取部署列表
- **功能**：获取所有部署记录。
- **HTTP方法**：GET
- **URL路径**：`/api/deployments`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "deployments list"
  }
  ```
- **状态码**：
  - 200: 成功

### POST /api/deployments - 创建部署
- **功能**：创建新的部署任务。
- **HTTP方法**：POST
- **URL路径**：`/api/deployments`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "create deployment"
  }
  ```
- **状态码**：
  - 200: 成功

**Section sources**
- [router.go](file://backend/internal/api/router.go#L79-L86)

## 任务接口 /task

### GET /api/tasks - 获取任务列表
- **功能**：获取所有定时任务列表。
- **HTTP方法**：GET
- **URL路径**：`/api/tasks`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "tasks list"
  }
  ```
- **状态码**：
  - 200: 成功

### POST /api/tasks - 创建任务
- **功能**：创建新的定时任务。
- **HTTP方法**：POST
- **URL路径**：`/api/tasks`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "create task"
  }
  ```
- **状态码**：
  - 200: 成功

**Section sources**
- [router.go](file://backend/internal/api/router.go#L88-L95)

## 监控接口 /monitor

### GET /api/monitor/dashboard - 获取仪表盘数据
- **功能**：获取系统仪表盘综合数据，包括统计信息、告警和最近活动。
- **HTTP方法**：GET
- **URL路径**：`/api/monitor/dashboard`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "stats": { /* 系统统计 */ },
      "alerts": [ /* 告警列表 */ ],
      "recent_activities": [ /* 最近活动 */ ]
    }
  }
  ```
- **状态码**：
  - 200: 成功
  - 500: 获取失败

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L100-L148)

### GET /api/monitor/stats - 获取系统统计
- **功能**：获取系统级统计信息（如服务器总数、在线数等）。
- **HTTP方法**：GET
- **URL路径**：`/api/monitor/stats`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": { /* 统计数据 */ }
  }
  ```
- **状态码**：
  - 200: 成功
  - 500: 获取失败

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L78-L98)

### GET /api/monitor/servers/:id/metrics - 获取服务器实时指标
- **功能**：获取指定服务器的实时监控指标（CPU、内存等）。
- **HTTP方法**：GET
- **URL路径**：`/api/monitor/servers/1/metrics`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": { /* 指标数据 */ }
  }
  ```
- **状态码**：
  - 200: 成功
  - 400: ID无效
  - 404: 获取失败

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L10-L38)

### GET /api/monitor/servers/:id/status - 获取服务器状态
- **功能**：获取指定服务器的运行状态。
- **HTTP方法**：GET
- **URL路径**：`/api/monitor/servers/1/status`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "server_id": 1,
      "status": "online"
    }
  }
  ```
- **状态码**：
  - 200: 成功
  - 400: ID无效
  - 500: 获取失败

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L40-L66)

### GET /api/monitor/servers/:id/history - 获取历史监控数据
- **功能**：获取指定服务器的历史监控数据。
- **HTTP方法**：GET
- **URL路径**：`/api/monitor/servers/1/history?time_range=1h&metric=cpu`
- **请求头**：`Authorization: Bearer <access_token>`
- **查询参数**：
  - `time_range`: 时间范围（1h, 6h, 24h, 7d）
  - `metric`: 指标类型（cpu, memory, disk, network）
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "获取成功",
    "data": {
      "server_id": 1,
      "metric": "cpu",
      "time_range": "1h",
      "data": [ /* 时间序列数据点 */ ]
    }
  }
  ```
- **状态码**：
  - 200: 成功
  - 400: ID无效

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L150-L257)

### POST /api/monitor/servers - 添加服务器到监控
- **功能**：将服务器添加到监控列表。
- **HTTP方法**：POST
- **URL路径**：`/api/monitor/servers`
- **请求头**：`Authorization: Bearer <access_token>`, `Content-Type: application/json`
- **请求体Schema**：
  ```json
  {
    "server_id": 1
  }
  ```
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "添加监控成功"
  }
  ```
- **状态码**：
  - 200: 成功
  - 400: 参数错误

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L178-L197)

### DELETE /api/monitor/servers/:id - 从监控移除服务器
- **功能**：将服务器从监控列表中移除。
- **HTTP方法**：DELETE
- **URL路径**：`/api/monitor/servers/1`
- **请求头**：`Authorization: Bearer <access_token>`
- **成功响应**（200 OK）：
  ```json
  {
    "code": 200,
    "message": "移除监控成功"
  }
  ```
- **状态码**：
  - 200: 成功
  - 400: ID无效

**Section sources**
- [monitor.go](file://backend/internal/api/monitor.go#L199-L220)

## 通用响应结构
所有API响应均遵循统一格式：
```json
{
  "code": 200,
  "message": "操作成功",
  "data": { /* 可选数据 */ }
}
```
- **code**：业务状态码（200表示成功）。
- **message**：响应消息。
- **data**：返回的数据，分页数据使用`PageResponse`结构。

**Section sources**
- [types.go](file://backend/internal/api/types.go#L3-L10)

## 错误处理与建议
- **400 Bad Request**：检查请求参数是否符合Schema，特别是必填字段和格式。
- **401 Unauthorized**：确保`Authorization`头包含有效的Bearer令牌，或重新登录获取新令牌。
- **404 Not Found**：确认资源ID是否存在。
- **500 Internal Error**：服务端内部错误，检查日志或联系管理员。
- **建议**：在客户端实现自动刷新令牌机制，当收到401时尝试用刷新令牌获取新访问令牌。