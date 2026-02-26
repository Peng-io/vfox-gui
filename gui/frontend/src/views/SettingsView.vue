<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { NCard, NForm, NFormItem, NInput, NButton, NSpin, NDivider, NIcon, NSwitch, NSpace, useMessage } from 'naive-ui'
import { SaveOutline, RefreshOutline, FolderOpenOutline } from '@vicons/ionicons5'

const message = useMessage()

const config = ref({
  proxy: {
    url: '',
    enable: false
  },
  storage: {
    sdkPath: ''
  },
  registry: {
    address: ''
  },
  cache: {
    availableHookDuration: ''
  },
  legacyVersionFile: {
    enable: false,
    strategy: ''
  }
})

const loading = ref(false)
const saving = ref(false)

onMounted(async () => {
  await loadConfig()
})

async function loadConfig() {
  loading.value = true
  try {
    const cfg = await window.go.main.App.GetConfig()
    
    // 代理设置
    if (cfg.proxy) {
      const proxy = cfg.proxy as { url: string; enable: boolean }
      config.value.proxy.url = proxy.url || ''
      config.value.proxy.enable = proxy.enable || false
    }
    
    // 存储设置
    if (cfg.storage) {
      const storage = cfg.storage as { sdkPath: string }
      config.value.storage.sdkPath = storage.sdkPath || ''
    }
    
    // 注册表设置
    if (cfg.registry) {
      const registry = cfg.registry as { address: string }
      config.value.registry.address = registry.address || ''
    }
    
    // 缓存设置
    if (cfg.cache) {
      const cache = cfg.cache as { availableHookDuration: string }
      config.value.cache.availableHookDuration = cache.availableHookDuration || ''
    }
    
    // 遗留版本文件设置
    if (cfg.legacyVersionFile) {
      const lvf = cfg.legacyVersionFile as { enable: boolean; strategy: string }
      config.value.legacyVersionFile.enable = lvf.enable || false
      config.value.legacyVersionFile.strategy = lvf.strategy || ''
    }
  } catch (e) {
    message.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

async function saveConfig() {
  saving.value = true
  try {
    const result = await window.go.main.App.SaveConfig({
      proxy: {
        url: config.value.proxy.url,
        enable: config.value.proxy.enable
      },
      storage: {
        sdkPath: config.value.storage.sdkPath
      },
      registry: {
        address: config.value.registry.address
      },
      cache: {
        availableHookDuration: config.value.cache.availableHookDuration
      },
      legacyVersionFile: {
        enable: config.value.legacyVersionFile.enable,
        strategy: config.value.legacyVersionFile.strategy
      }
    })
    
    if (result.success) {
      message.success(result.message || '配置已保存')
    } else {
      message.error(result.error || '保存配置失败')
    }
  } catch (e) {
    message.error('保存配置失败')
  } finally {
    saving.value = false
  }
}

async function selectSdkPath() {
  try {
    const path = await window.go.main.App.SelectDirectory()
    if (path) {
      config.value.storage.sdkPath = path
    }
  } catch (e) {
    // 用户取消选择
  }
}
</script>

<template>
  <div class="settings-view">
    <div class="header">
      <h1>设置</h1>
    </div>

    <n-spin :show="loading">
      <n-card class="settings-card">
        <n-form label-placement="left" label-width="140px">
          <!-- 代理设置 -->
          <h3>代理设置</h3>
          <n-form-item label="启用代理">
            <n-switch v-model:value="config.proxy.enable" />
          </n-form-item>
          <n-form-item label="代理地址">
            <n-input 
              v-model:value="config.proxy.url" 
              placeholder="例如：http://127.0.0.1:7890 或 socks5://127.0.0.1:1080"
              :disabled="!config.proxy.enable"
            />
          </n-form-item>

          <n-divider />

          <!-- 存储设置 -->
          <h3>存储设置</h3>
          <n-form-item label="SDK 安装路径">
            <n-input-group>
              <n-input 
                v-model:value="config.storage.sdkPath" 
                placeholder="自定义 SDK 安装路径（留空使用默认路径）"
                style="flex: 1"
              />
              <n-button @click="selectSdkPath">
                <template #icon>
                  <n-icon><FolderOpenOutline /></n-icon>
                </template>
                浏览
              </n-button>
            </n-input-group>
          </n-form-item>

          <n-divider />

          <!-- 注册表设置 -->
          <h3>注册表设置</h3>
          <n-form-item label="注册表地址">
            <n-input 
              v-model:value="config.registry.address" 
              placeholder="自定义注册表地址（留空使用默认值）"
            />
          </n-form-item>

          <n-divider />

          <!-- 缓存设置 -->
          <h3>缓存设置</h3>
          <n-form-item label="版本列表缓存">
            <n-input 
              v-model:value="config.cache.availableHookDuration" 
              placeholder="例如：12h（12小时）、30m（30分钟）"
            />
          </n-form-item>

          <n-divider />

          <!-- 遗留版本文件设置 -->
          <h3>遗留版本文件</h3>
          <n-form-item label="启用解析">
            <n-switch v-model:value="config.legacyVersionFile.enable" />
          </n-form-item>
          <n-form-item label="解析策略">
            <n-input 
              v-model:value="config.legacyVersionFile.strategy" 
              placeholder="遗留文件解析策略"
              :disabled="!config.legacyVersionFile.enable"
            />
          </n-form-item>

          <n-divider />

          <!-- 操作按钮 -->
          <n-space justify="end">
            <n-button @click="loadConfig" :loading="loading">
              <template #icon>
                <n-icon><RefreshOutline /></n-icon>
              </template>
              重置
            </n-button>
            <n-button type="primary" @click="saveConfig" :loading="saving">
              <template #icon>
                <n-icon><SaveOutline /></n-icon>
              </template>
              保存
            </n-button>
          </n-space>
        </n-form>
      </n-card>
    </n-spin>
  </div>
</template>

<style scoped>
.settings-view {
  max-width: 800px;
}

.header {
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  color: #f8fafc;
}

.settings-card {
  background-color: #1e293b;
  border: 1px solid #334155;
}

h3 {
  color: #f8fafc;
  margin-bottom: 16px;
}
</style>
