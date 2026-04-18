<template>
  <div
    class="mx-auto flex flex-col gap-5 text-zinc-800 dark:text-zinc-200 h-full pb-10"
  >
    <div
      class="flex flex-col xl:flex-row xl:items-center justify-between gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl shadow-sm border border-zinc-200 dark:border-zinc-800/80"
    >
      <h3
        class="text-lg font-bold tracking-tight m-0 flex items-center gap-2 shrink-0"
      >
        <el-icon class="text-indigo-500"><Connection /></el-icon> 隧道状态
      </h3>

      <div
        class="flex flex-col sm:flex-row sm:items-center sm:justify-end gap-3 w-full xl:w-auto xl:ml-auto"
      >
        <div
          class="flex flex-col sm:flex-row items-center gap-2 w-full sm:w-auto"
        >
          <div class="w-full sm:w-40 lg:w-40 shrink-0">
            <el-select v-model="sortOption" placeholder="默认排序">
              <el-option label="默认排序" value="default" />
              <el-option label="入站流量 ↑" value="trafficInAsc" />
              <el-option label="入站流量 ↓" value="trafficInDesc" />
              <el-option label="出站流量 ↑" value="trafficOutAsc" />
              <el-option label="出站流量 ↓" value="trafficOutDesc" />
              <el-option label="连接数 ↑" value="connsAsc" />
              <el-option label="连接数 ↓" value="connsDesc" />
            </el-select>
          </div>
          <div class="w-full sm:w-40 lg:w-48 shrink-0">
            <el-select
              :model-value="selectedClientKey"
              placeholder="所有客户端"
              clearable
              filterable
              @change="onClientFilterChange"
            >
              <el-option label="全部客户端" value="" />
              <el-option
                v-if="clientIDFilter && !selectedClientInList"
                :label="`${formatSafeText(userFilter ? userFilter + '.' + clientIDFilter : clientIDFilter)} (离线)`"
                :value="selectedClientKey"
                style="color: var(--el-color-warning); font-style: italic"
              />
              <el-option
                v-for="client in clientOptions"
                :key="client.key"
                :label="client.label"
                :value="client.key"
              />
            </el-select>
          </div>

          <div class="w-full flex-1 sm:w-48 lg:w-56">
            <el-input
              v-model="searchText"
              placeholder="搜索隧道..."
              clearable
              :prefix-icon="Search"
            />
          </div>
        </div>

        <div class="flex items-center gap-2 w-full sm:w-auto shrink-0">
          <button
            @click="fetchData"
            class="flex-1 sm:flex-none flex items-center justify-center gap-1.5 px-3 py-1.5 bg-zinc-100 hover:bg-zinc-200 dark:bg-zinc-800 dark:hover:bg-zinc-700 text-zinc-700 dark:text-zinc-300 text-sm font-medium rounded-lg transition-colors border-0 cursor-pointer whitespace-nowrap"
          >
            <el-icon><Refresh /></el-icon> 刷新
          </button>

          <el-popconfirm
            title="确定清理离线隧道？"
            width="200"
            @confirm="clearOfflineProxies"
          >
            <template #reference>
              <button
                class="flex-1 sm:flex-none flex items-center justify-center gap-1.5 px-3 py-1.5 bg-red-50 hover:bg-red-100 dark:bg-red-500/10 dark:hover:bg-red-500/20 text-red-600 dark:text-red-400 text-sm font-medium rounded-lg transition-colors border-0 cursor-pointer whitespace-nowrap"
              >
                <el-icon><Delete /></el-icon> 清理
              </button>
            </template>
          </el-popconfirm>
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 overflow-x-auto custom-scrollbar pb-1">
      <button
        v-for="t in proxyTypes"
        :key="t.value"
        class="px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider transition-colors whitespace-nowrap shrink-0 border-0 cursor-pointer"
        :class="
          activeType === t.value
            ? 'bg-indigo-500 text-white shadow-sm'
            : 'bg-white dark:bg-zinc-900 text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 border border-zinc-200 dark:border-zinc-800'
        "
        @click="activeType = t.value"
      >
        {{ t.label }}
      </button>
    </div>

    <div
      v-loading="loading"
      element-loading-background="rgba(0, 0, 0, 0.0)"
      class="min-h-[200px]"
    >
      <div v-if="filteredProxies.length > 0" class="flex flex-col gap-4">
        <ProxyCard
          v-for="proxy in filteredProxies"
          :key="proxy.name"
          :proxy="proxy"
        />
      </div>

      <div
        v-else-if="!loading"
        class="flex flex-col items-center justify-center py-20 text-zinc-400 dark:text-zinc-500 gap-4 bg-white/50 dark:bg-zinc-900/50 rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800"
      >
        <div
          class="w-16 h-16 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center"
        >
          <el-icon class="text-2xl text-zinc-300 dark:text-zinc-600"
            ><box
          /></el-icon>
        </div>
        <span class="text-sm font-medium">暂无隧道数据</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  Search,
  Refresh,
  Delete,
  Connection,
  Box,
} from '@element-plus/icons-vue'
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
import ProxyCard from '../components/ProxyCard.vue'
import {
  getProxiesByType,
  clearOfflineProxies as apiClearOfflineProxies,
} from '../api/proxy'
import { getServerInfo } from '../api/server'
import { getClients } from '../api/client'
import { Client } from '../utils/client'

const route = useRoute()
const router = useRouter()

const formatSafeText = (text: string | undefined | string[]) => {
  if (!text) return ''
  const str = String(text)
  const matchWithSuffix = str.match(/^.*-(\d+)\.(.+)$/)
  if (matchWithSuffix) {
    const uid = parseInt(matchWithSuffix[1]) - 10000
    return `UID:${uid} · ${matchWithSuffix[2]}`
  }
  const matchTokenOnly = str.match(/^.*-(\d+)$/)
  if (matchTokenOnly) {
    const uid = parseInt(matchTokenOnly[1]) - 10000
    return `UID:${uid}`
  }
  return str
}

const proxyTypes = [
  { label: 'TCP', value: 'tcp' },
  { label: 'UDP', value: 'udp' },
  { label: 'HTTP', value: 'http' },
  { label: 'HTTPS', value: 'https' },
  { label: 'TCPMUX', value: 'tcpmux' },
  { label: 'STCP', value: 'stcp' },
  { label: 'SUDP', value: 'sudp' },
]

const activeType = ref((route.params.type as string) || 'tcp')
const proxies = ref<BaseProxy[]>([])
const clients = ref<Client[]>([])
const loading = ref(false)
const searchText = ref('')
const sortOption = ref('default')
const clientIDFilter = ref((route.query.clientID as string) || '')
const userFilter = ref((route.query.user as string) || '')

const clientOptions = computed(() => {
  return clients.value
    .map((c) => {
      const rawLabel = c.user ? `${c.user}.${c.clientID}` : c.clientID
      return {
        key: c.key,
        clientID: c.clientID,
        user: c.user,
        label: formatSafeText(rawLabel), // 在这里对下拉显示的文字进行格式化
      }
    })
    .sort((a, b) => a.label.localeCompare(b.label))
})

// Compute selected client key for el-select v-model
const selectedClientKey = computed(() => {
  if (!clientIDFilter.value) return ''
  const client = clientOptions.value.find(
    (c) => c.clientID === clientIDFilter.value && c.user === userFilter.value,
  )
  // Return a synthetic key even if not found, so the select shows the filter is active
  return client?.key || `${userFilter.value}:${clientIDFilter.value}`
})

// Check if the filtered client exists in the client list
const selectedClientInList = computed(() => {
  if (!clientIDFilter.value) return true
  return clientOptions.value.some(
    (c) => c.clientID === clientIDFilter.value && c.user === userFilter.value,
  )
})

const filteredProxies = computed(() => {
  let result = proxies.value

  if (clientIDFilter.value) {
    result = result.filter(
      (p) => p.clientID === clientIDFilter.value && p.user === userFilter.value,
    )
  }

  if (searchText.value) {
    const search = searchText.value.toLowerCase()
    result = result.filter((p) => p.name.toLowerCase().includes(search))
  }

  // ====== 排序逻辑 ======
  if (sortOption.value !== 'default') {
    result = [...result].sort((a: any, b: any) => {
      const valAIn = a.trafficIn || 0
      const valBIn = b.trafficIn || 0
      const valAOut = a.trafficOut || 0
      const valBOut = b.trafficOut || 0
      const valAConns = a.conns || 0
      const valBConns = b.conns || 0

      switch (sortOption.value) {
        case 'trafficInAsc':
          return valAIn - valBIn
        case 'trafficInDesc':
          return valBIn - valAIn
        case 'trafficOutAsc':
          return valAOut - valBOut
        case 'trafficOutDesc':
          return valBOut - valAOut
        case 'connsAsc':
          return valAConns - valBConns
        case 'connsDesc':
          return valBConns - valAConns
        default:
          return 0
      }
    })
  }

  return result
})

const onClientFilterChange = (key: string) => {
  if (key) {
    const client = clientOptions.value.find((c) => c.key === key)
    if (client) {
      router.replace({
        query: { ...route.query, clientID: client.clientID, user: client.user },
      })
    }
  } else {
    const query = { ...route.query }
    delete query.clientID
    delete query.user
    router.replace({ query })
  }
}

const fetchClients = async () => {
  try {
    const json = await getClients()
    clients.value = json.map((data) => new Client(data))
  } catch {
    // Ignore errors when fetching clients
  }
}

// Server info cache
let serverInfo: {
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  subdomainHost: string
} | null = null

const fetchServerInfo = async () => {
  if (serverInfo) return serverInfo
  const res = await getServerInfo()
  serverInfo = res
  return serverInfo
}

const fetchData = async () => {
  loading.value = true
  proxies.value = []

  try {
    const type = activeType.value
    const json = await getProxiesByType(type)

    if (type === 'tcp') {
      proxies.value = json.proxies.map((p: any) => new TCPProxy(p))
    } else if (type === 'udp') {
      proxies.value = json.proxies.map((p: any) => new UDPProxy(p))
    } else if (type === 'http') {
      const info = await fetchServerInfo()
      if (info && info.vhostHTTPPort) {
        proxies.value = json.proxies.map(
          (p: any) => new HTTPProxy(p, info.vhostHTTPPort, info.subdomainHost),
        )
      }
    } else if (type === 'https') {
      const info = await fetchServerInfo()
      if (info && info.vhostHTTPSPort) {
        proxies.value = json.proxies.map(
          (p: any) =>
            new HTTPSProxy(p, info.vhostHTTPSPort, info.subdomainHost),
        )
      }
    } else if (type === 'tcpmux') {
      const info = await fetchServerInfo()
      if (info && info.tcpmuxHTTPConnectPort) {
        proxies.value = json.proxies.map(
          (p: any) =>
            new TCPMuxProxy(p, info.tcpmuxHTTPConnectPort, info.subdomainHost),
        )
      }
    } else if (type === 'stcp') {
      proxies.value = json.proxies.map((p: any) => new STCPProxy(p))
    } else if (type === 'sudp') {
      proxies.value = json.proxies.map((p: any) => new SUDPProxy(p))
    }
  } catch (error: any) {
    ElMessage({
      showClose: true,
      message: 'Failed to fetch proxies: ' + error.message,
      type: 'error',
    })
  } finally {
    loading.value = false
  }
}

const clearOfflineProxies = async () => {
  try {
    await apiClearOfflineProxies()
    ElMessage({
      message: 'Successfully cleared offline proxies',
      type: 'success',
    })
    fetchData()
  } catch (err: any) {
    ElMessage({
      message: 'Failed to clear offline proxies: ' + err.message,
      type: 'warning',
    })
  }
}

// Watch for type changes
watch(activeType, (newType) => {
  // Update route but preserve query params
  router.replace({ params: { type: newType }, query: route.query })
  fetchData()
})

// Watch for route query changes (client filter)
watch(
  () => [route.query.clientID, route.query.user],
  ([newClientID, newUser]) => {
    clientIDFilter.value = (newClientID as string) || ''
    userFilter.value = (newUser as string) || ''
  },
)

// Initial fetch
fetchData()
fetchClients()
</script>

<style scoped>
/* 隐藏横向滚动的滚动条但保留功能 */
.custom-scrollbar::-webkit-scrollbar {
  display: none;
}
.custom-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
