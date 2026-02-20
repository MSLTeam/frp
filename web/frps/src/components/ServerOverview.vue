<template>
  <div class="overview-page">
    <section class="metric-strip">
      <el-card v-for="item in overviewMetrics" :key="item.label" class="metric-card">
        <p class="metric-label">{{ item.label }}</p>
        <div class="metric-main">
          <span class="metric-value">{{ item.value }}</span>
          <span class="metric-unit">{{ item.unit }}</span>
        </div>
      </el-card>
    </section>

    <el-row :gutter="16">
      <el-col :xs="24" :lg="14">
        <el-card class="server-info-card">
          <template #header>
            <div class="card-head">
              <h3>服务器信息</h3>
              <p>节点配置与运行状态</p>
            </div>
          </template>

          <div class="info-grid">
            <div
              v-for="item in serverInfoRows"
              :key="item.label"
              class="info-item"
            >
              <span class="info-label">{{ item.label }}</span>
              <span class="info-value" v-if="!item.long">{{ item.value }}</span>
              <LongSpan v-else :content="String(item.value)" :length="36" />
            </div>
          </div>
        </el-card>
      </el-col>

      <el-col :xs="24" :lg="10">
        <el-card class="chart-card">
          <template #header>
            <div class="card-head compact">
              <h3>网络流量</h3>
              <p>今日</p>
            </div>
          </template>
          <div id="traffic" class="chart-container"></div>
        </el-card>

        <el-card class="chart-card">
          <template #header>
            <div class="card-head compact">
              <h3>隧道分布</h3>
              <p>当前</p>
            </div>
          </template>
          <div id="proxies" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { DrawProxyChart, DrawTrafficChart } from '../utils/chart'
import LongSpan from './LongSpan.vue'

interface ServerOverviewData {
  version: string
  bindPort: number
  kcpBindPort: number
  quicBindPort: number
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  subdomainHost: string
  maxPoolCount: number
  maxPortsPerClient: string
  allowPortsStr: string
  tlsForce: boolean
  heartbeatTimeout: number
  clientCounts: number
  curConns: number
  proxyCounts: number
}

const data = ref<ServerOverviewData>({
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
})

const overviewMetrics = computed(() => {
  return [
    { label: '客户端数量', value: data.value.clientCounts, unit: '个' },
    { label: '当前连接数', value: data.value.curConns, unit: '条' },
    { label: '隧道总数', value: data.value.proxyCounts, unit: '个' },
  ]
})

const serverInfoRows = computed(() => {
  const rows: Array<{ label: string; value: string | number | boolean; long?: boolean }> = [
    { label: '版本', value: data.value.version },
    { label: '服务端口', value: data.value.bindPort },
  ]

  if (data.value.kcpBindPort !== 0) {
    rows.push({ label: 'KCP 服务端口', value: data.value.kcpBindPort })
  }
  if (data.value.quicBindPort !== 0) {
    rows.push({ label: 'QUIC 服务端口', value: data.value.quicBindPort })
  }
  if (data.value.vhostHTTPPort !== 0) {
    rows.push({ label: 'HTTP 服务端口', value: data.value.vhostHTTPPort })
  }
  if (data.value.vhostHTTPSPort !== 0) {
    rows.push({ label: 'HTTPS 服务端口', value: data.value.vhostHTTPSPort })
  }
  if (data.value.tcpmuxHTTPConnectPort !== 0) {
    rows.push({ label: 'TCPMux HTTPConnect 端口', value: data.value.tcpmuxHTTPConnectPort })
  }

  rows.push({ label: '最大连接池数量', value: data.value.maxPoolCount })
  rows.push({ label: '客户端最大端口限制', value: data.value.maxPortsPerClient || '无限制' })
  rows.push({ label: '心跳包超时', value: `${data.value.heartbeatTimeout}s` })

  if (data.value.subdomainHost) {
    rows.push({ label: '子域名', value: data.value.subdomainHost, long: true })
  }
  if (data.value.allowPortsStr) {
    rows.push({ label: '允许端口', value: data.value.allowPortsStr, long: true })
  }
  if (data.value.tlsForce) {
    rows.push({ label: '强制 TLS', value: 'true' })
  }

  return rows
})

const fetchData = () => {
  fetch('../api/serverinfo', { credentials: 'include' })
    .then((res) => res.json())
    .then((json) => {
      data.value.version = json.version
      data.value.bindPort = json.bindPort
      data.value.kcpBindPort = json.kcpBindPort
      data.value.quicBindPort = json.quicBindPort
      data.value.vhostHTTPPort = json.vhostHTTPPort
      data.value.vhostHTTPSPort = json.vhostHTTPSPort
      data.value.tcpmuxHTTPConnectPort = json.tcpmuxHTTPConnectPort
      data.value.subdomainHost = json.subdomainHost
      data.value.maxPoolCount = json.maxPoolCount
      data.value.maxPortsPerClient = json.maxPortsPerClient
      if (data.value.maxPortsPerClient === '0') {
        data.value.maxPortsPerClient = '无限制'
      }
      data.value.allowPortsStr = json.allowPortsStr
      data.value.tlsForce = json.tlsForce
      data.value.heartbeatTimeout = json.heartbeatTimeout
      data.value.clientCounts = json.clientCounts
      data.value.curConns = json.curConns
      data.value.proxyCounts = 0
      if (json.proxyTypeCount != null) {
        if (json.proxyTypeCount.tcp != null) {
          data.value.proxyCounts += json.proxyTypeCount.tcp
        }
        if (json.proxyTypeCount.udp != null) {
          data.value.proxyCounts += json.proxyTypeCount.udp
        }
        if (json.proxyTypeCount.http != null) {
          data.value.proxyCounts += json.proxyTypeCount.http
        }
        if (json.proxyTypeCount.https != null) {
          data.value.proxyCounts += json.proxyTypeCount.https
        }
        if (json.proxyTypeCount.stcp != null) {
          data.value.proxyCounts += json.proxyTypeCount.stcp
        }
        if (json.proxyTypeCount.sudp != null) {
          data.value.proxyCounts += json.proxyTypeCount.sudp
        }
        if (json.proxyTypeCount.xtcp != null) {
          data.value.proxyCounts += json.proxyTypeCount.xtcp
        }
      }

      DrawTrafficChart('traffic', json.totalTrafficIn, json.totalTrafficOut)
      DrawProxyChart('proxies', json)
    })
    .catch(() => {
      ElMessage({
        showClose: true,
        message: '获取服务器信息失败！',
        type: 'warning',
      })
    })
}

fetchData()
</script>

<style scoped>
.overview-page {
  padding: 4px;
}

.metric-strip {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  margin-bottom: 14px;
}

.metric-card {
  position: relative;
  overflow: hidden;
}

.metric-card::before {
  content: '';
  position: absolute;
  right: -24px;
  top: -24px;
  width: 90px;
  height: 90px;
  border-radius: 999px;
  background: radial-gradient(circle, rgba(14, 165, 233, 0.18) 0%, transparent 70%);
}

.metric-label {
  margin: 0;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.4px;
  color: var(--frp-muted);
}

.metric-main {
  display: flex;
  align-items: baseline;
  gap: 6px;
  margin-top: 6px;
}

.metric-value {
  font-size: 34px;
  line-height: 1;
  font-weight: 700;
  color: #0369a1;
}

.metric-unit {
  font-size: 12px;
  font-weight: 600;
  color: var(--frp-muted);
}

.card-head {
  display: flex;
  flex-direction: column;
  gap: 2px;

  h3 {
    margin: 0;
    font-size: 17px;
    line-height: 1.2;
  }

  p {
    margin: 0;
    font-size: 12px;
    color: var(--frp-muted);
  }
}

.card-head.compact {
  h3 {
    font-size: 15px;
  }
}

.server-info-card,
.chart-card {
  margin-bottom: 14px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.info-item {
  min-width: 0;
  border: 1px solid var(--frp-border);
  background: var(--frp-surface-bg);
  border-radius: 12px;
  padding: 10px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.info-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--frp-muted);
}

.info-value {
  font-size: 14px;
  font-weight: 700;
  color: var(--frp-ink);
}

.chart-container {
  width: 100%;
  height: 260px;
}

:global(html.dark) .metric-value {
  color: #7dd3fc;
}

@media (max-width: 1200px) {
  .metric-strip {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 900px) {
  .metric-strip {
    grid-template-columns: 1fr;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .chart-container {
    height: 240px;
  }
}
</style>
