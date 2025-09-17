#!/bin/bash

# 启动开发环境脚本

echo "启动自动化运维平台开发环境..."

# 检查Docker是否运行
if ! docker info > /dev/null 2>&1; then
    echo "错误: Docker 未运行，请先启动 Docker"
    exit 1
fi

# 启动数据库和Redis
echo "启动数据库和缓存服务..."
docker-compose up -d mysql redis

# 等待数据库启动
echo "等待数据库启动..."
sleep 10

# 启动后端服务
echo "启动后端服务..."
cd backend
go run cmd/main.go &
BACKEND_PID=$!

# 启动前端服务
echo "启动前端服务..."
cd ../frontend
npm run dev &
FRONTEND_PID=$!

echo ""
echo "服务启动完成!"
echo "前端地址: http://localhost:3000"
echo "后端地址: http://localhost:8080"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 捕获中断信号，清理进程
trap "echo '正在停止服务...'; kill $BACKEND_PID $FRONTEND_PID; docker-compose down; exit" INT

# 等待进程结束
wait