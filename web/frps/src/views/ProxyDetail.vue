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
      <template v-if="fromClient">
        <router-link
          to="/clients"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >客户端</router-link
        >
        <span class="mx-2 text-zinc-300 dark:text-zinc-700">/</span>
        <router-link
          :to="`/clients/${fromClient}`"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >{{ formatSafeText(fromClient) }}</router-link
        >
        <span class="mx-2 text-zinc-300 dark:text-zinc-700">/</span>
      </template>
      <template v-else>
        <router-link
          to="/proxies"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >隧道</router-link
        >
        <span class="mx-2 text-zinc-300 dark:text-zinc-700">/</span>
        <router-link
          v-if="proxy?.clientID"
          :to="clientLink"
          class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
        >
          <router-link
            v-if="proxy?.clientID"
            :to="clientLink"
            class="hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
          >
            {{
              formatSafeText(
                proxy.user ? `${proxy.user}.${proxy.clientID}` : proxy.clientID,
              )
            }}
          </router-link>
        </router-link>
        <span
          v-if="proxy?.clientID"
          class="mx-2 text-zinc-300 dark:text-zinc-700"
          >/</span
        >
      </template>
      <span class="text-zinc-800 dark:text-zinc-200 font-semibold">{{
        formatSafeText(proxyName)
      }}</span>
    </nav>

    <div
      v-loading="loading"
      element-loading-background="rgba(0, 0, 0, 0.0)"
      class="min-h-[400px] flex flex-col gap-5"
    >
      <template v-if="proxy">
        <div
          class="p-5 sm:p-6 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col sm:flex-row items-start sm:items-center gap-5"
        >
          <div
            class="w-14 h-14 rounded-xl flex items-center justify-center text-xl shrink-0 border bg-zinc-50 border-zinc-200 text-zinc-700 dark:bg-zinc-800 dark:border-zinc-700 dark:text-zinc-300"
          >
            <el-icon><component :is="proxyIconConfig.icon" /></el-icon>
          </div>

          <div class="flex flex-col gap-2 min-w-0 flex-1">
            <div class="flex items-center flex-wrap gap-2.5">
              <h1
                class="text-xl sm:text-2xl font-bold tracking-tight m-0 break-all leading-none text-zinc-800 dark:text-zinc-100"
              >
                {{ formatSafeText(proxy.name) }}
              </h1>
              <span
                class="px-2 py-0.5 rounded text-[10px] font-bold uppercase tracking-widest bg-zinc-100 text-zinc-500 dark:bg-zinc-800 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700"
              >
                {{ proxy.type.toUpperCase() }}
              </span>
              <span
                class="px-2 py-0.5 rounded text-[10px] font-bold tracking-widest uppercase border"
                :class="
                  proxy.status === 'online'
                    ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
                    : 'bg-zinc-50 text-zinc-500 border-zinc-200 dark:bg-zinc-800/50 dark:text-zinc-400 dark:border-zinc-700'
                "
              >
                {{ proxy.status === 'online' ? '在线' : '离线' }}
              </span>
            </div>

            <div class="flex items-center gap-3 mt-1">
              <router-link
                v-if="proxy.clientID"
                :to="clientLink"
                class="inline-flex items-center gap-1.5 text-[13px] font-medium text-zinc-500 dark:text-zinc-400 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors"
              >
                <el-icon><Monitor /></el-icon>
                <span
                  >归属:
                  <span class="text-zinc-700 dark:text-zinc-300 font-mono">{{
                    formatSafeText(
                      proxy.user
                        ? `${proxy.user}.${proxy.clientID}`
                        : proxy.clientID,
                    )
                  }}</span></span
                >
              </router-link>
            </div>
          </div>
        </div>

        <div
          class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-5"
        >
          <div
            v-if="proxy.port"
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col gap-3"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 mt-1"
                >服务端口</span
              >
              <div class="text-zinc-400 dark:text-zinc-500">
                <el-icon><Connection /></el-icon>
              </div>
            </div>
            <span
              class="text-[28px] font-black font-mono text-zinc-800 dark:text-zinc-100 leading-none"
              >{{ proxy.port }}</span
            >
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col gap-3"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 mt-1"
                >当前连接数</span
              >
              <div class="text-zinc-400 dark:text-zinc-500">
                <el-icon><DataLine /></el-icon>
              </div>
            </div>
            <span
              class="text-[28px] font-black font-mono text-zinc-800 dark:text-zinc-100 leading-none"
              >{{ proxy.conns }}</span
            >
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col gap-3"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 mt-1"
                >入站流量</span
              >
              <div class="text-blue-500 dark:text-blue-400">
                <el-icon><Bottom /></el-icon>
              </div>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span
                class="text-[28px] font-black font-mono text-zinc-800 dark:text-zinc-100 leading-none"
                >{{ formatTrafficValue(proxy.trafficIn) }}</span
              >
              <span class="text-xs font-bold text-zinc-400">{{
                formatTrafficUnit(proxy.trafficIn)
              }}</span>
            </div>
          </div>

          <div
            class="p-5 bg-white dark:bg-zinc-900 rounded-2xl border border-zinc-200 dark:border-zinc-800 flex flex-col gap-3"
          >
            <div class="flex justify-between items-start">
              <span
                class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 mt-1"
                >出站流量</span
              >
              <div class="text-emerald-500 dark:text-emerald-400">
                <el-icon><Top /></el-icon>
              </div>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span
                class="text-[28px] font-black font-mono text-zinc-800 dark:text-zinc-100 leading-none"
                >{{ formatTrafficValue(proxy.trafficOut) }}</span
              >
              <span class="text-xs font-bold text-zinc-400">{{
                formatTrafficUnit(proxy.trafficOut)
              }}</span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-3 gap-5">
          <div
            class="lg:col-span-1 p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col"
          >
            <div class="flex items-center gap-2 mb-5">
              <span
                class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
              ></span>
              <h3
                class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
              >
                状态时间轴
              </h3>
            </div>
            <div class="flex flex-col gap-4">
              <div
                class="flex flex-col gap-1 border-l-2 border-zinc-200 dark:border-zinc-700 pl-3 py-1"
              >
                <span
                  class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                  >上次启动</span
                >
                <span
                  class="text-sm font-medium font-mono text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.lastStartTime || '-' }}</span
                >
              </div>
              <div
                class="flex flex-col gap-1 border-l-2 border-zinc-200 dark:border-zinc-700 pl-3 py-1"
              >
                <span
                  class="text-[10px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500"
                  >上次离线</span
                >
                <span
                  class="text-sm font-medium font-mono text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.lastCloseTime || '-' }}</span
                >
              </div>
            </div>
          </div>

          <div
            class="lg:col-span-2 p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col"
          >
            <div class="flex items-center gap-2 mb-5">
              <span
                class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
              ></span>
              <h3
                class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
              >
                配置参数
              </h3>
            </div>

            <div class="grid grid-cols-2 sm:grid-cols-3 gap-x-4 gap-y-5">
              <div class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Lock /></el-icon> 加密</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.encryption ? '已启用' : '未启用' }}</span
                >
              </div>

              <div class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Lightning /></el-icon> 压缩</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200"
                  >{{ proxy.compression ? '已启用' : '未启用' }}</span
                >
              </div>

              <div
                v-if="proxy.customDomains"
                class="flex flex-col gap-1 col-span-2 sm:col-span-1"
              >
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Link /></el-icon> 自定义域名</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.customDomains }}</span
                >
              </div>

              <div v-if="proxy.subdomain" class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Link /></el-icon> 子域名</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.subdomain }}</span
                >
              </div>

              <div v-if="proxy.locations" class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Location /></el-icon> 位置</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.locations }}</span
                >
              </div>

              <div v-if="proxy.hostHeaderRewrite" class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Tickets /></el-icon> HOST 重定向</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.hostHeaderRewrite }}</span
                >
              </div>

              <div v-if="proxy.multiplexer" class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Cpu /></el-icon> 多路复用</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.multiplexer }}</span
                >
              </div>

              <div v-if="proxy.routeByHTTPUser" class="flex flex-col gap-1">
                <span
                  class="text-[10px] font-semibold text-zinc-400 dark:text-zinc-500 uppercase flex items-center gap-1.5"
                  ><el-icon><Connection /></el-icon> 用户 HTTP 路由</span
                >
                <span
                  class="text-[13px] font-medium text-zinc-800 dark:text-zinc-200 break-all"
                  >{{ proxy.routeByHTTPUser }}</span
                >
              </div>
            </div>

            <div
              v-if="proxy.annotations && proxy.annotations.size > 0"
              class="flex flex-wrap gap-2 mt-4 pt-4 border-t border-zinc-100 dark:border-zinc-800/80"
            >
              <span
                v-for="[key, value] in proxy.annotations"
                :key="key"
                class="px-2 py-0.5 bg-zinc-100 dark:bg-zinc-800/50 rounded text-[11px] font-medium text-zinc-600 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700"
              >
                <span class="text-zinc-400">{{ key }}:</span> {{ value }}
              </span>
            </div>
          </div>
        </div>

        <div
          class="p-5 sm:p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col gap-4"
        >
          <div
            class="flex items-center gap-2.5 pb-4 border-b border-zinc-100 dark:border-zinc-800/80"
          >
            <span
              class="w-1.5 h-4 bg-zinc-800 dark:bg-zinc-200 rounded-full"
            ></span>
            <h3
              class="text-[15px] font-bold text-zinc-800 dark:text-zinc-200 tracking-wide m-0"
            >
              流量监控数据
            </h3>
          </div>
          <div class="w-full">
            <Traffic :proxy-name="proxyName" />
          </div>
        </div>
      </template>

      <div
        v-else-if="!loading"
        class="flex flex-col items-center justify-center py-20 text-zinc-400 gap-4 bg-white/50 dark:bg-zinc-900/50 rounded-2xl border border-dashed border-zinc-200 dark:border-zinc-800 mt-5"
      >
        <el-icon class="text-4xl text-zinc-300 dark:text-zinc-600"
          ><Inbox
        /></el-icon>
        <div class="flex flex-col items-center gap-1">
          <span class="text-base font-bold text-zinc-700 dark:text-zinc-200"
            >隧道未找到</span
          >
          <span class="text-[13px]">隧道不存在或已经被移除</span>
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
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import {
  ArrowLeft,
  Monitor,
  Connection,
  Link,
  Lock,
  Promotion,
  Grid,
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

<style scoped></style>
