<template>
  <div
    id="app"
    class="min-h-screen flex flex-col bg-zinc-50 dark:bg-black font-sans text-zinc-900 dark:text-zinc-100 transition-colors duration-300"
  >
    <header
      class="sticky top-0 z-[100] w-full bg-white/70 dark:bg-zinc-950/70 backdrop-blur-2xl border-b border-zinc-200/80 dark:border-zinc-800/80 shadow-sm transition-colors duration-300"
    >
      <div class="w-full px-4 sm:px-6 md:px-8">
        <div class="h-16 flex items-center justify-between gap-4">
          <div
            class="flex items-center gap-3 shrink-0 group cursor-pointer hover:opacity-80 transition-opacity"
          >
            <div
              class="flex items-center justify-center w-10 h-10 rounded-2xl bg-gradient-to-br from-zinc-100 to-white dark:from-zinc-800 dark:to-zinc-900 ring-1 ring-zinc-200/80 dark:ring-zinc-700/80 shadow-sm group-hover:ring-indigo-500/50 transition-all"
            >
              <img
                src="./assets/icons/msl-user.png"
                class="w-6 h-6 object-contain drop-shadow-sm"
                alt="MSL Logo"
              />
            </div>
            <div class="flex flex-col justify-center">
              <div class="flex items-center gap-2">
                <span
                  class="font-black text-[17px] text-zinc-900 dark:text-white tracking-tight leading-none"
                  >MSLFrp</span
                >
                <span
                  class="hidden sm:inline-block text-[9px] px-2 py-0.5 rounded-md bg-zinc-900 dark:bg-zinc-100 text-white dark:text-zinc-900 font-black tracking-widest uppercase leading-none shadow-sm"
                >
                  节点控制台
                </span>
              </div>
              <span
                class="text-[11px] font-semibold text-zinc-400 dark:text-zinc-500 leading-none mt-1.5 uppercase tracking-wider"
              >
                {{ currentRouteName || 'System' }}
              </span>
            </div>
          </div>

          <nav
            class="hidden md:flex items-center p-1 bg-zinc-100/80 dark:bg-zinc-800/80 rounded-xl border border-zinc-200/50 dark:border-zinc-700/50 shadow-inner"
          >
            <router-link
              to="/"
              class="px-6 py-1.5 text-[13px] font-bold text-zinc-500 dark:text-zinc-400 rounded-lg transition-all duration-300 select-none"
              active-class="!text-zinc-900 dark:!text-white bg-white dark:bg-zinc-700 shadow-sm ring-1 ring-zinc-200 dark:ring-zinc-600"
            >
              仪表盘
            </router-link>

            <router-link
              to="/clients"
              class="px-6 py-1.5 text-[13px] font-bold text-zinc-500 dark:text-zinc-400 rounded-lg transition-all duration-300 select-none"
              active-class="!text-zinc-900 dark:!text-white bg-white dark:bg-zinc-700 shadow-sm ring-1 ring-zinc-200 dark:ring-zinc-600"
            >
              客户端
            </router-link>

            <router-link
              to="/proxies"
              class="px-6 py-1.5 text-[13px] font-bold text-zinc-500 dark:text-zinc-400 rounded-lg transition-all duration-300 select-none"
              :class="{
                '!text-zinc-900 dark:!text-white bg-white dark:bg-zinc-700 shadow-sm ring-1 ring-zinc-200 dark:ring-zinc-600':
                  route.path.startsWith('/proxies'),
              }"
            >
              隧道
            </router-link>
          </nav>

          <div class="flex items-center gap-2.5 sm:gap-3 shrink-0">
            <a
              class="w-9 h-9 flex items-center justify-center rounded-xl bg-zinc-50 dark:bg-zinc-800/50 text-zinc-500 dark:text-zinc-400 hover:text-zinc-900 dark:hover:text-white hover:bg-zinc-100 dark:hover:bg-zinc-800 border border-zinc-200/80 dark:border-zinc-700/80 transition-all shadow-sm hover:shadow"
              href="https://github.com/MSLTeam/frp"
              target="_blank"
              aria-label="GitHub"
            >
              <GitHubIcon class="w-[18px] h-[18px] fill-current" />
            </a>

            <div
              class="h-4 w-px bg-zinc-200 dark:bg-zinc-700 hidden sm:block mx-1"
            ></div>

            <div
              class="flex items-center justify-center h-9 px-2.5 rounded-xl bg-zinc-50 dark:bg-zinc-800/50 border border-zinc-200/80 dark:border-zinc-700/80 shadow-sm"
            >
              <el-switch
                :model-value="isDark"
                @change="toggleTheme"
                inline-prompt
                :active-icon="Moon"
                :inactive-icon="Sunny"
                class="theme-switch !h-5"
              />
            </div>
          </div>
        </div>

        <nav
          class="md:hidden flex items-center gap-2 overflow-x-auto custom-scrollbar py-2.5 border-t border-zinc-100 dark:border-zinc-800/80"
        >
          <router-link
            to="/"
            class="px-4 py-2 text-xs font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 bg-zinc-100/50 dark:bg-zinc-800/50 border border-zinc-200/50 dark:border-zinc-700/50 rounded-lg whitespace-nowrap transition-all"
            active-class="!text-zinc-900 dark:!text-white !bg-white dark:!bg-zinc-800 !border-zinc-300 dark:!border-zinc-600 shadow-sm"
          >
            仪表盘
          </router-link>
          <router-link
            to="/clients"
            class="px-4 py-2 text-xs font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 bg-zinc-100/50 dark:bg-zinc-800/50 border border-zinc-200/50 dark:border-zinc-700/50 rounded-lg whitespace-nowrap transition-all"
            active-class="!text-zinc-900 dark:!text-white !bg-white dark:!bg-zinc-800 !border-zinc-300 dark:!border-zinc-600 shadow-sm"
          >
            客户端
          </router-link>
          <router-link
            to="/proxies"
            class="px-4 py-2 text-xs font-bold uppercase tracking-wider text-zinc-500 dark:text-zinc-400 bg-zinc-100/50 dark:bg-zinc-800/50 border border-zinc-200/50 dark:border-zinc-700/50 rounded-lg whitespace-nowrap transition-all"
            :class="{
              '!text-zinc-900 dark:!text-white !bg-white dark:!bg-zinc-800 !border-zinc-300 dark:!border-zinc-600 shadow-sm':
                route.path.startsWith('/proxies'),
            }"
          >
            隧道
          </router-link>
        </nav>
      </aside>

    <main
      id="content"
      class="flex-1 w-full px-4 sm:px-6 md:px-8 py-6 sm:py-8 box-border"
    >
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick } from 'vue'
import { useRoute } from 'vue-router'
import { useDark } from '@vueuse/core'
import { Moon, Sunny } from '@element-plus/icons-vue'
import GitHubIcon from './assets/icons/github.svg?component'

const route = useRoute()
const isDark = useDark()

const toggleTheme = (val: boolean | string | number) => {
  const isDarkVal = val as boolean

  if (!document.startViewTransition) {
    isDark.value = isDarkVal
    return
  }

  document.startViewTransition(async () => {
    isDark.value = isDarkVal
    await nextTick()
  })
}

const currentRouteName = computed(() => {
  if (route.path === '/') return '仪表盘'
  if (route.path.startsWith('/clients')) return '客户端'
  if (route.path.startsWith('/proxies')) return '隧道数'
  return ''
})
</script>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  display: none;
}
.custom-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

/* 隐藏横向滚动的滚动条但保留滚动功能 */
.custom-scrollbar::-webkit-scrollbar {
  display: none;
}
.custom-scrollbar {
  -ms-overflow-style: none;
  scrollbar-width: none;
}


.page-fade-enter-active,
.page-fade-leave-active {
  transition:
    opacity 0.25s ease,
    transform 0.3s cubic-bezier(0.2, 0.8, 0.2, 1),
    filter 0.25s ease;
}

/* 新页面进场前的初始状态：轻微向下偏移、缩小、模糊、透明 */
.page-fade-enter-from {
  opacity: 0;
  transform: translateY(10px) scale(0.99);
  filter: blur(4px);
}

/* 老页面退场后的最终状态：轻微向上偏移、缩小、模糊、透明 */
.page-fade-leave-to {
  opacity: 0;
  transform: translateY(-10px) scale(0.99);
  filter: blur(4px);
}
</style>
