# 自动化运维平台部署文档

## 系统要求

### 最低系统要求
- CPU: 2核心
- 内存: 4GB RAM
- 磁盘: 20GB 可用空间
- 操作系统: Linux (Ubuntu 18.04+), Windows 10+, macOS 10.15+

### 依赖软件
- Docker 20.10+
- Docker Compose 2.0+
- Node.js 18+ (开发环境)
- Go 1.21+ (开发环境)
- MySQL 8.0+
- Redis 6.0+

## 快速开始

### 1. 获取代码
```bash
git clone <repository-url>
cd qoder
```

### 2. 使用Docker部署（推荐）
```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

### 3. 访问应用
- 前端地址: http://localhost
- 后端API: http://localhost/api
- 默认管理员账号: admin / admin123

## 开发环境部署

### 1. 安装依赖
```bash
# 后端依赖
cd backend
go mod tidy

# 前端依赖
cd ../frontend
npm install
```

### 2. 配置环境
复制并修改配置文件：
```bash
cp backend/configs/config.yaml backend/configs/config.local.yaml
cp frontend/.env frontend/.env.local
```

### 3. 启动数据库
```bash
docker-compose up -d mysql redis
```

### 4. 启动后端服务
```bash
cd backend
go run cmd/main.go
```

### 5. 启动前端服务
```bash
cd frontend
npm run dev
```

## 生产环境部署

### 1. 环境准备
```bash
# 创建部署目录
mkdir -p /opt/devops-platform
cd /opt/devops-platform

# 复制代码和配置
cp -r /path/to/qoder/* .
```

### 2. 配置修改
修改 `backend/configs/config.yaml`：
```yaml
server:
  port: 8080
  mode: release

database:
  host: your-mysql-host
  username: your-mysql-user
  password: your-mysql-password
  database: devops_platform

redis:
  host: your-redis-host
  password: your-redis-password

jwt:
  secret: your-production-secret-key
```

### 3. SSL证书配置
如果使用HTTPS，修改 `docker/nginx.conf`：
```nginx
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
}
```

### 4. 启动服务
```bash
# 构建并启动
docker-compose -f docker-compose.prod.yml up -d

# 或使用Makefile
make docker-build
make docker-up
```

## 数据备份与恢复

### 数据备份
```bash
# 备份MySQL数据
docker exec devops_mysql mysqldump -u root -p devops_platform > backup.sql

# 备份Redis数据
docker exec devops_redis redis-cli SAVE
docker cp devops_redis:/data/dump.rdb ./redis-backup.rdb
```

### 数据恢复
```bash
# 恢复MySQL数据
docker exec -i devops_mysql mysql -u root -p devops_platform < backup.sql

# 恢复Redis数据
docker cp ./redis-backup.rdb devops_redis:/data/dump.rdb
docker restart devops_redis
```

## 监控与日志

### 日志查看
```bash
# 查看应用日志
docker-compose logs -f backend
docker-compose logs -f frontend

# 查看数据库日志
docker-compose logs -f mysql
docker-compose logs -f redis
```

### 性能监控
可以集成以下监控工具：
- Prometheus + Grafana
- ELK Stack (Elasticsearch, Logstash, Kibana)
- APM (Application Performance Monitoring)

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查数据库服务是否启动
   - 验证连接配置是否正确
   - 确认网络连通性

2. **前端无法访问后端API**
   - 检查后端服务是否正常运行
   - 验证API代理配置
   - 检查CORS设置

3. **JWT认证失败**
   - 确认JWT密钥配置正确
   - 检查token是否过期
   - 验证请求头格式

### 调试模式
```bash
# 开启调试日志
export LOG_LEVEL=debug

# 后端调试模式
cd backend
go run cmd/main.go --debug

# 前端开发模式
cd frontend
npm run dev
```

## 安全配置

### 1. 数据库安全
- 使用强密码
- 限制数据库访问IP
- 定期更新数据库版本

### 2. 应用安全
- 使用HTTPS
- 定期更换JWT密钥
- 实施访问频率限制

### 3. 服务器安全
- 定期更新系统补丁
- 配置防火墙规则
- 使用非root用户运行服务

## 版本更新

### 1. 备份数据
在更新前务必备份数据库和重要配置文件。

### 2. 更新代码
```bash
git pull origin main
```

### 3. 更新依赖
```bash
# 后端依赖
cd backend && go mod tidy

# 前端依赖
cd frontend && npm install
```

### 4. 数据库迁移
```bash
# 如果有数据库结构变更
make migrate
```

### 5. 重启服务
```bash
docker-compose down
docker-compose up -d
```

## 联系支持

如遇到部署问题，请提供以下信息：
- 操作系统版本
- Docker版本
- 错误日志
- 配置文件（隐藏敏感信息）

技术支持邮箱: support@example.com