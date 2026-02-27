<script setup lang="ts">
import { ref, h } from 'vue'
import { useRouter } from 'vue-router'
import { NLayout, NLayoutSider, NLayoutContent, NMenu, NIcon } from 'naive-ui'
import {
  HomeOutline,
  CubeOutline,
  ExtensionPuzzleOutline,
  SettingsOutline
} from '@vicons/ionicons5'

const router = useRouter()
const collapsed = ref(true)

const menuOptions = [
  {
    label: '首页',
    key: 'home',
    icon: () => h(NIcon, null, { default: () => h(HomeOutline) })
  },
  {
    label: '开发工具',
    key: 'sdks',
    icon: () => h(NIcon, null, { default: () => h(CubeOutline) })
  },
  {
    label: '插件',
    key: 'plugins',
    icon: () => h(NIcon, null, { default: () => h(ExtensionPuzzleOutline) })
  },
  {
    label: '设置',
    key: 'settings',
    icon: () => h(NIcon, null, { default: () => h(SettingsOutline) })
  }
]

const currentKey = ref('home')

function handleMenuSelect(key: string) {
  currentKey.value = key
  router.push({ name: key })
}
</script>

<template>
  <n-layout has-sider style="height: 100vh">
    <n-layout-sider
      bordered
      collapse-mode="width"
      :collapsed-width="64"
      :width="200"
      :collapsed="collapsed"
      show-trigger
      @collapse="collapsed = true"
      @expand="collapsed = false"
      :native-scrollbar="false"
      style="background-color: #1e293b"
    >
      <div class="logo">
        <span v-if="!collapsed" class="logo-text">vfox-gui</span>
        <span v-else class="logo-icon">🦊</span>
      </div>
      <n-menu
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="22"
        :options="menuOptions"
        :value="currentKey"
        @update:value="handleMenuSelect"
        style="background-color: #1e293b"
      />
    </n-layout-sider>
    <n-layout-content :native-scrollbar="false" style="background-color: #1b2636">
      <div class="content-wrapper">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </n-layout-content>
  </n-layout>
</template>

<style scoped>
.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom: 1px solid #334155;
}

.logo-text {
  font-size: 20px;
  font-weight: bold;
  color: #f8fafc;
}

.logo-icon {
  font-size: 28px;
}

.content-wrapper {
  padding: 20px;
  height: 100%;
  overflow-y: auto;
}
</style>
