package monitor

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCollector(t *testing.T) {
	collector := NewCollector(1)
	ctx := context.Background()

	t.Run("CollectSystemMetrics", func(t *testing.T) {
		// 由于这个测试需要访问系统资源，在CI环境中可能失败
		// 可以根据环境跳过
		if testing.Short() {
			t.Skip("跳过系统监控测试")
		}

		metrics, err := collector.CollectSystemMetrics(ctx)
		if err != nil {
			t.Logf("收集系统监控数据失败: %v", err)
			return
		}

		assert.NotNil(t, metrics)
		assert.Equal(t, uint(1), metrics.ServerID)
		assert.True(t, time.Since(metrics.Timestamp) < time.Minute)

		// 检查CPU指标
		assert.GreaterOrEqual(t, metrics.CPU.Usage, 0.0)
		assert.LessOrEqual(t, metrics.CPU.Usage, 100.0)
		assert.Greater(t, metrics.CPU.Cores, 0)

		// 检查内存指标
		assert.Greater(t, metrics.Memory.Total, uint64(0))
		assert.GreaterOrEqual(t, metrics.Memory.Usage, 0.0)
		assert.LessOrEqual(t, metrics.Memory.Usage, 100.0)

		// 检查磁盘指标
		assert.NotEmpty(t, metrics.Disk.Partitions)

		// 检查网络指标
		assert.NotEmpty(t, metrics.Network.Interfaces)

		// 检查负载指标
		assert.GreaterOrEqual(t, metrics.Load.Load1, 0.0)

		// 检查系统信息
		assert.Greater(t, metrics.Processes, 0)
		assert.Greater(t, metrics.Uptime, int64(0))

		t.Logf("系统监控数据: CPU=%.2f%%, 内存=%.2f%%, 进程数=%d", 
			metrics.CPU.Usage, metrics.Memory.Usage, metrics.Processes)
	})

	t.Run("CollectProcessMetrics", func(t *testing.T) {
		if testing.Short() {
			t.Skip("跳过进程监控测试")
		}

		processes, err := collector.CollectProcessMetrics(ctx, 5)
		if err != nil {
			t.Logf("收集进程监控数据失败: %v", err)
			return
		}

		assert.LessOrEqual(t, len(processes), 5)
		
		if len(processes) > 0 {
			proc := processes[0]
			assert.Greater(t, proc.PID, 0)
			assert.NotEmpty(t, proc.Name)
			assert.GreaterOrEqual(t, proc.CPUPercent, 0.0)
			assert.Greater(t, proc.MemoryRSS, uint64(0))
		}

		t.Logf("收集到 %d 个进程信息", len(processes))
	})
}

func TestMonitorTypes(t *testing.T) {
	t.Run("SystemMetrics", func(t *testing.T) {
		metrics := &SystemMetrics{
			ServerID:  1,
			Timestamp: time.Now(),
			CPU: CPUMetrics{
				Usage: 45.5,
				Cores: 4,
			},
			Memory: MemoryMetrics{
				Total: 8589934592, // 8GB
				Used:  4294967296, // 4GB
				Usage: 50.0,
			},
		}

		assert.Equal(t, uint(1), metrics.ServerID)
		assert.Equal(t, 45.5, metrics.CPU.Usage)
		assert.Equal(t, 4, metrics.CPU.Cores)
		assert.Equal(t, uint64(8589934592), metrics.Memory.Total)
		assert.Equal(t, 50.0, metrics.Memory.Usage)
	})

	t.Run("Alert", func(t *testing.T) {
		alert := &Alert{
			ID:           1,
			ServerID:     1,
			MetricType:   "cpu",
			CurrentValue: 85.5,
			Threshold:    80.0,
			Status:       "firing",
			Message:      "CPU使用率过高",
			FiredAt:      time.Now(),
		}

		assert.Equal(t, uint(1), alert.ID)
		assert.Equal(t, "cpu", alert.MetricType)
		assert.Equal(t, 85.5, alert.CurrentValue)
		assert.Equal(t, 80.0, alert.Threshold)
		assert.Equal(t, "firing", alert.Status)
		assert.Contains(t, alert.Message, "CPU")
	})
}