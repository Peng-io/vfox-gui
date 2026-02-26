import { defineStore } from 'pinia'
import type { SDKInfo, PluginInfo, VersionInfo, OperationResult } from '@/types/wails.d'

export const useSdkStore = defineStore('sdk', {
  state: () => ({
    sdks: [] as SDKInfo[],
    plugins: [] as PluginInfo[],
    currentVersions: {} as Record<string, string>,
    loading: false,
    error: null as string | null
  }),

  getters: {
    installedSdks: (state) => state.sdks.filter(s => s.installed),
    installedPlugins: (state) => state.plugins.filter(p => p.installed),
    availablePlugins: (state) => state.plugins.filter(p => !p.installed)
  },

  actions: {
    async loadInstalledSDKs() {
      this.loading = true
      this.error = null
      try {
        const sdks = await window.go.main.App.GetInstalledSDKs()
        this.sdks = sdks
      } catch (e) {
        this.error = String(e)
        console.error('Failed to load SDKs:', e)
      } finally {
        this.loading = false
      }
    },

    async loadAvailablePlugins() {
      this.loading = true
      this.error = null
      try {
        const plugins = await window.go.main.App.GetAvailablePlugins()
        this.plugins = plugins
      } catch (e) {
        this.error = String(e)
        console.error('Failed to load plugins:', e)
      } finally {
        this.loading = false
      }
    },

    async loadCurrentVersions() {
      try {
        const versions = await window.go.main.App.GetCurrentVersions()
        this.currentVersions = versions
      } catch (e) {
        console.error('Failed to load current versions:', e)
      }
    },

    async getAvailableVersions(sdkName: string): Promise<VersionInfo[]> {
      try {
        return await window.go.main.App.GetAvailableVersions(sdkName)
      } catch (e) {
        console.error('Failed to get available versions:', e)
        return []
      }
    },

    async installSDK(sdkName: string, version: string): Promise<OperationResult> {
      this.loading = true
      try {
        const result = await window.go.main.App.InstallSDK(sdkName, version)
        if (result.success) {
          await this.loadInstalledSDKs()
          await this.loadCurrentVersions()
        }
        return result
      } catch (e) {
        return { success: false, message: '', error: String(e) }
      } finally {
        this.loading = false
      }
    },

    async useSDK(sdkName: string, version: string, scope: string): Promise<OperationResult> {
      this.loading = true
      try {
        const result = await window.go.main.App.UseSDK(sdkName, version, scope)
        if (result.success) {
          await this.loadCurrentVersions()
        }
        return result
      } catch (e) {
        return { success: false, message: '', error: String(e) }
      } finally {
        this.loading = false
      }
    },

    async uninstallSDK(sdkName: string, version: string): Promise<OperationResult> {
      this.loading = true
      try {
        const result = await window.go.main.App.UninstallSDK(sdkName, version)
        if (result.success) {
          await this.loadInstalledSDKs()
        }
        return result
      } catch (e) {
        return { success: false, message: '', error: String(e) }
      } finally {
        this.loading = false
      }
    },

    async addPlugin(name: string, url: string = ''): Promise<OperationResult> {
      this.loading = true
      try {
        const result = await window.go.main.App.AddPlugin(name, url)
        if (result.success) {
          await this.loadAvailablePlugins()
        }
        return result
      } catch (e) {
        return { success: false, message: '', error: String(e) }
      } finally {
        this.loading = false
      }
    },

    async removePlugin(name: string): Promise<OperationResult> {
      this.loading = true
      try {
        const result = await window.go.main.App.RemovePlugin(name)
        if (result.success) {
          await this.loadAvailablePlugins()
          await this.loadInstalledSDKs()
        }
        return result
      } catch (e) {
        return { success: false, message: '', error: String(e) }
      } finally {
        this.loading = false
      }
    },

    async refresh() {
      try {
        await window.go.main.App.Refresh()
        await Promise.all([
          this.loadInstalledSDKs(),
          this.loadAvailablePlugins(),
          this.loadCurrentVersions()
        ])
      } catch (e) {
        console.error('Failed to refresh:', e)
      }
    }
  }
})
