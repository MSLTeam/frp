<template>
  <div id="app">
    <header class="grid-content header-color">
      <div class="header-content">
        <div class="brand">
          <img src="/favicon.ico" alt="Logo" class="logo" />
          <a href="#">MSL-Frp节点控制面板</a>
        </div>
        <div class="dark-switch">
          <el-switch
            v-model="darkmodeSwitch"
            inline-prompt
            active-text="🌙"
            inactive-text="☀"
            @change="toggleDark"
            style="
              --el-switch-on-color: #444452;
              --el-switch-off-color: #589ef8;
            "
          />
        </div>
      </div>
    </header>
    <section>
      <el-row>
        <el-col id="side-nav" :xs="24" :md="3">
          <el-menu
            default-active="/"
            mode="vertical"
            theme="light"
            class="el-menu--vertical"
            router="false"
            @select="handleSelect"
            style="height: calc(100vh - 60px);"
          >
            <el-menu-item index="/">
              <el-icon><home-filled /></el-icon>
              <span>总览</span>
            </el-menu-item>
            <el-sub-menu index="/proxies">
              <template #title>
                <el-icon><connection /></el-icon>
                <span>隧道</span>
              </template>
              <el-menu-item index="/proxies/tcp">TCP</el-menu-item>
              <el-menu-item index="/proxies/udp">UDP</el-menu-item>
              <el-menu-item index="/proxies/http">HTTP</el-menu-item>
              <el-menu-item index="/proxies/https">HTTPS</el-menu-item>
              <el-menu-item index="/proxies/tcpmux">TCPMUX</el-menu-item>
              <el-menu-item index="/proxies/stcp">STCP</el-menu-item>
              <el-menu-item index="/proxies/sudp">SUDP</el-menu-item>
            </el-sub-menu>
            <el-menu-item index="">
              <el-icon><user /></el-icon>
              <span>MSL用户中心</span>
            </el-menu-item>
          </el-menu>
        </el-col>

        <el-col :xs="24" :md="21">
          <div id="content">
            <router-view></router-view>
          </div>
        </el-col>
      </el-row>
    </section>
    <footer></footer>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useDark, useToggle } from '@vueuse/core'
import { HomeFilled, Connection, User } from '@element-plus/icons-vue'

const isDark = useDark()
const darkmodeSwitch = ref(isDark)
const toggleDark = useToggle(isDark)

const handleSelect = (key: string) => {
  if (key == '') {
    window.open('https://user.mslmc.net')
  }
}
</script>

<style>
body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, Helvetica Neue, sans-serif;
  min-height: 100vh;
}

header {
  width: 100%;
  height: 60px;
  background: linear-gradient(90deg, #585bff, #986ee2);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

html.dark header {
  background: linear-gradient(90deg, #6864c1, #4860a8);
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.brand {
  display: flex;
  align-items: center;
}

.brand a {
  color: #fff;
  font-size: 25px;
  font-weight: bold;
  text-decoration: none;
  margin-left: 10px;
  line-height: 60px; /* Ensures vertical centering */
}

.logo {
  height: 40px;
}

.dark-switch {
  display: flex;
  align-items: center;
}

#content {
  margin-top: 20px;
  padding-right: 40px;
}

.el-menu-item, .el-sub-menu__title {
  display: flex;
  align-items: center;
}

.el-menu-item span, .el-sub-menu__title span {
  margin-left: 10px;
}

.el-menu-item, .el-sub-menu__title {
  justify-content: flex-start;
}


</style>