<template>
  <div class="server-overview">
    <el-row :gutter="20" class="stats-row">
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          label="客户端"
          :value="data.clientCounts"
          type="clients"
          subtitle="已连接的客户端数量"
          to="/clients"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          label="隧道"
          :value="data.proxyCounts"
          type="proxies"
          subtitle="已连接的隧道数"
          to="/proxies/tcp"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          label="连接"
          :value="data.curConns"
          type="connections"
          subtitle="隧道被连接总数"
        />
      </el-col>
      <el-col :xs="24" :sm="12" :lg="6">
        <StatCard
          label="流量"
          :value="formatTrafficTotal()"
          type="traffic"
          subtitle="今日流量"
        />
      </el-col>
    </el-row>

    <el-row :gutter="20" class="charts-row">
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">节点流量</span>
              <el-tag size="small" type="info">今日</el-tag>
            </div>
          </template>
          <div class="traffic-summary">
            <div class="traffic-item in">
              <div class="traffic-icon">
                <el-icon><Download /></el-icon>
              </div>
              <div class="traffic-info">
                <div class="label">上传</div>
                <div class="value">
                  {{ formatFileSize(data.totalTrafficIn) }}
                </div>
              </div>
            </div>
            <div class="traffic-divider"></div>
            <div class="traffic-item out">
              <div class="traffic-icon">
                <el-icon><Upload /></el-icon>
              </div>
              <div class="traffic-info">
                <div class="label">下载</div>
                <div class="value">
                  {{ formatFileSize(data.totalTrafficOut) }}
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card class="chart-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">隧道类型</span>
              <el-tag size="small" type="info">实时</el-tag>
            </div>
          </template>
          <div class="proxy-types-grid">
            <div
              v-for="(count, type) in data.proxyTypeCounts"
              :key="type"
              class="proxy-type-item"
              v-show="count > 0"
            >
              <div class="proxy-type-name">{{ type.toUpperCase() }}</div>
              <div class="proxy-type-count">{{ count }}</div>
            </div>
            <div v-if="!hasActiveProxies" class="no-data">没有活跃的隧道</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="config-card" shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">节点配置</span>
          <el-tag size="small" type="success">v{{ data.version }}</el-tag>
        </div>
      </template>
      <div class="config-grid">
        <div class="config-item">
          <span class="config-label">服务端口</span>
          <span class="config-value">{{ data.bindPort }}</span>
        </div>
        <div class="config-item" v-if="data.kcpBindPort != 0">
          <span class="config-label">KCP 服务端口</span>
          <span class="config-value">{{ data.kcpBindPort }}</span>
        </div>
        <div class="config-item" v-if="data.quicBindPort != 0">
          <span class="config-label">QUIC 服务端口</span>
          <span class="config-value">{{ data.quicBindPort }}</span>
        </div>
        <div class="config-item" v-if="data.vhostHTTPPort != 0">
          <span class="config-label">HTTP 服务端口</span>
          <span class="config-value">{{ data.vhostHTTPPort }}</span>
        </div>
        <div class="config-item" v-if="data.vhostHTTPSPort != 0">
          <span class="config-label">HTTPS 服务端口</span>
          <span class="config-value">{{ data.vhostHTTPSPort }}</span>
        </div>
        <div class="config-item" v-if="data.tcpmuxHTTPConnectPort != 0">
          <span class="config-label">TCPMux 服务端口</span>
          <span class="config-value">{{ data.tcpmuxHTTPConnectPort }}</span>
        </div>
        <div class="config-item" v-if="data.subdomainHost != ''">
          <span class="config-label">子域名</span>
          <span class="config-value">{{ data.subdomainHost }}</span>
        </div>
        <div class="config-item">
          <span class="config-label">最大连接池</span>
          <span class="config-value">{{ data.maxPoolCount }}</span>
        </div>
        <div class="config-item">
          <span class="config-label">最大端口数/客户端</span>
          <span class="config-value">{{ data.maxPortsPerClient }}</span>
        </div>
        <div class="config-item" v-if="data.allowPortsStr != ''">
          <span class="config-label">可用连接端口</span>
          <span class="config-value">{{ data.allowPortsStr }}</span>
        </div>
        <div class="config-item" v-if="data.tlsForce">
          <span class="config-label">强制 TLS</span>
          <el-tag size="small" type="warning">已启用</el-tag>
        </div>
        <div class="config-item">
          <span class="config-label">心跳包超时</span>
          <span class="config-value">{{ data.heartbeatTimeout }}s</span>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { formatFileSize } from '../utils/format'
import { Download, Upload } from '@element-plus/icons-vue'
import StatCard from '../components/StatCard.vue'
import { getServerInfo } from '../api/server'

const data = ref({
  version: '',
  bindPort: 0,
  kcpBindPort: 0,
  quicBindPort: 0,
  vhostHTTPPort: 0,
  vhostHTTPSPort: 0,
  tcpmuxHTTPConnectPort: 0,
  subdomainHost: '',
  maxPoolCount: 0,
  maxPortsPerClient: '',
  allowPortsStr: '',
  tlsForce: false,
  heartbeatTimeout: 0,
  clientCounts: 0,
  curConns: 0,
  proxyCounts: 0,
  totalTrafficIn: 0,
  totalTrafficOut: 0,
  proxyTypeCounts: {} as Record<string, number>,
})

const hasActiveProxies = computed(() => {
  return Object.values(data.value.proxyTypeCounts).some((c) => c > 0)
})

const formatTrafficTotal = () => {
  const total = data.value.totalTrafficIn + data.value.totalTrafficOut
  return formatFileSize(total)
}

const fetchData = async () => {
  try {
    const json = await getServerInfo()
    data.value.version = json.version
    data.value.bindPort = json.bindPort
    data.value.kcpBindPort = json.kcpBindPort
    data.value.quicBindPort = json.quicBindPort
    data.value.vhostHTTPPort = json.vhostHTTPPort
    data.value.vhostHTTPSPort = json.vhostHTTPSPort
    data.value.tcpmuxHTTPConnectPort = json.tcpmuxHTTPConnectPort
    data.value.subdomainHost = json.subdomainHost
    data.value.maxPoolCount = json.maxPoolCount
    data.value.maxPortsPerClient = String(json.maxPortsPerClient)
    if (data.value.maxPortsPerClient == '0') {
      data.value.maxPortsPerClient = 'no limit'
    }
    data.value.allowPortsStr = json.allowPortsStr
    data.value.tlsForce = json.tlsForce
    data.value.heartbeatTimeout = json.heartbeatTimeout
    data.value.clientCounts = json.clientCounts
    data.value.curConns = json.curConns
    data.value.totalTrafficIn = json.totalTrafficIn
    data.value.totalTrafficOut = json.totalTrafficOut
    data.value.proxyTypeCounts = json.proxyTypeCount || {}

    data.value.proxyCounts = 0
    if (json.proxyTypeCount != null) {
      Object.values(json.proxyTypeCount).forEach((count: any) => {
        data.value.proxyCounts += count || 0
      })
    }
  } catch (err) {
    ElMessage({
      showClose: true,
      message: 'Get server info from frps failed!',
      type: 'error',
    })
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.server-overview {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.charts-row {
  margin-bottom: 20px;
}

.chart-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  height: 100%;
}

html.dark .chart-card {
  border-color: #3a3d5c;
  background: #27293d;
}

.config-card {
  border-radius: 12px;
  border: 1px solid #e4e7ed;
  margin-bottom: 20px;
}

html.dark .config-card {
  border-color: #3a3d5c;
  background: #27293d;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-title {
  font-size: 16px;
  font-weight: 500;
  color: #303133;
}

html.dark .card-title {
  color: #e5e7eb;
}

.traffic-summary {
  display: flex;
  align-items: center;
  justify-content: space-around;
  min-height: 120px;
  padding: 10px 0;
}

.traffic-item {
  display: flex;
  align-items: center;
  gap: 16px;
}

.traffic-icon {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
}

.traffic-item.in .traffic-icon {
  background: rgba(84, 112, 198, 0.1);
  color: #5470c6;
}

.traffic-item.out .traffic-icon {
  background: rgba(145, 204, 117, 0.1);
  color: #91cc75;
}

.traffic-info {
  display: flex;
  flex-direction: column;
}

.traffic-info .label {
  font-size: 14px;
  color: #909399;
}

.traffic-info .value {
  font-size: 24px;
  font-weight: 500;
  color: #303133;
}

html.dark .traffic-info .value {
  color: #e5e7eb;
}

.traffic-divider {
  width: 1px;
  height: 60px;
  background: #e4e7ed;
}

html.dark .traffic-divider {
  background: #3a3d5c;
}

.proxy-types-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(100px, 1fr));
  gap: 16px;
  min-height: 120px;
  align-content: center;
  padding: 10px 0;
}

.proxy-type-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
}

html.dark .proxy-type-item {
  background: #1e1e2d;
}

.proxy-type-name {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
  margin-bottom: 4px;
}

.proxy-type-count {
  font-size: 20px;
  font-weight: 500;
  color: #303133;
}

html.dark .proxy-type-count {
  color: #e5e7eb;
}

.no-data {
  grid-column: 1 / -1;
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #909399;
  font-size: 14px;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
  gap: 16px;
}

.config-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px;
  background: #f8f9fa;
  border-radius: 8px;
  transition: background 0.2s;
}

html.dark .config-item {
  background: #1e1e2d;
}

.config-label {
  font-size: 12px;
  color: #909399;
  font-weight: 500;
}

html.dark .config-label {
  color: #9ca3af;
}

.config-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;
  word-break: break-all;
}

html.dark .config-value {
  color: #e5e7eb;
}

@media (max-width: 768px) {
  .chart-container {
    height: 250px;
  }

  .config-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
