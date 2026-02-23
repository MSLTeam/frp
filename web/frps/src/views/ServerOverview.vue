<template>
  <div class="flex flex-col gap-6 w-full pb-10">
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
      <StatCard
        label="客户端"
        :value="data.clientCounts"
        type="clients"
        subtitle="已连接的客户端数量"
        to="/clients"
      />
      <StatCard
        label="隧道"
        :value="data.proxyCounts"
        type="proxies"
        subtitle="已连接的隧道数"
        to="/proxies/tcp"
      />
      <StatCard
        label="连接"
        :value="data.curConns"
        type="connections"
        subtitle="隧道被连接总数"
      />
      <StatCard
        label="流量"
        :value="formatTrafficTotal()"
        type="traffic"
        subtitle="今日流量"
      />
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
      <div
        class="p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col gap-6"
      >
        <div class="flex justify-between items-center">
          <div class="flex items-center gap-2.5">
            <span
              class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
            ></span>
            <h3
              class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
            >
              节点流量监控
            </h3>
          </div>
          <span
            class="text-[10px] font-bold px-2 py-0.5 border border-zinc-200 dark:border-zinc-700 text-zinc-500 dark:text-zinc-400 rounded uppercase"
            >Today</span
          >
        </div>

        <div
          class="flex-1 flex flex-col sm:flex-row items-center justify-center gap-8 sm:gap-16"
        >
          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 border border-blue-100 dark:border-blue-900 bg-blue-50/50 dark:bg-blue-500/5 text-blue-600 dark:text-blue-400 rounded-lg flex items-center justify-center text-lg"
            >
              <el-icon><Download /></el-icon>
            </div>
            <div class="flex flex-col">
              <span
                class="text-[11px] font-semibold text-zinc-500 dark:text-zinc-400 uppercase tracking-widest"
                >Upload</span
              >
              <span
                class="text-2xl font-black font-mono tracking-tight text-zinc-800 dark:text-zinc-100"
                >{{ formatFileSize(data.totalTrafficIn) }}</span
              >
            </div>
          </div>

          <div class="flex items-center gap-3">
            <div
              class="w-10 h-10 border border-emerald-100 dark:border-emerald-900 bg-emerald-50/50 dark:bg-emerald-500/5 text-emerald-600 dark:text-emerald-400 rounded-lg flex items-center justify-center text-lg"
            >
              <el-icon><Upload /></el-icon>
            </div>
            <div class="flex flex-col">
              <span
                class="text-[11px] font-semibold text-zinc-500 dark:text-zinc-400 uppercase tracking-widest"
                >Download</span
              >
              <span
                class="text-2xl font-black font-mono tracking-tight text-zinc-800 dark:text-zinc-100"
                >{{ formatFileSize(data.totalTrafficOut) }}</span
              >
            </div>
          </div>
        </div>
      </div>

      <div
        class="p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col gap-6"
      >
        <div class="flex justify-between items-center">
          <div class="flex items-center gap-2.5">
            <span
              class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
            ></span>
            <h3
              class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
            >
              活跃隧道分布
            </h3>
          </div>
          <span
            class="text-[10px] font-bold px-2 py-0.5 border border-zinc-200 dark:border-zinc-700 text-zinc-500 dark:text-zinc-400 rounded uppercase"
            >Live</span
          >
        </div>

        <div class="flex-1 flex flex-col justify-center">
          <div
            v-if="hasActiveProxies"
            class="grid grid-cols-3 sm:grid-cols-4 gap-3"
          >
            <div
              v-for="(count, type) in data.proxyTypeCounts"
              :key="type"
              v-show="count > 0"
              class="flex flex-col items-center justify-center p-3 border border-zinc-100 dark:border-zinc-800 bg-zinc-50 dark:bg-zinc-800/30 rounded-xl"
            >
              <span
                class="text-[10px] font-bold text-zinc-400 dark:text-zinc-500 uppercase tracking-wider mb-1"
                >{{ type }}</span
              >
              <span
                class="text-lg font-black text-zinc-800 dark:text-zinc-200 leading-none"
                >{{ count }}</span
              >
            </div>
          </div>
          <div
            v-else
            class="flex flex-col items-center justify-center h-full text-zinc-400 py-6 gap-2"
          >
            <el-icon class="text-2xl text-zinc-300 dark:text-zinc-600"
              ><Box
            /></el-icon>
            <span class="text-xs font-medium">暂无活跃隧道</span>
          </div>
        </div>
      </div>
    </div>

    <div
      class="p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl"
    >
      <div
        class="flex justify-between items-center mb-6 border-b border-zinc-100 dark:border-zinc-800/80 pb-4"
      >
        <div class="flex items-center gap-2.5">
          <span
            class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
          ></span>
          <h3
            class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
          >
            服务端口与配置参数
          </h3>
        </div>
        <span
          class="text-[11px] font-mono font-bold px-2.5 py-1 bg-zinc-100 dark:bg-zinc-800 text-zinc-700 dark:text-zinc-300 rounded-md"
        >
          v{{ data.version }}
        </span>
      </div>

      <div
        class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-x-4 gap-y-5"
      >
        <div class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >服务端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.bindPort }}</span
          >
        </div>

        <div v-if="data.kcpBindPort != 0" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >KCP 端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.kcpBindPort }}</span
          >
        </div>

        <div v-if="data.quicBindPort != 0" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >QUIC 端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.quicBindPort }}</span
          >
        </div>

        <div v-if="data.vhostHTTPPort != 0" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >HTTP 端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.vhostHTTPPort }}</span
          >
        </div>

        <div v-if="data.vhostHTTPSPort != 0" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >HTTPS 端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.vhostHTTPSPort }}</span
          >
        </div>

        <div v-if="data.tcpmuxHTTPConnectPort != 0" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >TCPMux 端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.tcpmuxHTTPConnectPort }}</span
          >
        </div>

        <div class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >最大连接池</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.maxPoolCount }}</span
          >
        </div>

        <div class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >单客户端最大端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.maxPortsPerClient }}</span
          >
        </div>

        <div class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >心跳包超时</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200"
            >{{ data.heartbeatTimeout }}s</span
          >
        </div>

        <div v-if="data.tlsForce" class="flex flex-col gap-1">
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >强制 TLS</span
          >
          <span
            class="text-[11px] font-bold text-zinc-600 dark:text-zinc-300 mt-0.5"
            >已启用</span
          >
        </div>

        <div
          v-if="data.subdomainHost != ''"
          class="flex flex-col gap-1 col-span-2 sm:col-span-1 lg:col-span-2"
        >
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >根域名 (Subdomain)</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200 truncate"
            :title="data.subdomainHost"
            >{{ data.subdomainHost }}</span
          >
        </div>

        <div
          v-if="data.allowPortsStr != ''"
          class="flex flex-col gap-1 col-span-2 sm:col-span-1 lg:col-span-2"
        >
          <span
            class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase"
            >可用连接端口</span
          >
          <span
            class="text-sm font-mono font-medium text-zinc-800 dark:text-zinc-200 truncate"
            :title="data.allowPortsStr"
            >{{ data.allowPortsStr }}</span
          >
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { formatFileSize } from '../utils/format'
import { Download, Upload, Box } from '@element-plus/icons-vue'
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

<style scoped></style>
