package service

import (
	"context"
	"errors"
	"fmt"

	"devops-platform/internal/model"
	"devops-platform/pkg/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// ServerService 服务器服务
type ServerService struct {
	db    *gorm.DB
	rdb   *redis.Client
	cache *cache.CacheService
	keys  *cache.CacheKeys
}

// NewServerService 创建服务器服务
func NewServerService(db *gorm.DB, rdb *redis.Client) *ServerService {
	cacheService := cache.NewCacheService(rdb, "devops")
	return &ServerService{
		db:    db,
		rdb:   rdb,
		cache: cacheService,
		keys:  cache.NewCacheKeys(),
	}
}

// GetByID 根据ID获取服务器
func (s *ServerService) GetByID(id uint) (*model.Server, error) {
	ctx := context.Background()
	
	// 先从缓存查找
	var server model.Server
	cacheKey := s.keys.ServerInfo(id)
	if err := s.cache.Get(ctx, cacheKey, &server); err == nil {
		return &server, nil
	}
	
	// 缓存中没有，从数据库查询
	err := s.db.First(&server, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("服务器不存在")
		}
		return nil, fmt.Errorf("查询服务器失败: %w", err)
	}
	
	// 将结果存入缓存
	s.cache.Set(ctx, cacheKey, &server, cache.TTLServerList)
	
	return &server, nil
}

// List 获取服务器列表
func (s *ServerService) List(userID uint) ([]model.Server, error) {
	ctx := context.Background()
	
	// 尝试从缓存获取
	cacheKey := s.keys.ServerList(userID)
	var servers []model.Server
	if err := s.cache.Get(ctx, cacheKey, &servers); err == nil {
		return servers, nil
	}
	
	// 从数据库查询
	// 这里可以根据用户权限过滤服务器
	err := s.db.Find(&servers).Error
	if err != nil {
		return nil, fmt.Errorf("查询服务器列表失败: %w", err)
	}
	
	// 缓存结果
	s.cache.Set(ctx, cacheKey, servers, cache.TTLServerList)
	
	return servers, nil
}

// Create 创建服务器
func (s *ServerService) Create(server *model.Server) error {
	if err := s.db.Create(server).Error; err != nil {
		return fmt.Errorf("创建服务器失败: %w", err)
	}
	
	// 清除相关缓存
	ctx := context.Background()
	s.InvalidateServerListCache(ctx)
	
	return nil
}

// Update 更新服务器
func (s *ServerService) Update(id uint, updates map[string]interface{}) error {
	result := s.db.Model(&model.Server{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新服务器失败: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.New("服务器不存在")
	}
	
	// 清除缓存
	ctx := context.Background()
	s.InvalidateServerCache(ctx, id)
	
	return nil
}

// Delete 删除服务器
func (s *ServerService) Delete(id uint) error {
	result := s.db.Delete(&model.Server{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除服务器失败: %w", result.Error)
	}
	
	if result.RowsAffected == 0 {
		return errors.New("服务器不存在")
	}
	
	// 清除缓存
	ctx := context.Background()
	s.InvalidateServerCache(ctx, id)
	
	return nil
}

// UpdateStatus 更新服务器状态
func (s *ServerService) UpdateStatus(id uint, status int) error {
	ctx := context.Background()
	
	// 更新数据库
	if err := s.db.Model(&model.Server{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		return fmt.Errorf("更新服务器状态失败: %w", err)
	}
	
	// 更新缓存中的状态
	statusKey := s.keys.ServerStatus(id)
	s.cache.Set(ctx, statusKey, status, cache.TTLServerList)
	
	// 清除服务器信息缓存以便重新加载
	s.InvalidateServerCache(ctx, id)
	
	return nil
}

// GetStatus 获取服务器状态
func (s *ServerService) GetStatus(id uint) (int, error) {
	ctx := context.Background()
	
	// 先从缓存获取状态
	statusKey := s.keys.ServerStatus(id)
	var status int
	if err := s.cache.Get(ctx, statusKey, &status); err == nil {
		return status, nil
	}
	
	// 从数据库获取
	var server model.Server
	if err := s.db.Select("status").First(&server, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("服务器不存在")
		}
		return 0, fmt.Errorf("查询服务器状态失败: %w", err)
	}
	
	// 缓存状态
	s.cache.Set(ctx, statusKey, server.Status, cache.TTLServerList)
	
	return server.Status, nil
}

// StoreMetrics 存储服务器监控数据
func (s *ServerService) StoreMetrics(serverID uint, metrics interface{}) error {
	ctx := context.Background()
	metricsKey := s.keys.ServerMetrics(serverID)
	return s.cache.Set(ctx, metricsKey, metrics, cache.TTLServerMetrics)
}

// GetMetrics 获取服务器监控数据
func (s *ServerService) GetMetrics(serverID uint) (interface{}, error) {
	ctx := context.Background()
	metricsKey := s.keys.ServerMetrics(serverID)
	
	var metrics interface{}
	if err := s.cache.Get(ctx, metricsKey, &metrics); err != nil {
		return nil, fmt.Errorf("获取服务器监控数据失败: %w", err)
	}
	
	return metrics, nil
}

// InvalidateServerCache 清除服务器缓存
func (s *ServerService) InvalidateServerCache(ctx context.Context, serverID uint) {
	keys := []string{
		s.keys.ServerInfo(serverID),
		s.keys.ServerStatus(serverID),
		s.keys.ServerMetrics(serverID),
	}
	
	for _, key := range keys {
		s.cache.Delete(ctx, key)
	}
	
	// 清除列表缓存
	s.InvalidateServerListCache(ctx)
}

// InvalidateServerListCache 清除服务器列表缓存
func (s *ServerService) InvalidateServerListCache(ctx context.Context) {
	// 实际项目中可以维护用户ID列表，这里简化处理
	// 清除所有可能的服务器列表缓存
	s.cache.Delete(ctx, "server:list:*")
}