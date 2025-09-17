package monitor

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
)

// Collector 监控数据收集器
type Collector struct {
	serverID uint
}

// NewCollector 创建监控数据收集器
func NewCollector(serverID uint) *Collector {
	return &Collector{
		serverID: serverID,
	}
}

// CollectSystemMetrics 收集系统监控指标
func (c *Collector) CollectSystemMetrics(ctx context.Context) (*SystemMetrics, error) {
	metrics := &SystemMetrics{
		ServerID:  c.serverID,
		Timestamp: time.Now(),
	}

	// 并发收集各项指标
	errChan := make(chan error, 6)

	go func() {
		var err error
		metrics.CPU, err = c.collectCPUMetrics(ctx)
		errChan <- err
	}()

	go func() {
		var err error
		metrics.Memory, err = c.collectMemoryMetrics(ctx)
		errChan <- err
	}()

	go func() {
		var err error
		metrics.Disk, err = c.collectDiskMetrics(ctx)
		errChan <- err
	}()

	go func() {
		var err error
		metrics.Network, err = c.collectNetworkMetrics(ctx)
		errChan <- err
	}()

	go func() {
		var err error
		metrics.Load, err = c.collectLoadMetrics(ctx)
		errChan <- err
	}()

	go func() {
		var err error
		metrics.Processes, metrics.Uptime, err = c.collectSystemInfo(ctx)
		errChan <- err
	}()

	// 等待所有收集完成
	for i := 0; i < 6; i++ {
		if err := <-errChan; err != nil {
			return nil, fmt.Errorf("收集监控数据失败: %w", err)
		}
	}

	return metrics, nil
}

// collectCPUMetrics 收集CPU指标
func (c *Collector) collectCPUMetrics(ctx context.Context) (CPUMetrics, error) {
	// 获取CPU使用率
	cpuPercent, err := cpu.PercentWithContext(ctx, time.Second, false)
	if err != nil {
		return CPUMetrics{}, err
	}

	// 获取CPU详细信息
	cpuTimes, err := cpu.TimesWithContext(ctx, false)
	if err != nil {
		return CPUMetrics{}, err
	}

	// 获取CPU核心数
	cores, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		cores = runtime.NumCPU()
	}

	var usage float64
	if len(cpuPercent) > 0 {
		usage = cpuPercent[0]
	}

	var cpuTime cpu.TimesStat
	if len(cpuTimes) > 0 {
		cpuTime = cpuTimes[0]
	}

	total := cpuTime.User + cpuTime.System + cpuTime.Idle + cpuTime.Iowait + cpuTime.Nice + cpuTime.Irq + cpuTime.Softirq

	return CPUMetrics{
		Usage:      usage,
		UserMode:   (cpuTime.User / total) * 100,
		SystemMode: (cpuTime.System / total) * 100,
		Idle:       (cpuTime.Idle / total) * 100,
		IOWait:     (cpuTime.Iowait / total) * 100,
		Cores:      cores,
	}, nil
}

// collectMemoryMetrics 收集内存指标
func (c *Collector) collectMemoryMetrics(ctx context.Context) (MemoryMetrics, error) {
	vmStat, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return MemoryMetrics{}, err
	}

	swapStat, err := mem.SwapMemoryWithContext(ctx)
	if err != nil {
		return MemoryMetrics{}, err
	}

	return MemoryMetrics{
		Total:     vmStat.Total,
		Used:      vmStat.Used,
		Available: vmStat.Available,
		Free:      vmStat.Free,
		Usage:     vmStat.UsedPercent,
		Buffers:   vmStat.Buffers,
		Cached:    vmStat.Cached,
		SwapTotal: swapStat.Total,
		SwapUsed:  swapStat.Used,
		SwapFree:  swapStat.Free,
	}, nil
}

// collectDiskMetrics 收集磁盘指标
func (c *Collector) collectDiskMetrics(ctx context.Context) (DiskMetrics, error) {
	// 获取磁盘分区信息
	partitions, err := disk.PartitionsWithContext(ctx, false)
	if err != nil {
		return DiskMetrics{}, err
	}

	var partitionMetrics []PartitionMetrics
	for _, partition := range partitions {
		usage, err := disk.UsageWithContext(ctx, partition.Mountpoint)
		if err != nil {
			continue
		}

		partitionMetrics = append(partitionMetrics, PartitionMetrics{
			Device:     partition.Device,
			Mountpoint: partition.Mountpoint,
			Filesystem: partition.Fstype,
			Total:      usage.Total,
			Used:       usage.Used,
			Available:  usage.Free,
			Usage:      usage.UsedPercent,
			Inodes:     usage.InodesTotal,
			InodesUsed: usage.InodesUsed,
			InodesFree: usage.InodesFree,
		})
	}

	// 获取磁盘IO统计
	ioStats, err := disk.IOCountersWithContext(ctx)
	if err != nil {
		return DiskMetrics{
			Partitions: partitionMetrics,
			IOStats:    DiskIOStats{},
		}, nil
	}

	var totalIOStats DiskIOStats
	for _, stat := range ioStats {
		totalIOStats.ReadBytes += stat.ReadBytes
		totalIOStats.WriteBytes += stat.WriteBytes
		totalIOStats.ReadOps += stat.ReadCount
		totalIOStats.WriteOps += stat.WriteCount
		totalIOStats.ReadTime += stat.ReadTime
		totalIOStats.WriteTime += stat.WriteTime
		totalIOStats.IOTime += stat.IoTime
	}

	return DiskMetrics{
		Partitions: partitionMetrics,
		IOStats:    totalIOStats,
	}, nil
}

// collectNetworkMetrics 收集网络指标
func (c *Collector) collectNetworkMetrics(ctx context.Context) (NetworkMetrics, error) {
	interfaces, err := net.IOCountersWithContext(ctx, true)
	if err != nil {
		return NetworkMetrics{}, err
	}

	var networkInterfaces []NetworkInterface
	for _, iface := range interfaces {
		// 跳过回环接口
		if iface.Name == "lo" || strings.HasPrefix(iface.Name, "lo") {
			continue
		}

		networkInterfaces = append(networkInterfaces, NetworkInterface{
			Name:        iface.Name,
			BytesRecv:   iface.BytesRecv,
			BytesSent:   iface.BytesSent,
			PacketsRecv: iface.PacketsRecv,
			PacketsSent: iface.PacketsSent,
			ErrorsRecv:  iface.Errin,
			ErrorsSent:  iface.Errout,
			DroppedRecv: iface.Dropin,
			DroppedSent: iface.Dropout,
		})
	}

	return NetworkMetrics{
		Interfaces: networkInterfaces,
	}, nil
}

// collectLoadMetrics 收集负载指标
func (c *Collector) collectLoadMetrics(ctx context.Context) (LoadMetrics, error) {
	loadStat, err := load.AvgWithContext(ctx)
	if err != nil {
		return LoadMetrics{}, err
	}

	return LoadMetrics{
		Load1:  loadStat.Load1,
		Load5:  loadStat.Load5,
		Load15: loadStat.Load15,
	}, nil
}

// collectSystemInfo 收集系统信息
func (c *Collector) collectSystemInfo(ctx context.Context) (int, int64, error) {
	// 获取运行时间
	hostInfo, err := host.InfoWithContext(ctx)
	if err != nil {
		return 0, 0, err
	}

	// 获取进程数量
	processes, err := process.PidsWithContext(ctx)
	if err != nil {
		return 0, int64(hostInfo.Uptime), err
	}

	return len(processes), int64(hostInfo.Uptime), nil
}

// CollectProcessMetrics 收集进程指标
func (c *Collector) CollectProcessMetrics(ctx context.Context, limit int) ([]ProcessMetrics, error) {
	pids, err := process.PidsWithContext(ctx)
	if err != nil {
		return nil, err
	}

	var processMetrics []ProcessMetrics
	count := 0

	for _, pid := range pids {
		if count >= limit {
			break
		}

		proc, err := process.NewProcessWithContext(ctx, pid)
		if err != nil {
			continue
		}

		name, err := proc.NameWithContext(ctx)
		if err != nil {
			continue
		}

		status, err := proc.StatusWithContext(ctx)
		if err != nil {
			status = []string{"unknown"}
		}

		var statusStr string
		if len(status) > 0 {
			statusStr = status[0]
		} else {
			statusStr = "unknown"
		}

		cpuPercent, err := proc.CPUPercentWithContext(ctx)
		if err != nil {
			cpuPercent = 0
		}

		memInfo, err := proc.MemoryInfoWithContext(ctx)
		if err != nil {
			continue
		}

		memPercent, err := proc.MemoryPercentWithContext(ctx)
		if err != nil {
			memPercent = 0
		}

		openFiles, err := proc.NumFDsWithContext(ctx)
		if err != nil {
			openFiles = 0
		}

		threads, err := proc.NumThreadsWithContext(ctx)
		if err != nil {
			threads = 0
		}

		createTime, err := proc.CreateTimeWithContext(ctx)
		if err != nil {
			createTime = 0
		}

		processMetrics = append(processMetrics, ProcessMetrics{
			PID:           int(pid),
			Name:          name,
			Status:        statusStr,
			CPUPercent:    cpuPercent,
			MemoryRSS:     memInfo.RSS,
			MemoryVMS:     memInfo.VMS,
			MemoryPercent: float64(memPercent),
			OpenFiles:     int(openFiles),
			Threads:       int(threads),
			CreateTime:    createTime,
		})

		count++
	}

	return processMetrics, nil
}

// CheckServiceStatus 检查服务状态
func (c *Collector) CheckServiceStatus(ctx context.Context, serviceName string) (*ServiceStatus, error) {
	// 这里实现服务状态检查逻辑
	// 不同操作系统可能需要不同的实现

	status := &ServiceStatus{
		Name:   serviceName,
		Status: "unknown",
	}

	// Linux systemd
	if runtime.GOOS == "linux" {
		return c.checkSystemdService(ctx, serviceName)
	}

	// Windows服务
	if runtime.GOOS == "windows" {
		return c.checkWindowsService(ctx, serviceName)
	}

	return status, nil
}

// checkSystemdService 检查systemd服务状态
func (c *Collector) checkSystemdService(ctx context.Context, serviceName string) (*ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "systemctl", "show", serviceName, "--property=ActiveState,LoadState,SubState,MainPID,ExecMainStartTimestamp")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	status := &ServiceStatus{
		Name: serviceName,
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ActiveState=") {
			state := strings.TrimPrefix(line, "ActiveState=")
			if state == "active" {
				status.Status = "running"
			} else if state == "inactive" {
				status.Status = "stopped"
			} else if state == "failed" {
				status.Status = "failed"
			}
		} else if strings.HasPrefix(line, "MainPID=") {
			pidStr := strings.TrimPrefix(line, "MainPID=")
			if pid, err := strconv.Atoi(pidStr); err == nil {
				status.PID = pid
			}
		} else if strings.HasPrefix(line, "ExecMainStartTimestamp=") {
			timeStr := strings.TrimPrefix(line, "ExecMainStartTimestamp=")
			if t, err := time.Parse("Mon 2006-01-02 15:04:05 MST", timeStr); err == nil {
				status.LastStarted = t
				status.Uptime = int64(time.Since(t).Seconds())
			}
		}
	}

	return status, nil
}

// checkWindowsService 检查Windows服务状态
func (c *Collector) checkWindowsService(ctx context.Context, serviceName string) (*ServiceStatus, error) {
	cmd := exec.CommandContext(ctx, "sc", "query", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	status := &ServiceStatus{
		Name:   serviceName,
		Status: "stopped",
	}

	if strings.Contains(string(output), "RUNNING") {
		status.Status = "running"
	} else if strings.Contains(string(output), "STOPPED") {
		status.Status = "stopped"
	}

	return status, nil
}
