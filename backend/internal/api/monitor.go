package api

import (
	"net/http"
	"strconv"

	"devops-platform/internal/monitor"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// MonitorHandler 监控处理器
type MonitorHandler struct {
	monitorService *monitor.Service
}

// NewMonitorHandler 创建监控处理器
func NewMonitorHandler(db *gorm.DB, rdb *redis.Client) *MonitorHandler {
	return &MonitorHandler{
		monitorService: monitor.NewService(db, rdb),
	}
}

// GetServerMetrics 获取服务器监控数据
func (h *MonitorHandler) GetServerMetrics(c *gin.Context) {
	serverIDStr := c.Param("id")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的服务器ID",
		})
		return
	}

	metrics, err := h.monitorService.GetServerMetrics(c.Request.Context(), uint(serverID))
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Code:    404,
			Message: "获取监控数据失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    metrics,
	})
}

// GetServerStatus 获取服务器状态
func (h *MonitorHandler) GetServerStatus(c *gin.Context) {
	serverIDStr := c.Param("id")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的服务器ID",
		})
		return
	}

	status, err := h.monitorService.GetServerStatus(c.Request.Context(), uint(serverID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取服务器状态失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data: map[string]interface{}{
			"server_id": serverID,
			"status":    status,
		},
	})
}

// GetSystemStats 获取系统统计信息
func (h *MonitorHandler) GetSystemStats(c *gin.Context) {
	stats, err := h.monitorService.GetSystemStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取系统统计失败: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    stats,
	})
}

// GetDashboardData 获取仪表盘数据
func (h *MonitorHandler) GetDashboardData(c *gin.Context) {
	ctx := c.Request.Context()

	// 获取系统统计
	stats, err := h.monitorService.GetSystemStats(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Code:    500,
			Message: "获取系统统计失败: " + err.Error(),
		})
		return
	}

	// 构建仪表盘数据
	dashboardData := map[string]interface{}{
		"stats": stats,
		"alerts": []map[string]interface{}{
			{
				"id":      1,
				"level":   "warning",
				"message": "服务器 server-01 CPU使用率过高 (85%)",
				"time":    "2024-01-15 10:30:00",
			},
			{
				"id":      2,
				"level":   "error",
				"message": "服务器 server-02 磁盘空间不足 (95%)",
				"time":    "2024-01-15 10:25:00",
			},
		},
		"recent_activities": []map[string]interface{}{
			{
				"id":          1,
				"type":        "deployment",
				"description": "部署 demo-app 到生产环境",
				"status":      "success",
				"time":        "2024-01-15 10:15:00",
			},
			{
				"id":          2,
				"type":        "task",
				"description": "执行数据库备份任务",
				"status":      "success",
				"time":        "2024-01-15 09:00:00",
			},
		},
	}

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    dashboardData,
	})
}

// AddServerToMonitor 添加服务器到监控
func (h *MonitorHandler) AddServerToMonitor(c *gin.Context) {
	var req struct {
		ServerID uint `json:"server_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "请求参数错误: " + err.Error(),
		})
		return
	}

	h.monitorService.AddServer(req.ServerID)

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "添加监控成功",
	})
}

// RemoveServerFromMonitor 从监控中移除服务器
func (h *MonitorHandler) RemoveServerFromMonitor(c *gin.Context) {
	serverIDStr := c.Param("id")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的服务器ID",
		})
		return
	}

	h.monitorService.RemoveServer(uint(serverID))

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "移除监控成功",
	})
}

// GetServerHistory 获取服务器历史监控数据
func (h *MonitorHandler) GetServerHistory(c *gin.Context) {
	serverIDStr := c.Param("id")
	serverID, err := strconv.ParseUint(serverIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Code:    400,
			Message: "无效的服务器ID",
		})
		return
	}

	// 获取查询参数
	timeRange := c.DefaultQuery("time_range", "1h") // 1h, 6h, 24h, 7d
	metric := c.DefaultQuery("metric", "cpu")        // cpu, memory, disk, network

	// 这里应该从数据库或时序数据库获取历史数据
	// 暂时返回模拟数据
	historyData := generateMockHistoryData(uint(serverID), timeRange, metric)

	c.JSON(http.StatusOK, Response{
		Code:    200,
		Message: "获取成功",
		Data:    historyData,
	})
}

// generateMockHistoryData 生成模拟历史数据
func generateMockHistoryData(serverID uint, timeRange, metric string) map[string]interface{} {
	// 生成模拟时间序列数据
	data := make([]map[string]interface{}, 0)

	for i := 0; i < 20; i++ {
		point := map[string]interface{}{
			"timestamp": "2024-01-15 " + strconv.Itoa(10+i/2) + ":" + strconv.Itoa((i%2)*30) + ":00",
		}

		switch metric {
		case "cpu":
			point["value"] = 20 + (i % 10) * 5
		case "memory":
			point["value"] = 40 + (i % 8) * 6
		case "disk":
			point["value"] = 60 + (i % 5) * 3
		case "network":
			point["value"] = (i % 15) * 10
		}

		data = append(data, point)
	}

	return map[string]interface{}{
		"server_id":  serverID,
		"metric":     metric,
		"time_range": timeRange,
		"data":       data,
	}
}