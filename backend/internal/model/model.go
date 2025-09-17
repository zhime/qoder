package model

import (
	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Server{},
		&Deployment{},
		&DeploymentLog{},
		&Task{},
		&TaskExecution{},
	)
}