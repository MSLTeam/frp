<template>
  <div class="server-overview-container">
    <el-row :gutter="20">
      <el-col :md="12">
        <el-card class="server-info-card">
          <div class="server-info-header">
            <h3>服务器信息</h3>
          </div>
          <el-form
            label-position="left"
            label-width="220px"
            class="server-info-form"
          >
            <el-form-item label="版本">
              <span class="info-value">{{ data.version }}</span>
            </el-form-item>
            <el-form-item label="服务端口">
              <span class="info-value">{{ data.bindPort }}</span>
            </el-form-item>
            <el-form-item label="KCP服务端口" v-if="data.kcpBindPort != 0">
              <span class="info-value">{{ data.kcpBindPort }}</span>
            </el-form-item>
            <el-form-item label="QUIC服务端口" v-if="data.quicBindPort != 0">
              <span class="info-value">{{ data.quicBindPort }}</span>
            </el-form-item>
            <el-form-item label="HTTP服务端口" v-if="data.vhostHTTPPort != 0">
              <span class="info-value">{{ data.vhostHTTPPort }}</span>
            </el-form-item>
            <el-form-item label="HTTPS服务端口" v-if="data.vhostHTTPSPort != 0">
              <span class="info-value">{{ data.vhostHTTPSPort }}</span>
            </el-form-item>
            <el-form-item
              label="TCPMux HTTPConnect服务端口"
              v-if="data.tcpmuxHTTPConnectPort != 0"
            >
              <span class="info-value">{{ data.tcpmuxHTTPConnectPort }}</span>
            </el-form-item>
            <el-form-item label="子域名" v-if="data.subdomainHost != ''">
              <LongSpan :content="data.subdomainHost" :length="30"></LongSpan>
            </el-form-item>
            <el-form-item label="最大连接池数量">
              <span class="info-value">{{ data.maxPoolCount }}</span>
            </el-form-item>
            <el-form-item label="客户端最大端口限制数">
              <span class="info-value">{{ data.maxPortsPerClient }}</span>
            </el-form-item>
            <el-form-item label="允许的端口" v-if="data.allowPortsStr != ''">
              <LongSpan class="info-value" :content="data.allowPortsStr" :length="30"></LongSpan>
            </el-form-item>
            <el-form-item label="强制TLS" v-if="data.tlsForce === true">
              <span class="info-value">{{ data.tlsForce }}</span>
            </el-form-item>
            <el-form-item label="心跳包超时">
              <span class="info-value">{{ data.heartbeatTimeout }}</span>
            </el-form-item>
            <el-form-item label="客户端数量">
              <span class="info-value">{{ data.clientCounts }}</span>
            </el-form-item>
            <el-form-item label="当前连接数">
              <span class="info-value">{{ data.curConns }}</span>
            </el-form-item>
            <el-form-item label="隧道数">
              <span class="info-value">{{ data.proxyCounts }}</span>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
      <el-col :md="12">
        <el-card class="chart-card">
          <div id="traffic" class="chart-container"></div>
        </el-card>
        <el-card class="chart-card">
          <div id="proxies" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { DrawTrafficChart, DrawProxyChart } from '../utils/chart'
import LongSpan from './LongSpan.vue'

let data = ref({
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
      if (data.value.maxPortsPerClient == '0') {
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

      // draw chart
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
.server-overview-container {
  padding-left: 20px;
}

.server-info-card {
  margin-bottom: 20px;
}

.chart-card {
  margin-bottom: 20px;
}

.server-info-header {
  margin-bottom: 20px;
}

.server-info-form {
  padding: 20px;
}

.info-value {
  font-weight: bold;
  color: #409eff;
}

.chart-container {
  width: 100%;
  height: 250px;
}

.el-card {
  border-radius: 12px;
}
</style>
