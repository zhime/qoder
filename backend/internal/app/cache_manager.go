package app

import (
	"context"
	"fmt"
	"time"

	"devops/internal/config"
	"devops/pkg/redis"

	redisPkg "github.com/redis/go-redis/v9"
)

type CacheManager struct {
	rdb    *redisPkg.Client
	config config.Redis
}

func NewCacheManager(cfg config.Redis) *CacheManager {
	return &CacheManager{
		config: cfg,
	}
}

func (cm *CacheManager) Initialize() (*redisPkg.Client, error) {
	rdb, err := redis.Init(cm.config)
	if err != nil {
		return nil, fmt.Errorf("redis初始化失败: %w", err)
	}

	cm.rdb = rdb

	if err := cm.ping(); err != nil {
		return nil, fmt.Errorf("redis连接测试失败: %w", err)
	}

	return rdb, nil
}

func (cm *CacheManager) ping() error {
	if cm.rdb == nil {
		return fmt.Errorf("redis客户端未初始化")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result := cm.rdb.Ping(ctx)
	return result.Err()
}

func (cm *CacheManager) GetClient() *redisPkg.Client {
	return cm.rdb
}

func (cm *CacheManager) Close() error {
	if cm.rdb == nil {
		return nil
	}

	return cm.rdb.Close()
}

func (cm *CacheManager) IsHealthy() bool {
	if cm.rdb == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := cm.rdb.Ping(ctx).Err()
	return err == nil
}
