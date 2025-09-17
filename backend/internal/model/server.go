package model

import (
	"time"

	"gorm.io/gorm"
)

// Server 服务器模型
type Server struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Host        string         `gorm:"size:100;not null" json:"host"`
	Port        int            `gorm:"default:22" json:"port"`
	Username    string         `gorm:"size:50;not null" json:"username"`
	Password    string         `gorm:"size:100" json:"-"`
	PrivateKey  string         `gorm:"type:text" json:"-"`
	Status      int            `gorm:"default:1" json:"status"`
	Environment string         `gorm:"size:20" json:"environment"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Deployments []Deployment `gorm:"foreignKey:ServerID" json:"-"`
	Tasks       []Task       `gorm:"foreignKey:ServerID" json:"-"`
}

// TableName 设置表名
func (Server) TableName() string {
	return "servers"
}