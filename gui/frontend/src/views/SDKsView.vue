<script setup lang="ts">
import { onMounted, ref, computed, h } from 'vue'
import { 
  NCard, NSpin, NButton, NIcon, NInput, NGrid, NGi, 
  NTag, NSpace, NModal, NSelect, NDataTable, NDropdown, NEmpty, useMessage
} from 'naive-ui'
import { 
  RefreshOutline, SearchOutline, DownloadOutline,
  PlayOutline
} from '@vicons/ionicons5'
import { useSdkStore } from '@/stores/sdk'
import type { VersionInfo } from '@/types/wails.d'

const message = useMessage()
const sdkStore = useSdkStore()

const searchQuery = ref('')
const showVersionsModal = ref(false)
const selectedSdk = ref<string | null>(null)
const selectedScope = ref<string>('global')
const availableVersions = ref<VersionInfo[]>([])
const loadingVersions = ref(false)

const scopeOptions = [
  { label: '全局', value: 'global' },
  { label: '项目', value: 'project' },
  { label: '会话', value: 'session' }
]

const filteredSdks = computed(() => {
  if (!searchQuery.value) return sdkStore.sdks
  const query = searchQuery.value.toLowerCase()
  return sdkStore.sdks.filter(s => s.name.toLowerCase().includes(query))
})

onMounted(async () => {
  await sdkStore.loadInstalledSDKs()
})

async function openVersionsModal(sdkName: string) {
  selectedSdk.value = sdkName
  showVersionsModal.value = true
  loadingVersions.value = true
  try {
    availableVersions.value = await sdkStore.getAvailableVersions(sdkName)
  } catch (e) {
    message.error('加载版本列表失败')
  } finally {
    loadingVersions.value = false
  }
}

async function installVersion(version: string) {
  if (!selectedSdk.value) return
  
  const result = await sdkStore.installSDK(selectedSdk.value, version)
  if (result.success) {
    message.success(result.message)
    await openVersionsModal(selectedSdk.value)
  } else {
    message.error(result.error || '安装失败')
  }
}

async function useVersion(sdkName: string, version: string) {
  const result = await sdkStore.useSDK(sdkName, version, selectedScope.value)
  if (result.success) {
    message.success(result.message)
  } else {
    message.error(result.error || '切换版本失败')
  }
}


</script>

<template>
  <div class="sdks-view">
    <div class="header">
      <h1>开发工具管理</h1>
      <n-space>
        <n-input v-model:value="searchQuery" placeholder="搜索开发工具..." clearable>
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
        <n-button @click="sdkStore.loadInstalledSDKs()" :loading="sdkStore.loading">
          <template #icon>
            <n-icon><RefreshOutline /></n-icon>
          </template>
          刷新
        </n-button>
      </n-space>
    </div>

    <n-spin :show="sdkStore.loading">
      <n-grid :cols="3" :x-gap="16" :y-gap="16" v-if="filteredSdks.length > 0">
        <n-gi v-for="sdk in filteredSdks" :key="sdk.name">
          <n-card :title="sdk.name" class="sdk-card card-hover">
            <template #header-extra>
              <n-tag :type="sdk.current ? 'success' : 'default'" size="small">
                {{ sdk.current || '未激活' }}
              </n-tag>
            </template>
            
            <p class="sdk-desc">{{ sdk.description || '暂无描述' }}</p>
            
            <div class="versions-list" v-if="sdk.versions.length > 0">
              <span class="label">已安装：</span>
              <n-space size="small">
                <n-tag v-for="v in sdk.versions.slice(0, 3)" :key="v" size="small" round>
                  {{ v }}
                </n-tag>
                <n-tag v-if="sdk.versions.length > 3" size="small" round>
                  还有 {{ sdk.versions.length - 3 }} 个
                </n-tag>
              </n-space>
            </div>

            <template #footer>
              <n-space justify="space-between">
                <n-button size="small" @click="openVersionsModal(sdk.name)">
                  <template #icon>
                    <n-icon><DownloadOutline /></n-icon>
                  </template>
                  安装
                </n-button>
                <n-dropdown
                  v-if="sdk.versions.length > 0"
                  :options="sdk.versions.map(v => ({ label: v, key: v }))"
                  @select="(key: string) => useVersion(sdk.name, key)"
                >
                  <n-button size="small">
                    <template #icon>
                      <n-icon><PlayOutline /></n-icon>
                    </template>
                    使用
                  </n-button>
                </n-dropdown>
              </n-space>
            </template>
          </n-card>
        </n-gi>
      </n-grid>
      <n-empty v-else description="暂无已安装的开发工具" />
    </n-spin>

    <!-- Versions Modal -->
    <n-modal v-model:show="showVersionsModal" preset="card" style="width: 600px" title="可用版本">
      <n-spin :show="loadingVersions">
        <div class="scope-selector">
          <span>作用域：</span>
          <n-select v-model:value="selectedScope" :options="scopeOptions" style="width: 150px" />
        </div>
        <n-data-table
          :columns="[
            { title: '版本', key: 'version' },
            { title: '状态', key: 'installed', render: (row: VersionInfo) => row.installed ? '已安装' : '可安装' },
            { 
              title: '操作', 
              key: 'actions',
              render: (row: VersionInfo) => {
                if (row.installed) {
                  return h(NButton, { 
                    size: 'small', 
                    type: 'primary',
                    onClick: () => useVersion(selectedSdk!, row.version)
                  }, { default: () => '使用' })
                }
                return h(NButton, { 
                  size: 'small',
                  onClick: () => installVersion(row.version)
                }, { default: () => '安装' })
              }
            }
          ]"
          :data="availableVersions"
          :max-height="400"
          virtual-scroll
        />
      </n-spin>
    </n-modal>
  </div>
</template>

<style scoped>
.sdks-view {
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

.sdk-card {
  background-color: #1e293b;
  border: 1px solid #334155;
}

.sdk-desc {
  color: #94a3b8;
  font-size: 14px;
  margin-bottom: 12px;
}

.versions-list {
  margin-top: 8px;
}

.versions-list .label {
  display: block;
  color: #64748b;
  font-size: 12px;
  margin-bottom: 4px;
}

.scope-selector {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 16px;
}
</style>
