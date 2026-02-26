<script setup lang="ts">
import { onMounted } from 'vue'
import { NCard, NSpin, NGrid, NGi, NStatistic, NButton, NIcon } from 'naive-ui'
import { RefreshOutline, CubeOutline, ExtensionPuzzleOutline } from '@vicons/ionicons5'
import { useSdkStore } from '@/stores/sdk'
import { useRouter } from 'vue-router'

const sdkStore = useSdkStore()
const router = useRouter()

onMounted(async () => {
  await Promise.all([
    sdkStore.loadInstalledSDKs(),
    sdkStore.loadAvailablePlugins(),
    sdkStore.loadCurrentVersions()
  ])
})
</script>

<template>
  <div class="home-view">
    <div class="header">
      <h1>仪表盘</h1>
      <n-button @click="sdkStore.refresh()" :loading="sdkStore.loading">
        <template #icon>
          <n-icon><RefreshOutline /></n-icon>
        </template>
        刷新
      </n-button>
    </div>

    <n-spin :show="sdkStore.loading">
      <n-grid :cols="4" :x-gap="20" :y-gap="20">
        <n-gi>
          <n-card class="stat-card card-hover" @click="router.push('/sdks')">
            <n-statistic label="已安装开发工具" :value="sdkStore.installedSdks.length">
              <template #prefix>
                <n-icon size="24"><CubeOutline /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="stat-card card-hover" @click="router.push('/plugins')">
            <n-statistic label="已安装插件" :value="sdkStore.installedPlugins.length">
              <template #prefix>
                <n-icon size="24"><ExtensionPuzzleOutline /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="stat-card">
            <n-statistic label="可用插件" :value="sdkStore.availablePlugins.length">
              <template #prefix>
                <n-icon size="24"><ExtensionPuzzleOutline /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
        <n-gi>
          <n-card class="stat-card">
            <n-statistic label="开发工具总数" :value="sdkStore.sdks.length">
              <template #prefix>
                <n-icon size="24"><CubeOutline /></n-icon>
              </template>
            </n-statistic>
          </n-card>
        </n-gi>
      </n-grid>

      <h2 style="margin-top: 30px; margin-bottom: 15px;">当前版本</h2>
      <n-grid :cols="3" :x-gap="16" :y-gap="16">
        <n-gi v-for="sdk in sdkStore.installedSdks" :key="sdk.name">
          <n-card size="small" class="version-card">
            <div class="sdk-item">
              <span class="sdk-name">{{ sdk.name }}</span>
              <span class="sdk-version">{{ sdk.current || '未设置' }}</span>
            </div>
          </n-card>
        </n-gi>
      </n-grid>
    </n-spin>
  </div>
</template>

<style scoped>
.home-view {
  max-width: 1200px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  color: #f8fafc;
}

.stat-card {
  cursor: pointer;
  background-color: #1e293b;
  border: 1px solid #334155;
}

.stat-card:hover {
  border-color: #3b82f6;
}

.version-card {
  background-color: #1e293b;
  border: 1px solid #334155;
}

.sdk-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.sdk-name {
  font-weight: 600;
  color: #f8fafc;
}

.sdk-version {
  color: #3b82f6;
  font-family: monospace;
}
</style>
