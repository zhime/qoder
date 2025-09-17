package service

import (
	"context"
	"errors"
	"fmt"

	"devops/internal/auth"
	"devops/internal/model"
	"devops/pkg/cache"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// UserServiceWithCache 带缓存的用户服务
type UserServiceWithCache struct {
	db    *gorm.DB
	rdb   *redis.Client
	cache *cache.CacheService
	keys  *cache.CacheKeys
}

// NewUserServiceWithCache 创建带缓存的用户服务
func NewUserServiceWithCache(db *gorm.DB, rdb *redis.Client) *UserServiceWithCache {
	cacheService := cache.NewCacheService(rdb, "devops")
	return &UserServiceWithCache{
		db:    db,
		rdb:   rdb,
		cache: cacheService,
		keys:  cache.NewCacheKeys(),
	}
}

// Login 用户登录验证
func (s *UserServiceWithCache) Login(username, password string) (*model.User, error) {
	var user model.User
	
	// 根据用户名或邮箱查询用户
	err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状�?	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !auth.CheckPassword(password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	// 缓存用户信息
	ctx := context.Background()
	cacheKey := s.keys.UserInfo(user.ID)
	s.cache.Set(ctx, cacheKey, &user, cache.TTLUserInfo)

	return &user, nil
}

// GetByID 根据ID获取用户（带缓存�?func (s *UserServiceWithCache) GetByID(id uint) (*model.User, error) {
	ctx := context.Background()
	
	// 先从缓存查找
	var user model.User
	cacheKey := s.keys.UserInfo(id)
	if err := s.cache.Get(ctx, cacheKey, &user); err == nil {
		return &user, nil
	}
	
	// 缓存中没有，从数据库查询
	err := s.db.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存�?)
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}
	
	// 将结果存入缓�?	s.cache.Set(ctx, cacheKey, &user, cache.TTLUserInfo)
	
	return &user, nil
}

// Create 创建用户
func (s *UserServiceWithCache) Create(username, email, password, role string) (*model.User, error) {
	// 检查用户名是否已存�?	var count int64
	s.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名或邮箱已存�?)
	}

	// 加密密码
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %w", err)
	}

	user := &model.User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Role:     role,
		Status:   1,
	}

	if err := s.db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 缓存新用户信�?	ctx := context.Background()
	cacheKey := s.keys.UserInfo(user.ID)
	s.cache.Set(ctx, cacheKey, user, cache.TTLUserInfo)

	// 清除相关缓存
	s.InvalidateUserListCache(ctx)

	return user, nil
}

// Update 更新用户
func (s *UserServiceWithCache) Update(id uint, updates map[string]interface{}) error {
	// 如果包含密码，需要加�?	if password, ok := updates["password"].(string); ok {
		hashedPassword, err := auth.HashPassword(password)
		if err != nil {
			return fmt.Errorf("密码加密失败: %w", err)
		}
		updates["password"] = hashedPassword
	}

	result := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("更新用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("用户不存�?)
	}

	// 清除用户缓存
	ctx := context.Background()
	s.InvalidateUserCache(ctx, id)

	return nil
}

// Delete 删除用户（软删除�?func (s *UserServiceWithCache) Delete(id uint) error {
	result := s.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("用户不存�?)
	}

	// 清除用户缓存
	ctx := context.Background()
	s.InvalidateUserCache(ctx, id)

	return nil
}

// List 获取用户列表
func (s *UserServiceWithCache) List(page, pageSize int) ([]model.User, int64, error) {
	ctx := context.Background()
	
	// 尝试从缓存获取用户列�?	cacheKey := fmt.Sprintf("user:list:%d:%d", page, pageSize)
	var cachedResult struct {
		Users []model.User `json:"users"`
		Total int64        `json:"total"`
	}
	
	if err := s.cache.Get(ctx, cacheKey, &cachedResult); err == nil {
		return cachedResult.Users, cachedResult.Total, nil
	}

	var users []model.User
	var total int64

	// 获取总数
	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 缓存结果
	cachedResult.Users = users
	cachedResult.Total = total
	s.cache.Set(ctx, cacheKey, cachedResult, cache.TTLUserInfo)

	return users, total, nil
}

// InvalidateUserCache 清除用户缓存
func (s *UserServiceWithCache) InvalidateUserCache(ctx context.Context, userID uint) {
	cacheKey := s.keys.UserInfo(userID)
	s.cache.Delete(ctx, cacheKey)
	
	// 清除相关缓存
	s.InvalidateUserListCache(ctx)
}

// InvalidateUserListCache 清除用户列表缓存
func (s *UserServiceWithCache) InvalidateUserListCache(ctx context.Context) {
	// 这里可以使用模式匹配删除所有用户列表缓�?	// 简化处理，实际项目中可以维护一个缓存键列表
	keys := []string{
		"user:list:1:10",
		"user:list:1:20",
		"user:list:1:50",
		// 可以根据实际情况添加更多
	}
	
	for _, key := range keys {
		s.cache.Delete(ctx, key)
	}
}

// GetOnlineUsers 获取在线用户数量
func (s *UserServiceWithCache) GetOnlineUsers(ctx context.Context) (int64, error) {
	return s.cache.GetSetMembers(ctx, s.keys.OnlineUsers())
}

// SetUserOnline 设置用户在线状�?func (s *UserServiceWithCache) SetUserOnline(ctx context.Context, userID uint) error {
	return s.cache.AddToSet(ctx, s.keys.OnlineUsers(), userID)
}

// SetUserOffline 设置用户离线状�?func (s *UserServiceWithCache) SetUserOffline(ctx context.Context, userID uint) error {
	return s.cache.RemoveFromSet(ctx, s.keys.OnlineUsers(), userID)
}
