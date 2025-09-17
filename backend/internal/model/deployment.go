package model

import (
	"time"

	"gorm.io/gorm"
)

// Deployment 部署模型
type Deployment struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	Name       string         `gorm:"size:100;not null" json:"name"`
	ServerID   uint           `gorm:"index;not null" json:"server_id"`
	Server     Server         `gorm:"foreignKey:ServerID" json:"server,omitempty"`
	Repository string         `gorm:"size:200" json:"repository"`
	Branch     string         `gorm:"size:50;default:main" json:"branch"`
	Path       string         `gorm:"size:200" json:"path"`
	Script     string         `gorm:"type:text" json:"script"`
	Status     int            `gorm:"default:0" json:"status"` // 0:待部署 1:部署中 2:部署成功 3:部署失败
	CreatedBy  uint           `gorm:"index;not null" json:"created_by"`
	User       User           `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Logs []DeploymentLog `gorm:"foreignKey:DeploymentID" json:"-"`
}

// TableName 设置表名
func (Deployment) TableName() string {
	return "deployments"
}

// DeploymentLog 部署日志模型
type DeploymentLog struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	DeploymentID uint           `gorm:"index;not null" json:"deployment_id"`
	Deployment   Deployment     `gorm:"foreignKey:DeploymentID" json:"-"`
	Level        string         `gorm:"size:20" json:"level"` // info, warn, error
	Message      string         `gorm:"type:text" json:"message"`
	CreatedAt    time.Time      `json:"created_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 设置表名
func (DeploymentLog) TableName() string {
	return "deployment_logs"
}