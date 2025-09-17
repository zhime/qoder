-- 创建数据库
CREATE DATABASE IF NOT EXISTS devops_platform CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE devops_platform;

-- 创建默认管理员用户
INSERT INTO users (id, username, email, password, role, status, created_at, updated_at) 
VALUES (1, 'admin', 'admin@example.com', '$2a$10$N.zmdr9k7uOCQb376NoUnuTGEkXTjjktb5N9gKDMfpMxBfg.a..ou', 'admin', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE username=username;

-- 默认密码是: admin123

-- 创建测试服务器
INSERT INTO servers (id, name, host, port, username, password, status, environment, description, created_at, updated_at)
VALUES (1, '测试服务器', '127.0.0.1', 22, 'root', '', 1, 'test', '本地测试服务器', NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;