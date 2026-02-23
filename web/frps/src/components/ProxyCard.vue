<template>
  <router-link
    :to="proxyLink"
    class="group block w-full bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm hover:shadow-md hover:ring-indigo-500/40 dark:hover:ring-indigo-500/50 transition-all duration-300 text-zinc-800 dark:text-zinc-200 cursor-pointer overflow-hidden"
  >
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between p-4 sm:p-5 gap-4"
    >
      <div class="flex flex-col gap-3 min-w-0 flex-1">
        <div class="flex items-center gap-3">
          <span
            class="text-base font-bold text-zinc-900 dark:text-zinc-100 truncate tracking-tight group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors"
          >
            {{ proxy.name }}
          </span>
          <span
            v-if="showType"
            class="shrink-0 px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800 text-[10px] font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700/50"
          >
            {{ proxy.type.toUpperCase() }}
          </span>
        </div>

        <div class="flex flex-wrap items-center gap-x-5 gap-y-2.5">
          <div v-if="proxy.port" class="flex items-baseline gap-1.5">
            <span class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
              >端口:</span
            >
            <span
              class="text-[13px] font-semibold font-mono text-zinc-700 dark:text-zinc-300"
              >{{ proxy.port }}</span
            >
          </div>
          <div class="flex items-baseline gap-1.5">
            <span class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
              >连接数:</span
            >
            <span
              class="text-[13px] font-semibold font-mono text-zinc-700 dark:text-zinc-300"
              >{{ proxy.conns }}</span
            >
          </div>
          <div
            v-if="proxy.clientID"
            class="flex items-baseline gap-1.5 min-w-0"
          >
            <span
              class="text-xs font-medium shrink-0 text-zinc-400 dark:text-zinc-500"
              >客户端:</span
            >
            <span
              class="text-[13px] font-semibold text-zinc-700 dark:text-zinc-300 break-all"
              :title="
                proxy.user ? `${proxy.user}.${proxy.clientID}` : proxy.clientID
              "
            >
              {{
                proxy.user ? `${proxy.user}.${proxy.clientID}` : proxy.clientID
              }}
            </span>
          </div>
        </div>
      </div>

      <div
        class="flex flex-row items-center justify-between sm:justify-end gap-6 sm:gap-8 w-full sm:w-auto pt-4 sm:pt-0 border-t border-zinc-100 dark:border-zinc-800/80 sm:border-0 shrink-0"
      >
        <div
          class="flex sm:flex-col items-center sm:items-end gap-4 sm:gap-1.5"
        >
          <div class="flex items-center gap-1.5">
            <el-icon class="text-sm text-emerald-500 dark:text-emerald-400"
              ><Top
            /></el-icon>
            <span
              class="text-[11px] font-mono font-medium text-zinc-500 dark:text-zinc-400 w-16 sm:text-right"
            >
              {{ formatFileSize(proxy.trafficOut) }}
            </span>
          </div>
          <div class="flex items-center gap-1.5">
            <el-icon class="text-sm text-blue-500 dark:text-blue-400"
              ><Bottom
            /></el-icon>
            <span
              class="text-[11px] font-mono font-medium text-zinc-500 dark:text-zinc-400 w-16 sm:text-right"
            >
              {{ formatFileSize(proxy.trafficIn) }}
            </span>
          </div>
        </div>

        <div
          class="flex items-center justify-center px-3 py-1.5 rounded-lg text-[11px] font-bold tracking-widest uppercase transition-colors shrink-0"
          :class="
            proxy.status === 'online'
              ? 'bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
              : 'bg-rose-50 text-rose-600 border border-rose-200 dark:bg-rose-500/10 dark:text-rose-400 dark:border-rose-500/20'
          "
        >
          {{ proxy.status === 'online' ? '在线' : '离线' }}
        </div>
      </div>
    </div>
  </router-link>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { Top, Bottom } from '@element-plus/icons-vue'
import { formatFileSize } from '../utils/format'
import type { BaseProxy } from '../utils/proxy'

interface Props {
  proxy: BaseProxy
  showType?: boolean
}

const props = defineProps<Props>()
const route = useRoute()

const proxyLink = computed(() => {
  const base = `/proxy/${props.proxy.name}`
  // If we're on a client detail page, pass client info
  if (route.name === 'ClientDetail' && route.params.key) {
    return `${base}?from=client&client=${route.params.key}`
  }
  return base
})
</script>

<style scoped></style>
