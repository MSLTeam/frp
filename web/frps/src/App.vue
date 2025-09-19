<template>
  <div class="app-container">
    <!-- 头部 -->
    <header class="app-header">
      <div class="header-content">
        <router-link to="/" class="brand">
          <img src="/favicon.ico" alt="Logo" class="logo" />
          <span class="title">MSLFrp 节点监控</span>
        </router-link>
        <el-switch
          v-model="darkmodeSwitch"
          inline-prompt
          :active-icon="Moon"
          :inactive-icon="Sunny"
          @change="toggleDark"
          class="theme-switch"
        />
      </div>
    </header>

    <!-- 主内容 -->
    <main class="app-main">
      <el-row class="main-wrapper" :gutter="0">
        <!-- 侧边导航 -->
        <el-col :xs="24" :md="4" class="nav-col">
          <el-scrollbar>
            <el-menu
              :default-active="$route.path"
              router
              class="side-nav"
              :collapse="isCollapse"
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
        </el-col>

        <!-- 内容区 -->
        <el-col :xs="24" :md="20" class="content-col">
          <el-scrollbar class="content-scroll">
            <router-view v-slot="{ Component }">
              <transition name="fade-slide" mode="out-in">
                <component :is="Component" class="page-content" />
              </transition>
            </router-view>
          </el-scrollbar>
        </el-col>
      </el-row>
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import {
  HomeFilled,
  Connection,
  User,
  Moon,
  Sunny,
} from '@element-plus/icons-vue'

// 响应式布局
const isCollapse = ref(false)
const proxyTypes = [
  { value: 'tcp', label: 'TCP' },
  { value: 'udp', label: 'UDP' },
  { value: 'http', label: 'HTTP' },
  { value: 'https', label: 'HTTPS' },
  { value: 'tcpmux', label: 'TCPMUX' },
  { value: 'stcp', label: 'STCP' },
  { value: 'sudp', label: 'SUDP' },
]

// 暗黑模式
const isDark = useDark()
const darkmodeSwitch = ref(isDark)
const toggleDark = useToggle(isDark)

const handleSelect = (key: string) => {
  if (key === '/user-center') {
    window.open('https://user.mslmc.net', '_blank')
  }
}
</script>

<style lang="scss" scoped>

.app-container {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--el-bg-color-page);
  margin: 0;
  overflow: hidden;
}

.app-header {
  height: 64px;
  background: linear-gradient(135deg, #6366f1 0%, #8b5cf6 100%);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;

  .dark & {
    background: linear-gradient(135deg, #4f46e5 0%, #7c3aed 100%);
  }

  .header-content {
    height: 100%;
    margin-right: 10px;
    padding: 0 24px;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .brand {
    display: flex;
    align-items: center;
    gap: 12px;
    text-decoration: none;

    .logo {
      height: 36px;
      border-radius: 6px;
    }

    .title {
      color: #fff;
      font-size: 20px;
      font-weight: 600;
      letter-spacing: 0.5px;
    }
  }

  .theme-switch {
    --el-switch-on-color: #475569;
    --el-switch-off-color: #e2e8f0;
  }
}

.app-main {
  flex: 1;
  overflow: hidden;

  .main-wrapper {
    height: 100%;
    margin: 0 !important;
  }
}

.nav-col {
  background: var(--el-bg-color-overlay);
  border-right: 1px solid var(--el-border-color-light);

  .side-nav {
    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      justify-content: flex-start;
      text-align: left;
      padding-left: 20px !important;

      .el-icon {
        margin-right: 12px;
      }
      }
    border-right: none;
    transition: width 0.2s;

    :deep(.el-menu-item),
    :deep(.el-sub-menu__title) {
      height: 48px;
      margin: 4px 12px;
      border-radius: 8px;
      transition: all 0.2s;

      &:hover {
        background: var(--el-color-primary-light-9);
      }

      &.is-active {
        background: var(--el-color-primary-light-8);
        color: var(--el-color-primary);
      }
    }
  }
}

.content-col {
  .content-scroll {
    height: calc(100vh - 64px);
    padding: 24px;
  }

  .page-content {
    border-radius: 12px;
    padding: 0px;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    min-height: calc(100% - 48px);
  }
}

// 过渡动画
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

@media (max-width: 768px) {
  .nav-col {
    display: none;
  }

  .content-col .content-scroll {
    padding: 16px;
  }

  .page-content {
    padding: 16px !important;
  }
}
</style>
