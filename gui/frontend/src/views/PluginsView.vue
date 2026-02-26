<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { 
  NCard, NSpin, NButton, NIcon, NInput, NGrid, NGi, 
  NTag, NSpace, NEmpty, useMessage
} from 'naive-ui'
import { RefreshOutline, SearchOutline, AddOutline, TrashOutline } from '@vicons/ionicons5'
import { useSdkStore } from '@/stores/sdk'

const message = useMessage()
const sdkStore = useSdkStore()

const searchQuery = ref('')
const activeTab = ref<'installed' | 'available'>('installed')

const filteredPlugins = computed(() => {
  const plugins = activeTab.value === 'installed' 
    ? sdkStore.installedPlugins 
    : sdkStore.availablePlugins
  
  if (!searchQuery.value) return plugins
  const query = searchQuery.value.toLowerCase()
  return plugins.filter(p => p.name.toLowerCase().includes(query))
})

onMounted(async () => {
  await sdkStore.loadAvailablePlugins()
})

async function addPlugin(name: string) {
  const result = await sdkStore.addPlugin(name)
  if (result.success) {
    message.success(result.message)
  } else {
    message.error(result.error || '添加插件失败')
  }
}

async function removePlugin(name: string) {
  const result = await sdkStore.removePlugin(name)
  if (result.success) {
    message.success(result.message)
  } else {
    message.error(result.error || '移除插件失败')
  }
}
</script>

<template>
  <div class="plugins-view">
    <div class="header">
      <h1>插件管理</h1>
      <n-space>
        <n-input v-model:value="searchQuery" placeholder="搜索插件..." clearable>
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
        <n-button @click="sdkStore.loadAvailablePlugins()" :loading="sdkStore.loading">
          <template #icon>
            <n-icon><RefreshOutline /></n-icon>
          </template>
          刷新
        </n-button>
      </n-space>
    </div>

    <div class="tabs">
      <n-button 
        :type="activeTab === 'installed' ? 'primary' : 'default'"
        @click="activeTab = 'installed'"
      >
        已安装 ({{ sdkStore.installedPlugins.length }})
      </n-button>
      <n-button 
        :type="activeTab === 'available' ? 'primary' : 'default'"
        @click="activeTab = 'available'"
      >
        可用 ({{ sdkStore.availablePlugins.length }})
      </n-button>
    </div>

    <n-spin :show="sdkStore.loading">
      <n-grid :cols="3" :x-gap="16" :y-gap="16" v-if="filteredPlugins.length > 0">
        <n-gi v-for="plugin in filteredPlugins" :key="plugin.name">
          <n-card :title="plugin.name" class="plugin-card card-hover">
            <template #header-extra>
              <n-tag :type="plugin.installed ? 'success' : 'default'" size="small">
                {{ plugin.installed ? '已安装' : '可安装' }}
              </n-tag>
            </template>
            
            <p class="plugin-desc">{{ plugin.description || '暂无描述' }}</p>
            <p class="plugin-version">版本：{{ plugin.version || '-' }}</p>
            
            <template #footer>
              <n-space justify="end">
                <n-button 
                  v-if="!plugin.installed"
                  size="small" 
                  type="primary"
                  @click="addPlugin(plugin.name)"
                >
                  <template #icon>
                    <n-icon><AddOutline /></n-icon>
                  </template>
                  添加
                </n-button>
                <n-button 
                  v-else
                  size="small" 
                  type="error"
                  @click="removePlugin(plugin.name)"
                >
                  <template #icon>
                    <n-icon><TrashOutline /></n-icon>
                  </template>
                  移除
                </n-button>
              </n-space>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-empty v-else :description="activeTab === 'installed' ? '暂无已安装的插件' : '暂无可用插件'" />
    </n-spin>
  </div>
</template>

<style scoped>
.plugins-view {
  max-width: 1400px;
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

.tabs {
  display: flex;
  gap: 10px;
  margin-bottom: 20px;
}

.plugin-card {
  background-color: #1e293b;
  border: 1px solid #334155;
}

.plugin-desc {
  color: #94a3b8;
  font-size: 14px;
  margin-bottom: 8px;
}

.plugin-version {
  color: #64748b;
  font-size: 12px;
}
</style>
