<template>
  <div
    class="w-full flex flex-col gap-6 text-zinc-900 dark:text-zinc-100 pb-10"
  >
    <nav
      class="flex items-center text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2"
    >
      <a
        class="flex items-center justify-center w-8 h-8 rounded-lg hover:bg-zinc-100 dark:hover:bg-zinc-800 cursor-pointer transition-colors mr-2"
        @click="goBack"
      >
        <el-icon class="text-lg"><ArrowLeft /></el-icon>
      </a>
      <template v-if="fromClient">
        <router-link
          to="/clients"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >客户端</router-link
        >
        <span class="mx-2.5 text-zinc-300 dark:text-zinc-700">/</span>
        <router-link
          :to="`/clients/${fromClient}`"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >{{ fromClient }}</router-link
        >
        <span class="mx-2.5 text-zinc-300 dark:text-zinc-700">/</span>
      </template>
      <template v-else>
        <router-link
          to="/proxies"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >隧道</router-link
        >
        <span class="mx-2.5 text-zinc-300 dark:text-zinc-700">/</span>
        <router-link
          v-if="proxy?.clientID"
          :to="clientLink"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
        >
          {{ proxy.user ? `${proxy.user}.${proxy.clientID}` : proxy.clientID }}
        </router-link>
        <span
          v-if="proxy?.clientID"
          class="mx-2.5 text-zinc-300 dark:text-zinc-700"
          >/</span
        >
      </template>
      <span class="text-zinc-900 dark:text-zinc-200 font-semibold">{{
        proxyName
      }}</span>
    </nav>

    <div
      v-loading="loading"
      element-loading-background="rgba(0, 0, 0, 0.0)"
      class="min-h-[400px] flex flex-col gap-6"
    >
      <template v-if="proxy">
        <div
          class="relative p-6 sm:p-8 bg-white dark:bg-zinc-900 rounded-3xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm overflow-hidden flex flex-col sm:flex-row items-start sm:items-center gap-6"
        >
          <div
            class="absolute -right-20 -top-20 w-64 h-64 opacity-20 dark:opacity-10 rounded-full blur-3xl pointer-events-none"
            :style="{ background: proxyIconConfig.gradient }"
          ></div>

          <div
            class="w-16 h-16 rounded-2xl flex items-center justify-center text-white text-3xl shadow-lg shrink-0 z-10"
            :style="{ background: proxyIconConfig.gradient }"
          >
            <el-icon><component :is="proxyIconConfig.icon" /></el-icon>
          </div>

          <div class="flex flex-col gap-3 z-10 min-w-0 flex-1">
            <div class="flex items-center flex-wrap gap-3">
              <h1
                class="text-2xl sm:text-3xl font-black tracking-tight m-0 break-all leading-none"
              >
                {{ proxy.name }}
              </h1>
              <span
                class="px-3 py-1 rounded-full bg-zinc-100 dark:bg-zinc-800 text-xs font-bold uppercase tracking-widest text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700/50"
              >
                {{ proxy.type.toUpperCase() }}
              </span>
              <span
                class="px-3 py-1 rounded-md text-[11px] font-bold tracking-widest uppercase border"
                :class="
                  proxy.status === 'online'
                    ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
                    : 'bg-zinc-50 text-zinc-500 border-zinc-200 dark:bg-zinc-800/50 dark:text-zinc-400 dark:border-zinc-700/50'
                "
              >
                {{ proxy.status === 'online' ? '在线' : '离线' }}
              </span>
            </div>

            <div class="flex items-center gap-4">
              <router-link
                v-if="proxy.clientID"
                :to="clientLink"
                class="inline-flex items-center gap-1.5 text-sm font-medium text-zinc-500 dark:text-zinc-400 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
              >
                <el-icon><Monitor /></el-icon>
                <span
                  >客户端:
                  <span class="text-zinc-700 dark:text-zinc-300 break-all">{{
                    proxy.user
                      ? `${proxy.user}.${proxy.clientID}`
                      : proxy.clientID
                  }}</span></span
                >
              </router-link>
            </div>
          </div>
        </div>

        <div
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6"
        >
          <div
            v-if="proxy.port"
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-4 group"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-xs font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >服务端口</span
              >
              <div
                class="w-10 h-10 rounded-xl bg-purple-500/10 text-purple-600 dark:text-purple-400 flex items-center justify-center text-xl transition-transform group-hover:scale-110"
              >
                <el-icon><Connection /></el-icon>
              </div>
            </div>
            <span
              class="text-2xl font-black font-mono text-zinc-900 dark:text-white leading-none"
              >{{ proxy.port }}</span
            >
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-4 group"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-xs font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >当前连接数</span
              >
              <div
                class="w-10 h-10 rounded-xl bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 flex items-center justify-center text-xl transition-transform group-hover:scale-110"
              >
                <el-icon><DataLine /></el-icon>
              </div>
            </div>
            <span
              class="text-2xl font-black font-mono text-zinc-900 dark:text-white leading-none"
              >{{ proxy.conns }}</span
            >
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-4 group"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-xs font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >入站流量</span
              >
              <div
                class="w-10 h-10 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center text-xl transition-transform group-hover:scale-110"
              >
                <el-icon><Bottom /></el-icon>
              </div>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span
                class="text-2xl font-black font-mono text-zinc-900 dark:text-white leading-none"
                >{{ formatTrafficValue(proxy.trafficIn) }}</span
              >
              <span class="text-sm font-bold text-zinc-400">{{
                formatTrafficUnit(proxy.trafficIn)
              }}</span>
            </div>
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-4 group"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-xs font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >出站流量</span
              >
              <div
                class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-xl transition-transform group-hover:scale-110"
              >
                <el-icon><Top /></el-icon>
              </div>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span
                class="text-2xl font-black font-mono text-zinc-900 dark:text-white leading-none"
                >{{ formatTrafficValue(proxy.trafficOut) }}</span
              >
              <span class="text-sm font-bold text-zinc-400">{{
                formatTrafficUnit(proxy.trafficOut)
              }}</span>
            </div>
          </div>
        </div>

        <div
          class="p-6 sm:p-8 bg-white dark:bg-zinc-900 rounded-3xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-6"
        >
          <div class="flex items-center gap-2 text-zinc-900 dark:text-white">
            <el-icon class="text-indigo-500 text-xl"><DataLine /></el-icon>
            <h2 class="text-lg font-bold tracking-tight m-0">状态时间轴</h2>
          </div>
          <div
            class="grid grid-cols-1 md:grid-cols-2 gap-4 bg-zinc-50 dark:bg-zinc-800/40 rounded-2xl p-5 border border-zinc-100 dark:border-zinc-700/50"
          >
            <div class="flex flex-col gap-1.5">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >上次启动时间</span
              >
              <span
                class="text-[15px] font-medium font-mono text-zinc-800 dark:text-zinc-200"
                >{{ proxy.lastStartTime || '-' }}</span
              >
            </div>
            <div class="flex flex-col gap-1.5">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                >上次离线时间</span
              >
              <span
                class="text-[15px] font-medium font-mono text-zinc-800 dark:text-zinc-200"
                >{{ proxy.lastCloseTime || '-' }}</span
              >
            </div>
          </div>
        </div>

        <div class="flex flex-col gap-5">
          <div
            class="flex items-center gap-2 px-2 text-zinc-900 dark:text-white"
          >
            <el-icon class="text-emerald-500 text-xl"><Setting /></el-icon>
            <h2 class="text-lg font-bold tracking-tight m-0">高级配置</h2>
          </div>

          <div
            class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
          >
            <div
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Lock /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >加密</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.encryption ? '已启用' : '未启用' }}</span
                >
              </div>
            </div>

            <div
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Lightning /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >压缩</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.compression ? '已启用' : '未启用' }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.customDomains"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm sm:col-span-2 lg:col-span-1"
            >
              <div
                class="w-10 h-10 rounded-xl bg-purple-500/10 text-purple-600 dark:text-purple-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Link /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >自定义域名</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.customDomains }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.subdomain"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-purple-500/10 text-purple-600 dark:text-purple-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Link /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >子域名</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.subdomain }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.locations"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Location /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >位置</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.locations }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.hostHeaderRewrite"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-orange-500/10 text-orange-600 dark:text-orange-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Tickets /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >HOST 重定向</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.hostHeaderRewrite }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.multiplexer"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-blue-500/10 text-blue-600 dark:text-blue-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Cpu /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >多路复用</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.multiplexer }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.routeByHTTPUser"
              class="flex items-center gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm"
            >
              <div
                class="w-10 h-10 rounded-xl bg-pink-500/10 text-pink-600 dark:text-pink-400 flex items-center justify-center text-lg shrink-0"
              >
                <el-icon><Connection /></el-icon>
              </div>
              <div class="flex flex-col gap-1 min-w-0">
                <span
                  class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                  >用户 HTTP 路由</span
                >
                <span
                  class="text-sm font-bold text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.routeByHTTPUser }}</span
                >
              </div>
            </div>
          </div>

          <div
            v-if="proxy.annotations && proxy.annotations.size > 0"
            class="flex flex-wrap gap-2 mt-2 px-2"
          >
            <span
              v-for="[key, value] in proxy.annotations"
              :key="key"
              class="inline-flex items-center px-3 py-1.5 bg-zinc-100 dark:bg-zinc-800 rounded-lg text-xs font-medium text-zinc-600 dark:text-zinc-300 border border-zinc-200 dark:border-zinc-700/50"
            >
              <span class="text-zinc-400 mr-1">{{ key }}:</span> {{ value }}
            </span>
          </div>
        </div>

        <div
          class="p-6 sm:p-8 bg-white dark:bg-zinc-900 rounded-3xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm flex flex-col gap-6"
        >
          <div
            class="flex items-center gap-2 text-zinc-900 dark:text-white pb-6 border-b border-zinc-100 dark:border-zinc-800/80"
          >
            <el-icon class="text-blue-500 text-xl"><DataLine /></el-icon>
            <h2 class="text-lg font-bold tracking-tight m-0">流量监控</h2>
          </div>
          <div class="w-full">
            <Traffic :proxy-name="proxyName" />
          </div>
        </div>
      </template>

      <div
        v-else-if="!loading"
        class="flex flex-col items-center justify-center py-20 text-zinc-400 dark:text-zinc-500 gap-5 bg-white/50 dark:bg-zinc-900/50 rounded-3xl border border-dashed border-zinc-200 dark:border-zinc-800 mt-10"
      >
        <div
          class="w-20 h-20 rounded-full bg-zinc-100 dark:bg-zinc-800 flex items-center justify-center shadow-inner"
        >
          <el-icon class="text-3xl text-zinc-300 dark:text-zinc-600"
            ><Inbox
          /></el-icon>
        </div>
        <div class="flex flex-col items-center gap-1">
          <span class="text-lg font-bold text-zinc-700 dark:text-zinc-200"
            >隧道未找到</span
          >
          <span class="text-sm font-medium">隧道不存在或已经被移除</span>
        </div>
        <button
          @click="goBack"
          class="mt-4 px-6 py-2.5 bg-zinc-900 dark:bg-white hover:bg-zinc-800 dark:hover:bg-zinc-100 text-white dark:text-zinc-900 text-sm font-bold rounded-xl transition-colors shadow-md shadow-zinc-900/20 dark:shadow-white/20 border-0 cursor-pointer"
        >
          返回列表
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Monitor,
  Connection,
  DataLine,
  Bottom,
  Top,
  Link,
  Lock,
  Promotion,
  Grid,
  Setting,
  Cpu,
  Lightning,
  Tickets,
  Location,
} from '@element-plus/icons-vue'
import { getProxyByName } from '../api/proxy'
import { getServerInfo } from '../api/server'
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
import Traffic from '../components/Traffic.vue'

const route = useRoute()
const router = useRouter()
const proxyName = computed(() => route.params.name as string)
const fromClient = computed(() => {
  if (route.query.from === 'client' && route.query.client) {
    return route.query.client as string
  }
  return null
})
const proxy = ref<BaseProxy | null>(null)
const loading = ref(true)

const goBack = () => {
  if (window.history.length > 1) {
    router.back()
  } else {
    router.push('/proxies')
  }
}

let serverInfo: {
  vhostHTTPPort: number
  vhostHTTPSPort: number
  tcpmuxHTTPConnectPort: number
  subdomainHost: string
} | null = null

const clientLink = computed(() => {
  if (!proxy.value) return ''
  const key = proxy.value.user
    ? `${proxy.value.user}.${proxy.value.clientID}`
    : proxy.value.clientID
  return `/clients/${key}`
})

const proxyIconConfig = computed(() => {
  const type = proxy.value?.type?.toLowerCase() || ''
  const configs: Record<string, { icon: any; gradient: string }> = {
    tcp: {
      icon: Connection,
      gradient: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
    },
    udp: {
      icon: Promotion,
      gradient: 'linear-gradient(135deg, #8b5cf6 0%, #6d28d9 100%)',
    },
    http: {
      icon: Link,
      gradient: 'linear-gradient(135deg, #22c55e 0%, #16a34a 100%)',
    },
    https: {
      icon: Lock,
      gradient: 'linear-gradient(135deg, #14b8a6 0%, #0d9488 100%)',
    },
    stcp: {
      icon: Lock,
      gradient: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
    },
    sudp: {
      icon: Lock,
      gradient: 'linear-gradient(135deg, #f97316 0%, #ea580c 100%)',
    },
    tcpmux: {
      icon: Grid,
      gradient: 'linear-gradient(135deg, #06b6d4 0%, #0891b2 100%)',
    },
    xtcp: {
      icon: Connection,
      gradient: 'linear-gradient(135deg, #ec4899 0%, #db2777 100%)',
    },
  }
  return (
    configs[type] || {
      icon: Connection,
      gradient: 'linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%)',
    }
  )
})

const formatTrafficValue = (bytes: number): string => {
  if (bytes === 0) return '0'
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const value = bytes / Math.pow(k, i)
  return value < 10 ? value.toFixed(1) : Math.round(value).toString()
}

const formatTrafficUnit = (bytes: number): string => {
  if (bytes === 0) return 'B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return units[i]
}

const fetchServerInfo = async () => {
  if (serverInfo) return serverInfo
  const res = await getServerInfo()
  serverInfo = res
  return serverInfo
}

const fetchProxy = async () => {
  const name = proxyName.value
  if (!name) {
    loading.value = false
    return
  }

  try {
    const data = await getProxyByName(name)
    const info = await fetchServerInfo()
    const type = data.conf?.type || ''

    if (type === 'tcp') {
      proxy.value = new TCPProxy(data)
    } else if (type === 'udp') {
      proxy.value = new UDPProxy(data)
    } else if (type === 'http' && info?.vhostHTTPPort) {
      proxy.value = new HTTPProxy(data, info.vhostHTTPPort, info.subdomainHost)
    } else if (type === 'https' && info?.vhostHTTPSPort) {
      proxy.value = new HTTPSProxy(
        data,
        info.vhostHTTPSPort,
        info.subdomainHost,
      )
    } else if (type === 'tcpmux' && info?.tcpmuxHTTPConnectPort) {
      proxy.value = new TCPMuxProxy(
        data,
        info.tcpmuxHTTPConnectPort,
        info.subdomainHost,
      )
    } else if (type === 'stcp') {
      proxy.value = new STCPProxy(data)
    } else if (type === 'sudp') {
      proxy.value = new SUDPProxy(data)
    } else {
      proxy.value = new BaseProxy(data)
      proxy.value.type = type
    }
  } catch (error: any) {
    ElMessage.error('Failed to fetch proxy: ' + error.message)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchProxy()
})
</script>

<style scoped>

</style>
