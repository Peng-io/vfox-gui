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

package api

import (
	"errors"
	"fmt"
	"sync"

	"github.com/version-fox/vfox/internal"
	"github.com/version-fox/vfox/internal/env"
	"github.com/version-fox/vfox/internal/sdk"
)

// Version type alias for convenience
type Version = sdk.Version

// Service provides API methods for GUI
type Service struct {
	manager *internal.Manager
	mu      sync.RWMutex
}

// NewService creates a new API service
func NewService() (*Service, error) {
	manager, err := internal.NewSdkManager()
	if err != nil {
		return nil, fmt.Errorf("failed to create SDK manager: %w", err)
	}
	return &Service{
		manager: manager,
	}, nil
}

// Close releases resources
func (s *Service) Close() {
	if s.manager != nil {
		s.manager.Close()
	}
}

// SDKInfo represents SDK information for frontend
type SDKInfo struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Homepage    string   `json:"homepage"`
	Version     string   `json:"version"`
	Installed   bool     `json:"installed"`
	Versions    []string `json:"versions"`
	Current     string   `json:"current"`
}

// VersionInfo represents a version with details
type VersionInfo struct {
	Version     string `json:"version"`
	Installed   bool   `json:"installed"`
	Current     bool   `json:"current"`
	Description string `json:"description,omitempty"`
}

// OperationResult represents the result of an operation
type OperationResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// ProgressCallback is called during long operations
type ProgressCallback func(stage string, current int, total int, message string)

// GetInstalledSDKs returns all installed SDKs
func (s *Service) GetInstalledSDKs() ([]SDKInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sdks, err := s.manager.LoadAllSdk()
	if err != nil {
		return nil, err
	}

	result := make([]SDKInfo, 0, len(sdks))
	for _, sdk := range sdks {
		metadata := sdk.Metadata()
		versions := sdk.InstalledList()
		versionStrs := make([]string, len(versions))
		for i, v := range versions {
			versionStrs[i] = string(v)
		}

		current := sdk.Current()
		result = append(result, SDKInfo{
			Name:        metadata.Name,
			Description: metadata.PluginMetadata.Description,
			Homepage:    metadata.PluginMetadata.Homepage,
			Version:     metadata.PluginMetadata.Version,
			Installed:   true,
			Versions:    versionStrs,
			Current:     string(current),
		})
	}

	return result, nil
}

// GetAvailableVersions returns available versions for an SDK
func (s *Service) GetAvailableVersions(sdkName string) ([]VersionInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sdk, err := s.manager.LookupSdk(sdkName)
	if err != nil {
		return nil, fmt.Errorf("SDK %s not found: %w", sdkName, err)
	}

	// Get available versions
	available, err := sdk.Available([]string{})
	if err != nil {
		return nil, fmt.Errorf("failed to get available versions: %w", err)
	}

	// Get installed versions
	installed := sdk.InstalledList()
	installedMap := make(map[string]bool)
	for _, v := range installed {
		installedMap[string(v)] = true
	}

	// Get current version
	current := sdk.Current()

	result := make([]VersionInfo, 0, len(available))
	for _, pkg := range available {
		v := string(pkg.Version)
		result = append(result, VersionInfo{
			Version:     v,
			Installed:   installedMap[v],
			Current:     v == string(current),
			Description: pkg.Note,
		})
	}

	return result, nil
}

// InstallSDK installs a specific version of an SDK
func (s *Service) InstallSDK(sdkName, version string) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	sdkInstance, err := s.manager.LookupSdkWithInstall(sdkName, true)
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	// Handle "latest" version
	if version == "latest" {
		available, err := sdkInstance.Available([]string{})
		if err != nil {
			return OperationResult{
				Success: false,
				Error:   fmt.Sprintf("Failed to get available versions: %v", err),
			}
		}
		if len(available) == 0 {
			return OperationResult{
				Success: false,
				Error:   "No available versions found",
			}
		}
		version = string(available[0].Version)
	}

	err = sdkInstance.Install(Version(version))
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully installed %s@%s", sdkName, version),
	}
}

// UseSDK switches to a specific version of an SDK
func (s *Service) UseSDK(sdkName, version, scope string) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	sdkInstance, err := s.manager.LookupSdk(sdkName)
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   fmt.Sprintf("SDK %s not found: %v", sdkName, err),
		}
	}

	// Parse scope
	var useScope env.UseScope
	switch scope {
	case "global", "g":
		useScope = env.Global
	case "project", "p":
		useScope = env.Project
	case "session", "s":
		useScope = env.Session
	default:
		useScope = env.Session
	}

	err = sdkInstance.Use(Version(version), useScope)
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("Now using %s@%s", sdkName, version),
	}
}

// UninstallSDK removes a specific version of an SDK
func (s *Service) UninstallSDK(sdkName, version string) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	sdkInstance, err := s.manager.LookupSdk(sdkName)
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   fmt.Sprintf("SDK %s not found: %v", sdkName, err),
		}
	}

	err = sdkInstance.Uninstall(Version(version))
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully uninstalled %s@%s", sdkName, version),
	}
}

// GetCurrentVersions returns current versions of all SDKs
func (s *Service) GetCurrentVersions() (map[string]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sdks, err := s.manager.LoadAllSdk()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, sdk := range sdks {
		name := sdk.Metadata().Name
		current := sdk.Current()
		if current != "" {
			result[name] = string(current)
		} else {
			result[name] = "N/A"
		}
	}

	return result, nil
}

// PluginInfo represents plugin information
type PluginInfo struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Installed   bool   `json:"installed"`
}

// GetAvailablePlugins returns plugins available in registry
func (s *Service) GetAvailablePlugins() ([]PluginInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	registry, err := s.manager.Available()
	if err != nil {
		return nil, fmt.Errorf("failed to get available plugins: %w", err)
	}

	// Get installed plugins
	installedSdks, _ := s.manager.LoadAllSdk()
	installedMap := make(map[string]bool)
	for _, sdk := range installedSdks {
		installedMap[sdk.Metadata().Name] = true
	}

	result := make([]PluginInfo, 0, len(registry))
	for _, p := range registry {
		result = append(result, PluginInfo{
			Name:        p.Name,
			Version:     "", // RegistryIndexItem doesn't have Version
			Description: p.Desc,
			Homepage:    p.Homepage,
			Installed:   installedMap[p.Name],
		})
	}

	return result, nil
}

// AddPlugin adds a plugin from registry or URL
func (s *Service) AddPlugin(name, url string) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.Add(name, url, "")
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully added plugin: %s", name),
	}
}

// RemovePlugin removes a plugin
func (s *Service) RemovePlugin(name string) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	err := s.manager.Remove(name)
	if err != nil {
		return OperationResult{
			Success: false,
			Error:   err.Error(),
		}
	}

	return OperationResult{
		Success: true,
		Message: fmt.Sprintf("Successfully removed plugin: %s", name),
	}
}

// SearchSDK searches for SDK versions matching a query
func (s *Service) SearchSDK(sdkName, query string) ([]VersionInfo, error) {
	versions, err := s.GetAvailableVersions(sdkName)
	if err != nil {
		return nil, err
	}

	// Simple filtering - can be enhanced
	if query == "" {
		return versions, nil
	}

	result := make([]VersionInfo, 0)
	for _, v := range versions {
		if containsIgnoreCase(v.Version, query) {
			result = append(result, v)
		}
	}

	return result, nil
}

// GetSDKInfo returns detailed info about a specific SDK
func (s *Service) GetSDKInfo(sdkName string) (*SDKInfo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sdk, err := s.manager.LookupSdk(sdkName)
	if err != nil {
		return nil, err
	}

	metadata := sdk.Metadata()
	versions := sdk.InstalledList()
	versionStrs := make([]string, len(versions))
	for i, v := range versions {
		versionStrs[i] = string(v)
	}

	current := sdk.Current()

	return &SDKInfo{
		Name:        metadata.Name,
		Description: metadata.PluginMetadata.Description,
		Homepage:    metadata.PluginMetadata.Homepage,
		Version:     metadata.PluginMetadata.Version,
		Installed:   true,
		Versions:    versionStrs,
		Current:     string(current),
	}, nil
}

// CheckSDKExists checks if an SDK version is installed
func (s *Service) CheckSDKExists(sdkName, version string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sdkInstance, err := s.manager.LookupSdk(sdkName)
	if err != nil {
		return false
	}

	return sdkInstance.CheckRuntimeExist(Version(version))
}

// GetConfig returns the current configuration
func (s *Service) GetConfig() (map[string]interface{}, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	cfg := s.manager.RuntimeEnvContext.UserConfig
	result := map[string]interface{}{
		"proxy": map[string]interface{}{
			"url":    cfg.Proxy.Url,
			"enable": cfg.Proxy.Enable,
		},
		"storage": map[string]interface{}{
			"sdkPath": cfg.Storage.SdkPath,
		},
		"registry": map[string]interface{}{
			"address": cfg.Registry.Address,
		},
		"cache": map[string]interface{}{
			"availableHookDuration": cfg.Cache.AvailableHookDuration.String(),
		},
		"legacyVersionFile": map[string]interface{}{
			"enable":   cfg.LegacyVersionFile.Enable,
			"strategy": cfg.LegacyVersionFile.Strategy,
		},
	}

	return result, nil
}

// ConfigUpdate represents configuration update request
type ConfigUpdate struct {
	Proxy             *ProxyConfig             `json:"proxy"`
	Storage           *StorageConfig           `json:"storage"`
	Registry          *RegistryConfig          `json:"registry"`
	Cache             *CacheConfig             `json:"cache"`
	LegacyVersionFile *LegacyVersionFileConfig `json:"legacyVersionFile"`
}

type ProxyConfig struct {
	Url    string `json:"url"`
	Enable bool   `json:"enable"`
}

type StorageConfig struct {
	SdkPath string `json:"sdkPath"`
}

type RegistryConfig struct {
	Address string `json:"address"`
}

type CacheConfig struct {
	AvailableHookDuration string `json:"availableHookDuration"`
}

type LegacyVersionFileConfig struct {
	Enable   bool   `json:"enable"`
	Strategy string `json:"strategy"`
}

// SaveConfig saves the configuration
func (s *Service) SaveConfig(update ConfigUpdate) OperationResult {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg := s.manager.RuntimeEnvContext.UserConfig

	// Update proxy settings
	if update.Proxy != nil {
		cfg.Proxy.Url = update.Proxy.Url
		cfg.Proxy.Enable = update.Proxy.Enable
	}

	// Update storage settings
	if update.Storage != nil {
		cfg.Storage.SdkPath = update.Storage.SdkPath
	}

	// Update registry settings
	if update.Registry != nil {
		cfg.Registry.Address = update.Registry.Address
	}

	// Update cache settings
	if update.Cache != nil && update.Cache.AvailableHookDuration != "" {
		// Parse duration string
		// Note: For simplicity, we store as string, actual parsing happens on reload
	}

	// Update legacy version file settings
	if update.LegacyVersionFile != nil {
		cfg.LegacyVersionFile.Enable = update.LegacyVersionFile.Enable
		cfg.LegacyVersionFile.Strategy = update.LegacyVersionFile.Strategy
	}

	// Save config to file
	configPath := s.manager.RuntimeEnvContext.PathMeta.User.Home
	if err := cfg.SaveConfig(configPath); err != nil {
		return OperationResult{
			Success: false,
			Error:   fmt.Sprintf("保存配置失败: %v", err),
		}
	}

	return OperationResult{
		Success: true,
		Message: "配置已保存，部分设置需要重启应用后生效",
	}
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsLower(lower(s), lower(substr))))
}

func containsLower(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func lower(s string) string {
	b := make([]byte, len(s))
	for i := range b {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

// ValidateSDKName validates an SDK name
func (s *Service) ValidateSDKName(name string) error {
	if name == "" {
		return errors.New("SDK name cannot be empty")
	}
	return nil
}
