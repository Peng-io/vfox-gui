# vfox-gui

[![GitHub License](https://img.shields.io/github/license/version-fox/vfox?style=for-the-badge)](LICENSE)
[![GitHub Release](https://img.shields.io/github/v/release/version-fox/vfox?display_name=tag&style=for-the-badge)](https://github.com/version-fox/vfox/releases)

[Version Fox](https://github.com/version-fox/vfox) 的图形化桌面应用，让你通过直观的可视化界面管理 SDK 版本和插件，无需命令行操作。

> **原项目地址**：[vfox - https://github.com/version-fox/vfox](https://github.com/version-fox/vfox)

**当前版本：v1.0.6**

## 功能特性

- **仪表盘**：实时查看已安装的 SDK 数量、插件数量，以及当前激活的 SDK 版本概览
- **SDK 管理**：以卡片形式展示已安装的 SDK，支持安装、卸载、切换版本，并可选择作用域（全局 / 项目 / 会话）
- **插件管理**：浏览注册中心中的可用插件，一键安装或移除已安装的插件
- **设置中心**：图形化配置代理、存储路径、镜像源、缓存时长等 vfox 核心参数
- **深色主题**：基于 Naive UI 构建的现代化深色界面
- **跨平台**：支持 Windows、Linux 和 macOS

## 技术栈

| 层级 | 技术 |
|------|------|
| 桌面框架 | [Wails v2](https://wails.io/) |
| 后端语言 | Go 1.24+ |
| 前端框架 | Vue 3 + TypeScript |
| 状态管理 | Pinia |
| 路由 | Vue Router 4 |
| UI 组件库 | Naive UI (Dark Theme) |
| 构建工具 | Vite 5 |
| 核心引擎 | vfox SDK Manager（直接嵌入，无需外部依赖） |

## 前置要求

- [Go 1.24+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation/)

## 快速开始

### 1. 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 2. 克隆仓库并进入 GUI 目录

```bash
git clone https://github.com/version-fox/vfox.git
cd vfox/gui
```

### 3. 开发模式（支持热重载）

```bash
make dev
# 或
wails dev
```

### 4. 生产构建

```bash
# 为当前平台构建
make build

# 为指定平台构建
make build-windows
make build-linux
make build-darwin

# 为所有平台构建
make build-all
```

构建产物位于 `gui/build/bin/` 目录下。

## 项目结构

```
gui/
├── api/
│   └── service.go            # API 服务层，封装 vfox internal.Manager
├── app.go                    # Wails 应用入口与前端绑定
├── main.go                   # 程序入口
├── wails.json                # Wails 配置
├── go.mod                    # Go 模块（通过 replace 指令引用 vfox 核心）
├── Makefile                  # 构建脚本
├── frontend/                 # Vue 3 前端应用
│   ├── src/
│   │   ├── components/       # 通用组件
│   │   │   └── AppLayout.vue # 侧边栏布局
│   │   ├── views/            # 页面视图
│   │   │   ├── HomeView.vue  # 仪表盘
│   │   │   ├── SDKsView.vue  # SDK 管理
│   │   │   ├── PluginsView.vue # 插件管理
│   │   │   └── SettingsView.vue # 设置
│   │   ├── stores/           # Pinia 状态管理
│   │   ├── router/           # 路由配置
│   │   ├── styles/           # 全局样式
│   │   └── types/            # TypeScript 类型定义（Wails 绑定）
│   ├── package.json
│   └── vite.config.ts
```

## 架构说明

GUI 应用通过 `go.mod` 中的 `replace` 指令直接引用 vfox 核心代码库（`../`），将 vfox 的 `internal.Manager`、`internal/sdk`、`internal/env` 等内部包嵌入到桌面应用中。这意味着 **GUI 是一个自包含的桌面应用，不依赖外部 vfox 二进制文件**。

```
┌─────────────────────────────────────┐
│          Wails Desktop App          │
│  ┌───────────────────────────────┐  │
│  │       Vue 3 Frontend          │  │
│  │  (Naive UI + Pinia + Router)  │  │
│  └──────────────┬────────────────┘  │
│                 │ Wails Bridge       │
│  ┌──────────────▼────────────────┐  │
│  │       Go Backend (app.go)     │  │
│  │  ┌─────────────────────────┐  │  │
│  │  │  vfox internal.Manager  │  │  │
│  │  │  (SDK / Plugin / Config)│  │  │
│  │  └─────────────────────────┘  │  │
│  └───────────────────────────────┘  │
└─────────────────────────────────────┘
```

## API 参考

后端通过 Wails 绑定向前端暴露以下方法：

### SDK 操作

| 方法 | 说明 |
|------|------|
| `GetInstalledSDKs()` | 获取所有已安装的 SDK |
| `GetAvailableVersions(sdkName)` | 获取指定 SDK 的可用版本 |
| `InstallSDK(sdkName, version)` | 安装指定版本 |
| `UseSDK(sdkName, version, scope)` | 切换到指定版本（global / project / session） |
| `UninstallSDK(sdkName, version)` | 卸载指定版本 |
| `GetCurrentVersions()` | 获取各 SDK 当前激活的版本 |
| `SearchSDK(sdkName, query)` | 按关键字搜索版本 |
| `GetSDKInfo(sdkName)` | 获取 SDK 详细信息 |
| `CheckSDKExists(sdkName, version)` | 检查某版本是否已安装 |

### 插件操作

| 方法 | 说明 |
|------|------|
| `GetAvailablePlugins()` | 从注册中心获取可用插件（5 分钟缓存） |
| `AddPlugin(name, url)` | 安装插件（从注册中心或 URL） |
| `RemovePlugin(name)` | 移除已安装的插件 |

### 配置与工具

| 方法 | 说明 |
|------|------|
| `GetConfig()` | 获取当前配置 |
| `SaveConfig(update)` | 保存配置 |
| `Refresh()` | 重新加载所有 SDK 数据 |
| `SelectDirectory()` | 打开原生目录选择器 |
| `ShowAboutDialog()` | 显示关于对话框 |

## 许可证

[Apache License 2.0](LICENSE) - Copyright (C) 2026 Han Li and [contributors](https://github.com/version-fox/vfox/graphs/contributors)
