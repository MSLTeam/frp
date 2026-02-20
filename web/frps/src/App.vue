<template>
  <div class="app-shell">
    <div class="bg-orb orb-a"></div>
    <div class="bg-orb orb-b"></div>

    <header class="app-topbar">
      <div class="brand-group">
        <el-button
          class="mobile-nav-btn"
          circle
          plain
          @click="mobileNavVisible = true"
        >
          <el-icon><Menu /></el-icon>
        </el-button>

        <router-link to="/" class="brand-link">
          <div class="brand-badge">FRPS</div>
          <div class="brand-copy">
            <h1>MSLFrp 节点控制台</h1>
            <p>实时监控 · 流量治理 · 连接可视化</p>
          </div>
        </router-link>
      </div>

      <div class="topbar-actions">
        <div class="route-pill">
          <el-icon><DataBoard /></el-icon>
          <span>{{ currentRouteTitle }}</span>
        </div>

        <el-switch
          v-model="darkmodeSwitch"
          inline-prompt
          :active-icon="Moon"
          :inactive-icon="Sunny"
          class="theme-switch"
        />
      </div>
    </header>

    <main class="app-body">
      <aside class="sidebar-card">
        <el-scrollbar>
          <el-menu
            :default-active="$route.path"
            router
            class="side-nav"
            @select="handleSelect"
          >
            <el-menu-item index="/">
              <el-icon><HomeFilled /></el-icon>
              <span>总览看板</span>
            </el-menu-item>

            <el-sub-menu index="/proxies">
              <template #title>
                <el-icon><Connection /></el-icon>
                <span>隧道管理</span>
              </template>
              <el-menu-item
                v-for="item in proxyTypes"
                :key="item.value"
                :index="`/proxies/${item.value}`"
              >
                {{ item.label }}
              </el-menu-item>
            </el-sub-menu>

            <el-menu-item index="/user-center">
              <el-icon><User /></el-icon>
              <span>用户中心</span>
            </el-menu-item>
          </el-menu>
        </el-scrollbar>

        <div class="side-footnote">
          <el-icon><Guide /></el-icon>
          <span>监控刷新与数据统计由 FRPS API 实时驱动</span>
        </div>
      </aside>

      <section class="content-stage">
        <router-view v-slot="{ Component }">
          <transition name="panel-rise" mode="out-in">
            <component :is="Component" class="page-panel" />
          </transition>
        </router-view>
      </section>
    </main>

    <el-drawer
      v-model="mobileNavVisible"
      direction="ltr"
      size="260px"
      :with-header="false"
      class="mobile-drawer"
    >
      <div class="drawer-head">
        <div class="brand-badge">FRPS</div>
        <div class="drawer-title">MSLFrp 导航</div>
      </div>
      <el-menu
        :default-active="$route.path"
        router
        class="mobile-menu"
        @select="handleSelect"
      >
        <el-menu-item index="/">
          <el-icon><HomeFilled /></el-icon>
          <span>总览看板</span>
        </el-menu-item>

        <el-sub-menu index="/proxies">
          <template #title>
            <el-icon><Connection /></el-icon>
            <span>隧道管理</span>
          </template>
          <el-menu-item
            v-for="item in proxyTypes"
            :key="item.value"
            :index="`/proxies/${item.value}`"
          >
            {{ item.label }}
          </el-menu-item>
        </el-sub-menu>

        <el-menu-item index="/user-center">
          <el-icon><User /></el-icon>
          <span>用户中心</span>
        </el-menu-item>
      </el-menu>
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { useRoute } from 'vue-router'
import {
  HomeFilled,
  Connection,
  User,
  Moon,
  Sunny,
  Menu,
  DataBoard,
  Guide,
} from '@element-plus/icons-vue'

const route = useRoute()
const mobileNavVisible = ref(false)

const proxyTypes = [
  { value: 'tcp', label: 'TCP' },
  { value: 'udp', label: 'UDP' },
  { value: 'http', label: 'HTTP' },
  { value: 'https', label: 'HTTPS' },
  { value: 'tcpmux', label: 'TCPMUX' },
  { value: 'stcp', label: 'STCP' },
  { value: 'sudp', label: 'SUDP' },
]

const routeTitleMap: Record<string, string> = {
  '/': '总览看板',
  '/proxies/tcp': '隧道管理 / TCP',
  '/proxies/udp': '隧道管理 / UDP',
  '/proxies/http': '隧道管理 / HTTP',
  '/proxies/https': '隧道管理 / HTTPS',
  '/proxies/tcpmux': '隧道管理 / TCPMUX',
  '/proxies/stcp': '隧道管理 / STCP',
  '/proxies/sudp': '隧道管理 / SUDP',
}

const currentRouteTitle = computed(() => {
  return routeTitleMap[route.path] || 'MSLFrp 控制台'
})

const isDark = useDark()
const toggleDark = useToggle(isDark)
const darkmodeSwitch = computed({
  get: () => isDark.value,
  set: (val: boolean) => toggleDark(val),
})

const handleSelect = (key: string) => {
  if (key === '/user-center') {
    window.open('https://user.mslmc.net', '_blank')
  }
  mobileNavVisible.value = false
}
</script>

<style lang="scss" scoped>
.app-shell {
  --accent: #0ea5e9;
  --accent-strong: #0284c7;
  --accent-soft: #e0f2fe;
  --accent-alt: #14b8a6;
  --warm: #f59e0b;

  min-height: 100vh;
  position: relative;
  overflow: hidden;
  background:
    radial-gradient(circle at 10% -10%, rgba(20, 184, 166, 0.18), transparent 35%),
    radial-gradient(circle at 95% 0%, rgba(14, 165, 233, 0.2), transparent 32%),
    linear-gradient(160deg, #f2f7ff 0%, #f4fbf8 55%, #f8fafc 100%);
  color: var(--el-text-color-primary);
  font-family: 'Space Grotesk', 'Noto Sans SC', 'Source Han Sans SC',
    'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;
  padding: 18px;
  box-sizing: border-box;
}

:global(html.dark) .app-shell {
  background:
    radial-gradient(circle at 8% -18%, rgba(20, 184, 166, 0.22), transparent 34%),
    radial-gradient(circle at 100% -5%, rgba(14, 165, 233, 0.2), transparent 32%),
    linear-gradient(180deg, #0b1220 0%, #0f172a 55%, #111827 100%);
}

.bg-orb {
  position: absolute;
  pointer-events: none;
  border-radius: 999px;
  filter: blur(12px);
}

.orb-a {
  width: 220px;
  height: 220px;
  right: -60px;
  top: 80px;
  background: rgba(14, 165, 233, 0.12);
}

.orb-b {
  width: 180px;
  height: 180px;
  left: -40px;
  bottom: -20px;
  background: rgba(20, 184, 166, 0.15);
}

.app-topbar {
  position: relative;
  z-index: 2;
  height: 82px;
  border-radius: 22px;
  padding: 0 22px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: rgba(255, 255, 255, 0.84);
  border: 1px solid rgba(255, 255, 255, 0.65);
  box-shadow:
    0 10px 40px rgba(14, 165, 233, 0.12),
    inset 0 1px 0 rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(8px);
}

:global(html.dark) .app-topbar {
  background: rgba(15, 23, 42, 0.76);
  border-color: rgba(148, 163, 184, 0.16);
  box-shadow:
    0 10px 38px rgba(2, 6, 23, 0.45),
    inset 0 1px 0 rgba(148, 163, 184, 0.08);
}

.brand-group {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
}

.mobile-nav-btn {
  display: none;
}

.brand-link {
  display: flex;
  align-items: center;
  text-decoration: none;
  gap: 12px;
  min-width: 0;
}

.brand-badge {
  width: 46px;
  height: 46px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  font-size: 13px;
  font-weight: 700;
  letter-spacing: 0.7px;
  color: #fff;
  background: linear-gradient(145deg, var(--accent) 0%, var(--accent-alt) 100%);
  box-shadow: 0 8px 20px rgba(14, 165, 233, 0.35);
}

.brand-copy {
  min-width: 0;

  h1 {
    margin: 0;
    font-size: 21px;
    line-height: 1.1;
    color: #0f172a;
    white-space: nowrap;
  }

  p {
    margin: 4px 0 0;
    font-size: 12px;
    letter-spacing: 0.45px;
    color: #475569;
    white-space: nowrap;
  }
}

:global(html.dark) .brand-copy {
  h1 {
    color: #e5e7eb;
  }

  p {
    color: #94a3b8;
  }
}

.topbar-actions {
  display: flex;
  align-items: center;
  gap: 14px;
}

.route-pill {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  font-weight: 600;
  color: #0f766e;
  background: rgba(20, 184, 166, 0.12);
  border: 1px solid rgba(20, 184, 166, 0.22);
  border-radius: 999px;
  padding: 8px 12px;
}

:global(html.dark) .route-pill {
  color: #5eead4;
  background: rgba(20, 184, 166, 0.15);
  border-color: rgba(45, 212, 191, 0.26);
}

.theme-switch {
  --el-switch-on-color: var(--accent-strong);
  --el-switch-off-color: #94a3b8;
}

.app-body {
  position: relative;
  z-index: 2;
  margin-top: 16px;
  height: calc(100vh - 134px);
  display: grid;
  grid-template-columns: 280px 1fr;
  gap: 16px;
}

.sidebar-card {
  background: rgba(255, 255, 255, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.7);
  border-radius: 20px;
  box-shadow: 0 8px 30px rgba(15, 23, 42, 0.08);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

:global(html.dark) .sidebar-card {
  background: rgba(15, 23, 42, 0.72);
  border-color: rgba(148, 163, 184, 0.14);
  box-shadow: 0 8px 28px rgba(2, 6, 23, 0.35);
}

.side-nav {
  border-right: none;
  background: transparent;
  padding: 14px 10px;

  :deep(.el-menu-item),
  :deep(.el-sub-menu__title) {
    margin: 6px 6px;
    border-radius: 12px;
    height: 44px;
    line-height: 44px;
    transition: all 0.2s ease;

    .el-icon {
      margin-right: 10px;
      font-size: 16px;
    }
  }

  :deep(.el-menu-item:hover),
  :deep(.el-sub-menu__title:hover) {
    background: rgba(14, 165, 233, 0.1);
  }

  :deep(.el-menu-item.is-active) {
    color: #fff;
    background: linear-gradient(135deg, var(--accent) 0%, var(--accent-alt) 100%);
    box-shadow: 0 8px 18px rgba(14, 165, 233, 0.28);
  }
}

.side-footnote {
  margin: 8px 14px 14px;
  border-radius: 12px;
  padding: 10px 12px;
  font-size: 12px;
  line-height: 1.45;
  color: #3f3f46;
  background: linear-gradient(145deg, rgba(245, 158, 11, 0.12), rgba(14, 165, 233, 0.1));
  border: 1px solid rgba(245, 158, 11, 0.2);
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

:global(html.dark) .side-footnote {
  color: #f8fafc;
  border-color: rgba(245, 158, 11, 0.25);
  background: linear-gradient(145deg, rgba(245, 158, 11, 0.16), rgba(14, 165, 233, 0.15));
}

.content-stage {
  min-width: 0;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.55);
  border: 1px solid rgba(255, 255, 255, 0.66);
  box-shadow: 0 8px 30px rgba(15, 23, 42, 0.08);
  backdrop-filter: blur(6px);
  padding: 16px;
  overflow: auto;
}

:global(html.dark) .content-stage {
  background: rgba(15, 23, 42, 0.58);
  border-color: rgba(148, 163, 184, 0.14);
  box-shadow: 0 8px 28px rgba(2, 6, 23, 0.35);
}

.page-panel {
  min-height: 100%;
}

:deep(.mobile-drawer .el-drawer) {
  background: linear-gradient(180deg, #f8fbff 0%, #f5fdf9 100%);
}

.drawer-head {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 8px 10px 14px;
}

.drawer-title {
  font-size: 16px;
  font-weight: 700;
}

.mobile-menu {
  border-right: none;
  background: transparent;
}

.panel-rise-enter-active,
.panel-rise-leave-active {
  transition: all 0.25s ease;
}

.panel-rise-enter-from {
  opacity: 0;
  transform: translateY(10px);
}

.panel-rise-leave-to {
  opacity: 0;
  transform: translateY(-10px);
}

@media (max-width: 1024px) {
  .app-shell {
    padding: 12px;
  }

  .app-topbar {
    height: 74px;
    border-radius: 16px;
    padding: 0 14px;
  }

  .mobile-nav-btn {
    display: inline-flex;
  }

  .brand-copy p,
  .route-pill {
    display: none;
  }

  .brand-copy h1 {
    font-size: 18px;
  }

  .app-body {
    grid-template-columns: 1fr;
    height: calc(100vh - 110px);
  }

  .sidebar-card {
    display: none;
  }

  .content-stage {
    padding: 12px;
    border-radius: 16px;
  }
}

@media (max-width: 640px) {
  .brand-badge {
    width: 40px;
    height: 40px;
    font-size: 11px;
  }

  .brand-copy h1 {
    font-size: 16px;
  }

  .topbar-actions {
    gap: 8px;
  }
}
</style>
