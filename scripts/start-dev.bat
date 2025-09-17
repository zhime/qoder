@echo off
echo 启动自动化运维平台开发环境...

REM 检查Docker是否运行
docker info >nul 2>&1
if errorlevel 1 (
    echo 错误: Docker 未运行，请先启动 Docker
    pause
    exit /b 1
)

REM 启动数据库和Redis
echo 启动数据库和缓存服务...
docker-compose up -d mysql redis

REM 等待数据库启动
echo 等待数据库启动...
timeout /t 10 /nobreak >nul

REM 启动后端服务
echo 启动后端服务...
start "Backend Server" cmd /k "cd backend && go run cmd/main.go"

REM 等待后端启动
timeout /t 5 /nobreak >nul

REM 启动前端服务
echo 启动前端服务...
start "Frontend Server" cmd /k "cd frontend && npm run dev"

echo.
echo 服务启动完成!
echo 前端地址: http://localhost:3000
echo 后端地址: http://localhost:8080
echo.
echo 按任意键停止所有服务...
pause >nul

REM 停止服务
taskkill /f /im node.exe >nul 2>&1
taskkill /f /im go.exe >nul 2>&1
docker-compose down
echo 所有服务已停止