package monitor

import (
	"time"
)

// SystemMetrics 系统监控指标
type SystemMetrics struct {
	ServerID    uint      `json:"server_id"`
	Timestamp   time.Time `json:"timestamp"`
	CPU         CPUMetrics `json:"cpu"`
	Memory      MemoryMetrics `json:"memory"`
	Disk        DiskMetrics `json:"disk"`
	Network     NetworkMetrics `json:"network"`
	Load        LoadMetrics `json:"load"`
	Processes   int       `json:"processes"`
	Uptime      int64     `json:"uptime"`
}

// CPUMetrics CPU监控指标
type CPUMetrics struct {
	Usage     float64 `json:"usage"`      // CPU使用�?(%)
	UserMode  float64 `json:"user_mode"`  // 用户模式使用�?(%)
	SystemMode float64 `json:"system_mode"` // 系统模式使用�?(%)
	Idle      float64 `json:"idle"`       // 空闲�?(%)
	IOWait    float64 `json:"iowait"`     // IO等待�?(%)
	Cores     int     `json:"cores"`      // CPU核心�?}

// MemoryMetrics 内存监控指标
type MemoryMetrics struct {
	Total     uint64  `json:"total"`      // 总内�?(bytes)
	Used      uint64  `json:"used"`       // 已使用内�?(bytes)
	Available uint64  `json:"available"`  // 可用内存 (bytes)
	Free      uint64  `json:"free"`       // 空闲内存 (bytes)
	Usage     float64 `json:"usage"`      // 内存使用�?(%)
	Buffers   uint64  `json:"buffers"`    // 缓冲�?(bytes)
	Cached    uint64  `json:"cached"`     // 缓存 (bytes)
	SwapTotal uint64  `json:"swap_total"` // 交换区总大�?(bytes)
	SwapUsed  uint64  `json:"swap_used"`  // 交换区使用量 (bytes)
	SwapFree  uint64  `json:"swap_free"`  // 交换区空闲量 (bytes)
}

// DiskMetrics 磁盘监控指标
type DiskMetrics struct {
	Partitions []PartitionMetrics `json:"partitions"`
	IOStats    DiskIOStats        `json:"io_stats"`
}

// PartitionMetrics 分区监控指标
type PartitionMetrics struct {
	Device     string  `json:"device"`     // 设备�?	Mountpoint string  `json:"mountpoint"` // 挂载�?	Filesystem string  `json:"filesystem"` // 文件系统类型
	Total      uint64  `json:"total"`      // 总空�?(bytes)
	Used       uint64  `json:"used"`       // 已使用空�?(bytes)
	Available  uint64  `json:"available"`  // 可用空间 (bytes)
	Usage      float64 `json:"usage"`      // 使用�?(%)
	Inodes     uint64  `json:"inodes"`     // inode总数
	InodesUsed uint64  `json:"inodes_used"` // 已使用inode�?	InodesFree uint64  `json:"inodes_free"` // 空闲inode�?}

// DiskIOStats 磁盘IO统计
type DiskIOStats struct {
	ReadBytes   uint64 `json:"read_bytes"`   // 读取字节�?	WriteBytes  uint64 `json:"write_bytes"`  // 写入字节�?	ReadOps     uint64 `json:"read_ops"`     // 读操作次�?	WriteOps    uint64 `json:"write_ops"`    // 写操作次�?	ReadTime    uint64 `json:"read_time"`    // 读取时间 (ms)
	WriteTime   uint64 `json:"write_time"`   // 写入时间 (ms)
	IOTime      uint64 `json:"io_time"`      // IO时间 (ms)
}

// NetworkMetrics 网络监控指标
type NetworkMetrics struct {
	Interfaces []NetworkInterface `json:"interfaces"`
}

// NetworkInterface 网络接口指标
type NetworkInterface struct {
	Name      string `json:"name"`       // 接口名称
	BytesRecv uint64 `json:"bytes_recv"` // 接收字节�?	BytesSent uint64 `json:"bytes_sent"` // 发送字节数
	PacketsRecv uint64 `json:"packets_recv"` // 接收数据包数
	PacketsSent uint64 `json:"packets_sent"` // 发送数据包�?	ErrorsRecv  uint64 `json:"errors_recv"`  // 接收错误�?	ErrorsSent  uint64 `json:"errors_sent"`  // 发送错误数
	DroppedRecv uint64 `json:"dropped_recv"` // 接收丢弃�?	DroppedSent uint64 `json:"dropped_sent"` // 发送丢弃数
}

// LoadMetrics 系统负载指标
type LoadMetrics struct {
	Load1  float64 `json:"load1"`  // 1分钟平均负载
	Load5  float64 `json:"load5"`  // 5分钟平均负载
	Load15 float64 `json:"load15"` // 15分钟平均负载
}

// ProcessMetrics 进程监控指标
type ProcessMetrics struct {
	PID        int     `json:"pid"`         // 进程ID
	Name       string  `json:"name"`        // 进程名称
	Status     string  `json:"status"`      // 进程状�?	CPUPercent float64 `json:"cpu_percent"` // CPU使用�?	MemoryRSS  uint64  `json:"memory_rss"`  // 物理内存使用�?	MemoryVMS  uint64  `json:"memory_vms"`  // 虚拟内存使用�?	MemoryPercent float64 `json:"memory_percent"` // 内存使用�?	OpenFiles  int     `json:"open_files"`  // 打开文件�?	Threads    int     `json:"threads"`     // 线程�?	CreateTime int64   `json:"create_time"` // 创建时间
}

// ServiceStatus 服务状�?type ServiceStatus struct {
	Name        string    `json:"name"`
	Status      string    `json:"status"`      // running, stopped, failed
	Enabled     bool      `json:"enabled"`     // 是否开机启�?	LastStarted time.Time `json:"last_started"`
	Uptime      int64     `json:"uptime"`
	PID         int       `json:"pid"`
}

// AlertRule 告警规则
type AlertRule struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	ServerID    uint    `json:"server_id"`
	MetricType  string  `json:"metric_type"`  // cpu, memory, disk, network
	Condition   string  `json:"condition"`    // >, <, >=, <=, ==
	Threshold   float64 `json:"threshold"`
	Duration    int     `json:"duration"`     // 持续时间（秒�?	Enabled     bool    `json:"enabled"`
	NotifyEmail bool    `json:"notify_email"`
	NotifyWebhook bool  `json:"notify_webhook"`
	WebhookURL  string  `json:"webhook_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Alert 告警信息
type Alert struct {
	ID          uint      `json:"id"`
	RuleID      uint      `json:"rule_id"`
	ServerID    uint      `json:"server_id"`
	ServerName  string    `json:"server_name"`
	MetricType  string    `json:"metric_type"`
	CurrentValue float64  `json:"current_value"`
	Threshold   float64   `json:"threshold"`
	Status      string    `json:"status"`      // firing, resolved
	Message     string    `json:"message"`
	FiredAt     time.Time `json:"fired_at"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
	NotifiedAt  *time.Time `json:"notified_at,omitempty"`
}
