package monitor

import (
	"context"
	"fmt"
	"sync"
	"time"

	"devops/pkg/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Service 监控服务
type Service struct {
	db         *gorm.DB
	rdb        *redis.Client
	cache      *cache.CacheService
	keys       *cache.CacheKeys
	collectors map[uint]*Collector
	mu         sync.RWMutex
	stopChan   chan struct{}
	wg         sync.WaitGroup
}

// NewService 创建监控服务
func NewService(db *gorm.DB, rdb *redis.Client) *Service {
	cacheService := cache.NewCacheService(rdb, "devops")
	return &Service{
		db:         db,
		rdb:        rdb,
		cache:      cacheService,
		keys:       cache.NewCacheKeys(),
		collectors: make(map[uint]*Collector),
		stopChan:   make(chan struct{}),
	}
}

// StartMonitoring 启动监控
func (s *Service) StartMonitoring(ctx context.Context, interval time.Duration) error {
	// 获取所有需要监控的服务器
	servers, err := s.getActiveServers()
	if err != nil {
		return fmt.Errorf("获取服务器列表失败: %w", err)
	}

	// 为每台服务器创建收集器
	s.mu.Lock()
	for _, serverID := range servers {
		if _, exists := s.collectors[serverID]; !exists {
			s.collectors[serverID] = NewCollector(serverID)
		}
	}
	s.mu.Unlock()

	// 启动定时收集任务
	s.wg.Add(1)
	go s.collectMetricsLoop(ctx, interval)

	return nil
}

// StopMonitoring 停止监控
func (s *Service) StopMonitoring() {
	close(s.stopChan)
	s.wg.Wait()
}

// collectMetricsLoop 监控数据收集循环
func (s *Service) collectMetricsLoop(ctx context.Context, interval time.Duration) {
	defer s.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		case <-ticker.C:
			s.collectAllMetrics(ctx)
		}
	}
}

// collectAllMetrics 收集所有服务器的监控数据
func (s *Service) collectAllMetrics(ctx context.Context) {
	s.mu.RLock()
	collectors := make(map[uint]*Collector)
	for id, collector := range s.collectors {
		collectors[id] = collector
	}
	s.mu.RUnlock()

	// 并发收集各服务器监控数据
	var wg sync.WaitGroup
	for serverID, collector := range collectors {
		wg.Add(1)
		go func(sid uint, c *Collector) {
			defer wg.Done()
			s.collectServerMetrics(ctx, sid, c)
		}(serverID, collector)
	}

	wg.Wait()
}

// collectServerMetrics 收集单台服务器的监控数据
func (s *Service) collectServerMetrics(ctx context.Context, serverID uint, collector *Collector) {
	// 收集系统指标
	metrics, err := collector.CollectSystemMetrics(ctx)
	if err != nil {
		fmt.Printf("收集服务器%d 监控数据失败: %v\n", serverID, err)
		return
	}

	// 存储到缓存
	metricsKey := s.keys.ServerMetrics(serverID)
	if err := s.cache.Set(ctx, metricsKey, metrics, cache.TTLServerMetrics); err != nil {
		fmt.Printf("存储服务器%d 监控数据到缓存失败: %v\n", serverID, err)
	}

	// 检查告警规则
	s.checkAlerts(ctx, serverID, metrics)

	// 可选：持久化到数据库（用于历史数据分析）
	// s.persistMetrics(ctx, metrics)
}

// GetServerMetrics 获取服务器监控数据
func (s *Service) GetServerMetrics(ctx context.Context, serverID uint) (*SystemMetrics, error) {
	metricsKey := s.keys.ServerMetrics(serverID)

	var metrics SystemMetrics
	if err := s.cache.Get(ctx, metricsKey, &metrics); err != nil {
		return nil, fmt.Errorf("获取服务器监控数据失败: %w", err)
	}

	return &metrics, nil
}

// GetServerStatus 获取服务器在线状态
func (s *Service) GetServerStatus(ctx context.Context, serverID uint) (string, error) {
	// 检查最近的监控数据时间
	metrics, err := s.GetServerMetrics(ctx, serverID)
	if err != nil {
		return "offline", nil
	}

	// 如果监控数据超过5分钟，认为离线
	if time.Since(metrics.Timestamp) > 5*time.Minute {
		return "offline", nil
	}

	return "online", nil
}

// AddServer 添加服务器到监控
func (s *Service) AddServer(serverID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.collectors[serverID]; !exists {
		s.collectors[serverID] = NewCollector(serverID)
	}
}

// RemoveServer 从监控中移除服务器
func (s *Service) RemoveServer(serverID uint) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.collectors, serverID)

	// 清除缓存
	ctx := context.Background()
	metricsKey := s.keys.ServerMetrics(serverID)
	s.cache.Delete(ctx, metricsKey)
}

// GetSystemStats 获取系统统计信息
func (s *Service) GetSystemStats(ctx context.Context) (map[string]interface{}, error) {
	// 尝试从缓存获取
	statsKey := s.keys.SystemStats()
	var stats map[string]interface{}
	if err := s.cache.Get(ctx, statsKey, &stats); err == nil {
		return stats, nil
	}

	// 重新计算统计信息
	s.mu.RLock()
	totalServers := len(s.collectors)
	s.mu.RUnlock()

	onlineServers := 0
	offlineServers := 0

	// 检查每台服务器状态
	for serverID := range s.collectors {
		status, _ := s.GetServerStatus(ctx, serverID)
		if status == "online" {
			onlineServers++
		} else {
			offlineServers++
		}
	}

	stats = map[string]interface{}{
		"total_servers":   totalServers,
		"online_servers":  onlineServers,
		"offline_servers": offlineServers,
		"timestamp":       time.Now(),
	}

	// 缓存统计信息
	s.cache.Set(ctx, statsKey, stats, 5*time.Minute)

	return stats, nil
}

// checkAlerts 检查告警规则
func (s *Service) checkAlerts(ctx context.Context, serverID uint, metrics *SystemMetrics) {
	// 这里实现告警规则检查逻辑
	// 简化实现，实际项目中需要从数据库获取告警规则
	// CPU使用率告警
	if metrics.CPU.Usage > 80 {
		s.triggerAlert(ctx, serverID, "cpu", metrics.CPU.Usage, 80, "CPU使用率过高")
	}

	// 内存使用率告警
	if metrics.Memory.Usage > 85 {
		s.triggerAlert(ctx, serverID, "memory", metrics.Memory.Usage, 85, "内存使用率过高")
	}

	// 磁盘使用率告警
	for _, partition := range metrics.Disk.Partitions {
		if partition.Usage > 90 {
			s.triggerAlert(ctx, serverID, "disk", partition.Usage, 90,
				fmt.Sprintf("磁盘 %s 使用率过高", partition.Mountpoint))
		}
	}
}

// triggerAlert 触发告警
func (s *Service) triggerAlert(ctx context.Context, serverID uint, metricType string, currentValue, threshold float64, message string) {
	alert := Alert{
		ServerID:     serverID,
		MetricType:   metricType,
		CurrentValue: currentValue,
		Threshold:    threshold,
		Status:       "firing",
		Message:      message,
		FiredAt:      time.Now(),
	}

	// 存储告警到缓存
	alertKey := fmt.Sprintf("alert:%d:%s", serverID, metricType)
	s.cache.Set(ctx, alertKey, alert, 24*time.Hour)

	// 这里可以实现告警通知逻辑
	fmt.Printf("告警: 服务器%d %s\n", serverID, message)
}

// getActiveServers 获取需要监控的活跃服务器列表
func (s *Service) getActiveServers() ([]uint, error) {
	// 简化实现，实际应该从数据库查询
	// var servers []uint
	// err := s.db.Model(&model.Server{}).Where("status = ?", 1).Pluck("id", &servers).Error
	// return servers, err

	// 暂时返回示例数据
	return []uint{1}, nil
}
