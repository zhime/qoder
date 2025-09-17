package cache

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setupTestRedis(t *testing.T) *redis.Client {
	// 使用内存Redis进行测试，或者跳过测试
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   15, // 使用测试数据库
	})

	// 测试连接
	ctx := context.Background()
	if err := rdb.Ping(ctx).Err(); err != nil {
		t.Skip("Redis not available, skipping cache tests")
	}

	// 清空测试数据库
	rdb.FlushDB(ctx)

	return rdb
}

func TestCacheService(t *testing.T) {
	rdb := setupTestRedis(t)
	defer rdb.Close()

	cache := NewCacheService(rdb, "test")
	ctx := context.Background()

	t.Run("SetAndGet", func(t *testing.T) {
		key := "test_key"
		value := map[string]interface{}{
			"name": "test",
			"age":  25,
		}

		// 设置缓存
		err := cache.Set(ctx, key, value, time.Minute)
		assert.NoError(t, err)

		// 获取缓存
		var result map[string]interface{}
		err = cache.Get(ctx, key, &result)
		assert.NoError(t, err)
		assert.Equal(t, "test", result["name"])
		assert.Equal(t, float64(25), result["age"]) // JSON数字默认为float64
	})

	t.Run("Exists", func(t *testing.T) {
		key := "exist_key"

		// 键不存在
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)

		// 设置键
		err = cache.Set(ctx, key, "value", time.Minute)
		assert.NoError(t, err)

		// 键存在
		exists, err = cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("Delete", func(t *testing.T) {
		key := "delete_key"

		// 设置键
		err := cache.Set(ctx, key, "value", time.Minute)
		assert.NoError(t, err)

		// 删除键
		err = cache.Delete(ctx, key)
		assert.NoError(t, err)

		// 检查键是否被删除
		exists, err := cache.Exists(ctx, key)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("HashOperations", func(t *testing.T) {
		key := "hash_key"
		fields := map[string]interface{}{
			"field1": "value1",
			"field2": "value2",
		}

		// 设置哈希
		err := cache.SetHash(ctx, key, fields)
		assert.NoError(t, err)

		// 获取单个字段
		value, err := cache.GetHash(ctx, key, "field1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", value)

		// 获取所有字段
		all, err := cache.GetAllHash(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, "value1", all["field1"])
		assert.Equal(t, "value2", all["field2"])
	})

	t.Run("IncrementOperations", func(t *testing.T) {
		key := "counter_key"

		// 递增
		count, err := cache.Increment(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), count)

		// 按值递增
		count, err = cache.IncrementBy(ctx, key, 5)
		assert.NoError(t, err)
		assert.Equal(t, int64(6), count)

		// 递减
		count, err = cache.Decrement(ctx, key)
		assert.NoError(t, err)
		assert.Equal(t, int64(5), count)
	})

	t.Run("SetOperations", func(t *testing.T) {
		key := "set_key"

		// 添加到集合
		err := cache.AddToSet(ctx, key, "member1", "member2")
		assert.NoError(t, err)

		// 检查成员
		isMember, err := cache.IsSetMember(ctx, key, "member1")
		assert.NoError(t, err)
		assert.True(t, isMember)

		// 获取所有成员
		members, err := cache.GetSetMembers(ctx, key)
		assert.NoError(t, err)
		assert.Contains(t, members, "member1")
		assert.Contains(t, members, "member2")

		// 从集合删除
		err = cache.RemoveFromSet(ctx, key, "member1")
		assert.NoError(t, err)

		// 检查成员是否被删除
		isMember, err = cache.IsSetMember(ctx, key, "member1")
		assert.NoError(t, err)
		assert.False(t, isMember)
	})

	t.Run("SortedSetOperations", func(t *testing.T) {
		key := "zset_key"

		// 添加到有序集合
		err := cache.AddToSortedSet(ctx, key, 100, "member1")
		assert.NoError(t, err)
		err = cache.AddToSortedSet(ctx, key, 200, "member2")
		assert.NoError(t, err)

		// 获取范围
		members, err := cache.GetSortedSetRange(ctx, key, 0, -1)
		assert.NoError(t, err)
		assert.Equal(t, []string{"member1", "member2"}, members)
	})
}

func TestCacheKeys(t *testing.T) {
	keys := NewCacheKeys()

	t.Run("UserKeys", func(t *testing.T) {
		userID := uint(123)

		userInfoKey := keys.UserInfo(userID)
		assert.Equal(t, "user:info:123", userInfoKey)

		permissionsKey := keys.UserPermissions(userID)
		assert.Equal(t, "user:permissions:123", permissionsKey)
	})

	t.Run("ServerKeys", func(t *testing.T) {
		serverID := uint(456)
		userID := uint(123)

		serverInfoKey := keys.ServerInfo(serverID)
		assert.Equal(t, "server:info:456", serverInfoKey)

		serverListKey := keys.ServerList(userID)
		assert.Equal(t, "server:list:123", serverListKey)

		metricsKey := keys.ServerMetrics(serverID)
		assert.Equal(t, "metrics:456", metricsKey)
	})

	t.Run("TaskKeys", func(t *testing.T) {
		taskID := uint(789)

		nextRunKey := keys.TaskNextRun()
		assert.Equal(t, "task:next_run", nextRunKey)

		executionKey := keys.TaskExecution(taskID)
		assert.Equal(t, "task:execution:789", executionKey)

		lockKey := keys.TaskLock(taskID)
		assert.Equal(t, "task:lock:789", lockKey)
	})
}
