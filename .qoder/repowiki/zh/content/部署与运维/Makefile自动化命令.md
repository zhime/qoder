# Makefile自动化命令

<cite>
**本文档引用文件**  
- [Makefile](file://Makefile)
- [docker-compose.yml](file://docker-compose.yml)
- [docker/Dockerfile.backend](file://docker/Dockerfile.backend)
- [docker/Dockerfile.frontend](file://docker/Dockerfile.frontend)
- [scripts/start-dev.sh](file://scripts/start-dev.sh)
- [frontend/package.json](file://frontend/package.json)
- [backend/go.mod](file://backend/go.mod)
</cite>

## 目录
1. [简介](#简介)
2. [核心构建命令](#核心构建命令)
3. [开发与测试命令](#开发与测试命令)
4. [Docker相关命令](#docker相关命令)
5. [辅助与维护命令](#辅助与维护命令)
6. [常用组合命令示例](#常用组合命令示例)
7. [自定义任务扩展](#自定义任务扩展)

## 简介
`Makefile` 是 qoder 项目中用于自动化构建、测试、部署的核心脚本。通过定义一系列标准化的 `make` 命令，开发者可以一键完成从代码编译、依赖安装到容器化部署的全流程操作，极大提升了开发效率和部署一致性。

该 Makefile 覆盖了后端 Go 服务编译、前端 Vue 项目构建、单元测试运行、Docker 镜像打包与容器编排等关键环节，实现了开发、测试、生产环境的一体化管理。

## 核心构建命令

### build：全栈构建
`make build` 是主构建命令，负责编译整个项目的前后端代码。

- **build-backend**：进入 `backend` 目录，使用 `go build` 编译主程序，输出二进制文件至 `bin/devops-backend`
- **build-frontend**：进入 `frontend` 目录，执行 `npm run build`，生成生产级静态资源至 `dist/` 目录

此命令是部署前的标准准备步骤，确保前后端代码均已编译就绪。

**Section sources**
- [Makefile](file://Makefile#L18-L25)

### build-backend：后端编译
单独执行后端 Go 程序的编译任务。使用标准 `go build` 命令，不启用 CGO，生成静态可执行文件，便于跨平台部署。

**Section sources**
- [Makefile](file://Makefile#L20-L22)

### build-frontend：前端构建
调用 `npm run build` 执行前端构建流程。根据 `package.json` 定义，该命令会执行 TypeScript 类型检查并使用 Vite 构建生产环境资源包。

**Section sources**
- [Makefile](file://Makefile#L23-L25)
- [frontend/package.json](file://frontend/package.json#L6-L8)

## 开发与测试命令

### dev：启动开发环境
`make dev` 启动完整的本地开发环境。在非 Windows 系统中，调用 `scripts/start-dev.sh` 脚本，自动执行以下流程：

1. 检查 Docker 运行状态
2. 使用 `docker-compose` 启动 MySQL 和 Redis 服务
3. 等待数据库就绪（10秒）
4. 在后台启动 Go 后端服务
5. 在后台启动 Vue 前端开发服务器
6. 提供 Ctrl+C 信号捕获，自动清理进程并关闭容器

该命令实现了“一键启动”开发环境，极大简化了本地开发配置。

**Section sources**
- [Makefile](file://Makefile#L34-L42)
- [scripts/start-dev.sh](file://scripts/start-dev.sh#L1-L44)

### test：运行测试套件
`make test` 执行项目所有测试，包含前后端测试任务。

- **test-backend**：在 `backend` 目录下运行 `go test -v ./...`，执行所有 Go 单元测试
- **test-frontend**：当前为空操作，提示“前端测试暂未实现”

**Section sources**
- [Makefile](file://Makefile#L27-L32)

## Docker相关命令

### docker-build：构建Docker镜像
使用项目根目录下的两个 Dockerfile 分别构建前后端镜像：

- 后端镜像 `devops-backend:latest`：基于 `golang:1.21-alpine` 多阶段构建，最终使用 `alpine:latest` 运行静态编译的 Go 程序
- 前端镜像 `devops-frontend:latest`：基于 `node:18-alpine` 构建，使用 `nginx:alpine` 作为运行时，托管构建后的静态文件

**Section sources**
- [Makefile](file://Makefile#L44-L48)
- [docker/Dockerfile.backend](file://docker/Dockerfile.backend#L1-L18)
- [docker/Dockerfile.frontend](file://docker/Dockerfile.frontend#L1-L16)

### docker-up：启动Docker容器
执行 `docker-compose up -d`，根据 `docker-compose.yml` 定义启动四个服务：

- **mysql**：MySQL 8.0 数据库，初始化脚本挂载 `init.sql`
- **redis**：Redis 7 缓存服务，启用持久化
- **backend**：后端服务，依赖数据库和 Redis，暴露 8080 端口
- **frontend**：前端服务，依赖后端，通过 Nginx 暴露 80 端口

**Section sources**
- [Makefile](file://Makefile#L50-L53)
- [docker-compose.yml](file://docker-compose.yml#L1-L61)

### docker-down：停止Docker容器
执行 `docker-compose down`，停止并移除所有由 `docker-up` 启动的容器、网络，但保留卷数据（如数据库）。

**Section sources**
- [Makefile](file://Makefile#L55-L58)

## 辅助与维护命令

### clean：清理构建产物
清除所有生成的文件和缓存：

- 删除 `bin/` 目录（后端二进制）
- 删除 `frontend/dist/` 目录（前端构建产物）
- 在 `backend` 目录执行 `go clean` 清理 Go 编译缓存

**Section sources**
- [Makefile](file://Makefile#L36-L40)

### install-deps：安装依赖
分别安装前后端依赖：

- 后端：执行 `go mod tidy` 整理 Go 模块依赖
- 前端：执行 `npm install` 安装 Node.js 包

**Section sources**
- [Makefile](file://Makefile#L60-L65)
- [backend/go.mod](file://backend/go.mod#L1-L71)

### lint：代码检查
执行前后端代码静态检查：

- 后端：`go vet` 和 `gofmt -l` 检查代码规范
- 前端：`npm run lint` 执行 ESLint 检查

**Section sources**
- [Makefile](file://Makefile#L67-L72)
- [frontend/package.json](file://frontend/package.json#L9-L11)

### migrate：数据库迁移
执行 `go run cmd/migrate.go` 运行数据库迁移脚本。尽管文件未在当前上下文中找到，但根据调用方式，该命令应负责执行数据库 schema 的版本升级。

**Section sources**
- [Makefile](file://Makefile#L74-L76)

## 常用组合命令示例
以下是一些高频使用的 `make` 命令组合：

```bash
# 清理旧构建并重新构建运行
make clean build run

# 安装依赖后启动开发环境
make install-deps dev

# 构建Docker镜像并启动容器
make docker-build docker-up

# 测试全流程（清理、构建、测试）
make clean build test
```

这些组合命令可添加到 CI/CD 流程中，实现自动化集成。

## 自定义任务扩展
可在 `Makefile` 中添加新的目标以扩展功能。例如：

```makefile
# 自定义：生成代码
generate:
	@echo "生成代码..."
	cd backend && go generate ./...

# 自定义：运行性能测试
bench:
	@echo "运行性能测试..."
	cd backend && go test -bench=. ./...
```

新任务应遵循现有格式，使用 `.PHONY` 声明，并在 `help` 目标中添加说明。通过这种方式，可灵活扩展自动化流程，适应项目演进需求。