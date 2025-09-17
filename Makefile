# 自动化运维平台 Makefile

.PHONY: help build test clean dev docker-build docker-up docker-down

# 默认目标
help:
	@echo "可用命令:"
	@echo "  make build         - 构建后端和前端"
	@echo "  make test          - 运行所有测试"
	@echo "  make dev           - 启动开发环境"
	@echo "  make clean         - 清理构建文件"
	@echo "  make docker-build  - 构建Docker镜像"
	@echo "  make docker-up     - 启动Docker容器"
	@echo "  make docker-down   - 停止Docker容器"

# 构建
build: build-backend build-frontend

build-backend:
	@echo "构建后端..."
	cd backend && go build -o ../bin/devops-backend cmd/main.go

build-frontend:
	@echo "构建前端..."
	cd frontend && npm run build

# 测试
test: test-backend test-frontend

test-backend:
	@echo "运行后端测试..."
	cd backend && go test -v ./...

test-frontend:
	@echo "运行前端测试..."
	@echo "前端测试暂未实现"

# 开发环境
dev:
	@echo "启动开发环境..."
ifeq ($(OS),Windows_NT)
	@scripts/start-dev.bat
else
	@chmod +x scripts/start-dev.sh
	@scripts/start-dev.sh
endif

# 清理
clean:
	@echo "清理构建文件..."
	rm -rf bin/
	rm -rf frontend/dist/
	cd backend && go clean

# Docker
docker-build:
	@echo "构建Docker镜像..."
	docker build -f docker/Dockerfile.backend -t devops-backend:latest .
	docker build -f docker/Dockerfile.frontend -t devops-frontend:latest .

docker-up:
	@echo "启动Docker容器..."
	docker-compose up -d

docker-down:
	@echo "停止Docker容器..."
	docker-compose down

# 安装依赖
install-deps:
	@echo "安装后端依赖..."
	cd backend && go mod tidy
	@echo "安装前端依赖..."
	cd frontend && npm install

# 代码检查
lint:
	@echo "后端代码检查..."
	cd backend && go vet ./...
	cd backend && gofmt -l .
	@echo "前端代码检查..."
	cd frontend && npm run lint

# 数据库迁移
migrate:
	@echo "数据库迁移..."
	cd backend && go run cmd/migrate.go