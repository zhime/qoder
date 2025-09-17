package app

import (
	"fmt"
	"os"

	"devops/internal/config"
	"devops/pkg/logger"
)

type ConfigManager struct {
	config *config.Config
}

func NewConfigManager() *ConfigManager {
	return &ConfigManager{}
}

func (cm *ConfigManager) Load() (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("load config failed: %w", err)
	}

	if err := cm.validate(cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	cm.config = cfg
	return cfg, nil
}

func (cm *ConfigManager) validate(cfg *config.Config) error {
	if cfg.Server.Port <= 0 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Database.DSN == "" {
		return fmt.Errorf("database DSN cannot be empty")
	}

	if cfg.Redis.Addr == "" {
		return fmt.Errorf("redis address cannot be empty")
	}

	return nil
}

func (cm *ConfigManager) GetConfig() *config.Config {
	return cm.config
}

func (cm *ConfigManager) InitLogger() error {
	if cm.config == nil {
		return fmt.Errorf("config not loaded, cannot init logger")
	}

	logger.Init(cm.config.Log)
	return nil
}

// Reload 重新加载配置文件
func (cm *ConfigManager) Reload() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("reload config failed: %w", err)
	}

	if err := cm.validate(cfg); err != nil {
		return fmt.Errorf("config validation failed after reload: %w", err)
	}

	cm.config = cfg
	return nil
}

// ValidateEnvironment 验证运行环境
func (cm *ConfigManager) ValidateEnvironment() error {
	if cm.config == nil {
		return fmt.Errorf("config not loaded")
	}

	// 检查必要的环境变量
	if cm.config.Server.Mode == "production" {
		if os.Getenv("JWT_SECRET") == "" {
			return fmt.Errorf("JWT_SECRET environment variable is required in production mode")
		}
	}

	return nil
}

// IsProduction 检查是否为生产环境
func (cm *ConfigManager) IsProduction() bool {
	if cm.config == nil {
		return false
	}
	return cm.config.Server.Mode == "release" || cm.config.Server.Mode == "production"
}

// IsDevelopment 检查是否为开发环境
func (cm *ConfigManager) IsDevelopment() bool {
	if cm.config == nil {
		return true
	}
	return cm.config.Server.Mode == "debug" || cm.config.Server.Mode == "development"
}

// GetVersion 获取配置版本信息
func (cm *ConfigManager) GetVersion() string {
	return "1.0.0" // 配置管理器版本
}

// String 实现Stringer接口
func (cm *ConfigManager) String() string {
	if cm.config == nil {
		return "ConfigManager{config: nil}"
	}
	return fmt.Sprintf("ConfigManager{port: %d, mode: %s}", cm.config.Server.Port, cm.config.Server.Mode)
}

// IsConfigLoaded 检查配置是否已加载
func (cm *ConfigManager) IsConfigLoaded() bool {
	return cm.config != nil
}