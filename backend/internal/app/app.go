package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"devops/internal/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Application 应用主体
type Application struct {
	config *config.Config
	db     *gorm.DB
	rdb    *redis.Client
	logger *zap.Logger

	// 管理器组件
	configMgr   *ConfigManager
	databaseMgr *DatabaseManager
	cacheMgr    *CacheManager
	serverMgr   *ServerManager

	// 关闭通道
	shutdownCh chan struct{}
}

// New 创建应用实例
func New() *Application {
	return &Application{
		shutdownCh: make(chan struct{}),
	}
}

// Run 运行应用
func (app *Application) Run() error {
	// 第一步：初始化配置
	if err := app.initConfig(); err != nil {
		return fmt.Errorf("配置初始化失败: %w", err)
	}

	// 第二步：初始化日志
	if err := app.initLogger(); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}

	// 第三步：初始化数据库
	if err := app.initDatabase(); err != nil {
		return fmt.Errorf("数据库初始化失败: %w", err)
	}

	// 第四步：初始化缓存
	if err := app.initCache(); err != nil {
		log.Printf("缓存初始化失败: %v, 继续运行但缓存功能不可用", err)
	}

	// 第五步：初始化服务器
	if err := app.initServer(); err != nil {
		return fmt.Errorf("服务器初始化失败: %w", err)
	}

	// 第六步：启动服务器
	if err := app.startServer(); err != nil {
		return fmt.Errorf("服务器启动失败: %w", err)
	}

	// 第七步：等待关闭信号
	app.waitForShutdown()

	return nil
}

// initConfig 初始化配置
func (app *Application) initConfig() error {
	app.configMgr = NewConfigManager()

	config, err := app.configMgr.Load()
	if err != nil {
		return err
	}

	app.config = config
	return nil
}

// initLogger 初始化日志
func (app *Application) initLogger() error {
	return app.configMgr.InitLogger()
}

// initDatabase 初始化数据库
func (app *Application) initDatabase() error {
	app.databaseMgr = NewDatabaseManager(app.config.Database)

	db, err := app.databaseMgr.Initialize()
	if err != nil {
		return err
	}

	app.db = db

	// 执行数据库迁移
	if err := app.databaseMgr.Migrate(); err != nil {
		return err
	}

	return nil
}

// initCache 初始化缓存
func (app *Application) initCache() error {
	app.cacheMgr = NewCacheManager(app.config.Redis)

	rdb, err := app.cacheMgr.Initialize()
	if err != nil {
		return err
	}

	app.rdb = rdb
	return nil
}

// initServer 初始化服务器
func (app *Application) initServer() error {
	app.serverMgr = NewServerManager(app.config.Server)

	return app.serverMgr.Initialize(app.db, app.rdb, app.config)
}

// startServer 启动服务器
func (app *Application) startServer() error {
	// 在goroutine中启动服务器
	go func() {
		if err := app.serverMgr.Start(); err != nil {
			log.Printf("服务器启动失败: %v", err)
			app.shutdownCh <- struct{}{}
		}
	}()

	return nil
}

// waitForShutdown 等待关闭信号
func (app *Application) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		log.Println("接收到关闭信号")
	case <-app.shutdownCh:
		log.Println("应用内部触发关闭")
	}

	app.gracefulShutdown()
}

// gracefulShutdown 优雅关闭
func (app *Application) gracefulShutdown() {
	log.Println("开始优雅关闭...")

	// 5秒超时关闭
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 关闭HTTP服务器
	if app.serverMgr != nil {
		if err := app.serverMgr.Shutdown(ctx); err != nil {
			log.Printf("服务器关闭失败: %v", err)
		}
	}

	// 关闭数据库连接
	if app.databaseMgr != nil {
		if err := app.databaseMgr.Close(); err != nil {
			log.Printf("数据库关闭失败: %v", err)
		}
	}

	// 关闭Redis连接
	if app.cacheMgr != nil {
		if err := app.cacheMgr.Close(); err != nil {
			log.Printf("缓存关闭失败: %v", err)
		}
	}

	log.Println("应用已退出")
}

// GetConfig 获取配置实例
func (app *Application) GetConfig() *config.Config {
	return app.config
}

// GetDB 获取数据库实例
func (app *Application) GetDB() *gorm.DB {
	return app.db
}

// GetRedis 获取Redis实例
func (app *Application) GetRedis() *redis.Client {
	return app.rdb
}
