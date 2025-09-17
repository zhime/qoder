# 自动化运维平台

一个集服务器监控、应用部署、任务调度、日志管理于一体的综合性运维管理系统。

## 技术栈

### 前端
- Vue 3 - 前端框架
- Vue Router - 路由管理
- Pinia - 状态管理
- Element Plus - UI组件库
- Axios - HTTP客户端
- ECharts - 数据可视化

### 后端
- Golang - 后端开发语言
- Gin - Web框架
- GORM - ORM框架
- JWT - 身份验证
- Viper - 配置管理
- Zap - 日志框架
- MySQL - 主数据库
- Redis - 缓存数据库

## 项目结构

```
├── backend/                 # 后端代码
│   ├── cmd/                # 主程序入口
│   ├── internal/           # 内部包
│   │   ├── api/           # API 处理器
│   │   ├── auth/          # 认证相关
│   │   ├── config/        # 配置管理
│   │   ├── model/         # 数据模型
│   │   ├── service/       # 业务逻辑
│   │   └── middleware/    # 中间件
│   ├── pkg/               # 公共包
│   └── configs/           # 配置文件
├── frontend/              # 前端代码
│   ├── src/              # 源代码
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面
│   │   ├── router/       # 路由
│   │   ├── store/        # 状态管理
│   │   └── api/          # API服务
│   └── public/           # 静态资源
├── scripts/              # 脚本文件
├── docs/                 # 文档
└── docker/              # Docker配置
```

## 快速开始

### 后端启动
```bash
cd backend
go mod tidy
go run cmd/main.go
```

### 前端启动
```bash
cd frontend
npm install
npm run dev
```

## API文档

- 认证相关：`/api/auth/*`
- 服务器管理：`/api/servers/*`
- 部署管理：`/api/deployments/*`
- 任务调度：`/api/tasks/*`
- 用户管理：`/api/users/*`

## 许可证

MIT License