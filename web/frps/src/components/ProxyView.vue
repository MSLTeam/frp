<template>
  <div class="proxy-view-page">
    <el-card class="proxy-card">
      <template #header>
        <div class="proxy-head">
          <div class="head-title-wrap">
            <h3>{{ proxyType.toUpperCase() }} 隧道</h3>
            <p>隧道运行状态与流量分析</p>
          </div>

          <div class="head-actions">
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

        <div class="proxy-summary">
          <div class="summary-chip">
            <span class="chip-label">总数</span>
            <strong>{{ proxies.length }}</strong>
          </div>
          <div class="summary-chip success">
            <span class="chip-label">在线</span>
            <strong>{{ onlineCount }}</strong>
          </div>
          <div class="summary-chip danger">
            <span class="chip-label">离线</span>
            <strong>{{ offlineCount }}</strong>
          </div>
        </div>
      </template>

      <el-table
        :data="proxies"
        :default-sort="{ prop: 'name', order: 'ascending' }"
        stripe
        class="proxy-table"
      >
        <el-table-column type="expand">
          <template #default="props">
            <ProxyViewExpand :row="props.row" :proxyType="proxyType" />
          </template>
        </el-table-column>

        <el-table-column label="隧道名字" sortable>
          <template #default="scope">
            {{ tunnelName(scope.row.name) }}
          </template>
        </el-table-column>

        <el-table-column label="用户 UID" sortable>
          <template #default="scope">
            {{ uid(scope.row.name) }}
          </template>
        </el-table-column>

        <el-table-column label="远程端口" prop="port" sortable />
        <el-table-column label="连接数" prop="conns" sortable />

        <el-table-column
          label="流量(入)"
          prop="trafficIn"
          :formatter="formatTrafficIn"
          sortable
        />

        <el-table-column
          label="流量(出)"
          prop="trafficOut"
          :formatter="formatTrafficOut"
          sortable
        />

        <el-table-column label="客户端版本" prop="clientVersion" sortable />

        <el-table-column label="状态" prop="status" sortable>
          <template #default="scope">
            <el-tag v-if="scope.row.status === 'online'" type="success">在线</el-tag>
            <el-tag v-else type="danger">离线</el-tag>
          </template>
        </el-table-column>

        <el-table-column label="操作" width="130">
          <template #default="scope">
            <el-button
              type="primary"
              plain
              @click="dialogVisibleName = scope.row.name; dialogVisible = true"
            >
              流量统计
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-dialog
        v-model="dialogVisible"
        destroy-on-close
        :title="dialogVisibleName"
        width="760px"
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
import { computed, ref } from 'vue'
import * as Humanize from 'humanize-plus'
import type { TableColumnCtx } from 'element-plus'
import type { BaseProxy } from '../utils/proxy.js'
import { ElMessage } from 'element-plus'
import ProxyViewExpand from './ProxyViewExpand.vue'
import Traffic from './Traffic.vue'

const props = defineProps<{
  proxies: BaseProxy[]
  proxyType: string
}>()

const emit = defineEmits(['refresh'])

const dialogVisible = ref(false)
const dialogVisibleName = ref('')

const onlineCount = computed(() => {
  return props.proxies.filter((item) => item.status === 'online').length
})

const offlineCount = computed(() => {
  return props.proxies.length - onlineCount.value
})

const tunnelName = (name: string) => {
  return name.split('.')[1] || ''
}

const uid = (name: string) => {
  const raw = Number(name.split('.')[0]?.split('-')[1])
  if (Number.isNaN(raw)) {
    return '-'
  }
  return raw - 10000
}

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
.proxy-view-page {
  padding: 4px;
}

.proxy-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.head-title-wrap {
  h3 {
    margin: 0;
    font-size: 18px;
    line-height: 1.2;
  }

  p {
    margin: 3px 0 0;
    font-size: 12px;
    color: var(--frp-muted);
  }
}

.head-actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.proxy-summary {
  margin-top: 10px;
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.summary-chip {
  border-radius: 999px;
  border: 1px solid var(--frp-border);
  background: var(--frp-surface-bg);
  padding: 6px 10px;
  font-size: 12px;
  display: inline-flex;
  align-items: center;
  gap: 8px;

  strong {
    font-size: 14px;
    line-height: 1;
  }
}

.summary-chip.success {
  border-color: rgba(34, 197, 94, 0.35);
  background: rgba(34, 197, 94, 0.12);
}

.summary-chip.danger {
  border-color: rgba(239, 68, 68, 0.32);
  background: rgba(239, 68, 68, 0.1);
}

.chip-label {
  color: var(--frp-muted);
  font-weight: 600;
}

.proxy-table {
  margin-top: 4px;
}

@media (max-width: 900px) {
  .proxy-head {
    flex-direction: column;
    align-items: stretch;
  }

  .head-actions {
    justify-content: flex-start;
  }
}
</style>
