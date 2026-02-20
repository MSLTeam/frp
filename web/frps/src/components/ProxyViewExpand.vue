<template>
  <el-form
    label-position="left"
    label-width="auto"
    inline
    class="proxy-table-expand"
  >
    <el-form-item label="完整隧道名">
      <span>{{ row.name }}</span>
    </el-form-item>
    <el-form-item label="类型">
      <span>{{ row.type }}</span>
    </el-form-item>
    <el-form-item label="加密">
      <span>{{ row.encryption }}</span>
    </el-form-item>
    <el-form-item label="压缩">
      <span>{{ row.compression }}</span>
    </el-form-item>
    <el-form-item label="上次启动时间">
      <span>{{ row.lastStartTime }}</span>
    </el-form-item>
    <el-form-item label="上次关闭时间">
      <span>{{ row.lastCloseTime }}</span>
    </el-form-item>

    <div v-if="proxyType === 'http' || proxyType === 'https'" class="expand-group">
      <el-form-item label="域名列表">
        <span>{{ row.customDomains }}</span>
      </el-form-item>
      <el-form-item label="子域名">
        <span>{{ row.subdomain }}</span>
      </el-form-item>
      <el-form-item label="位置">
        <span>{{ row.locations }}</span>
      </el-form-item>
      <el-form-item label="域名重写">
        <span>{{ row.hostHeaderRewrite }}</span>
      </el-form-item>
    </div>

    <div v-else-if="proxyType === 'tcpmux'" class="expand-group">
      <el-form-item label="多路复用">
        <span>{{ row.multiplexer }}</span>
      </el-form-item>
      <el-form-item label="用户路由">
        <span>{{ row.routeByHTTPUser }}</span>
      </el-form-item>
      <el-form-item label="域名列表">
        <span>{{ row.customDomains }}</span>
      </el-form-item>
      <el-form-item label="子域名">
        <span>{{ row.subdomain }}</span>
      </el-form-item>
    </div>

    <div v-else class="expand-group">
      <el-form-item label="连接地址">
        <span>{{ row.addr }}</span>
      </el-form-item>
    </div>
  </el-form>

  <div v-if="row.annotations && row.annotations.size > 0" class="annotations-wrap">
    <el-divider />
    <el-text class="title-text" size="large">Annotations</el-text>
    <ul>
      <li v-for="item in annotationsArray()" :key="item.key">
        <span class="annotation-key">{{ item.key }}</span>
        <span class="annotation-value">{{ item.value }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup lang="ts">
const props = defineProps<{
  row: any
  proxyType: string
}>()

const annotationsArray = (): Array<{ key: string; value: string }> => {
  const array: Array<{ key: string; value: any }> = []
  if (props.row.annotations) {
    props.row.annotations.forEach((value: any, key: string) => {
      array.push({ key, value })
    })
  }
  return array
}
</script>

<style scoped>
.expand-group {
  display: contents;
}

.annotations-wrap {
  margin-top: 6px;
}

ul {
  list-style-type: none;
  padding: 0;
  margin: 8px 0 0;
}

ul li {
  display: grid;
  grid-template-columns: minmax(180px, 280px) 1fr;
  gap: 10px;
  padding: 8px 10px;
  border-radius: 10px;
  border: 1px solid var(--frp-border);
  background: var(--frp-surface-bg);
  margin-bottom: 8px;
}

.annotation-key {
  font-weight: 700;
  color: var(--frp-muted);
  overflow-wrap: anywhere;
}

.annotation-value {
  overflow-wrap: anywhere;
}

.title-text {
  color: var(--frp-muted);
  font-weight: 700;
}

@media (max-width: 900px) {
  ul li {
    grid-template-columns: 1fr;
  }
}
</style>
