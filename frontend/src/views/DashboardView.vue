<template>
  <div class="dashboard">
    <el-container class="dashboard-container">
      <el-header class="dashboard-header">
        <div class="header-left">
          <h1>自动化运维平台</h1>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <span class="user-dropdown">
              <el-avatar :size="32" src="https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png" />
              <span class="username">{{ userStore.user?.username }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人信息</el-dropdown-item>
                <el-dropdown-item command="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      
      <el-container>
        <el-aside class="dashboard-sidebar">
          <el-menu
            :default-active="currentRoute"
            router
            class="sidebar-menu"
          >
            <el-menu-item index="/dashboard">
              <el-icon><Monitor /></el-icon>
              <span>仪表盘</span>
            </el-menu-item>
            <el-menu-item index="/servers">
              <el-icon><Server /></el-icon>
              <span>服务器管理</span>
            </el-menu-item>
            <el-menu-item index="/deployments">
              <el-icon><Upload /></el-icon>
              <span>部署管理</span>
            </el-menu-item>
            <el-menu-item index="/tasks">
              <el-icon><Timer /></el-icon>
              <span>任务调度</span>
            </el-menu-item>
            <el-menu-item v-if="userStore.user?.role === 'admin'" index="/users">
              <el-icon><User /></el-icon>
              <span>用户管理</span>
            </el-menu-item>
          </el-menu>
        </el-aside>
        
        <el-main class="dashboard-main">
          <div class="stats-grid">
            <el-card class="stat-card" v-loading="loading">
              <div class="stat-content">
                <div class="stat-icon server-icon">
                  <el-icon size="24"><Server /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-number">{{ dashboardData?.stats?.total_servers || 0 }}</div>
                  <div class="stat-label">服务器总数</div>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card" v-loading="loading">
              <div class="stat-content">
                <div class="stat-icon deploy-icon">
                  <el-icon size="24"><Upload /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-number">{{ dashboardData?.stats?.online_servers || 0 }}</div>
                  <div class="stat-label">在线服务器</div>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card" v-loading="loading">
              <div class="stat-content">
                <div class="stat-icon task-icon">
                  <el-icon size="24"><Timer /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-number">{{ dashboardData?.stats?.offline_servers || 0 }}</div>
                  <div class="stat-label">离线服务器</div>
                </div>
              </div>
            </el-card>
            
            <el-card class="stat-card" v-loading="loading">
              <div class="stat-content">
                <div class="stat-icon alert-icon">
                  <el-icon size="24"><Warning /></el-icon>
                </div>
                <div class="stat-info">
                  <div class="stat-number">{{ dashboardData?.alerts?.length || 0 }}</div>
                  <div class="stat-label">告警数量</div>
                </div>
              </div>
            </el-card>
          </div>
          
          <div class="charts-grid">
            <el-card class="chart-card">
              <template #header>
                <span>系统监控</span>
              </template>
              <div class="chart-placeholder">
                <p>监控图表将在这里显示</p>
              </div>
            </el-card>
            
            <el-card class="chart-card">
              <template #header>
                <span>最近部署</span>
              </template>
              <div class="recent-deployments" v-loading="loading">
                <template v-if="dashboardData?.recent_activities?.length">
                  <el-timeline>
                    <el-timeline-item 
                      v-for="activity in dashboardData.recent_activities" 
                      :key="activity.id"
                      :timestamp="activity.time"
                    >
                      <p>{{ activity.description }}</p>
                      <el-tag 
                        size="small" 
                        :type="activity.status === 'success' ? 'success' : 'danger'"
                      >
                        {{ activity.status === 'success' ? '成功' : '失败' }}
                      </el-tag>
                    </el-timeline-item>
                  </el-timeline>
                </template>
                <template v-else>
                  <p style="text-align: center; color: #999;">暂无数据</p>
                </template>
              </div>
            </el-card>
          </div>
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ArrowDown, Monitor, Server, Upload, Timer, User, Warning } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { getDashboardData, type DashboardData } from '@/api/monitor'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const currentRoute = computed(() => route.path)
const dashboardData = ref<DashboardData | null>(null)
const loading = ref(true)

// 获取仪表盘数据
const fetchDashboardData = async () => {
  try {
    loading.value = true
    const response = await getDashboardData()
    dashboardData.value = response.data
  } catch (error) {
    console.error('获取仪表盘数据失败:', error)
    ElMessage.error('获取仪表盘数据失败')
  } finally {
    loading.value = false
  }
}

// 组件挂载时获取数据
onMounted(() => {
  fetchDashboardData()
  // 每30秒刷新一次
  setInterval(fetchDashboardData, 30000)
})

const handleCommand = async (command: string) => {
  switch (command) {
    case 'profile':
      ElMessage.info('个人信息页面开发中...')
      break
    case 'logout':
      try {
        await ElMessageBox.confirm('确定要退出登录吗？', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        })
        userStore.logout()
        router.push('/login')
        ElMessage.success('已退出登录')
      } catch {
        // 用户取消操作
      }
      break
  }
}
</script>

<style scoped>
.dashboard {
  height: 100vh;
}

.dashboard-container {
  height: 100%;
}

.dashboard-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  padding: 0 20px;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
}

.header-left h1 {
  margin: 0;
  color: #333;
  font-size: 20px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  cursor: pointer;
  padding: 8px;
  border-radius: 4px;
  transition: all 0.3s;
}

.user-dropdown:hover {
  background-color: #f5f5f5;
}

.username {
  margin: 0 8px;
  font-size: 14px;
}

.dashboard-sidebar {
  width: 200px;
  background: #001529;
}

.sidebar-menu {
  border: none;
  background: #001529;
}

.sidebar-menu .el-menu-item {
  color: #fff;
}

.sidebar-menu .el-menu-item:hover {
  background-color: #1890ff;
}

.sidebar-menu .el-menu-item.is-active {
  background-color: #1890ff;
}

.dashboard-main {
  padding: 20px;
  background-color: #f0f2f5;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-bottom: 20px;
}

.stat-card {
  border: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.stat-content {
  display: flex;
  align-items: center;
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16px;
  color: white;
}

.server-icon {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.deploy-icon {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.task-icon {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.alert-icon {
  background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
}

.stat-number {
  font-size: 24px;
  font-weight: bold;
  color: #333;
}

.stat-label {
  font-size: 14px;
  color: #666;
  margin-top: 4px;
}

.charts-grid {
  display: grid;
  grid-template-columns: 2fr 1fr;
  gap: 20px;
}

.chart-card {
  border: none;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

.chart-placeholder {
  height: 300px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 16px;
}

.recent-deployments {
  padding: 10px 0;
}
</style>