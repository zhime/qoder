import request from './request'

export interface ServerMetrics {
  server_id: number
  timestamp: string
  cpu: {
    usage: number
    cores: number
  }
  memory: {
    total: number
    used: number
    usage: number
  }
  disk: {
    partitions: Array<{
      device: string
      mountpoint: string
      total: number
      used: number
      usage: number
    }>
  }
  network: {
    interfaces: Array<{
      name: string
      bytes_recv: number
      bytes_sent: number
    }>
  }
  load: {
    load1: number
    load5: number
    load15: number
  }
  processes: number
  uptime: number
}

export interface SystemStats {
  total_servers: number
  online_servers: number
  offline_servers: number
  timestamp: string
}

export interface DashboardData {
  stats: SystemStats
  alerts: Array<{
    id: number
    level: string
    message: string
    time: string
  }>
  recent_activities: Array<{
    id: number
    type: string
    description: string
    status: string
    time: string
  }>
}

// 获取仪表盘数据
export const getDashboardData = () => {
  return request.get<DashboardData>('/monitor/dashboard')
}

// 获取系统统计
export const getSystemStats = () => {
  return request.get<SystemStats>('/monitor/stats')
}

// 获取服务器监控数据
export const getServerMetrics = (serverId: number) => {
  return request.get<ServerMetrics>(`/monitor/servers/${serverId}/metrics`)
}

// 获取服务器状态
export const getServerStatus = (serverId: number) => {
  return request.get<{ server_id: number; status: string }>(`/monitor/servers/${serverId}/status`)
}

// 获取服务器历史数据
export const getServerHistory = (serverId: number, timeRange = '1h', metric = 'cpu') => {
  return request.get(`/monitor/servers/${serverId}/history`, {
    params: { time_range: timeRange, metric }
  })
}

// 添加服务器到监控
export const addServerToMonitor = (serverId: number) => {
  return request.post('/monitor/servers', { server_id: serverId })
}

// 从监控中移除服务器
export const removeServerFromMonitor = (serverId: number) => {
  return request.delete(`/monitor/servers/${serverId}`)
}