package app

import (
	"fmt"
	"time"

	"devops/internal/config"
	"devops/internal/model"
	"devops/pkg/database"

	"gorm.io/gorm"
)

// DatabaseManager 数据库管理器
type DatabaseManager struct {
	db     *gorm.DB
	config config.Database
}

// NewDatabaseManager 创建数据库管理器实例
func NewDatabaseManager(cfg config.Database) *DatabaseManager {
	return &DatabaseManager{
		config: cfg,
	}
}

// Initialize 初始化数据库连接
func (dm *DatabaseManager) Initialize() (*gorm.DB, error) {
	db, err := database.Init(dm.config)
	if err != nil {
		return nil, fmt.Errorf("数据库初始化失败: %w", err)
	}

	dm.db = db

	// 设置连接池配�?	if err := dm.configureConnectionPool(); err != nil {
		return nil, fmt.Errorf("数据库连接池配置失败: %w", err)
	}

	return db, nil
}

// configureConnectionPool 配置数据库连接池
func (dm *DatabaseManager) configureConnectionPool() error {
	sqlDB, err := dm.db.DB()
	if err != nil {
		return err
	}

	// 设置最大打开连接�?	sqlDB.SetMaxOpenConns(25)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)
	// 设置连接最大生存时�?	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return nil
}

// Migrate 执行数据库迁�?func (dm *DatabaseManager) Migrate() error {
	if dm.db == nil {
		return fmt.Errorf("数据库未初始化，无法执行迁移")
	}

	if err := model.AutoMigrate(dm.db); err != nil {
		return fmt.Errorf("数据库迁移失�? %w", err)
	}

	return nil
}

// GetDB 获取数据库实�?func (dm *DatabaseManager) GetDB() *gorm.DB {
	return dm.db
}

// Close 关闭数据库连�?func (dm *DatabaseManager) Close() error {
	if dm.db == nil {
		return nil
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Close()
}

// Ping 检测数据库连接状�?func (dm *DatabaseManager) Ping() error {
	if dm.db == nil {
		return fmt.Errorf("数据库未初始�?)
	}

	sqlDB, err := dm.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.Ping()
}
