package services

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// ConfigService handles configuration loading and saving
type ConfigService struct {
	config     Config
	configPath string
}

// NewConfigService creates a new ConfigService instance
func NewConfigService() (*ConfigService, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(homeDir, ".brb", "config.json")

	service := &ConfigService{
		configPath: configPath,
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, err
	}

	// Create empty config file if it doesn't exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		emptyConfig := Config{
			Browsers:          []BrowserConfig{},
			DefaultBrowserURL: "",
		}
		service.SetConfig(emptyConfig)
		if err := service.Save(); err != nil {
			return nil, err
		}
	}

	_ = service.Load() // Errors shown by menu on startup

	return service, nil
}

// Load loads the configuration from disk
// If the file doesn't exist, config remains empty (no error)
func (cs *ConfigService) Load() error {
	if cs.configPath == "" {
		// No path set, config stays empty
		return nil
	}

	data, err := os.ReadFile(cs.configPath)
	if err != nil {
		// File doesn't exist - that's okay, config stays empty
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	// If file is empty, use empty config
	if len(data) == 0 {
		cs.config = Config{
			Browsers:          []BrowserConfig{},
			DefaultBrowserURL: "",
		}
		return nil
	}

	if err := json.Unmarshal(data, &cs.config); err != nil {
		log.Printf("Invalid config at %s: %v", cs.configPath, err)
		cs.config = Config{
			Browsers:          []BrowserConfig{},
			DefaultBrowserURL: "",
		}
		return fmt.Errorf("invalid JSON in config file: %w", err)
	}

	return nil
}

// Save saves the configuration to disk
func (cs *ConfigService) Save() error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(cs.configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cs.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(cs.configPath, data, 0644)
}

// GetConfig returns the current configuration
func (cs *ConfigService) GetConfig() Config {
	return cs.config
}

// SetConfig sets the configuration
func (cs *ConfigService) SetConfig(config Config) {
	cs.config = config
}

// GetConfigPath returns the config file path
func (cs *ConfigService) GetConfigPath() string {
	return cs.configPath
}

// SetDefaultBrowser sets the default browser and saves the configuration
// It reloads the config first to ensure we don't lose any existing browser configurations
func (cs *ConfigService) SetDefaultBrowser(browserPath string) error {
	if err := cs.Load(); err != nil {
		return fmt.Errorf("cannot load config file: %w", err)
	}
	cs.config.DefaultBrowserURL = browserPath
	if err := cs.Save(); err != nil {
		return fmt.Errorf("cannot save config: %w", err)
	}
	return nil
}
