<template>
  <div class="server-view">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>服务器管理</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加服务器
          </el-button>
        </div>
      </template>
      
      <el-table :data="servers" style="width: 100%">
        <el-table-column prop="name" label="服务器名称" />
        <el-table-column prop="host" label="主机地址" />
        <el-table-column prop="environment" label="环境" />
        <el-table-column prop="status" label="状态">
          <template #default="scope">
            <el-tag :type="scope.row.status === 1 ? 'success' : 'danger'">
              {{ scope.row.status === 1 ? '在线' : '离线' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作">
          <template #default="scope">
            <el-button size="small" @click="handleEdit(scope.row)">编辑</el-button>
            <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'

const servers = ref([
  { id: 1, name: '生产服务器1', host: '192.168.1.100', environment: 'prod', status: 1 },
  { id: 2, name: '测试服务器1', host: '192.168.1.101', environment: 'test', status: 1 },
  { id: 3, name: '开发服务器1', host: '192.168.1.102', environment: 'dev', status: 0 }
])

const handleAdd = () => {
  ElMessage.info('添加服务器功能开发中...')
}

const handleEdit = (row: any) => {
  ElMessage.info(`编辑服务器 ${row.name} 功能开发中...`)
}

const handleDelete = (row: any) => {
  ElMessage.info(`删除服务器 ${row.name} 功能开发中...`)
}
</script>

<style scoped>
.server-view {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>