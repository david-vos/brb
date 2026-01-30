package services

import (
	"browserRedirectBar/src/services"
	"testing"
)

func TestConfigService_LoadAndSave(t *testing.T) {

	// Create a config service - we'll need to test through NewConfigService
	// but that uses the default path. For now, test SetConfig/GetConfig

	service := &services.ConfigService{}

	// Test config
	testConfig := services.Config{
		Browsers: []services.BrowserConfig{
			{
				Patterns:      []string{"github.com"},
				RegexPatterns: []string{"^https://.*\\.internal\\..*"},
				BrowserURL:    "/Applications/Chrome.app",
			},
		},
		DefaultBrowserURL: "/Applications/Safari.app",
	}

	// Test SetConfig/GetConfig
	service.SetConfig(testConfig)
	loadedConfig := service.GetConfig()

	// Verify loaded config matches
	if len(loadedConfig.Browsers) != len(testConfig.Browsers) {
		t.Errorf("Loaded config has %d browsers, want %d", len(loadedConfig.Browsers), len(testConfig.Browsers))
	}

	if loadedConfig.DefaultBrowserURL != testConfig.DefaultBrowserURL {
		t.Errorf("DefaultBrowserURL = %q, want %q", loadedConfig.DefaultBrowserURL, testConfig.DefaultBrowserURL)
	}

	if len(loadedConfig.Browsers) > 0 {
		loadedBrowser := loadedConfig.Browsers[0]
		testBrowser := testConfig.Browsers[0]

		if !stringsEqual(loadedBrowser.Patterns, testBrowser.Patterns) {
			t.Errorf("Patterns = %v, want %v", loadedBrowser.Patterns, testBrowser.Patterns)
		}

		if !stringsEqual(loadedBrowser.RegexPatterns, testBrowser.RegexPatterns) {
			t.Errorf("RegexPatterns = %v, want %v", loadedBrowser.RegexPatterns, testBrowser.RegexPatterns)
		}

		if loadedBrowser.BrowserURL != testBrowser.BrowserURL {
			t.Errorf("BrowserURL = %q, want %q", loadedBrowser.BrowserURL, testBrowser.BrowserURL)
		}
	}

	// Note: Full Save/Load testing requires access to private configPath field
	// which is set by NewConfigService. We test SetConfig/GetConfig here.
}

func TestConfigService_Load_NoFile(t *testing.T) {
	// Test that Load() works when config file doesn't exist
	service := &services.ConfigService{}

	// Use a non-existent path
	err := service.Load()
	if err != nil {
		t.Errorf("Load() should not error on non-existent file, got: %v", err)
	}

	// Config should be empty when file doesn't exist
	config := service.GetConfig()
	if len(config.Browsers) != 0 {
		t.Errorf("Config should be empty when file doesn't exist, got %d browsers", len(config.Browsers))
	}
	if config.DefaultBrowserURL != "" {
		t.Errorf("DefaultBrowserURL should be empty when file doesn't exist, got %q", config.DefaultBrowserURL)
	}
}

// Helper function to compare string slices
func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
