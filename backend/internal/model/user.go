package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email     string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password  string         `gorm:"size:100;not null" json:"-"`
	Role      string         `gorm:"size:20;default:user" json:"role"`
	Status    int            `gorm:"default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Deployments []Deployment `gorm:"foreignKey:CreatedBy" json:"-"`
	Tasks       []Task       `gorm:"foreignKey:CreatedBy" json:"-"`
}

// TableName 设置表名
func (User) TableName() string {
	return "users"
}
