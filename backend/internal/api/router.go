package api

import (
	"devops-platform/internal/config"
	"devops-platform/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// NewRouter 创建新的路由器
func NewRouter(db *gorm.DB, rdb *redis.Client, cfg *config.Config) *gin.Engine {
	router := gin.New()

	// 全局中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORS())

	// 创建处理器
	authHandler := NewAuthHandler(db, rdb, cfg.JWT.Secret)
	userHandler := NewUserHandler(db, rdb)
	monitorHandler := NewMonitorHandler(db, rdb)

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Service is running",
		})
	})

	// API路由组
	api := router.Group("/api")
	{
		// 认证路由（无需JWT验证）
		auth := api.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 需要JWT验证的路由
		protected := api.Group("")
		protected.Use(middleware.JWTAuth(cfg.JWT.Secret))
		{
			// 认证用户信息
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.POST("/auth/logout", authHandler.Logout)

			// 用户相关
			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.GetByID)
				// 管理员权限
				adminUsers := users.Group("")
				adminUsers.Use(middleware.RequireRole("admin"))
				{
					adminUsers.POST("", userHandler.Create)
					adminUsers.PUT("/:id", userHandler.Update)
					adminUsers.DELETE("/:id", userHandler.Delete)
				}
			}

			// 服务器相关
			servers := protected.Group("/servers")
			{
				servers.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "servers list"})
				})
				servers.POST("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "create server"})
				})
			}

			// 部署相关
			deployments := protected.Group("/deployments")
			{
				deployments.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "deployments list"})
				})
				deployments.POST("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "create deployment"})
				})
			}

			// 任务相关
			tasks := protected.Group("/tasks")
			{
				tasks.GET("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "tasks list"})
				})
				tasks.POST("", func(c *gin.Context) {
					c.JSON(200, gin.H{"message": "create task"})
				})
			}

			// 监控相关
			monitor := protected.Group("/monitor")
			{
				monitor.GET("/dashboard", monitorHandler.GetDashboardData)
				monitor.GET("/stats", monitorHandler.GetSystemStats)
				monitor.GET("/servers/:id/metrics", monitorHandler.GetServerMetrics)
				monitor.GET("/servers/:id/status", monitorHandler.GetServerStatus)
				monitor.GET("/servers/:id/history", monitorHandler.GetServerHistory)
				monitor.POST("/servers", monitorHandler.AddServerToMonitor)
				monitor.DELETE("/servers/:id", monitorHandler.RemoveServerFromMonitor)
			}
		}
	}

	return router
}