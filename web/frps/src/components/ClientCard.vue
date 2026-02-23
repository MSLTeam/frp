<template>
  <div
    class="group block w-full bg-white dark:bg-zinc-900 rounded-2xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm hover:shadow-md hover:ring-indigo-500/40 dark:hover:ring-indigo-500/50 transition-all duration-300 text-zinc-800 dark:text-zinc-200 cursor-pointer overflow-hidden"
    @click="viewDetail"
  >
    <div
      class="flex flex-col sm:flex-row sm:items-center justify-between p-4 sm:p-5 gap-4"
    >
      <div class="flex items-start sm:items-center gap-4 min-w-0 flex-1">
        <div
          class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 transition-colors duration-300"
          :class="
            client.online
              ? 'bg-emerald-50 dark:bg-emerald-500/10 group-hover:bg-emerald-100 dark:group-hover:bg-emerald-500/20'
              : 'bg-zinc-100 dark:bg-zinc-800 group-hover:bg-zinc-200 dark:group-hover:bg-zinc-700'
          "
        >
          <span
            class="w-2.5 h-2.5 rounded-full ring-2 transition-all duration-300"
            :class="
              client.online
                ? 'bg-emerald-500 ring-emerald-200 dark:ring-emerald-500/30'
                : 'bg-zinc-400 dark:bg-zinc-500 ring-transparent'
            "
          ></span>
        </div>

        <div class="flex flex-col gap-2.5 min-w-0 flex-1">
          <div class="flex items-center gap-3">
            <span
              class="text-base font-bold text-zinc-900 dark:text-zinc-100 break-all leading-tight group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors"
            >
              {{ formatSafeText(client.displayName) }}
            </span>
            <span
              v-if="client.hostname"
              class="shrink-0 px-2 py-0.5 rounded-md bg-zinc-100 dark:bg-zinc-800 text-[10px] font-bold tracking-wider text-zinc-500 dark:text-zinc-400 border border-zinc-200 dark:border-zinc-700/50 truncate max-w-[120px]"
              :title="client.hostname"
            >
              {{ client.hostname }}
            </span>
          </div>

          <div class="flex flex-wrap items-center gap-x-5 gap-y-2">
            <div v-if="client.ip" class="flex items-baseline gap-1.5">
              <span class="text-xs font-medium text-zinc-400 dark:text-zinc-500"
                >IP:</span
              >
              <span
                class="text-[13px] font-semibold font-mono text-zinc-700 dark:text-zinc-300"
                >{{ client.ip }}</span
              >
            </div>
            <div class="flex items-center gap-1.5">
              <el-icon class="text-xs text-zinc-400 dark:text-zinc-500"
                ><DataLine
              /></el-icon>
              <span
                class="text-[12px] font-medium text-zinc-500 dark:text-zinc-400"
              >
                {{
                  client.online
                    ? client.lastConnectedAgo
                    : client.disconnectedAgo
                }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div
        class="flex flex-row items-center justify-between sm:justify-end gap-6 sm:gap-4 w-full sm:w-auto pt-4 sm:pt-0 border-t border-zinc-100 dark:border-zinc-800/80 sm:border-0 shrink-0"
      >
        <div
          class="flex items-center justify-center px-3 py-1.5 rounded-lg text-[11px] font-bold tracking-widest uppercase transition-colors shrink-0"
          :class="
            client.online
              ? 'bg-emerald-50 text-emerald-600 border border-emerald-200 dark:bg-emerald-500/10 dark:text-emerald-400 dark:border-emerald-500/20'
              : 'bg-zinc-50 text-zinc-500 border border-zinc-200 dark:bg-zinc-800/50 dark:text-zinc-400 dark:border-zinc-700/50'
          "
        >
          {{ client.online ? '在线' : '离线' }}
        </div>

        <el-icon
          class="text-lg text-zinc-400 dark:text-zinc-500 group-hover:text-zinc-900 dark:group-hover:text-zinc-100 group-hover:translate-x-1 transition-all duration-300"
        >
          <ArrowRight />
        </el-icon>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { DataLine, ArrowRight } from '@element-plus/icons-vue'
import type { Client } from '../utils/client'

interface Props {
  client: Client
}

const props = defineProps<Props>()
const router = useRouter()

const formatSafeText = (text: string | undefined) => {
  if (!text) return ''

  const matchWithSuffix = text.match(/^.*-(\d+)\.(.+)$/)
  if (matchWithSuffix) {
    const uid = parseInt(matchWithSuffix[1]) - 10000
    return `UID:${uid} · ${matchWithSuffix[2]}`
  }

  const matchTokenOnly = text.match(/^.*-(\d+)$/)
  if (matchTokenOnly) {
    const uid = parseInt(matchTokenOnly[1]) - 10000
    return `UID:${uid}`
  }

  return text
}

const viewDetail = () => {
  router.push({
    name: 'ClientDetail',
    params: { key: props.client.key },
  })
}
</script>

<style scoped></style>
