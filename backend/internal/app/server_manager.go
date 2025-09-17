package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"devops/internal/api"
	"devops/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServerManager 服务器管理器
type ServerManager struct {
	server *http.Server
	router *gin.Engine
	config config.Server
}

// NewServerManager 创建服务器管理器实例
func NewServerManager(cfg config.Server) *ServerManager {
	return &ServerManager{
		config: cfg,
	}
}

// Initialize 初始化HTTP服务器
func (sm *ServerManager) Initialize(db *gorm.DB, rdb *redis.Client, cfg *config.Config) error {
	// 设置Gin模式
	gin.SetMode(sm.config.Mode)

	// 初始化路由
	router := api.NewRouter(db, rdb, cfg)
	sm.router = router

	// 创建HTTP服务器
	server := &http.Server{
		Addr:           fmt.Sprintf(":%d", sm.config.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(sm.config.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(sm.config.WriteTimeout) * time.Second,
		IdleTimeout:    time.Duration(sm.config.IdleTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1MB
	}
	sm.server = server

	return nil
}

// Start 启动HTTP服务器
func (sm *ServerManager) Start() error {
	if sm.server == nil {
		return fmt.Errorf("服务器未初始化")
	}

	log.Printf("服务器启动在端口 %d", sm.config.Port)

	if err := sm.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	return nil
}

// Shutdown 优雅关闭服务器
func (sm *ServerManager) Shutdown(ctx context.Context) error {
	if sm.server == nil {
		return nil
	}

	log.Println("正在关闭服务器...")

	if err := sm.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("服务器关闭失败: %w", err)
	}

	log.Println("服务器已关闭")
	return nil
}

// GetServer 获取HTTP服务器实例
func (sm *ServerManager) GetServer() *http.Server {
	return sm.server
}

// GetRouter 获取Gin路由器实例
func (sm *ServerManager) GetRouter() *gin.Engine {
	return sm.router
}

// IsRunning 检查服务器是否正在运行
func (sm *ServerManager) IsRunning() bool {
	return sm.server != nil
}

// GetListenAddr 获取服务器监听地址
func (sm *ServerManager) GetListenAddr() string {
	if sm.server == nil {
		return ""
	}
	return sm.server.Addr
}
