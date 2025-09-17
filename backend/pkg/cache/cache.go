package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService 缓存服务
type CacheService struct {
	client *redis.Client
	prefix string
}

// NewCacheService 创建缓存服务
func NewCacheService(client *redis.Client, prefix string) *CacheService {
	return &CacheService{
		client: client,
		prefix: prefix,
	}
}

// buildKey 构建缓存键
func (c *CacheService) buildKey(key string) string {
	return fmt.Sprintf("%s:%s", c.prefix, key)
}

// Set 设置缓存
func (c *CacheService) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("序列化数据失败: %w", err)
	}

	return c.client.Set(ctx, c.buildKey(key), data, expiration).Err()
}

// Get 获取缓存
func (c *CacheService) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := c.client.Get(ctx, c.buildKey(key)).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

// Delete 删除缓存
func (c *CacheService) Delete(ctx context.Context, keys ...string) error {
	cacheKeys := make([]string, len(keys))
	for i, key := range keys {
		cacheKeys[i] = c.buildKey(key)
	}
	return c.client.Del(ctx, cacheKeys...).Err()
}

// Exists 检查缓存是否存在
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, c.buildKey(key)).Result()
	return count > 0, err
}

// SetHash 设置哈希缓存
func (c *CacheService) SetHash(ctx context.Context, key string, fields map[string]interface{}) error {
	return c.client.HMSet(ctx, c.buildKey(key), fields).Err()
}

// GetHash 获取哈希缓存
func (c *CacheService) GetHash(ctx context.Context, key string, field string) (string, error) {
	return c.client.HGet(ctx, c.buildKey(key), field).Result()
}

// GetAllHash 获取所有哈希字段
func (c *CacheService) GetAllHash(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, c.buildKey(key)).Result()
}

// DeleteHashField 删除哈希字段
func (c *CacheService) DeleteHashField(ctx context.Context, key string, fields ...string) error {
	return c.client.HDel(ctx, c.buildKey(key), fields...).Err()
}

// SetExpire 设置过期时间
func (c *CacheService) SetExpire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, c.buildKey(key), expiration).Err()
}

// Increment 递增
func (c *CacheService) Increment(ctx context.Context, key string) (int64, error) {
	return c.client.Incr(ctx, c.buildKey(key)).Result()
}

// IncrementBy 按指定值递增
func (c *CacheService) IncrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.IncrBy(ctx, c.buildKey(key), value).Result()
}

// Decrement 递减
func (c *CacheService) Decrement(ctx context.Context, key string) (int64, error) {
	return c.client.Decr(ctx, c.buildKey(key)).Result()
}

// DecrementBy 按指定值递减
func (c *CacheService) DecrementBy(ctx context.Context, key string, value int64) (int64, error) {
	return c.client.DecrBy(ctx, c.buildKey(key), value).Result()
}

// AddToSet 添加到集合
func (c *CacheService) AddToSet(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SAdd(ctx, c.buildKey(key), members...).Err()
}

// RemoveFromSet 从集合删除
func (c *CacheService) RemoveFromSet(ctx context.Context, key string, members ...interface{}) error {
	return c.client.SRem(ctx, c.buildKey(key), members...).Err()
}

// IsSetMember 检查是否为集合成员
func (c *CacheService) IsSetMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, c.buildKey(key), member).Result()
}

// GetSetMembers 获取集合所有成员
func (c *CacheService) GetSetMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, c.buildKey(key)).Result()
}

// AddToSortedSet 添加到有序集合
func (c *CacheService) AddToSortedSet(ctx context.Context, key string, score float64, member interface{}) error {
	return c.client.ZAdd(ctx, c.buildKey(key), redis.Z{
		Score:  score,
		Member: member,
	}).Err()
}

// GetSortedSetRange 获取有序集合范围
func (c *CacheService) GetSortedSetRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, c.buildKey(key), start, stop).Result()
}

// RemoveFromSortedSet 从有序集合删除
func (c *CacheService) RemoveFromSortedSet(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, c.buildKey(key), members...).Err()
}

// FlushAll 清空所有缓存（谨慎使用）
func (c *CacheService) FlushAll(ctx context.Context) error {
	return c.client.FlushAll(ctx).Err()
}
