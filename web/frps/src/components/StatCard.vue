<template>
  <div
    class="relative group p-6 bg-white dark:bg-zinc-900 border border-zinc-200 dark:border-zinc-800 rounded-2xl flex flex-col transition-colors duration-200"
    :class="{
      'cursor-pointer hover:border-zinc-300 dark:hover:border-zinc-700 hover:bg-zinc-50/80 dark:hover:bg-zinc-800/40':
        !!to,
    }"
    @click="handleClick"
  >
    <div class="flex justify-between items-start mb-4">
      <span class="text-sm font-semibold text-zinc-500 dark:text-zinc-400 mt-1">
        {{ label }}
      </span>

      <div
        class="w-10 h-10 rounded-xl flex items-center justify-center shrink-0 border"
        :class="{
          'bg-indigo-50/50 border-indigo-100 dark:bg-indigo-500/10 dark:border-indigo-500/20 text-indigo-600 dark:text-indigo-400':
            type === 'clients',
          'bg-rose-50/50 border-rose-100 dark:bg-rose-500/10 dark:border-rose-500/20 text-rose-600 dark:text-rose-400':
            type === 'proxies',
          'bg-blue-50/50 border-blue-100 dark:bg-blue-500/10 dark:border-blue-500/20 text-blue-600 dark:text-blue-400':
            type === 'connections',
          'bg-emerald-50/50 border-emerald-100 dark:bg-emerald-500/10 dark:border-emerald-500/20 text-emerald-600 dark:text-emerald-400':
            type === 'traffic',
        }"
      >
        <component :is="iconComponent" class="w-5 h-5" />
      </div>
    </div>

    <div class="mb-4">
      <span
        class="text-3xl sm:text-4xl font-black font-mono tracking-tight text-zinc-900 dark:text-zinc-50 leading-none"
      >
        {{ value }}
      </span>
    </div>

    <div
      class="flex items-center justify-between mt-auto pt-4 border-t border-zinc-100 dark:border-zinc-800/80"
    >
      <span
        v-if="subtitle"
        class="text-xs font-medium text-zinc-400 dark:text-zinc-500 truncate mr-2"
      >
        {{ subtitle }}
      </span>
      <div v-else class="flex-1"></div>

      <el-icon
        v-if="to"
        class="text-zinc-400 dark:text-zinc-500 text-base transition-transform duration-300 group-hover:translate-x-1 group-hover:text-zinc-800 dark:group-hover:text-zinc-200 shrink-0"
      >
        <ArrowRight />
      </el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import {
  User,
  Connection,
  DataAnalysis,
  Promotion,
  ArrowRight,
} from '@element-plus/icons-vue'

interface Props {
  label: string
  value: string | number
  type?: 'clients' | 'proxies' | 'connections' | 'traffic'
  subtitle?: string
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'clients',
})

const router = useRouter()

const iconComponent = computed(() => {
  switch (props.type) {
    case 'clients':
      return User
    case 'proxies':
      return Connection
    case 'connections':
      return DataAnalysis
    case 'traffic':
      return Promotion
    default:
      return User
  }
})

const handleClick = () => {
  if (props.to) {
    router.push(props.to)
  }
}
</script>

<style scoped></style>
