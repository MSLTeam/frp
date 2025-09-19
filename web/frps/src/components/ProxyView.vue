<template>
  <div class="container">
    <el-card style="width: 100%; border-radius: 12px; box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);">
      <template #header>
        <div style="display: flex; justify-content: space-between; align-items: center;">
          <h3 style="margin: 0; font-size: 18px; font-weight: 600;">{{ proxyType }}</h3>
          <div class="flex items-center">
            <el-popconfirm
              title="确定清除全部离线隧道?"
              @confirm="clearOfflineProxies"
            >
              <template #reference>
                <el-button type="danger" plain>清理离线隧道</el-button>
              </template>
            </el-popconfirm>
            <el-button type="primary" @click="$emit('refresh')">刷新</el-button>
          </div>
        </div>
      </template>
      <div>
        <el-table
          :data="proxies"
          :default-sort="{ prop: 'name', order: 'ascending' }"
          style="width: 100%; border-radius: 8px;"
          stripe
        >
          <el-table-column type="expand">
            <template #default="props">
              <ProxyViewExpand :row="props.row" :proxyType="proxyType" />
            </template>
          </el-table-column>
          <el-table-column label="隧道名字" sortable>
            <template #default="scope">
              {{ scope.row.name.split('.')[1] || '' }}
            </template>
          </el-table-column>
          <el-table-column label="用户UID" sortable>
            <template #default="scope">
              {{ Number(scope.row.name.split('.')[0].split('-')[1]) - 10000 }}
            </template>
          </el-table-column>
          <el-table-column label="远程端口" prop="port" sortable> </el-table-column>
          <el-table-column label="连接数量" prop="conns" sortable>
          </el-table-column>
          <el-table-column
            label="流量(入)"
            prop="trafficIn"
            :formatter="formatTrafficIn"
            sortable
          >
          </el-table-column>
          <el-table-column
            label="流量(出)"
            prop="trafficOut"
            :formatter="formatTrafficOut"
            sortable
          >
          </el-table-column>
          <el-table-column label="客户端版本" prop="clientVersion" sortable>
          </el-table-column>
          <el-table-column label="状态" prop="status" sortable>
            <template #default="scope">
              <el-tag v-if="scope.row.status === 'online'" type="success">{{
                scope.row.status === 'online' ? '在线' : scope.row.status
              }}</el-tag>
              <el-tag v-else type="danger">{{ scope.row.status === 'offline' ? '离线' : scope.row.status }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作">
            <template #default="scope">
              <el-button
                type="primary"
                :name="scope.row.name"
                style="margin-bottom: 10px; border-radius: 6px;"
                @click="dialogVisibleName = scope.row.name; dialogVisible = true"
              >流量统计
              </el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <el-dialog
        v-model="dialogVisible"
        destroy-on-close="true"
        :title="dialogVisibleName"
        width="700px"
        :modal="true"
        :close-on-click-modal="false"
        :close-on-press-escape="false"
      >
        <Traffic :proxyName="dialogVisibleName" />
      </el-dialog>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import * as Humanize from 'humanize-plus'
import type { TableColumnCtx } from 'element-plus'
import type { BaseProxy } from '../utils/proxy.js'
import { ElMessage } from 'element-plus'
import ProxyViewExpand from './ProxyViewExpand.vue'
import { ref } from 'vue'

defineProps<{
  proxies: BaseProxy[]
  proxyType: string
}>()

const emit = defineEmits(['refresh'])

const dialogVisible = ref(false)
const dialogVisibleName = ref("")

const formatTrafficIn = (row: BaseProxy, _: TableColumnCtx<BaseProxy>) => {
  return Humanize.fileSize(row.trafficIn)
}

const formatTrafficOut = (row: BaseProxy, _: TableColumnCtx<BaseProxy>) => {
  return Humanize.fileSize(row.trafficOut)
}

const clearOfflineProxies = () => {
  fetch('../api/proxies?status=offline', {
    method: 'DELETE',
    credentials: 'include',
  })
    .then((res) => {
      if (res.ok) {
        ElMessage({
          message: '成功清理离线隧道！',
          type: 'success',
        })
        emit('refresh')
      } else {
        ElMessage({
          message: '无法清理离线隧道: ' + res.status + ' ' + res.statusText,
          type: 'warning',
        })
      }
    })
    .catch((err) => {
      ElMessage({
        message: '无法清理离线隧道: ' + err.message,
        type: 'warning',
      })
    })
}
</script>

<style scoped>
.container {
  padding: 20px;
}

.el-button {
  border-radius: 6px;
}

.el-table {
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.el-dialog {
  border-radius: 12px;
}

</style>
