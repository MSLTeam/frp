<template>
  <div
    class="w-full flex flex-col gap-5 text-zinc-900 dark:text-zinc-100 pb-10"
  >
    <nav
      class="flex items-center text-[13px] font-medium text-zinc-500 dark:text-zinc-400 mb-1"
    >
      <a
        class="flex items-center justify-center w-7 h-7 rounded hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer transition-colors mr-1.5"
        @click="goBack"
      >
        <el-icon class="text-base"><ArrowLeft /></el-icon>
      </a>
      <router-link
        to="/clients"
        class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
        >客户端</router-link
      >
      <span class="mx-2 text-zinc-300 dark:text-zinc-700">/</span>
      <span class="text-zinc-800 dark:text-zinc-200 font-semibold">{{
        client?.displayName || route.params.key
      }}</span>
    </nav>

    <div
      v-loading="loading"
      element-loading-background="rgba(0, 0, 0, 0.0)"
      class="min-h-[400px] flex flex-col gap-5"
    >
      <template v-if="client">
        <div
          class="p-5 sm:p-6 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col sm:flex-row items-start sm:items-center justify-between gap-5"
        >
          <div class="flex items-center gap-4 min-w-0">
            <div
              class="w-12 h-12 rounded-xl flex items-center justify-center text-xl font-bold shrink-0 border bg-zinc-50 border-zinc-200 text-zinc-700 dark:bg-zinc-800 dark:border-zinc-700 dark:text-zinc-300"
            >
              {{ client.displayName.charAt(0).toUpperCase() }}
            </div>

            <div class="flex flex-col min-w-0">
              <h1
                class="text-xl sm:text-2xl font-bold tracking-tight m-0 truncate text-zinc-800 dark:text-zinc-100 leading-tight"
              >
                {{ client.displayName }}
              </h1>
              <div
                class="flex items-center gap-2.5 mt-1.5 text-[13px] font-medium text-zinc-500 dark:text-zinc-400"
              >
                <span v-if="client.ip" class="flex items-center font-mono">
                  {{ client.ip }}
                </span>
                <span
                  v-if="client.ip && client.hostname"
                  class="w-px h-3 bg-zinc-300 dark:bg-zinc-700"
                ></span>
                <span
                  v-if="client.hostname"
                  class="truncate max-w-[150px] sm:max-w-[250px]"
                  :title="client.hostname"
                >
                  {{ client.hostname }}
                </span>
              </div>
            </div>
          </div>

          <div class="shrink-0">
            <span
              class="px-2.5 py-1 rounded-md text-[11px] font-bold tracking-widest uppercase border"
              :class="
                client.online
                  ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
                  : 'bg-zinc-50 text-zinc-500 border-zinc-200 dark:bg-zinc-800/50 dark:text-zinc-400 dark:border-zinc-700'
              "
            >
              {{ client.online ? '在线' : '离线' }}
            </span>
          </div>
        </div>

        <div
          class="grid grid-cols-2 lg:grid-cols-4 gap-4 bg-white dark:bg-zinc-900 rounded-2xl p-5 border border-zinc-200 dark:border-zinc-800"
        >
          <div class="flex flex-col gap-1.5">
            <span
              class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
              >连接数</span
            >
            <span
              class="text-lg font-black font-mono text-zinc-800 dark:text-zinc-200 leading-none"
              >{{ totalConnections }}</span
            >
          </div>
          <div class="flex flex-col gap-1.5 min-w-0">
            <span
              class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
              >运行 ID</span
            >
            <span
              class="text-lg font-black font-mono text-zinc-800 dark:text-zinc-200 leading-none truncate"
              :title="client.runID"
              >{{ client.runID || '-' }}</span
            >
          </div>
          <div class="flex flex-col gap-1.5">
            <span
              class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
              >首次连接</span
            >
            <span
              class="text-lg font-black font-mono text-zinc-800 dark:text-zinc-200 leading-none"
              >{{ client.firstConnectedAgo || '-' }}</span
            >
          </div>
          <div class="flex flex-col gap-1.5">
            <span
              class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
              >{{ client.online ? '已连接' : '离线时长' }}</span
            >
            <span
              class="text-lg font-black font-mono text-zinc-800 dark:text-zinc-200 leading-none"
              >{{
                client.online ? client.lastConnectedAgo : client.disconnectedAgo
              }}</span
            >
          </div>
        </div>

        <div
          class="p-5 sm:p-6 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col gap-5"
        >
          <div
            class="flex flex-col sm:flex-row sm:items-center justify-between gap-4 pb-4 border-b border-zinc-100 dark:border-zinc-800/80"
          >
            <div class="flex items-center gap-2.5">
              <span
                class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
              ></span>
              <h3
                class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
              >
                所属隧道
              </h3>
              <span
                class="px-2 py-0.5 rounded text-[10px] font-bold text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700 bg-zinc-50 dark:bg-zinc-800/50"
              >
                {{ filteredProxies.length }}
              </span>
            </div>

            <div class="w-full sm:w-56 shrink-0">
              <el-input
                v-model="proxySearch"
                placeholder="搜索隧道..."
                :prefix-icon="Search"
                clearable
                class="custom-el-input w-full"
              />
            </div>
          </div>

          <div class="w-full">
            <div
              v-if="proxiesLoading"
              class="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500 gap-3"
            >
              <el-icon class="is-loading text-2xl"><Loading /></el-icon>
              <span class="text-[11px] font-bold tracking-widest uppercase"
                >加载中...</span
              >
            </div>

            <div
              v-else-if="filteredProxies.length > 0"
              class="flex flex-col gap-3"
            >
              <ProxyCard
                v-for="proxy in filteredProxies"
                :key="proxy.name"
                :proxy="proxy"
                show-type
              />
            </div>

            <div
              v-else-if="clientProxies.length > 0"
              class="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500 gap-3 bg-zinc-50 dark:bg-zinc-800/30 rounded-xl border border-dashed border-zinc-200 dark:border-zinc-700/50"
            >
              <el-icon class="text-3xl text-zinc-300 dark:text-zinc-600"
                ><Search
              /></el-icon>
              <span class="text-[13px] font-medium"
                >未找到匹配的隧道 "{{ proxySearch }}"</span
              >
            </div>

            <div
              v-else
              class="flex flex-col items-center justify-center py-12 text-zinc-400 dark:text-zinc-500 gap-3 bg-zinc-50 dark:bg-zinc-800/30 rounded-xl border border-dashed border-zinc-200 dark:border-zinc-700/50"
            >
              <span class="text-[13px] font-medium">该客户端暂无隧道数据</span>
            </div>
          </div>
        </div>
      </template>

      <div
        v-else-if="!loading"
        class="flex flex-col items-center justify-center py-20 text-zinc-400 gap-4 bg-white/50 dark:bg-zinc-900/50 rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800 mt-5"
      >
        <el-icon class="text-4xl text-zinc-300 dark:text-zinc-600"
          ><Search
        /></el-icon>
        <div class="flex flex-col items-center gap-1">
          <span class="text-base font-bold text-zinc-700 dark:text-zinc-200"
            >客户端未找到</span
          >
          <span class="text-[13px]">此客户端不存在或已被移除</span>
        </div>
        <button
          @click="goBack"
          class="mt-2 px-5 py-1.5 bg-zinc-100 dark:bg-zinc-800 hover:bg-zinc-200 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300 text-[13px] font-medium rounded-lg transition-colors border-0 cursor-pointer"
        >
          返回列表
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ArrowLeft, Loading, Search } from '@element-plus/icons-vue'
import { Client } from '../utils/client'
import { getClient } from '../api/client'
import { getProxiesByType } from '../api/proxy'
import {
  BaseProxy,
  TCPProxy,
  UDPProxy,
  HTTPProxy,
  HTTPSProxy,
  TCPMuxProxy,
  STCPProxy,
  SUDPProxy,
} from '../utils/proxy'
import { getServerInfo } from '../api/server'
import ProxyCard from '../components/ProxyCard.vue'

const route = useRoute()
const router = useRouter()
const client = ref<Client | null>(null)
const loading = ref(true)

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/clients')
  }
}
const proxiesLoading = ref(false)
const allProxies = ref<BaseProxy[]>([])
const proxySearch = ref('')

let serverInfo: {
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  subdomainHost: string
} | null = null

const clientProxies = computed(() => {
  if (!client.value) return []
  return allProxies.value.filter(
    (p) =>
      p.clientID === client.value!.clientID && p.user === client.value!.user,
  )
})

const filteredProxies = computed(() => {
  if (!proxySearch.value) return clientProxies.value
  const search = proxySearch.value.toLowerCase()
  return clientProxies.value.filter(
    (p) =>
      p.name.toLowerCase().includes(search) ||
      p.type.toLowerCase().includes(search),
  )
})

const totalConnections = computed(() => {
  return clientProxies.value.reduce((sum, p) => sum + p.conns, 0)
})

const fetchServerInfo = async () => {
  if (serverInfo) return serverInfo
  const res = await getServerInfo()
  serverInfo = res
  return serverInfo
}

const fetchClient = async () => {
  const key = route.params.key as string
  if (!key) {
    loading.value = false
    return
  }
  try {
    const data = await getClient(key)
    client.value = new Client(data)
  } catch (error: any) {
    ElMessage.error('Failed to fetch client: ' + error.message)
  } finally {
    loading.value = false
  }
}

const fetchProxies = async () => {
  proxiesLoading.value = true
  const proxyTypes = ['tcp', 'udp', 'http', 'https', 'tcpmux', 'stcp', 'sudp']
  const proxies: BaseProxy[] = []
  try {
    const info = await fetchServerInfo()
    for (const type of proxyTypes) {
      try {
        const json = await getProxiesByType(type)
        if (!json.proxies) continue
        if (type === 'tcp') {
          proxies.push(...json.proxies.map((p: any) => new TCPProxy(p)))
        } else if (type === 'udp') {
          proxies.push(...json.proxies.map((p: any) => new UDPProxy(p)))
        } else if (type === 'http' && info?.vhostHTTPPort) {
          proxies.push(
            ...json.proxies.map(
              (p: any) =>
                new HTTPProxy(p, info.vhostHTTPPort, info.subdomainHost),
            ),
          )
        } else if (type === 'https' && info?.vhostHTTPSPort) {
          proxies.push(
            ...json.proxies.map(
              (p: any) =>
                new HTTPSProxy(p, info.vhostHTTPSPort, info.subdomainHost),
            ),
          )
        } else if (type === 'tcpmux' && info?.tcpmuxHTTPConnectPort) {
          proxies.push(
            ...json.proxies.map(
              (p: any) =>
                new TCPMuxProxy(
                  p,
                  info.tcpmuxHTTPConnectPort,
                  info.subdomainHost,
                ),
            ),
          )
        } else if (type === 'stcp') {
          proxies.push(...json.proxies.map((p: any) => new STCPProxy(p)))
        } else if (type === 'sudp') {
          proxies.push(...json.proxies.map((p: any) => new SUDPProxy(p)))
        }
      } catch {
        // Ignore
      }
    }
    allProxies.value = proxies
  } catch {
    // Ignore
  } finally {
    proxiesLoading.value = false
  }
}

onMounted(() => {
  fetchClient()
  fetchProxies()
})
</script>

<style scoped>
.client-detail-page {
}

/* Breadcrumb */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  margin-bottom: 24px;
}

.breadcrumb-link {
  display: flex;
  align-items: center;
  color: var(--text-secondary);
  cursor: pointer;
  transition: color 0.2s;
  margin-right: 4px;
}

.breadcrumb-link:hover {
  color: var(--text-primary);
}

.breadcrumb-item {
  color: var(--text-secondary);
  text-decoration: none;
  transition: color 0.2s;
}

.breadcrumb-item:hover {
  color: var(--el-color-primary);
}

.breadcrumb-separator {
  color: var(--el-border-color);
}

.breadcrumb-current {
  color: var(--text-primary);
  font-weight: 500;
}

/* Card Base */
.header-card,
.proxies-card {
  background: var(--el-bg-color);
  border: 1px solid var(--header-border);
  border-radius: 12px;
  margin-bottom: 16px;
}

/* Header Card */
.header-main {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 24px;
}

.header-left {
  display: flex;
  gap: 16px;
  align-items: center;
}

.client-avatar {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  font-weight: 500;
  flex-shrink: 0;
}

.client-info {
  min-width: 0;
}

.client-name {
  font-size: 20px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 4px 0;
  line-height: 1.3;
}

.client-meta {
  display: flex;
  gap: 12px;
  font-size: 14px;
  color: var(--text-secondary);
}

.status-badge {
  padding: 6px 12px;
  border-radius: 6px;
  font-size: 13px;
  font-weight: 500;
}

.status-badge.online {
  background: rgba(34, 197, 94, 0.1);
  color: #16a34a;
}

.status-badge.offline {
  background: var(--hover-bg);
  color: var(--text-secondary);
}

html.dark .status-badge.online {
  background: rgba(34, 197, 94, 0.15);
  color: #4ade80;
}

/* Info Section */
.info-section {
  display: flex;
  flex-wrap: wrap;
  gap: 16px 32px;
  padding: 16px 24px;
}

.info-item {
  display: flex;
  align-items: baseline;
  gap: 8px;
}

.info-label {
  font-size: 13px;
  color: var(--text-secondary);
}

.info-label::after {
  content: ':';
}

.info-value {
  font-size: 13px;
  color: var(--text-primary);
  font-weight: 500;
  word-break: break-all;
}

/* Proxies Card */
.proxies-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  gap: 16px;
}

.proxies-title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.proxies-title h2 {
  font-size: 15px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0;
}

.proxies-count {
  font-size: 13px;
  font-weight: 500;
  color: var(--text-secondary);
  background: var(--hover-bg);
  padding: 4px 10px;
  border-radius: 6px;
}

.proxy-search {
  width: 200px;
}

.proxy-search :deep(.el-input__wrapper) {
  border-radius: 6px;
}

.proxies-body {
  padding: 16px;
}

.proxies-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 40px;
  color: var(--text-secondary);
}

.empty-state {
  text-align: center;
  padding: 40px;
  color: var(--text-secondary);
}

.empty-state p {
  margin: 0;
}

/* Not Found */
.not-found {
  text-align: center;
  padding: 60px 20px;
}

.not-found h2 {
  font-size: 18px;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0 0 8px;
}

.not-found p {
  font-size: 14px;
  color: var(--text-secondary);
  margin: 0 0 20px;
}

/* Responsive */
@media (max-width: 640px) {
  .header-main {
    flex-direction: column;
    gap: 16px;
  }

  .header-right {
    align-self: flex-start;
  }
}
</style>
