# vfox-gui

A graphical user interface for [vfox](https://github.com/version-fox/vfox) (Version Fox), built with Wails and Vue 3.

## Features

- **SDK Management**: View, install, and manage multiple SDK versions
- **Plugin Management**: Browse and install plugins from the registry
- **Version Switching**: Easily switch between SDK versions with different scopes
- **Modern UI**: Clean, dark-themed interface built with Naive UI
- **Cross-Platform**: Works on Windows, Linux, and macOS

## Prerequisites

- [Go 1.24+](https://golang.org/dl/)
- [Node.js 18+](https://nodejs.org/)
- [Wails CLI v2](https://wails.io/docs/gettingstarted/installation/)

## Installation

### Install Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### Clone and Build

```bash
cd gui

# Install frontend dependencies
cd frontend && npm install && cd ..

# Run in development mode
wails dev

# Build for production
wails build
```

## Project Structure

```
gui/
├── api/
│   └── service.go        # API service layer
├── app.go                # Wails app bindings
├── main.go               # Application entry point
├── wails.json            # Wails configuration
├── go.mod                # Go module file
├── frontend/
│   ├── src/
│   │   ├── components/   # Vue components
│   │   ├── views/        # Page views
│   │   ├── stores/       # Pinia stores
│   │   ├── router/       # Vue Router config
│   │   ├── styles/       # CSS styles
│   │   └── types/        # TypeScript types
│   ├── package.json
│   └── vite.config.ts
└── Makefile              # Build automation
```

## Development

```bash
# Run development server with hot reload
make dev

# Or directly
wails dev
```

## Building

```bash
# Build for current platform
make build

# Build for specific platforms
make build-windows
make build-linux
make build-darwin

# Build for all platforms
make build-all
```

## API Reference

The GUI exposes the following methods through Wails bindings:

### SDK Operations

| Method | Description |
|--------|-------------|
| `GetInstalledSDKs()` | Get all installed SDKs |
| `GetAvailableVersions(sdkName)` | Get available versions for an SDK |
| `InstallSDK(sdkName, version)` | Install a specific version |
| `UseSDK(sdkName, version, scope)` | Switch to a version |
| `UninstallSDK(sdkName, version)` | Remove a version |

### Plugin Operations

| Method | Description |
|--------|-------------|
| `GetAvailablePlugins()` | Get plugins from registry |
| `AddPlugin(name, url)` | Add a plugin |
| `RemovePlugin(name)` | Remove a plugin |

### Configuration

| Method | Description |
|--------|-------------|
| `GetConfig()` | Get current configuration |
| `Refresh()` | Reload all data |

## License

Apache License 2.0
