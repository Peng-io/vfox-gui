# vfox-gui

[Version Fox](https://github.com/version-fox/vfox) 的图形用户界面，基于 Wails 和 Vue 3 构建。

**当前版本：v1.0.6**

## 功能特性

- **SDK 管理**：查看、安装和管理多个 SDK 版本
- **插件管理**：浏览和安装注册中心中的插件
- **版本切换**：轻松在不同作用域间切换 SDK 版本
- **现代化 UI**：简洁的深色主题界面，采用 Naive UI 构建
- **跨平台支持**：支持 Windows、Linux 和 macOS

## 前置要求

- [Go 1.24+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation/)

## 安装步骤

### 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 克隆并构建

```bash
cd gui

# 安装前端依赖
cd frontend && npm install && cd ..

# 开发模式运行
wails dev

# 生产环境构建
wails build
```

## 项目结构

```
gui/
├── api/
│   └── service.go        # API 服务层
├── app.go                # Wails 应用绑定
├── main.go               # 应用入口
├── wails.json            # Wails 配置文件
├── go.mod                # Go 模块文件
├── frontend/
│   ├── src/
│   │   ├── components/   # Vue 组件
│   │   ├── views/        # 页面视图
│   │   ├── stores/       # Pinia 状态管理
│   │   ├── router/       # Vue Router 配置
│   │   ├── styles/       # CSS 样式
│   │   └── types/        # TypeScript 类型定义
│   ├── package.json
│   └── vite.config.ts
└── Makefile              # 构建自动化脚本
```

## 开发指南

```bash
# 运行开发服务器（支持热重载）
make dev

# 或直接运行
wails dev
```

## 构建说明

```bash
# 为当前平台构建
make build

# 为特定平台构建
make build-windows
make build-linux
make build-darwin

# 为所有平台构建
make build-all
```

## API 参考

GUI 通过 Wails 绑定暴露以下方法：

### SDK 操作

| 方法 | 说明 |
|--------|-------------|
| `GetInstalledSDKs()` | 获取所有已安装的 SDK |
| `GetAvailableVersions(sdkName)` | 获取指定 SDK 的可用版本 |
| `InstallSDK(sdkName, version)` | 安装指定版本 |
| `UseSDK(sdkName, version, scope)` | 切换到指定版本（带作用域） |
| `UninstallSDK(sdkName, version)` | 卸载指定版本 |

### 插件操作

| 方法 | 说明 |
|--------|-------------|
| `GetAvailablePlugins()` | 从注册中心获取可用插件 |
| `AddPlugin(name, url)` | 添加插件 |
| `RemovePlugin(name)` | 移除插件 |

### 配置管理

| 方法 | 说明 |
|--------|-------------|
| `GetConfig()` | 获取当前配置 |
| `Refresh()` | 重新加载所有数据 |

## 许可证

Apache License 2.0
