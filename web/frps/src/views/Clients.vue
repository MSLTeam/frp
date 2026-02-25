<template>
  <div
    class="mx-auto flex flex-col gap-5 text-zinc-800 dark:text-zinc-200 h-full pb-10"
  >
    <div
      class="flex flex-col xl:flex-row xl:items-center justify-between gap-4 p-5 bg-white dark:bg-zinc-900 rounded-2xl shadow-sm border border-zinc-200 dark:border-zinc-800/80"
    >
      <div class="flex flex-col gap-1 shrink-0">
        <h3
          class="text-lg font-bold tracking-tight m-0 flex items-center gap-2"
        >
          <el-icon class="text-indigo-500"><Monitor /></el-icon> 客户端状态
        </h3>
      </div>

      <div
        class="flex flex-col sm:flex-row sm:items-center sm:justify-end gap-3 w-full xl:w-auto xl:ml-auto"
      >
        <div class="w-full sm:w-64 lg:w-72 shrink-0">
          <el-input
            v-model="searchText"
            placeholder="搜索客户端..."
            clearable
            :prefix-icon="Search"
          />
        </div>
      </div>
    </div>

    <div class="flex items-center gap-2 overflow-x-auto custom-scrollbar pb-1">
      <button
        v-for="tab in statusTabs"
        :key="tab.value"
        class="flex items-center gap-2 px-4 py-1.5 rounded-lg text-xs font-bold uppercase tracking-wider transition-colors whitespace-nowrap shrink-0 border-0 cursor-pointer"
        :class="
          statusFilter === tab.value
            ? 'bg-indigo-500 text-white shadow-sm'
            : 'bg-white dark:bg-zinc-900 text-zinc-500 dark:text-zinc-400 hover:bg-zinc-100 dark:hover:bg-zinc-800 border border-zinc-200 dark:border-zinc-800'
        "
        @click="statusFilter = tab.value"
      >
        <span
          class="w-2 h-2 rounded-full"
          :class="{
            'bg-zinc-200 dark:bg-zinc-600':
              tab.value === 'all' && statusFilter !== 'all',
            'bg-zinc-100': tab.value === 'all' && statusFilter === 'all',
            'bg-emerald-500':
              tab.value === 'online' && statusFilter !== 'online',
            'bg-emerald-300':
              tab.value === 'online' && statusFilter === 'online',
            'bg-zinc-400':
              tab.value === 'offline' && statusFilter !== 'offline',
            'bg-zinc-300':
              tab.value === 'offline' && statusFilter === 'offline',
          }"
        ></span>
        {{ tab.label }}
        <span class="opacity-80 ml-0.5">{{ tab.count }}</span>
      </button>
    </div>

    <div
      v-loading="loading"
      element-loading-background="rgba(0, 0, 0, 0.0)"
      class="min-h-[200px]"
    >
      <div v-if="filteredClients.length > 0" class="flex flex-col gap-4">
        <ClientCard
          v-for="client in filteredClients"
          :key="client.key"
          :client="client"
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
            ><Monitor
          /></el-icon>
        </div>
        <span class="text-sm font-medium">未找到任何客户端</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Monitor } from '@element-plus/icons-vue'
import { Client } from '../utils/client'
import ClientCard from '../components/ClientCard.vue'
import { getClients } from '../api/client'

const clients = ref<Client[]>([])
const loading = ref(false)
const searchText = ref('')
const statusFilter = ref<'all' | 'online' | 'offline'>('all')

let refreshTimer: number | null = null

const stats = computed(() => {
  const total = clients.value.length
  const online = clients.value.filter((c) => c.online).length
  const offline = total - online
  return { total, online, offline }
})

const statusTabs = computed(() => [
  { value: 'all' as const, label: '全部', count: stats.value.total },
  { value: 'online' as const, label: '在线', count: stats.value.online },
  { value: 'offline' as const, label: '离线', count: stats.value.offline },
])

const filteredClients = computed(() => {
  let result = clients.value

  // Filter by status
  if (statusFilter.value === 'online') {
    result = result.filter((c) => c.online)
  } else if (statusFilter.value === 'offline') {
    result = result.filter((c) => !c.online)
  }

  // Filter by search text
  if (searchText.value) {
    result = result.filter((c) => c.matchesFilter(searchText.value))
  }

  // Sort: online first, then by display name
  result.sort((a, b) => {
    if (a.online !== b.online) {
      return a.online ? -1 : 1
    }
    return a.displayName.localeCompare(b.displayName)
  })

  return result
})

const fetchData = async (showLoading = true) => {
  if (showLoading) {
    loading.value = true
  }

  try {
    const json = await getClients()
    clients.value = json.map((data) => new Client(data))
  } catch (error: any) {
    ElMessage({
      showClose: true,
      message: 'Failed to fetch clients: ' + error.message,
      type: 'error',
    })
  } finally {
    if (showLoading) {
      loading.value = false
    }
  }
}

const startAutoRefresh = () => {
  refreshTimer = window.setInterval(() => {
    fetchData(false)
  }, 5000)
}

const stopAutoRefresh = () => {
  if (refreshTimer !== null) {
    window.clearInterval(refreshTimer)
    refreshTimer = null
  }
}

onMounted(() => {
  fetchData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
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
