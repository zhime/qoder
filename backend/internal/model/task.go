package model

import (
	"time"

	"gorm.io/gorm"
)

// Task 任务模型
type Task struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"size:100;not null" json:"name"`
	Command   string         `gorm:"type:text;not null" json:"command"`
	CronExpr  string         `gorm:"size:50" json:"cron_expr"`
	ServerID  uint           `gorm:"index;not null" json:"server_id"`
	Server    Server         `gorm:"foreignKey:ServerID" json:"server,omitempty"`
	Status    int            `gorm:"default:1" json:"status"` // 0:禁用 1:启用
	LastRun   *time.Time     `json:"last_run"`
	NextRun   *time.Time     `json:"next_run"`
	CreatedBy uint           `gorm:"index;not null" json:"created_by"`
	User      User           `gorm:"foreignKey:CreatedBy" json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Executions []TaskExecution `gorm:"foreignKey:TaskID" json:"-"`
}

// TableName 设置表名
func (Task) TableName() string {
	return "tasks"
}

// TaskExecution 任务执行记录模型
type TaskExecution struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TaskID    uint           `gorm:"index;not null" json:"task_id"`
	Task      Task           `gorm:"foreignKey:TaskID" json:"-"`
	Status    int            `gorm:"default:0" json:"status"` // 0:运行中 1:成功 2:失败
	Output    string         `gorm:"type:text" json:"output"`
	Error     string         `gorm:"type:text" json:"error"`
	StartTime time.Time      `json:"start_time"`
	EndTime   *time.Time     `json:"end_time"`
	Duration  int            `json:"duration"` // 执行时长（秒）
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 设置表名
func (TaskExecution) TableName() string {
	return "task_executions"
}
