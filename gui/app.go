/*
 *    Copyright 2026 Han Li and contributors
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package main

import (
	"context"
	"fmt"

	"github.com/version-fox/vfox/gui/api"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct holds the application state
type App struct {
	ctx     context.Context
	service *api.Service
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize the API service
	service, err := api.NewService()
	if err != nil {
		runtime.LogFatal(ctx, fmt.Sprintf("Failed to initialize service: %v", err))
		return
	}
	a.service = service
}

// shutdown is called when the app is closing
func (a *App) shutdown(ctx context.Context) {
	if a.service != nil {
		a.service.Close()
	}
}

// GetInstalledSDKs returns all installed SDKs
func (a *App) GetInstalledSDKs() ([]api.SDKInfo, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetInstalledSDKs()
}

// GetAvailableVersions returns available versions for an SDK
func (a *App) GetAvailableVersions(sdkName string) ([]api.VersionInfo, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetAvailableVersions(sdkName)
}

// InstallSDK installs a specific version of an SDK
func (a *App) InstallSDK(sdkName, version string) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.InstallSDK(sdkName, version)
}

// UseSDK switches to a specific version of an SDK
func (a *App) UseSDK(sdkName, version, scope string) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.UseSDK(sdkName, version, scope)
}

// UninstallSDK removes a specific version of an SDK
func (a *App) UninstallSDK(sdkName, version string) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.UninstallSDK(sdkName, version)
}

// GetCurrentVersions returns current versions of all SDKs
func (a *App) GetCurrentVersions() (map[string]string, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetCurrentVersions()
}

// GetAvailablePlugins returns plugins available in registry
func (a *App) GetAvailablePlugins() ([]api.PluginInfo, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetAvailablePlugins()
}

// AddPlugin adds a plugin from registry or URL
func (a *App) AddPlugin(name, url string) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.AddPlugin(name, url)
}

// RemovePlugin removes a plugin
func (a *App) RemovePlugin(name string) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.RemovePlugin(name)
}

// SearchSDK searches for SDK versions matching a query
func (a *App) SearchSDK(sdkName, query string) ([]api.VersionInfo, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.SearchSDK(sdkName, query)
}

// GetSDKInfo returns detailed info about a specific SDK
func (a *App) GetSDKInfo(sdkName string) (*api.SDKInfo, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetSDKInfo(sdkName)
}

// CheckSDKExists checks if an SDK version is installed
func (a *App) CheckSDKExists(sdkName, version string) bool {
	if a.service == nil {
		return false
	}
	return a.service.CheckSDKExists(sdkName, version)
}

// GetConfig returns the current configuration
func (a *App) GetConfig() (map[string]interface{}, error) {
	if a.service == nil {
		return nil, fmt.Errorf("service not initialized")
	}
	return a.service.GetConfig()
}

// SaveConfig saves the configuration
func (a *App) SaveConfig(update api.ConfigUpdate) api.OperationResult {
	if a.service == nil {
		return api.OperationResult{
			Success: false,
			Error:   "service not initialized",
		}
	}
	return a.service.SaveConfig(update)
}

// Refresh reloads all SDK data
func (a *App) Refresh() error {
	if a.service == nil {
		return fmt.Errorf("service not initialized")
	}

	// Close existing service
	a.service.Close()

	// Create new service
	service, err := api.NewService()
	if err != nil {
		return err
	}
	a.service = service
	return nil
}

// ShowAboutDialog shows the about dialog
func (a *App) ShowAboutDialog() {
	runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   "关于 vfox-gui",
		Message: "vfox 图形界面 - vfox (Version Fox) 的图形化前端\n\n版本：1.0.0\n\n跨平台 SDK 版本管理器，配备现代化图形界面。",
	})
}

// SelectDirectory opens a directory picker dialog
func (a *App) SelectDirectory() (string, error) {
	return runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "选择目录",
	})
}
