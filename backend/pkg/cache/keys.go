package cache

import (
	"fmt"
	"time"
)

// 缓存键前缀
const (
	PrefixUser       = "user"
	PrefixServer     = "server"
	PrefixDeployment = "deployment"
	PrefixTask       = "task"
	PrefixMetrics    = "metrics"
	PrefixSession    = "session"
)

// 缓存过期时间
const (
	TTLUserInfo      = 1 * time.Hour      // 用户信息缓存1小时
	TTLServerList    = 10 * time.Minute   // 服务器列表缓存10分钟
	TTLServerMetrics = 5 * time.Minute    // 服务器监控数据缓存5分钟
	TTLDeployStatus  = 30 * time.Minute   // 部署状态缓存30分钟
	TTLTaskNextRun   = 0                  // 任务执行队列永不过期
	TTLSession       = 24 * time.Hour     // 会话缓存24小时
	TTLRefreshToken  = 7 * 24 * time.Hour // 刷新令牌缓存7天
)

// CacheKeys 缓存键生成器
type CacheKeys struct{}

// NewCacheKeys 创建缓存键生成器
func NewCacheKeys() *CacheKeys {
	return &CacheKeys{}
}

// UserInfo 用户信息缓存键
func (k *CacheKeys) UserInfo(userID uint) string {
	return fmt.Sprintf("%s:info:%d", PrefixUser, userID)
}

// UserPermissions 用户权限缓存键
func (k *CacheKeys) UserPermissions(userID uint) string {
	return fmt.Sprintf("%s:permissions:%d", PrefixUser, userID)
}

// ServerList 服务器列表缓存键
func (k *CacheKeys) ServerList(userID uint) string {
	return fmt.Sprintf("%s:list:%d", PrefixServer, userID)
}

// ServerInfo 服务器信息缓存键
func (k *CacheKeys) ServerInfo(serverID uint) string {
	return fmt.Sprintf("%s:info:%d", PrefixServer, serverID)
}

// ServerMetrics 服务器监控数据缓存键
func (k *CacheKeys) ServerMetrics(serverID uint) string {
	return fmt.Sprintf("%s:%d", PrefixMetrics, serverID)
}

// ServerStatus 服务器状态缓存键
func (k *CacheKeys) ServerStatus(serverID uint) string {
	return fmt.Sprintf("%s:status:%d", PrefixServer, serverID)
}

// DeploymentStatus 部署状态缓存键
func (k *CacheKeys) DeploymentStatus(deploymentID uint) string {
	return fmt.Sprintf("%s:status:%d", PrefixDeployment, deploymentID)
}

// DeploymentLogs 部署日志缓存键
func (k *CacheKeys) DeploymentLogs(deploymentID uint) string {
	return fmt.Sprintf("%s:logs:%d", PrefixDeployment, deploymentID)
}

// TaskNextRun 任务执行队列缓存键
func (k *CacheKeys) TaskNextRun() string {
	return fmt.Sprintf("%s:next_run", PrefixTask)
}

// TaskExecution 任务执行状态缓存键
func (k *CacheKeys) TaskExecution(taskID uint) string {
	return fmt.Sprintf("%s:execution:%d", PrefixTask, taskID)
}

// TaskLock 任务锁缓存键
func (k *CacheKeys) TaskLock(taskID uint) string {
	return fmt.Sprintf("%s:lock:%d", PrefixTask, taskID)
}

// UserSession 用户会话缓存键
func (k *CacheKeys) UserSession(sessionID string) string {
	return fmt.Sprintf("%s:session:%s", PrefixSession, sessionID)
}

// RefreshToken 刷新令牌缓存键
func (k *CacheKeys) RefreshToken(userID uint, tokenID string) string {
	return fmt.Sprintf("%s:refresh:%d:%s", PrefixSession, userID, tokenID)
}

// LoginAttempts 登录尝试次数缓存键
func (k *CacheKeys) LoginAttempts(ip string) string {
	return fmt.Sprintf("%s:login_attempts:%s", PrefixUser, ip)
}

// APIRateLimit API限流缓存键
func (k *CacheKeys) APIRateLimit(userID uint, endpoint string) string {
	return fmt.Sprintf("rate_limit:%d:%s", userID, endpoint)
}

// OnlineUsers 在线用户缓存键
func (k *CacheKeys) OnlineUsers() string {
	return "online_users"
}

// SystemStats 系统统计缓存键
func (k *CacheKeys) SystemStats() string {
	return "system_stats"
}
