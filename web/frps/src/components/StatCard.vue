<template>
  <div
    class="relative group p-6 bg-white dark:bg-zinc-900 rounded-3xl ring-1 ring-zinc-200 dark:ring-zinc-800 shadow-sm transition-all duration-500 overflow-hidden flex flex-col"
    :class="{
      'cursor-pointer hover:-translate-y-1 hover:shadow-xl hover:shadow-zinc-200/50 dark:hover:shadow-none':
        !!to,
    }"
    @click="handleClick"
  >
    <div
      v-if="!!to"
      class="absolute -right-10 -top-10 w-40 h-40 rounded-full blur-3xl transition-all opacity-0 group-hover:opacity-100 pointer-events-none"
      :class="{
        'bg-indigo-500/10': type === 'clients',
        'bg-rose-500/10': type === 'proxies',
        'bg-blue-500/10': type === 'connections',
        'bg-emerald-500/10': type === 'traffic',
      }"
    ></div>

    <div class="relative flex items-center justify-between z-10">
      <div class="flex items-center gap-4">
        <div
          class="w-14 h-14 rounded-2xl flex items-center justify-center text-white shrink-0 shadow-inner"
          :class="{
            'bg-gradient-to-br from-indigo-500 to-purple-600 dark:from-indigo-400 dark:to-purple-500':
              type === 'clients',
            'bg-gradient-to-br from-rose-400 to-red-500 dark:from-rose-500 dark:to-red-600':
              type === 'proxies',
            'bg-gradient-to-br from-blue-400 to-cyan-500 dark:from-blue-500 dark:to-cyan-600':
              type === 'connections',
            'bg-gradient-to-br from-emerald-400 to-teal-500 dark:from-emerald-500 dark:to-teal-600':
              type === 'traffic',
          }"
        >
          <component :is="iconComponent" class="w-7 h-7" />
        </div>

        <div class="flex flex-col">
          <span
            class="text-3xl font-black tracking-tight text-zinc-900 dark:text-white leading-none mb-1.5"
            >{{ value }}</span
          >
          <span
            class="text-[11px] font-bold uppercase tracking-widest text-zinc-400 dark:text-zinc-500 leading-none"
            >{{ label }}</span
          >
        </div>
      </div>

      <el-icon
        v-if="to"
        class="text-zinc-400 dark:text-zinc-500 text-xl transition-transform duration-300 group-hover:translate-x-1.5 group-hover:text-zinc-900 dark:group-hover:text-white shrink-0"
      >
        <ArrowRight />
      </el-icon>
    </div>

    <div
      v-if="subtitle"
      class="relative z-10 mt-5 pt-3 border-t border-zinc-100 dark:border-zinc-800/80 text-xs font-medium text-zinc-500 dark:text-zinc-400"
    >
      {{ subtitle }}
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

<style scoped>

</style>
