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

// UserService 用户服务
type UserService struct {
	db    *gorm.DB
	rdb   *redis.Client
	cache *cache.CacheService
	keys  *cache.CacheKeys
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB, rdb *redis.Client) *UserService {
	cacheService := cache.NewCacheService(rdb, "devops")
	return &UserService{
		db:    db,
		rdb:   rdb,
		cache: cacheService,
		keys:  cache.NewCacheKeys(),
	}
}

// Login 用户登录验证
func (s *UserService) Login(username, password string) (*model.User, error) {
	var user model.User

	// 根据用户名或邮箱查询用户
	err := s.db.Where("username = ? OR email = ?", username, username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if !auth.CheckPassword(password, user.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}

// GetByID 根据ID获取用户
func (s *UserService) GetByID(id uint) (*model.User, error) {
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
			return nil, errors.New("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 将结果存入缓存
	s.cache.Set(ctx, cacheKey, &user, cache.TTLUserInfo)

	return &user, nil
}

// Create 创建用户
func (s *UserService) Create(username, email, password, role string) (*model.User, error) {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if count > 0 {
		return nil, errors.New("用户名或邮箱已存在")
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

	return user, nil
}

// Update 更新用户
func (s *UserService) Update(id uint, updates map[string]interface{}) error {
	// 如果包含密码，需要加密
	if password, ok := updates["password"].(string); ok {
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
		return errors.New("用户不存在")
	}

	return nil
}

// Delete 删除用户（软删除）
func (s *UserService) Delete(id uint) error {
	result := s.db.Delete(&model.User{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除用户失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("用户不存在")
	}

	return nil
}

// List 获取用户列表
func (s *UserService) List(page, pageSize int) ([]model.User, int64, error) {
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

	return users, total, nil
}
